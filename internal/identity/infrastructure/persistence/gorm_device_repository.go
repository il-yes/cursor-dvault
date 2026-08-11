package identity_persistence

import (
	"context"
	"errors"

	identity_domain "vault-app/internal/identity/domain"

	"gorm.io/gorm"
)

type GormDeviceRepository struct {
	db *gorm.DB
}

func NewGormDeviceRepository(db *gorm.DB) *GormDeviceRepository {
	return &GormDeviceRepository{db: db}
}

func (r *GormDeviceRepository) Save(ctx context.Context, d *identity_domain.Device) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *GormDeviceRepository) FindByID(ctx context.Context, id string) (*identity_domain.Device, error) {
	var d identity_domain.Device
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, identity_domain.ErrDeviceNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *GormDeviceRepository) ListByVaultID(ctx context.Context, vaultID string) ([]*identity_domain.Device, error) {
	var devices []*identity_domain.Device
	if err := r.db.WithContext(ctx).Where("vault_id = ?", vaultID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *GormDeviceRepository) Update(ctx context.Context, d *identity_domain.Device) error {
	res := r.db.WithContext(ctx).Save(d)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return identity_domain.ErrDeviceNotFound
	}
	return nil
}

// Ensure interface satisfaction at compile-time
var _ identity_domain.DeviceRepository = (*GormDeviceRepository)(nil)
