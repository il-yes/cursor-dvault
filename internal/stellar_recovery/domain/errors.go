package stellar_recovery_domain

import "errors"

var (
	ErrNotFound     = errors.New("stellar_recovery: not found")
	ErrInvalidState = errors.New("stellar_recovery: invalid state")
	ErrValidation   = errors.New("stellar_recovery: validation failed")
	ErrInvalidKey   = errors.New("stellar_recovery: invalid key")
	ErrKeyVerification = errors.New("stellar_recovery: key verification failed")
	ErrUserNotFound   = errors.New("stellar_recovery: user not found")
	ErrVaultNotFound  = errors.New("stellar_recovery: vault not found")
	ErrKeyAlreadyUsed = errors.New("stellar_recovery: key already used")
)