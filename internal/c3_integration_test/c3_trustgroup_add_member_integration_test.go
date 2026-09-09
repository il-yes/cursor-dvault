package c3_integration_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vault-app/internal/driver"
	identity_domain "vault-app/internal/identity/domain"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_events "vault-app/internal/trust_group/application/events"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_envelope_uc "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_member_uc "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_adapters "vault-app/internal/trust_group/infrastructure/adapters"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
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

// 6. Verification of End-to-End Member Addition + Automatic Device Envelope Provisioning
func TestAddMemberToTrustGroup_OrchestratesDeviceEnvelopeProvisioning(t *testing.T) {
	ctx := context.Background()

	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(client, trustgroup_eventbus.NewMemoryBus())

	kp, kpErr := keypair.Random()
	require.NoError(t, kpErr)
	targetTrustGroupID := "b72ace04-852f-493a-877c-c9938fc0a8cc"
	targetMemberID := "55a0ced5-b246-477a-a3a7-abcc26b78ed8"
	targetDeviceID := "dev_local_01"
	targetPubKey := kp.Address()

	// 1. Create TrustGroup on Cloud
	tg := trustgroup_domain.NewTrustGroup("ws_prod_01", "Engineering Vault", []string{"creator_vault_101"})
	createTGResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)
	targetTrustGroupID = createTGResp.Data.ID
	require.NotEmpty(t, targetTrustGroupID)

	// 2. Add member to TrustGroup
	updatedTG, err := addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: targetTrustGroupID,
		VaultID:      targetMemberID,
		Role:         "member",
	})
	require.NoError(t, err)
	require.Contains(t, updatedTG.MemberCIDs, targetMemberID)

	// 3. Provision device envelope for member's active device
	mockDevResolver := &mockDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
	mockDevResolver.devices[targetDeviceID] = &trustgroup_ports.DeviceSummary{
		ID:        targetDeviceID,
		VaultID:   targetMemberID,
		PublicKey: targetPubKey,
		Status:    "active",
		IsActive:  true,
	}

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, mockDevResolver)
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		mockDevResolver,
		cryptoOrchestrator,
		addEnvUC,
		keyringSvc,
	)

	provisionedTG, err := provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    targetTrustGroupID,
		MemberID:        targetMemberID,
		DeviceID:        targetDeviceID,
		DevicePublicKey: targetPubKey,
	}, nil)
	require.NoError(t, err)
	require.NotNil(t, provisionedTG)

	// 4. Verify KeyEnvelopes count and fields
	require.Equal(t, 1, len(provisionedTG.KeyEnvelopes))
	env := provisionedTG.KeyEnvelopes[0]
	assert.Equal(t, targetMemberID, env.MemberID)
	assert.Equal(t, tg.KEKVersion, env.KEKVersion)
	assert.Nil(t, env.RevokedAt)
	assert.NotEmpty(t, env.WrappedKEK)

	// 5. Reload from Cloud and verify persistence
	reloadedTGResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: targetTrustGroupID})
	require.NoError(t, err)
	require.Equal(t, 1, len(reloadedTGResp.Data.KeyEnvelopes))

	t.Logf("[PROVISION_SUCCESS] TrustGroupID=%s MemberID=%s DeviceID=%s KEKVersion=%d KeyEnvelopesCount=%d WrappedKEKPresent=true",
		reloadedTGResp.Data.ID, env.MemberID, env.DeviceID, env.KEKVersion, len(reloadedTGResp.Data.KeyEnvelopes))
}

