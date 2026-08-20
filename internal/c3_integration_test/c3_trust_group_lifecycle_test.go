package c3_integration_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

// ---------------------------------------------------------------------------
// Phase 2: TrustGroup Membership & Key Lifecycle Test Matrix
// ---------------------------------------------------------------------------

// 1. Member Revoked: Historical ThreadEvent remains, but resolution is DENIED
func TestC3TrustGroupLifecycle_MemberRevoked_ResolutionDenied(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAlice, _ := keypair.Random()
	userAliceID := "user_alice_rev"
	deviceAliceLaptop := "dev_alice_laptop"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_tg_1", "Council", []string{userAliceID})
	tg.KEKVersion = 1

	thread := thread_domain.NewThread("ch_tg_1", "doc", "Council Minutes", "")
	_, _ = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})

	rawContent := []byte(`{"minutes":"Confidential Council Minutes"}`)
	assetCID := "bafybeicouncilminutes2026"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_minutes",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceLaptop, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
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

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(repo)

	createResp, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
		Envelopes:    prepared.Envelopes,
	})
	require.NoError(t, err)

	repo.shareEntries[createResp.ShareEntry.ID] = createResp.ShareEntry

	originalEvt, err := appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createResp.ShareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_share_"+createResp.ShareEntry.ID)
	require.NoError(t, err)

	// Verify Alice can resolve BEFORE revocation
	res, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createResp.ShareEntry.ID,
		CallerUserID: userAliceID,
		DeviceID:     deviceAliceLaptop,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, res.Plaintext)

	// --- MUTATION: REVOKE ALICE FROM TRUSTGROUP ---
	tg.MemberCIDs = []string{} // Alice removed from MemberCIDs
	repo.trustGroups[tg.ID] = *tg

	// INVARIANT 1: Historical ThreadEvent remains byte-for-byte unchanged in timeline
	evts, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, evts, 1)
	assert.Equal(t, originalEvt.ID, evts[0].ID)

	// INVARIANT 2: Alice's resolution attempt is DENIED with ErrUnauthorizedMember
	_, err = resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createResp.ShareEntry.ID,
		CallerUserID: userAliceID,
		DeviceID:     deviceAliceLaptop,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
}

// 2. Device Revoked vs Active Device for Same Member
func TestC3TrustGroupLifecycle_DeviceRevoked_DistinctFromMemberIdentity(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAliceLaptop, _ := keypair.Random()
	kpAliceDesktop, _ := keypair.Random()

	userAliceID := "user_alice_multi_device"
	deviceLaptopID := "dev_laptop"
	deviceDesktopID := "dev_desktop"

	repo.seeds[userAliceID] = kpAliceLaptop.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_tg_2", "Board", []string{userAliceID})
	tg.KEKVersion = 1

	rawContent := []byte(`{"data":"Multi-device payload"}`)
	assetCID := "bafybeimultidevice2026"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_multi",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAliceLaptop.Address(), IsActive: true},
			{DeviceID: deviceDesktopID, MemberID: userAliceID, PublicKey: kpAliceDesktop.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	repo.assets[assetCID] = prepared.EncryptedData

	for _, envReq := range prepared.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	repo.trustGroups[tg.ID] = *tg

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	createResp, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
		Envelopes:    prepared.Envelopes,
	})
	require.NoError(t, err)
	repo.shareEntries[createResp.ShareEntry.ID] = createResp.ShareEntry

	// --- MUTATION: REVOKE LAPTOP ENVELOPE ONLY ---
	now := time.Now()
	for i := range tg.KeyEnvelopes {
		if tg.KeyEnvelopes[i].DeviceID == deviceLaptopID {
			tg.KeyEnvelopes[i].RevokedAt = &now
		}
	}
	repo.trustGroups[tg.ID] = *tg

	// Revoked Laptop Device -> DENIED
	_, err = resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createResp.ShareEntry.ID,
		CallerUserID: userAliceID,
		DeviceID:     deviceLaptopID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound, "Revoked device envelope must return ErrKeyEnvelopeNotFound")

	// Active Desktop Device -> ALLOWED
	repo.seeds[userAliceID] = kpAliceDesktop.Seed() // Switch seed to Desktop keypair
	resDesktop, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createResp.ShareEntry.ID,
		CallerUserID: userAliceID,
		DeviceID:     deviceDesktopID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, resDesktop.Plaintext, "Active device for authorized member MUST succeed")
}

