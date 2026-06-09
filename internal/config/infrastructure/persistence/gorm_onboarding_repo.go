package persistence

import (
	app_config_domain "vault-app/internal/config/domain"

	"gorm.io/gorm"
)



type GormOnboardingConfigRepository struct {
	db *gorm.DB
	
}	

// NewGormOnboardingConfigRepo creates a new repository instance
func NewGormOnboardingConfigRepo(db *gorm.DB) *GormOnboardingConfigRepository {
	return &GormOnboardingConfigRepository{db: db}
}

func (r *GormOnboardingConfigRepository) Create(dc *app_config_domain.OnboardingConfig) error {
	sqlDB, err := r.ToSqlDB(dc)
	if err != nil {
		return err
	}
	return r.db.Create(&sqlDB).Error
}

func (r *GormOnboardingConfigRepository) Find(id string) (*app_config_domain.OnboardingConfig, error) {
	var dc *OnboardingConfigSqlDB
	err := r.db.First(&dc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return dc.ToDomain()
}

func (r *GormOnboardingConfigRepository) FindByUserID(userID string, opts map[string]interface{}) (*app_config_domain.OnboardingConfig, error) {
	db := r.db
	if opts != nil {
		if optDB, ok := opts["db_transaction"]; ok {
			db = optDB.(*gorm.DB)
		}
	}

	var row OnboardingConfigSqlDB
	err := db.Where("user_id = ?", userID).First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToDomain()
}

func (r *GormOnboardingConfigRepository) Update(id string, dc *app_config_domain.OnboardingConfig, opts map[string]interface{}) error {
	db := r.db
	if opts != nil {
		if optDB, ok := opts["db_transaction"]; ok {
			db = optDB.(*gorm.DB)
		}
	}

	sqlDB, err := r.ToSqlDB(dc)
	if err != nil {
		return err
	}

	return db.Where("user_id = ?", id).Updates(sqlDB).Error
}	

func (r *GormOnboardingConfigRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&OnboardingConfigSqlDB{}).Error
}

func (r *GormOnboardingConfigRepository) FindAll(userID string) ([]*app_config_domain.OnboardingConfig, error) {
	var deviceConfigs []*app_config_domain.OnboardingConfig
	err := r.db.Where("user_id = ?", userID).Find(&deviceConfigs).Error
	return deviceConfigs, err
}



// MAPPER
type OnboardingConfigSqlDB struct {
	UserID             string   `json:"user_id" gorm:"column:user_id"`
	Packs              app_config_domain.StringSlice `json:"packs" gorm:"type:json"`               // ["compliance_critical","client_data"]
	UseCases           app_config_domain.StringSlice `json:"use_cases" gorm:"type:json"`           // optional UI labels if you want
	InstalledSeeds     app_config_domain.StringSlice `json:"installed_seeds" gorm:"type:json"`     // ["devsecops.incident.v1", ...]
	Completed          bool     `json:"completed" gorm:"column:completed"`
	PacksApplied       bool     `json:"packs_applied" gorm:"column:packs_applied"`
}


// SQL DB MODEL TO DOMAIN MODEL
func (sqlDB *OnboardingConfigSqlDB) ToDomain() (*app_config_domain.OnboardingConfig, error) {
	return &app_config_domain.OnboardingConfig{
		UserID:         sqlDB.UserID,
		Packs:          sqlDB.Packs,
		UseCases:       sqlDB.UseCases,
		InstalledSeeds: sqlDB.InstalledSeeds,
		Completed:      sqlDB.Completed,
		PacksApplied:   sqlDB.PacksApplied,
	}, nil
}

// DOMAIN MODEL TO SQL DB MODEL
func (d *GormOnboardingConfigRepository) ToSqlDB(dc *app_config_domain.OnboardingConfig) (*OnboardingConfigSqlDB, error) {
	return &OnboardingConfigSqlDB{
		UserID:         dc.UserID,
		Packs:          dc.Packs,
		UseCases:       dc.UseCases,
		InstalledSeeds: dc.InstalledSeeds,
		Completed:      dc.Completed,
		PacksApplied:   dc.PacksApplied,
	}, nil
}
