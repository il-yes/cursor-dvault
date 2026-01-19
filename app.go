package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	utils "vault-app/internal"
	share_application "vault-app/internal/application/use_cases"
	"vault-app/internal/auth"
	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_persistence "vault-app/internal/auth/infrastructure/persistence"
	auth_ui "vault-app/internal/auth/ui"
	billing_infrastructure_eventbus "vault-app/internal/billing/infrastructure/eventbus"
	billing_ui "vault-app/internal/billing/ui"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	app_config_ui "vault-app/internal/config/ui"
	share_domain "vault-app/internal/domain/shared"
	"vault-app/internal/driver"
	"vault-app/internal/handlers"
	identity_commands "vault-app/internal/identity/application/commands"
	identity_domain "vault-app/internal/identity/domain"
	identity_infrastructure_eventbus "vault-app/internal/identity/infrastructure/eventbus"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	onboarding_infrastructure_eventbus "vault-app/internal/onboarding/infrastructure/eventbus"
	onboarding_persistence "vault-app/internal/onboarding/infrastructure/persistence"
	onboarding_ui_wails "vault-app/internal/onboarding/ui/wails"
	"vault-app/internal/registry"
	shared "vault-app/internal/shared/stellar"
	stellar_recovery_usecase "vault-app/internal/stellar_recovery/application/use_case"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
	stellar "vault-app/internal/stellar_recovery/infrastructure"
	"vault-app/internal/stellar_recovery/infrastructure/events"
	stellar_recovery_persistence "vault-app/internal/stellar_recovery/infrastructure/persistence"
	"vault-app/internal/stellar_recovery/infrastructure/token"
	stellar_recovery_ui_api "vault-app/internal/stellar_recovery/ui/api"
	payments "vault-app/internal/stripe"
	subscription_usecase "vault-app/internal/subscription/application/usecase"
	subscription_domain "vault-app/internal/subscription/domain"
	subscription_infrastructure "vault-app/internal/subscription/infrastructure"
	subscription_infrastructure_eventbus "vault-app/internal/subscription/infrastructure/eventbus"
	subscription_persistence "vault-app/internal/subscription/infrastructure/persistence"
	subscription_ui_wails "vault-app/internal/subscription/ui/wails"
	"vault-app/internal/tracecore"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"
	vault_ui "vault-app/internal/vault/ui"

	// "vault-app/internal/logger/logger"
	"vault-app/internal/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
	"github.com/stripe/stripe-go/v74"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type CoreApp interface {
	SignIn(req handlers.LoginRequest) (*handlers.LoginResponse, error)
	SignUp(setup handlers.OnBoarding) (*handlers.OnBoardingResponse, error)
	// etc...
}

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	Domain string
	// Jwt auth
	auth        auth.Auth
	JWTSecret   string
	JWTIssuer   string
	JWTAudience string
	APIKey      string
	ANCHORA_SECRET string
}

type App struct {
	config   config
	Logger   logger.Logger
	version  string
	DB       models.DBModel
	ctx      context.Context
	sessions map[string]*models.VaultSession
	NowUTC   func() string

	// Core handlers
	AppConfigHandler          *app_config_ui.AppConfigHandler
	Auth                      *handlers.AuthHandler
	AuthHandler               *auth_ui.AuthHandler
	ConnectWithStellarHandler *stellar_recovery_ui_api.StellarRecoveryHandler
	EntryRegistry             *registry.EntryRegistry
	Identity                  *identity_ui.IdentityHandler
	OnBoardingHandler         *onboarding_ui_wails.OnBoardingHandler
	StellarRecoveryHandler    *stellar_recovery_ui_api.StellarRecoveryHandler
	SubscriptionHandler       *subscription_ui_wails.SubscriptionHandler
	Vault                     *vault_ui.VaultHandler
	Vaults                    *handlers.VaultHandler

	// New: Global state
	// RuntimeContext *models.VaultRuntimeContext
	RuntimeContext *vault_session.RuntimeContext
	cancel         context.CancelFunc
}

