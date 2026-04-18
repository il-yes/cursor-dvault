package onboarding_usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	app_config_commands "vault-app/internal/config/application/commands"
	app_config_domain "vault-app/internal/config/domain"
	identity_ui "vault-app/internal/identity/ui"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	"vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"
	vaults_domain "vault-app/internal/vault/domain"

	"vault-app/internal/logger/logger"

	"gorm.io/gorm"

	identity_app "vault-app/internal/identity/application"
	identity_domain "vault-app/internal/identity/domain"
)

// --------------------------------------------------------------------------------------------------
// FAKE
// --------------------------------------------------------------------------------------------------
// ============= fakeIdentity =======================================================
type fakeIdentity struct {
	called bool
	err    error
}

func (f *fakeIdentity) Registers(
	req identity_ui.OnboardRequest) (*identity_domain.User, error) {
	f.called = true
	if f.err != nil {
		return nil, f.err
	}
	return &identity_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}

// ============= fakeVault =======================================================
type fakeVault struct { // vault_commands.CreateVaultCommandHandler
	called                 bool
	err                    error
	updateFn               func(v *vaults_domain.Vault) error
	saveFn                 func(v *vaults_domain.Vault) error
	existingVault          *vaults_domain.Vault
	updateCalled           bool
	saveCalled             bool
	updateError            error
	saveError              error
	deleteCalled           bool
	deleteError            error
	CreateVaultFunc        func(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error)
	UpdateVaultFunc        func(v *vaults_domain.Vault) error
	SaveVaultFunc          func(v *vaults_domain.Vault) error
	GetLatestByUserIDFunc  func(userID string) (*vaults_domain.Vault, error)
	GetVaultFunc           func(string) (*vaults_domain.Vault, error)
	DeleteVaultFunc        func(string) error
	GetByUserIDAndNameFunc func(string, string) (*vaults_domain.Vault, error)
	UpdateVaultCIDFunc     func(vaultID, cid string) error
}

func (f *fakeVault) CreateVault(v vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error) {
	f.called = true
	if f.CreateVaultFunc != nil {
		return f.CreateVaultFunc(v)
	}
	return &vault_commands.CreateVaultResult{}, f.err
}

func (f *fakeVault) UpdateVault(v *vaults_domain.Vault) error {
	f.called = true
	return f.updateFn(v)
}

func (m *fakeVault) SaveVault(v *vaults_domain.Vault) error {
	if m.saveFn != nil {
		return m.saveFn(v)
	}
	m.existingVault = v
	return nil
}

