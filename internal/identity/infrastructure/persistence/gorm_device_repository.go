package identity_persistence

import (
	"context"
	"errors"
	"fmt"

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
	err := r.db.WithContext(ctx).Create(d).Error
	fmt.Printf("[DIAGNOSTIC][GormDeviceRepository.Save] DB_Ptr=%p Saved Device ID=%s VaultID=%s Status=%s err=%v\n", r.db, d.ID, d.VaultID, d.Status, err)
	return err
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
	var totalCount int64
	r.db.WithContext(ctx).Model(&identity_domain.Device{}).Count(&totalCount)

	var allDevices []*identity_domain.Device
	r.db.WithContext(ctx).Find(&allDevices)

	fmt.Printf("[DIAGNOSTIC][GormDeviceRepository.ListByVaultID] DB_Ptr=%p RequestedVaultID=%s TotalDevicesInDB=%d\n", r.db, vaultID, totalCount)
	for i, d := range allDevices {
		if d != nil {
			fmt.Printf("  -> DB Row[%d]: ID=%s VaultID=%s Status=%s RevokedAt=%v IsActive=%t\n", i, d.ID, d.VaultID, d.Status, d.RevokedAt, d.IsActive())
		}
	}

	var devices []*identity_domain.Device
	err := r.db.WithContext(ctx).Where("vault_id = ?", vaultID).Find(&devices).Error
	fmt.Printf("[DIAGNOSTIC][GormDeviceRepository.ListByVaultID] Query Result: vault_id=%s MatchingRowsCount=%d err=%v\n", vaultID, len(devices), err)
	if err != nil {
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
