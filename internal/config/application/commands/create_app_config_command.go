package app_config_commands

import (
	app_config_domain "vault-app/internal/config/domain"
)


type CreateAppConfigCommandInput struct {
	AppConfig *app_config_domain.AppConfig
}


type CreateAppConfigCommandOutput struct {
	AppConfig *app_config_domain.AppConfig
}


type CreateAppConfigCommand struct {
	appConfigRepository app_config_domain.AppConfigRepository
}

func NewCreateAppConfigCommand(appConfigRepository app_config_domain.AppConfigRepository) *CreateAppConfigCommand {
	return &CreateAppConfigCommand{appConfigRepository: appConfigRepository}
}
	
func (c *CreateAppConfigCommand) Execute(input *CreateAppConfigCommandInput) (*CreateAppConfigCommandOutput, error) {
	if err := c.appConfigRepository.CreateAppConfig(input.AppConfig); err != nil {
		return nil, err
	}	
	return &CreateAppConfigCommandOutput{AppConfig: input.AppConfig}, nil
}	
