package billing_persistence

import (
	"context"
	billing_domain "vault-app/internal/billing/domain"

	"gorm.io/gorm"
)


type GormBillingRepository struct {
	db *gorm.DB
}

func NewGormBillingRepository(db *gorm.DB) *GormBillingRepository {
	return &GormBillingRepository{db: db}
}
func (r *GormBillingRepository) TableName() string {
	return "billing_instrument"
}

func (r *GormBillingRepository) Save(ctx context.Context, pm *billing_domain.BillingInstrument) error {
	return r.db.Create(pm).Error
}
func (r *GormBillingRepository) FindByID(ctx context.Context, id string) (*billing_domain.BillingInstrument, error) {
	var pm billing_domain.BillingInstrument
	if err := r.db.First(&pm, id).Error; err != nil {
		return nil, err
	}
	return &pm, nil
}
func (r *GormBillingRepository) FindByUserID(ctx context.Context, id string) (*billing_domain.BillingInstrument, error) {
	var pm billing_domain.BillingInstrument
	if err := r.db.Where("user_id = ?", id).First(&pm).Error; err != nil {
		return nil, err
	}
	return &pm, nil
}

	