// NewApp creates a new App instance (required by Wails)
func NewApp() *App {
	startTime := time.Now()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := loadConfig()
	utils.LogPretty("cfg", cfg)

	// Use auto-init from env
	appLogger := logger.NewFromEnv()
	appLogger.Info("🚀 Starting D-Vault initialization...")

	// Pick DSN
	dsn := cfg.db.dsn
	if dsn == "" {
		dsn = "sqlite3.db"
	}
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	cfg.stripe.key = os.Getenv("STRIPE_SECRET")
	stripe.Key = os.Getenv("STRIPE_SECRET")
	appLogger.Info("✅ Stripe Key from env: ", stripe.Key)
	appLogger.Info("✅ Stripe Key hardcoded: ", stripe.Key)

	// read from command line
	// flag.StringVar(&dsn, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	// flag.StringVar(&cfg.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	// flag.StringVar(&cfg.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	// flag.StringVar(&cfg.JWTAudience, "jwt-audience", "example.com", "signing audience")
	// flag.StringVar(&cfg.Domain, "domain", "example.com", "domain")
	// flag.Parse()

	// database
	db, err := driver.InitDatabase(dsn, *appLogger)
	if err != nil {
		appLogger.Error("❌ Failed to init DB: %v", err)
		os.Exit(1)
	}
	appLogger.Info("✅ Local DB ready")
	cfg.auth = auth.Auth{
		Issuer:        cfg.JWTIssuer,
		Audience:      cfg.JWTAudience,
		Secret:        cfg.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
	}
	authV2 := auth_domain.Auth{
		Issuer:        cfg.JWTIssuer,
		Audience:      cfg.JWTAudience,
		Secret:        cfg.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
	}

	// Init services
	appLogger.Info("🔧 Initializing IPFS client...")
	ipfs := blockchain.NewIPFSClient("localhost:5001")
	appLogger.Info("✅ IPFS client initialized (connection will be tested on first use)")

	sessions := make(map[string]*models.VaultSession)
	sessionsV2 := make(map[string]*vault_session.Session)

	ctx, cancel := context.WithCancel(context.Background())
	runtimeCtx := &vault_session.RuntimeContext{
		AppConfig: app_config.AppConfig{
			// Load from file/env or defaults
			Branch:           "main",
			EncryptionPolicy: "AES-256-GCM",
			Blockchain: app_config.BlockchainConfig{
				Stellar: app_config.StellarConfig{
					Network:    "testnet",
					HorizonURL: "https://horizon-testnet.stellar.org",
					Fee:        100,
				},
				IPFS: app_config.IPFSConfig{
					APIEndpoint: "http://localhost:5001",
					GatewayURL:  "https://ipfs.io/ipfs/",
				},
			},
		},
		SessionSecrets: make(map[string]string),
		// LoadedEntries:  []models.VaultEntry{},
	}

	// vaultRepo := vaults_persistence.NewGormVaultRepository(db.DB)
	// sessionRepo := vaults_persistence.NewGormSessionRepository(db.DB)
	// sessionV2 := vault_session.NewManager(sessionRepo, vaultRepo, appLogger, ctx, ipfs)

	appLogger.Info("🔧 Initializing Tracecore client...")
	tcClient := tracecore.NewTracecoreClient(os.Getenv("ANKHORA_URL"), os.Getenv("TRACECORE_TOKEN"))
	appLogger.Info("✅ Tracecore client initialized")

	reg := registry.NewRegistry(appLogger)
	reg.RegisterDefinitions([]registry.EntryDefinition{
		{
			Type:    "login",
			Factory: func() models.VaultEntry { return &vaults_domain.LoginEntry{} },
			Handler: vault_ui.NewLoginHandler(*db, appLogger),
		},
		{
			Type:    "card",
			Factory: func() models.VaultEntry { return &vaults_domain.CardEntry{} },
			Handler: vault_ui.NewCardHandler(*db, appLogger),
		},
		{
			Type:    "note",
			Factory: func() models.VaultEntry { return &vaults_domain.NoteEntry{} },
			Handler: vault_ui.NewNoteHandler(*db, appLogger),
		},
		{
			Type:    "identity",
			Factory: func() models.VaultEntry { return &vaults_domain.IdentityEntry{} },
			Handler: vault_ui.NewIdentityHandler(*db, appLogger),
		},
		{
			Type:    "sshkey",
			Factory: func() models.VaultEntry { return &vaults_domain.SSHKeyEntry{} },
			Handler: vault_ui.NewSSHKeyHandler(*db, appLogger),
		},
	})
	// Legacy vault handler
	vaults := handlers.NewVaultHandler(*db, ipfs, reg, sessions, appLogger, tcClient, *runtimeCtx)
	onboardingUserRepo := onboarding_persistence.NewGormUserRepository(db.DB)
	auth := handlers.NewAuthHandler(*db, vaults, ipfs, appLogger, tcClient, cfg.auth, onboardingUserRepo)

	// Recovery context implementations
	userRepo := stellar_recovery_persistence.NewGormUserRepository(db.DB)
	vaultStellarRepo := stellar_recovery_persistence.NewGormVaultRepository(db.DB)
	stellarRecoverySubRepo := stellar_recovery_persistence.NewGormSubscriptionRepository(db.DB)
	verifier := stellar.NewStellarKeyAdapter()
	eventDisp := events.NewLocalDispatcher()
	tokenGen := token.NewSimpleTokenGen()

	checkUC := stellar_recovery_usecase.NewCheckKeyUseCase(userRepo, vaultStellarRepo, stellarRecoverySubRepo, verifier)
	recoverUC := stellar_recovery_usecase.NewRecoverVaultUseCase(userRepo, vaultStellarRepo, stellarRecoverySubRepo, verifier, eventDisp, tokenGen)
	importUC := stellar_recovery_usecase.NewImportKeyUseCase(userRepo, verifier)
	loginAdapter := shared.NewStellarLoginAdapter(db)

	connectUC := stellar_recovery_usecase.NewConnectWithStellarUseCase(
		loginAdapter,
		vaultStellarRepo,
		stellarRecoverySubRepo,
	)

	subscriptionSubRepo := subscription_persistence.NewSubscriptionRepository(db.DB, appLogger)
	userSubscriptionRepo := subscription_persistence.NewUserSubscriptionRepository(db.DB, appLogger)
	getRecommendedTierUC := onboarding_usecase.GetRecommendedTierUseCase{Db: db.DB}
	onboardingBus := onboarding_infrastructure_eventbus.NewMemoryBus()
	stellarService := blockchain.StellarService{}
	onboardingCreateAccountUC := onboarding_usecase.NewCreateAccountUseCase(&stellarService, onboardingUserRepo, onboardingBus, appLogger)
	stellarRecoveryHandler := stellar_recovery_ui_api.NewStellarRecoveryHandler(checkUC, recoverUC, importUC, connectUC)
	onboardingSetupPaymentUseCase := onboarding_usecase.NewSetupPaymentAndActivateUseCase(onboardingUserRepo, userSubscriptionRepo, subscriptionSubRepo, onboardingBus, *tcClient)
	onBoardingHandler := onboarding_ui_wails.NewOnBoardingHandler(&getRecommendedTierUC, onboardingCreateAccountUC, onboardingSetupPaymentUseCase, db.DB)
	subscriptionBus := subscription_infrastructure_eventbus.NewMemoryBus()
	createSubscriptionUC := subscription_usecase.NewCreateSubscriptionUseCase(subscriptionSubRepo, subscriptionBus, tcClient)
	subscriptionHandler := subscription_ui_wails.NewSubscriptionHandler(*createSubscriptionUC, subscriptionSubRepo)

	// runtimeCtxV2 := &vault_session.RuntimeContext{
	// 	AppConfig: app_config.AppConfig{
	// 		// Load from file/env or defaults
	// 		Branch:           "main",
	// 		EncryptionPolicy: "AES-256-GCM",
	// 		Blockchain: app_config.BlockchainConfig{
	// 			Stellar: app_config.StellarConfig{
	// 				Network:    "testnet",
	// 				HorizonURL: "https://horizon-testnet.stellar.org",
	// 				Fee:        100,
	// 			},
	// 			IPFS: app_config.IPFSConfig{
	// 				APIEndpoint: "http://localhost:5001",
	// 				GatewayURL:  "https://ipfs.io/ipfs/",
	// 			},
	// 		},
	// 	},
	// 	SessionSecrets: make(map[string]string),
	// }

	cryptoService := blockchain.CryptoService{}
	authRepository := auth_persistence.NewGormAuthRepository(db.DB)
	authTokenService := auth_usecases.NewTokenService(authV2, authRepository, db.DB)

	// AppConfig
	appConfigHandler := app_config_ui.NewAppConfigHandler(db.DB, *appLogger)

	// Vault
	vaultHandler := vault_ui.NewVaultHandler(reg, *appLogger, ctx, ipfs, &cryptoService, db.DB)

	// Identity
	identityMemoryBus := identity_infrastructure_eventbus.NewMemoryEventBus()
	identityHandler := identity_ui.NewIdentityHandler(db.DB, authTokenService, identityMemoryBus, onboardingUserRepo)

	// Auth
	tokenUC := auth_usecases.NewGenerateTokensUseCase(authRepository, authTokenService)
	authHandler := auth_ui.NewAuthHandler(identityHandler, tokenUC, db.DB)

	// Stripe webhook listener
	go func() {
		port := "4242" // your webhook port
		http.HandleFunc("/stripe-webhook", payments.WebhookHandler)

		log.Printf("🚀 Stripe webhook listener running on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("❌ Stripe webhook server failed: %v", err)
		}
	}()

	// ⚡ Restore sessions asynchronously to speed up startup
	/*
		go func() {
			appLogger.Info("🔄 Restoring sessions in background...")
			storedSessions, err := Db.db.GetAllSessions()
			if err != nil {
				appLogger.Error("❌ Failed to load stored sessions: %v", err)
				return
			}

			for _, s := range storedSessions {
				sessions[s.UserID] = s
				if len(s.PendingCommits) > 0 {
					for _, commit := range s.PendingCommits {
						if err := vaults.QueuePendingCommits(s.UserID, commit); err != nil {
							appLogger.Error("❌ Failed to queue commit for user %d: %v", s.UserID, err)
						}
					}
				}
			}
			appLogger.Info("✅ Restored %d sessions from DB", len(storedSessions))
		}()
	*/
	go func() {
		sessionDBModel := vaults_persistence.NewSessionDBModel(db.DB)
		appLogger.Info("🔄 Restoring sessions in background...")
		storedSessions, err := sessionDBModel.FindAll()
		if err != nil {
			appLogger.Error("❌ Failed to load stored sessions: %v", err)
			return
		}

		for _, s := range storedSessions {
			sessionsV2[s.UserID] = s
			if len(s.PendingCommits) > 0 {
				for _, commit := range s.PendingCommits {
					if err := vaults.QueuePendingCommits(s.UserID, commit); err != nil {
						appLogger.Error("❌ Failed to queue commit for user %d: %v", s.UserID, err)
					}
				}
			}
		}
		appLogger.Info("✅ Restored %d sessions from DB", len(storedSessions))
	}()

	// Event bus (single memory bus for subscription domain)
	subscriptionService := subscription_infrastructure.NewSubscriptionService()
	// ===== New: core activator (business logic) =====
	// Note: pass a Stellar port implementation if you have one, otherwise nil
	activator := subscription_usecase.NewSubscriptionActivator(
		subscriptionSubRepo, // repo
		subscriptionBus,
		subscriptionService, // vault port (implements ActivationVaultPort)
	)

	// ===== New: listener which only forwards SubscriptionCreated -> activator =====
	createdListener := subscription_usecase.NewSubscriptionCreatedListener(appLogger, activator, subscriptionBus)
	go createdListener.Listen(ctx)

	// ===== New: monitor for post-activation side effects (email, metrics...) =====
	billingBus := billing_infrastructure_eventbus.NewMemoryBus()
	billingHandler := billing_ui.NewBillingHandler(db.DB, billingBus)

	initializeVaultHandler := vault_commands.NewInitializeVaultCommandHandler(db.DB)
	createIpfsCommandHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(
		vaultHandler.VaultRepository, &cryptoService, ipfs,
	)
	createVaultCommand := vault_commands.NewCreateVaultCommandHandler(
		initializeVaultHandler, createIpfsCommandHandler, vaultHandler.VaultRepository,
	)

	activationMonitor := subscription_usecase.NewSubscriptionActivationMonitor(
		appLogger,
		subscriptionBus,
		userSubscriptionRepo,
		subscriptionSubRepo,
		createVaultCommand,
		&stellarService,
		onboardingUserRepo,
		onboardingBus,
		identityHandler,
		billingHandler,
	)
	go activationMonitor.Listen(ctx)

	// ===== New: vault monitor =====
	vaultOpenedListener := vault_commands.NewVaultOpenedListener(appLogger, vaultHandler.EventBus, vaultHandler)
	go vaultOpenedListener.Listen(ctx)
	appLogger.Info("Vault opened listener started")

	// Start pending commit worker
	vaults.StartPendingCommitWorker(ctx, 2*time.Minute)

	elapsed := time.Since(startTime)
	appLogger.Info("✅ D-Vault initialized successfully in %v", elapsed)

	return &App{
		AppConfigHandler:       appConfigHandler,
		Auth:                   auth,
		AuthHandler:            authHandler,
		cancel:                    cancel,
		ConnectWithStellarHandler: stellarRecoveryHandler,
		config:                 cfg,
		DB:                     *db,
		EntryRegistry:          reg,
		NowUTC:                    func() string { return time.Now().Format(time.RFC3339) },
		Identity:               identityHandler,
		Logger:                 *appLogger,
		OnBoardingHandler:      onBoardingHandler,
		sessions:               sessions,
		StellarRecoveryHandler: stellarRecoveryHandler,
		SubscriptionHandler:    subscriptionHandler,
		RuntimeContext:         runtimeCtx,
		// RuntimeContextV2:          runtimeCtxV2,
		Vault:                  vaultHandler,  // internal/vault/ui/vault_handler.go
		Vaults:                 vaults,        // internal/handlers/vault_handler.go
		version:                version,
	}

}

