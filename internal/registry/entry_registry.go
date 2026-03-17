package registry

import (
	"encoding/json"
	"fmt"
	"vault-app/internal/logger/logger"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

// Catalog
type EntryHandler interface {
	Add(userID string, entry any) (*vaults_domain.VaultPayload, error)
	Edit(userID string, entry any) (*vaults_domain.VaultPayload, error)
	Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error)
	Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error)
	SetSession(session *vault_session.Session)
	SetVaultRepository(vaultRepository vaults_domain.VaultRepository)
	EditWithAttachments(userID string, entry any, attachments []vault_dto.SelectedAttachment) (*vaults_domain.VaultPayload, error)
}

type EntryRegistry struct {
	logger   logger.Logger
	handlers map[string]EntryHandler
	Vault    *vault_session.Session
	Session  *vault_session.Session
}

type EntryDefinition struct {
	Type    string
	Factory func() vaults_domain.VaultEntry
	Handler EntryHandler
}

func (r *EntryRegistry) RegisterDefinitions(defs []EntryDefinition) {
	for _, def := range defs {
		if def.Factory != nil {
			entryFactories[def.Type] = def.Factory
		}
		if def.Handler != nil {
			r.Register(def.Type, def.Handler)
		}
	}
}
func NewRegistry(logger *logger.Logger) *EntryRegistry {
	logger.Info("🔧 Registry - Initializing Registry...")
	return &EntryRegistry{
		logger:   *logger,
		handlers: make(map[string]EntryHandler),
	}
}

func (r *EntryRegistry) HandlerFor(entryType string) (EntryHandler, error) {
	h, ok := r.handlers[entryType]
	if !ok {
		return nil, fmt.Errorf("🔧 Registry - handler not found for type: %s", entryType)
	}
	return h, nil
}

var entryFactories = map[string]func() vaults_domain.VaultEntry{}

func (r *EntryRegistry) UnmarshalEntry(entryType string, raw []byte) (vaults_domain.VaultEntry, error) {
	factory, ok := entryFactories[entryType]
	if !ok {
		return nil, fmt.Errorf("🔧 Registry - unknown entry type: %s", entryType)
	}
	entry := factory()
	if err := json.Unmarshal(raw, entry); err != nil {
		return nil, err
	}
	return entry, nil
}
func (r *EntryRegistry) Register(entryType string, handler EntryHandler) {
	r.handlers[entryType] = handler
}
func (r *EntryRegistry) RegisterEntryType(name string, factory func() vaults_domain.VaultEntry) {
	entryFactories[name] = factory
}