func (f *fakeVault) GetLatestByUserID(userID string) (*vaults_domain.Vault, error) {
	if userID == "user-1" { //"test_user" {
		return &vaults_domain.Vault{
			UserID: userID,
			Name:   "test_vault_name",
		}, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (f *fakeVault) GetVault(string) (*vaults_domain.Vault, error) {
	panic("not used")
}
func (f *fakeVault) DeleteVault(string) error {
	f.deleteCalled = true
	return f.deleteError
}
func (f *fakeVault) GetByUserIDAndName(string, string) (*vaults_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vaults_domain.ErrVaultNotFound
}
func (f *fakeVault) UpdateVaultCID(vaultID, cid string) error {
	f.updateCalled = true
	return f.updateError
}

// ============= fakeBilling =======================================================
type fakeBilling struct {
	called bool
	err    error
}

func (f *fakeBilling) Onboard(ctx context.Context,
	req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
	f.called = true
	return &billing_ui_handlers.AddPaymentMethodResponse{}, f.err
}
func (f *fakeBilling) AddPaymentMethod(
	ctx context.Context,
	userID string,
	method string,
	payload string,
) (string, error) {
	f.called = true
	return "pay-123", f.err
}

// ============= fakeUserRepo =======================================================
type fakeUserRepo struct{}

func (f *fakeUserRepo) Create(u *onboarding_domain.User) (*onboarding_domain.User, error) {
	return u, nil
}
func (f *fakeUserRepo) Update(u *onboarding_domain.User) error {
	return nil
}
func (f *fakeUserRepo) FindByID(id string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserRepo) Delete(id string) error {
	return nil
}
func (f *fakeUserRepo) FindByEmail(email string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserRepo) GetByID(id string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserRepo) GetByEmail(email string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserRepo) List() ([]onboarding_domain.User, error) {
	return []onboarding_domain.User{
		{
			ID:    "user-123",
			Email: "user@example.com",
		},
	}, nil
}

// ============= fakeStellarService =======================================================
type fakeStellarService struct{}

func (f *fakeStellarService) GenerateKeypair() (string, string, error) {
	return "PUB", "SECRET", nil
}
func (f *fakeStellarService) CreateAccount(stellarPublicKey string) (*blockchain.CreateAccountRes, error) {
	return &blockchain.CreateAccountRes{}, nil
}
func (f *fakeStellarService) GetPublicKey() (string, error) {
	return "PUB", nil
}
func (f *fakeStellarService) CreateKeypair() (string, string, string, error) {
	return "PUB", "SECRET", "", nil
}

// ============= fakeUserService =======================================================
type fakeUserService struct{}

func (f *fakeUserService) Create(*onboarding_domain.User) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserService) FindByEmail(email string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}

// ============= fakeBus =======================================================
type fakeBus struct {
	called bool
}

func (f *fakeBus) PublishCreated(ctx context.Context, event onboarding_application_events.AccountCreatedEvent) error {
	f.called = true
	return nil
}
func (f *fakeBus) Publish(event onboarding_events.OnboardingEventBus) {
	f.called = true
}
func (f *fakeBus) PublishSubscriptionActivated(ctx context.Context, event onboarding_application_events.SubscriptionActivatedEvent) error {
	f.called = true
	return nil
}
func (f *fakeBus) SubscribeToAccountCreation(func(onboarding_application_events.AccountCreatedEvent)) error {
	f.called = true
	return nil
}
func (f *fakeBus) SubscribeToSubscriptionActivated(onboarding_application_events.SubscriptionActivatedEvent) error {
	f.called = true
	return nil
}
func (f *fakeBus) SubscribeToSubscriptionActivation(func(onboarding_application_events.SubscriptionActivatedEvent)) error {
	f.called = true
	return nil
}

// ============= fakeAppConfigHandler =======================================================
type fakeAppConfigHandler struct {
	InitAppConfigFunc         func(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error)
	InitUserConfigFunc        func(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error)
	GetAppConfigByUserIDFunc  func(ctx context.Context, userID string) (*app_config_commands.CreateAppConfigCommandOutput, error)
	GetUserConfigByUserIDFunc func(userID string) (*app_config_domain.UserConfig, error)
}

func (f *fakeAppConfigHandler) InitAppConfig(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error) {
	return &app_config_commands.CreateAppConfigCommandOutput{
		AppConfig: &app_config_domain.AppConfig{}, // ✅ CRITICAL
	}, nil
}

func (f *fakeAppConfigHandler) InitUserConfig(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error) {
	return &app_config_commands.CreateUserConfigCommandOutput{
		UserConfig: &app_config_domain.UserConfig{},
	}, nil
}

func (f *fakeAppConfigHandler) GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error) {
	return &app_config_domain.AppConfig{}, nil
}

func (f *fakeAppConfigHandler) GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error) {
	return &app_config_domain.UserConfig{}, nil
}

// ============= fakeAppStateRepo =======================================================
type fakeAppStateRepo struct {
	called                    bool
	SaveFunc                  func(state *onboarding_domain.AppState) error
	GetFunc                   func() (*onboarding_domain.AppState, error)
	UpdateFunc                func(appState *onboarding_domain.AppState) error
	GetAppConfigByUserIDFunc  func(ctx context.Context, userID string) (*app_config_domain.AppConfig, error)
	GetUserConfigByUserIDFunc func(userID string) (*app_config_domain.UserConfig, error)
}