func (a *App) CheckPaymentOnResume() {
	// status, err := a.SubscriptionRepo.GetStatusForUser()
	// if err != nil {
	// 	return
	// }

	// if status == "active" {
	// 	runtime.EventsEmit(a.ctx, "payment:success")
	// }
}

func (a *App) NotifyPaymentSuccess(subID string) {
	a.Logger.Info("✅ Subscription created successfully: %v", subID)
	runtime.EventsEmit(a.ctx, "payment:success", subID)
}

// -----------------------------
// OnBoarding
// -----------------------------
type CheckKeyResponse struct {
	ID               string  `json:"id"`
	CreatedAt        string  `json:"created_at"`
	SubscriptionTier string  `json:"subscription_tier"`
	StorageUsedGB    float64 `json:"storage_used_gb"`
	LastSyncedAt     string  `json:"last_synced_at"`
	Ok               bool    `json:"ok"` // exported!
}

func (a *App) CheckStellarKeyForVault(stellarKey string) (*CheckKeyResponse, error) {
	res, err := a.StellarRecoveryHandler.CheckVault(context.Background(), stellarKey)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return &CheckKeyResponse{Ok: false}, nil
	}

	return &CheckKeyResponse{
		ID:               res.ID,
		CreatedAt:        res.CreatedAt,
		SubscriptionTier: res.SubscriptionTier,
		StorageUsedGB:    res.StorageUsedGB,
		LastSyncedAt:     res.LastSyncedAt,
		Ok:               true,
	}, nil
}
func (a *App) CreateAccount(req onboarding_usecase.AccountCreationRequest) (*onboarding_ui_wails.AccountCreationResponse, error) {
	utils.LogPretty("CreateAccount req", req)
	return a.OnBoardingHandler.CreateAccount(req)
}

