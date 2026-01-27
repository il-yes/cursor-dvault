package stellar_recovery_usecase

import (
	"context"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

type CheckKeyUseCase struct {
	users   stellar_recovery_domain.UserRepository
	vaults  stellar_recovery_domain.VaultRepository
	subs    stellar_recovery_domain.SubscriptionRepository
	verifier stellar_recovery_domain.StellarKeyVerifier
}

func NewCheckKeyUseCase(
	users stellar_recovery_domain.UserRepository,
	vaults stellar_recovery_domain.VaultRepository,
	subs stellar_recovery_domain.SubscriptionRepository,
	verifier stellar_recovery_domain.StellarKeyVerifier,
) *CheckKeyUseCase {
	return &CheckKeyUseCase{users: users, vaults: vaults, subs: subs, verifier: verifier}
}

type CheckKeyResult struct {
	VaultExists bool
	PublicKey   string
	Vault       *stellar_recovery_domain.Vault
	Subscription *stellar_recovery_domain.Subscription
}

func (uc *CheckKeyUseCase) Execute(ctx context.Context, pub string) (*CheckKeyResult, error) {
	// 1. -------------- Validate public key --------------
	// pub, err := uc.verifier.ParseSecret(secret)
	if pub == "" {
		// invalid secret format -> return VaultExists false and the parsed public key empty
		return &CheckKeyResult{
			VaultExists: false,
			PublicKey:   "",
		}, nil
	}

	// 2. -------------- Check if user exists --------------
	user, err := uc.users.GetByStellarPublicKey(ctx, pub)
	if err != nil {
		// no user found -> vault does not exist
		return &CheckKeyResult{
			VaultExists: false,
			PublicKey:   pub,
		}, nil
	}

	// 3. -------------- Get vault and subscription --------------
	vault, _ := uc.vaults.GetByUserID(ctx, user.ID)
	sub, _ := uc.subs.GetActiveByUserID(ctx, user.ID)

	// 4. -------------- Fires Checked Vault Event --------------

	return &CheckKeyResult{
		VaultExists: true,
		PublicKey:   pub,
		Vault:       vault,
		Subscription: sub,
	}, nil
}

