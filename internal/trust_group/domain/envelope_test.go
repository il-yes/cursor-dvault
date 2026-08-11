package trustgroup_domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	trustgroup_domain "vault-app/internal/trust_group/domain"
)

func TestTrustGroup_AddEnvelope_Success(t *testing.T) {
	tg := trustgroup_domain.NewTrustGroup("channel-001", "OEM Group", []string{"member-1"})
	require.Equal(t, uint64(1), tg.KEKVersion)
	require.Empty(t, tg.KeyEnvelopes)

	env := trustgroup_domain.TrustGroupKeyEnvelope{
		MemberID:   "member-1",
		DeviceID:   "device-001",
		KEKVersion: 1,
		WrappedKEK: "wrapped-kek-value",
	}

	err := tg.AddEnvelope(env)
	require.NoError(t, err)
	require.Len(t, tg.KeyEnvelopes, 1)
	require.Equal(t, "member-1", tg.KeyEnvelopes[0].MemberID)
	require.Equal(t, "device-001", tg.KeyEnvelopes[0].DeviceID)
	require.Equal(t, uint64(1), tg.KeyEnvelopes[0].KEKVersion)
	require.NotEmpty(t, tg.KeyEnvelopes[0].ID)
	require.Equal(t, tg.ID, tg.KeyEnvelopes[0].TrustGroupID)
	require.False(t, tg.KeyEnvelopes[0].CreatedAt.IsZero())
	require.True(t, tg.IsDirty)
}

func TestTrustGroup_AddEnvelope_StaleKEKVersion(t *testing.T) {
	tg := trustgroup_domain.NewTrustGroup("channel-001", "OEM Group", []string{"member-1"})
	tg.KEKVersion = 2

	env := trustgroup_domain.TrustGroupKeyEnvelope{
		MemberID:   "member-1",
		DeviceID:   "device-001",
		KEKVersion: 1, // Stale
		WrappedKEK: "wrapped-kek-value",
	}

	err := tg.AddEnvelope(env)
	require.Error(t, err)
	require.ErrorIs(t, err, trustgroup_domain.ErrStaleKEKVersion)
	require.Empty(t, tg.KeyEnvelopes)
}

func TestTrustGroup_AddEnvelope_DuplicateActiveDeviceEnvelope(t *testing.T) {
	tg := trustgroup_domain.NewTrustGroup("channel-001", "OEM Group", []string{"member-1"})

	env1 := trustgroup_domain.TrustGroupKeyEnvelope{
		MemberID:   "member-1",
		DeviceID:   "device-001",
		KEKVersion: 1,
		WrappedKEK: "wrapped-kek-value-1",
	}
	err := tg.AddEnvelope(env1)
	require.NoError(t, err)

	// Attempting to add duplicate active envelope for device-001 at KEKVersion 1
	env2 := trustgroup_domain.TrustGroupKeyEnvelope{
		MemberID:   "member-1",
		DeviceID:   "device-001",
		KEKVersion: 1,
		WrappedKEK: "wrapped-kek-value-2",
	}
	err = tg.AddEnvelope(env2)
	require.Error(t, err)
	require.ErrorIs(t, err, trustgroup_domain.ErrDuplicateKeyEnvelope)
	require.Len(t, tg.KeyEnvelopes, 1)
}

func TestTrustGroup_AddEnvelope_RevokedDeviceAllowsNewEnvelope(t *testing.T) {
	tg := trustgroup_domain.NewTrustGroup("channel-001", "OEM Group", []string{"member-1"})

	revokedTime := time.Now().Add(-1 * time.Hour)
	env1 := trustgroup_domain.TrustGroupKeyEnvelope{
		MemberID:   "member-1",
		DeviceID:   "device-001",
		KEKVersion: 1,
		WrappedKEK: "wrapped-kek-value-old",
		RevokedAt:  &revokedTime,
	}
	err := tg.AddEnvelope(env1)
	require.NoError(t, err)

	// New active envelope for same device is allowed since the previous one is revoked
	env2 := trustgroup_domain.TrustGroupKeyEnvelope{
		MemberID:   "member-1",
		DeviceID:   "device-001",
		KEKVersion: 1,
		WrappedKEK: "wrapped-kek-value-new",
	}
	err = tg.AddEnvelope(env2)
	require.NoError(t, err)
	require.Len(t, tg.KeyEnvelopes, 2)
}

func TestTrustGroup_AddEnvelope_MissingFields(t *testing.T) {
	tg := trustgroup_domain.NewTrustGroup("channel-001", "OEM Group", []string{"member-1"})

	tests := []struct {
		name          string
		envelope      trustgroup_domain.TrustGroupKeyEnvelope
		expectedError error
	}{
		{
			name: "missing member id",
			envelope: trustgroup_domain.TrustGroupKeyEnvelope{
				DeviceID:   "device-001",
				KEKVersion: 1,
				WrappedKEK: "wrapped-kek",
			},
			expectedError: trustgroup_domain.ErrMemberIDRequired,
		},
		{
			name: "missing device id",
			envelope: trustgroup_domain.TrustGroupKeyEnvelope{
				MemberID:   "member-1",
				KEKVersion: 1,
				WrappedKEK: "wrapped-kek",
			},
			expectedError: trustgroup_domain.ErrDeviceIDRequired,
		},
		{
			name: "missing wrapped kek",
			envelope: trustgroup_domain.TrustGroupKeyEnvelope{
				MemberID:   "member-1",
				DeviceID:   "device-001",
				KEKVersion: 1,
			},
			expectedError: trustgroup_domain.ErrWrappedKEKRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tg.AddEnvelope(tt.envelope)
			require.Error(t, err)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
