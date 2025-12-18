package driver

import (
	"vault-app/internal/auth"
	// auth_domain "vault-app/internal/auth/domain"
	// auth_persistence "vault-app/internal/auth/infrastructure/persistence"
	app_config "vault-app/internal/config"
	share_domain "vault-app/internal/domain/shared"
	identity_domain "vault-app/internal/identity/domain"
	share_infrastructure "vault-app/internal/infrastructure/share"
	"vault-app/internal/models"
	onboarding_persistence "vault-app/internal/onboarding/infrastructure/persistence"
	subscription_persistence "vault-app/internal/subscription/infrastructure/persistence"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Folder{},
		&models.VaultCID{},
		&models.VaultContent{},
		&models.LoginEntry{},
		&models.CardEntry{},
		&models.IdentityEntry{},
		&models.NoteEntry{},
		&models.SSHKeyEntry{},

		// App & User Configs
		&app_config.AppConfig{},
		&app_config.CommitRule{},
		&app_config.UserConfig{},
		&app_config.SharingRule{},
		&app_config.SharingConfig{}, // if used for advanced sharing
		&models.UserSession{},
		&auth.TokenPairs{},

		// Sharing
		&share_infrastructure.ShareEntryModel{},
		&share_infrastructure.RecipientModel{},
		&share_domain.AuditLog{},

		// Onboarding
		&onboarding_persistence.UserDB{},
		&subscription_persistence.SubscriptionMapper{},
		&subscription_persistence.UserSubscriptionMapper{},

		// Auth

		// Identity
		&identity_domain.User{},


		// vault
		&vaults_persistence.VaultMapper{},
		&vaults_persistence.SessionMapper{},
	)
}	