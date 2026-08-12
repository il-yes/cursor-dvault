package trustgroup_usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"


	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_usecases "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

// Fake DeviceResolver for tests
type fakeDeviceResolver struct {
	devices map[string]*trustgroup_ports.DeviceSummary
}

func newFakeDeviceResolver() *fakeDeviceResolver {
	return &fakeDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
}

func (r *fakeDeviceResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	d, ok := r.devices[deviceID]
	if !ok {
		return nil, nil
	}
	return d, nil
}

// Fake TrustGroupRepository for tests
type fakeTrustGroupRepo struct {
	groups map[string]*trustgroup_domain.TrustGroup
}

func newFakeTrustGroupRepo() *fakeTrustGroupRepo {
	return &fakeTrustGroupRepo{groups: make(map[string]*trustgroup_domain.TrustGroup)}
}

func (r *fakeTrustGroupRepo) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *fakeTrustGroupRepo) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, nil
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *fakeTrustGroupRepo) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *fakeTrustGroupRepo) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func TestAddTrustGroupKeyEnvelope_Success(t *testing.T) {
	ctx := context.Background()
	repo := newFakeTrustGroupRepo()
	resolver := newFakeDeviceResolver()
	uc := trustgroup_usecases.NewAddTrustGroupKeyEnvelopeUseCase(repo, resolver)

	// Setup TrustGroup
	tg := trustgroup_domain.NewTrustGroup("chan-1", "Eng Team", []string{"member-vault-1"})
	_, err := repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// Setup Device
	resolver.devices["dev-1"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-1",
		VaultID:   "member-vault-1",
		PublicKey: "pubkey-1",
		Status:    "active",
		IsActive:  true,
	}

	req := trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
		TrustGroupID: tg.ID,
		MemberID:     "member-vault-1",
		DeviceID:     "dev-1",
		KEKVersion:   1,
		WrappedKEK:   "opaque-wrapped-kek-bytes",
	}

	updatedTg, err := uc.Execute(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, updatedTg)
	assert.Len(t, updatedTg.KeyEnvelopes, 1)
	assert.Equal(t, "dev-1", updatedTg.KeyEnvelopes[0].DeviceID)
	assert.Equal(t, "member-vault-1", updatedTg.KeyEnvelopes[0].MemberID)
	assert.Equal(t, uint64(1), updatedTg.KeyEnvelopes[0].KEKVersion)
	assert.Equal(t, "opaque-wrapped-kek-bytes", updatedTg.KeyEnvelopes[0].WrappedKEK)
}

func TestAddTrustGroupKeyEnvelope_MultipleDevicesForSameMember(t *testing.T) {
	ctx := context.Background()
	repo := newFakeTrustGroupRepo()
	resolver := newFakeDeviceResolver()
	uc := trustgroup_usecases.NewAddTrustGroupKeyEnvelopeUseCase(repo, resolver)

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Eng Team", []string{"member-vault-1"})
	_, err := repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// Device 1 (laptop)
	resolver.devices["dev-laptop"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-laptop",
		VaultID:   "member-vault-1",
		PublicKey: "pub-laptop",
		Status:    "active",
		IsActive:  true,
	}
	// Device 2 (mobile)
	resolver.devices["dev-mobile"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-mobile",
		VaultID:   "member-vault-1",
		PublicKey: "pub-mobile",
		Status:    "active",
		IsActive:  true,
	}

	// Add envelope for laptop
	_, err = uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
		TrustGroupID: tg.ID,
		MemberID:     "member-vault-1",
		DeviceID:     "dev-laptop",
		KEKVersion:   1,
		WrappedKEK:   "wrapped-laptop",
	})
	require.NoError(t, err)

	// Add envelope for mobile
	updatedTg, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
		TrustGroupID: tg.ID,
		MemberID:     "member-vault-1",
		DeviceID:     "dev-mobile",
		KEKVersion:   1,
		WrappedKEK:   "wrapped-mobile",
	})
	require.NoError(t, err)
	assert.Len(t, updatedTg.KeyEnvelopes, 2)
}

func TestAddTrustGroupKeyEnvelope_Failures(t *testing.T) {
	ctx := context.Background()
	repo := newFakeTrustGroupRepo()
	resolver := newFakeDeviceResolver()
	uc := trustgroup_usecases.NewAddTrustGroupKeyEnvelopeUseCase(repo, resolver)

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Eng Team", []string{"member-vault-1"})
	_, err := repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	resolver.devices["dev-active"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-active",
		VaultID:   "member-vault-1",
		PublicKey: "pub-active",
		Status:    "active",
		IsActive:  true,
	}

	resolver.devices["dev-revoked"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-revoked",
		VaultID:   "member-vault-1",
		PublicKey: "pub-revoked",
		Status:    "revoked",
		IsActive:  false,
	}

	resolver.devices["dev-other-vault"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-other-vault",
		VaultID:   "member-vault-OTHER",
		PublicKey: "pub-other",
		Status:    "active",
		IsActive:  true,
	}

	t.Run("TrustGroup does not exist", func(t *testing.T) {
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: "nonexistent-tg",
			MemberID:     "member-vault-1",
			DeviceID:     "dev-active",
			KEKVersion:   1,
			WrappedKEK:   "wrapped",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrTrustGroupNotFound)
	})

	t.Run("Member not in TrustGroup", func(t *testing.T) {
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "unknown-member",
			DeviceID:     "dev-active",
			KEKVersion:   1,
			WrappedKEK:   "wrapped",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrMemberNotInTrustGroup)
	})

	t.Run("Device does not exist", func(t *testing.T) {
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "member-vault-1",
			DeviceID:     "dev-ghost",
			KEKVersion:   1,
			WrappedKEK:   "wrapped",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrDeviceNotFound)
	})

	t.Run("Device is revoked", func(t *testing.T) {
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "member-vault-1",
			DeviceID:     "dev-revoked",
			KEKVersion:   1,
			WrappedKEK:   "wrapped",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrDeviceRevoked)
	})

	t.Run("Device belongs to another vault", func(t *testing.T) {
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "member-vault-1",
			DeviceID:     "dev-other-vault",
			KEKVersion:   1,
			WrappedKEK:   "wrapped",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrDeviceMemberMismatch)
	})

	t.Run("Stale KEKVersion", func(t *testing.T) {
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "member-vault-1",
			DeviceID:     "dev-active",
			KEKVersion:   999, // Stale/mismatched version
			WrappedKEK:   "wrapped",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrStaleKEKVersion)
	})

	t.Run("Empty WrappedKEK", func(t *testing.T) {
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "member-vault-1",
			DeviceID:     "dev-active",
			KEKVersion:   1,
			WrappedKEK:   "",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrWrappedKEKRequired)
	})

	t.Run("Duplicate active device envelope", func(t *testing.T) {
		// First addition
		_, err := uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "member-vault-1",
			DeviceID:     "dev-active",
			KEKVersion:   1,
			WrappedKEK:   "wrapped-1",
		})
		require.NoError(t, err)

		// Second addition for same active device and version
		_, err = uc.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "member-vault-1",
			DeviceID:     "dev-active",
			KEKVersion:   1,
			WrappedKEK:   "wrapped-2",
		})
		assert.ErrorIs(t, err, trustgroup_domain.ErrDuplicateKeyEnvelope)
	})
}


