package stellar_recovery_usecase

import (
	"context"
	"vault-app/internal/handlers"
	"vault-app/internal/models"
	shared "vault-app/internal/shared/stellar"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)




type ConnectWithStellarResult struct {
    Password     string
    User         *models.User
    Vault        *stellar_recovery_domain.Vault
    Subscription *stellar_recovery_domain.Subscription
}

// func NewConnectWithStellarUseCase(stellarAdapter *shared.StellarLoginAdapter) *ConnectWithStellarUseCase {
// 	return &ConnectWithStellarUseCase{stellarAdapter: stellarAdapter}
// }	
type ConnectWithStellarUseCase struct {
    StellarPort shared.StellarLoginPort
    VaultRepo   stellar_recovery_domain.VaultRepository
    SubRepo     stellar_recovery_domain.SubscriptionRepository
}


func NewConnectWithStellarUseCase(
    stellarPort shared.StellarLoginPort,
    vaultRepo stellar_recovery_domain.VaultRepository,
    subRepo stellar_recovery_domain.SubscriptionRepository,
) *ConnectWithStellarUseCase {
    return &ConnectWithStellarUseCase{
        StellarPort: stellarPort,
        VaultRepo:   vaultRepo,
        SubRepo:     subRepo,
    }
}


func (uc *ConnectWithStellarUseCase) Execute(ctx context.Context, req handlers.LoginRequest) (*ConnectWithStellarResult, error) {
    password, user, err := uc.StellarPort.RecoverPassword(ctx, shared.RecoverPasswordInput{
        PublicKey: req.PublicKey,
        SignedMessage: req.SignedMessage , // frontend later provides this
        Signature: req.Signature,
    })
    if err != nil {
        return nil, err
    }

    vault, _ := uc.VaultRepo.GetByUserID(ctx, user.ID)
    sub, _ := uc.SubRepo.GetActiveByUserID(ctx, user.ID)

    return &ConnectWithStellarResult{
        Password:     password,
        User:         user,
        Vault:        vault,
        Subscription: sub,
    }, nil
}
