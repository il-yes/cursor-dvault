package c3_integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_ui "vault-app/internal/collaboration/ui"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

func TestC3EntryShared_RecipientReadFlow_NoDecryption(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	senderID := "vault_sender_uuid"
	recipientID := "vault_recipient_uuid"
	channelID := "ch_contract_execution_1"

	// 1. Setup TrustGroup & Thread
	tg := trustgroup_domain.NewTrustGroup(channelID, "Legal & Finance Group", []string{senderID, recipientID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	thread := thread_domain.NewThread(channelID, "contract", "Executive Agreement", "v1.0")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	// 2. Sender creates C3 ShareEntry (without performing decryption)
	shareEntry := c3_asset_domain.ShareEntry{
		ID:           "se_authoritative_999",
		AssetCID:     "bafybeicontractasset2026",
		TrustGroupID: tg.ID,
		WrappedDEK:   "wrapped_dek_string_content",
		KEKVersion:   1,
		CreatedBy:    senderID,
		Status:       c3_asset_domain.ShareEntryStatusActive,
		Metadata: map[string]string{
			"title":  "Executive Agreement Document",
			"vendor": "Acme Corp",
		},
	}
	repo.shareEntries[shareEntry.ID] = shareEntry

	// 3. Sender appends entry.shared event referencing shareEntry.ID
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}

	appendedEvt, err := appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", refPayload, "evt_share_"+shareEntry.ID)
	require.NoError(t, err)
	require.NotNil(t, appendedEvt)

	// 4. Recipient fetches thread events timeline (ListThreadEvents)
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(repo)
	recipientEvents, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, recipientEvents, 1)

	recEvt := recipientEvents[0]
	assert.Equal(t, "entry.shared", string(recEvt.Type))

	// 5. Recipient extracts the real ShareEntry ID from EventResourceRef
	extractedShareEntryID := recEvt.Payload.ShareEntryID
	extractedTrustGroupID := recEvt.Payload.TrustGroupID
	assert.Equal(t, shareEntry.ID, extractedShareEntryID, "Recipient must extract the real ShareEntry ID")
	assert.Equal(t, tg.ID, extractedTrustGroupID, "Recipient must extract the TrustGroup ID")

	// 6. Recipient resolves the ShareEntry via the C3 ShareEntry read path (GetShareEntry)
	// (NO decryption is invoked here)
	resolvedShare, err := repo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{ShareEntryID: extractedShareEntryID})
	require.NoError(t, err)
	require.NotNil(t, resolvedShare)

	// 7. Recipient verifies C3 ShareEntry metadata
	assert.Equal(t, shareEntry.ID, resolvedShare.Data.ID)
	assert.Equal(t, tg.ID, resolvedShare.Data.TrustGroupID)
	assert.Equal(t, c3_asset_domain.ShareEntryStatusActive, resolvedShare.Data.Status)
	assert.Equal(t, senderID, resolvedShare.Data.CreatedBy)
	assert.Equal(t, "bafybeicontractasset2026", resolvedShare.Data.AssetCID)
	assert.Equal(t, uint64(1), resolvedShare.Data.KEKVersion)
	assert.Equal(t, "Executive Agreement Document", resolvedShare.Data.Metadata["title"])
}