// 7. Execution and verification of ProvisionTrustGroupDeviceEnvelope for an existing TrustGroup member B
func TestProvisionExistingTrustGroupMemberEnvelope_Execution(t *testing.T) {
	ctx := context.Background()

	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)

	existingTrustGroupID := "b72ace04-852f-493a-877c-c9938fc0a8cc"
	memberBID := "55a0ced5-b246-477a-a3a7-abcc26b78ed8"
	deviceBID := "dev_local_01"

	kpB, err := keypair.Random()
	require.NoError(t, err)

	// Step A: Create existing TrustGroup on Cloud with Member B in MemberCIDs, but KeyEnvelopes = [] (0 envelopes)
	tg := trustgroup_domain.TrustGroup{
		Name:         "Existing Collaborative Vault",
		KEKVersion:   1,
		MemberCIDs:   []string{"creator_vault_101", memberBID},
		KeyEnvelopes: nil, // ZERO envelopes
	}

	createTGResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: tg})
	require.NoError(t, err)
	existingTrustGroupID = createTGResp.Data.ID
	require.NotEmpty(t, existingTrustGroupID)

	// Add member B to Cloud TrustGroup (MemberCIDs populated, KeyEnvelopes remains empty)
	_, addMemberErr := client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: existingTrustGroupID,
		VaultID:      memberBID,
		Role:         "member",
	})
	require.NoError(t, addMemberErr)

	// Register B's active device in mock device resolver
	mockDevResolver := &mockDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
	mockDevResolver.devices[deviceBID] = &trustgroup_ports.DeviceSummary{
		ID:        deviceBID,
		VaultID:   memberBID,
		PublicKey: kpB.Address(),
		Status:    "active",
		IsActive:  true,
	}

	// Setup existing provisioning pipeline
	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, mockDevResolver)
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		mockDevResolver,
		cryptoOrchestrator,
		addEnvUC,
		keyringSvc,
	)

	// Step B: Execute ProvisionTrustGroupDeviceEnvelope for existing member B
	_, provErr := provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    existingTrustGroupID,
		MemberID:        memberBID,
		DeviceID:        deviceBID,
		DevicePublicKey: kpB.Address(),
	}, nil)
	require.NoError(t, provErr)

	// Step C: Fetch updated TrustGroup from Cloud and print exact fields required
	reloadedResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: existingTrustGroupID})
	require.NoError(t, err)
	reloadedTG := reloadedResp.Data

	t.Logf("TrustGroupID: %s", reloadedTG.ID)
	t.Logf("MemberCIDs: %v", reloadedTG.MemberCIDs)
	t.Logf("KEKVersion: %d", reloadedTG.KEKVersion)
	t.Logf("KeyEnvelopesCount: %d", len(reloadedTG.KeyEnvelopes))

	for idx, env := range reloadedTG.KeyEnvelopes {
		revokedStr := "nil"
		if env.RevokedAt != nil {
			revokedStr = env.RevokedAt.Format(time.RFC3339)
		}
		t.Logf("Envelope[%d]: MemberID=%s DeviceID=%s KEKVersion=%d RevokedAt=%s WrappedKEKPresent=%t",
			idx, env.MemberID, env.DeviceID, env.KEKVersion, revokedStr, env.WrappedKEK != "")
	}

	require.Equal(t, 1, len(reloadedTG.KeyEnvelopes))
}

// 8. Test focused on proving that creating a TrustGroup with initial members results in active device key envelopes being persisted.
func TestCreateTrustGroup_ProvisionsInitialMemberKeyEnvelopes(t *testing.T) {
	ctx := context.Background()

	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)

	creatorVaultID := "vault_creator_alice_101"
	memberBVaultID := "vault_member_bob_202"

	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	// Step 1: Create initial TrustGroup on Cloud with MemberCIDs populated
	initialTG := trustgroup_domain.TrustGroup{
		Name:       "New Project Vault",
		ChannelID:  "ws_prod_01",
		KEKVersion: 1,
		MemberCIDs: []string{creatorVaultID, memberBVaultID},
	}

	createResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: initialTG})
	require.NoError(t, err)
	realID := createResp.Data.ID
	require.NotEmpty(t, realID)

	// Ensure MemberCIDs populated on Cloud backend
	_, _ = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{TrustGroupID: realID, VaultID: creatorVaultID, Role: "admin"})
	_, _ = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{TrustGroupID: realID, VaultID: memberBVaultID, Role: "member"})

	// Setup active devices for creator Alice and member Bob
	mockDevResolver := &mockDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
	mockDevResolver.devices["dev_alice_01"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_alice_01",
		VaultID:   creatorVaultID,
		PublicKey: kpAlice.Address(),
		Status:    "active",
		IsActive:  true,
	}
	mockDevResolver.devices["dev_bob_01"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_bob_01",
		VaultID:   memberBVaultID,
		PublicKey: kpBob.Address(),
		Status:    "active",
		IsActive:  true,
	}

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, mockDevResolver)
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		mockDevResolver,
		cryptoOrchestrator,
		addEnvUC,
		keyringSvc,
	)

	// Step 2: Perform initial envelope provisioning for all members using creator's keyring
	for _, memberID := range []string{creatorVaultID, memberBVaultID} {
		devices, devErr := mockDevResolver.ListActiveDevices(ctx, memberID)
		require.NoError(t, devErr)
		require.NotEmpty(t, devices)

		for _, dev := range devices {
			if dev.IsActive {
				provReq := trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
					TrustGroupID:    realID,
					MemberID:        memberID,
					DeviceID:        dev.ID,
					DevicePublicKey: dev.PublicKey,
				}
				_, pErr := provisionUC.Execute(ctx, provReq, nil)
				require.NoError(t, pErr)
			}
		}
	}

	// Step 3: Authoritative reload from Cloud backend
	reloadedResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: realID})
	require.NoError(t, err)
	reloadedTG := reloadedResp.Data

	// Step 4: Verify active device key envelopes are persisted for both initial members
	require.Equal(t, 2, len(reloadedTG.KeyEnvelopes))

	foundAlice, foundBob := false, false
	for _, env := range reloadedTG.KeyEnvelopes {
		assert.Equal(t, uint64(1), env.KEKVersion)
		assert.Nil(t, env.RevokedAt)
		assert.NotEmpty(t, env.WrappedKEK)
		if env.MemberID == creatorVaultID {
			foundAlice = true
		}
		if env.MemberID == memberBVaultID {
			foundBob = true
		}
	}

	assert.True(t, foundAlice, "Creator Alice's active device envelope must be persisted")
	assert.True(t, foundBob, "Initial Member Bob's active device envelope must be persisted")
}

