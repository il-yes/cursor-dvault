package onboarding_usecase

import subscription_domain "vault-app/internal/subscription/domain"



type GetTierFeaturesUseCase struct {
	
}

func NewGetTierFeaturesUseCase() *GetTierFeaturesUseCase {
	return &GetTierFeaturesUseCase{}
}


// GetTierFeatures returns feature comparison for pricing page
func (a *GetTierFeaturesUseCase) Execute() map[string]subscription_domain.SubscriptionFeatures {
    return map[string]subscription_domain.SubscriptionFeatures{
        "free": {
            StorageGB:        5,
            CloudBackup:      false,
            MobileApps:       false,
            SharingLimit:     5,
            Support:          "community",
        },
        "pro": {
            StorageGB:        100,
            CloudBackup:      true,
            MobileApps:       true,
            UnlimitedSharing: true,
            Support:          "email_24_48h",
        },
        "pro_plus": {
            StorageGB:         200,
            CloudBackup:       true,
            MobileApps:        true,
            UnlimitedSharing:  true,
            VersionHistory:    true,
            Telemetry:         false,
            AnonymousAccount:  true,
            CryptoPayments:    true,
            EncryptedPayments: true,
            Support:           "encrypted_chat_12h",
        },
        "business": {
            StorageGB:        1024,
            CloudBackup:      true,
            MobileApps:       true,
            UnlimitedSharing: true,
            VersionHistory:   true,
            Telemetry:        false,
            CryptoPayments:   true,
            EncryptedPayments: true,
            APIAccess:        true,
            Tracecore:        true,
            SSO:              true,
            TeamFeatures:     true,
            Support:          "24_7_live",
        },
    }
}