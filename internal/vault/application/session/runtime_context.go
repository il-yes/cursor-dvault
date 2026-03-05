// internal/vault/application/session/runtime_context.go
package vault_session

import (
	"time"
	app_config_domain "vault-app/internal/config/domain"
)

type RuntimeContext struct {
    AppConfig      app_config_domain.AppConfig
    UserConfig     app_config_domain.UserConfig
    SessionSecrets map[string]string
    WorkingBranch  string
}


func NewRuntimeContext() *RuntimeContext {
    return &RuntimeContext{
        AppConfig:      app_config_domain.AppConfig{},
        UserConfig:     app_config_domain.UserConfig{},
        SessionSecrets: map[string]string{},
        WorkingBranch:  "main",
    }
}

func NewRuntimeContextFrom(appConfig app_config_domain.AppConfig, userConfig app_config_domain.UserConfig, sessionSecrets map[string]string, workingBranch string) *RuntimeContext {
    return &RuntimeContext{
        AppConfig:      appConfig,
        UserConfig:     userConfig,
        SessionSecrets: sessionSecrets,
        WorkingBranch:  workingBranch,
    }
}
func (rc *RuntimeContext) GetAppConfig() app_config_domain.AppConfig {
	return rc.AppConfig
}
func (rc *RuntimeContext) GetUserConfig() app_config_domain.UserConfig {
	return rc.UserConfig
}
func (rc *RuntimeContext) GetSessionSecrets() map[string]string {
	return rc.SessionSecrets
}
func (rc *RuntimeContext) GetWorkingBranch() string {
	return rc.WorkingBranch
}
func (rc *RuntimeContext) SetAppConfig(appConfig app_config_domain.AppConfig) {
    rc.AppConfig = appConfig
}

func (rc *RuntimeContext) SetUserConfig(userConfig app_config_domain.UserConfig) {
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

// func (rc *RuntimeContext) ToFormerRuntimeContext() *models.VaultRuntimeContext {
// 	return &models.VaultRuntimeContext{
// 		CurrentUser:    rc.UserConfig,
// 		AppSettings:    rc.AppConfig,
// 		SessionSecrets: rc.SessionSecrets,
// 		WorkingBranch:  rc.WorkingBranch,
// 	}
// }
func (rc *RuntimeContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}