package app_config_commands

import app_config_domain "vault-app/internal/config/domain"

type CreateVaultConfigInput struct {
	VaultConfig app_config_domain.VaultConfigBeta	
}
type CreateVaultConfigOutput struct {
	VaultConfig app_config_domain.VaultConfigBeta
	Error error
}


type CreateVaultConfigCommand struct {
	VaultConfigRepository app_config_domain.VaultConfigRepository
}

func NewCreateVaultConfigCommand(repo app_config_domain.VaultConfigRepository) *CreateVaultConfigCommand {
	return &CreateVaultConfigCommand{VaultConfigRepository: repo}
}

func (c *CreateVaultConfigCommand) Execute(input CreateVaultConfigInput) CreateVaultConfigOutput {
	vaultConfig, err := c.VaultConfigRepository.Create(&input.VaultConfig)
	if err != nil {
		return CreateVaultConfigOutput{Error: err}
	}
	return CreateVaultConfigOutput{VaultConfig: *vaultConfig}
}