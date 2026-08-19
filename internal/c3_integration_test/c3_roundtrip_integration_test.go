package c3_integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// In-Memory Test Infrastructure for End-to-End C3 Round-Trip
// ---------------------------------------------------------------------------

type roundTripRepo struct {
	shareEntries map[string]c3_asset_domain.ShareEntry
	trustGroups  map[string]trustgroup_domain.TrustGroup
	threads      map[string]thread_domain.Thread
	events       map[string][]thread_domain.ThreadEvent
	assets       map[string][]byte // AssetCID -> Encrypted Data
	seeds        map[string]string // UserID -> Device Seed
	keyrings     map[string]*vaults_domain.VaultKeyring
}

func newRoundTripRepo() *roundTripRepo {
	return &roundTripRepo{
		shareEntries: make(map[string]c3_asset_domain.ShareEntry),
		trustGroups:  make(map[string]trustgroup_domain.TrustGroup),
		threads:      make(map[string]thread_domain.Thread),
		events:       make(map[string][]thread_domain.ThreadEvent),
		assets:       make(map[string][]byte),
		seeds:        make(map[string]string),
		keyrings:     make(map[string]*vaults_domain.VaultKeyring),
	}
}

// ShareEntryRepository methods
func (r *roundTripRepo) CreateShareEntry(_ context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.shareEntries[req.ShareEntry.ID] = req.ShareEntry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry}, nil
}
func (r *roundTripRepo) GetShareEntry(_ context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	se, ok := r.shareEntries[req.ShareEntryID]
	if !ok {
		return nil, nil
	}
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: se}, nil
}
func (r *roundTripRepo) UpdateShareEntry(_ context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.shareEntries[req.ShareEntry.ID] = req.ShareEntry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry}, nil
}
func (r *roundTripRepo) DeleteShareEntry(_ context.Context, _ *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}

// TrustGroupRepository methods
func (r *roundTripRepo) GetTrustGroup(_ context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.trustGroups[req.TrustGroupID]
	if !ok {
		return nil, nil
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: tg}, nil
}
func (r *roundTripRepo) UpdateTrustGroup(_ context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.trustGroups[req.TrustGroup.ID] = req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}
func (r *roundTripRepo) CreateTrustGroup(_ context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.trustGroups[req.TrustGroup.ID] = req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}
func (r *roundTripRepo) GetTrustGroupMember(_ context.Context, _ *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (r *roundTripRepo) ListTrustGroups(_ context.Context, _ *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *roundTripRepo) DeleteTrustGroup(_ context.Context, _ *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *roundTripRepo) AddMemberToTrustGroup(_ context.Context, _ *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *roundTripRepo) RemoveMemberFromTrustGroup(_ context.Context, _ *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *roundTripRepo) RotateTrustGroupKEK(_ context.Context, _ *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

// ThreadRepository methods
func (r *roundTripRepo) CreateThread(_ context.Context, req *thread_domain.CreateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	r.threads[req.Thread.ID] = req.Thread
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: req.Thread}, nil
}
func (r *roundTripRepo) GetThread(_ context.Context, req *thread_domain.GetThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	th, ok := r.threads[req.ThreadID]
	if !ok {
		return nil, thread_domain.ErrThreadNotFound
	}
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: th}, nil
}
func (r *roundTripRepo) ListThreads(_ context.Context, _ *thread_domain.ListThreadsRequest) (*tracecore_types.CloudResponse[[]thread_domain.Thread], error) {
	return nil, nil
}
func (r *roundTripRepo) UpdateThread(_ context.Context, req *thread_domain.UpdateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	r.threads[req.Thread.ID] = req.Thread
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: req.Thread}, nil
}
func (r *roundTripRepo) ListThreadEvents(_ context.Context, req *thread_domain.ListThreadEventsRequest) (*tracecore_types.CloudResponse[[]thread_domain.ThreadEvent], error) {
	evts := r.events[req.ThreadID]
	return &tracecore_types.CloudResponse[[]thread_domain.ThreadEvent]{Data: evts}, nil
}
func (r *roundTripRepo) AppendThreadEvent(_ context.Context, req *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
	evt := thread_domain.ThreadEvent{
		ID:             "evt_" + uuid.NewString()[:8],
		ThreadID:       req.ThreadID,
		Type:           thread_domain.ThreadEventType(req.EventType),
		Payload:        req.Payload,
		IdempotencyKey: req.IdempotencyKey,
		Cursor:         uint64(len(r.events[req.ThreadID]) + 1),
		CreatedAt:      time.Now(),
	}
	r.events[req.ThreadID] = append(r.events[req.ThreadID], evt)
	return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{Data: evt}, nil
}

// Application Ports implementation
func (r *roundTripRepo) FetchEncryptedAsset(_ context.Context, cid string) ([]byte, error) {
	data, ok := r.assets[cid]
	if !ok {
		return nil, errors.New("asset CID not found in storage")
	}
	return data, nil
}

func (r *roundTripRepo) GetDeviceSeed(_ context.Context, userID string) (string, error) {
	seed, ok := r.seeds[userID]
	if !ok {
		return "", errors.New("seed not found for user")
	}
	return seed, nil
}

