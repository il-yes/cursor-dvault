package identity_domain

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")

    ErrDeviceVaultIDRequired   = errors.New("device vault ID is required")
    ErrDevicePublicKeyRequired = errors.New("device public key is required")
    ErrDeviceKeyTypeRequired   = errors.New("device key type is required")
    ErrDeviceAlreadyRevoked    = errors.New("device is already revoked")
    ErrDeviceNotFound          = errors.New("device not found")
)

