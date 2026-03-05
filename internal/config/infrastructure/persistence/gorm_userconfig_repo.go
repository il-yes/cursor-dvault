package persistence

import (
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormUserConfigRepository struct {
	db *gorm.DB
}

func NewGormUserConfigRepository(db *gorm.DB) *GormUserConfigRepository {
	return &GormUserConfigRepository{db: db}
}	
func (GormUserConfigRepository) TableName() string {
	return "config_user"
}
func (r *GormUserConfigRepository) CreateUserConfig(userConfig *app_config_domain.UserConfig) error {
	utils.LogPretty("GormUserConfigRepository - CreateUserConfig - userConfig", userConfig)
	 // CRITICAL: Set IDs if empty
    if userConfig.ID == "" {
        userConfig.ID = uuid.NewString()
    }
    
    for i := range userConfig.SharingRules {
        if userConfig.SharingRules[i].ID == "" {
            userConfig.SharingRules[i].ID = uuid.NewString()
        }
        userConfig.SharingRules[i].UserConfigID = userConfig.ID
    }
	userConfigMapper := r.toUserConfigMapper(userConfig)
	return r.db.Create(userConfigMapper).Error
}

func (r *GormUserConfigRepository) GetUserConfig(id string) (*app_config_domain.UserConfig, error) {
	var userConfig UserConfigMapper
	if err := r.db.First(&userConfig, "id = ?", id).Error; err != nil {
		return nil, err
	}
	userConfigDomain := r.toUserConfigModel(&userConfig)
	return userConfigDomain, nil
}	

func (r *GormUserConfigRepository) UpdateUserConfig(userConfig *app_config_domain.UserConfig) error {
	
	return r.db.Save(r.toUserConfigMapper(userConfig)).Error
}

func (r *GormUserConfigRepository) DeleteUserConfig(id string) error {
	return r.db.Delete(&UserConfigMapper{}, "id = ?", id).Error
}	
	

// -------- Mappers --------
type UserConfigMapper struct {
	ID               string               `json:"id" gorm:"primaryKey;autoIncrement:false;size:36;uniqueIndex"`
	Role             string               `json:"role" gorm:"column:role"`
	Signature        string               `json:"signature" gorm:"column:signature"`
	ConnectedOrgs    []string             `json:"connected_orgs" gorm:"type:json;serializer:json"`
	StellarAccount   app_config_domain.StellarAccountConfig `json:"stellar_account" gorm:"embedded;embeddedPrefix:stellar_"`
	SharingRules     []app_config_domain.SharingRule        `json:"sharing_rules" gorm:"foreignKey:UserConfigID;constraint:OnDelete:CASCADE"`
	TwoFactorEnabled bool                 `json:"two_factor_enabled" yaml:"two_factor_enabled" gorm:"column:two_factor_enabled"`
}
func (r *GormUserConfigRepository) toUserConfigMapper(userConfig *app_config_domain.UserConfig) *UserConfigMapper {
	return &UserConfigMapper{
		ID:             userConfig.ID,
		Role:           userConfig.Role,
		Signature:      userConfig.Signature,
		ConnectedOrgs:  userConfig.ConnectedOrgs,
		StellarAccount: userConfig.StellarAccount,
		SharingRules:   userConfig.SharingRules,
		TwoFactorEnabled: userConfig.TwoFactorEnabled,
	}
}

func (r *GormUserConfigRepository) toUserConfigModel(userConfig *UserConfigMapper) *app_config_domain.UserConfig {
	return &app_config_domain.UserConfig{
		ID:             userConfig.ID,
		Role:           userConfig.Role,
		Signature:      userConfig.Signature,
		ConnectedOrgs:  userConfig.ConnectedOrgs,
		StellarAccount: userConfig.StellarAccount,
		SharingRules:   userConfig.SharingRules,
		TwoFactorEnabled: userConfig.TwoFactorEnabled,
	}
}