func (f *fakeAppStateRepo) Get() (*onboarding_domain.AppState, error) {
	if f.GetFunc != nil {
		return f.GetFunc()
	}
	f.called = true
	return &onboarding_domain.AppState{}, nil
}
func (f *fakeAppStateRepo) Update(appState *onboarding_domain.AppState) error {
	if f.UpdateFunc != nil {
		return f.UpdateFunc(appState)
	}
	return nil
}
func (f *fakeAppStateRepo) Save(appState *onboarding_domain.AppState) error {
	if f.SaveFunc != nil {
		return f.SaveFunc(appState)
	}
	return nil
}
func (f *fakeAppStateRepo) GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error) {
	if f.GetAppConfigByUserIDFunc != nil {
		return f.GetAppConfigByUserIDFunc(ctx, userID)
	}
	return &app_config_domain.AppConfig{}, nil
}
func (f *fakeAppStateRepo) GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error) {
	if f.GetUserConfigByUserIDFunc != nil {
		return f.GetUserConfigByUserIDFunc(userID)
	}
	return &app_config_domain.UserConfig{}, nil
}
func (f *fakeAppStateRepo) InitAppConfig(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error) {
	return &app_config_commands.CreateAppConfigCommandOutput{
		AppConfig: &app_config_domain.AppConfig{}, // ✅ CRITICAL
	}, nil
}
func (f *fakeAppStateRepo) InitUserConfig(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error) {
	return &app_config_commands.CreateUserConfigCommandOutput{
		UserConfig: &app_config_domain.UserConfig{},
	}, nil
}

func getFullConfig(userID string) app_config.StorageConfig {
	return app_config_domain.DefaultStorageConfig()

	// return app_config_domain.AppConfig{
	// 	UserID:       userID,
	// 	RepoID:       "my-repo-id",
	// 	Branch:       "DefaultBranch",
	// 	DefaultPhase: "DefaultDefaultPhase",
	// 	VaultSettings: app_config_domain.VaultConfig{
	// 		MaxEntries:       app_config_domain.DefaultMaxEntries,
	// 		AutoSyncEnabled:  true,
	// 		EncryptionScheme: "AES-256-GCM",
	// 	},
	// 	Blockchain: app_config_domain.BlockchainConfig{
	// 		Stellar: app_config_domain.StellarConfig{
	// 			Network:    "testnet",
	// 			HorizonURL: "https://horizon-testnet.stellar.org",
	// 			Fee:        100,
	// 		},
	// 		IPFS: app_config_domain.IPFSConfig{
	// 			APIEndpoint: "http://localhost:5001",
	// 			GatewayURL:  "https://ipfs.io/ipfs/",
	// 		},
	// 	},
	// 	Storage: app_config.StorageConfig{
	// 		Mode: app_config.StorageLocal, // ← production default

	// 		LocalIPFS: app_config.IPFSConfig{
	// 			APIEndpoint: "http://localhost:5001",
	// 			GatewayURL:  "https://ipfs.io/ipfs/",
	// 		},

	// 		PrivateIPFS: app_config.IPFSConfig{
	// 			APIEndpoint: "http://192.168.1.10:5001",
	// 			GatewayURL:  "http://192.168.1.10:8080/ipfs/",
	// 		},

	// 		Cloud: app_config.CloudConfig{
	// 			BaseURL: "https://ankhora.io/back",
	// 		},

	// 		EnterpriseS3: app_config.S3Config{
	// 			Region:   "us-east-1",
	// 			Bucket:   "ankhora-enterprise",
	// 			Endpoint: "https://s3.us-east-1.amazonaws.com",
	// 		},
	// 	},
	// }

}

// // ======= fakeInitVaultHandler ===============
type fakeInitVaultHandler struct {
	result    *vault_commands.InitializeVaultResult
	err       error
	called    bool
	executeFn func(cmd vault_commands.InitializeVaultCommand) (*vault_commands.InitializeVaultResult, error)
}

func (f *fakeInitVaultHandler) Execute(cmd vault_commands.InitializeVaultCommand) (*vault_commands.InitializeVaultResult, error) {
	f.called = true
	return f.result, f.err
}

// ======= fakeIdentity ===============
type MockIdentity struct {
	RegistersFunc func(req identity_ui.OnboardRequest) (*identity_domain.User, error)
}

func (m *MockIdentity) Registers(req identity_ui.OnboardRequest) (*identity_domain.User, error) {
	return m.RegistersFunc(req)
}

// ======= MockBilling ===============
type MockBilling struct {
	OnboardFunc func(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error)
}

func (m *MockBilling) Onboard(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
	return m.OnboardFunc(ctx, req)
}

// ======= MockVaultService ===============
// type MockVaultService struct { 	// vault_commands.CreateVaultCommandHandler
// 	CreateVaultFunc func(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error)
// }

// func (m *MockVaultService) CreateVault(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error) {
// 	if m.CreateVaultFunc != nil {
// 		return m.CreateVaultFunc(cmd)
// 	}
// 	return nil, fmt.Errorf("CreateVaultFunc not implemented")
// }

