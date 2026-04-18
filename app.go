package main

import (
	"context"
	"net"
	"os/exec"

	// "encoding/base64"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	share_application_dto "vault-app/internal/application"
	share_application "vault-app/internal/application/use_cases"
	"vault-app/internal/auth"
	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_persistence "vault-app/internal/auth/infrastructure/persistence"
	auth_ui "vault-app/internal/auth/ui"
	billing_domain "vault-app/internal/billing/domain"
	billing_ui "vault-app/internal/billing/ui"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	app_config_dto "vault-app/internal/config/application/dto"
	app_config_domain "vault-app/internal/config/domain"

	// "vault-app/internal/config/infrastructure/persistence"
	app_config_ui "vault-app/internal/config/ui"
	share_domain "vault-app/internal/domain/shared"
	"vault-app/internal/driver"
	"vault-app/internal/handlers"
	identity_commands "vault-app/internal/identity/application/commands"
	identity_dtos "vault-app/internal/identity/application/dtos"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	onboarding_persistence "vault-app/internal/onboarding/infrastructure/persistence"
	onboarding_ui_wails "vault-app/internal/onboarding/ui/wails"
	"vault-app/internal/registry"
	shared "vault-app/internal/shared/stellar"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
	"vault-app/internal/stellar_recovery/infrastructure/events"
	"vault-app/internal/stellar_recovery/infrastructure/token"
	stellar_recovery_ui_api "vault-app/internal/stellar_recovery/ui/api"
	payments "vault-app/internal/stripe"
	subscription_domain "vault-app/internal/subscription/domain"
	subscription_persistence "vault-app/internal/subscription/infrastructure/persistence"
	subscription_ui_wails "vault-app/internal/subscription/ui/wails"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	utils "vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"
	vault_ui "vault-app/internal/vault/ui"

	// "vault-app/internal/logger/logger"
	"vault-app/internal/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
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
	IsOnborded       bool
	Domain           string
	Branch           string
	EncryptionPolicy string

	// Jwt auth
	auth           auth.Auth
	JWTSecret      string
	JWTIssuer      string
	JWTAudience    string
	APIKey         string
	ANCHORA_SECRET string

	// Stripe
	stripe struct {
		secret string
		key    string
	}

	// Stellar
	StellarNetwork     string
	StellarHorizonURL  string
	StellarAssetCode   string
	StellarAssetIssuer string

	// IPFS
	IPFSClient  string
	IPFSGateway string
	IPFSNetwork string

	// Tracecore
	TracecoreURL   string
	TracecoreToken string

	// Cloud
	CloudURL      string
	CloudBackURL  string
	CloudFrontURL string

	KEYRING_PATH string
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
	BillingHandler            *billing_ui.BillingHandler
	ConnectWithStellarHandler *stellar_recovery_ui_api.StellarRecoveryHandler
	EntryRegistry             *registry.EntryRegistry
	Identity                  *identity_ui.IdentityHandler
	OnBoardingHandler         *onboarding_ui_wails.OnBoardingHandler
	StellarService            *blockchain.StellarService
	StellarRecoveryHandler    *stellar_recovery_ui_api.StellarRecoveryHandler
	SubscriptionHandler       *subscription_ui_wails.SubscriptionHandler
	Vault                     *vault_ui.VaultHandler
	Vaults                    *handlers.VaultHandler

	// New: Global state
	RuntimeContext *vault_session.RuntimeContext
	cancel         context.CancelFunc
}

