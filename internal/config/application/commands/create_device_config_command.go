package app_config_commands

import app_config_domain "vault-app/internal/config/domain"

type CreateDeviceConfigInput struct {
	Device app_config_domain.DeviceConfig	
}
type CreateDeviceConfigOutput struct {
	Device app_config_domain.DeviceConfig
	Error error
}

type CreateDeviceConfigCommand struct {
	DeviceConfigRepository app_config_domain.DeviceConfigRepository
}

func NewCreateDeviceConfigCommand(repo app_config_domain.DeviceConfigRepository) *CreateDeviceConfigCommand {
	return &CreateDeviceConfigCommand{DeviceConfigRepository: repo}
}

func (c *CreateDeviceConfigCommand) Execute(input CreateDeviceConfigInput) CreateDeviceConfigOutput {
	if err := c.DeviceConfigRepository.Create(&input.Device); err != nil {
		return CreateDeviceConfigOutput{Error: err}
	}
	return CreateDeviceConfigOutput{Device: input.Device}
}