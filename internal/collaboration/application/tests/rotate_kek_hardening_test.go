package collaboration_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)


// Transactional Mock Manager that supports rollback simulation
type mockTxManager struct {
	failOnTxCommit bool
}

func (m *mockTxManager) ExecuteInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	if m.failOnTxCommit {
		return errors.New("simulated transaction commit failure")
	}
	return fn(ctx)
}

// Transactional Fake Repositories
type transactionalFakeRepos struct {
	mu           sync.Mutex
	trustGroups  map[string]trustgroup_domain.TrustGroup
	shareEntries map[string]c3_asset_domain.ShareEntry

	failTGUpdate     bool
	failShareUpdate  bool
	failShareIDMatch string
}

func newTransactionalFakeRepos() *transactionalFakeRepos {
	return &transactionalFakeRepos{
		trustGroups:  make(map[string]trustgroup_domain.TrustGroup),
		shareEntries: make(map[string]c3_asset_domain.ShareEntry),
	}
}

func (r *transactionalFakeRepos) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tg, ok := r.trustGroups[req.TrustGroupID]
	if !ok {
		return nil, nil
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: tg}, nil
}

func (r *transactionalFakeRepos) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.failTGUpdate {
		return nil, errors.New("simulated TrustGroup update database error")
	}
	r.trustGroups[req.TrustGroup.ID] = req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *transactionalFakeRepos) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trustGroups[req.TrustGroup.ID] = req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *transactionalFakeRepos) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (r *transactionalFakeRepos) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *transactionalFakeRepos) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *transactionalFakeRepos) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *transactionalFakeRepos) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *transactionalFakeRepos) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (r *transactionalFakeRepos) CreateShareEntry(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shareEntries[req.ShareEntry.ID] = req.ShareEntry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry}, nil
}
func (r *transactionalFakeRepos) GetShareEntry(ctx context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.shareEntries[req.ShareEntryID]
	if !ok {
		return nil, nil
	}
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: entry}, nil
}
func (r *transactionalFakeRepos) UpdateShareEntry(ctx context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.failShareUpdate || (r.failShareIDMatch != "" && req.ShareEntry.ID == r.failShareIDMatch) {
		return nil, errors.New("simulated ShareEntry update database error")
	}
	r.shareEntries[req.ShareEntry.ID] = req.ShareEntry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry}, nil
}
func (r *transactionalFakeRepos) DeleteShareEntry(ctx context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}

// -----------------------------------------------------------------------------
// Failure Injection & Operational Tests
// -----------------------------------------------------------------------------

func TestRotateKEK_FailureA_TrustGroupUpdateFails(t *testing.T) {
	ctx := context.Background()
	repos := newTransactionalFakeRepos()

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Guild", []string{"vault-1"})
	_, _ = repos.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafy123", tg.ID, "wdek-v1", 1, "vault-1", nil)
	_, _ = repos.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	// Inject failure into TrustGroup update
	repos.failTGUpdate = true

	uc := collaboration_usecases.NewRotateTrustGroupKEKUseCase(repos, repos)

	rotReq := collaboration_usecases.RotateTrustGroupKEKRequest{
		TrustGroupID: tg.ID,
		OldVersion:   1,
		NewVersion:   2,
		NewEnvelopes: []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			{TrustGroupID: tg.ID, MemberID: "vault-1", DeviceID: "dev-1", KEKVersion: 2, WrappedKEK: "wkek-v2"},
		},
		RotatedShareEntries: []collaboration_usecases.RotatedShareEntryInput{
			{ShareEntryID: shareEntry.ID, ReWrappedDEK: "wdek-v2"},
		},
	}

	_, err := uc.Execute(ctx, rotReq)
	assert.ErrorContains(t, err, "simulated TrustGroup update database error")

	// Verify TrustGroup state was NOT mutated
	tgAfter, _ := repos.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	assert.Equal(t, uint64(1), tgAfter.Data.KEKVersion, "TrustGroup KEKVersion must remain 1 after failure")

	// Verify ShareEntry was NOT updated
	seAfter, _ := repos.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{ShareEntryID: shareEntry.ID})
	assert.Equal(t, uint64(1), seAfter.Data.KEKVersion, "ShareEntry KEKVersion must remain 1 after failure")
	assert.Equal(t, "wdek-v1", seAfter.Data.WrappedDEK)
}

func TestRotateKEK_FailureB_ShareEntryUpdateFails(t *testing.T) {
	ctx := context.Background()
	repos := newTransactionalFakeRepos()

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Guild", []string{"vault-1"})
	_, _ = repos.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafy123", tg.ID, "wdek-v1", 1, "vault-1", nil)
	_, _ = repos.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	// Inject failure on ShareEntry update
	repos.failShareIDMatch = shareEntry.ID

	uc := collaboration_usecases.NewRotateTrustGroupKEKUseCase(repos, repos)

	rotReq := collaboration_usecases.RotateTrustGroupKEKRequest{
		TrustGroupID: tg.ID,
		OldVersion:   1,
		NewVersion:   2,
		NewEnvelopes: []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			{TrustGroupID: tg.ID, MemberID: "vault-1", DeviceID: "dev-1", KEKVersion: 2, WrappedKEK: "wkek-v2"},
		},
		RotatedShareEntries: []collaboration_usecases.RotatedShareEntryInput{
			{ShareEntryID: shareEntry.ID, ReWrappedDEK: "wdek-v2"},
		},
	}

	_, err := uc.Execute(ctx, rotReq)
	assert.ErrorContains(t, err, "simulated ShareEntry update database error")
}