// 3. Existing ShareEntry Readable After Full KEK Rotation & DEK Re-wrapping (Option A Proof)
func TestC3TrustGroupLifecycle_ExistingShareEntry_ReadableAfterKEKRotation(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "/tmp/keyz_test_rot", nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	kpBob, _ := keypair.Random()
	userBobID := "user_bob_rot_proof"
	deviceBobID := "dev_bob_m1"

	repo.seeds[userBobID] = kpBob.Seed()
	bobKeyring := &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}
	repo.keyrings[userBobID] = bobKeyring

	tg := trustgroup_domain.NewTrustGroup("ch_rot_proof", "Security Team", []string{userBobID})
	tg.KEKVersion = 1

	rawContent := []byte(`{"content":"Confidential Asset to survive KEK rotation"}`)
	assetCID := "bafybeirotateproof2026"

	// Prepare V1 Asset
	prepV1 := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_rot_proof",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
		Keyring: bobKeyring,
	}

	preparedV1, err := orchestrator.PrepareCollaborativeAsset(ctx, prepV1)
	require.NoError(t, err)
	repo.assets[assetCID] = preparedV1.EncryptedData

	for _, envReq := range preparedV1.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	repo.trustGroups[tg.ID] = *tg

	shareEntry, _ := c3_asset_domain.NewShareEntry(assetCID, tg.ID, base64.StdEncoding.EncodeToString(preparedV1.WrappedDEK), 1, userBobID, nil)
	shareEntry.ID = "se_rot_proof_100"
	repo.shareEntries[shareEntry.ID] = shareEntry

	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)
	rotateKEKUC := collaboration_usecases.NewRotateTrustGroupKEKUseCase(repo, repo)

	// Verify Bob can resolve under V1 BEFORE rotation
	resV1, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceBobID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, resV1.Plaintext)

	// --- EXECUTE FULL KEK ROTATION TO V2 WITH DEK RE-WRAPPING ---
	rotPayload := trustgroup_orchestrator.RotateTrustGroupKEKPayload{
		TrustGroupID: tg.ID,
		OldVersion:   1,
		NewVersion:   2,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
		Assets: []trustgroup_orchestrator.RotateCollaborativeAssetInput{
			{ShareEntryID: shareEntry.ID, WrappedDEK: preparedV1.WrappedDEK},
		},
		Keyring: bobKeyring,
	}

	rotCryptoResult, err := orchestrator.RotateTrustGroupKEK(ctx, rotPayload)
	require.NoError(t, err)
	require.Len(t, rotCryptoResult.RotatedAssets, 1)

	// Persist KEK v2 rotation & re-wrapped DEKs via RotateTrustGroupKEKUseCase
	var rotatedInputs []collaboration_usecases.RotatedShareEntryInput
	for _, ra := range rotCryptoResult.RotatedAssets {
		rotatedInputs = append(rotatedInputs, collaboration_usecases.RotatedShareEntryInput{
			ShareEntryID: ra.ShareEntryID,
			ReWrappedDEK: base64.StdEncoding.EncodeToString(ra.ReWrappedDEK),
		})
	}

	_, err = rotateKEKUC.Execute(ctx, collaboration_usecases.RotateTrustGroupKEKRequest{
		RequestID:           "rot_req_001",
		TrustGroupID:        tg.ID,
		OldVersion:          1,
		NewVersion:          2,
		NewEnvelopes:        rotCryptoResult.NewEnvelopes,
		RotatedShareEntries: rotatedInputs,
	})
	require.NoError(t, err)

	// INVARIANT VERIFICATION: AssetCID payload ciphertext in storage remains UNTOUCHED
	assert.Equal(t, preparedV1.EncryptedData, repo.assets[assetCID], "Raw ciphertext in storage MUST remain untouched")

	// INVARIANT VERIFICATION: ShareEntry.WrappedDEK and KEKVersion updated to V2
	updatedShare := repo.shareEntries[shareEntry.ID]
	assert.Equal(t, uint64(2), updatedShare.KEKVersion, "ShareEntry.KEKVersion MUST be updated to 2")
	assert.NotEqual(t, base64.StdEncoding.EncodeToString(preparedV1.WrappedDEK), updatedShare.WrappedDEK, "WrappedDEK MUST be re-wrapped")

	// INVARIANT VERIFICATION: Bob can resolve existing ShareEntry under KEK v2
	resV2, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceBobID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, resV2.Plaintext, "Existing ShareEntry MUST remain readable after KEK rotation & DEK re-wrapping")
}