func GetConfig(userID string, vaultName string) (*app_config_domain.Config, error) {
	res, err := app_config_domain.InitConfigFromVault(userID, vaultName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ======= fakeCreateIPFSPayloadHandler ===============
type fakeCreateIPFSPayloadHandler struct {
	ipfs      vault_commands.IpfsServiceInterface
	result    *vault_commands.CreateIPFSPayloadCommandResult
	err       error
	called    bool
	executeFn func(ctx context.Context, vc app_config_domain.VaultContext, cmd vault_commands.CreateIPFSPayloadCommand) (*vault_commands.CreateIPFSPayloadCommandResult, error)
}

func (f *fakeCreateIPFSPayloadHandler) Execute(ctx context.Context, vc app_config_domain.VaultContext, cmd vault_commands.CreateIPFSPayloadCommand) (*vault_commands.CreateIPFSPayloadCommandResult, error) {
	f.called = true
	// return f.result, f.err
	return f.executeFn(ctx, vc, cmd)
}
func (f *fakeCreateIPFSPayloadHandler) SetIpfsService(i vault_commands.IpfsServiceInterface) {
	f.ipfs = i
}

type OnboardTestBuilder struct {
	t *testing.T

	// deps
	vaultHandler    *fakeVault
	stellar         *MockStellarService
	userRepo        *MockUserService
	bus             *MockBus
	identity        onboarding_usecase.IdentityHandlerInterface
	billing         onboarding_usecase.BillingHandlerInterface
	appConfig       *fakeAppConfigHandler
	appStateRepo    *fakeAppStateRepo
	CreateVaultFunc func(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error)

	uc *onboarding_usecase.OnboardUseCase
}

func (b *OnboardTestBuilder) WithCreateVaultFunc(fn func(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error)) *OnboardTestBuilder {
	b.CreateVaultFunc = fn
	return b
}
func (b *OnboardTestBuilder) WithIdentity(i onboarding_usecase.IdentityHandlerInterface) *OnboardTestBuilder {
	b.identity = i
	b.uc.IdentityHandler = i
	return b
}
func (b *OnboardTestBuilder) WithBilling(i onboarding_usecase.BillingHandlerInterface) *OnboardTestBuilder {
	b.billing = i
	b.uc.BillingHandler = i
	return b
}

//	func (b *OnboardTestBuilder) WithAppConfig(i onboarding_usecase.AppConfigHandlerInterface) *OnboardTestBuilder {
//		b.appConfig = i
//		b.uc.AppConfigHandler = i
//		return b
//	}
//
//	func (b *OnboardTestBuilder) WithAppStateRepo(i onboarding_usecase.AppStateRepositoryInterface) *OnboardTestBuilder {
//		b.appStateRepo = i
//		b.uc.AppStateRepository = i
//		return b
//	}
func (b *OnboardTestBuilder) WithVaultHandler(i *fakeVault) *OnboardTestBuilder {
	b.vaultHandler = i
	b.uc.Vault = i
	return b
}
func NewOnboardTestBuilder(t *testing.T) *OnboardTestBuilder {
	t.Helper()
	userId := "user-123"
	// vaultName := "vault-1"
	// storageConfig, err := GetConfig(userId, vaultName)
	// if err != nil {
	// 	t.Fatalf("failed to get config: %v", err)
	// }

	// ---------- DEFAULT SAFE MOCKS ----------
	appStateRepo := &fakeAppStateRepo{
		SaveFunc: func(state *onboarding_domain.AppState) error {
			return nil
		},
	}

	// ✅ User repo
	userRepo := &MockUserService{
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return &onboarding_domain.User{
				ID:       "onboard-1",
				Email:    email,
				Password: "hashed",
			}, nil
		},
	}

	// ✅ Identity
	identity := &MockIdentity{
		RegistersFunc: func(req identity_ui.OnboardRequest) (*identity_domain.User, error) {
			return &identity_domain.User{
				ID:    userId,
				Email: req.Email,
			}, nil
		},
	}

	// ✅ VaultHandler
	vaultHandler := &fakeVault{
		CreateVaultFunc: func(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error) {
			return &vault_commands.CreateVaultResult{
				Vault: &vaults_domain.Vault{ID: "vault-1"},
			}, nil
		},
	}

	// ✅ AppConfig (CRITICAL FIX)
	appConfig := &fakeAppConfigHandler{
		InitAppConfigFunc: func(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error) {
			return &app_config_commands.CreateAppConfigCommandOutput{
				AppConfig: &app_config_domain.AppConfig{},
			}, nil
		},
		InitUserConfigFunc: func(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error) {
			return &app_config_commands.CreateUserConfigCommandOutput{
				UserConfig: &app_config_domain.UserConfig{},
			}, nil
		},
		GetAppConfigByUserIDFunc: func(ctx context.Context, userID string) (*app_config_commands.CreateAppConfigCommandOutput, error) {
			return &app_config_commands.CreateAppConfigCommandOutput{
				AppConfig: &app_config_domain.AppConfig{}, // 🔥 NEVER NIL
			}, nil
		},
		GetUserConfigByUserIDFunc: func(userID string) (*app_config_domain.UserConfig, error) {
			return &app_config_domain.UserConfig{}, nil
		},
	}

	// ✅ Billing
	billing := &MockBilling{
		OnboardFunc: func(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
			return &billing_ui_handlers.AddPaymentMethodResponse{}, nil
		},
	}

	// ✅ Bus
	bus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	// ✅ Stellar (not always used but safe)
	stellar := &MockStellarService{}

	uc := onboarding_usecase.NewOnboardUseCase(
		vaultHandler,
		stellar,
		userRepo,
		bus,
		&logger.Logger{},
		identity,
		billing,
		appConfig,
		appStateRepo,
	)

	return &OnboardTestBuilder{
		t:            t,
		vaultHandler: vaultHandler,
		stellar:      stellar,
		userRepo:     userRepo,
		bus:          bus,
		identity:     identity,
		billing:      billing,
		appConfig:    appConfig,
		appStateRepo: appStateRepo,
		uc:           uc,
	}
}

/*
	------------------------------
	  Fake ID Generator

--------------------------------
*/
func fakeIDGen() string {
	return "test-user-id"
}

/*
	------------------------------
	  Fake In-Memory Repo

--------------------------------
*/
type fakeRepo struct {
	usersByID    map[string]*identity_domain.User
	usersByEmail map[string]*identity_domain.User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		usersByID:    map[string]*identity_domain.User{},
		usersByEmail: map[string]*identity_domain.User{},
	}
}
func (r *fakeRepo) Save(ctx context.Context, u *identity_domain.User) error {
	r.usersByID[u.ID] = u
	if u.Email != "" {
		r.usersByEmail[u.Email] = u
	}
	return nil
}
func (r *fakeRepo) Update(ctx context.Context, u *identity_domain.User) error {
	return nil
}
func (r *fakeRepo) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	return r.usersByID[id], nil
}
func (r *fakeRepo) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	return r.usersByEmail[email], nil
}