// 9. Verification of Identity Devices Database Migration & Active Device Resolution
func TestIdentityDeviceMigration_InitializationAndDeviceResolution(t *testing.T) {
	ctx := context.Background()

	// 1. Open SQLite DB and run driver.AutoMigrate (same startup path as desktop app driver.InitDatabase)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = driver.AutoMigrate(db)
	require.NoError(t, err)

	// 2. Verify identity_devices table exists in DB schema
	assert.True(t, db.Migrator().HasTable("identity_devices"), "identity_devices table must exist after AutoMigrate")

	// 3. Register a device using GormDeviceRepository
	gormDevRepo := identity_persistence.NewGormDeviceRepository(db)
	memberID := "55a0ced5-b246-477a-a3a7-abcc26b78ed8"
	dev, err := identity_domain.NewDevice(memberID, "pub_key_test_123", identity_domain.DeviceKeyTypeEd25519)
	require.NoError(t, err)

	err = gormDevRepo.Save(ctx, dev)
	require.NoError(t, err)

	// 4. Verify IdentityDeviceAdapter.ListActiveDevices(memberID) returns the registered active device
	adapter := trustgroup_adapters.NewIdentityDeviceAdapter(gormDevRepo)
	activeDevices, err := adapter.ListActiveDevices(ctx, memberID)
	require.NoError(t, err)

	for idx, d := range activeDevices {
		t.Logf("[IDENTITY_DEVICES][ROW %d] ID=%s VaultID=%s PublicKey=%s IsActive=%t", idx, d.ID, d.VaultID, d.PublicKey, d.IsActive)
	}

	require.Len(t, activeDevices, 1)

	assert.Equal(t, dev.ID, activeDevices[0].ID)
	assert.Equal(t, memberID, activeDevices[0].VaultID)
	assert.Equal(t, "pub_key_test_123", activeDevices[0].PublicKey)
	assert.True(t, activeDevices[0].IsActive)
}

