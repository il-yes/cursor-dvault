package app_config_worker

import (
	"context"
	"fmt"
	"time"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/logger/logger"
	"vault-app/internal/utils"
	vault_dto "vault-app/internal/vault/application/dto"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Pseudocode (high level)
// Load config with a DB lock (or transaction + row lock)
// If cfg.Onboarding.Completed == true and packs already applied → return success (idempotent)
// Fetch packs from cloud registry (by IDs)
// For each pack’s template refs:
// if template already installed (in config) → skip
// else fetch template definition from registry
// compute seeded fields = deepMerge(template.fields, seed.overrides)
// create encrypted blob entry (draft)
// mark template as installed
// Update config:
// append Packs (unique)
// append InstalledTemplates (unique)
// set Completed = true (or set PacksApplied = true if you want a separate flag)
// Commit transaction

type SeedWorkItem struct {
	PackID     string
	SeedID     string
	TemplateID string
	FolderID   string
	Overrides  map[string]any
}

type TemplateDTO struct {
	TemplateID    string         `json:"template_id"`
	RecordType    string         `json:"record_type"`
	SchemaVersion int            `json:"schema_version"`
	Fields        map[string]any `json:"fields"`
}
type TemplateRef struct {
	SeedID     string         `json:"seed_id"`
	TemplateID string         `json:"template_id"`
	FolderID   string         `json:"folder_id,omitempty"`
	Overrides  map[string]any `json:"overrides,omitempty"`
}
type PackDTO struct {
	ID        string
	Templates []TemplateRef `json:"templates"`
}

type CreateEntryInput struct {
	UserID    string	`json:"user_id"`
	VaultName string	`json:"vault_name"`
	FolderID  string	`json:"folder_id"`

	// Merged fields from template.fields + seed.overrides
	SeedFields map[string]any	`json:"seed_fields"`

	// Info to reconstruct where this seed came from
	TemplateID    string	`json:"template_id"`
	SchemaVersion int	`json:"schema_version"`
	RecordType    string	`json:"record_type"`
	PackID        string	`json:"pack_id"`
	SeedID        string	`json:"seed_id"`
}

// ====================================================================================================
// Interfaces
// ====================================================================================================
type RegistryClient interface {
	GetPack(ctx context.Context, packID string) (*PackDTO, error)
	GetTemplate(ctx context.Context, templateID string) (*TemplateDTO, error)
}

type EntryHandlerInterface interface {
	AddEntryFor(userID string, entry any) (*vaults_domain.VaultEntry, error)
	VaultPayloadAddEntry(req vault_dto.VaultPayloadAddEntryRequest) error 
}

type OnboardingConfigGetter interface {
	GetOnboardingConfigByUserID(ctx context.Context, userID string, opts map[string]interface{}) (*app_config_domain.OnboardingConfig, error)
	UpdateOnboardingConfig(userId string, cfg *app_config_domain.OnboardingConfig, opts map[string]interface{}) error
}

// ====================================================================================================
// Worker
// ====================================================================================================

type ApplyOnboardingPacksWorker struct {
	db            *gorm.DB
	registry      RegistryClient
	configHandler OnboardingConfigGetter
	entryHandler  EntryHandlerInterface
	clock         func() time.Time
	logger        *logger.Logger
}

func NewApplyOnboardingPacksWorker(
	db *gorm.DB,
	registry RegistryClient,
	configHandler OnboardingConfigGetter,
	entryHandler EntryHandlerInterface,
	clock func() time.Time,
	logger *logger.Logger,
) *ApplyOnboardingPacksWorker {
	return &ApplyOnboardingPacksWorker{
		db:            db,
		registry:      registry,
		configHandler: configHandler,
		entryHandler:  entryHandler,
		clock:         clock,
		logger:        logger,
	}
}

func (w *ApplyOnboardingPacksWorker) HandleOnboardingCompleted(
	ctx context.Context,
	userID,
	vaultName string,
	userOnboardingID string,
) error {
	w.logger.Info("Installing onboarding packs for user %s and vault %s", userID, vaultName)

	// 1. Load config WITHOUT transaction first (no FOR UPDATE)
	cfgPreview, err := w.loadConfig(ctx, userOnboardingID)
	if err != nil {
		w.logger.Error("Failed to load config for user %s and vault %s: %v", userID, vaultName, err)
		return err
	}
	w.logger.Info("cfgPreview %v for user %s", cfgPreview, userOnboardingID)

	selectedPackIDs := cfgPreview.Packs
	if len(selectedPackIDs) == 0 {
		w.logger.Warn("No packs selected for user %s and vault %s", userID, vaultName)
		return nil
	}
	w.logger.Info("Selected packs for user %s and vault %s: %v", userID, vaultName, selectedPackIDs)

	// 2. Prefetch seeds from packs (outside DB txn)
	w.logger.Info("Prefetching seeds for user %s and vault %s", userID, vaultName)
	seeds, err := w.prefetchSeeds(ctx, selectedPackIDs)
	if err != nil {
		w.logger.Error("Failed to prefetch seeds for user %s and vault %s: %v", userID, vaultName, err)
		return err
	}

	if len(seeds) == 0 {
		w.logger.Warn("No seeds found for user %s and vault %s", userID, vaultName)
		return nil
	}

	// 3. Prefetch templates (outside DB txn)
	w.logger.Info("Prefetching templates for user %s and vault %s", userID, vaultName)
	templatesByID, err := w.prefetchTemplates(ctx, seeds)
	if err != nil {
		w.logger.Error("Failed to prefetch templates for user %s and vault %s: %v", userID, vaultName, err)
		return err
	}

	// 4. Real DB transaction: lock config + install missed seeds
	w.logger.Info("Installing onboarding packs for user %s and vault %s", userID, vaultName)
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 4a. SELECT FOR UPDATE on onboarding config
		cfg, err := w.configHandler.GetOnboardingConfigByUserID(ctx, userOnboardingID, map[string]interface{}{"db_transaction": tx})
		if err != nil {
			return err
		}

		// 4b. Convert cfg.InstalledSeeds → SeedKey set
		installedKeys := make(map[string]bool, len(cfg.InstalledSeeds))
		for _, kStr := range cfg.InstalledSeeds {
			installedKeys[kStr] = true
		}

		// 4c. For each seed in selected packs
		for _, item := range seeds {
			seedKey := app_config_domain.NewSeedKey(item.PackID, item.SeedID).String()

			if installedKeys[seedKey] {
				w.logger.Info("Seed already installed, skipping",
					map[string]any{
						"user_id":     userID,
						"pack_id":     item.PackID,
						"seed_id":     item.SeedID,
						"template_id": item.TemplateID,
					})
				continue
			}

			// 4d. Fetch template definition
			tpl, ok := templatesByID[item.TemplateID]
			if !ok {
				return fmt.Errorf("template not found: %s", item.TemplateID)
			}
			utils.LogPretty("ApplyOnboardingPacksWorker - HandleOnboardingCompleted - tpl", tpl)

			// 4e. Deep‑merge template.fields + seed.overrides
			seededFields := DeepMergeCopy(tpl.Fields, item.Overrides)

			// 4f. Create encrypted entry
			err := w.CreateEncryptedEntry(ctx, CreateEntryInput{
				UserID:     userID,
				VaultName:  vaultName,
				FolderID:   item.FolderID,
				SeedFields: seededFields,
				TemplateID: tpl.TemplateID,
				RecordType: tpl.RecordType,
				PackID:     item.PackID,
				SchemaVersion: tpl.SchemaVersion,
				SeedID:     item.SeedID,
			})
			if err != nil {
				return fmt.Errorf("failed to create seed entry for %s/%s: %w", item.PackID, item.SeedID, err)
			}

			// 4g. Mark seed as installed
			installedKeys[seedKey] = true
		}
		

		// 4h. Rebuild InstalledSeeds slice from set (optional: keep in DB as JSON list of strings)
		var newKeys []string
		for k := range installedKeys {
			newKeys = append(newKeys, k)
		}
		newKeys = uniqueSeedKeysFromStrings(newKeys) // trivial if you already deduped
		cfg.InstalledSeeds = newKeys

		// 4i. Optionally mark pack application as done
		// This is cosmetic; true idempotency is on SeedKey
		cfg.PacksApplied = true

		// 4j. Save onboarding config
		w.logger.Info("cfg %v for user %s", cfg, userOnboardingID)
		return w.configHandler.UpdateOnboardingConfig(cfg.UserID, cfg, map[string]interface{}{"db_transaction": tx})
	})
}

