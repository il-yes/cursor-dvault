package app_config_tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	app_config_worker "vault-app/internal/config/application/worker"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/logger/logger"
	vault_dto "vault-app/internal/vault/application/dto"
	vaults_domain "vault-app/internal/vault/domain"
)

// ------------------------------------------------------------------------------------------------------------
// MOCKS
// ------------------------------------------------------------------------------------------------------------
// MockRegistryClient ===================================================================================
type MockRegistryClient struct {
	mock.Mock
}

func (m *MockRegistryClient) GetPack(ctx context.Context, packID string) (*app_config_worker.PackDTO, error) {
	args := m.Called(ctx, packID)
	pack, _ := args.Get(0).(*app_config_worker.PackDTO)
	return pack, args.Error(1)
}

func (m *MockRegistryClient) GetTemplate(ctx context.Context, templateID string) (*app_config_worker.TemplateDTO, error) {
	args := m.Called(ctx, templateID)
	tpl, _ := args.Get(0).(*app_config_worker.TemplateDTO)
	return tpl, args.Error(1)
}

// MockEntryHandler =====================================================================================
type MockEntryHandler struct {
	mock.Mock
}

func (m *MockEntryHandler) AddEntryFor(userID string, entry any) (*vaults_domain.VaultEntry, error) {
	args := m.Called(userID, entry)
	ve, _ := args.Get(0).(*vaults_domain.VaultEntry)
	return ve, args.Error(1)
}
func (m *MockEntryHandler) VaultPayloadAddEntry(req vault_dto.VaultPayloadAddEntryRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

// MockOnboardingConfigGetter =============================================================================
type MockOnboardingConfigGetter struct {
	mock.Mock
}

func (m *MockOnboardingConfigGetter) GetOnboardingConfigByUserID(ctx context.Context, userID string, opts map[string]interface{}) (*app_config_domain.OnboardingConfig, error) {
	args := m.Called(ctx, userID, opts)
	cfg, _ := args.Get(0).(*app_config_domain.OnboardingConfig)
	return cfg, args.Error(1)
}

func (m *MockOnboardingConfigGetter) UpdateOnboardingConfig(userID string, cfg *app_config_domain.OnboardingConfig, opts map[string]interface{}) error {
	args := m.Called(userID, cfg, opts)
	return args.Error(0)
}

// ============================================================================================================
// TESTS
// ============================================================================================================
func TestApplyOnboardingPacksWorker_HandleOnboardingCompleted(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	assert.NoError(t, err)
	logger := logger.NewFromEnv()
	userOnboardingID := "user_onboard_123"

	t.Run("returns immediately when no packs are configured", func(t *testing.T) {
		registry := &MockRegistryClient{}
		entries := &MockEntryHandler{}
		cfgGetter := &MockOnboardingConfigGetter{}

		cfgGetter.
			On("GetOnboardingConfigByUserID", mock.Anything, "user-123", (map[string]interface{})(nil)).
			Return(&app_config_domain.OnboardingConfig{
				UserID:         "user-123",
				Packs:          nil,
				InstalledSeeds: nil,
				PacksApplied:   false,
				Completed:      false,
			}, nil)

		worker := app_config_worker.NewApplyOnboardingPacksWorker(
			db,
			registry,
			cfgGetter,
			entries,
			time.Now,
			logger,
		)

		err := worker.HandleOnboardingCompleted(context.Background(), "user-123", "vault-1", userOnboardingID)
		assert.NoError(t, err)

		cfgGetter.AssertExpectations(t)
		registry.AssertExpectations(t)
		entries.AssertExpectations(t)
	})

	t.Run("creates entries and updates onboarding config when a new pack seed is found", func(t *testing.T) {
		registry := &MockRegistryClient{}
		entries := &MockEntryHandler{}
		cfgGetter := &MockOnboardingConfigGetter{}

		cfgGetter.
			On("GetOnboardingConfigByUserID", mock.Anything, "user-123", (map[string]interface{})(nil)).
			Return(&app_config_domain.OnboardingConfig{
				UserID:         "user-123",
				Packs:          []string{"pack-1"},
				InstalledSeeds: []string{},
				PacksApplied:   false,
				Completed:      false,
			}, nil)

		cfgGetter.
			On("GetOnboardingConfigByUserID", mock.Anything, "user-123",
				mock.MatchedBy(func(opts map[string]interface{}) bool {
					_, ok := opts["db_transaction"]
					return ok
				}),
			).
			Return(&app_config_domain.OnboardingConfig{
				UserID:         "user-123",
				Packs:          []string{"pack-1"},
				InstalledSeeds: []string{},
				PacksApplied:   false,
				Completed:      false,
			}, nil)

		registry.
			On("GetPack", mock.Anything, "pack-1").
			Return(&app_config_worker.PackDTO{
				ID: "pack-1",
				Templates: []app_config_worker.TemplateRef{
					{
						SeedID:     "seed-1",
						TemplateID: "tmpl-1",
						FolderID:   "Incidents",
						Overrides:  map[string]any{"status": "draft"},
					},
				},
				},  nil)

		registry.
			On("GetTemplate", mock.Anything, "tmpl-1").
			Return(&app_config_worker.TemplateDTO{
				TemplateID:    "tmpl-1",
				RecordType:    "incident",
				SchemaVersion: 1,
				Fields: map[string]any{
					"status":   "detected",
					"severity": "medium",
				},
			}, nil)

		entry := vaults_domain.NoteEntry{
			BaseEntry: vaults_domain.BaseEntry{
				ID:   "entry-1",
				Type: "note",
			},
		}
		entries.
			On("AddEntryFor", "user-123", mock.Anything).
			Return(&entry, nil)

		cfgGetter.
			On("Update", "user-123", mock.MatchedBy(func(cfg *app_config_domain.OnboardingConfig) bool {
				return cfg.UserID == "user-123" &&
					len(cfg.Packs) == 1 &&
					cfg.Packs[0] == "pack-1" &&
					len(cfg.InstalledSeeds) == 1 &&
					cfg.InstalledSeeds[0] == "pack-1:seed-1" &&
					cfg.PacksApplied
			}), mock.Anything).
			Return(nil)

		worker := app_config_worker.NewApplyOnboardingPacksWorker(
			db,
			registry,
			cfgGetter,
			entries,
			time.Now,
			logger,
		)

		err := worker.HandleOnboardingCompleted(context.Background(), "user-123", "vault-1", userOnboardingID)
		assert.NoError(t, err)

		cfgGetter.AssertExpectations(t)
		registry.AssertExpectations(t)
		entries.AssertExpectations(t)
	})

	t.Run("returns error when entry creation fails", func(t *testing.T) {
		registry := &MockRegistryClient{}
		entries := &MockEntryHandler{}
		cfgGetter := &MockOnboardingConfigGetter{}

		cfgGetter.
			On("GetOnboardingConfigByUserID", mock.Anything, "user-123", (map[string]interface{})(nil)).
			Return(&app_config_domain.OnboardingConfig{
				UserID:         "user-123",
				Packs:          []string{"pack-1"},
				InstalledSeeds: []string{},
				PacksApplied:   false,
				Completed:      false,
			}, nil)

		cfgGetter.
			On("GetOnboardingConfigByUserID", mock.Anything, "user-123",
				mock.MatchedBy(func(opts map[string]interface{}) bool {
					_, ok := opts["db_transaction"]
					return ok
				}),
			).
			Return(&app_config_domain.OnboardingConfig{
				UserID:         "user-123",
				Packs:          []string{"pack-1"},
				InstalledSeeds: []string{},
				PacksApplied:   false,
				Completed:      false,
			}, nil)

		registry.
			On("GetPack", mock.Anything, "pack-1").
			Return(&app_config_worker.PackDTO{
				ID: "pack-1",
				Templates: []app_config_worker.TemplateRef{
					{
						SeedID:     "seed-1",
						TemplateID: "tmpl-1",
						FolderID:   "Incidents",
						Overrides:  nil,
					},
				},
			}, nil)

		registry.
			On("GetTemplate", mock.Anything, "tmpl-1").
			Return(&app_config_worker.TemplateDTO{
				TemplateID:    "tmpl-1",
				RecordType:    "incident",
				SchemaVersion: 1,
				Fields: map[string]any{
					"severity": "medium",
				},
			}, nil)

		entries.
			On("AddEntryFor", "user-123", mock.Anything).
			Return((*vaults_domain.VaultEntry)(nil), errors.New("failed to create entry"))

		worker := app_config_worker.NewApplyOnboardingPacksWorker(
			db,
			registry,
			cfgGetter,
			entries,
			time.Now,
			logger,
		)

		err := worker.HandleOnboardingCompleted(context.Background(), "user-123", "vault-1", userOnboardingID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create seed entry")

		cfgGetter.AssertExpectations(t)
		registry.AssertExpectations(t)
		entries.AssertExpectations(t)
		cfgGetter.AssertNotCalled(t, "UpdateOnboardingConfig", mock.Anything, mock.Anything, mock.Anything)
	})
}