func (a *App) RecoverVaultWithKey(stellarKey string) (*stellar_recovery_domain.RecoveredVault, error) {
	return a.StellarRecoveryHandler.RecoverVault(context.Background(), stellarKey)
}

func (a *App) ImportVaultWithKey(stellarKey string) (*stellar_recovery_domain.ImportedKey, error) {
	return a.StellarRecoveryHandler.ImportKey(context.Background(), stellarKey)
}

// in waiting for applying full ddd above
func (a *App) ConnectWithStellar(req handlers.LoginRequest) (*CheckKeyResponse, error) {
	response, err := a.StellarRecoveryHandler.ConnectWithStellar(context.Background(), req)
	fmt.Println("ConnectWithStellar req", response)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil // means no vault found
	}

	res := &CheckKeyResponse{
		ID:               response.ID,
		CreatedAt:        response.CreatedAt,
		SubscriptionTier: response.SubscriptionTier,
		StorageUsedGB:    response.StorageUsedGB,
		LastSyncedAt:     response.LastSyncedAt,
		Ok:               true,
	}
	utils.LogPretty("ConnectWithStellar res", res)
	return res, nil
}

type OnboardingStep1Response struct {
	Identity identity_domain.IdentityChoice `json:"identity"`
}

func (a *App) GetRecommendedTier(identity identity_domain.IdentityChoice) (OnboardingStep1Response, error) {
	choice := a.OnBoardingHandler.GetRecommendedTier(identity)
	return OnboardingStep1Response{Identity: identity_domain.IdentityChoice(choice)}, nil
}

