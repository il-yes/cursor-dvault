package identity_persistence

import (
	"context"
	"sync"

	identity_domain "vault-app/internal/identity/domain"
)

type MemoryDeviceRepository struct {
	mu   sync.RWMutex
	byID map[string]*identity_domain.Device
}

func NewMemoryDeviceRepository() *MemoryDeviceRepository {
	return &MemoryDeviceRepository{
		byID: make(map[string]*identity_domain.Device),
	}
}

func (r *MemoryDeviceRepository) Save(ctx context.Context, d *identity_domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		return identity_domain.ErrDeviceNotFound
	}
	r.byID[d.ID] = d
	return nil
}

func (r *MemoryDeviceRepository) FindByID(ctx context.Context, id string) (*identity_domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.byID[id]
	if !ok {
		return nil, identity_domain.ErrDeviceNotFound
	}
	return d, nil
}

func (r *MemoryDeviceRepository) ListByVaultID(ctx context.Context, vaultID string) ([]*identity_domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var devices []*identity_domain.Device
	for _, d := range r.byID {
		if d.VaultID == vaultID {
			devices = append(devices, d)
		}
	}
	return devices, nil
}

func (r *MemoryDeviceRepository) Update(ctx context.Context, d *identity_domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[d.ID]; !ok {
		return identity_domain.ErrDeviceNotFound
	}
	r.byID[d.ID] = d
	return nil
}

// Ensure interface satisfaction at compile-time
var _ identity_domain.DeviceRepository = (*MemoryDeviceRepository)(nil)
