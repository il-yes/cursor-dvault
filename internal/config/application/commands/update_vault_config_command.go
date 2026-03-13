package app_config_commands

import (
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/logger/logger"
)

type UpdateVaultConfigCommandInput struct {
	VaultConfig *app_config_domain.VaultConfigBeta
}

type UpdateVaultConfigCommandOutput struct {
	VaultConfig *app_config_domain.VaultConfigBeta
}

type UpdateVaultConfigCommand struct {
	vaultConfigRepository    app_config_domain.VaultConfigRepository
	CreateVaultConfigCommand *CreateVaultConfigCommand
	Logger                   *logger.Logger
}

func NewUpdateVaultConfigCommand(
	vaultConfigRepository app_config_domain.VaultConfigRepository,
	createVaultConfigCommand *CreateVaultConfigCommand,
	logger *logger.Logger,
) *UpdateVaultConfigCommand {

	return &UpdateVaultConfigCommand{
		vaultConfigRepository:    vaultConfigRepository,
		CreateVaultConfigCommand: createVaultConfigCommand,
		Logger:                   logger,
	}
}

func (c *UpdateVaultConfigCommand) Execute(input *UpdateVaultConfigCommandInput) (*UpdateVaultConfigCommandOutput, error) {
	if err := c.vaultConfigRepository.Update(input.VaultConfig.ID, input.VaultConfig); err != nil {
		return nil, err
	}
	return &UpdateVaultConfigCommandOutput{VaultConfig: input.VaultConfig}, nil
}

func (c *UpdateVaultConfigCommand) CreateIfNotExists(vaultConfig *app_config_domain.VaultConfigBeta) error {
	c.Logger.LogPretty("UpdateVaultConfigCommand: CreateIfNotExists - vaultConfig", vaultConfig)
	// - Check if existing ---------------------------------------------------------------
	existingVaultCfg, err := c.vaultConfigRepository.FindByUserIDAndVaultName(vaultConfig.UserID, vaultConfig.VaultName)
	if err != nil {
		c.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to get vault config: %v", err)
	}

	if existingVaultCfg.ID == "" {
		// - Create new vault config ------------------------------------------------------
		c.Logger.Info("AppConfigHandler: OnChangingSettings - vault config not found => create new one")
		output := c.CreateVaultConfigCommand.Execute(CreateVaultConfigInput{
			VaultConfig: *vaultConfig,
		})
		if output.Error != nil {
			c.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to create vault config: %v", output.Error)
			return output.Error
		}
		c.Logger.LogPretty("AppConfigHandler: OnChangingSettings - vaultCfg created", output.VaultConfig)
	} else {
		// - Update existing vault config --------------------------------------------------
		c.Logger.Info("AppConfigHandler: OnChangingSettings - vault config found => update")
		_, err := c.Execute(&UpdateVaultConfigCommandInput{
			VaultConfig: &existingVaultCfg,
		})
		return err
	}
	return nil
}