// 0. Get Tier Features
func (a *App) GetTierFeatures() (map[string]onboarding_domain.SubscriptionFeatures, error) {
	return a.OnBoardingHandler.GetTierFeatures(), nil
}

// Step 2: Use Case (conditional based on Step 1)
type UseCaseResponse struct {
	UseCases []string `json:"use_cases"` // ["passwords", "financial", "medical", etc.]
}

func (a *App) SetupPaymentAndActivate(req onboarding_usecase.PaymentSetupRequest) (*subscription_domain.Subscription, error) {
	utils.LogPretty("SetupPaymentAndActivate req", req)
	return a.OnBoardingHandler.SetupPaymentAndActivate(req)
}

// Response with session ID
type CreateCheckoutResponse struct {
	SessionID string `json:"sessionId"`
	URL       string `json:"url"`
}

// GetCheckoutURL returns the cloud backend checkout page URL
func (a *App) GetCheckoutURL(plan string) (CreateCheckoutResponse, error) {
	// -----------------------------
	// 0. Generate Session ID
	// -----------------------------
	sessionID := uuid.New().String()

	// -----------------------------
	// 1. Generate Checkout URL
	// -----------------------------
	baseURL := "http://localhost:4002/checkout?session_id=" // your cloud page URL
	url := baseURL + sessionID

	res := CreateCheckoutResponse{
		SessionID: sessionID,
		URL:       url,
	}
	return res, nil
}

func (a *App) OpenURL(rawURL string) error {
	runtime.BrowserOpenURL(a.ctx, rawURL)
	return nil
}

// Poll backend for payment status
func (a *App) PollPaymentStatus(sessionID string, plainPassword string) (string, error) {
	// 0. ------------- Poll backend for payment status -----------------
	// fmt.Println("🔁 Polling session:", sessionID)
	url := "http://localhost:4001/api/payment-status/" + sessionID
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 1. ------------- Check response status -------------
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("poll failed %d: %s", resp.StatusCode, string(body))
	}

	var r struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	// 2. 	------------- Check payment status -------------
	if r.Status == "paid" {
		go func() {
			if err := a.OnPaymentConfirmation(sessionID, plainPassword); err != nil {
				a.Logger.Error("Payment confirmation failed:", err)
			}
		}()
		return "paid", nil
	}

	return "unpaid", nil

}

func (a *App) OnPaymentConfirmation(sessionID string, plainPassword string) error {
	log.Println("Deep link received:", sessionID)
	// 0. ------------- OnPaymentConfirmation -------------
	response, err := a.SubscriptionHandler.CreateSubscription(a.ctx, sessionID, plainPassword)
	if err != nil {
		return err
	}
	a.Logger.Info("✅ Subscription created successfully: %v", response)

	// 1. ------------- Notify frontend -------------
	runtime.EventsEmit(a.ctx, "payment:success", response.Subscription)
	return nil
}

// Simple method to open a URL in the system browser
func (a *App) OpenGoogle() {
	if a.ctx == nil {
		log.Println("❌ Context not set!")
		return
	}

	// Opens default browser to Google
	runtime.BrowserOpenURL(a.ctx, "http://localhost:4002/checkout")
}

// -----------------------------
// Connexion
// -----------------------------
func (a *App) Sign(req handlers.LoginRequest) (*handlers.LoginResponse, error) {
	return a.Auth.Login(req)
}
func (a *App) SignInWithStellar(req handlers.LoginRequest) (*handlers.LoginResponse, error) {
	return a.Auth.Login(req)
}

func (a *App) SignUp(setup handlers.OnBoarding) (*handlers.OnBoardingResponse, error) {
	return a.Auth.OnBoarding(setup)
}
func (a *App) SignOut(userID string) error {
	a.Logger.Info("App - SignOut userID", userID)
	// a.Auth.Logout(userID)
	if err := a.Vault.LogoutUser(userID); err != nil {
		a.Logger.Error("❌ SignOut failed for user %s: %v", userID, err)
		return err
	}
	a.Logger.Info("✅ User %s signed out", userID)

	return nil
}
func (a *App) CheckSession(userID string) (*auth.TokenPairs, error) {
	utils.LogPretty("CheckSession userID", userID)
	tokenPair, err := a.AuthHandler.GenerateTokenPair(userID)
	if err != nil {
		return nil, err
	}
	return tokenPair.ToFormerModel(), nil
	// return a.Auth.RefreshToken(userID) // same logic you already wrote
}
func (a *App) CheckEmail(email string) (*handlers.CheckEmailResponse, error) {
	return a.Auth.CheckEmail(email)
}
func (a *App) SaveSessionTest(jwtToken string) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	a.Logger.Info("App - SaveSessionTest userID", claims.UserID)
	// return a.Vault.SaveSession(claims.UserID)
	return nil
}