func DeepMergeCopy(base map[string]any, overrides map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range base {
		out[k] = v
	}
	if overrides == nil {
		return out
	}
	return deepMerge(out, overrides)
}

func deepMerge(dst, src map[string]any) map[string]any {
	// dst is mutated
	for k, v := range src {
		vm, ok := v.(map[string]any)
		if !ok {
			dst[k] = v
			continue
		}
		if existing, ok := dst[k].(map[string]any); ok {
			dst[k] = deepMerge(existing, vm)
		} else {
			dst[k] = deepMerge(map[string]any{}, vm)
		}
	}
	return dst
}

func uniqueSeedKeysFromStrings(s []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func (w *ApplyOnboardingPacksWorker) prefetchSeeds(
	ctx context.Context,
	packIDs []string,
) ([]*SeedWorkItem, error) {
	var seeds []*SeedWorkItem
	for _, packID := range packIDs {
		pack, err := w.registry.GetPack(ctx, packID)
		if err != nil {
			return nil, err
		}
		for _, templateRef := range pack.Templates {
			seed := &SeedWorkItem{
				PackID:     packID,
				SeedID:     templateRef.TemplateID,
				TemplateID: templateRef.TemplateID,
				FolderID:   templateRef.FolderID,
				Overrides:  templateRef.Overrides,
			}
			seeds = append(seeds, seed)
		}
	}
	return seeds, nil
}

func (w *ApplyOnboardingPacksWorker) prefetchTemplates(
	ctx context.Context,
	seeds []*SeedWorkItem,
) (map[string]*TemplateDTO, error) {
	templates := make(map[string]*TemplateDTO)
	for _, seed := range seeds {
		if _, ok := templates[seed.TemplateID]; ok {
			continue // dedup
		}
		tpl, err := w.registry.GetTemplate(ctx, seed.TemplateID)
		if err != nil {
			return nil, err
		}
		templates[seed.TemplateID] = tpl
	}
	return templates, nil
}

func (w *ApplyOnboardingPacksWorker) loadConfig(
	ctx context.Context,
	userOnboardingID string,
) (*app_config_domain.OnboardingConfig, error) {
	utils.LogPretty("ApplyOnboardingPacksWorker.loadConfig", userOnboardingID)
	var opt map[string]interface{}
	cfg, err := w.configHandler.GetOnboardingConfigByUserID(ctx, userOnboardingID, opt)
	if err != nil {
		// Create a minimal default (no packs, no installed seeds, not applied)
		return &app_config_domain.OnboardingConfig{
			UserID:         userOnboardingID,
			Packs:          []string{},
			InstalledSeeds: []string{},
			UseCases:       []string{},
			PacksApplied:   false,
			Completed:      false,
		}, nil
	}
	return cfg, nil
}

// Stored in go runtime during the session
func (w *ApplyOnboardingPacksWorker) CreateEncryptedEntry(ctx context.Context, in CreateEntryInput) error {
	typeE := "card"
	seedEntry := vaults_domain.CardEntry{
		BaseEntry: vaults_domain.BaseEntry{
			ID:              uuid.NewString(),
			CID:             "",
			EntryName:       "Specimen_" + uuid.NewString(),
			FolderID:        in.FolderID,
			Type:            "card",
			TemplateID:      in.TemplateID,
			SchemaVersion:   in.SchemaVersion,
			RecordType:      in.RecordType,
			AdditionnalNote: in.TemplateID,
			CustomFields:    in.SeedFields,
			Trashed:         false,
			IsDraft:         false,
			IsDirty:         false,
			IsFavorite:      false,
			Attachments:     []vaults_domain.Attachment{},
			AttachmentCIDs:  []string{},
		},
	}
	if err := w.entryHandler.VaultPayloadAddEntry(vault_dto.VaultPayloadAddEntryRequest{
		UserID:     in.UserID,
		EntryType: typeE,
		Entry:      seedEntry,		
	}); err != nil {
		return err
	}
	utils.LogPretty("ApplyOnboardingPacksWorker.CreateEncryptedEntry", seedEntry)

	return nil
}