// NewApp creates a new App instance (required by Wails)
func NewApp() *App {
	startTime := time.Now()
	utils.LogPretty("Local IP Address", GetLocalIP())

	// -----------------------------------
	// Initialize
	// -----------------------------------
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := loadConfig()

	// Use auto-init from env
	appLogger := logger.NewFromEnv()
	appLogger.Info("🚀 Starting D-Vault initialization...")
	appLogger.Info("App Version: %s", version)
	appLogger.LogPretty("* App Config ***", cfg)

	// Pick DSN
	dsn := cfg.db.dsn
	if dsn == "" {
		dsn = "sqlite3.db"
	}

	// -------------------------------------------------------------------------------------------------
	// Database
	// -------------------------------------------------------------------------------------------------
	db, err := driver.InitDatabase(dsn, *appLogger)
	if err != nil {
		appLogger.Error("❌ Failed to init DB: %v", err)
		os.Exit(1)
	}
	appLogger.Info("✅ Local DB ready")

	// legacy
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

	// -------------------------------------------------------------------------------------------------
	// Blockchain
	// -------------------------------------------------------------------------------------------------
	ipfs := blockchain.NewIPFSClient(cfg.IPFSClient)
	appLogger.Info("✅ IPFS client initialized (connection will be tested on first use)")

	// -------------------------------------------------------------------------------------------------
	// Sessions
	// -------------------------------------------------------------------------------------------------
	sessions := make(map[string]*models.VaultSession) // legacy
	sessionsV2 := make(map[string]*vault_session.Session)

	// -------------------------------------------------------------------------------------------------
	// Context - Background
	// -------------------------------------------------------------------------------------------------
	ctx, cancel := context.WithCancel(context.Background())

	// -------------------------------------------------------------------------------------------------
	// Legacy - Runtime Context
	// -------------------------------------------------------------------------------------------------
	runtimeCtxLegacy := &vault_session.RuntimeContext{
		AppConfig: app_config_domain.AppConfig{
			// Load from file/env or defaults
			Branch:           cfg.Branch,
			EncryptionPolicy: cfg.EncryptionPolicy,
			Blockchain: app_config_domain.BlockchainConfig{
				Stellar: app_config_domain.StellarConfig{
					Network:    cfg.StellarNetwork,
					HorizonURL: cfg.StellarHorizonURL,
					Fee:        100,
				},
				IPFS: app_config_domain.IPFSConfig{
					APIEndpoint: cfg.IPFSClient,
					GatewayURL:  cfg.IPFSGateway,
				},
			},
		},
		SessionSecrets: make(map[string]string),
		// LoadedEntries:  []models.VaultEntry{},
	}

	// -------------------------------------------------------------------------------------------------
	// Tracecore
	// -------------------------------------------------------------------------------------------------
	tracecoreClient := tracecore.NewTracecoreClient(cfg.TracecoreURL, cfg.TracecoreToken, cfg.CloudFrontURL, cfg.CloudBackURL)

	// -------------------------------------------------------------------------------------------------
	// Registry
	// -------------------------------------------------------------------------------------------------
	reg := registry.NewRegistry(appLogger)
	reg.RegisterDefinitions([]registry.EntryDefinition{
		{
			Type:    "login",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.LoginEntry{} },
			Handler: vault_ui.NewLoginHandler(*db, appLogger),
		},
		{
			Type:    "card",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.CardEntry{} },
			Handler: vault_ui.NewCardHandler(*db, appLogger),
		},
		{
			Type:    "note",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.NoteEntry{} },
			Handler: vault_ui.NewNoteHandler(*db, appLogger),
		},
		{
			Type:    "identity",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.IdentityEntry{} },
			Handler: vault_ui.NewIdentityHandler(*db, appLogger),
		},
		{
			Type:    "sshkey",
			Factory: func() vaults_domain.VaultEntry { return &vaults_domain.SSHKeyEntry{} },
			Handler: vault_ui.NewSSHKeyHandler(*db, appLogger),
		},
	})
	appLogger.Info("✅ Registry initialized")

	// -------------------------------------------------------------------------------------------------
	// Legacy - Vault Handler
	// -------------------------------------------------------------------------------------------------
	vaults := handlers.NewVaultHandler(*db, ipfs, reg, sessions, appLogger, tracecoreClient, *runtimeCtxLegacy)
	onboardingUserRepo := onboarding_persistence.NewGormUserRepository(db.DB)
	auth := handlers.NewAuthHandler(*db, vaults, ipfs, appLogger, tracecoreClient, cfg.auth, onboardingUserRepo)

	stellarService := blockchain.NewStellarService(appLogger)

	// -------------------------------------------------------------------------------------------------
	// AppConfig
	// -------------------------------------------------------------------------------------------------
	appConfigHandler := app_config_ui.NewAppConfigHandler(db.DB, *appLogger)
	appLogger.Info("AppConfigHandler - NewAppConfigHandler - appConfigHandler", appConfigHandler)

	// -------------------------------------------------------------------------------------------------
	// Crypto Service
	// -------------------------------------------------------------------------------------------------
	cryptoService := blockchain.CryptoService{}

	// -------------------------------------------------------------------------------------------------
	// Vault
	// -------------------------------------------------------------------------------------------------
	vaultHandler := vault_ui.NewVaultHandler(
		reg,
		*appLogger,
		ctx,
		ipfs,
		&cryptoService,
		db.DB,
		*tracecoreClient,
		cfg.KEYRING_PATH,
	)

	// -------------------------------------------------------------------------------------------------
	// Subscription
	// -------------------------------------------------------------------------------------------------
	userSubscriptionRepo := subscription_persistence.NewUserSubscriptionRepository(db.DB, appLogger)
	subscriptionSubRepo := subscription_persistence.NewSubscriptionRepository(db.DB, appLogger)

	// -------------------------------------------------------------------------------------------------
	// Onboarding
	// -------------------------------------------------------------------------------------------------
	onBoardingHandler := onboarding_ui_wails.NewOnBoardingHandler(
		stellarService,
		userSubscriptionRepo,
		subscriptionSubRepo,
		tracecoreClient,
		db.DB,
		appLogger,
		*vaultHandler.KeyringService,
	)

	// -------------------------------------------------------------------------------------------------
	// Auth Infrastructure
	// -------------------------------------------------------------------------------------------------
	authRepository := auth_persistence.NewGormAuthRepository(db.DB)
	authTokenService := auth_usecases.NewTokenService(authV2, authRepository, db.DB)

	// -------------------------------------------------------------------------------------------------
	// Identity
	// -------------------------------------------------------------------------------------------------
	identityHandler := identity_ui.NewIdentityHandler(db.DB, authTokenService, onBoardingHandler.UserRepo)

	// -------------------------------------------------------------------------------------------------
	// Auth
	// -------------------------------------------------------------------------------------------------
	tokenUC := auth_usecases.NewGenerateTokensUseCase(authRepository, authTokenService)
	authHandler := auth_ui.NewAuthHandler(identityHandler, tokenUC, db.DB)

	// -------------------------------------------------------------------------------------------------
	// Subscription
	// -------------------------------------------------------------------------------------------------
	subscriptionHandler := subscription_ui_wails.NewSubscriptionHandler(
		db.DB,
		tracecoreClient,
		vaultHandler.CreateVaultCommandHandler,
		stellarService,
		onBoardingHandler.UserRepo,
		onBoardingHandler.Bus,
		identityHandler,
		*appConfigHandler,
		*appLogger,
	)

	// -------------------------------------------------------------------------------------------------
	// Billing
	// -------------------------------------------------------------------------------------------------
	billingHandler := billing_ui.NewBillingHandler(db.DB, &subscriptionHandler.SubscriptionSyncService, tracecoreClient)
	subscriptionHandler.SetBillingHandler(*billingHandler)
	appLogger.Info("billingHandler", billingHandler)
	appLogger.Info("subscriptionHandler", subscriptionHandler)

	// -------------------------------------------------------------------------------------------------
	// Stellar Recovery
	// -------------------------------------------------------------------------------------------------
	eventDisp := events.NewLocalDispatcher()
	tokenGen := token.NewSimpleTokenGen()
	loginAdapter := shared.NewStellarLoginAdapter(db)
	stellarRecoveryHandler := stellar_recovery_ui_api.NewStellarRecoveryHandler(db.DB, eventDisp, tokenGen, loginAdapter)

	// -------------------------------------------------------------------------------------------------
	// Stripe webhook listener
	// -------------------------------------------------------------------------------------------------
	go func() {
		port := "4242" // your webhook port
		http.HandleFunc("/stripe-webhook", payments.WebhookHandler)

		log.Printf("🚀 Stripe webhook listener running on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("❌ Stripe webhook server failed: %v", err)
		}
	}()

	// -------------------------------------------------------------------------------------------------
	// ⚡ Restore sessions asynchronously to speed up startup
	// -------------------------------------------------------------------------------------------------
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
	// ===== New: core activator (business logic) =====

	// ===== New: listener which only forwards SubscriptionCreated -> activator =====
	go subscriptionHandler.CreateListener.Listen(ctx)

	// ===== New: monitor for post-activation side effects (email, metrics...) =====
	go subscriptionHandler.MonitorActivationService.Listen(ctx)

	// ===== New: vault monitor =====
	vaultHandler.InitializeVaultOpenedListener()
	go vaultHandler.VaultOpenedListener.Listen(ctx)
	appLogger.Info("✅ Vault opened listener started")

	// Start pending commit worker
	vaults.StartPendingCommitWorker(ctx, 2*time.Minute)

	elapsed := time.Since(startTime)
	appLogger.Info("✅ D-Vault initialized successfully in %v", elapsed)

	// Startup:
	// ResetAndMigrate(db.DB) // Run ONCE on prod startup

	return &App{
		AppConfigHandler:          appConfigHandler,
		Auth:                      auth,
		BillingHandler:            billingHandler,
		AuthHandler:               authHandler,
		cancel:                    cancel,
		ConnectWithStellarHandler: stellarRecoveryHandler,
		config:                    cfg,
		DB:                        *db,
		EntryRegistry:             reg,
		NowUTC:                    func() string { return time.Now().Format(time.RFC3339) },
		Identity:                  identityHandler,
		Logger:                    *appLogger,
		OnBoardingHandler:         onBoardingHandler,
		sessions:                  sessions, // TODO: remove legacy sessions
		StellarService:            stellarService,
		StellarRecoveryHandler:    stellarRecoveryHandler,
		SubscriptionHandler:       subscriptionHandler,
		RuntimeContext:            runtimeCtxLegacy,
		Vault:                     vaultHandler, // internal/vault/ui/vault_handler.go
		Vaults:                    vaults,       // internal/handlers/vault_handler.go legacy
		version:                   version,
	}
}

// func (a *App) TestUploadIPFS(jwtToken string) error {

// 	vaultName := "tamboo"
// 	entryType := "login"

// 	// =========================
// 	// 1. LOAD FILE
// 	// =========================
// 	filePath := "./img-001.png"

// 	fileBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		log.Println("❌ failed to read file:", err)
// 	}

// 	log.Println("📦 Original file size:", len(fileBytes))

// 	// =========================
// 	// 2. LOCAL UPLOAD (OPTIONAL)
// 	// =========================
// 	raw := json.RawMessage{}

// 	attachments := vault_dto.SelectedAttachments{}

// 	_, errUpload := a.UploadAttachments(
// 		jwtToken,
// 		vaultName,
// 		entryType,
// 		raw,
// 		attachments,
// 	)

// 	if errUpload != nil {
// 		log.Println("❌ UploadAttachments error:", errUpload)
// 	}

