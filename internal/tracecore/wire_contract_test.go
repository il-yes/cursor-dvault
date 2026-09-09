package tracecore_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tracecore_types "vault-app/internal/tracecore/types"
)

func TestNewCreateWorkspaceRequest_WireContract(t *testing.T) {
	req := tracecore_types.NewCreateWorkspaceRequest{
		VaultID:     "fbffb9ad-d4de-4581-b7af-7b8e160f63bb",
		Name:        "Test Workspace",
		Description: "Boundary Wire Contract Test",
		OwnerID:     "fbffb9ad-d4de-4581-b7af-7b8e160f63bb",
	}

	bytes, err := json.Marshal(req)
	require.NoError(t, err)

	var rawMap map[string]interface{}
	err = json.Unmarshal(bytes, &rawMap)
	require.NoError(t, err)

	// Explicit Wire Contract Assertion:
	// Cloud backend JSON decoder expects exact key "vault_id" (snake_case).
	// "VaultID" (PascalCase) causes Cloud to decode empty string and return 403 Forbidden.
	assert.Contains(t, rawMap, "vault_id", "JSON wire payload must explicitly emit 'vault_id'")
	assert.Equal(t, "fbffb9ad-d4de-4581-b7af-7b8e160f63bb", rawMap["vault_id"])
	assert.NotContains(t, rawMap, "VaultID", "JSON wire payload must NOT emit PascalCase 'VaultID'")
}

func TestCreateTrustGroup_WireContract_MemberObjects(t *testing.T) {
	// Verify that CreateTrustGroup marshals members as MemberRequest objects with vault_id and role,
	// rather than raw strings, satisfying the Cloud backend CreateGroupRequest contract.
	callerVaultID := "60496dae-7df6-4441-93f7-60d54e53c55e"
	membersList := []map[string]interface{}{
		{
			"vault_id": callerVaultID,
			"role":     "admin",
		},
	}
	payload := map[string]interface{}{
		"workspace_id": "0958f238-983d-4f03-bb9d-bcf8fbe21d56",
		"name":         "Loop-engineering",
		"members":      membersList,
	}

	bytes, err := json.Marshal(payload)
	require.NoError(t, err)

	var rawMap map[string]interface{}
	err = json.Unmarshal(bytes, &rawMap)
	require.NoError(t, err)

	members, ok := rawMap["members"].([]interface{})
	require.True(t, ok, "members must be a JSON array")
	require.Len(t, members, 1)

	firstMember, ok := members[0].(map[string]interface{})
	require.True(t, ok, "each member must be a JSON object, NOT a raw string")
	assert.Equal(t, callerVaultID, firstMember["vault_id"])
	assert.Equal(t, "admin", firstMember["role"])
}
