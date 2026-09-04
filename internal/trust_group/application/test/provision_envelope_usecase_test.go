package trustgroup_usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_usecases "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

func TestProvisionTrustGroupDeviceEnvelope_Success(t *testing.T) {
	ctx := context.Background()
	repo := newFakeTrustGroupRepo()
	resolver := newFakeDeviceResolver()

	addEnvUC := trustgroup_usecases.NewAddTrustGroupKeyEnvelopeUseCase(repo, resolver)
	provisionUC := trustgroup_usecases.NewProvisionTrustGroupDeviceEnvelopeUseCase(repo, resolver, nil, addEnvUC, nil)

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Eng Team", []string{"vault_alice", "vault_bob"})
	tg.ID = "tg_authoritative_101"
	_, err := repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	validStellarPubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"

	resolver.devices["dev_bob_laptop"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_bob_laptop",
		VaultID:   "vault_bob",
		PublicKey: validStellarPubKey,
		Status:    "active",
		IsActive:  true,
	}

	req := trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    "tg_authoritative_101",
		MemberID:        "vault_bob",
		DeviceID:        "dev_bob_laptop",
		DevicePublicKey: validStellarPubKey,
	}

	updatedTg, err := provisionUC.Execute(ctx, req, nil)
	require.NoError(t, err)
	require.NotNil(t, updatedTg)
	require.Len(t, updatedTg.KeyEnvelopes, 1)

	env := updatedTg.KeyEnvelopes[0]
	assert.Equal(t, "tg_authoritative_101", env.TrustGroupID)
	assert.Equal(t, "vault_bob", env.MemberID)
	assert.Equal(t, "dev_bob_laptop", env.DeviceID)
	assert.Equal(t, uint64(1), env.KEKVersion)
	assert.NotEmpty(t, env.WrappedKEK)
}

func TestProvisionTrustGroupDeviceEnvelope_MemberNotInTrustGroup(t *testing.T) {
	ctx := context.Background()
	repo := newFakeTrustGroupRepo()
	resolver := newFakeDeviceResolver()

	addEnvUC := trustgroup_usecases.NewAddTrustGroupKeyEnvelopeUseCase(repo, resolver)
	provisionUC := trustgroup_usecases.NewProvisionTrustGroupDeviceEnvelopeUseCase(repo, resolver, nil, addEnvUC, nil)

	tg := trustgroup_domain.NewTrustGroup("chan-1", "Eng Team", []string{"vault_alice"})
	tg.ID = "tg_authoritative_101"
	_, err := repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	validStellarPubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"

	resolver.devices["dev_bob_laptop"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev_bob_laptop",
		VaultID:   "vault_bob",
		PublicKey: validStellarPubKey,
		Status:    "active",
		IsActive:  true,
	}

	req := trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    "tg_authoritative_101",
		MemberID:        "vault_bob",
		DeviceID:        "dev_bob_laptop",
		DevicePublicKey: validStellarPubKey,
	}

	_, err = provisionUC.Execute(ctx, req, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, trustgroup_domain.ErrMemberNotInTrustGroup)
}