func (r *fakeRepo) FindByPublicKey(ctx context.Context, publicKey string) (*identity_domain.User, error) {
	return r.usersByEmail[publicKey], nil
}

/*
-------------------------------

	Fake Event Bus

--------------------------------
*/
type fakeEventBus struct {
	published       []identity_app.UserRegistered
	publishedLogins []identity_app.UserLoggedIn
}

func newFakeEventBus() *fakeEventBus {
	return &fakeEventBus{
		published:       []identity_app.UserRegistered{},
		publishedLogins: []identity_app.UserLoggedIn{},
	}
}
func (b *fakeEventBus) PublishUserRegistered(ctx context.Context, e identity_app.UserRegistered) error {
	b.published = append(b.published, e)
	return nil
}
func (b *fakeEventBus) SubscribeToUserRegistered(handler identity_app.UserRegisteredHandler) error {
	return nil
}
func (b *fakeEventBus) PublishUserLoggedIn(ctx context.Context, e identity_app.UserLoggedIn) error {
	return nil
}
func (b *fakeEventBus) SubscribeToUserLoggedIn(handler identity_app.UserLoggedInHandler) error {
	return nil
}

// --------------------------------------------------------------------------------------------------
// TESTS
// --------------------------------------------------------------------------------------------------
func TestOnboardUseCase_Success(t *testing.T) {
	builder := NewOnboardTestBuilder(t)

	ctx := context.Background()

	res, err := builder.uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email:                "user@example.com",
		VaultName:            "my-vault",
		Password:             "password",
		IsAnonymous:          false,
		Identity:             "team",
		Tier:                 "pro",
		SubscriptionID:       "sub-123",
		UserSubscriptionID:   "user-sub-1",
		PaymentMethod:        "stripe",
		EncryptedPaymentData: "encrypted",
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if res.UserID != "user-123" {
		t.Fatalf("wrong userID: %s", res.UserID)
	}
}