func TestC3EntryShared_RecipientReadFlow_WithDecryption(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	// 1. Setup Keypairs for Vault A (Alice/Sender) & Vault B (Bob/Recipient)
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	kpBob, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_uuid"
	userBobID := "user_bob_uuid"
	deviceAliceID := "device_laptop_alice"
	deviceBobID := "device_laptop_bob"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "vault_alice"}

	repo.seeds[userBobID] = kpBob.Seed()
	repo.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "vault_bob"}

	// 2. Create TrustGroup containing both Alice & Bob
	tg := trustgroup_domain.NewTrustGroup("ch_contract_1", "Executive Board", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	// 3. Create Thread
	thread := thread_domain.NewThread("ch_contract_1", "contract", "Confidential Settlement", "v1.0")
	_, err = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	// 4. Vault A prepares & encrypts asset payload for TrustGroup members
	rawOriginalContent := []byte("secret construction document")
	assetCID := "bafybeisettlement2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_settlement",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	require.Len(t, prepared.Envelopes, 2)

	// Persist encrypted asset into storage repository
	repo.assets[assetCID] = prepared.EncryptedData

	// Add key envelopes for all devices to TrustGroup
	for _, envReq := range prepared.Envelopes {
		err = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
		require.NoError(t, err)
	}
	repo.trustGroups[tg.ID] = *tg

	// Wire Use Cases & Handlers
	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(repo)

	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, resolveCollabShareUC, appendThreadEventUC)

	// 5. Vault A creates C3 ShareEntry and appends entry.shared event
	wrappedDEKStr := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

	shareReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   wrappedDEKStr,
		Metadata:     map[string]string{"title": "Confidential Settlement"},
	}

	createResp, err := createCollabShareUC.Execute(ctx, shareReq)
	require.NoError(t, err)
	createdShareEntryID := createResp.ShareEntry.ID
	repo.shareEntries[createdShareEntryID] = createResp.ShareEntry

	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createdShareEntryID,
		TrustGroupID: tg.ID,
	}
	_, err = appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", refPayload, "evt_share_"+createdShareEntryID)
	require.NoError(t, err)

	// 6. Vault B (Recipient) receives timeline events via ListThreadEvents
	recipientEvents, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, recipientEvents, 1)

	recEvt := recipientEvents[0]
	assert.Equal(t, "entry.shared", string(recEvt.Type))

	// Extract share_entry_id from event payload
	extractedShareEntryID := recEvt.Payload.ShareEntryID
	require.Equal(t, createdShareEntryID, extractedShareEntryID)

	// 7. Vault B resolves ShareEntry metadata via GetShareEntry
	shareEntryFromGet, err := repo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{ShareEntryID: extractedShareEntryID})
	require.NoError(t, err)
	require.NotNil(t, shareEntryFromGet)

	// 8. Vault B executes full cryptographic decryption via ResolveCollaborativeShare
	resolvedDTO, err := collabHandler.ResolveCollaborativeShare(ctx, userBobID, userBobID, extractedShareEntryID)
	require.NoError(t, err)
	require.NotNil(t, resolvedDTO)

	// 9. ASSERT ALL INVARIANTS
	require.Equal(t, rawOriginalContent, resolvedDTO.Plaintext, "Recipient plaintext MUST match sender original plaintext")
	require.Equal(t, createdShareEntryID, recEvt.Payload.ShareEntryID, "Event ShareEntryID MUST match created ShareEntry ID")
	require.Equal(t, tg.ID, shareEntryFromGet.Data.TrustGroupID, "ShareEntry TrustGroupID MUST match TrustGroup ID")
	require.Equal(t, assetCID, shareEntryFromGet.Data.AssetCID, "ShareEntry AssetCID MUST match Asset CID")
}

