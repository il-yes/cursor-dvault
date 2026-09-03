package c3_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_usecases "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_member_usecases "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
)

type memoryInvitationOrchestrationRepo struct {
	groups      map[string]*trustgroup_domain.TrustGroup
	invitations map[string]*channel_domain.Invitation
	shares      map[string]*c3_asset_domain.ShareEntry
}

func newMemoryInvitationOrchestrationRepo() *memoryInvitationOrchestrationRepo {
	return &memoryInvitationOrchestrationRepo{
		groups:      make(map[string]*trustgroup_domain.TrustGroup),
		invitations: make(map[string]*channel_domain.Invitation),
		shares:      make(map[string]*c3_asset_domain.ShareEntry),
	}
}

func (r *memoryInvitationOrchestrationRepo) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryInvitationOrchestrationRepo) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	var matched []trustgroup_domain.TrustGroup
	for _, tg := range r.groups {
		if req.ChannelID == "" || tg.ChannelID == req.ChannelID {
			matched = append(matched, *tg)
		}
	}
	return &tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup]{Data: matched}, nil
}

func (r *memoryInvitationOrchestrationRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	alreadyFound := false
	for _, cid := range tg.MemberCIDs {
		if cid == req.VaultID {
			alreadyFound = true
			break
		}
	}
	if !alreadyFound {
		tg.MemberCIDs = append(tg.MemberCIDs, req.VaultID)
	}
	r.groups[req.TrustGroupID] = tg
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryInvitationOrchestrationRepo) AddTrustGroupKeyEnvelope(ctx context.Context, req *trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	newEnv := trustgroup_domain.TrustGroupKeyEnvelope{
		ID:         "env_" + req.MemberID + "_" + req.DeviceID,
		MemberID:   req.MemberID,
		DeviceID:   req.DeviceID,
		KEKVersion: req.KEKVersion,
		WrappedKEK: req.WrappedKEK,
		CreatedAt:  time.Now(),
	}
	tg.KeyEnvelopes = append(tg.KeyEnvelopes, newEnv)
	r.groups[req.TrustGroupID] = tg
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryInvitationOrchestrationRepo) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}
func (r *memoryInvitationOrchestrationRepo) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (r *memoryInvitationOrchestrationRepo) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}
func (r *memoryInvitationOrchestrationRepo) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryInvitationOrchestrationRepo) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryInvitationOrchestrationRepo) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func TestInvitationAcceptanceOrchestration_FullLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryInvitationOrchestrationRepo()

	userA_VaultID := "vault_alice_777"
	userB_VaultID := "vault_bob_888"
	deviceB_ID := "dev_bob_desktop_01"

	kpBob, err := keypair.Random()
	require.NoError(t, err)

	channelID := "ch_legal_dept"
	tgID := "tg_legal_secure"

	// 1. User A creates TrustGroup & Channel
	tg := trustgroup_domain.NewTrustGroup(channelID, "Legal TrustGroup", []string{userA_VaultID})
	tg.ID = tgID
	_, err = repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	mockDevResolver := &mockDeviceResolver{
		devices: map[string]*trustgroup_ports.DeviceSummary{
			deviceB_ID: {ID: deviceB_ID, VaultID: userB_VaultID, PublicKey: kpBob.Address(), IsActive: true},
		},
	}

	// Setup usecases
	addMemberUC := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(repo, trustgroup_eventbus.NewMemoryBus())
	addEnvelopeUC := trustgroup_usecases.NewAddTrustGroupKeyEnvelopeUseCase(repo, mockDevResolver)

	provisionEnvelopeUC := trustgroup_usecases.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		repo,
		mockDevResolver,
		nil,
		addEnvelopeUC,
		nil,
	)

	// TEST 2: Before acceptance, User B is NOT in TrustGroup
	tgBefore, err := repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	assert.NotContains(t, tgBefore.Data.MemberCIDs, userB_VaultID, "User B must not be in TrustGroup before acceptance")

	// TEST 4: Execute Invitation Acceptance Orchestration (Step 3 + Step 4)
	// 3. Add Member
	_, err = addMemberUC.Execute(ctx, trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: tgID,
		VaultID:      userB_VaultID,
		Role:         "member",
	})
	require.NoError(t, err)

	// TEST 3: User B is member, but envelope not yet attached
	tgMid, err := repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	assert.Contains(t, tgMid.Data.MemberCIDs, userB_VaultID, "User B must be in MemberCIDs after member addition")
	assert.Empty(t, tgMid.Data.KeyEnvelopes, "KeyEnvelopes must be empty before envelope provisioning")

	// 4. Provision Envelope
	_, err = provisionEnvelopeUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tgID,
		MemberID:        userB_VaultID,
		DeviceID:        deviceB_ID,
		DevicePublicKey: kpBob.Address(),
	}, nil)
	require.NoError(t, err)

	// TEST 1: Final verification of full acceptance pipeline
	tgAfter, err := repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	assert.Contains(t, tgAfter.Data.MemberCIDs, userB_VaultID, "User B VaultID must be in MemberCIDs")
	require.Len(t, tgAfter.Data.KeyEnvelopes, 1, "User B must have exactly 1 device key envelope")
	assert.Equal(t, userB_VaultID, tgAfter.Data.KeyEnvelopes[0].MemberID)
	assert.Equal(t, deviceB_ID, tgAfter.Data.KeyEnvelopes[0].DeviceID)

	// TEST 5: Idempotency check - re-running provisioning must skip cleanly without error or duplicates
	_, err = provisionEnvelopeUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tgID,
		MemberID:        userB_VaultID,
		DeviceID:        deviceB_ID,
		DevicePublicKey: kpBob.Address(),
	}, nil)
	require.NoError(t, err)

	tgIdempotent, err := repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	assert.Len(t, tgIdempotent.Data.KeyEnvelopes, 1, "KeyEnvelopes must remain exactly 1 after idempotent retry")
}

type mockDeviceResolver struct {
	devices map[string]*trustgroup_ports.DeviceSummary
}

func (m *mockDeviceResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	if dev, ok := m.devices[deviceID]; ok {
		return dev, nil
	}
	return nil, nil
}
func (m *mockDeviceResolver) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	var list []trustgroup_ports.DeviceSummary
	for _, dev := range m.devices {
		if dev.VaultID == memberID && dev.IsActive {
			list = append(list, *dev)
		}
	}
	return list, nil
}
