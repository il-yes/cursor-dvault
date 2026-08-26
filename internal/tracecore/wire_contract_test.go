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