// 	// =========================
// 	// 3. ENCRYPT
// 	// =========================
// 	encrypted, err := a.EncryptAttachment(
// 		jwtToken,
// 		fileBytes,
// 		"vaultPassword",
// 	)

// 	if err != nil {
// 		log.Println("❌ Encrypt error:", err)
// 	}

// 	log.Println("🔐 Encrypted size:", len(encrypted))

// 	// =========================
// 	// 4. IPFS UPLOAD
// 	// =========================
// 	cid, err := a.UploadAttachmentToIPFSWithEncryption(
// 		jwtToken,
// 		encrypted,
// 	)

// 	if err != nil {
// 		log.Println("❌ IPFS upload error:", err)
// 	}

// 	log.Println("🌐 CID:", cid)

// 	// =========================
// 	// 5. FETCH FROM IPFS
// 	// =========================
// 	ipfsFile, err := a.GetIPFSFile(jwtToken, cid)
// 	if err != nil {
// 		log.Println("❌ IPFS fetch error:", err)
// 	}

// 	log.Println("📥 IPFS size:", len(ipfsFile))

// 	// =========================
// 	// 6. DECRYPT
// 	// =========================
// 	decrypted, err := a.DecryptAttachment(
// 		jwtToken,
// 		ipfsFile,
// 		"vaultPassword",
// 	)

// 	if err != nil {
// 		log.Println("❌ Decrypt error:", err)
// 	}

// 	log.Println("🔓 Decrypted size:", len(decrypted))

// 	// =========================
// 	// 7. VERIFY
// 	// =========================
// 	if bytes.Equal(fileBytes, decrypted) {
// 		log.Println("✅ SUCCESS: encryption roundtrip OK")
// 	} else {
// 		log.Println("❌ FAILURE: decrypted != original")
// 		log.Println("Original:", len(fileBytes))
// 		log.Println("Decrypted:", len(decrypted))
// 	}
// 	return nil
// }

