package c3_integration_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type memoryPersistenceTestRepo struct {
	groups map[string]*trustgroup_domain.TrustGroup
	shares map[string]*c3_asset_domain.ShareEntry
}

func newMemoryPersistenceTestRepo() *memoryPersistenceTestRepo {
	return &memoryPersistenceTestRepo{
		groups: make(map[string]*trustgroup_domain.TrustGroup),
		shares: make(map[string]*c3_asset_domain.ShareEntry),
	}
}

func (r *memoryPersistenceTestRepo) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *memoryPersistenceTestRepo) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryPersistenceTestRepo) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *memoryPersistenceTestRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	tg.MemberCIDs = append(tg.MemberCIDs, req.VaultID)
	r.groups[req.TrustGroupID] = tg
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryPersistenceTestRepo) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (r *memoryPersistenceTestRepo) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryPersistenceTestRepo) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryPersistenceTestRepo) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryPersistenceTestRepo) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func TestTrustGroup_MembershipPersistenceAndReload(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryPersistenceTestRepo()

	userA_VaultID := "vault_alice_101"
	userB_VaultID := "vault_bob_202"

	// 1. Create TrustGroup with User A
	tg := trustgroup_domain.NewTrustGroup("ch_legal_01", "Legal Team", []string{userA_VaultID})
	tg.ID = "b72ace04-852f-493a-877c-c9938fc0a8cc"
	_, err := repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// 2. Add User B to TrustGroup
	addResp, err := repo.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: tg.ID,
		VaultID:      userB_VaultID,
		Role:         "member",
	})
	require.NoError(t, err)
	require.NotNil(t, addResp)
	assert.Contains(t, addResp.Data.MemberCIDs, userB_VaultID)

	// 3. Reload TrustGroup through SAME repository path used by ResolveCollaborativeShare
	loadedResp, err := repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	require.NoError(t, err)
	require.NotNil(t, loadedResp)

	// 4. Assert User B's VaultID remains in MemberCIDs
	assert.Contains(t, loadedResp.Data.MemberCIDs, userB_VaultID, "User B VaultID must remain in MemberCIDs after persistence & reload")

	// 5. Assert ShareEntry TrustGroupID matches the updated TrustGroup ID
	shareEntry := c3_asset_domain.ShareEntry{
		ID:           "964decad-98d6-43a7-b552-78d06158589f",
		TrustGroupID: "b72ace04-852f-493a-877c-c9938fc0a8cc",
	}
	assert.Equal(t, loadedResp.Data.ID, shareEntry.TrustGroupID, "ShareEntry.TrustGroupID must match the updated TrustGroup ID")
}

func TestAddMemberToTrustGroup_JSONSerializationContract(t *testing.T) {
	req := trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: "b72ace04-852f-493a-877c-c9938fc0a8cc",
		VaultID:      "31e1866f-f82b-469e-9382-96fb9f311a9a",
		Role:         "member",
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	bodyStr := string(body)
	assert.Contains(t, bodyStr, `"vault_id":`)
	assert.Contains(t, bodyStr, `"role":`)
	assert.NotContains(t, bodyStr, `"trust_group_id":`)
	assert.NotContains(t, bodyStr, `"member_id":`)
}