func TestC3EntryShared_RecipientReadFlow_UnauthorizedRecipient_MetadataVisible_DecryptionDenied(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	// 1. Setup Keypairs for Vault A (Alice/Member) & Vault C (Charlie/Unauthorized Non-Member)
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	kpCharlie, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_uuid"
	userCharlieID := "user_charlie_uuid"
	deviceAliceID := "device_laptop_alice"
	deviceCharlieID := "device_laptop_charlie"
	_ = deviceCharlieID

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "vault_alice"}

	repo.seeds[userCharlieID] = kpCharlie.Seed()
	repo.keyrings[userCharlieID] = &vaults_domain.VaultKeyring{UserID: userCharlieID, VaultID: "vault_charlie"}

	// 2. TrustGroup contains ONLY Alice (Charlie is NOT a member)
	tg := trustgroup_domain.NewTrustGroup("ch_contract_restricted", "Restricted Group", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	// 3. Thread
	thread := thread_domain.NewThread("ch_contract_restricted", "contract", "Restricted Contract", "")
	_, err = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	// 4. Prepare & Encrypt Asset for Alice only
	rawOriginalContent := []byte(`{"secret":"Restricted Board Minutes"}`)
	assetCID := "bafybeirestricted2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_restricted",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	repo.assets[assetCID] = prepared.EncryptedData

	envReq := prepared.Envelopes[0]
	err = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
		TrustGroupID: envReq.TrustGroupID,
		MemberID:     envReq.MemberID,
		DeviceID:     envReq.DeviceID,
		KEKVersion:   envReq.KEKVersion,
		WrappedKEK:   envReq.WrappedKEK,
	})
	require.NoError(t, err)
	repo.trustGroups[tg.ID] = *tg

	// Wire Use Cases & Handler
	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)

	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, resolveCollabShareUC, appendThreadEventUC)

	// 5. Alice creates C3 ShareEntry and appends entry.shared event
	wrappedDEKStr := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)
	shareReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   wrappedDEKStr,
		Metadata:     map[string]string{"title": "Restricted Board Minutes"},
	}

	createResp, err := createCollabShareUC.Execute(ctx, shareReq)
	require.NoError(t, err)
	createdShareEntryID := createResp.ShareEntry.ID
	repo.shareEntries[createdShareEntryID] = createResp.ShareEntry

	// 6. INVARIANT 1: GetShareEntry SUCCEEDS for Charlie (Metadata is visible)
	shareEntryFromGet, err := repo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{ShareEntryID: createdShareEntryID})
	require.NoError(t, err, "GetShareEntry metadata read MUST succeed even for non-member recipient")
	require.NotNil(t, shareEntryFromGet)
	assert.Equal(t, createdShareEntryID, shareEntryFromGet.Data.ID)
	assert.Equal(t, tg.ID, shareEntryFromGet.Data.TrustGroupID)
	assert.Equal(t, assetCID, shareEntryFromGet.Data.AssetCID)
	assert.Equal(t, userAliceID, shareEntryFromGet.Data.CreatedBy)

	// 7. INVARIANT 2: ResolveCollaborativeShare FAILS for Charlie (Decryption is DENIED)
	resolvedDTO, err := collabHandler.ResolveCollaborativeShare(ctx, userCharlieID, userCharlieID, createdShareEntryID)
	assert.Error(t, err, "Decryption resolution MUST fail for unauthorized recipient")
	assert.Nil(t, resolvedDTO, "NO plaintext must be returned to unauthorized recipient")
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember, "Must return ErrUnauthorizedMember")
}

