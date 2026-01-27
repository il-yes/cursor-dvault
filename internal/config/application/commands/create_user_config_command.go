package app_config_commands

import (
	app_config_domain "vault-app/internal/config/domain"
)


type CreateUserConfigCommandInput struct {
	UserConfig *app_config_domain.UserConfig
}


type CreateUserConfigCommandOutput struct {
	UserConfig *app_config_domain.UserConfig
}


type CreateUserConfigCommand struct {
	userConfigRepository app_config_domain.UserConfigRepository
}

func NewCreateUserConfigCommand(userConfigRepository app_config_domain.UserConfigRepository) *CreateUserConfigCommand {
	return &CreateUserConfigCommand{userConfigRepository: userConfigRepository}
}
	
func (c *CreateUserConfigCommand) Execute(input *CreateUserConfigCommandInput) (*CreateUserConfigCommandOutput, error) {
	if err := c.userConfigRepository.CreateUserConfig(input.UserConfig); err != nil {
		return nil, err
	}	
	return &CreateUserConfigCommandOutput{UserConfig: input.UserConfig}, nil
}	