func (r *roundTripRepo) GetVaultKeyring(_ context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	kr, ok := r.keyrings[userID]
	if !ok {
		return &vaults_domain.VaultKeyring{UserID: userID, VaultID: "vault_rt"}, nil
	}
	return kr, nil
}

// ---------------------------------------------------------------------------
// End-to-End Vertical Slice Integration Test
// ---------------------------------------------------------------------------

func TestC3CollaborativeShare_WriteReadRoundTrip(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	// 1. Setup Alice's Sovereign Device Keypair & TrustGroup
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_uuid"
	deviceLaptopID := "device_laptop_m2"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "vault_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_rt_100", "Executive Board", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	// 2. Setup Thread Timeline
	thread := thread_domain.NewThread("ch_rt_100", "financial_plan", "Q4 Strategy", "Confidential")
	_, err = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	// 3. WRITE FLOW: Prepare & Encrypt Asset Payload (Simulating UI / Domain Write)
	rawOriginalContent := []byte(`{"title":"Q4 Growth Strategy","budget":5000000,"status":"APPROVED"}`)
	assetCID := "bafybeigq4strategy2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_q4_plan",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	require.Len(t, prepared.Envelopes, 1)

	// Persist encrypted asset payload into storage provider (AssetContentResolver)
	repo.assets[assetCID] = prepared.EncryptedData

	// Add key envelope to TrustGroup
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
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(repo)

	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, resolveCollabShareUC, appendThreadEventUC)

	// 4. WRITE EXECUTION: CreateCollaborativeShare -> Persist ShareEntry & Append ThreadEvent
	wrappedDEKStr := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

	shareReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   wrappedDEKStr,
		Envelopes:    prepared.Envelopes,
		Metadata:     map[string]string{"classification": "top_secret"},
	}

	createResp, err := createCollabShareUC.Execute(ctx, shareReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	createdShareEntryID := createResp.ShareEntry.ID
	repo.shareEntries[createdShareEntryID] = createResp.ShareEntry

	// Append Thread Event
	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createdShareEntryID,
		TrustGroupID: tg.ID,
	}
	_, err = appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", refPayload, "evt_share_"+createdShareEntryID)
	require.NoError(t, err)

	// 5. READ TIMELINE: Query Thread Events (Simulating UI opening Thread timeline)
	timelineEvents, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, timelineEvents, 1)

	threadEvt := timelineEvents[0]
	assert.Equal(t, thread_domain.ResourceShareEntry, threadEvt.Payload.RefType)
	assert.Equal(t, createdShareEntryID, threadEvt.Payload.ShareEntryID)
	assert.Equal(t, tg.ID, threadEvt.Payload.TrustGroupID)

	// SECURITY CHECK A: Verify Thread timeline contains NO plaintext or crypto secrets
	evtJsonBytes, err := json.Marshal(threadEvt)
	require.NoError(t, err)
	evtJsonStr := string(evtJsonBytes)
	assert.False(t, strings.Contains(evtJsonStr, "Q4 Growth Strategy"), "Thread timeline MUST NOT contain plaintext")
	assert.False(t, strings.Contains(evtJsonStr, "wrapped_dek"), "Thread timeline MUST NOT contain wrapped DEK")
	assert.False(t, strings.Contains(evtJsonStr, "wrapped_kek"), "Thread timeline MUST NOT contain wrapped KEK")

	// 6. READ EXECUTION: Resolve Collaborative Share via Handler/App Boundary
	resolvedDTO, err := collabHandler.ResolveCollaborativeShare(ctx, userAliceID, threadEvt.Payload.ShareEntryID, deviceLaptopID)
	require.NoError(t, err)
	require.NotNil(t, resolvedDTO)

	// 7. ASSERTIONS: Verify exact 100% Round-Trip Equality
	assert.Equal(t, createdShareEntryID, resolvedDTO.ShareEntryID)
	assert.Equal(t, tg.ID, resolvedDTO.TrustGroupID)
	assert.Equal(t, userAliceID, resolvedDTO.CreatedBy)
	assert.Equal(t, rawOriginalContent, resolvedDTO.Plaintext, "ROUND-TRIP INVARIANT: Decrypted plaintext MUST match original payload")

	// SECURITY CHECK B: Verify Resolved DTO contains NO cryptographic secrets when serialized
	dtoBytes, err := json.Marshal(resolvedDTO)
	require.NoError(t, err)
	dtoStr := string(dtoBytes)
	assert.False(t, strings.Contains(dtoStr, "wrapped_dek"))
	assert.False(t, strings.Contains(dtoStr, "wrapped_kek"))
	assert.False(t, strings.Contains(dtoStr, "device_seed"))

	// SECURITY CHECK C: Verify ShareEntry access descriptor contains NO plaintext
	seBytes, err := json.Marshal(createResp.ShareEntry)
	require.NoError(t, err)
	seStr := string(seBytes)
	assert.False(t, strings.Contains(seStr, "Q4 Growth Strategy"), "ShareEntry access descriptor MUST NOT store unencrypted plaintext")

	// SECURITY CHECK D: Verify storage provider contains ONLY encrypted ciphertext
	rawStoredCiphertext := repo.assets[assetCID]
	assert.NotEqual(t, rawOriginalContent, rawStoredCiphertext)
	assert.False(t, strings.Contains(string(rawStoredCiphertext), "Q4 Growth Strategy"), "Storage provider MUST store ciphertext only")
}
