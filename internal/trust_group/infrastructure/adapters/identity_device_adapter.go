package trustgroup_adapters

import (
	"context"
	"errors"

	identity_domain "vault-app/internal/identity/domain"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
)

type IdentityDeviceAdapter struct {
	deviceRepo identity_domain.DeviceRepository
}

func NewIdentityDeviceAdapter(deviceRepo identity_domain.DeviceRepository) *IdentityDeviceAdapter {
	return &IdentityDeviceAdapter{deviceRepo: deviceRepo}
}

func (a *IdentityDeviceAdapter) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	if a.deviceRepo == nil {
		return nil, errors.New("device repository is nil")
	}

	dev, err := a.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, identity_domain.ErrDeviceNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if dev == nil {
		return nil, nil
	}

	return &trustgroup_ports.DeviceSummary{
		ID:        dev.ID,
		VaultID:   dev.VaultID,
		PublicKey: dev.PublicKey,
		KeyType:   dev.KeyType,
		Status:    dev.Status,
		IsActive:  dev.IsActive(),
	}, nil
}

var _ trustgroup_ports.DeviceResolver = (*IdentityDeviceAdapter)(nil)