func TestC3ProductionReadPath_ProductInvariantMatrix(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	// 1. Setup Identities
	kpAlice, _ := keypair.Random()
	kpBob, _ := keypair.Random()
	kpCharlie, _ := keypair.Random()

	userAliceID := "usr_alice_owner"
	userBobID := "usr_bob_recipient"
	userCharlieID := "usr_charlie_eavesdropper"

	deviceAliceLaptop := "dev_alice_m2"
	deviceBobLaptop := "dev_bob_thinkpad"
	deviceCharlieLaptop := "dev_charlie_dell"
	_ = deviceCharlieLaptop

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	repo.seeds[userBobID] = kpBob.Seed()
	repo.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}

	repo.seeds[userCharlieID] = kpCharlie.Seed()
	repo.keyrings[userCharlieID] = &vaults_domain.VaultKeyring{UserID: userCharlieID, VaultID: "v_charlie"}

	// 2. TrustGroup contains Alice & Bob (Charlie excluded)
	tg := trustgroup_domain.NewTrustGroup("ch_production_1", "Sovereign Board", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	thread := thread_domain.NewThread("ch_production_1", "contract", "Sovereign Asset Shares", "")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	// 3. Alice prepares & encrypts raw payload
	rawOriginalContent := []byte("top_secret_sovereign_blueprint_2026")
	assetCID := "bafybeisovereignblueprint2026"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_blueprint",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceLaptop, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
			{DeviceID: deviceBobLaptop, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	repo.assets[assetCID] = prepared.EncryptedData

	for _, envReq := range prepared.Envelopes {
		err = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
		require.NoError(t, err)
	}
	repo.trustGroups[tg.ID] = *tg

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(repo)

	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, resolveCollabShareUC, appendThreadEventUC)

	// 4. Alice creates ShareEntry & appends event
	createResp, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
		Metadata:     map[string]string{"title": "Sovereign Blueprint"},
	})
	require.NoError(t, err)
	createdShareEntryID := createResp.ShareEntry.ID
	repo.shareEntries[createdShareEntryID] = createResp.ShareEntry

	_, err = appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createdShareEntryID,
		TrustGroupID: tg.ID,
	}, "evt_share_"+createdShareEntryID)
	require.NoError(t, err)

	// =========================================================================
	// MATRIX STEP A: Authorized Recipient (Bob) Flow
	// =========================================================================
	bobEvents, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, bobEvents, 1)

	bobEvt := bobEvents[0]
	assert.Equal(t, createdShareEntryID, bobEvt.Payload.ShareEntryID)

	bobMeta, err := repo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{ShareEntryID: bobEvt.Payload.ShareEntryID})
	require.NoError(t, err)
	require.NotNil(t, bobMeta)
	assert.Equal(t, tg.ID, bobMeta.Data.TrustGroupID)

	bobResolved, err := collabHandler.ResolveCollaborativeShare(ctx, userBobID, userBobID, createdShareEntryID)
	require.NoError(t, err, "Bob MUST be authorized to resolve and decrypt plaintext")
	require.NotNil(t, bobResolved)
	require.Equal(t, rawOriginalContent, bobResolved.Plaintext, "Bob MUST receive exact original plaintext")

	// =========================================================================
	// MATRIX STEP B: Unauthorized Recipient (Charlie) Flow
	// =========================================================================
	charlieMeta, err := repo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{ShareEntryID: createdShareEntryID})
	require.NoError(t, err, "Charlie can read ShareEntry metadata")
	require.NotNil(t, charlieMeta)

	charlieResolved, err := collabHandler.ResolveCollaborativeShare(ctx, userCharlieID, userCharlieID, createdShareEntryID)
	assert.Error(t, err, "Charlie MUST be denied decryption")
	assert.Nil(t, charlieResolved, "NO plaintext must be returned to Charlie")
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)

	// =========================================================================
	// MATRIX STEP C: Revoked Envelope Flow
	// =========================================================================
	now := time.Now()
	tgAlice := repo.trustGroups[tg.ID]
	for i := range tgAlice.KeyEnvelopes {
		tgAlice.KeyEnvelopes[i].RevokedAt = &now
	}
	repo.trustGroups[tg.ID] = tgAlice

	aliceOldDeviceResolved, err := collabHandler.ResolveCollaborativeShare(ctx, userAliceID, userAliceID, createdShareEntryID)
	assert.Error(t, err, "Revoked envelope MUST be denied decryption")
	assert.Nil(t, aliceOldDeviceResolved)
	assert.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound)
}

func TestC3EntryShared_RealCloudPayloadShape_Regression(t *testing.T) {
	// Test 1: JSON payload with nested share_entry_ref unmarshals into EventResourceRef
	nestedJSON := `{"share_entry_ref":{"share_entry_id":"se_real_cloud_999","trust_group_id":"tg_real_cloud_111"}}`
	var ref1 thread_domain.EventResourceRef
	err := json.Unmarshal([]byte(nestedJSON), &ref1)
	require.NoError(t, err)
	assert.Equal(t, "se_real_cloud_999", ref1.ShareEntryID)
	assert.Equal(t, "tg_real_cloud_111", ref1.TrustGroupID)
	assert.Equal(t, thread_domain.ResourceShareEntry, ref1.RefType)

	// Test 2: JSON payload with nested resource_ref unmarshals into EventResourceRef
	resourceRefJSON := `{"resource_ref":{"share_entry_id":"se_real_cloud_888","trust_group_id":"tg_real_cloud_222"}}`
	var ref2 thread_domain.EventResourceRef
	err = json.Unmarshal([]byte(resourceRefJSON), &ref2)
	require.NoError(t, err)
	assert.Equal(t, "se_real_cloud_888", ref2.ShareEntryID)
	assert.Equal(t, "tg_real_cloud_222", ref2.TrustGroupID)
	assert.Equal(t, thread_domain.ResourceShareEntry, ref2.RefType)

	// Test 3: JSON payload with share_id unmarshals into EventResourceRef
	shareIDJSON := `{"share_id":"se_real_cloud_777","trust_group_id":"tg_real_cloud_333"}`
	var ref3 thread_domain.EventResourceRef
	err = json.Unmarshal([]byte(shareIDJSON), &ref3)
	require.NoError(t, err)
	assert.Equal(t, "se_real_cloud_777", ref3.ShareEntryID)
	assert.Equal(t, "tg_real_cloud_333", ref3.TrustGroupID)

	// Test 4: JSON payload with entry_id unmarshals into EventResourceRef
	entryIDJSON := `{"entry_id":"se_real_cloud_666"}`
	var ref4 thread_domain.EventResourceRef
	err = json.Unmarshal([]byte(entryIDJSON), &ref4)
	require.NoError(t, err)
	assert.Equal(t, "se_real_cloud_666", ref4.ShareEntryID)
}

