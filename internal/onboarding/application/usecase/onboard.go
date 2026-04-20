package onboarding_usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	app_config_commands "vault-app/internal/config/application/commands"
	app_config_domain "vault-app/internal/config/domain"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_domain "vault-app/internal/onboarding/domain"
	"vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// This use case orchestrates creating the user profiles, adding payment method, (creating subscription), and creating vault.
// It uses application ports (interfaces) to interact with other bounded contexts.

// ----------- Interfaces -----------
type IdentityRegisterPort interface {
	RegisterIdentity(ctx context.Context, req identity_usecase.RegisterRequest) (*identity_domain.User, error)
}

type BillingPort interface {
	AddPaymentMethod(ctx context.Context, userID string, method billing_domain.PaymentMethod, encryptedPayload string) (*billing_usecase.AddPaymentMethodResponse, error)
}

type SubscriptionPort interface {
	CreateSubscription(ctx context.Context, userID string, tier string) (string, error)
}

type VaultPort interface {
	CreateVault(v vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error)
}

type IdentityHandlerInterface interface {
	Registers(req identity_ui.OnboardRequest) (*identity_domain.User, error)
}

type BillingHandlerInterface interface {
	Onboard(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error)
}

type AppConfigHandlerInterface interface {
	GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error)
	InitAppConfig(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error)
	InitUserConfig(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error)
	GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error)
}

// ----------- Request -----------
type RegisterRequest struct {
	Email            string
	Password         string
	IsAnonymous      bool
	StellarPublicKey string
}

// ----------- UseCase Request -----------
type OnboardRequest struct {
	Identity             string
	VaultName            string
	Email                string
	Password             string
	IsAnonymous          bool
	Tier                 string
	PaymentMethod        string
	EncryptedPaymentData string
	SubscriptionID       string
	UserSubscriptionID   string
}

// ----------- UseCase Result -----------
type OnboardResult struct {
	UserID         string
	StellarKey     string
	SubscriptionID string
}

type OnboardUseCase struct {
	StellarService           StellarServiceInterface
	OnBoardingUserRepository UserServiceInterface
	Bus                      onboarding_application_events.OnboardingEventBus
	Logger                   Logger
	IdentityHandler          IdentityHandlerInterface
	BillingHandler           BillingHandlerInterface
	Vault                    VaultPort
	AppConfigHandler         AppConfigHandlerInterface
	AppStateRepo             onboarding_domain.AppStateRepository
	gormDb                   *gorm.DB
}

// ----------- UseCase Constructor -----------
func NewOnboardUseCase(
	v VaultPort,
	stellarService StellarServiceInterface,
	userService UserServiceInterface,
	eventBus onboarding_application_events.OnboardingEventBus,
	logger Logger,
	identityHandler IdentityHandlerInterface,
	billingHandler BillingHandlerInterface,
	appConfigHandler AppConfigHandlerInterface,
	appStateRepo onboarding_domain.AppStateRepository, // ✅ injected
) *OnboardUseCase {

	return &OnboardUseCase{
		Vault:                    v,
		StellarService:           stellarService,
		OnBoardingUserRepository: userService,
		Bus:                      eventBus,
		Logger:                   logger,
		IdentityHandler:          identityHandler,
		BillingHandler:           billingHandler,
		AppConfigHandler:         appConfigHandler,
		AppStateRepo:             appStateRepo,
	}
}

