package channel_domain

import "errors"


var (
	ErrInvalidName   = errors.New("Invalid Channel name")
	ErrRepositoryNil = errors.New("repository is nil")
	ErrSignatureMissing = errors.New("Signature is required")
	ErrChannelOwnerRequired = errors.New("owner id is required")
	ErrChannelNotFound = errors.New("channel not found")
	ErrChannelIDRequired = errors.New("channel id is required")
	ErrChannelNameRequired = errors.New("channel name is required")
	ErrVaultIDRequired = errors.New("vault id is required")
	ErrChannelBusRequired = errors.New("Event bus is nil")
	ErrRequestRequired = errors.New("Request is nil")
	ErrRepositoryResponse = errors.New("channel repository returned nil response")
	ErrWorkspaceIDRequired = errors.New("workspace id is required")
	ErrChannelNotArchivable = errors.New("only active channels can be archived")
	ErrChannelNotModifiable = errors.New("channel is not modifiable")
	ErrSlotAlreadyExists    = errors.New("slot already exists")
	ErrSlotNotFound         = errors.New("slot not found")
	ErrSlotIdentityMismatch = errors.New("slot identity cannot be changed")
	ErrChannelRepositoryRequired = errors.New("channel repository is required")
	ErrSlotIDRequired       = errors.New("slot ID is required")

	ErrAssignmentIDRequired = errors.New("assignment ID is required")
	ErrAssignmentNotFound = errors.New("assignment not found")
	ErrAssignmentIdentityMismatch = errors.New("assignment identity cannot be changed")
	ErrAssignmentRepositoryRequired = errors.New("assignment repository is required")

)