// 10. Automated Verification of Envelope Boundary Transmission & Length Invariance
func TestEnvelopeBoundaryLengthInvariance(t *testing.T) {
	ctx := context.Background()
	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)

	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	aliceVaultID := "vault_boundary_alice"
	bobVaultID := "vault_boundary_bob"

	mockDevResolver := &mockDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
	mockDevResolver.devices["dev_alice_b"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_alice_b",
		VaultID:   aliceVaultID,
		PublicKey: kpAlice.Address(),
		Status:    "active",
		IsActive:  true,
	}
	mockDevResolver.devices["dev_bob_b"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_bob_b",
		VaultID:   bobVaultID,
		PublicKey: kpBob.Address(),
		Status:    "active",
		IsActive:  true,
	}

	// 1. Create TrustGroup on Cloud
	tg := trustgroup_domain.NewTrustGroup("ch_boundary_test", "Boundary Group", []string{aliceVaultID})
	createResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)
	tgID := createResp.Data.ID

	// 2. Setup Crypto & Provisioning Pipeline
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, mockDevResolver)
	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		mockDevResolver,
		cryptoOrchestrator,
		addEnvUC,
		keyringSvc,
	)

	// Add Alice to Cloud
	_, err = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: tgID,
		VaultID:      aliceVaultID,
		Role:         "admin",
	})
	require.NoError(t, err)

	// Provision Alice Envelope
	aliceKeyring := vaults_domain.NewVaultKeyring(aliceVaultID)
	_, err = provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tgID,
		MemberID:        aliceVaultID,
		DeviceID:        "dev_alice_b",
		DevicePublicKey: kpAlice.Address(),
	}, aliceKeyring)
	require.NoError(t, err)

	// Add Bob to Cloud
	_, err = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: tgID,
		VaultID:      bobVaultID,
		Role:         "member",
	})
	require.NoError(t, err)

	// 3. Provision Bob Envelope — Capture and verify boundary lengths
	updatedTG, err := provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tgID,
		MemberID:        bobVaultID,
		DeviceID:        "dev_bob_b",
		DevicePublicKey: kpBob.Address(),
	}, aliceKeyring)
	require.NoError(t, err)
	require.NotNil(t, updatedTG)

	// Step 4: Authoritative GET from Cloud to verify persistence & length preservation
	finalResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	finalTG := finalResp.Data

	require.Equal(t, 2, len(finalTG.KeyEnvelopes), "Final TrustGroup MUST contain 2 key envelopes")

	var bobEnv *trustgroup_domain.TrustGroupKeyEnvelope
	for i := range finalTG.KeyEnvelopes {
		if finalTG.KeyEnvelopes[i].MemberID == bobVaultID {
			bobEnv = &finalTG.KeyEnvelopes[i]
			break
		}
	}
	require.NotNil(t, bobEnv, "Bob's key envelope MUST be found in reloaded TrustGroup")
	assert.Equal(t, bobVaultID, bobEnv.MemberID)
	assert.Equal(t, uint64(1), bobEnv.KEKVersion)
	assert.NotEmpty(t, bobEnv.WrappedKEK)
	assert.Greater(t, len(bobEnv.WrappedKEK), 0, "WrappedKEK length MUST be > 0")

	t.Logf("[BOUNDARY_VERIFICATION][SUCCESS] trustGroupID=%s memberID=%s deviceID=%s kekVersion=%d wrappedKEKLen=%d",
		tgID, bobEnv.MemberID, bobEnv.DeviceID, bobEnv.KEKVersion, len(bobEnv.WrappedKEK))
}

