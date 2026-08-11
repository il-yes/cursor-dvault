package identity_domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	DeviceStatusActive  = "active"
	DeviceStatusRevoked = "revoked"
)

const (
	DeviceKeyTypeEd25519   = "ed25519"
	DeviceKeyTypeSecp256k1 = "secp256k1"
	DeviceKeyTypeRSA       = "rsa"
)

type Device struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	VaultID   string     `json:"vault_id" gorm:"index;not null"`
	PublicKey string     `json:"public_key" gorm:"not null"`
	KeyType   string     `json:"key_type" gorm:"not null"`
	Status    string     `json:"status" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func (Device) TableName() string {
	return "identity_devices"
}

func NewDevice(vaultID, publicKey, keyType string) (*Device, error) {
	if vaultID == "" {
		return nil, ErrDeviceVaultIDRequired
	}
	if publicKey == "" {
		return nil, ErrDevicePublicKeyRequired
	}
	if keyType == "" {
		return nil, ErrDeviceKeyTypeRequired
	}

	return &Device{
		ID:        uuid.New().String(),
		VaultID:   vaultID,
		PublicKey: publicKey,
		KeyType:   keyType,
		Status:    DeviceStatusActive,
		CreatedAt: time.Now(),
		RevokedAt: nil,
	}, nil
}

func (d *Device) IsActive() bool {
	return d.Status == DeviceStatusActive && d.RevokedAt == nil
}

func (d *Device) Revoke() error {
	if !d.IsActive() || d.Status == DeviceStatusRevoked || d.RevokedAt != nil {
		return ErrDeviceAlreadyRevoked
	}
	now := time.Now()
	d.Status = DeviceStatusRevoked
	d.RevokedAt = &now
	return nil
}