// 4. Revoked Member Cannot Read Historical ShareEntry After KEK Rotation
func TestC3TrustGroupLifecycle_RevokedMember_CannotReadAfterKEKRotation(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "/tmp/keyz_test_rot_rev", nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	kpAlice, _ := keypair.Random()
	kpBob, _ := keypair.Random()

	userAliceID := "user_alice_rev_proof"
	userBobID := "user_bob_rem_proof"
	deviceAliceID := "dev_alice_m1"
	deviceBobID := "dev_bob_m1"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.seeds[userBobID] = kpBob.Seed()
	bobKeyring := &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}
	repo.keyrings[userBobID] = bobKeyring

	tg := trustgroup_domain.NewTrustGroup("ch_rot_rev_proof", "Guild", []string{userAliceID, userBobID})
	tg.KEKVersion = 1

	rawContent := []byte(`{"content":"Guild Vault Strategy V1"}`)
	assetCID := "bafybeiguildstrategy2026"

	prepV1 := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_guild_strategy",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
		Keyring: bobKeyring,
	}

	preparedV1, err := orchestrator.PrepareCollaborativeAsset(ctx, prepV1)
	require.NoError(t, err)
	repo.assets[assetCID] = preparedV1.EncryptedData

	for _, envReq := range preparedV1.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	repo.trustGroups[tg.ID] = *tg

	shareEntry, _ := c3_asset_domain.NewShareEntry(assetCID, tg.ID, base64.StdEncoding.EncodeToString(preparedV1.WrappedDEK), 1, userAliceID, nil)
	shareEntry.ID = "se_guild_strat_100"
	repo.shareEntries[shareEntry.ID] = shareEntry

	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)
	rotateKEKUC := collaboration_usecases.NewRotateTrustGroupKEKUseCase(repo, repo)

	// --- EXECUTE KEK ROTATION TO V2 WHILE REVOKING ALICE ---
	rotPayload := trustgroup_orchestrator.RotateTrustGroupKEKPayload{
		TrustGroupID: tg.ID,
		OldVersion:   1,
		NewVersion:   2,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
			// Alice is excluded (IsActive == false or omitted)
		},
		Assets: []trustgroup_orchestrator.RotateCollaborativeAssetInput{
			{ShareEntryID: shareEntry.ID, WrappedDEK: preparedV1.WrappedDEK},
		},
		Keyring: bobKeyring,
	}

	rotCryptoResult, err := orchestrator.RotateTrustGroupKEK(ctx, rotPayload)
	require.NoError(t, err)

	var rotatedInputs []collaboration_usecases.RotatedShareEntryInput
	for _, ra := range rotCryptoResult.RotatedAssets {
		rotatedInputs = append(rotatedInputs, collaboration_usecases.RotatedShareEntryInput{
			ShareEntryID: ra.ShareEntryID,
			ReWrappedDEK: base64.StdEncoding.EncodeToString(ra.ReWrappedDEK),
		})
	}

	_, err = rotateKEKUC.Execute(ctx, collaboration_usecases.RotateTrustGroupKEKRequest{
		RequestID:           "rot_req_rev_002",
		TrustGroupID:        tg.ID,
		OldVersion:          1,
		NewVersion:          2,
		RevokedMemberID:     userAliceID, // Remove Alice from TrustGroup.MemberCIDs
		NewEnvelopes:        rotCryptoResult.NewEnvelopes,
		RotatedShareEntries: rotatedInputs,
	})
	require.NoError(t, err)

	// ASSERTION 1: Revoked Alice attempting to resolve migrated ShareEntry -> DENIED (ErrUnauthorizedMember)
	_, err = resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userAliceID,
		DeviceID:     deviceAliceID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember, "Revoked member MUST be denied access after rotation")

	// ASSERTION 2: Remaining Member Bob resolving migrated ShareEntry -> SUCCESS
	resBob, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceBobID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, resBob.Plaintext, "Remaining authorized member MUST be able to resolve re-wrapped ShareEntry")
}