// -----------------------------
// AppState
// -----------------------------
func (a *App) GetAppState() (*onboarding_domain.AppState, error) {
	appState, err := a.OnBoardingHandler.GetAppState()
	if err != nil {
		return onboarding_domain.NewAppState(), nil
	}
	return appState, nil
}
func (a *App) CompleteOnboarding() error {
	return a.OnBoardingHandler.CompleteOnboarding()
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

func (a *App) SetupFreeAndActivate(req onboarding_usecase.FreeSetupRequest) (*tracecore.FreeCheckoutResponse, error) {
	utils.LogPretty("SetupFreeAndActivate req", req)

	response, err := a.OnBoardingHandler.SetupFreeAndActivate(req)
	if err != nil {
		utils.LogPretty("SetupFreeAndActivate err", err)
		return nil, err
	}
	utils.LogPretty("SetupFreeAndActivate response", response)

	return response, nil
}

// Response with session ID
type CreateCheckoutResponse struct {
	SessionID string `json:"sessionId"`
	URL       string `json:"url"`
}

// GetCheckoutURL returns the cloud backend checkout page URL
func (a *App) GetCheckoutURL(identity identity_domain.IdentityChoice, isAnonymous bool, rail string, email, tier, plan string) (CreateCheckoutResponse, error) {
	// -----------------------------
	// 0. Generate Session ID
	// -----------------------------
	sessionID := uuid.New().String()
	periodMonths := "1"

	// -----------------------------
	// 1. Generate Checkout URL
	// -----------------------------
	baseURL := a.config.CloudFrontURL + "/checkout" // your cloud page URL
	url := fmt.Sprintf("%s?session_id=%s&identity=%s&rail=%s&email=%s&tier=%s&plan=%s&period_months=%s&isAnonymous=%t", baseURL, sessionID, identity, rail, email, tier, plan, periodMonths, isAnonymous)

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
func (a *App) OpenFileInDefaultApp(path string) error {
    // On macOS
    cmd := exec.Command("open", path)
    return cmd.Run()
}

// Poll backend for payment status
func (a *App) PollPaymentStatus(sessionID string, email string, plainPassword string) (string, error) {
	// 0. ------------- Poll backend for payment status -----------------
	// fmt.Println("🔁 Polling session:", sessionID)
	url := a.config.CloudBackURL + "/billing/payment-status/" + sessionID
	a.Logger.Info("Polling payment status:", url)
	resp, err := http.Get(url)
	if err != nil {
		a.Logger.Error("Polling payment status failed:", err)
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
	if r.Status == "active" || r.Status == "paid" {
		go func() {
			if err := a.OnPaymentConfirmation(sessionID, email, plainPassword); err != nil {
				a.Logger.Error("Payment confirmation failed:", err)
			}
		}()
		return "paid", nil
	}

	return "unpaid", nil
}

func (a *App) OnPaymentConfirmation(sessionID string, email string, plainPassword string) error {
	a.Logger.Info("Deep link received:", sessionID)
	// 0. ------------- OnPaymentConfirmation -------------
	response, err := a.SubscriptionHandler.CreateSubscription(a.ctx, sessionID, email, plainPassword)
	if err != nil {
		a.Logger.Error("OnPaymentConfirmation - Payment confirmation failed:", err)
		return err
	}
	a.Logger.Info("✅ Subscription created successfully: %v", response)

	// 1. ------------- Notify frontend -------------
	runtime.EventsEmit(a.ctx, "payment:success", response.Subscription)
	return nil
}

// -----------------------------
// Connexion Legagcy
// -----------------------------
func (a *App) Sign(req handlers.LoginRequest) (*handlers.LoginResponse, error) {
	return a.Auth.Login(req)
}
func (a *App) SignUp(setup handlers.OnBoarding) (*handlers.OnBoardingResponse, error) {
	return a.Auth.OnBoarding(setup)
}
func (a *App) SignOut(userID string) error {
	a.Logger.Info("App - SignOut userID", userID)
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
func (a *App) CheckUserEmail(jwtToken string, email string) (*tracecore_types.User, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("❌ App - CheckUserEmail - failed to require auth %s: %v", claims.UserID, err)
		return nil, err
	}
	return a.Auth.CheckUserEmail(email, "token")
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
// Connexion (identity)
// -----------------------------
func (a *App) SignInWithStellar(req handlers.LoginRequest) (*vault_dto.LoginResponse, error) {
	a.Logger.Info("App - SignInWithStellar req", req)
	return a.SignIn(req)
}
func (a *App) SignInWithIdentity(req handlers.LoginRequest) (*vault_dto.LoginResponse, error) {
	a.Logger.Info("App - SignInWithIdentity req", req)
	return a.SignIn(req)
}
func (a *App) SignIn(req handlers.LoginRequest) (*vault_dto.LoginResponse, error) {
	cmd := identity_commands.LoginCommand{
		Email:         req.Email,
		Password:      req.Password,
		PublicKey:     req.PublicKey,
		SignedMessage: req.SignedMessage,
		Signature:     req.Signature,
	}
	// --------- Identity login ---------
	if a.Identity == nil {
		a.Logger.Error("❌ App - SignIn - identity is not initialized")
		return nil, errors.New("App - SignIn - identity is not initialized")
	}
	a.Logger.Info("App - SignIn - identity is initialized")

	result, err := a.Identity.Login(cmd)
	if err != nil {
		a.Logger.Error("❌ App - SignIn - failed to identify user %s: %v", result.User.ID, err)
		return nil, err
	}
	a.Logger.Info("Identity login successful: %v", result)

	// --------- Session Warm Up ---------
	session, err := a.Vault.PrepareSession(result.User.ID)
	if err != nil {
		a.Logger.Error("❌ App - SignIn - failed to get session for user %s: %v", result.User.ID, err)
	}
	if session == nil {
		a.Logger.Error("❌ App - SignIn - failed to get session for user %s: %v", result.User.ID, err)
		// return	 nil, err
	} else {
		a.Logger.Info("Session fetched successfully: %v", session)
	}

	a.Logger.LogPretty("Session provisionned successfully - Runtime from session: ", session.Runtime)

	// --------- Find user onboarding ---------
	userOnboarding, err := a.OnBoardingHandler.FindUsersUseCase.FindByEmail(result.User.Email)
	if err != nil {
		a.Logger.Error("❌ App - SignIn - failed to find user onboarding for user %s: %v", result.User.ID, err)
		return nil, err
	}
	a.Logger.Info("User onboarding found successfully: %v", userOnboarding)

	// --------- Open vault ---------
	vaultRes, err := a.Vault.Open(
		context.Background(),
		vault_commands.OpenVaultCommand{
			UserID:   result.User.ID,
			Password: req.Password,
			Session:  session,
			UserOnboardingID: userOnboarding.ID,
		},
		a.AppConfigHandler,
	)
	if err != nil {
		a.Logger.Error("❌ App - SignIn - failed to open vault for user %s: %v", result.User.ID, err)
		return nil, err
	}
	a.Logger.Info(
		"Vault opened successfully for user %s (reused=%v)",
		result.User.ID,
		vaultRes.ReusedExisting,
	)

	loginRes := &vault_dto.LoginResponse{
		User:                *result.User,
		Tokens:              result.Tokens,
		SessionID:           session.UserID,
		Vault:               *vaultRes.Content,
		VaultRuntimeContext: *vaultRes.RuntimeContext,
		LastCID:             vaultRes.LastCID,
		Dirty:               session.Dirty,
	}

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
		"SharedEntries":       []vaults_domain.VaultEntry{},
		"VaultRuntimeContext": userSession.Runtime,
		"LastCID":             userSession.LastCID,
		"Dirty":               userSession.Dirty,
	}
	return &GetSessionResponse{Data: response}, nil
}
func (a *App) GetConfig(vaultName string, jwtToken string) (*app_config_domain.Config, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	a.Logger.Info("App - GetConfig - vaultName", vaultName)
	return a.AppConfigHandler.GetConfig(claims.UserID, vaultName)
}

func (a *App) EditConfig(vaultName string, s *app_config_dto.Settings, jwtToken string) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	a.Logger.LogPretty("App - EditConfig - settings", s)

	return a.AppConfigHandler.EditSettings(claims.UserID, vaultName, s)
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

func (a *App) AccessDecryptVaultEntry(jwtToken string, entry tracecore_types.AccessCryptoShareRequest) (*tracecore_types.CloudResponse[tracecore_types.DecryptCryptoShareResponse], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	fmt.Println("userID", claims.UserID)
	// 1. Access encrypted entry ==============================
	entry.IPAddress = GetLocalIP()
	res, err := a.Vault.AccessEncryptedEntry(a.ctx, claims.UserID, entry, *a.Auth.TracecoreClient)
	if err != nil {
		return nil, err
	}
	a.Logger.LogPretty("App - AccessEncryptedEntry - res", res)

	// UserConfig - Get stellar private key from user config
	userConfig, err := a.AppConfigHandler.GetUserConfigByUserID(claims.UserID)
	if err != nil {
		return nil, err
	}
	// 2. Decrypt ==============================
	stellarAccount := userConfig.StellarAccount
	req := tracecore_types.DecryptCryptoShareRequest{
		EncryptedKey:        res.Data.EncryptedKey,
		EncryptedPayload:    res.Data.EncryptedPayload,
		RecipientPrivateKey: stellarAccount.PrivateKey,
	}
	// TODO: check if thisshoild be made by the user and not the cloud
	response, err := a.Vault.DecryptVaultEntry(context.Background(), req, *a.Auth.TracecoreClient)
	if err != nil {
		return nil, err
	}
	a.Logger.LogPretty("App - DecryptVaultEntry - response", response)
	// 3. Apply access policy from AppConfig ==============================
	appConfig, err := a.AppConfigHandler.GetAppConfigByUserID(context.Background(), claims.UserID)
	if err != nil {
		return nil, err
	}
	response.Data.ExpiresIn = appConfig.AccessPolicyDuration
	utils.LogPretty("App - DecryptVaultEntry - Final response", response.Data.ExpiresIn)

	return response, nil
}

// -----------------------------
// Vault Crud
// -----------------------------
func (a *App) AddEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AddEntry - error: %v", err)
		return nil, err
	}
	res, err := a.Vault.AddEntry(claims.UserID, entryType, raw)
	if err != nil {
		a.Logger.Error("App - AddEntry - error: %v", err)
		return nil, err
	}
	utils.LogPretty("App - AddEntry - res", res)
	return res, nil
}
func (a *App) EditEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - EditEntry - error: %v", err)
		return nil, err
	}
	isSyncMode := false
	res, err := a.Vault.UpdateEntry(claims.UserID, entryType, raw, isSyncMode)
	if err != nil {
		a.Logger.Error("App - EditEntry - error: %v", err)
		return nil, err
	}
	// utils.LogPretty("App - EditEntry - res", res)
	return res, nil
}
func (a *App) TrashEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - TrashEntry - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - TrashEntry - payload", map[string]interface{}{"raw": raw, "entryType": entryType})

	res, err := a.Vault.TrashEntry(claims.UserID, entryType, raw)
	if err != nil {
		a.Logger.Error("App - TrashEntry - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - TrashEntry - res", res)
	return res, nil
}
func (a *App) RestoreEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - RestoreEntry - error: %v", err)
		return nil, err
	}

	res, err := a.Vault.RestoreEntry(claims.UserID, entryType, raw)
	if err != nil {
		a.Logger.Error("App - RestoreEntry - error: %v", err)
		return nil, err
	}

	a.Logger.LogPretty("App - RestoreEntry - res", res)
	return res, nil
}
func (a *App) DeleteEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - DeleteEntry - error: %v", err)
		return nil, err
	}
	// TODO: Implement permanent deletion via API (different from trash)
	res, err := a.Vault.TrashEntry(claims.UserID, entryType, raw)
	if err != nil {
		a.Logger.Error("App - DeleteEntry - error: %v", err)
		return nil, err
	}
	utils.LogPretty("App - DeleteEntry - res", res)
	return res, nil
}
func (a *App) CreateFolder(name string, jwtToken string) (*vaults_domain.VaultPayload, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - CreateFolder - error: %v", err)
		return nil, err
	}
	return a.Vault.CreateFolder(claims.UserID, name)
}
func (a *App) GetFoldersByVault(vaultCID string, jwtToken string) ([]vaults_domain.Folder, error) {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - GetFoldersByVault - error: %v", err)
		return nil, err
	}
	return a.Vault.GetFoldersByVault(vaultCID)
}
func (a *App) UpdateFolder(id string, newName string, isDraft bool, jwtToken string) (*vaults_domain.Folder, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UpdateFolder - error: %v", err)
		return nil, err
	}
	return a.Vault.UpdateFolder(claims.UserID, newName, isDraft)
}
func (a *App) DeleteFolder(id string, jwtToken string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - DeleteFolder - error: %v", err)
		return "", err
	}
	a.Vault.DeleteFolder(claims.UserID, id)
	return fmt.Sprintf("Folder deleted %s successfuly", id), nil
}