// -----------------------------
// Connexion (identity) - V2
// -----------------------------
func (a *App) SignInWithIdentity(req handlers.LoginRequest) (*handlers.LoginResponse, error) {
	cmd := identity_commands.LoginCommand{
		Email:         req.Email,
		Password:      req.Password,
		PublicKey:     req.PublicKey,
		SignedMessage: req.SignedMessage,
		Signature:     req.Signature,
	}
	// --------- Identity login ---------
	result, err := a.Identity.Login(cmd)
	if err != nil {
		a.Logger.Error("❌ App - SignInWithIdentity - failed to identify user %s: %v", result.User.ID, err)
		return nil, err
	}
	a.Logger.Info("Identity login successful: %v", result)

	// --------- Session Warm Up ---------
	session, err := a.Vault.PrepareSession(result.User.ID)
	if err != nil {
		a.Logger.Error("❌ App - SignInWithIdentity - failed to get session for user %s: %v", result.User.ID, err)
		// return	 nil, err
	}
	if session == nil {
		a.Logger.Error("❌ App - SignInWithIdentity - failed to get session for user %s: %v", result.User.ID, err)
		// return	 nil, err
	} else {
		a.Logger.Info("Session fetched successfully: %v", session)
	}

	// --------- Open vault ---------
	vaultRes, err := a.Vault.Open(
		context.Background(),
		vault_commands.OpenVaultCommand{
			UserID:   result.User.ID,
			Password: req.Password,
			Session:  session,
		},
		*a.AppConfigHandler,
	)
	if err != nil {
		return nil, err
	}
	a.Logger.Info(
		"Vault opened successfully for user %s (reused=%v)",
		result.User.ID,
		vaultRes.ReusedExisting,
	)

	// --------- Vault response converter for v1 ---------
	var formerVault *models.VaultPayload
	if vaultRes.Content != nil {
		formerVault = vaultRes.Content.ToFormerVaultPayload()
	}
	var formerRuntimeContext *models.VaultRuntimeContext
	if vaultRes.RuntimeContext != nil {
		formerRuntimeContext = vaultRes.RuntimeContext.ToFormerRuntimeContext()
	}

	utils.LogPretty("SignInWithIdentity - formerVault", formerVault.Name)

	// --------- Login response converter for v1 ---------
	loginRes := &handlers.LoginResponse{
		User:                *result.User.ToFormerUser(),
		Tokens:              result.Tokens.ToFormerModel(),
		SessionID:           session.UserID,
		Vault:               formerVault,
		VaultRuntimeContext: formerRuntimeContext,
		LastCID:             vaultRes.LastCID,
		Dirty:               session.Dirty,
	}
	utils.LogPretty("SignInWithIdentity - loginRes", loginRes.Vault)
	return loginRes, nil
}

type GetSessionResponse struct {
	Data  map[string]interface{}
	Error error
}

func (a *App) GetSession(userID string) (*GetSessionResponse, error) {
	if a.Vault.SessionManager == nil {
		return &GetSessionResponse{Error: errors.New("session manager not initialized")}, nil
	}

	userSession, err := a.Vault.GetSession(userID)
	if err != nil {
		return &GetSessionResponse{Error: err}, nil
	}
	user, err := a.Identity.FindUserById(a.ctx, userID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"User":                user,
		"role":                "user",
		"Vault":               userSession.Vault,
		"SharedEntries":       []models.VaultEntry{},
		"VaultRuntimeContext": userSession.Runtime,
		"LastCID":             userSession.LastCID,
		"Dirty":               userSession.Dirty,
	}
	return &GetSessionResponse{Data: response}, nil
}

// -----------------------------
// JWT Token
// -----------------------------
func (a *App) RefreshToken(userID string) (*auth.TokenPairs, error) {
	// token, err := a.Auth.RefreshToken(userID)
	utils.LogPretty("App - RefreshToken - userID", userID)
	token, err := a.AuthHandler.GenerateTokenPair(userID)
	if err != nil {
		a.Logger.Error("App - RefreshToken - error", err)
		return nil, err
	}
	utils.LogPretty("App - RefreshToken - token", token)

	return token.ToFormerModel(), nil
}
func (a *App) RequireAuth(jwtToken string) (*auth.Claims, error) {
	utils.LogPretty("App - RequireAuth ", jwtToken)
	claims, err := a.AuthHandler.VerifyToken(jwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	utils.LogPretty("claims", claims)
	return claims.ToFormerModel(), nil
}
func (a *App) RequestChallenge(req blockchain.ChallengeRequest) (blockchain.ChallengeResponse, error) {
	return a.Auth.RequestChallenge(req)
}
func (a *App) AuthVerify(req blockchain.SignatureVerification) (string, error) {
	return a.Auth.AuthVerify(&req)
}

// -----------------------------
// Vault Crud
// -----------------------------
func (a *App) GetVault(userID string) (map[string]interface{}, error) {
	user, err := a.Identity.FindUserById(a.ctx, userID)
	if err != nil {
		return nil, err
	}
	session, err := a.Vault.GetSession(userID)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"User":                user,
		"role":                "user",
		"Vault":               session.Vault,
		"SharedEntries":       []models.VaultEntry{},
		"VaultRuntimeContext": *session.Runtime,
		"LastCID":             session.LastCID,
		"Dirty":               session.Dirty,
	}
	return response, nil
}
func (a *App) AddEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AddEntry - error: %v", err)
		return nil, err
	}
	res, err := a.Vault.AddEntry(claims.UserID, entryType, raw)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("App - AddEntry - res", res)
	return res, nil
}

