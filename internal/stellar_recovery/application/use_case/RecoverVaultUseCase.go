package stellar_recovery_usecase

import (
	"context"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)



type RecoverVaultUseCase struct {
	users       stellar_recovery_domain.UserRepository
	vaults      stellar_recovery_domain.VaultRepository
	subs        stellar_recovery_domain.SubscriptionRepository
	verifier    stellar_recovery_domain.StellarKeyVerifier
	events      stellar_recovery_domain.EventDispatcher
	tokenGen    stellar_recovery_domain.TokenGenerator
}

func NewRecoverVaultUseCase(
	users stellar_recovery_domain.UserRepository,
	vaults stellar_recovery_domain.VaultRepository,
	subs stellar_recovery_domain.SubscriptionRepository,
	verifier stellar_recovery_domain.StellarKeyVerifier,
	events stellar_recovery_domain.EventDispatcher,
	tokenGen stellar_recovery_domain.TokenGenerator,
) *RecoverVaultUseCase {
	return &RecoverVaultUseCase{users, vaults, subs, verifier, events, tokenGen}
}


func (uc *RecoverVaultUseCase) Execute(ctx context.Context, secret string) (*stellar_recovery_domain.RecoveredVault, error) {
	// 1. -------------- Parse public key --------------
	pub, err := uc.verifier.ParseSecret(secret)
	if err != nil {
		return nil, stellar_recovery_domain.ErrInvalidKey
	}

	// 2. -------------- Lookup user --------------
	user, err := uc.users.GetByStellarPublicKey(ctx, pub)
	if err != nil || user == nil {
		return nil, stellar_recovery_domain.ErrUserNotFound
	}

	// 3. -------------- Verify ownership --------------
	if err := uc.verifier.VerifyOwnership(secret, user.StellarPublicKey); err != nil {
		return nil, stellar_recovery_domain.ErrKeyVerification
	}

	// 4. -------------- Load vault --------------
	vault, err := uc.vaults.GetByUserID(ctx, user.ID)
	if err != nil || vault == nil {
		return nil, stellar_recovery_domain.ErrVaultNotFound
	}

	// 5. -------------- Optionally load subscription --------------
	sub, _ := uc.subs.GetActiveByUserID(ctx, user.ID)

	// 6. -------------- Publish domain event --------------
	evt := stellar_recovery_domain.NewVaultRecovered(user.ID, vault.ID, pub)
	_ = uc.events.Dispatch(evt)

	// 7. -------------- Generate session token --------------
	token, _ := uc.tokenGen.NewSessionToken(user.ID)

	return &stellar_recovery_domain.RecoveredVault{
		User:         user,
		Vault:        vault,
		Subscription: sub,
		SessionToken: token,
	}, nil
}
