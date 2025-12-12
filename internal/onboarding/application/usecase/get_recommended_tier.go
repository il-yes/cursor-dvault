package onboarding_usecase

import (
	identity_domain "vault-app/internal/identity/domain"
	subscription_domain "vault-app/internal/subscription/domain"

	"gorm.io/gorm"
)




type GetRecommendedTierUseCase struct {
	Db *gorm.DB

}	




// GetRecommendedTier returns recommended tier based on identity choice
func (a *GetRecommendedTierUseCase) Execute(identity identity_domain.IdentityChoice) subscription_domain.SubscriptionTier {
    switch identity {
	case identity_domain.IdentityPersonal:
        return subscription_domain.TierPro
    case identity_domain.IdentityAnonymous:
        return subscription_domain.TierProPlus
    case identity_domain.IdentityTeam:
        return subscription_domain.TierBusiness
    case identity_domain.IdentityCompliance:
        return subscription_domain.TierBusiness
    default:
        return subscription_domain.TierFree
    }
}
