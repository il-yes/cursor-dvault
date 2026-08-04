package channel_domain


var (
	ErrInvalidName   = "Invalid Channel name"
	ErrRepositoryNil = "repository is nil"
	ErrSignatureMissing = "Signature is required"
	ErrChannelOwnerRequired = "owner id is required"
	ErrChannelIDRequired = "channel id is required"
	ErrChannelNameRequired = "channel name is required"
	ErrVaultIDRequired = "vault id is required"
	ErrChannelBusRequired = "Event bus is nil"
	ErrRequestRequired = "Request is nil"
	ErrRepositoryResponse = "channel repository returned nil response"
)