// -----------------------------
// Cloud Services
// -----------------------------
func (a *App) SynchronizeVault(jwtToken string, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - SynchronizeVault - error: %v", err)
		return "", err
	}

	// Get Vault ==============================
	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		a.Logger.Error("App - SynchronizeVault - error: %v", err)
		return "", err
	}
	a.Logger.LogPretty("App - SynchronizeVault - vault", vault)

	// Sync Vault ==============================
	input := vault_dto.SynchronizeVaultRequest{
		UserID:   claims.UserID,
		Password: password,
		Vault:    *vault,
	}
	a.Vaults.Ctx = a.ctx

	return a.Vault.SyncVault(a.ctx, input, *a.Auth.TracecoreClient)
}
func (a *App) EncryptFile(jwtToken string, fileData string, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - EncryptFile - error: %v", err)
		return "", err
	}

	// Emit start progress
	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
		"percent": 0,
		"stage":   "encrypting",
	})
	a.Logger.LogPretty("App - EncryptFile - fileData", fileData)

	// Real AES-256-GCM encryption with progress
	encryptedPath, err := a.Vaults.EncryptFile(claims.UserID, []byte(fileData), password)
	if err != nil {
		a.Logger.Error("App - EncryptFile - error: %v", err)
		return "", err
	}

	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
		"percent": 70,
		"stage":   "encrypted",
	})

	return encryptedPath, nil
}
func (a *App) EncryptAttachment(jwtToken string, data []byte, password string) ([]byte, error) {
	_, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	return a.Vault.EncryptAttachment(data, password)
}
func (a *App) DecryptAttachment(jwtToken string, data []byte, password string) ([]byte, error) {
	_, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	a.Logger.Info("DecryptAttachment processing....")

	return a.Vault.DecryptAttachment(data, password)
}

// func (a *App) DecryptAttachmentBase64(jwtToken string, data string, password string) (string, error) {
// 	_, err := a.Auth.RequireAuth(jwtToken)
// 	if err != nil {
// 		return "", err
// 	}

// 	data, errR 	:= a.Vault.DecryptAttachmentBase64(data, password)
// 	if errR != nil {
// 		return "", err
// 	}

// 	return base64.StdEncoding.EncodeToString(data), nil
// }

func (a *App) GetIPFSFile(jwtToken string, cid string, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}

	a.Logger.Info("App - GetIPFSFile processing")
	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	
	ipfsQuery, err := a.Vault.GetIPFSFile(vault_ui.GetIPFSFileRequest{
		UserID: claims.UserID,
		CID:    cid,
		Password: password,
		Vault: *vault,
	})
	if err != nil {
		a.Logger.Error("App - GetIPFSFile - error: %v", err)
		return "", err
	}
	

	return base64.StdEncoding.EncodeToString(ipfsQuery), nil
}
func (a *App) UploadToIPFS(jwtToken string, filePath string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadToIPFS - error: %v", err)
		return "", err
	}

	// Simulate upload progress (integrate with your IPFS client for real progress)
	current := 70
	for i := 1; i <= 20; i++ {
		current += 1
		runtime.EventsEmit(a.ctx, "progress-update", current)
		time.Sleep(50 * time.Millisecond) // Simulate; use actual IPFS progress
	}

	cid, err := a.Vault.UploadToIPFS(claims.UserID, filePath)
	runtime.EventsEmit(a.ctx, "progress-update", 95) // Near complete
	if err != nil {
		a.Logger.Error("App - UploadToIPFS - error: %v", err)
		return "", err
	}
	return cid, nil
}
func (a *App) DownloadAttachment(jwtToken string,  password string, cid string, ext string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}

	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}

	return a.Vault.DownloadAttachment(context.Background() , vault_ui.DownloadAttachmentRequest{
		UserID: claims.UserID,
		Vault: *vault,
		CID: cid,
		Password: password,
		Ext: ext,
	})
}

func (a *App) UploadAttachmentToIPFS(jwtToken string, data []uint8, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	// Get Vault ==============================
	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	a.Logger.LogPretty("App - UploadAttachmentToIPFS - vault", vault)

	filePath, err := a.Vault.UploadAttachementToIPFS(claims.UserID, vault_ui.UploadAttachRequest{
		Data:               data,
		VaultName:          vault.Name,
		UserSubscriptionID: vault.UserSubscriptionID,
		Password:           password,
	})
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	return filePath, nil
}
func (a *App) UploadAttachmentToIPFSWithEncryption(jwtToken string, data []uint8, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	// Get Vault ==============================
	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	a.Logger.LogPretty("App - UploadAttachmentToIPFS - vault", vault)

	filePath, err := a.Vault.UploadAttachementToIPFSWithEncryption(claims.UserID, vault_ui.UploadAttachRequest{
		Data:               data,
		VaultName:          vault.Name,
		UserSubscriptionID: vault.UserSubscriptionID,
		Password:           password,
	})
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	return filePath, nil
}

func (a *App) CreateStellarCommit(jwtToken string, cid string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - CreateStellarCommit - error: %v", err)
		return "", err
	}

	// Quick commit with final progress
	runtime.EventsEmit(a.ctx, "progress-update", 100)
	return a.Vaults.CreateStellarCommit(claims.UserID, cid)
}
func (a *App) EncryptVault(jwtToken string, password string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - EncryptVault - error: %v", err)
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
		a.Logger.Error("App - EncryptVault - error: %v", err)
		return "", err
	}

	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
		"percent": 70,
		"stage":   "encrypted",
	})

	return encryptedPath, nil
}
func (a *App) UploadAvatar(jwtToken string, vaultName string, avatar []byte) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadAvatar - error: %v", err)
		return "", err
	}
	a.Logger.Info("App - UploadAvatar - vaultName", vaultName)
	return a.Vault.UploadAvatar(claims.UserID, vaultName, avatar)
}

// -----------------------------
// Vault Avatar - Vault Attachments
// -----------------------------
func (a *App) GetVaultAvatar(jwtToken string, vaultName string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - GetVaultAvatar - error: %v", err)
		return "", err
	}
	vault, err := a.Vault.GetVault(claims.UserID, vaultName)
	if err != nil {
		a.Logger.Error("App - GetVaultAvatar - error: %v", err)
		return "", err
	}
	a.Logger.LogPretty("App - GetVaultAvatar - vault", vault)
	return vault.Avatar, nil
}
func (a *App) LoadAvatar(jwtToken string, vaultName string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - LoadAvatar - error: %v", err)
		return "", err
	}

	avatar, err := a.Vault.LoadAvatar(claims.UserID, vaultName)
	if err != nil {
		a.Logger.Error("App - LoadAvatar - error: %v", err)
		return "", err
	}
	return avatar, nil
}
// upload First local (then ipfs)
func (a *App) UploadAttachments(jwtToken string, vaultName string, entryType string, raw json.RawMessage, attachments vault_dto.SelectedAttachments) (*vaults_domain.VaultEntry, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadAttachment - error: %v", err)
		return nil, err
	}
	a.Logger.Info("App - UploadAttachment - vaultName", vaultName)
	ve, err := a.Vault.UpdateEntryWithAttachments(claims.UserID, entryType, raw, vaultName, attachments)
	if err != nil {
		a.Logger.Error("App - UploadAttachment - error: %v", err)
		return nil, err
	}

	return ve, nil
}
func (a *App) LoadAttachment(jwtToken string, vaultName string, hash string) (string, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - LoadAttachment - error: %v", err)
		return "", err
	}
	a.Logger.Info("App - LoadAttachment - vaultName", vaultName)
	return a.Vault.LoadAttachment(claims.UserID, vaultName, hash)
}

