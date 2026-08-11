package trustgroup_domain

import "errors"

var (
	ErrInvalidName   = errors.New("Invalid TrustGroup name")
	ErrRepositoryNil = errors.New("repository is nil")
	ErrSignatureMissing = errors.New("Signature is required")
	ErrTrustGroupOwnerRequired = errors.New("owner id is required")
	ErrTrustGroupIDRequired = errors.New("trustGroup id is required")
	ErrTrustGroupNameRequired = errors.New("trustGroup name is required")
	ErrVaultIDRequired = errors.New("vault id is required")
	ErrTrustGroupBusRequired = errors.New("Event bus is nil")
	ErrRequestRequired = errors.New("Request is nil")
	ErrRepositoryResponse = errors.New("trustGroup repository returned nil response")
	ErrChannelIDRequired = errors.New("channel id is required")
	ErrTrustGroupNotFound = errors.New("trust group not found")
	ErrTrustGroupMemberNotFound = errors.New("trust group member not found")
	ErrStaleKEKVersion = errors.New("stale KEK version")
	ErrDuplicateKeyEnvelope = errors.New("duplicate active key envelope for device")
	ErrKEKVersionRequired = errors.New("kek version is required")
	ErrDeviceIDRequired = errors.New("device id is required")
	ErrMemberIDRequired = errors.New("member id is required")
	ErrWrappedKEKRequired = errors.New("wrapped kek is required")
	ErrDeviceNotFound = errors.New("device not found")
	ErrDeviceRevoked = errors.New("device is revoked")
	ErrDeviceMemberMismatch = errors.New("device does not belong to expected member or vault")
	ErrMemberNotInTrustGroup = errors.New("member does not belong to trust group")
)