func TestOnboardUseCase_AnonymousMissingKey(t *testing.T) {
	ctx := context.Background()

	builder := NewOnboardTestBuilder(t)
	appStateRepo := &fakeAppStateRepo{
		SaveFunc: func(state *onboarding_domain.AppState) error {
			return nil
		},
	}
	builder.uc.AppStateRepo = appStateRepo

	res, err := builder.uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email:                "x@x.com",
		VaultName:            "my-vault",
		Password:             "",
		IsAnonymous:          true,
		Identity:             "team",
		Tier:                 "pro",
		SubscriptionID:       "sub-123",
		UserSubscriptionID:   "user-sub-1",
		PaymentMethod:        "stripe",
		EncryptedPaymentData: "encrypted",
	})

	if err == nil {
		t.Fatalf("expected error")
	}
	utils.LogPretty("res", res)
}

func TestOnboardUseCase_IdentityFails(t *testing.T) {
	ctx := context.Background()

	builder := NewOnboardTestBuilder(t).
		WithIdentity(&MockIdentity{
			RegistersFunc: func(req identity_ui.OnboardRequest) (*identity_domain.User, error) {
				return nil, errors.New("identity-fail")
			},
		})

	// 🔥 inject into use case
	builder.uc.IdentityHandler = builder.identity

	builder.vaultHandler.CreateVaultFunc = func(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error) {
		t.Fatalf("vault should NOT be called")
		return nil, nil
	}

	res, err := builder.uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email:              "x@x.com",
		VaultName:          "my-vault",
		Password:           "password",
		IsAnonymous:        false,
		Identity:           "team",
		Tier:               "pro",
		SubscriptionID:     "sub-123",
		UserSubscriptionID: "user-sub-1",
	})

	if err == nil || err.Error() != "identity-fail" {
		t.Fatalf("expected identity-fail, got %v", err)
	}

	if res != nil {
		t.Fatalf("expected nil result")
	}
}

func TestOnboardUseCase_VaultFails(t *testing.T) {
	ctx := context.Background()

	builder := NewOnboardTestBuilder(t).
		WithVaultHandler(&fakeVault{
			CreateVaultFunc: func(cmd vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error) {
				return nil, errors.New("vault-fail")
			},
		})

	builder.uc.Vault = builder.vaultHandler

	res, err := builder.uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email:                "x@x.com",
		VaultName:            "my-vault",
		Password:             "password",
		IsAnonymous:          false,
		Identity:             "team",
		Tier:                 "pro",
		SubscriptionID:       "sub-123",
		UserSubscriptionID:   "user-sub-1",
		PaymentMethod:        "stripe",
		EncryptedPaymentData: "encrypted",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "vault-fail") {
		t.Fatalf("unexpected error: %v", err)
	}
	utils.LogPretty("res", res)
}
func TestOnboardUseCase_BillingFails(t *testing.T) {
	ctx := context.Background()

	builder := NewOnboardTestBuilder(t).
		WithBilling(&MockBilling{
			OnboardFunc: func(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
				return nil, errors.New("billing-fail")
			},
		})

	builder.uc.BillingHandler = builder.billing

	builder.appStateRepo.SaveFunc = func(state *onboarding_domain.AppState) error {
		t.Fatalf("appStateRepo save should NOT be called")
		return nil
	}

	res, err := builder.uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email:                "x@x.com",
		VaultName:            "my-vault",
		Password:             "password",
		IsAnonymous:          false,
		Identity:             "team",
		Tier:                 "pro",
		SubscriptionID:       "sub-123",
		UserSubscriptionID:   "user-sub-1",
		PaymentMethod:        "stripe",
		EncryptedPaymentData: "encrypted",
	})
	if err == nil || err.Error() != "billing-fail" {
		t.Fatalf("expected billing-fail")
	}
	utils.LogPretty("res", res)
}
