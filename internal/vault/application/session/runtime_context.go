// internal/vault/application/session/runtime_context.go
package vault_session

import (
	app_config "vault-app/internal/config"
	"vault-app/internal/models"
)

type RuntimeContext struct {
    AppConfig      app_config.AppConfig
    UserConfig     app_config.UserConfig
    SessionSecrets map[string]string
    WorkingBranch  string
}


func NewRuntimeContext() *RuntimeContext {
    return &RuntimeContext{
        AppConfig:      app_config.AppConfig{},
        UserConfig:     app_config.UserConfig{},
        SessionSecrets: map[string]string{},
        WorkingBranch:  "main",
    }
}

func NewRuntimeContextFrom(appConfig app_config.AppConfig, userConfig app_config.UserConfig, sessionSecrets map[string]string, workingBranch string) *RuntimeContext {
    return &RuntimeContext{
        AppConfig:      appConfig,
        UserConfig:     userConfig,
        SessionSecrets: sessionSecrets,
        WorkingBranch:  workingBranch,
    }
}

func (rc *RuntimeContext) SetAppConfig(appConfig app_config.AppConfig) {
    rc.AppConfig = appConfig
}

func (rc *RuntimeContext) SetUserConfig(userConfig app_config.UserConfig) {
    rc.UserConfig = userConfig
}

func (rc *RuntimeContext) SetSessionSecrets(sessionSecrets map[string]string) {
    rc.SessionSecrets = sessionSecrets
}

func (rc *RuntimeContext) SetWorkingBranch(workingBranch string) {
    rc.WorkingBranch = workingBranch
}

func (rc *RuntimeContext) IsMultiActorMode() bool {
	return len(rc.AppConfig.Actors) > 1
}

func (rc *RuntimeContext) CurrentUserID() string {
	return rc.UserConfig.ID
}

func (rc *RuntimeContext) ToFormerRuntimeContext() *models.VaultRuntimeContext {
	return &models.VaultRuntimeContext{
		CurrentUser:    rc.UserConfig,
		AppSettings:    rc.AppConfig,
		SessionSecrets: rc.SessionSecrets,
		WorkingBranch:  rc.WorkingBranch,
	}
}