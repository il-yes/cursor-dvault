package stellar_recovery_usecase

import (
	"context"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

type ImportKeyUseCase struct {
	UserRepo   stellar_recovery_domain.UserRepository
	KeyService stellar_recovery_domain.StellarKeyValidator
}

func NewImportKeyUseCase(userRepo stellar_recovery_domain.UserRepository, keyService stellar_recovery_domain.StellarKeyValidator) *ImportKeyUseCase {
	return &ImportKeyUseCase{
		UserRepo:   userRepo,
		KeyService: keyService,
	}
}	

func (uc *ImportKeyUseCase) Execute(ctx context.Context, secret string) (*stellar_recovery_domain.ImportedKey, error) {
	pub, err := uc.KeyService.ParseSecret(secret)
	if err != nil {
		return nil, stellar_recovery_domain.ErrInvalidKey
	}

	user, _ := uc.UserRepo.GetByStellarPublicKey(ctx, pub)
	if user != nil {
		return nil, stellar_recovery_domain.ErrKeyAlreadyUsed
	}

	return &stellar_recovery_domain.ImportedKey{
		StellarPublic: pub,
		StellarSecret: secret,
		CanCreate:     true,
	}, nil
}
