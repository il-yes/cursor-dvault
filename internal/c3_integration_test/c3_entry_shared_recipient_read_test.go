package c3_integration_test

import (
	"context"
	"encoding/base64"
	"testing"

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
		Envelopes:    prepared.Envelopes,
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
	resolvedDTO, err := collabHandler.ResolveCollaborativeShare(ctx, userBobID, extractedShareEntryID, deviceBobID)
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
		Envelopes:    prepared.Envelopes,
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
	resolvedDTO, err := collabHandler.ResolveCollaborativeShare(ctx, userCharlieID, createdShareEntryID, deviceCharlieID)
	assert.Error(t, err, "Decryption resolution MUST fail for unauthorized recipient")
	assert.Nil(t, resolvedDTO, "NO plaintext must be returned to unauthorized recipient")
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember, "Must return ErrUnauthorizedMember")
}
