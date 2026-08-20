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
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Phase 1: ShareEntry Lifecycle Tests
// ---------------------------------------------------------------------------

func TestC3ShareEntryLifecycle_Revoked_HistoricalThreadPreserved(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_uuid"
	deviceLaptopID := "device_laptop_m2"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "vault_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_life_1", "Executive Board", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	thread := thread_domain.NewThread("ch_life_1", "contract", "Merger Terms", "Draft v1")
	_, err = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	rawOriginalContent := []byte(`{"terms":"Confidential Merger Agreement"}`)
	assetCID := "bafybeimergerterms2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_merger",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
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

	wrappedDEKStr := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

	shareReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   wrappedDEKStr,
		Envelopes:    prepared.Envelopes,
		Metadata:     map[string]string{"title": "Merger Agreement"},
	}

	createResp, err := createCollabShareUC.Execute(ctx, shareReq)
	require.NoError(t, err)

	createdShareEntryID := createResp.ShareEntry.ID
	repo.shareEntries[createdShareEntryID] = createResp.ShareEntry

	// Append Thread Event
	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createdShareEntryID,
		TrustGroupID: tg.ID,
	}
	originalEvt, err := appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", refPayload, "evt_share_"+createdShareEntryID)
	require.NoError(t, err)

	// Verify ACTIVE state resolution works
	activeDTO, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createdShareEntryID,
		CallerUserID: userAliceID,
		DeviceID:     deviceLaptopID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawOriginalContent, activeDTO.Plaintext)

	// --- MUTATION: REVOKE ShareEntry ---
	shareEntry := repo.shareEntries[createdShareEntryID]
	shareEntry.Status = c3_asset_domain.ShareEntryStatusRevoked
	repo.shareEntries[createdShareEntryID] = shareEntry

	// INVARIANT 1: ThreadEvent remains byte-for-byte UNCHANGED in timeline
	timelineEvents, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, timelineEvents, 1)

	currentEvt := timelineEvents[0]
	assert.Equal(t, originalEvt.ID, currentEvt.ID, "ThreadEvent ID must be unchanged")
	assert.Equal(t, originalEvt.Payload.RefType, currentEvt.Payload.RefType, "ResourceRef RefType must be unchanged")
	assert.Equal(t, originalEvt.Payload.ShareEntryID, currentEvt.Payload.ShareEntryID, "ResourceRef ShareEntryID must be unchanged")
	assert.Equal(t, originalEvt.Payload.TrustGroupID, currentEvt.Payload.TrustGroupID, "ResourceRef TrustGroupID must be unchanged")

	// INVARIANT 2: Resolution of Revoked ShareEntry is DENIED immediately
	_, err = resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createdShareEntryID,
		CallerUserID: userAliceID,
		DeviceID:     deviceLaptopID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryRevoked, "Accessing revoked ShareEntry must return ErrShareEntryRevoked")
}

func TestC3ShareEntryLifecycle_HardDeleted_HistoricalThreadPreserved(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_uuid"
	deviceLaptopID := "device_laptop_m2"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "vault_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_life_2", "Board", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	thread := thread_domain.NewThread("ch_life_2", "note", "Deleted Note", "")
	_, err = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	rawOriginalContent := []byte(`{"note":"Temporary Note"}`)
	assetCID := "bafybeitempnote2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_temp_note",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
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

	wrappedDEKStr := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

	shareReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   wrappedDEKStr,
		Envelopes:    prepared.Envelopes,
		Metadata:     map[string]string{"title": "Temporary Note"},
	}

	createResp, err := createCollabShareUC.Execute(ctx, shareReq)
	require.NoError(t, err)

	createdShareEntryID := createResp.ShareEntry.ID
	repo.shareEntries[createdShareEntryID] = createResp.ShareEntry

	// Append Thread Event
	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createdShareEntryID,
		TrustGroupID: tg.ID,
	}
	originalEvt, err := appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", refPayload, "evt_share_"+createdShareEntryID)
	require.NoError(t, err)

	// --- MUTATION: HARD DELETE ShareEntry ---
	delete(repo.shareEntries, createdShareEntryID)

	// INVARIANT 1: ThreadEvent remains byte-for-byte UNCHANGED in timeline
	timelineEvents, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, timelineEvents, 1)

	currentEvt := timelineEvents[0]
	assert.Equal(t, originalEvt.ID, currentEvt.ID, "ThreadEvent ID must be unchanged")
	assert.Equal(t, originalEvt.Payload.RefType, currentEvt.Payload.RefType, "ResourceRef RefType must be unchanged")
	assert.Equal(t, originalEvt.Payload.ShareEntryID, currentEvt.Payload.ShareEntryID, "ResourceRef ShareEntryID must be unchanged")

	// INVARIANT 2: Resolution of Hard Deleted ShareEntry returns ErrShareEntryNotFound
	_, err = resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createdShareEntryID,
		CallerUserID: userAliceID,
		DeviceID:     deviceLaptopID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryNotFound, "Accessing deleted ShareEntry must return ErrShareEntryNotFound")
}

func TestC3ShareEntryLifecycle_RevokedBeforeTrustGroupOrCryptoResolution(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	shareEntry := c3_asset_domain.ShareEntry{
		ID:           "se_revoked_early",
		AssetCID:     "cid_never_fetched",
		TrustGroupID: "tg_never_queried",
		WrappedDEK:   "dek_never_unwrapped",
		KEKVersion:   1,
		CreatedBy:    "user_alice",
		Status:       c3_asset_domain.ShareEntryStatusRevoked,
	}

	repo.shareEntries[shareEntry.ID] = shareEntry

	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	_, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: "user_alice",
		DeviceID:     "dev_laptop",
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryRevoked)
}
