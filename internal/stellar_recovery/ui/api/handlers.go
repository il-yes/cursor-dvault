package stellar_recovery_ui_api

import (
	"context"
	"fmt"
	"time"
	"vault-app/internal/handlers"
	stellar_recovery_usecase "vault-app/internal/stellar_recovery/application/use_case"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

type StellarRecoveryHandler struct {
	checkUC   *stellar_recovery_usecase.CheckKeyUseCase
	recoverUC *stellar_recovery_usecase.RecoverVaultUseCase
	importUC  *stellar_recovery_usecase.ImportKeyUseCase
	connectUC *stellar_recovery_usecase.ConnectWithStellarUseCase
}

func NewStellarRecoveryHandler(
	checkUC *stellar_recovery_usecase.CheckKeyUseCase,
	recoverUC *stellar_recovery_usecase.RecoverVaultUseCase,
	importUC *stellar_recovery_usecase.ImportKeyUseCase,
	connectUC *stellar_recovery_usecase.ConnectWithStellarUseCase,
) *StellarRecoveryHandler {
	return &StellarRecoveryHandler{
		checkUC:   checkUC,
		recoverUC: recoverUC,
		importUC:  importUC,
		connectUC: connectUC,
	}
}
type CheckKeyResponse struct {
	ID               string  `json:"id"`
	CreatedAt        string  `json:"created_at"`
	SubscriptionTier string  `json:"subscription_tier"`
	StorageUsedGB    float64 `json:"storage_used_gb"`
	LastSyncedAt     string  `json:"last_synced_at"`
	Ok               bool    `json:"ok"`
}

func (h *StellarRecoveryHandler) CheckVault(ctx context.Context, stellarKey string) (*CheckKeyResponse, error) {
	result, err := h.checkUC.Execute(ctx, stellarKey)
	if err != nil {
		return nil, err
	}

	if result == nil || !result.VaultExists {
		return &CheckKeyResponse{Ok: false}, nil
	}

	v := result.Vault
	subTier := "Free"
	if result.Subscription != nil {
		subTier = result.Subscription.Tier
	}

	lastSynced := ""
	if v.LastSyncedAt != nil {
		lastSynced = v.LastSyncedAt.Format(time.RFC3339)
	}

	return &CheckKeyResponse{
		ID:               v.ID,
		CreatedAt:        v.CreatedAt.Format(time.RFC3339),
		SubscriptionTier: subTier,
		StorageUsedGB:    v.StorageUsedGB,
		LastSyncedAt:     lastSynced,
		Ok:               true,
	}, nil
}

func (h *StellarRecoveryHandler) RecoverVault(ctx context.Context, stellarKey string) (*stellar_recovery_domain.RecoveredVault, error) {
	return h.recoverUC.Execute(ctx, stellarKey)
}


func (h *StellarRecoveryHandler) ImportKey(ctx context.Context, stellarKey string) (*stellar_recovery_domain.ImportedKey, error) {
	return h.importUC.Execute(ctx, stellarKey)
}

func (h *StellarRecoveryHandler) ConnectWithStellar(ctx context.Context, req handlers.LoginRequest) (*CheckKeyResponse, error) {
    result, err := h.connectUC.Execute(ctx, req)
    if err != nil {
        return nil, err
    }
	fmt.Println("ConnectWithStellar result", result)

    // no vault → onboarding continues
    if result.Vault == nil {
        return nil, nil
    }

    v := result.Vault

    subTier := "Free"
    if result.Subscription != nil {
        subTier = result.Subscription.Tier
    }

    lastSync := ""
    if v.LastSyncedAt != nil {
        lastSync = v.LastSyncedAt.Format("2006-01-02")
    }

    return &CheckKeyResponse{
        ID:               v.ID,
        CreatedAt:        v.CreatedAt.Format("2006-01-02"),
        SubscriptionTier: subTier,
        StorageUsedGB:    v.StorageUsedGB,
        LastSyncedAt:     lastSync,
        Ok:               true,
    }, nil
}