func (uc *OnboardUseCase) Execute(ctx context.Context, req OnboardRequest) (*OnboardResult, error) {
	if req.SubscriptionID == "" {
		return nil, errors.New("subscription ID is required")
	}
	if req.IsAnonymous {
		return nil, errors.New("anonymous requested but no stellar public key provided")
	}
	var secretKey string
	var err error

	// 1. ------------- Fetch onboarded account ------------------
	onboardUser, err := uc.OnBoardingUserRepository.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	uc.Logger.Info("onboardUser", onboardUser)

	// 2. ------------- Identity registration ------------------
	userIdentity, err := uc.IdentityHandler.Registers(identity_ui.OnboardRequest{
		Email:            req.Email,
		Password:         onboardUser.Password,
		IsAnonymous:      req.IsAnonymous,
		StellarPublicKey: onboardUser.StellarPublicKey,
	})
	if err != nil {
		return nil, err
	}
	utils.LogPretty("OnboardUseCase - Execute - userIdentity", userIdentity)

	// 3. ------------- App Configs creation ------------------
	configs, err := app_config_domain.InitConfigFromVault(userIdentity.ID, req.VaultName)
	if err != nil {
		uc.Logger.Error("OnboardUseCase - Execute - Failed to create app config: %v", err)
		return nil, err
	}
	// appConfig, err := uc.AppConfigHandler.InitAppConfig(&app_config_commands.CreateAppConfigCommandInput{
	// 	AppConfig: configs.App,
	// })
	// if err != nil {
	// 	uc.Logger.Error("OnboardUseCase - Execute - Failed to create app config: %v", err)
	// 	return nil, err
	// }
	// userConfig, err := uc.AppConfigHandler.InitUserConfig(&app_config_commands.CreateUserConfigCommandInput{
	// 	UserConfig: configs.User,
	// })
	// if err != nil {
	// 	uc.Logger.Error("OnboardUseCase - Execute - Failed to create user config: %v", err)
	// 	return nil, err
	// }
	
	configs.Subscription.ID = req.SubscriptionID
	configs.Subscription.UserID = req.UserSubscriptionID
	configs.Subscription.BaseVaultConfig = app_config_domain.BaseVaultConfig{
			ID:        req.SubscriptionID,
			UserID:    req.UserSubscriptionID,
			VaultName: req.VaultName,
	}
	configs.Subscription.Plan = req.Tier

	deviceName, err := uc.GetDeviceName()
	if err != nil {
		uc.Logger.Error("OnboardUseCase - Execute - Failed to get device name: %v", err)
		return nil, err
	}

	configs.Devices = []app_config_domain.DeviceConfig{
		{
			BaseVaultConfig: app_config_domain.BaseVaultConfig{
				ID:        uuid.NewString(),
				UserID:    onboardUser.ID,
				VaultName: req.VaultName,
			},
			DeviceID:   uuid.NewString(),
			DeviceName: deviceName,
		},
	}
	

	appConfigSaved, err := uc.AppConfigHandler.GetAppConfigByUserID(ctx, userIdentity.ID)
	if err != nil {
		uc.Logger.Error("OnboardUseCase - Execute - Failed to get app config: %v", err)
		return nil, err
	}
	if appConfigSaved == nil {
		uc.Logger.Error("OnboardUseCase - Execute - App config not found")
		return nil, errors.New("app config not found")
	}
	utils.LogPretty("OnboardUseCase - Execute - appConfig", appConfigSaved)
	userConfigSaved, err := uc.AppConfigHandler.GetUserConfigByUserID(userIdentity.ID)
	if err != nil {
		uc.Logger.Error("OnboardUseCase - Execute - Failed to get user config: %v", err)
		return nil, err
	}
	if userConfigSaved == nil {
		uc.Logger.Error("OnboardUseCase - Execute - User config not found")
		return nil, errors.New("user config not found")
	}
	utils.LogPretty("OnboardUseCase - Execute - userConfig", userConfigSaved)
	


	if uc.Vault == nil {
		uc.Logger.Error("OnboardUseCase - Execute - Vault is nil")
		return nil, errors.New("vault is nil")
	}
	fmt.Println("uc.Vault", uc.Vault)
	if appConfigSaved == nil {
		uc.Logger.Error("OnboardUseCase - Execute - App config is nil")
		return nil, errors.New("app config is nil")
	}
	if userIdentity == nil {
		uc.Logger.Error("OnboardUseCase - Execute - User identity is nil")
		return nil, errors.New("user identity is nil")
	}
	utils.LogPretty("OnboardUseCase - Execute - req", req)
	// 3. ------------- Vault creation ------------------
	result, err := uc.Vault.CreateVault(vault_commands.CreateVaultCommand{
		UserID:             userIdentity.ID,
		VaultName:          req.VaultName,
		Password:           req.Password,
		UserSubscriptionID: req.UserSubscriptionID,
		AppConfig:          *appConfigSaved,
		UserOnboarding:  onboardUser,
	})
	if err != nil {
		return nil, err
	}
	uc.Logger.Info("vault created", result)

	// 4. ------------- Billing registration optional) ------------------
	if req.PaymentMethod != "" {
		if _, err := uc.BillingHandler.Onboard(ctx, billing_ui_handlers.AddPaymentMethodRequest{
			UserID:           userIdentity.ID,
			Method:           req.PaymentMethod,
			EncryptedPayload: req.EncryptedPaymentData,
		}); err != nil {
			return nil, err
		}
	}

	// 5. ------------- (Optional event...) ------------------
	appState := onboarding_domain.NewAppState()
	appState.SetHasVault(true)
	if uc.AppStateRepo == nil {
		return nil, errors.New("appStateRepo is nil")
	}

	if err := uc.AppStateRepo.Save(appState); err != nil {
		return nil, err
	}

	return &OnboardResult{UserID: userIdentity.ID, StellarKey: secretKey, SubscriptionID: req.SubscriptionID}, nil
}


func (uc *OnboardUseCase) GetDeviceName() (string, error) {
	deviceName, err := os.Hostname()
	if err != nil || deviceName == "" {
		deviceName = "unknown"
	}

	// Trim and sanitize if you want shorter names
	deviceName = strings.ReplaceAll(deviceName, ".", "-")
	if len(deviceName) > 64 {
		deviceName = deviceName[:64]
	}
	return deviceName, nil
}