func (a *App) GetVault(userID string) (map[string]interface{}, error) {
	user, err := a.Identity.FindUserById(a.ctx, userID)
	if err != nil {
		a.Logger.Error("App - GetVault - error: %v", err)
		return nil, err
	}
	session, err := a.Vault.GetSession(userID)
	if err != nil {
		a.Logger.Error("App - GetVault - error: %v", err)
		return nil, err
	}

	response := map[string]interface{}{
		"User":                user,
		"role":                "user",
		"Vault":               session.Vault,
		"SharedEntries":       []vaults_domain.VaultEntry{},
		"VaultRuntimeContext": *session.Runtime,
		"LastCID":             session.LastCID,
		"Dirty":               session.Dirty,
	}
	return response, nil
}
func (a *App) GetVaultFromCloud(jwtToken string, vaultName string) (*tracecore_types.Vault, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
		return nil, err
	}
	// FETCH SUBSCRIPTION
	sub, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - GetVaultFromCloud - sub", sub)

	response, err := a.Vault.GetVaultFromCloud(sub.ID)
	if err != nil {
		a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
		return nil, err
	}
	utils.LogPretty("App - GetVaultFromCloud - response", response.Data)
	return &response.Data, nil
}
func (a *App) GetSubscriptionFromCloud(jwtToken string, vaultName string) (*subscription_domain.Subscription, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - GetSubscriptionFromCloud - error: %v", err)
		return nil, err
	}
	a.Logger.Info("App - GetSubscriptionFromCloud - claims", claims)

	// FETCH SUBSCRIPTION
	sub, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - GetVaultFromCloud - sub cloud", sub)

	subCloud, err := a.Vault.TracecoreClient.GetSubscriptionByID(context.Background(), sub.ID)
	if err != nil {
		a.Logger.Error("App - GetSubscriptionFromCloud - error: %v", err)
		return nil, err
	}
	utils.LogPretty("App - GetSubscriptionFromCloud - response", subCloud)
	return subCloud, nil
}

// -----------------------------
// Link shares
// -----------------------------
type CreateLinkShareOutput struct {
	Data   *share_domain.LinkShare `json:"data"`
	Status string                  `json:"status"`
	Error  string                  `json:"error"`
	Code   string                  `json:"code"`
}

func (a *App) CreateLinkShare(payload share_application_dto.LinkShareCreateRequest, jwtToken string) (*CreateLinkShareOutput, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - CreateLinkShare - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - CreateLinkShare - payload", payload)
	output := CreateLinkShareOutput{}
	output.Data, err = a.Vaults.CreateLinkShare(claims.Email, payload)
	if err != nil {
		a.Logger.Error("App - CreateLinkShare - error: %v", err)
		output.Error = err.Error()
		output.Code = "500"
		return nil, err
	}
	output.Code = "200"
	output.Status = "success"
	return &output, nil
}

type ListLinkSharesByMeResponse struct {
	Data       *[]tracecore.WailsLinkShare `json:"data"`
	Status     string                      `json:"status"`
	Error      string                      `json:"error"`
	StatusCode string                      `json:"status_code"`
}

func (a *App) ListLinkSharesByMe(jwtToken string) (*ListLinkSharesByMeResponse, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListLinkSharesByMe - error: %v", err)
		return nil, err
	}
	res := ListLinkSharesByMeResponse{}
	res.Data, err = a.Vaults.ListLinkSharesByMe(claims.Email)
	if err != nil {
		a.Logger.Error("App - ListLinkSharesByMe - error: %v", err)
		res.Error = err.Error()
		res.StatusCode = "500"
		return nil, err
	}
	res.StatusCode = "200"
	res.Status = "success"

	a.Logger.LogPretty("App - ListLinkSharesByMe - res", res)
	return &res, nil
}
func (a *App) ListLinkSharesWithMe(jwtToken string) (*[]tracecore.WailsLinkShare, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListLinkSharesWithMe - error: %v", err)
		return nil, err
	}
	return a.Vaults.ListLinkSharesWithMe(claims.Email)
}

// func (a *App) DeleteLinkShare(jwtToken string, shareID string) (string, error) {
// 	claims, err := a.RequireAuth(jwtToken)
// 	if err != nil {
// 		return "", err
// 	}
// 	return a.Vaults.DeleteLinkShare(claims.UserID, shareID)
// }

// -----------------------------
// Cryptographic shares
// -----------------------------
type CreateShareInput struct {
	Payload  handlers.CreateShareEntryPayload `json:"payload"`
	JwtToken string                           `json:"jwtToken"`
}

func (a *App) CreateShare(input CreateShareInput) (*share_domain.ShareEntry, error) {
	claims, err := a.RequireAuth(input.JwtToken)
	if err != nil {
		a.Logger.Error("App - CreateShare - error: %v", err)
		return nil, err
	}
	return a.Vaults.CreateShareEntry(context.Background(), input.Payload, claims.UserID, claims.Email, *a.AppConfigHandler, a.config.ANCHORA_SECRET)
}

// Cryptographic share by me
func (a *App) ListSharedEntries(jwtToken string) (*[]share_domain.ShareEntry, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListSharedEntries - error: %v", err)
		return nil, fmt.Errorf("ListSharedEntries - auth failed: %w", err)
	}

	entries, err := a.Vaults.ListSharedEntries(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - ListSharedEntries - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - ListSharedEntries - Cryptographic entries: %v", entries)

	return &entries, nil
}

// Cryptographic share with by me
func (a *App) ListReceivedShares(jwtToken string) (*[]share_domain.ShareEntry, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListReceivedShares - error: %v", err)
		return nil, err
	}

	entries, err := a.Vaults.ListReceivedShares(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - ListReceivedShares - error: %v", err)
		return nil, err
	}

	return &entries, nil // Wails wants pointer
}
func (a *App) GetShareForAccept(jwt, shareID string) (*share_domain.ShareAcceptData, error) {
	claims, err := a.Auth.RequireAuth(jwt)
	if err != nil {
		a.Logger.Error("App - GetShareForAccept - error: %v", err)
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
		a.Logger.Error("App - RejectShare - error: %v", err)
		return nil, err
	}

	return a.Vaults.RejectShare(context.Background(), claims.UserID, shareID)
}
func (a *App) AddReceiver(jwtToken string, payload share_application.AddReceiverInput) (*share_application.AddReceiverResult, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AddReceiver - error: %v", err)
		return nil, err
	}

	return a.Vaults.AddReceiver(context.Background(), claims.UserID, payload)
}

func (a *App) AddRecipient(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AddRecipient - error: %v", err)
		return nil, err
	}

	var addRecipRequest share_application_dto.AddRecipientRequest
	if err := json.Unmarshal(raw, &addRecipRequest); err != nil {
		a.Logger.Error("App - AddRecipient - error: %v", err)
		return nil, err
	}
	return a.Vaults.AddRecipient(context.Background(), claims.UserID, addRecipRequest, *a.AppConfigHandler, a.config.ANCHORA_SECRET)
}

func (a *App) UpdateRecipient(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UpdateRecipient - error: %v", err)
		return nil, err
	}

	var updateRecipRequest share_application_dto.UpdateRecipientRequest
	if err := json.Unmarshal(raw, &updateRecipRequest); err != nil {
		a.Logger.Error("App - UpdateRecipient - error: %v", err)
		return nil, err
	}
	return a.Vaults.UpdateRecipient(context.Background(), claims.UserID, updateRecipRequest)
}

func (a *App) RevokeRecipient(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - RevokeRecipient - error: %v", err)
		return nil, err
	}

	var revokeRecipRequest share_application_dto.UpdateRecipientRequest
	if err := json.Unmarshal(raw, &revokeRecipRequest); err != nil {
		a.Logger.Error("App - RevokeRecipient - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - RevokeRecipient - request: %v", revokeRecipRequest)
	return a.Vaults.RevokeRecipient(context.Background(), claims.UserID, revokeRecipRequest)
}

