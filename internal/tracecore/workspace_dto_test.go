package tracecore_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
)

func TestCloudWorkspaceDTOMapping_FullPayload(t *testing.T) {
	cloudJSON := `{
		"status": 200,
		"data": [
			{
				"ID": "ws_12345",
				"VaultID": "v_67890",
				"Name": "Engineering Workspace",
				"Description": "Primary workspace for engineering",
				"Status": "active",
				"OwnerID": "usr_70",
				"CreatedAt": "2026-08-18T14:00:00Z",
				"UpdatedAt": "2026-08-18T14:30:00Z",
				"IsDraft": false,
				"IsDirty": true
			}
		],
		"message": "success",
		"success": true
	}`

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.CloudWorkspaceDTO]
	err := json.Unmarshal([]byte(cloudJSON), &cloudResp)
	require.NoError(t, err)
	require.Len(t, cloudResp.Data, 1)

	dto := cloudResp.Data[0]
	mapped := tracecore_types.MapCloudWorkspaceToTypes(&dto)
	require.NotNil(t, mapped)

	// Explicit invariant check: VaultID must NEVER disappear
	assert.NotEmpty(t, mapped.VaultID, "VaultID must never silently disappear")
	assert.Equal(t, "v_67890", mapped.VaultID)
	assert.Equal(t, "ws_12345", mapped.ID)
	assert.Equal(t, "Engineering Workspace", mapped.Name)
	assert.Equal(t, "Primary workspace for engineering", mapped.Description)
	assert.Equal(t, "active", mapped.Status)
	assert.Equal(t, "usr_70", mapped.OwnerID)
	assert.False(t, mapped.IsDraft)
	assert.True(t, mapped.IsDirty)

	expectedCreated, _ := time.Parse(time.RFC3339, "2026-08-18T14:00:00Z")
	expectedUpdated, _ := time.Parse(time.RFC3339, "2026-08-18T14:30:00Z")
	assert.True(t, mapped.CreatedAt.Equal(expectedCreated))
	assert.True(t, mapped.UpdatedAt.Equal(expectedUpdated))
}

func TestCloudWorkspaceDTOMapping_ZeroValuesAndUnexpectedFields(t *testing.T) {
	// Payload with extra unknown fields and zero-value timestamps
	cloudJSON := `{
		"status": 200,
		"data": [
			{
				"ID": "ws_zero",
				"VaultID": "v_zero",
				"Name": "Minimal Workspace",
				"Description": "",
				"Status": "active",
				"OwnerID": "usr_99",
				"UnexpectedField": "ignored_value",
				"ExtraNumber": 42
			}
		],
		"message": "success",
		"success": true
	}`

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.CloudWorkspaceDTO]
	err := json.Unmarshal([]byte(cloudJSON), &cloudResp)
	require.NoError(t, err)
	require.Len(t, cloudResp.Data, 1)

	mapped := tracecore_types.MapCloudWorkspaceToTypes(&cloudResp.Data[0])
	require.NotNil(t, mapped)

	assert.Equal(t, "ws_zero", mapped.ID)
	assert.Equal(t, "v_zero", mapped.VaultID)
	assert.Equal(t, "Minimal Workspace", mapped.Name)
	assert.Equal(t, "", mapped.Description)
	assert.Equal(t, "active", mapped.Status)
	assert.Equal(t, "usr_99", mapped.OwnerID)
	assert.True(t, mapped.CreatedAt.IsZero())
}

func TestCloudWorkspaceDTOMapping_NilMapper(t *testing.T) {
	assert.Nil(t, tracecore_types.MapCloudWorkspaceToTypes(nil))
}

func TestMapHTTPStatusToError(t *testing.T) {
	assert.True(t, errors.Is(tracecore.MapHTTPStatusToError(http.StatusUnauthorized, ""), tracecore.ErrCloudUnauthorized))
	assert.True(t, errors.Is(tracecore.MapHTTPStatusToError(http.StatusForbidden, ""), tracecore.ErrVaultForbidden))
	assert.True(t, errors.Is(tracecore.MapHTTPStatusToError(http.StatusNotFound, ""), tracecore.ErrResourceNotFound))
	assert.True(t, errors.Is(tracecore.MapHTTPStatusToError(http.StatusBadRequest, "invalid signature"), tracecore.ErrRegistrationBadRequest))
	assert.True(t, errors.Is(tracecore.MapHTTPStatusToError(http.StatusInternalServerError, "database error"), tracecore.ErrCloudServerError))
}
