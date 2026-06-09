package app_config_commands

import (
	app_config_domain "vault-app/internal/config/domain"
	)

type CreateOnboardingCommand struct {
	UserID         string
	Packs          []string `json:"packs" gorm:"type:json"` // optional UI labels if you want
	UseCases       []string `json:"use_cases" gorm:"type:json"`       // optional UI labels if you want
	InstalledSeeds []string `json:"installed_seeds" gorm:"type:json"` // ["devsecops.incident.v1", ...]
    Completed      bool     `json:"completed" gorm:"column:completed"`
    PacksApplied bool     `json:"packs_applied" gorm:"column:packs_applied"`
}

type CreateOnboardingCommandHandler struct {
	onboardingConfigRepo app_config_domain.OnboardingConfigRepository
}

func NewCreateOnboardingCommandHandler(onboardingConfigRepo app_config_domain.OnboardingConfigRepository) *CreateOnboardingCommandHandler {
	return &CreateOnboardingCommandHandler{onboardingConfigRepo}
}

func (h *CreateOnboardingCommandHandler) Handle(onboardingConfig *app_config_domain.OnboardingConfig) error {
	// onboardingConfig := &app_config_domain.OnboardingConfig{
	// 	UserID:         cmd.UserID,
	// 	Packs:          cmd.Packs,
	// 	UseCases:       cmd.UseCases,
	// 	InstalledSeeds: cmd.InstalledSeeds,
	// 	Completed:      cmd.Completed,
	// 	PacksApplied:   cmd.PacksApplied,
	// }
	return h.onboardingConfigRepo.Create(onboardingConfig)
}