func TestC3WritePath_CreateCollaborativeShare_AppendsValidEventRef(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	userAliceID := "user_alice_writer"

	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_write_path_test", "Write Path Group", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	thread := thread_domain.NewThread("ch_write_path_test", "contract", "Write Path Thread", "")
	_, err = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(repo)

	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, resolveCollabShareUC, appendThreadEventUC)

	// Execute REAL CreateCollaborativeShare operation on CollaborationHandler (which triggers AppendThreadEvent)
	shareRefDTO, err := collabHandler.CreateCollaborativeShare(
		ctx,
		userAliceID,
		thread.ID,
		tg.ID,
		"bafybeirealwritepathcid2026",
		"target_v_bob",
		"Production write path test",
		"d3JhcHBlZF9kZWtfYmFzZTY0",
		1,
	)
	require.NoError(t, err)
	require.NotNil(t, shareRefDTO)

	// Assertions on returned ShareEntryRefDTO
	require.NotEmpty(t, shareRefDTO.ShareEntryID, "createdShareEntry.ID MUST NOT be empty")
	require.NotEmpty(t, shareRefDTO.TrustGroupID, "createdShareEntry.TrustGroupID MUST NOT be empty")
	assert.Equal(t, tg.ID, shareRefDTO.TrustGroupID)

	// Retrieve actual appended thread events
	events, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, events, 1, "Creation MUST append exactly 1 event")

	appendedEvt := events[0]
	assert.Equal(t, "entry.shared", string(appendedEvt.Type))
	assert.Equal(t, shareRefDTO.ShareEntryID, appendedEvt.Payload.ShareEntryID, "Appended event ShareEntryID MUST match created ShareEntry ID")
	assert.Equal(t, shareRefDTO.TrustGroupID, appendedEvt.Payload.TrustGroupID, "Appended event TrustGroupID MUST match created TrustGroup ID")
	assert.Equal(t, thread_domain.ResourceShareEntry, appendedEvt.Payload.RefType, "RefType MUST be share_entry")
}

