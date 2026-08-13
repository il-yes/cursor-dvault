package c3_integration_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	tracecore_types "vault-app/internal/tracecore/types"
)

func TestC3SecurityBoundary_NoSecretsInPayloadRef(t *testing.T) {
	forbiddenSubstrings := []string{
		"DEK",
		"KEK",
		"PrivateKey",
		"SecretKey",
		"Seed",
		"Password",
		"private_key",
		"secret_key",
		"encrypted_payload",
		"decrypted_content",
	}

	payloadRef := tracecore_types.PayloadRefDTO{
		CID:         "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
		ContentHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Size:        1048576,
		AssetType:   "contract_draft",
		Name:        "Master_Services_Agreement_v2.pdf",
	}

	eventDTO := tracecore_types.ThreadEventDTO{
		ID:       "evt_123456789",
		ThreadID: "thread_987654321",
		Type:     "entry.shared",
		Cursor:   1,
		Payload: map[string]any{
			"notes":       "Countersigned contract draft committed to vault.",
			"payload_ref": payloadRef,
		},
		CreatedAt: time.Now(),
	}

	marshaledBytes, err := json.Marshal(eventDTO)
	if err != nil {
		t.Fatalf("Failed to marshal ThreadEventDTO: %v", err)
	}

	marshaledStr := string(marshaledBytes)

	for _, forbidden := range forbiddenSubstrings {
		if strings.Contains(marshaledStr, forbidden) {
			t.Errorf("SECURITY VIOLATION: Serialized ThreadEventDTO contains forbidden key/secret substring: %q", forbidden)
		}
	}
}

func TestC3SecurityBoundary_NoSecretsInShareEntryRef(t *testing.T) {
	forbiddenSubstrings := []string{
		"wrapped_dek",
		"wrappedDEK",
		"kek_version",
		"DEK",
		"KEK",
		"PrivateKey",
		"SecretKey",
		"Seed",
		"Password",
		"private_key",
		"secret_key",
		"encrypted_payload",
		"decrypted_content",
	}

	shareEntryRef := tracecore_types.ShareEntryRefDTO{
		ShareEntryID: "se_9988776655",
		TrustGroupID: "tg_legal_counsel",
		AssetCID:     "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
		CreatedBy:    "user_vault_legal",
		Status:       "active",
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	eventDTO := tracecore_types.ThreadEventDTO{
		ID:       "evt_share_001",
		ThreadID: "thread_contract_01",
		Type:     "entry.shared",
		Cursor:   2,
		Payload: map[string]any{
			"notes":           "Collaborative share entry created for legal review.",
			"share_entry_ref": shareEntryRef,
		},
		CreatedAt: time.Now(),
	}

	marshaledBytes, err := json.Marshal(eventDTO)
	if err != nil {
		t.Fatalf("Failed to marshal ThreadEventDTO: %v", err)
	}

	marshaledStr := string(marshaledBytes)

	for _, forbidden := range forbiddenSubstrings {
		if strings.Contains(marshaledStr, forbidden) {
			t.Errorf("SECURITY VIOLATION: Serialized C3 ShareEntryRef contains forbidden secret/key substring: %q", forbidden)
		}
	}
}
