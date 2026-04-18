package vaults_domain

import "errors"

var (
	ErrVaultNotFound = errors.New("vault not found")
	ErrInvalidKey = errors.New("invalid key")
)
	