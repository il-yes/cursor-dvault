package trustgroup_ports

import "context"

type DeviceSummary struct {
	ID        string
	VaultID   string
	PublicKey string
	KeyType   string
	Status    string
	IsActive  bool
}

type DeviceResolver interface {
	GetDevice(ctx context.Context, deviceID string) (*DeviceSummary, error)
	ListActiveDevices(ctx context.Context, memberID string) ([]DeviceSummary, error)
}
