package collaboration_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserBIdentityTrace struct {
	JwtSubCallerIdentityID string
	InviteeVaultID         string
	PersistedMemberCID     string
	RetrievedMemberCID     string
	DeviceVaultID          string
	DeviceIdentityID       string
}

func EvaluateUserBIdentityMapping(trace UserBIdentityTrace) (bool, []string) {
	var discrepancies []string

	caller := strings.TrimSpace(trace.JwtSubCallerIdentityID)
	invitee := strings.TrimSpace(trace.InviteeVaultID)
	persisted := strings.TrimSpace(trace.PersistedMemberCID)
	retrieved := strings.TrimSpace(trace.RetrievedMemberCID)
	devVault := strings.TrimSpace(trace.DeviceVaultID)
	devIdent := strings.TrimSpace(trace.DeviceIdentityID)

	if caller != persisted {
		discrepancies = append(discrepancies, fmt.Sprintf("callerIdentityID %q != persistedMemberCID %q", caller, persisted))
	}
	if caller != retrieved {
		discrepancies = append(discrepancies, fmt.Sprintf("callerIdentityID %q != retrievedMemberCID %q", caller, retrieved))
	}
	if invitee != persisted {
		discrepancies = append(discrepancies, fmt.Sprintf("inviteeVaultID %q != persistedMemberCID %q", invitee, persisted))
	}
	if devVault != caller && devIdent != caller {
		discrepancies = append(discrepancies, fmt.Sprintf("device VaultID/IdentityID (%q/%q) != callerIdentityID %q", devVault, devIdent, caller))
	}

	isExactMatch := (caller == retrieved)
	return isExactMatch, discrepancies
}

func TestDiagnostic_UserBIdentityMappingComparison(t *testing.T) {
	t.Run("Authoritative matching identities trace", func(t *testing.T) {
		trace := UserBIdentityTrace{
			JwtSubCallerIdentityID: "vault_user_b_777",
			InviteeVaultID:         "vault_user_b_777",
			PersistedMemberCID:     "vault_user_b_777",
			RetrievedMemberCID:     "vault_user_b_777",
			DeviceVaultID:          "vault_user_b_777",
			DeviceIdentityID:       "vault_user_b_777",
		}

		matched, discrepancies := EvaluateUserBIdentityMapping(trace)
		assert.True(t, matched)
		assert.Empty(t, discrepancies)
	})

	t.Run("Detect URI scheme discrepancy e.g. ankhora:// prefix", func(t *testing.T) {
		trace := UserBIdentityTrace{
			JwtSubCallerIdentityID: "vault_user_b_777",
			InviteeVaultID:         "ankhora://vault_user_b_777",
			PersistedMemberCID:     "ankhora://vault_user_b_777",
			RetrievedMemberCID:     "ankhora://vault_user_b_777",
			DeviceVaultID:          "vault_user_b_777",
			DeviceIdentityID:       "vault_user_b_777",
		}

		matched, discrepancies := EvaluateUserBIdentityMapping(trace)
		assert.False(t, matched)
		assert.NotEmpty(t, discrepancies)
		t.Logf("Detected identity format discrepancy: %v", discrepancies)
	})
}
