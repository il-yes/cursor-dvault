package identity_domain

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
)
