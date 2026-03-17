package app_config_commands

import (
	"errors"
	"vault-app/internal/config/domain"
)



type CreateSubscriptionConfigInput struct {
	SubscriptionConfig *app_config_domain.SubscriptionConfig	
}
type CreateSubscriptionConfigOutput struct {
	SubscriptionConfig *app_config_domain.SubscriptionConfig
	Error error
}



type CreateSubscriptionConfigCommand struct {
	SubscriptionConfigRepository app_config_domain.SubscriptionConfigRepository
}

func NewCreateSubscriptionConfigCommand(subscriptionConfigRepository app_config_domain.SubscriptionConfigRepository) *CreateSubscriptionConfigCommand {
	return &CreateSubscriptionConfigCommand{
		SubscriptionConfigRepository: subscriptionConfigRepository,
	}
}

func (c *CreateSubscriptionConfigCommand) Execute(input CreateSubscriptionConfigInput) *CreateSubscriptionConfigOutput {
	if input.SubscriptionConfig == nil {
		return &CreateSubscriptionConfigOutput{
			SubscriptionConfig: nil,
			Error:              errors.New("subscription config is required"),
		}
	}
	if c.SubscriptionConfigRepository == nil {
		return &CreateSubscriptionConfigOutput{
			SubscriptionConfig: nil,
			Error:              errors.New("subscription config repository is required"),
		}
	}

	if err := c.SubscriptionConfigRepository.Create(input.SubscriptionConfig); err != nil {
		return &CreateSubscriptionConfigOutput{
			SubscriptionConfig: nil,
			Error:              err,
		}
	}
	return &CreateSubscriptionConfigOutput{
		SubscriptionConfig: input.SubscriptionConfig,
		Error:              nil,
	}
}
