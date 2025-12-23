package registry

import (
	"encoding/json"
	"fmt"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	vault_session "vault-app/internal/vault/application/session"
)

type EntryHandler interface {
	Add(userID string, entry any) (*any, error)
	Edit(userID string, entry any) (*any, error)
	Trash(userID string, entryID string) error
	Restore(userID string, entryID string) error
	SetVault(vault *vault_session.Session)
}

type EntryRegistry struct {
	logger logger.Logger
	handlers map[string]EntryHandler
	Vault *vault_session.Session
}


type EntryDefinition struct {
	Type     string
	Factory  func() models.VaultEntry
	Handler  EntryHandler
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
	return &EntryRegistry{
		logger: *logger,
		handlers: make(map[string]EntryHandler),
	}
}

func (r *EntryRegistry) HandlerFor(entryType string) (EntryHandler, error) {
	h, ok := r.handlers[entryType]
	if !ok {
		return nil, fmt.Errorf("handler not found for type: %s", entryType)
	}
	return h, nil
}

var entryFactories = map[string]func() models.VaultEntry{}

func(r *EntryRegistry) UnmarshalEntry(entryType string, raw []byte) (models.VaultEntry, error) {
    factory, ok := entryFactories[entryType]
    if !ok {
        return nil, fmt.Errorf("unknown entry type: %s", entryType)
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
func (r *EntryRegistry) RegisterEntryType(name string, factory func() models.VaultEntry) {
	entryFactories[name] = factory
}


