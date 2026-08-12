package identity_usecase

import (
	"context"

	identity_eventbus "vault-app/internal/identity/application"
	identity_domain "vault-app/internal/identity/domain"
)

type CreateDeviceRequest struct {
	VaultID   string
	PublicKey string
	KeyType   string
}

type CreateDeviceUseCase struct {
	repo identity_domain.DeviceRepository
	bus  identity_eventbus.EventBus
}

func NewCreateDeviceUseCase(repo identity_domain.DeviceRepository, bus identity_eventbus.EventBus) *CreateDeviceUseCase {
	return &CreateDeviceUseCase{
		repo: repo,
		bus:  bus,
	}
}

func (uc *CreateDeviceUseCase) Execute(ctx context.Context, req CreateDeviceRequest) (*identity_domain.Device, error) {
	dev, err := identity_domain.NewDevice(req.VaultID, req.PublicKey, req.KeyType)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, dev); err != nil {
		return nil, err
	}

	if uc.bus != nil {
		domainEvt := identity_domain.NewDeviceCreated(dev)
		_ = uc.bus.PublishDeviceCreated(ctx, identity_eventbus.DeviceCreated{
			DeviceID:   domainEvt.DeviceID,
			VaultID:    domainEvt.VaultID,
			PublicKey:  domainEvt.PublicKey,
			KeyType:    domainEvt.KeyType,
			OccurredAt: domainEvt.OccurredAt,
		})
	}

	return dev, nil
}

type GetDeviceUseCase struct {
	repo identity_domain.DeviceRepository
}

func NewGetDeviceUseCase(repo identity_domain.DeviceRepository) *GetDeviceUseCase {
	return &GetDeviceUseCase{repo: repo}
}

func (uc *GetDeviceUseCase) Execute(ctx context.Context, id string) (*identity_domain.Device, error) {
	if id == "" {
		return nil, identity_domain.ErrDeviceNotFound
	}
	return uc.repo.FindByID(ctx, id)
}

type ListDevicesUseCase struct {
	repo identity_domain.DeviceRepository
}

func NewListDevicesUseCase(repo identity_domain.DeviceRepository) *ListDevicesUseCase {
	return &ListDevicesUseCase{repo: repo}
}

func (uc *ListDevicesUseCase) Execute(ctx context.Context, vaultID string) ([]*identity_domain.Device, error) {
	if vaultID == "" {
		return nil, identity_domain.ErrDeviceVaultIDRequired
	}
	return uc.repo.ListByVaultID(ctx, vaultID)
}

type RevokeDeviceUseCase struct {
	repo identity_domain.DeviceRepository
	bus  identity_eventbus.EventBus
}

func NewRevokeDeviceUseCase(repo identity_domain.DeviceRepository, bus identity_eventbus.EventBus) *RevokeDeviceUseCase {
	return &RevokeDeviceUseCase{
		repo: repo,
		bus:  bus,
	}
}

func (uc *RevokeDeviceUseCase) Execute(ctx context.Context, id string) error {
	if id == "" {
		return identity_domain.ErrDeviceNotFound
	}

	dev, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if dev == nil {
		return identity_domain.ErrDeviceNotFound
	}

	if err := dev.Revoke(); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, dev); err != nil {
		return err
	}

	if uc.bus != nil {
		domainEvt := identity_domain.NewDeviceRevoked(dev)
		_ = uc.bus.PublishDeviceRevoked(ctx, identity_eventbus.DeviceRevoked{
			DeviceID:   domainEvt.DeviceID,
			VaultID:    domainEvt.VaultID,
			RevokedAt:  domainEvt.RevokedAt,
			OccurredAt: domainEvt.OccurredAt,
		})
	}

	return nil
}
