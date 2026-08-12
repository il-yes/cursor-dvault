package identity_domain

import "time"

// Domain event
type UserRegistered struct {
	UserID    string
	IsAnonymous bool
	OccurredAt int64
}

func NewUserRegistered(u *User) UserRegistered {
	return UserRegistered{UserID: u.ID, IsAnonymous: u.IsAnonymous, OccurredAt: time.Now().UnixNano()}
}

type DeviceCreated struct {
	DeviceID   string
	VaultID    string
	PublicKey  string
	KeyType    string
	OccurredAt time.Time
}

func NewDeviceCreated(d *Device) DeviceCreated {
	return DeviceCreated{
		DeviceID:   d.ID,
		VaultID:    d.VaultID,
		PublicKey:  d.PublicKey,
		KeyType:    d.KeyType,
		OccurredAt: time.Now(),
	}
}

type DeviceRevoked struct {
	DeviceID   string
	VaultID    string
	RevokedAt  time.Time
	OccurredAt time.Time
}

func NewDeviceRevoked(d *Device) DeviceRevoked {
	revTime := time.Now()
	if d.RevokedAt != nil {
		revTime = *d.RevokedAt
	}
	return DeviceRevoked{
		DeviceID:   d.ID,
		VaultID:    d.VaultID,
		RevokedAt:  revTime,
		OccurredAt: time.Now(),
	}
}