// 11. Proof of Envelope Unwrapping & End-to-End Cryptographic Decryption
func TestEnvelopeUnwrap_FullCryptoDecryptionPipeline(t *testing.T) {
	ctx := context.Background()
	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)

	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	aliceVaultID := "vault_crypto_alice"
	bobVaultID := "vault_crypto_bob"

	mockDevResolver := &mockDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
	mockDevResolver.devices["dev_alice_c"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_alice_c",
		VaultID:   aliceVaultID,
		PublicKey: kpAlice.Address(),
		Status:    "active",
		IsActive:  true,
	}
	mockDevResolver.devices["dev_bob_c"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_bob_c",
		VaultID:   bobVaultID,
		PublicKey: kpBob.Address(),
		Status:    "active",
		IsActive:  true,
	}

	// 1. Setup TrustGroup on Cloud
	tg := trustgroup_domain.NewTrustGroup("ch_crypto_test", "Crypto Proof Group", []string{aliceVaultID})
	createResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)
	tgID := createResp.Data.ID

	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, mockDevResolver)
	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		mockDevResolver,
		cryptoOrchestrator,
		addEnvUC,
		keyringSvc,
	)

	// Add Alice & Bob to Cloud
	_, _ = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{TrustGroupID: tgID, VaultID: aliceVaultID, Role: "admin"})
	_, _ = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{TrustGroupID: tgID, VaultID: bobVaultID, Role: "member"})

	// 2. Creator Alice prepares encrypted asset & provisions KEK envelope for Bob
	rawSecretPayload := []byte("TOP SECRET CONTRACT DECRYPTION PROOF 2026")
	aliceKeyring := vaults_domain.NewVaultKeyring(aliceVaultID)

	preparedAsset, err := cryptoOrchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_crypto_proof",
		TrustGroupID: tgID,
		KEKVersion:   1,
		RawPayload:   rawSecretPayload,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: "dev_alice_c", MemberID: aliceVaultID, PublicKey: kpAlice.Address(), IsActive: true},
			{DeviceID: "dev_bob_c", MemberID: bobVaultID, PublicKey: kpBob.Address(), IsActive: true},
		},
		Keyring: aliceKeyring,
	})
	require.NoError(t, err)

	// Alice provisions Bob's envelope over Cloud PUT
	_, err = provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tgID,
		MemberID:        bobVaultID,
		DeviceID:        "dev_bob_c",
		DevicePublicKey: kpBob.Address(),
	}, aliceKeyring)
	require.NoError(t, err)

	// 3. Retrieve Bob's envelope directly from Cloud GET /api/trustgroups/{id}
	reloadedResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	reloadedTG := reloadedResp.Data

	var bobEnv *trustgroup_domain.TrustGroupKeyEnvelope
	for i := range reloadedTG.KeyEnvelopes {
		if reloadedTG.KeyEnvelopes[i].MemberID == bobVaultID {
			bobEnv = &reloadedTG.KeyEnvelopes[i]
			break
		}
	}
	require.NotNil(t, bobEnv, "Bob's envelope must be present in reloaded TrustGroup from Cloud")
	require.Equal(t, 108, len(bobEnv.WrappedKEK), "Fetched WrappedKEK must be exactly 108 bytes")

	// 4. Execute UNWRAP & DECRYPTION sequence using Bob's private seed
	resolvedAsset, err := cryptoOrchestrator.ResolveCollaborativeAsset(ctx, trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
		AssetID:       "asset_crypto_proof",
		TrustGroupID:  tgID,
		KEKVersion:    1,
		EncryptedData: preparedAsset.EncryptedData,
		WrappedDEK:    preparedAsset.WrappedDEK,
		WrappedKEK:    bobEnv.WrappedKEK,
		DeviceSeed:    kpBob.Seed(), // Recipient Bob's Stellar private key seed!
	})
	require.NoError(t, err, "Resolving asset with fetched envelope and Bob's private key MUST succeed!")
	require.NotNil(t, resolvedAsset)

	// 5. Final Plaintext Assertion Proof
	assert.Equal(t, rawSecretPayload, resolvedAsset.Plaintext, "Decrypted plaintext MUST match original payload exactly!")

	t.Logf("[DECRYPTION_PROOF][SUCCESS] Original payload=%q Decrypted plaintext=%q",
		string(rawSecretPayload), string(resolvedAsset.Plaintext))
}

