package services

import (
	"fmt"
	"vault-app/internal/models"
	"vault-app/internal/tracecore"
)

const (
	CREATE_ENTRY = "create_entry"
	SHARE_ENTRY = "share_entry"
	ACCESS_ENTRY = "access_entry"
)

type CommitPayloadFactory struct {
	Action string
	From string						 		// public key of owner
	Entry     models.VaultEntry   // unique entry reference
	TargetUser  string   // optional: public key user receiving share/revoke
	Permissions []string // optional: for read/share/revoke
	Expiry      string
}

func NewCommitPayloadFactory(action string, from string, entry models.VaultEntry, targetUser string, perm []string, expiry string) *CommitPayloadFactory {

	return &CommitPayloadFactory{
		Action: action,
		From: from,
		Entry: entry,
		TargetUser: targetUser,	
		Permissions: perm,		
		Expiry: expiry,
	}
}

func (c *CommitPayloadFactory) BuildPostEntryPayload() tracecore.CommitMetadata {
	var cp tracecore.CommitMetadata

	cp.Message = fmt.Sprintf("Create entry: %s", c.Entry.GetTypeName()) //"Create entry: car_unlock_token"
	cp.Content = map[string]any{
		"entry_id":   fmt.Sprintf("%s:entry: %s", c.From, c.Entry.GetName()), // "alice:entry:car_unlock_001",
		"entry_type": c.Entry.GetTypeName(),                                // "credential",
		"entry_name": c.Entry.GetName(),                                    // "car_unlock_token",
	}
	cp.Context= map[string]string{
		"phase": "vault_entry",
		"stage": "creation",
	}
	cp.StatusChange = tracecore.StatusChange{
		Old: "none",
		New: "created",
	}

	return cp
}
func (c *CommitPayloadFactory) BuildShareEntryPayload() tracecore.CommitMetadata {
	var cp tracecore.CommitMetadata

	cp.Message = fmt.Sprintf("Share entry: %s  with %s", c.Entry.GetTypeName(), c.TargetUser) //"Share entry: car_unlock_token with john_pub"
	cp.Content = map[string]any{
		"entry_id":    fmt.Sprintf("%s:share_entry: %s", c.From, c.Entry.GetName()), // "alice:entry:car_unlock_001",
		"target_user": c.TargetUser,                                                       // "credential",
		"permissions": c.Permissions,
		"expiry":      c.Expiry, // "2025-12-31T23:59:59Z",
	}
	cp.Context = map[string]string{
		"phase": "vault_entry",
		"stage": "creation",
	}
	cp.StatusChange = tracecore.StatusChange{
		Old: "none",
		New: "created",
	}

	return cp
}
func (c *CommitPayloadFactory) BuildAccessEntryPayload() tracecore.CommitMetadata {
	var cp tracecore.CommitMetadata

	cp.Message = fmt.Sprintf("Access entry: %s by %s", c.Entry.GetTypeName(), c.TargetUser) //  "Access entry: car_unlock_token by john_pub"
	cp.Content = map[string]any{
		"entry_id":     fmt.Sprintf("%s:share_entry: %s", c.From, c.Entry.GetName()), // "alice:shared_entry:car_unlock_001",
		"target_user":  c.TargetUser,                                                              // "credential",
		"shared_entry": c.Entry.GetId(),
	}
	cp.Context = map[string]string{
		"phase":    "vault_entry",
		"stage":    "access",
		"substage": "owner_offline",
	}
	cp.StatusChange = tracecore.StatusChange{
		Old: "idle",
		New: "accessed",
	}

	return cp
}
func  (c *CommitPayloadFactory) Build() (*tracecore.CommitMetadata, error) {
	switch c.Action {
	case CREATE_ENTRY:
		fmt.Println("%s action", CREATE_ENTRY)
		response := c.BuildPostEntryPayload()
		return &response, nil
	case SHARE_ENTRY:
		fmt.Println("%s action", SHARE_ENTRY)
		response := c.BuildShareEntryPayload()
		return &response, nil
	case ACCESS_ENTRY:
		fmt.Println("%s action", ACCESS_ENTRY)
		response := c.BuildAccessEntryPayload()
		return &response, nil
	}
	return nil, fmt.Errorf("failed to execute %s: unknown action", c.Action)
}