func (a *App) RevokeShare(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - RevokeShare - error: %v", err)
		return nil, err
	}

	var revokeShareRequest share_application_dto.UpdateRecipientRequest
	if err := json.Unmarshal(raw, &revokeShareRequest); err != nil {
		a.Logger.Error("App - RevokeShare - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - RevokeShare - request: %v", revokeShareRequest)
	return a.Vaults.RevokeShare(context.Background(), claims.UserID, revokeShareRequest, *a.AppConfigHandler)
}

// -----------------------------
// Vault Config
// -----------------------------
type GenerateApiKeyInput struct {
	Password string `json:"password"`
	JwtToken string `json:"jwtToken"`
}
type GenerateApiKeyOutput struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func (a *App) GenerateApiKey(input GenerateApiKeyInput) (*GenerateApiKeyOutput, error) {
	claims, err := a.Auth.RequireAuth(input.JwtToken)
	if err != nil {
		a.Logger.Error("❌ GenerateApiKey - Failed to authenticate user: %v", err)
		return nil, err
	}
	userID := claims.UserID
	a.Logger.Info("GenerateApiKey - user id %s:", userID)

	// -------------------------------------------------------------------------------------------------
	// Stellar - Create Stellar account keypair with no friendbot funding
	// -------------------------------------------------------------------------------------------------
	account, err := a.StellarService.OnGenerateApiKey(input.Password)
	if err != nil {
		a.Logger.Error("❌ GenerateApiKey - Stellar account creation failed: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("✅ GenerateApiKey - Stellar account Keypair created & funded: %s", account)

	// -------------------------------------------------------------------------------------------------
	// Identity - save stellar public key to user identity
	// -------------------------------------------------------------------------------------------------
	identityUser, err := a.Identity.OnGenerateApiKey(a.ctx, userID, account.PublicKey)
	if err != nil {
		a.Logger.Error("❌ App - GenerateApiKey - failed to find user %s: %v", userID, err)
		return nil, err
	}
	a.Logger.LogPretty("✅ App - GenerateApiKey - identity user updated: %s", identityUser)

	// -------------------------------------------------------------------------------------------------
	// UserConfig - Save stellarAccount in user config
	// -------------------------------------------------------------------------------------------------
	stellarAccount := app_config_domain.NewStellarAccountConfigOnGeneratedApiKey(account)
	updatedUserCfg, err := a.AppConfigHandler.OnGenerateApiKey(userID, stellarAccount)
	if err != nil {
		a.Logger.Error("❌ GenerateApiKey - Failed to update user config: %v", err)
		return nil, err
	}
	a.Logger.Debug("✅ GenerateApiKey - User config updated: %s", updatedUserCfg)

	// -------------------------------------------------------------------------------------------------
	// Vault - save UserConfig to user vault
	// -------------------------------------------------------------------------------------------------
	if err := a.Vault.OnGenerateApiKey(context.Background(), vault_ui.OnGenerateApiKeyParams{
		UserID:     userID,
		UserConfig: *updatedUserCfg,
	}); err != nil {
		a.Logger.Error("❌ App - GenerateApiKey - failed to find user %s: %v", userID, err)
		return nil, err
	}
	a.Logger.LogPretty("✅ App - GenerateApiKey - user session updated: %s", updatedUserCfg)

	// -------------------------------------------------------------------------------------------------
	// Ankhora cloud - add public key to customer
	// -------------------------------------------------------------------------------------------------
	response, err := a.Vault.TracecoreClient.AddPublicKeyToCustomer(context.Background(), tracecore_types.AddPublicKeyToCustomerRequest{
		PublicKey: account.PublicKey,
		Email:     identityUser.Email,
	})
	if err != nil {
		a.Logger.Error("❌ App - GenerateApiKey - failed to add public key to customer: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("✅ App - GenerateApiKey - public key added to customer: %s", response)

	return &GenerateApiKeyOutput{
		PublicKey:  account.PublicKey,
		PrivateKey: account.PrivateKey,
	}, nil
}

func (a *App) IsVaultDirty(jwtToken string) (bool, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return false, err
	}
	return a.Vault.IsMarkedDirty(claims.UserID), nil
}

// -----------------------------
// Billing - Subscription
// -----------------------------
// GetPendingPaymentRequests returns all pending payment requests for current user
func (a *App) GetPendingPaymentRequests(jwtToken string) ([]*billing_domain.PaymentRequest, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.BillingHandler.GetPendingPaymentRequestsByUserID(a.ctx, claims.UserID)
}

type ClientPaymentRequest struct {
	PaymentRequestID      string `json:"payment_request_id"`
	StripePaymentMethodID string `json:"stripe_payment_method_id"`
}

// ProcessEncryptedPayment processes payment using decrypted card data
func (a *App) ProcessEncryptedPayment(req *ClientPaymentRequest) error {
	// return a.billingService.HandleClientInitiatedPayment(a.ctx, req)
	fmt.Println("✅ ProcessEncryptedPayment")
	return nil
}

type Subscription struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
}

// GetSubscriptionDetails returns current subscription details
func (a *App) GetSubscriptionDetails(jwtToken string) (*Subscription, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	fmt.Println("userID", claims.UserID)
	// return a.subscriptionService.GetActiveSubscription(a.ctx, claims.UserID)
	return &Subscription{
		ID:          "1",
		Amount:      10,
		Currency:    "USD",
		Description: "Monthly subscription",
		Status:      "active",
		UserID:      claims.UserID,
	}, nil
}

// CancelSubscription cancels current subscription
func (a *App) CancelSubscription(jwtToken string, reason string) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	fmt.Println("userID", claims.UserID)
	// return a.subscriptionService.CancelSubscription(a.ctx, claims.UserID, reason)
	return nil
}

type UpdatePaymentMethodRequest struct {
	UserID        string `json:"user_id"`
	PaymentMethod string `json:"payment_method"`
}

// UpdatePaymentMethod updates payment method for subscription
func (a *App) UpdatePaymentMethod(jwtToken string, req *UpdatePaymentMethodRequest) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	fmt.Println("userID", claims.UserID)
	req.UserID = claims.UserID
	// return a.subscriptionService.UpdatePaymentMethod(a.ctx, req)
	return nil
}

type PaymentHistory struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
}

// GetBillingHistory returns payment history
func (a *App) GetBillingHistory(jwtToken string, limit int) (*tracecore_types.CloudResponse[[]tracecore_types.PaymentHistory], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	fmt.Println("userID", claims.UserID)

	// FETCH SUBSCRIPTION
	sub, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
		return nil, err
	}

	response, err := a.BillingHandler.GetPaymentHistory(a.ctx, sub.ID, limit)
	if err != nil {
		return nil, err
	}
	a.Logger.LogPretty("✅ App - GetBillingHistory - payment history fetched: %v", response)
	return response, nil
}

type Receipt struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
}

// DownloadReceipt downloads blockchain-verified receipt
func (a *App) DownloadReceipt(jwtToken string, paymentID string) (*Receipt, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	fmt.Println("userID", claims.UserID)
	// return a.billingService.GenerateReceipt(a.ctx, claims.UserID, paymentID)
	return nil, nil
}

type StorageUsage struct {
	Used   int    `json:"used"`
	Quota  int    `json:"quota"`
	UserID string `json:"user_id"`
}

