package trustgroup_adapters

import (
	"context"
	"errors"
	"fmt"

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

func (a *IdentityDeviceAdapter) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	if a.deviceRepo == nil {
		return nil, errors.New("device repository is nil")
	}

	devices, err := a.deviceRepo.ListByVaultID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	var active []trustgroup_ports.DeviceSummary
	for _, dev := range devices {
		if dev != nil && dev.IsActive() {
			active = append(active, trustgroup_ports.DeviceSummary{
				ID:        dev.ID,
				VaultID:   dev.VaultID,
				PublicKey: dev.PublicKey,
				KeyType:   dev.KeyType,
				Status:    dev.Status,
				IsActive:  true,
			})
		}
	}
	fmt.Printf("[DIAGNOSTIC][IdentityDeviceAdapter.ListActiveDevices] memberID=%s totalDevices=%d activeDevices=%d\n", memberID, len(devices), len(active))
	return active, nil
}

var _ trustgroup_ports.DeviceResolver = (*IdentityDeviceAdapter)(nil)
