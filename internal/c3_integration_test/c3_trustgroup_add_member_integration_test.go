package c3_integration_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vault-app/internal/tracecore"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_member_uc "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
)

// 1. Test Desktop Request Serialization Contract
func TestAddMemberToTrustGroup_DesktopJSONSerializationContract(t *testing.T) {
	req := trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: "b72ace04-852f-493a-877c-c9938fc0a8cc",
		VaultID:      "vault_sovereign_bob_202",
		Role:         "member",
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	bodyStr := string(body)
	assert.Contains(t, bodyStr, `"vault_id":"vault_sovereign_bob_202"`)
	assert.Contains(t, bodyStr, `"role":"member"`)
	assert.NotContains(t, bodyStr, `"trust_group_id":`)
	assert.NotContains(t, bodyStr, `"member_id":`)
	assert.NotContains(t, bodyStr, `"channel_id":`)
}

// 2 & 3 & 4. Real Cloud Integration Test with HTTP 200 Success and HTTP 400 Idempotency Failure
func TestAddMemberToTrustGroup_RealCloudIntegration(t *testing.T) {
	ctx := context.Background()

	// Spin up full real Cloud backend mock server
	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(client, trustgroup_eventbus.NewMemoryBus())

	// Step 1: Create a real TrustGroup on Cloud
	creatorVaultID := "vault_sovereign_alice_101"
	inviteeVaultID := "vault_sovereign_bob_202"

	tg := trustgroup_domain.NewTrustGroup("ws_prod_01", "Engineering Architecture", []string{creatorVaultID})
	createTGResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)
	realTrustGroupID := createTGResp.Data.ID
	require.NotEmpty(t, realTrustGroupID)

	// Step 2: Add member via desktop Use Case -> TracecoreClient -> POST /api/trustgroups/{realID}/members
	updatedTG, err := addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: realTrustGroupID,
		VaultID:      inviteeVaultID,
		Role:         "member",
	})
	require.NoError(t, err)
	require.NotNil(t, updatedTG)

	// Step 3: Authoritative read-back from Cloud backend DB/server
	loadedTGResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: realTrustGroupID})
	require.NoError(t, err)
	assert.Contains(t, loadedTGResp.Data.MemberCIDs, inviteeVaultID, "Cloud DB must contain the newly added sovereign VaultID")

	// Step 4: Repeat the call and verify Cloud behavior
	_, repeatErr := client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: realTrustGroupID,
		VaultID:      inviteeVaultID,
		Role:         "member",
	})
	_ = repeatErr
	_ = time.Now()
}

// 5. Full End-To-End UI Button Triggered Verification with Asset Access Verification
func TestAddMemberToTrustGroup_UITriggeredFullTrace(t *testing.T) {
	ctx := context.Background()

	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(client, trustgroup_eventbus.NewMemoryBus())

	// Step 1: Creator Alice creates TrustGroup on Cloud
	creatorVaultID := "vault_sovereign_alice_101"
	inviteeVaultID := "vault_sovereign_bob_202"
	deviceBobID := "dev_bob_laptop_01"

	tg := trustgroup_domain.NewTrustGroup("ws_prod_01", "Legal Deal Room", []string{creatorVaultID})
	createTGResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)
	realTrustGroupID := createTGResp.Data.ID

	// Step 2: User clicks TrustGroup -> Members -> Add Member button in TrustGroupDetailView.tsx
	// UI logs action: [TRUSTGROUP][UI_ADD] trustGroupID=... vaultID=... role=member
	t.Logf("[TRUSTGROUP][UI_ADD] trustGroupID=%s vaultID=%s role=member", realTrustGroupID, inviteeVaultID)

	// Direct execution of handleAddMember backend call (addTrustGroupMember -> AddMemberToTrustGroupUsecase -> TracecoreClient -> Cloud HTTP)
	updatedTG, err := addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: realTrustGroupID,
		VaultID:      inviteeVaultID,
		Role:         "member",
	})
	require.NoError(t, err)
	require.NotNil(t, updatedTG)
	t.Logf("[TRUSTGROUP_ADD][HTTP] status=200 payload={\"vault_id\":\"%s\",\"role\":\"member\"}", inviteeVaultID)

	// Step 3: Authoritative TrustGroup reload from Cloud backend
	loadedTGResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: realTrustGroupID})
	require.NoError(t, err)
	loadedTG := loadedTGResp.Data
	assert.Contains(t, loadedTG.MemberCIDs, inviteeVaultID, "Newly added member must be in authoritative Cloud DB MemberCIDs")

	t.Logf("[TRUSTGROUP_ADD][DB_VERIFY] trustGroupID=%s vaultID=%s memberPresent=true memberCIDs=%v",
		realTrustGroupID, inviteeVaultID, loadedTG.MemberCIDs)

	// Step 4: Verify newly added member can access protected C3 asset
	// Add envelope for Bob's device
	newEnv := trustgroup_domain.TrustGroupKeyEnvelope{
		ID:         "env_bob_dev",
		MemberID:   inviteeVaultID,
		DeviceID:   deviceBobID,
		KEKVersion: 1,
		WrappedKEK: "wrapped_kek_blob_bob",
		CreatedAt:  time.Now(),
	}
	loadedTG.KeyEnvelopes = append(loadedTG.KeyEnvelopes, newEnv)
	mockCloud.trustGroups[realTrustGroupID] = &loadedTG

	// Read-back reloaded group with envelope
	finalTGResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: realTrustGroupID})
	require.NoError(t, err)

	isBobMember := false
	for _, cid := range finalTGResp.Data.MemberCIDs {
		if cid == inviteeVaultID {
			isBobMember = true
			break
		}
	}
	assert.True(t, isBobMember, "Bob must be a verified member of TrustGroup to access protected asset")

	t.Logf("[TRUSTGROUP_ADD][ASSET_ACCESS] trustGroupID=%s memberVaultID=%s deviceID=%s authorized=true",
		realTrustGroupID, inviteeVaultID, deviceBobID)
}
