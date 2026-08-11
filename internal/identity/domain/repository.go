package identity_domain	

import "context"

// Repository interface for identity
type UserRepository interface {
	Save(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	FindByPublicKey(ctx context.Context, publicKey string) (*User, error)
}

type DeviceRepository interface {
	Save(ctx context.Context, d *Device) error
	FindByID(ctx context.Context, id string) (*Device, error)
	ListByVaultID(ctx context.Context, vaultID string) ([]*Device, error)
	Update(ctx context.Context, d *Device) error
}