func (a *App) EditEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	res, err := a.Vault.UpdateEntry(claims.UserID, entryType, raw)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("App - EditEntry - res", res)
	return res, nil
}

func (a *App) TrashEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	res, err := a.Vault.TrashEntry(claims.UserID, entryType, raw)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("App - TrashEntry - res", res)
	return res, nil
}

func (a *App) RestoreEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	res, err := a.Vault.RestoreEntry(claims.UserID, entryType, raw)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("App - RestoreEntry - res", res)
	return res, nil
}

func (a *App) CreateFolder(name string, jwtToken string) (*vaults_domain.VaultPayload, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vault.CreateFolder(claims.UserID, name)
}
func (a *App) GetFoldersByVault(vaultCID string, jwtToken string) ([]vaults_domain.Folder, error) {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vault.GetFoldersByVault(vaultCID)
}
func (a *App) UpdateFolder(id string, newName string, isDraft bool, jwtToken string) (*vaults_domain.Folder, error) {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vault.UpdateFolder(id, newName, isDraft)
}
func (a *App) DeleteFolder(userID string, id string, jwtToken string) (string, error) {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Folder deleted %d successfuly", id), a.Vault.DeleteFolder(userID, id)
}

// TODO: ToggleFavorite entries

func (a *App) SynchronizeVault(jwtToken string, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}
	a.Vaults.Ctx = a.ctx
	return a.Vaults.SyncVault(claims.UserID, password)
}

func (a *App) EncryptFile(jwtToken string, fileData string, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}

	// Emit start progress
	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
		"percent": 0,
		"stage":   "encrypting",
	})

	// Real AES-256-GCM encryption with progress
	encryptedPath, err := a.Vaults.EncryptFile(claims.UserID, []byte(fileData), password)
	if err != nil {
		return "", err
	}

	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
		"percent": 70,
		"stage":   "encrypted",
	})

	return encryptedPath, nil
}
func (a *App) UploadToIPFS(jwtToken string, filePath string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}

	// Simulate upload progress (integrate with your IPFS client for real progress)
	current := 70
	for i := 1; i <= 20; i++ {
		current += 1
		runtime.EventsEmit(a.ctx, "progress-update", current)
		time.Sleep(50 * time.Millisecond) // Simulate; use actual IPFS progress
	}

	cid, err := a.Vaults.UploadToIPFS(claims.UserID, filePath)
	runtime.EventsEmit(a.ctx, "progress-update", 95) // Near complete
	if err != nil {
		return "", err
	}
	return cid, nil
}
func (a *App) CreateStellarCommit(jwtToken string, cid string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}

	// Quick commit with final progress
	runtime.EventsEmit(a.ctx, "progress-update", 100)
	return a.Vaults.CreateStellarCommit(claims.UserID, cid)
}
func (a *App) EncryptVault(jwtToken string, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}

	// Emit start progress
	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
		"percent": 0,
		"stage":   "encrypting",
	})

	// Real AES-256-GCM encryption with progress
	encryptedPath, err := a.Vaults.EncryptVault(claims.UserID, password)
	if err != nil {
		return "", err
	}

	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
		"percent": 70,
		"stage":   "encrypted",
	})

	return encryptedPath, nil
}

type CreateShareInput struct {
	Payload  handlers.CreateShareEntryPayload `json:"payload"`
	JwtToken string                           `json:"jwtToken"`
}

func (a *App) CreateShare(input CreateShareInput) (*share_domain.ShareEntry, error) {
	claims, err := a.Auth.RequireAuth(input.JwtToken)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID
	return a.Vaults.CreateShareEntry(context.Background(), input.Payload, userID, *a.AppConfigHandler, a.config.ANCHORA_SECRET)
}