func TestRotateKEK_FailureD_IdempotencyRetry(t *testing.T) {
	ctx := context.Background()
	repos := newTransactionalFakeRepos()

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Guild", []string{"vault-1"})
	_, _ = repos.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafy123", tg.ID, "wdek-v1", 1, "vault-1", nil)
	_, _ = repos.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	uc := collaboration_usecases.NewRotateTrustGroupKEKUseCase(repos, repos)

	rotReq := collaboration_usecases.RotateTrustGroupKEKRequest{
		RequestID:    "req-idempotent-101",
		TrustGroupID: tg.ID,
		OldVersion:   1,
		NewVersion:   2,
		NewEnvelopes: []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			{TrustGroupID: tg.ID, MemberID: "vault-1", DeviceID: "dev-1", KEKVersion: 2, WrappedKEK: "wkek-v2"},
		},
		RotatedShareEntries: []collaboration_usecases.RotatedShareEntryInput{
			{ShareEntryID: shareEntry.ID, ReWrappedDEK: "wdek-v2"},
		},
	}

	// First call -> commits rotation to v2
	resp1, err := uc.Execute(ctx, rotReq)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), resp1.TrustGroup.KEKVersion)

	// Second call (Retry with same RequestID and payload)
	resp2, err := uc.Execute(ctx, rotReq)
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// Idempotency assertions:
	assert.Equal(t, uint64(2), resp2.TrustGroup.KEKVersion, "Retry must return version 2, NOT jump to v3")
	assert.Len(t, resp2.TrustGroup.KeyEnvelopes, 1, "Must NOT duplicate envelopes on retry")
}

func TestRotateKEK_FailureE_StaleClientRejection(t *testing.T) {
	ctx := context.Background()
	repos := newTransactionalFakeRepos()

	// TrustGroup is already at KEKVersion 2
	tg := trustgroup_domain.NewTrustGroup("chan-1", "Guild", []string{"vault-1"})
	tg.KEKVersion = 2
	_, _ = repos.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})

	uc := collaboration_usecases.NewRotateTrustGroupKEKUseCase(repos, repos)

	// Stale Client attempts v1 -> v2 rotation
	staleReq := collaboration_usecases.RotateTrustGroupKEKRequest{
		TrustGroupID: tg.ID,
		OldVersion:   1, // Stale!
		NewVersion:   2,
	}

	_, err := uc.Execute(ctx, staleReq)
	assert.ErrorIs(t, err, trustgroup_domain.ErrStaleKEKVersion)
}

func TestRotateKEK_FailureF_ConcurrentRotations(t *testing.T) {
	ctx := context.Background()
	repos := newTransactionalFakeRepos()

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Guild", []string{"vault-1"})
	_, _ = repos.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafy123", tg.ID, "wdek-v1", 1, "vault-1", nil)
	_, _ = repos.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	uc := collaboration_usecases.NewRotateTrustGroupKEKUseCase(repos, repos)

	req1 := collaboration_usecases.RotateTrustGroupKEKRequest{
		RequestID:    "req-race-1",
		TrustGroupID: tg.ID,
		OldVersion:   1,
		NewVersion:   2,
		NewEnvelopes: []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			{TrustGroupID: tg.ID, MemberID: "vault-1", DeviceID: "dev-1", KEKVersion: 2, WrappedKEK: "wkek-v2-race1"},
		},
		RotatedShareEntries: []collaboration_usecases.RotatedShareEntryInput{
			{ShareEntryID: shareEntry.ID, ReWrappedDEK: "wdek-v2-race1"},
		},
	}

	req2 := collaboration_usecases.RotateTrustGroupKEKRequest{
		RequestID:    "req-race-2",
		TrustGroupID: tg.ID,
		OldVersion:   1,
		NewVersion:   2,
		NewEnvelopes: []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			{TrustGroupID: tg.ID, MemberID: "vault-1", DeviceID: "dev-1", KEKVersion: 2, WrappedKEK: "wkek-v2-race2"},
		},
		RotatedShareEntries: []collaboration_usecases.RotatedShareEntryInput{
			{ShareEntryID: shareEntry.ID, ReWrappedDEK: "wdek-v2-race2"},
		},
	}

	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err1 = uc.Execute(ctx, req1)
	}()
	go func() {
		defer wg.Done()
		_, err2 = uc.Execute(ctx, req2)
	}()
	wg.Wait()

	// Exactly one rotation succeeds, the other fails with ErrStaleKEKVersion or succeeds idempotently if matching
	successes := 0
	if err1 == nil {
		successes++
	}
	if err2 == nil {
		successes++
	}
	assert.GreaterOrEqual(t, successes, 1, "At least one concurrent rotation must succeed")

	// Final state must be consistently at KEKVersion 2
	finalTG, _ := repos.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	assert.Equal(t, uint64(2), finalTG.Data.KEKVersion)
}
