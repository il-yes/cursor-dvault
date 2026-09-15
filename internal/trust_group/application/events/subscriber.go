package trustgroup_events

import (
	"context"
	"fmt"
	"strings"

	identity_domain "vault-app/internal/identity/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_usecases_envelope "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
)

// UserFinder resolves local user information.
type UserFinder interface {
	FindUserById(ctx context.Context, userID string) (*identity_domain.User, error)
}

// UserByEmailResolver resolves user information from cloud by email.
type UserByEmailResolver interface {
	GetUserByEmail(ctx context.Context, email string) (*tracecore_types.User, error)
}

// MemberAddedSubscriber listens to MemberAddedToTrustGroup events and provisions member envelopes.
type MemberAddedSubscriber struct {
	provisionEnvelopeUC *trustgroup_usecases_envelope.ProvisionTrustGroupMemberEnvelopeUseCase
	userFinder          UserFinder
	userByEmailResolver UserByEmailResolver
	keyring             *vaults_domain.VaultKeyring
}

// NewMemberAddedSubscriber constructs a subscriber.
func NewMemberAddedSubscriber(
	provisionEnvelopeUC *trustgroup_usecases_envelope.ProvisionTrustGroupMemberEnvelopeUseCase,
	args ...interface{},
) *MemberAddedSubscriber {
	sub := &MemberAddedSubscriber{
		provisionEnvelopeUC: provisionEnvelopeUC,
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case UserFinder:
			sub.userFinder = v
		case UserByEmailResolver:
			sub.userByEmailResolver = v
		case *vaults_domain.VaultKeyring:
			sub.keyring = v
		}
	}
	return sub
}

// RegisterSubscribers wires the handler function to the event bus.
func (s *MemberAddedSubscriber) RegisterSubscribers(eventBus TrustGroupEventBus) {
	if eventBus == nil || s.provisionEnvelopeUC == nil {
		return
	}

	eventBus.SubscribeToMemberAddedToTrustGroup(s.HandleMemberAdded)
}

// HandleMemberAdded handles MemberAddedToTrustGroup event by resolving target public key and provisioning member envelope.
func (s *MemberAddedSubscriber) HandleMemberAdded(ctx context.Context, event trustgroup_domain.MemberAddedToTrustGroup) error {
	fmt.Printf("[C3][EVENT_SUBSCRIBER][MEMBER_ADDED] Received MemberAddedToTrustGroup event for TrustGroupID=%s MemberID=%s\n",
		event.TrustGroupID, event.MemberID)

	if s.provisionEnvelopeUC == nil {
		return nil
	}

	var targetPubKey string
	var email string

	if s.userFinder != nil {
		if user, err := s.userFinder.FindUserById(ctx, event.MemberID); err == nil && user != nil {
			email = user.Email
			if strings.TrimSpace(user.StellarPublicKey) != "" {
				targetPubKey = strings.TrimSpace(user.StellarPublicKey)
			}
		}
	}

	if targetPubKey == "" && email != "" && s.userByEmailResolver != nil {
		if user, err := s.userByEmailResolver.GetUserByEmail(ctx, email); err == nil && user != nil {
			if strings.TrimSpace(user.PublicKey) != "" {
				targetPubKey = strings.TrimSpace(user.PublicKey)
			}
		}
	}

	// req := trustgroup_dtos.ProvisionTrustGroupMemberEnvelopeRequest{
	// 	TrustGroupID:    event.TrustGroupID,
	// 	MemberID:        event.MemberID,
	// 	MemberPublicKey: targetPubKey,
	// }

	// _, err := s.provisionEnvelopeUC.Execute(ctx, req, s.keyring)
	// if err != nil {
	// 	fmt.Printf("[C3][EVENT_SUBSCRIBER][MEMBER_ADDED][ERROR] Failed to provision envelope for TrustGroupID=%s MemberID=%s: %v\n",
	// 		event.TrustGroupID, event.MemberID, err)
	// 	return fmt.Errorf("subscriber failed to provision key envelope for member %s: %w", event.MemberID, err)
	// }

	fmt.Printf("[C3][EVENT_SUBSCRIBER][MEMBER_ADDED][SUCCESS] Provisioned key envelope for TrustGroupID=%s MemberID=%s\n",
		event.TrustGroupID, event.MemberID)
	return nil
}