func (a *App) ListSharedEntries(jwtToken string) (*[]share_domain.ShareEntry, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, fmt.Errorf("ListSharedEntries - auth failed: %w", err)
	}
	utils.LogPretty("ListSharedEntries - claims", claims)

	entries, err := a.Vaults.ListSharedEntries(context.Background(), claims.UserID)
	if err != nil {
		return nil, err
	}

	return &entries, nil
}
func (a *App) ListReceivedShares(jwtToken string) (*[]share_domain.ShareEntry, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	entries, err := a.Vaults.ListReceivedShares(context.Background(), claims.UserID)
	if err != nil {
		return nil, err
	}

	return &entries, nil // Wails wants pointer
}
func (a *App) GetShareForAccept(jwt, shareID string) (*share_domain.ShareAcceptData, error) {
	claims, err := a.Auth.RequireAuth(jwt)
	if err != nil {
		return nil, err
	}

	return a.Vaults.GetShareForAccept(
		context.Background(),
		claims.UserID,
		shareID,
	)
}
func (a *App) RejectShare(jwtToken string, shareID string) (*share_application.RejectShareResult, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	return a.Vaults.RejectShare(context.Background(), claims.UserID, shareID)
}
func (a *App) AddReceiver(jwtToken string, payload share_application.AddReceiverInput) (*share_application.AddReceiverResult, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("❌ Failed to authenticate user: %v", err)
		return nil, err
	}

	return a.Vaults.AddReceiver(context.Background(), claims.UserID, payload)
}
// -----------------------------
// Vault Config
// -----------------------------
type GenerateApiKeyInput struct {
	Password string `json:"password"`
	JwtToken string `json:"jwtToken"`
}	
type GenerateApiKeyOutput struct {
	PublicKey   string `json:"public_key"`
	PrivateKey  string `json:"private_key"`
}
func (a *App) GenerateApiKey(input GenerateApiKeyInput) (*GenerateApiKeyOutput, error) {
	claims, err := a.Auth.RequireAuth(input.JwtToken)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID
	utils.LogPretty("userID", userID)

	var stellarAccount *app_config.StellarAccountConfig
	stellarService := blockchain.StellarService{Logger: &a.Logger}
	res, err := stellarService.CreatAccountWithFriendbotFunding(input.Password)
	if err != nil {
		a.Logger.Warn("⚠️ Stellar account creation failed: %v", err)
		return nil, err	
	}
	utils.LogPretty("res", res)
	// Encrypt the user password with Stellar secret
	salt, nonce, ct, err := blockchain.EncryptPasswordWithStellarSecure(input.Password, res.PrivateKey)
	if err != nil {
		a.Logger.Warn("⚠️ Failed to encrypt password with Stellar secret: %v", err)
		return nil, err
	}

	// TODO: encrypt the Stellar private key before storing (server-side master key or KMS)
	utils.LogPretty("anchorSecret", a.config.ANCHORA_SECRET)
	if a.config.ANCHORA_SECRET == "" {
		a.Logger.Warn("⚠️ Anchora secret not found")
		return nil, errors.New("anchora secret not found")
	}
	encryptedPrivateKey, err := blockchain.Encrypt([]byte(res.PrivateKey), a.config.ANCHORA_SECRET)
	if err != nil {
		a.Logger.Warn("⚠️ Failed to encrypt Stellar private key: %v", err)
		return nil, err
	}

	stellarAccount = &app_config.StellarAccountConfig{
		PublicKey:   res.PublicKey,
		PrivateKey:  string(encryptedPrivateKey),
		EncSalt:     salt,
		EncNonce:    nonce,
		EncPassword: ct,
	}
	a.Logger.Info("✅ Stellar account created: %s - txID: %s", stellarAccount.PublicKey, res.TxID)
	
	// Handle stellar config 
	userCfg, err := a.Vaults.DB.GetUserConfigByUserID(userID)
	if userCfg == nil {
		userCfg = &app_config.UserConfig{}
	}
	userCfg.StellarAccount = *stellarAccount

	// Save user config
	savedUserCfg, err := a.Vaults.DB.SaveUserConfig(*userCfg)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("savedUserCfg", savedUserCfg)	

	return &GenerateApiKeyOutput{
		PublicKey:   res.PublicKey,
		PrivateKey:  string(encryptedPrivateKey),
		}, nil
}

// FlushAllSessions persists and clears all active sessions.
func (a *App) FlushAllSessions() {
	a.Vault.SessionsMu.Lock()
	defer a.Vault.SessionsMu.Unlock()
	// -----------------------------
	// 0. ENFORCE INVARIANTS (NON-NEGOTIABLE)
	// -----------------------------
	if !a.Vault.HasSession() {
		a.Logger.Info("No sessions to flush")
		return
	}

	a.Logger.Info("💾 Flushing %d active sessions...", len(a.Vault.GetAllSessions()))
	// -----------------------------
	// 1. FLUSH ALL SESSIONS
	// -----------------------------
	for userID := range a.Vault.GetAllSessions() {
		a.Vault.SessionManager.EndSession(userID)
	}
	a.Logger.Info("✨ All sessions flushed and cleared")
}

func (a *App) IsVaultDirty(jwtToken string) (bool, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return false, err
	}
	return a.Vault.IsMarkedDirty(claims.UserID), nil
}

func (a *App) FetchUsers() ([]models.UserDTO, error) {
	users, err := a.OnBoardingHandler.FetchUsers()
	if err != nil {
		return nil, err
	}
	var userDTOs []models.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, models.UserDTO{
			ID:              user.ID,
			Email:           user.Email,
			Role:            "user",
			LastConnectedAt: time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return userDTOs, nil
}

// Initialize Stripe secret key
func (a *App) InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET")
	if stripe.Key == "" {
		log.Fatal("❌ STRIPE_SECRET missing in .env")
	}
	log.Println("✅ Stripe initialized")
}

// -----------------------------
// Helpers
// -----------------------------
// Wails needs this to generate Entries struct in TypeScript
func (a *App) DummyExposeEntries(e models.Entries) models.Entries {
	return e
}
func mapToStruct(input map[string]interface{}, out any) error {
	return mapstructure.Decode(input, out)
}
func loadConfig() config {
	return config{
		db: struct{ dsn string }{
			dsn: os.Getenv("DB_DSN"), // or default
		},
		ANCHORA_SECRET: os.Getenv("ANCHORA_SECRET"),
	}
}
func (a *App) startup(ctx context.Context) {
	fmt.Println("App has started.")
	a.ctx = ctx
	a.InitStripe()
	// Fired every time the app regains focus
	runtime.EventsOn(ctx, "wails:window:focus", func(_ ...interface{}) {
		go a.CheckPaymentOnResume()
	})

}