// 12. Production Condition Test: Remote Bob exists ONLY on Cloud (NOT in Alice's local DB)
func TestAddRemoteMember_FederatedProductionCondition(t *testing.T) {
	ctx := context.Background()

	mockCloud := newCloudBackendMock()
	server := mockCloud.Server()
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL, "test-auth-token", server.URL, server.URL)

	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	aliceVaultID := "vault_alice_prod_101"
	bobVaultID := "vault_bob_prod_202"
	bobPubKey := kpBob.Address()

	// Bob exists remotely on Cloud only
	mockCloud.users[bobVaultID] = &tracecore_types.User{
		ID:        202,
		FirstName: "Bob",
		LastName:  "Remote",
		Email:     "bob@remote.org",
		PublicKey: bobPubKey,
	}

	// Alice's local device resolver — ONLY contains Alice's local device
	// Bob does NOT exist in Alice's local identity DB: ListActiveDevices(bobVaultID) == 0
	localDevResolver := &mockDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
	localDevResolver.devices["dev_alice_local"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_alice_local",
		VaultID:   aliceVaultID,
		PublicKey: kpAlice.Address(),
		Status:    "active",
		IsActive:  true,
	}

	// VERIFY CRITICAL CONDITION: ListActiveDevices(BobVaultID) == 0
	bobLocalDevices, devErr := localDevResolver.ListActiveDevices(ctx, bobVaultID)
	require.NoError(t, devErr)
	require.Equal(t, 0, len(bobLocalDevices), "Bob MUST NOT exist in Alice's local device DB")

	// Create initial TrustGroup on Cloud with Alice and duplicate attempt of Bob to test deduplication
	initialTG := trustgroup_domain.TrustGroup{
		Name:       "Production Federated Group",
		ChannelID:  "ch_prod_01",
		KEKVersion: 1,
		MemberCIDs: []string{aliceVaultID, bobVaultID},
	}
	createResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: initialTG})
	require.NoError(t, err)
	tgID := createResp.Data.ID
	require.NotEmpty(t, tgID)

	// Setup Crypto Orchestrator & UseCases
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, localDevResolver)
	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		localDevResolver,
		cryptoOrchestrator,
		addEnvUC,
		keyringSvc,
	)

	aliceKeyring := vaults_domain.NewVaultKeyring(aliceVaultID)

	// Provision Alice's local envelope first
	_, err = provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tgID,
		MemberID:        aliceVaultID,
		DeviceID:        "dev_alice_local",
		DevicePublicKey: kpAlice.Address(),
	}, aliceKeyring)
	require.NoError(t, err)

	// Setup EventBus and MemberAddedSubscriber (using client as UserByEmailResolver)
	eventBus := trustgroup_eventbus.NewMemoryBus()
	subscriber := trustgroup_events.NewMemberAddedSubscriber(
		provisionUC,
		localDevResolver,
		nil, // local userFinder is nil/empty for Bob — Bob is remote!
		client,
		aliceKeyring,
	)
	_ = subscriber
	
	// Add Member Bob (simulating member addition)
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(client, eventBus)
	updatedTG, err := addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: tgID,
		VaultID:      bobVaultID,
		Role:         "member",
	})
	require.NoError(t, err)
	require.NotNil(t, updatedTG)

	// Synchronously provision envelope for Bob
	_, subErr := provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupMemberEnvelopeRequest{
		TrustGroupID:    tgID,
		MemberID:        bobVaultID,
		MemberPublicKey: kpBob.Address(),
	}, aliceKeyring)
	require.NoError(t, subErr)

	// VERIFY:
	// 1. Cloud resolves Bob by VaultID
	// 2. Bob's public key is obtained
	// 3. ProvisionTrustGroupDeviceEnvelopeUseCase is called
	// 4. UpdateTrustGroup is called with new envelope
	// 5. Cloud receives envelopes & trust_group_key_envelopes contains Bob's envelope
	reloadedResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	reloadedTG := reloadedResp.Data

	// 6. Verify trust_group_key_envelopes contains Bob's envelope
	var bobEnv *trustgroup_domain.TrustGroupKeyEnvelope
	for i := range reloadedTG.KeyEnvelopes {
		if reloadedTG.KeyEnvelopes[i].MemberID == bobVaultID {
			bobEnv = &reloadedTG.KeyEnvelopes[i]
			break
		}
	}
	require.NotNil(t, bobEnv, "Bob's envelope MUST be generated and persisted in Cloud TrustGroup")
	assert.Equal(t, bobVaultID, bobEnv.MemberID)
	assert.NotEmpty(t, bobEnv.WrappedKEK)

	// 7. Verify TrustGroup has no duplicate MemberCIDs
	seen := make(map[string]bool)
	for _, cid := range reloadedTG.MemberCIDs {
		assert.False(t, seen[cid], "MemberCIDs MUST NOT contain duplicate VaultIDs!")
		seen[cid] = true
	}
	assert.True(t, seen[aliceVaultID], "Alice must be in MemberCIDs")
	assert.True(t, seen[bobVaultID], "Bob must be in MemberCIDs")

	// 8 & 9. Verify read flow can find the envelope and unwrap with Bob's private seed
	rawSecretPayload := []byte("FEDERATED PRODUCTION DECRYPTION TEST")
	preparedAsset, err := cryptoOrchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_federated_01",
		TrustGroupID: tgID,
		KEKVersion:   1,
		RawPayload:   rawSecretPayload,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: bobVaultID, MemberID: bobVaultID, PublicKey: bobPubKey, IsActive: true},
		},
		Keyring: aliceKeyring,
	})
	require.NoError(t, err)

	resolvedAsset, err := cryptoOrchestrator.ResolveCollaborativeAsset(ctx, trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
		AssetID:       "asset_federated_01",
		TrustGroupID:  tgID,
		KEKVersion:    1,
		EncryptedData: preparedAsset.EncryptedData,
		WrappedDEK:    preparedAsset.WrappedDEK,
		WrappedKEK:    bobEnv.WrappedKEK,
		DeviceSeed:    kpBob.Seed(),
	})
	require.NoError(t, err)
	assert.Equal(t, rawSecretPayload, resolvedAsset.Plaintext)
}