// GetStorageUsage returns current storage usage vs quota
func (a *App) GetStorageUsage(jwtToken string, tier subscription_domain.SubscriptionTier) (*tracecore_types.CloudResponse[tracecore_types.StorageUsageResponse], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------------------------------------------------
	// UserConfig - Get user config
	// -------------------------------------------------------------------------------------------------
	// userConfig, err := a.AppConfigHandler.GetUserConfigByUserID(claims.UserID)
	// if err != nil {
	// 	a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
	// 	return nil, err
	// }

	// // -------------------------------------------------------------------------------------------------
	// //  Build user challenge - signature
	// // -------------------------------------------------------------------------------------------------
	// challenge, err := a.RequestChallenge(blockchain.ChallengeRequest{PublicKey: userConfig.StellarAccount.PublicKey})
	// if err != nil {
	// 	a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
	// 	return nil, err
	// }
	// a.Logger.LogPretty("✅ App - GetVaultFromCloud - challenge: %v", challenge)
	// signature, err := blockchain.SignActorWithStellarPrivateKey(userConfig.StellarAccount.PrivateKey, challenge.Challenge)
	// if err != nil {
	// 	a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
	// 	return nil, err
	// }
	// a.Logger.LogPretty("✅ App - GetVaultFromCloud - signature: %v", signature)

	// -------------------------------------------------------------------------------------------------
	//  Get user vault
	// -------------------------------------------------------------------------------------------------
	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		a.Logger.Error("App - SynchronizeVault - error: %v", err)
		return nil, err
	}

	return a.SubscriptionHandler.GetStorageUsage(
		a.ctx,
		vault.UserSubscriptionID,
		tier,
	)
}

type UpgradeRequest struct {
	UserID        string `json:"user_id"`
	NewTier       string `json:"new_tier"`
	PaymentMethod string `json:"payment_method"`
}

// UpgradeSubscription upgrades to a higher tier
func (a *App) UpgradeSubscription(jwtToken string, req *UpgradeRequest) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	req.UserID = claims.UserID
	return a.SubscriptionHandler.HandleUpgrade(a.ctx, req.UserID, req.NewTier, req.PaymentMethod)
}

// ReactivateSubscription reactivates a cancelled subscription
func (a *App) ReactivateSubscription(jwtToken string, tier string, paymentMethod string) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	fmt.Println("userID", claims.UserID)
	// return a.subscriptionService.ReactivateSubscription(a.ctx, claims.UserID, tier, paymentMethod)
	return nil
}

func (a *App) getCurrentUserID() string {
	// Get user ID from session/context
	// This would be set during authentication
	return a.ctx.Value("user_id").(string)
}
func (a *App) GetUserVaultKey(jwtToken string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}
	fmt.Println("userID", claims.UserID)
	// return a.subscriptionService.GetUserVaultKey(a.ctx, claims.UserID)
	return "", nil
}

// -----------------------------
// User
// -----------------------------
func (a *App) EditUserInfos(jwtToken string, req *identity_dtos.EditUserInfosRequest) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	fmt.Println("userID", claims.UserID)
	// return a.Identity.EditUserInfos(a.ctx, claims.UserID, req)
	return nil
}

func (a *App) FetchUsers() ([]models.UserDTO, error) {
	users, err := a.OnBoardingHandler.FetchUsers()
	if err != nil {
		a.Logger.Error("APP - FetchUsers -failed to load all vault users")
		return nil, err
	}
	var userDTOs []models.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, models.UserDTO{
			ID:              user.ID,
			Email:           user.Email,
			Role:            "user",
			LastConnectedAt: time.Now().Format("2006-01-02 15:04:05"), // Should be from the db not hardcoded
		})
	}

	return userDTOs, nil
}

// -----------------------------
// Helpers
// -----------------------------

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

// Wails needs this to generate Entries struct in TypeScript
func (a *App) DummyExposeEntries(e models.Entries) models.Entries {
	return e
}
func loadConfig() config {
	return config{
		db: struct{ dsn string }{
			dsn: os.Getenv("DB_DSN"), // or default
		},
		stripe: struct {
			secret string
			key    string
		}{
			secret: os.Getenv("STRIPE_SECRET"),
			key:    os.Getenv("STRIPE_SECRET"),
		},
		StellarNetwork:     os.Getenv("STELLAR_NETWORK"),
		StellarHorizonURL:  os.Getenv("STELLAR_HORIZON_URL"),
		StellarAssetCode:   os.Getenv("STELLAR_ASSET_CODE"),
		StellarAssetIssuer: os.Getenv("STELLAR_ASSET_ISSUER"),
		IPFSClient:         os.Getenv("IPFS_CLIENT"),
		IPFSGateway:        os.Getenv("IPFS_GATEWAY"),
		IPFSNetwork:        os.Getenv("IPFS_NETWORK"),
		Branch:             os.Getenv("BRANCH"),
		EncryptionPolicy:   os.Getenv("ENCRYPTION_POLICY"),
		TracecoreURL:       os.Getenv("TRACECORE_URL"),
		TracecoreToken:     os.Getenv("TRACECORE_TOKEN"),
		CloudBackURL:       os.Getenv("CLOUD_BACK_URL"),
		CloudFrontURL:      os.Getenv("CLOUD_FRONT_URL"),
		ANCHORA_SECRET:     os.Getenv("ANCHORA_SECRET"),
		KEYRING_PATH:        os.Getenv("KEYRING_PATH"),
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

func (a *App) startup(ctx context.Context) {
	fmt.Println("App has started.")
	a.ctx = ctx
	// Fired every time the app regains focus
	runtime.EventsOn(ctx, "wails:window:focus", func(_ ...interface{}) {
		go a.CheckPaymentOnResume()
	})
}

// Simple method to open a URL in the system browser
func (a *App) OpenGoogle() {
	if a.ctx == nil {
		log.Println("❌ Context not set!")
		return
	}

	// Opens default browser to Google
	runtime.BrowserOpenURL(a.ctx, "http://164.90.213.173:4002/checkout")
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return "unknown"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Local IP:", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}

	return "unknown"
}

// main.go - FORCE clean schema
func ResetAndMigrate(db *gorm.DB) error {
	// Drop problematic tables
	// db.Migrator().DropTable(&persistence.UserConfigMapper{})
	// db.Migrator().DropTable(&app_config_domain.SharingRule{})
	db.Migrator().DropTable(&onboarding_domain.AppState{})

	// Recreate with correct types
	return db.AutoMigrate(
		// &persistence.UserConfigMapper{},
		// &app_config_domain.SharingRule{},
		// &app_config_domain.StellarAccountConfig{},
		&onboarding_domain.AppState{},
	)
}
func (a *App) SetStorageMode(JwtToken string, mode string) {
	claims, err := a.Auth.RequireAuth(JwtToken)
	if err != nil {
		a.Logger.Error("❌ GenerateApiKey - Failed to authenticate user: %v", err)
	}
	appCfg, err := a.AppConfigHandler.GetAppConfigByUserID(context.Background(), claims.UserID)
	if err != nil {
		a.Logger.Error("❌ GenerateApiKey - Failed to authenticate user: %v", err)
	}
	appCfg.Storage.Mode = app_config.StorageMode(mode)

	a.AppConfigHandler.UpdateAppConfig(appCfg)
	a.Vault.UpdateAppConfig(claims.UserID, *appCfg)

	appCfgUpdated, err := a.AppConfigHandler.GetAppConfigByUserID(context.Background(), claims.UserID)
	if err != nil {
		a.Logger.Error("❌ GenerateApiKey - Failed to authenticate user: %v", err)
	}
	utils.LogPretty("appCfgUpdated", appCfgUpdated)
}