func TestC3WritePath_MockTrustGroupID_Rejected_RealTrustGroup_Accepted(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	userAliceID := "user_alice_tg_test"
	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	// 1. Create a REAL TrustGroup via TrustGroup domain API
	realTG := trustgroup_domain.NewTrustGroup("ch_tg_real_test", "Real Persisted Group", []string{userAliceID})
	realTG.KEKVersion = 1
	repo.trustGroups[realTG.ID] = *realTG

	thread := thread_domain.NewThread("ch_tg_real_test", "contract", "TG Real Thread", "")
	_, err = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, nil, appendThreadEventUC)

	// 2. Attempt CreateCollaborativeShare with MOCK / non-existent TrustGroup IDs -> MUST FAIL
	mockTrustGroupIDs := []string{"tg_legal_counsel", "tg_finance", "mock_group_123", "ch_tg_real_test"}
	for _, mockID := range mockTrustGroupIDs {
		_, err := collabHandler.CreateCollaborativeShare(
			ctx,
			userAliceID,
			thread.ID,
			mockID,
			"bafybeicidmocktg",
			"v_target",
			"Attempt with mock TrustGroup",
			"d3JhcHBlZA==",
			1,
		)
		assert.Error(t, err, "Creation MUST fail when TrustGroupID '%s' is not an authoritative persisted TrustGroup", mockID)
	}

	// 3. Attempt CreateCollaborativeShare with REAL persisted TrustGroup ID -> MUST SUCCEED
	shareRefDTO, err := collabHandler.CreateCollaborativeShare(
		ctx,
		userAliceID,
		thread.ID,
		realTG.ID,
		"bafybeicidrealtg",
		"v_target",
		"Attempt with real TrustGroup",
		"d3JhcHBlZA==",
		1,
	)
	require.NoError(t, err)
	require.NotNil(t, shareRefDTO)
	assert.Equal(t, realTG.ID, shareRefDTO.TrustGroupID, "Returned TrustGroupID MUST match real persisted TrustGroup ID unchanged")
}

func TestC3EntryShared_CloudContract_POST_GET_RoundTrip(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	thread := thread_domain.NewThread("ch_cloud_contract", "contract", "Contract Event Roundtrip", "")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listUC := thread_usecase.NewListThreadEventsUsecase(repo)

	// 1. POST entry.shared event using exact payload shape: share_entry_id + trust_group_id
	const expectedShareID = "22081c47-29e5-4b03-bd1e-626d3f6f87e3"
	const expectedTrustGroupID = "ffc46329-6b01-4259-a101-22b6ccd48251"

	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: expectedShareID,
		TrustGroupID: expectedTrustGroupID,
	}

	appended, err := appendUC.Execute(ctx, thread.ID, "entry.shared", refPayload)
	require.NoError(t, err)
	require.NotNil(t, appended)

	// Verify appended event domain representation
	assert.Equal(t, thread_domain.ResourceShareEntry, appended.Payload.RefType)
	assert.Equal(t, expectedShareID, appended.Payload.ShareEntryID)
	assert.Equal(t, expectedTrustGroupID, appended.Payload.TrustGroupID)

	// 2. GET /ListThreadEvents -> verify authoritative persistence round-trip
	events, err := listUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, events, 1)

	fetchedEvent := events[0]
	assert.Equal(t, string(thread_domain.EventEntryShared), string(fetchedEvent.Type))
	assert.Equal(t, thread_domain.ResourceShareEntry, fetchedEvent.Payload.RefType)
	assert.Equal(t, expectedShareID, fetchedEvent.Payload.ShareEntryID, "share_entry_id MUST survive POST -> persistence -> GET round trip")
	assert.Equal(t, expectedTrustGroupID, fetchedEvent.Payload.TrustGroupID, "trust_group_id MUST survive POST -> persistence -> GET round trip")

	// 3. POST storage_asset event -> verify storage assets continue to use CID / ContentHash / Size
	storagePayload := thread_domain.EventResourceRef{
		RefType:     thread_domain.ResourceStorageAsset,
		CID:         "bafybeistorageasset2026",
		ContentHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Size:        2048,
		AssetType:   "document",
	}

	appendedStorage, err := appendUC.Execute(ctx, thread.ID, "asset.created", storagePayload)
	require.NoError(t, err)
	require.NotNil(t, appendedStorage)

	assert.Equal(t, thread_domain.ResourceStorageAsset, appendedStorage.Payload.RefType)
	assert.Equal(t, "bafybeistorageasset2026", appendedStorage.Payload.CID)
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", appendedStorage.Payload.ContentHash)
	assert.Equal(t, int64(2048), appendedStorage.Payload.Size)

	// Verify all events via GET
	allEvents, err := listUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, allEvents, 2)

	assert.Equal(t, expectedShareID, allEvents[0].Payload.ShareEntryID)
	assert.Equal(t, "bafybeistorageasset2026", allEvents[1].Payload.CID)
}
