package main

import (
	"context"
	"path/filepath"
	"net"
	"os/exec"
	// "encoding/base64"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"

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
	app_config_worker "vault-app/internal/config/application/worker"
	app_config_domain "vault-app/internal/config/domain"
	app_config_ui "vault-app/internal/config/ui"
	share_domain "vault-app/internal/domain/shared"
	"vault-app/internal/driver"
	"vault-app/internal/handlers"
	identity_commands "vault-app/internal/identity/application/commands"
	identity_dtos "vault-app/internal/identity/application/dtos"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	notification_center_usecases "vault-app/internal/notification_center/application/use_cases"
	notification_center_domain "vault-app/internal/notification_center/domain"
	notification_center_ui "vault-app/internal/notification_center/ui"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	realtime_client_handlers "vault-app/internal/realtime_client/application/handlers"
	realtime_client_application_services "vault-app/internal/realtime_client/application/services"
	realtime_client_infrastructure_websocket "vault-app/internal/realtime_client/infrastructure/websocket"
	share_entry_application_dto "vault-app/internal/share_entry/application"
	share_entry_use_cases "vault-app/internal/share_entry/application/use_cases"
	share_entry_domain "vault-app/internal/share_entry/domain"
	share_entry_infrastructure "vault-app/internal/share_entry/infrastructure"
	sahre_entry_ui_wails "vault-app/internal/share_entry/ui/wails"
	onboarding_ui_wails "vault-app/internal/onboarding/ui/wails"
	"vault-app/internal/registry"
	shared_realtime "vault-app/internal/shared/realtime"
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
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vault_use_cases "vault-app/internal/vault/application/usecases"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"
	vault_ui "vault-app/internal/vault/ui"
	// "vault-app/internal/logger/logger"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_domain "vault-app/internal/channel/domain"
	channel_eventbus "vault-app/internal/channel/infrastructure/eventbus"
	channel_ui "vault-app/internal/channel/ui"
	collaboration_ui "vault-app/internal/collaboration/ui"
	"vault-app/internal/models"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_infrastructure_eventbus "vault-app/internal/thread/infrastructure/eventbus"
	thread_ui "vault-app/internal/thread/ui"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_infrastructure_eventbus "vault-app/internal/workspace/infrastructure/eventbus"
	workspace_ui "vault-app/internal/workspace/ui"

	_ "github.com/mattn/go-sqlite3"
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
	CloudURL                  string
	CloudBackURL              string
	CloudFrontURL             string
	ANKHORA_WEBSOCKET_GATEWAY string

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
	AppConfigHandler *app_config_ui.AppConfigHandler
	// Auth                      *handlers.AuthHandler
	AuthHandler               *auth_ui.AuthHandler
	BillingHandler            *billing_ui.BillingHandler
	ConnectWithStellarHandler *stellar_recovery_ui_api.StellarRecoveryHandler
	EntryRegistry             *registry.EntryRegistry
	CryptographicShareHandler *sahre_entry_ui_wails.CryptographicShareHandler
	LinkShareHandler          *sahre_entry_ui_wails.LinkShareHandler
	Identity                  *identity_ui.IdentityHandler
	OnBoardingHandler         *onboarding_ui_wails.OnBoardingHandler
	NotificationCenterHandler *notification_center_ui.NotificationHandler
	ShareEntryHandler         *sahre_entry_ui_wails.ShareEntryHandler
	StellarService            *blockchain.StellarService
	StellarRecoveryHandler    *stellar_recovery_ui_api.StellarRecoveryHandler
	SubscriptionHandler       *subscription_ui_wails.SubscriptionHandler
	Vault                     *vault_ui.VaultHandler
	// Vaults                    *handlers.VaultHandler

	// C3 Handlers
	WorkspaceHandler     *workspace_ui.WorkspaceHandler
	ChannelHandler       *channel_ui.ChannelHandler
	ThreadHandler        *thread_ui.ThreadHandler
	CollaborationHandler *collaboration_ui.CollaborationHandler

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
	// Load .env deterministically: search from executable directory upward, then project root, then cwd
	
	// 1. Try executable directory and parents (works for .app bundle and go run)
	execDir := filepath.Dir(os.Args[0])
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(execDir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			if err := godotenv.Load(candidate); err == nil {
				log.Printf("✅ Loaded .env from %s", candidate)
				break
			}
		}
		execDir = filepath.Dir(execDir)
		if execDir == "/" || execDir == "." {
			break
		}
	}

	// 2. Try known project root (development)
	projectRoot := "/Users/apple/sites/ankhora-dvault"
	if _, err := os.Stat(filepath.Join(projectRoot, ".env")); err == nil {
		if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err == nil {
			log.Printf("✅ Loaded .env from project root %s", projectRoot)
		}
	}

	// 3. Fallback: current working directory
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found in executable hierarchy or project root; trying cwd")
		_ = godotenv.Load()
	}

	// Set deterministic defaults from application config (ANKHORA_CLOUD_BACK_URL is authoritative)
	if os.Getenv("TRACECORE_URL") == "" {
		os.Setenv("TRACECORE_URL", "http://localhost:4001/api")
	}
	// Do NOT set a fake TRACECORE_TOKEN. An empty token means
	// no Authorization header is sent, which is correct until
	// Cloud authentication actually succeeds.
	if os.Getenv("ANKHORA_CLOUD_BACK_URL") == "" {
		os.Setenv("ANKHORA_CLOUD_BACK_URL", "http://localhost:4001/api")
	}
	if os.Getenv("CLOUD_BACK_URL") == "" {
		os.Setenv("CLOUD_BACK_URL", "http://localhost:4001/api")
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
	// Use ANKHORA_CLOUD_BACK_URL (via cfg.CloudBackURL) as the authoritative Cloud base URL
	tracecoreClient := tracecore.NewTracecoreClient(cfg.CloudBackURL, cfg.TracecoreToken, cfg.CloudFrontURL, cfg.CloudBackURL)

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
	// vaults := handlers.NewVaultHandler(*db, ipfs, reg, sessions, appLogger, tracecoreClient, *runtimeCtxLegacy)
	// onboardingUserRepo := onboarding_persistence.NewGormUserRepository(db.DB)
	// auth := handlers.NewAuthHandler(*db, vaults, ipfs, appLogger, tracecoreClient, cfg.auth, onboardingUserRepo)

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
		tracecoreClient,
		cfg.KEYRING_PATH,
	)

	appConfigHandler.SetVaultHandler(*vaultHandler)

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

	appConfigHandler.SetOnboardingHandler(*onBoardingHandler)

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
	// Share Entry
	// -------------------------------------------------------------------------------------------------
	entrySnapshotService := share_entry_infrastructure.NewEntrySnapshotService(
		*appLogger,
		vaultHandler,
	)
	cryptographicRepository := share_entry_infrastructure.NewGormShareRepository(db.DB)
	evtDispatcher := share_entry_infrastructure.InitializeEventDispatcher()
	cryptoAESUC := share_entry_use_cases.NewShareUseCaseAES(cryptographicRepository, tracecoreClient, evtDispatcher, &vault_infrastructure_crypto.AESService{}, entrySnapshotService)
	cryptographicShareHandler := sahre_entry_ui_wails.NewCryptographicShareHandler(*cryptoAESUC, *db.DB, evtDispatcher, appLogger, tracecoreClient)

	linkShareRepository := share_entry_infrastructure.NewGormShareRepository(db.DB)
	linkShareUC := share_entry_use_cases.NewLinkShareUseCase(linkShareRepository, tracecoreClient, evtDispatcher, &blockchain.CryptoService{})
	linkShareHandler := sahre_entry_ui_wails.NewLinkShareHandler(*linkShareUC, appLogger)

	shareEntryHandler := sahre_entry_ui_wails.NewShareEntryHandler(
		tracecoreClient,
		*appLogger,
		db.DB,
		evtDispatcher,
		vaultHandler,
	)
	appLogger.Info("shareEntryHandler", shareEntryHandler)

	// -------------------------------------------------------------------------------------------------
	// Notification Center
	// // -------------------------------------------------------------------------------------------------
	notificationUsecase := notification_center_usecases.NewNotificationUseCase(tracecoreClient)
	notificationHandler := notification_center_ui.NewNotificationHandler(*notificationUsecase)
	appLogger.Info("notificationHandler", notificationHandler)

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
			// Trace cloud token in restored session (Runtime is gorm:"-" so may be nil)
			if s.Runtime != nil && s.Runtime.SessionSecrets != nil {
				if cloudJWT, ok := s.Runtime.SessionSecrets["cloud_jwt"]; ok && cloudJWT != "" {
					appLogger.Info("☁️ [CLOUD-TRACE] Session restore: user=%s cloud_token_present=true cloud_token_length=%d",
						s.UserID, len(cloudJWT))
				} else {
					appLogger.Info("☁️ [CLOUD-TRACE] Session restore: user=%s cloud_jwt=absent", s.UserID)
				}
			} else {
				appLogger.Info("☁️ [CLOUD-TRACE] Session restore: user=%s cloud_token_present=false (runtime context not persisted)", s.UserID)
			}
			if len(s.PendingCommits) > 0 {
				// for _, commit := range s.PendingCommits {
				// 	// if err := vaults.QueuePendingCommits(s.UserID, commit); err != nil {
				// 	// 	appLogger.Error("❌ Failed to queue commit for user %d: %v", s.UserID, err)
				// 	// }
				// }
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
	vaultHandler.InitializeVaultOpenedListener(appConfigHandler)
	go vaultHandler.VaultOpenedListener.Listen(ctx)
	appLogger.Info("✅ Vault opened listener started")

	// ===== New: onboarding packs =====
	packWorker := app_config_worker.NewApplyOnboardingPacksWorker(
		db.DB,
		tracecoreClient,
		appConfigHandler,
		vaultHandler,
		nil,
		appLogger,
	)
	appConfigHandler.ApplyOnboardingPacksWorker = packWorker

	// Start pending commit worker
	// vaults.StartPendingCommitWorker(ctx, 2*time.Minute)

	elapsed := time.Since(startTime)
	appLogger.Info("✅ D-Vault initialized successfully in %v", elapsed)

	// Startup:
	// ResetAndMigrate(db.DB) // Run ONCE on prod startup

	// -------------------------------------------------------------------------------------------------
	// C3 Handlers Initialization
	// -------------------------------------------------------------------------------------------------
	workspaceBus := workspace_infrastructure_eventbus.NewMemoryBus()
	createWorkspaceUC := workspace_usecase.NewCreateWorkspaceUsecase(tracecoreClient, workspaceBus)
	listWorkspaceUC := workspace_usecase.NewListWorkspaceUsecase(tracecoreClient, workspaceBus)
	workspaceHandler := workspace_ui.NewWorkspaceHandler(createWorkspaceUC, listWorkspaceUC)

	channelBus := channel_eventbus.NewMemoryEventBus()
	createChannelUC := channel_usecase.NewCreateChannelUsecase(tracecoreClient, channelBus)
	listChannelUC := channel_usecase.NewListChannelUsecase(tracecoreClient)
	getChannelUC := channel_usecase.NewGetChannelUsecase(tracecoreClient)
	updateChannelUC := channel_usecase.NewUpdateChannelUsecase(tracecoreClient)
	deleteChannelUC := channel_usecase.NewDeleteChannelUsecase(tracecoreClient)
	activateChannelUC := channel_usecase.NewActivateChannelUsecase(tracecoreClient)
	revokeChannelUC := channel_usecase.NewRevokeChannelUsecase(tracecoreClient)
	addParticipantUC := channel_usecase.NewAddParticipantUsecase(tracecoreClient)
	listParticipantsUC := channel_usecase.NewListParticipantsUsecase(tracecoreClient)
	inviteToChannelUC := channel_usecase.NewInviteToChannelUsecase(tracecoreClient)
	acceptInvitationUC := channel_usecase.NewAcceptChannelInvitationUsecase(tracecoreClient)
	channelHandler := channel_ui.NewChannelHandler(createChannelUC, listChannelUC, getChannelUC, updateChannelUC, deleteChannelUC, activateChannelUC, revokeChannelUC, addParticipantUC, listParticipantsUC, inviteToChannelUC, acceptInvitationUC)

	threadBus := thread_infrastructure_eventbus.NewMemoryBus()
	createThreadUC := thread_usecase.NewCreateThreadUsecase(tracecoreClient, threadBus)
	listThreadsUC := thread_usecase.NewListThreadsUsecase(tracecoreClient)
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(tracecoreClient)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(tracecoreClient)
	threadHandler := thread_ui.NewThreadHandler(createThreadUC, listThreadsUC, listThreadEventsUC, appendThreadEventUC)

	collaborationHandler := collaboration_ui.NewCollaborationHandler(nil, appendThreadEventUC)

	application := &App{
		AppConfigHandler: appConfigHandler,
		// Auth:                      nil, // auth,
		BillingHandler:            billingHandler,
		AuthHandler:               authHandler,
		cancel:                    cancel,
		ConnectWithStellarHandler: stellarRecoveryHandler,
		config:                    cfg,
		DB:                        *db,
		EntryRegistry:             reg,
		CryptographicShareHandler: &cryptographicShareHandler,
		LinkShareHandler:          &linkShareHandler,
		NotificationCenterHandler: notificationHandler,
		NowUTC:                    func() string { return time.Now().Format(time.RFC3339) },
		Identity:                  identityHandler,
		Logger:                    *appLogger,
		OnBoardingHandler:         onBoardingHandler,
		sessions:                  sessions, // TODO: remove legacy sessions
		ShareEntryHandler:         &shareEntryHandler,
		StellarService:            stellarService,
		StellarRecoveryHandler:    stellarRecoveryHandler,
		SubscriptionHandler:       subscriptionHandler,
		RuntimeContext:            runtimeCtxLegacy,
		Vault:                     vaultHandler, // internal/vault/ui/vault_handler.go
		WorkspaceHandler:          workspaceHandler,
		ChannelHandler:            channelHandler,
		ThreadHandler:             threadHandler,
		CollaborationHandler:      collaborationHandler,
		// Vaults:                    nil,          // vaults, // internal/handlers/vault_handler.go legacy
		version: version,
	}

	// ===== New: vault share created =====
	vaultListener := vault_use_cases.NewVaultOnShareCreatedListener(
		application,
		vaultHandler,
		identityHandler.IdentityUserRepo,
		appLogger,
		shareEntryHandler.EventDispatcher,
	)
	go vaultListener.Listen(ctx)
	appLogger.Info("✅ Vault share created listener started")

	return application
}

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

type CheckoutContext struct {
	Identity     identity_domain.IdentityChoice `json:"identity"`
	IsAnonymous  bool                           `json:"isAnonymous"`
	Rail         string                         `json:"rail"`
	Email        string                         `json:"email"`
	Tier         string                         `json:"tier"`
	Plan         string                         `json:"plan"`
	PeriodMonths string                         `json:"periodMonths"`
	Mode         string                         `json:"mode"`
	Prorate      float64                        `json:"prorate,omitempty"`
}

// GetCheckoutURL returns the cloud backend checkout page URL
func (a *App) GetCheckoutURL(ctx CheckoutContext) (CreateCheckoutResponse, error) {
	// -----------------------------
	// 0. Generate Session ID
	// -----------------------------
	sessionID := uuid.New().String()
	periodMonths := "1"

	// -----------------------------
	// 1. Generate Checkout URL
	// -----------------------------
	baseURL := a.config.CloudFrontURL + "/checkout" // your cloud page URL
	url := fmt.Sprintf(
		"%s?session-id=%s&identity=%s&rail=%s&email=%s&tier=%s&plan=%s&period-months=%s&is-anonymous=%t&mode=%s",
		baseURL, sessionID, ctx.Identity, ctx.Rail, ctx.Email, ctx.Tier, ctx.Plan, periodMonths, ctx.IsAnonymous, ctx.Mode,
	)
	if ctx.Prorate > 0 {
		url += fmt.Sprintf("&prorate=%f", ctx.Prorate)
	}

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
//
//	func (a *App) Sign(req handlers.LoginRequest) (*handlers.LoginResponse, error) {
//		return a.Auth.Login(req)
//	}
//
//	func (a *App) SignUp(setup handlers.OnBoarding) (*handlers.OnBoardingResponse, error) {
//		return a.Auth.OnBoarding(setup)
//	}
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

type CheckEmailResponse struct {
	Status      string   `json:"status"`
	AuthMethods []string `json:"auth_methods,omitempty"`
}

func (a *App) CheckEmail(email string) (*CheckEmailResponse, error) {

	/* user, err := ah.TracecoreClient.GetUserByEmail(ctx, email)
	user, err := ah.DB.GetUserByEmail(email)
	utils.LogPretty("user in checkemail", user)
	// ----------------------------
	// Case 1: User does NOT exist
	// ----------------------------
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CheckEmailResponse{
				Status:      "NEW_USER",
				AuthMethods: []string{},
			}, nil
		}

		// real error -> bubble up
		return &CheckEmailResponse{}, err
	}

	// ----------------------------
	// Case 2: User exists
	// ----------------------------
	*/
	authMethods := []string{"password"}

	// if user.PublicKey != "" {
	// 	authMethods = append(authMethods, "stellar")
	// }

	return &CheckEmailResponse{
		Status:      "EXISTS",
		AuthMethods: authMethods,
	}, nil
}

func (a *App) CheckUserEmail(email string, token string) (*tracecore_types.User, error) {

	user, err := a.Vault.TracecoreClient.GetUserByEmail(context.Background(), email)
	if err != nil {
		a.Logger.Error("❌ Vault - CheckUserEmail Failed to get user by email: %v", err)
		return nil, err
	}
	// user, err := ah.DB.GetUserByEmail(email)
	utils.LogPretty("user in checkemail", user)

	return user, nil
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
	a.Logger.Info("Identity login successful: %v", result.User)
	log.Println("[CLOUD-TRACE] A: identity login succeeded")

	// --------- Cloud Authentication ---------
	// Authenticate with Ankhora Cloud to obtain a 26-char bearer token.
	// This token is required for Cloud API calls (workspaces, channels, etc.).
	log.Println("[CLOUD-TRACE] B: Cloud authentication starting")
	if req.Email != "" && req.Password != "" {
		cloudLoginResp, cloudErr := a.Vault.TracecoreClient.Login(
			context.Background(),
			tracecore_types.LoginRequest{
				Email:    req.Email,
				Password: req.Password,
			},
		)
		if cloudErr != nil {
			log.Printf("[CLOUD-TRACE] C: Cloud authentication FAILED: %v", cloudErr)
			a.Logger.Warn("☁️ Cloud authentication failed: %v — proceeding with local session only", cloudErr)
		} else if cloudLoginResp != nil &&
			cloudLoginResp.AuthenticationToken != nil &&
			cloudLoginResp.AuthenticationToken.Token != "" {
			cloudToken := cloudLoginResp.AuthenticationToken.Token
			log.Printf("[CLOUD-TRACE] C: Cloud authentication result: authentication_token_present=true token_length=%d", len(cloudToken))
			a.Vault.TracecoreClient.SetToken(cloudToken)
			log.Printf("[CLOUD-TRACE] D: TracecoreClient.SetToken: token_present=true token_length=%d", len(cloudToken))
			a.Logger.Info("☁️ [CLOUD-AUTH] Cloud token set: length=%d fingerprint=%s",
				len(cloudToken), tracecore.TraceTokenFingerprint(cloudToken))
		} else {
			log.Printf("[CLOUD-TRACE] C: Cloud authentication returned nil/empty token resp=%v auth=%v", cloudLoginResp != nil, cloudLoginResp != nil && cloudLoginResp.AuthenticationToken != nil)
		}
	} else {
		log.Printf("[CLOUD-TRACE] C: SKIPPED — email or password empty: email_empty=%v password_empty=%v", req.Email == "", req.Password == "")
	}

	// --------- Session Warm Up ---------
	session, err := a.Vault.PrepareSession(result.User.ID)
	if err != nil {
		a.Logger.Error("❌ App - SignIn - failed to get session for user %s: %v", result.User.ID, err)
	}
	if session == nil {
		a.Logger.Error("❌ App - SignIn - failed to get session for user %s: %v", result.User.ID, err)
		// return	 nil, err
	} else {
		a.Logger.Info("Session fetched successfully")
	}

	// --------- Find user onboarding ---------
	userOnboarding, err := a.OnBoardingHandler.FindUsersUseCase.FindByEmail(result.User.Email)
	if err != nil {
		a.Logger.Error("❌ App - SignIn - failed to find user onboarding for user %s: %v", result.User.ID, err)
		return nil, err
	}
	a.Logger.Info("User onboarding found successfully: %v", userOnboarding)

	subscription, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), userOnboarding.Email)
	if err != nil {
		a.Logger.Error("❌ App - SignIn - failed to get subscription for user %s: %v", result.User.ID, err)
		return nil, err
	}
	// a.Logger.Info("Subscription fetched successfully: %v", subscription)

	// --------- Open vault ---------
	vaultRes, err := a.Vault.Open(
		context.Background(),
		vault_commands.OpenVaultCommand{
			UserID:           result.User.ID,
			Password:         req.Password,
			Session:          session,
			UserOnboardingID: userOnboarding.ID,
			Subscription:     *subscription,
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

	// ---------- Persist Cloud token in session for future restoration --------- //
	if cloudToken := a.Vault.TracecoreClient.Token; cloudToken != "" && cloudToken != "atokentochange" {
		if vaultRes.RuntimeContext != nil {
			if vaultRes.RuntimeContext.SessionSecrets == nil {
				vaultRes.RuntimeContext.SessionSecrets = make(map[string]string)
			}
			vaultRes.RuntimeContext.SessionSecrets["cloud_jwt"] = cloudToken
			a.Logger.Info("☁️ [CLOUD-AUTH] Persisted cloud token in session: length=%d", len(cloudToken))
		}
	}

	// ---------- Connect to real-time --------- //
	a.ConnectToRealtime(*result.User)

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

// RestoreCloudTokenForUser restores the Cloud bearer token for a user from their session
func (a *App) RestoreCloudTokenForUser(userID string) error {
	if a.Vault.SessionManager == nil {
		return nil
	}
	userSession, err := a.Vault.GetSession(userID)
	if err != nil {
		return err
	}
	if userSession.Runtime != nil && userSession.Runtime.SessionSecrets != nil {
		if cloudJWT, ok := userSession.Runtime.SessionSecrets["cloud_jwt"]; ok && cloudJWT != "" {
			a.Vault.TracecoreClient.SetToken(cloudJWT)
			a.Logger.Info("☁️ [CLOUD-TRACE] RestoreCloudTokenForUser: user=%s token_present=true token_length=%d",
				userID, len(cloudJWT))
		} else {
			a.Logger.Info("☁️ [CLOUD-TRACE] RestoreCloudTokenForUser: user=%s cloud_jwt=absent", userID)
		}
	} else {
		a.Logger.Info("☁️ [CLOUD-TRACE] RestoreCloudTokenForUser: user=%s runtime_context=nil", userID)
	}
	return nil
}

func (a *App) GetSession(userID string) (*GetSessionResponse, error) {
	if a.Vault.SessionManager == nil {
		return &GetSessionResponse{Error: errors.New("session manager not initialized")}, nil
	}

	userSession, err := a.Vault.GetSession(userID)
	if err != nil {
		return &GetSessionResponse{Error: err}, nil
	}

	// Restore Cloud bearer token from session if available
	if err := a.RestoreCloudTokenForUser(userID); err != nil {
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

	// vault, err := a.Vault.GetVault(claims.UserID, vaultName)
	// if err != nil {
	// 	a.Logger.Error("App - GetVaultAvatar - error: %v", err)
	// 	return nil, err
	// }

	// subscription, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	// if err != nil {
	// 	a.Logger.Error("App - GetSubscriptionByUserID - error: %v", err)
	// 	return nil, err
	// }
	// // a.Logger.LogPretty("App - GetSubscriptionByUserID - sub", subscription)
	// a.AppConfigHandler.VaultHandler = a.Vault

	_, _, config, err := a.GetAllConfigs(claims.UserID, claims.Email)
	if err != nil {
		a.Logger.Error("App - GetConfig - error: %v", err)
		return nil, err
	}

	return config, nil
}

func (a *App) GetAllConfigs(userID string, email string) (*subscription_domain.Subscription, *vaults_domain.Vault, *app_config_domain.Config, error) {
	a.Logger.LogPretty("App - GetAllConfigs - userID", userID)
	a.Logger.LogPretty("App - GetAllConfigs - email", email)

	vault, err := a.Vault.GetLatestByUserID(userID)
	if err != nil {
		a.Logger.Error("App - GetAllConfigs - error: %v", err)
		return nil, nil, nil, err
	}

	subscription, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), email)
	if err != nil {
		a.Logger.Error("App - GetAllConfigs - error: %v", err)
		return nil, nil, nil, err
	}
	// a.Logger.LogPretty("App - GetSubscriptionByUserID - sub", subscription)
	a.AppConfigHandler.VaultHandler = a.Vault

	appConfig, err := a.AppConfigHandler.GetConfig(userID, *vault, subscription)
	if err != nil {
		a.Logger.Error("App - GetAllConfigs - error: %v", err)
		return nil, nil, nil, err
	}

	return subscription, vault, appConfig, nil
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
	return claims.ToFormerModel(), nil
}

// func (a *App) RequestChallenge(req blockchain.ChallengeRequest) (blockchain.ChallengeResponse, error) {
// 	return a.Auth.RequestChallenge(req)
// }
// func (a *App) AuthVerify(req blockchain.SignatureVerification) (string, error) {
// 	return a.Auth.AuthVerify(&req)
// }

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
	utils.LogPretty("App - SynchronizeVault - jwtToken", jwtToken) // ✅ log
	claims, err := a.RequireAuth(jwtToken)
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

	// Get configs ==============================
	cfgs, err := a.GetConfig(vault.Name, jwtToken)
	if err != nil {
		a.Logger.Error("App - SynchronizeVault - error: %v", err)
		return "", err
	}
	a.Logger.LogPretty("App - SynchronizeVault - cfgs", cfgs)

	// Get User onboarding ==============================
	userOnboarding, err := a.OnBoardingHandler.UserRepo.FindByEmail(claims.Email)
	if err != nil {
		a.Logger.Error("App - GetConfigByUserID - error: %v", err)
		return "", err
	}

	// Sync Vault ==============================
	input := vault_dto.SynchronizeVaultRequest{
		UserID:         claims.UserID,
		Password:       password,
		Vault:          *vault,
		UserOnboarding: userOnboarding.ID,
		Configs:        *cfgs,
	}
	a.Vault.Ctx = a.ctx

	res, err := a.Vault.SyncVault(a.ctx, input, *&a.Vault.TracecoreClient)
	if err != nil {
		a.Logger.Error("App - SynchronizeVault - error: %v", err)
		return "", err
	}
	utils.LogPretty("App - SynchronizeVault - res", res)
	return res, err
}

// func (a *App) EncryptFile(jwtToken string, fileData string, password string) (string, error) {
// 	claims, err := a.RequireAuth(jwtToken)
// 	if err != nil {
// 		a.Logger.Error("App - EncryptFile - error: %v", err)
// 		return "", err
// 	}

// 	// Emit start progress
// 	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
// 		"percent": 0,
// 		"stage":   "encrypting",
// 	})
// 	a.Logger.LogPretty("App - EncryptFile - fileData", fileData)

// 	// Real AES-256-GCM encryption with progress
// 	encryptedPath, err := a.Vault.EncryptFile(claims.UserID, []byte(fileData), password)
// 	if err != nil {
// 		a.Logger.Error("App - EncryptFile - error: %v", err)
// 		return "", err
// 	}

// 	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
// 		"percent": 70,
// 		"stage":   "encrypted",
// 	})

//		return encryptedPath, nil
//	}
func (a *App) EncryptAttachment(jwtToken string, data []byte, password string) ([]byte, error) {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	return a.Vault.EncryptAttachment(data, password)
}
func (a *App) DecryptAttachment(jwtToken string, data []byte, password string) ([]byte, error) {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	a.Logger.Info("DecryptAttachment processing....")

	return a.Vault.DecryptAttachment(data, password)
}

func (a *App) GetIPFSFile(jwtToken string, cid string, password string, vaultName string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}

	// Get Configs ==============================
	config, err := a.GetConfig(vaultName, jwtToken)
	if err != nil {
		return "", err
	}

	// Get user onboarding ==============================
	userOnboarding, err := a.OnBoardingHandler.UserRepo.FindByEmail(claims.Email)
	if err != nil {
		a.Logger.Error("App - GetIPFSFile - error: %v", err)
		return "", err
	}

	ipfsQuery, err := a.Vault.GetIPFSFile(vault_queries.GetIPFSDataQuerry{
		UserID:           claims.UserID,
		CID:              cid,
		Password:         password,
		Configs:          *config,
		VaultName:        config.Vaults.VaultName,
		UserOnboardingID: userOnboarding.ID,
	})
	if err != nil {
		a.Logger.Error("App - GetIPFSFile - error: %v", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ipfsQuery), nil
}
func (a *App) UploadToIPFS(jwtToken string, filePath string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
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

// From IPFS
func (a *App) DownloadAttachment(jwtToken string, password string, cid string, ext string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}

	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}

	return a.Vault.DownloadAttachment(context.Background(), vault_ui.DownloadAttachmentRequest{
		UserID:   claims.UserID,
		Vault:    *vault,
		CID:      cid,
		Password: password,
		Ext:      ext,
	})
}

type DecryptCryptoShareResponse struct {
	Payload      string            `json:"payload"`
	ExpiresIn    int64             `json:"expires_in,omitempty"`
	Attachments  map[string]string `json:"attachments"`
	EncryptedKey string            `json:"encrypted_key"`
}

func (a *App) AccessDecryptVaultEntry(jwtToken string, entry tracecore_types.AccessCryptoShareRequest) (*tracecore_types.CloudResponse[tracecore_types.DecryptCryptoShareResponse], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	fmt.Println("userID", claims.UserID)

	// 0. Get user vault ==============================
	// vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	// 1. Access encrypted entry ==============================
	entry.IPAddress = GetLocalIP()
	res, err := a.Vault.AccessEncryptedEntry(a.ctx, claims.UserID, entry, a.Vault.TracecoreClient)
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
	a.Logger.LogPretty("App - DecryptVaultEntry - stellarAccount", stellarAccount)

	response, err := a.Vault.DecryptVaultEntry(context.Background(), req, a.Vault.TracecoreClient)
	if err != nil {
		return nil, err
	}
	a.Logger.LogPretty("App - DecryptVaultEntry - Entry Snapshot decrypted response", response)
	// TODO decrypt attachement from the response
	a.Logger.Info("payload hex (first 32):", hex.EncodeToString([]byte(response.Data.Payload)[:32]))

	utils.LogPretty("AccessDecryptVaultEntry - BEFORE DownloadAttachment loop", "start")
	time.Sleep(100 * time.Millisecond) // give goroutines a chance to settle
	utils.LogPretty("AccessDecryptVaultEntry - BEGINNING DownloadAttachment loop", "entered")

	// 2. Parse JSON payload into VaultNode / EntrySnapshot
	var snapshot share_domain.EntrySnapshot
	if err := json.Unmarshal([]byte(response.Data.Payload), &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Payload: %w", err)
	}

	// 👇 add this
	appConfig, err := a.AppConfigHandler.GetAppConfigByUserID(context.Background(), claims.UserID)
	if err != nil {
		return nil, err
	}
	response.Data.ExpiresIn = appConfig.AccessPolicyDuration
	utils.LogPretty("App - DecryptVaultEntry - Config", appConfig)
	// utils.LogPretty("App - DecryptVaultEntry - Final response", response.Data)

	var finalResponse tracecore_types.DecryptCryptoShareResponse

	finalResponse.Payload = response.Data.Payload
	finalResponse.ExpiresIn = response.Data.ExpiresIn + 300 // 5 minutes
	finalResponse.EncryptedKey = res.Data.EncryptedKey

	return &tracecore_types.CloudResponse[tracecore_types.DecryptCryptoShareResponse]{
		Data: finalResponse,
	}, nil
}

func (a *App) DownloadShareAttachement(req vault_dto.DownloadShareAttachmentRequest, jwtToken string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - DownloadShareAttachement - error: %v", err)
		return "", err
	}
	// 0. Get user config - Get stellar private key from user config ==============================
	userConfig, err := a.AppConfigHandler.GetUserConfigByUserID(claims.UserID)
	if err != nil {
		return "", err
	}

	// 1. Get sstellar masterkey ==============================
	stellarAccount := userConfig.StellarAccount

	// 2. Get user vault ==============================
	vault, err := a.Vault.VaultRepository.GetLatestByUserID(claims.UserID)
	if err != nil {
		utils.LogPretty("App - DownloadShareAttachement - GetVaultByName failed", err)
		return "", err
	}

	// 3. Get app config ==============================
	cfgs, err := a.GetConfig(vault.Name, jwtToken)
	if err != nil {
		return "", fmt.Errorf("❌ DownloadShareAttachement: failed to get app config %w", err)
	}

	// 4. Download attachment ==============================
	res, err := a.Vault.DownloadAttachment(context.Background(), vault_ui.DownloadAttachmentRequest{
		UserID:       claims.UserID,
		Vault:        *vault,
		CID:          req.AttachmentCID,
		Ext:          req.FileExtension,
		Password:     "password",
		PrivateKey:   stellarAccount.PrivateKey,
		EncryptedKey: req.EncryptedKey,
		Configs:      cfgs,
		IsShared:     true,
	})
	if err != nil {
		a.Logger.Error("App - DownloadShareAttachement - error: %v", err)
		return "", err
	}
	utils.LogPretty("App - DownloadShareAttachement - res", res)
	return res, nil
}

func (a *App) UploadAttachmentToIPFS(jwtToken string, data []uint8, password string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
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

	// Get Configs ==============================
	configs, err := a.GetConfig(vault.Name, jwtToken)
	if err != nil {
		return "", err
	}

	filePath, err := a.Vault.UploadAttachementToIPFS(claims.UserID, vault_ui.UploadAttachRequest{
		Configs:            *configs,
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
	claims, err := a.RequireAuth(jwtToken)
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

	// Get Subscription ==============================
	sub, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}

	// Get user onboarding ==============================
	userOnboarding, err := a.OnBoardingHandler.FindUsersUseCase.FindByEmail(claims.Email)
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}
	utils.LogPretty("App - UploadAttachmentToIPFS - userOnboarding", userOnboarding)

	// Get Configs ==============================
	configs, err := a.GetConfig(vault.Name, jwtToken)
	if err != nil {
		return "", err
	}

	filePath, err := a.Vault.UploadAttachementToIPFSWithEncryption(claims.UserID, vault_ui.UploadAttachRequest{
		Configs:            *configs,
		Data:               data,
		VaultName:          vault.Name,
		UserSubscriptionID: sub.UserID, // TODO: replace with configs.Subscription.UserID
		Password:           password,
		UserOnboarding:     userOnboarding.ID,
	})
	if err != nil {
		a.Logger.Error("App - UploadAttachmentToIPFS - error: %v", err)
		return "", err
	}

	return filePath, nil
}

func (a *App) AddAttachements(jwtToken string, req vault_dto.AddAttachementsRequest) ([]*vaults_domain.Attachment, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AddAttachement - error: %v", err)
		return nil, err
	}

	configs, err := a.GetConfig(req.VaultName, jwtToken)
	if err != nil {
		a.Logger.Error("App - AddAttachement - error: %v", err)
		return nil, err
	}

	userOnboarding, err := a.OnBoardingHandler.FindUsersUseCase.FindByEmail(claims.Email)
	if err != nil {
		a.Logger.Error("App - AddAttachement - error: %v", err)
		return nil, err
	}

	var AddAttachementsResponse []*vaults_domain.Attachment

	for _, att := range req.Attachments {
		input := vault_dto.AddAttachementRequest{
			UserID:           claims.UserID,
			Data:             att.Data,
			Password:         req.Password,
			EntryID:          req.EntryID,
			VaultName:        req.VaultName,
			UserOnboardingID: userOnboarding.ID,
			Configs:          *configs,
			Name:             att.Name,
			Size:             att.Size,
			Ext:              att.Ext,
		}

		attachment, err := a.Vault.AddAttachement(context.Background(), input)
		if err != nil {
			a.Logger.Error("App - AddAttachements - error: %v", err)
			return nil, err
		}

		AddAttachementsResponse = append(AddAttachementsResponse, attachment)

	}

	return AddAttachementsResponse, nil
}
func (a *App) UpdateAttachment(jwtToken string, attachment vaults_domain.Attachment) (*vaults_domain.Attachment, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UpdateAttachment - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - UpdateAttachment - attachment", attachment)

	res, err := a.Vault.UpdateAttachment(context.Background(), claims.UserID, attachment)
	if err != nil {
		a.Logger.Error("App - UpdateAttachment - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - UpdateAttachment - res", res)
	return res, nil
}

// func (a *App) CreateStellarCommit(jwtToken string, cid string) (string, error) {
// 	claims, err := a.RequireAuth(jwtToken)
// 	if err != nil {
// 		a.Logger.Error("App - CreateStellarCommit - error: %v", err)
// 		return "", err
// 	}

// 	// Quick commit with final progress
// 	runtime.EventsEmit(a.ctx, "progress-update", 100)
// 	return a.Vaults.CreateStellarCommit(claims.UserID, cid)
// }
// func (a *App) EncryptVault(jwtToken string, password string) (string, error) {
// 	claims, err := a.RequireAuth(jwtToken)
// 	if err != nil {
// 		a.Logger.Error("App - EncryptVault - error: %v", err)
// 		return "", err
// 	}

// 	// Emit start progress
// 	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
// 		"percent": 0,
// 		"stage":   "encrypting",
// 	})

// 	// Real AES-256-GCM encryption with progress
// 	encryptedPath, err := a.Vaults.EncryptVault(claims.UserID, password)
// 	if err != nil {
// 		a.Logger.Error("App - EncryptVault - error: %v", err)
// 		return "", err
// 	}

// 	runtime.EventsEmit(a.ctx, "progress-update", map[string]interface{}{
// 		"percent": 70,
// 		"stage":   "encrypted",
// 	})

//		return encryptedPath, nil
//	}
func (a *App) UploadAvatar(jwtToken string, vaultName string, avatar []byte) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UploadAvatar - error: %v", err)
		return "", err
	}
	a.Logger.Info("App - UploadAvatar - vaultName", vaultName)
	return a.Vault.UploadAvatar(claims.UserID, vaultName, avatar)
}

func (a *App) GetPacks(packID string) (*app_config_worker.PackDTO, error) {
	res, err := a.Vault.TracecoreClient.GetPack(context.Background(), packID)
	if err != nil {
		a.Logger.LogPretty("GetPacks - res", err)
		return nil, err
	}
	a.Logger.LogPretty("GetPacks - res", res)
	return res, nil
}
func (a *App) GetTemplate(templateID string) (*app_config_worker.TemplateDTO, error) {
	res, err := a.Vault.TracecoreClient.GetTemplate(context.Background(), templateID)
	if err != nil {
		a.Logger.LogPretty("GetGetTemplatePacks - res", err)
		return nil, err
	}
	a.Logger.LogPretty("GetTemplate - res", res)
	return res, nil
}
func (a *App) SimulatePackWorker(jwtToken string, vaultName string) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - SimulatePackWorker - Auth error: %v", err)
		return err
	}

	userOnboarding, err := a.OnBoardingHandler.FindUsersUseCase.FindByEmail(claims.Email)
	if err != nil {
		a.Logger.Error("App - SimulatePackWorker - userOnboarding error: %v", err)
		return err
	}
	a.Logger.Info("App - SimulatePackWorker - userOnboarding: ", userOnboarding)

	a.AppConfigHandler.OnApplyOnboardingPacks(claims.UserID, vaultName, userOnboarding.ID)
	return nil
}

// -----------------------------
// Vault Avatar - Vault Attachments
// -----------------------------
func (a *App) GetVaultAvatar(jwtToken string, vaultName string) (string, error) {
	claims, err := a.RequireAuth(jwtToken)
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
	claims, err := a.RequireAuth(jwtToken)
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
	claims, err := a.RequireAuth(jwtToken)
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
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - LoadAttachment - error: %v", err)
		return "", err
	}
	a.Logger.Info("App - LoadAttachment - vaultName", vaultName)
	res, err := a.Vault.LoadAttachment(claims.UserID, vaultName, hash, "string")
	if err != nil {
		a.Logger.Error("App - LoadAttachment - error: %v", err)
		return "", err
	}
	return res.Hash, nil
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

	// FETCH SUBSCRIPTION
	sub, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - GetVaultFromCloud - error: %v", err)
		return nil, err
	}

	res, err := a.Vault.TracecoreClient.GetSubscriptionByUserID(context.Background(), sub.UserID)
	if err != nil {
		a.Logger.Error("App - GetSubscriptionFromCloud - error: %v", err)
		return nil, err
	}

	// Get configs ==============================
	cfgs, err := a.GetConfig(vaultName, jwtToken)
	if err != nil {
		a.Logger.Error("App - GetSubscriptionFromCloud - error: %v", err)
		return nil, err
	}

	cfgs.Subscription = &app_config_domain.SubscriptionConfig{
		BaseVaultConfig: app_config_domain.BaseVaultConfig{
			ID:        res.Data.ID,
			UserID:    res.Data.UserID,
			VaultName: vaultName,
		},
		Plan: res.Data.Tier,
		Features: app_config_domain.FeatureFlags{
			TracecoreEnabled:        res.Data.Features.Tracecore,
			CloudBackupEnabled:      res.Data.Features.CloudBackup,
			ThreatDetectionEnabled:  res.Data.Features.ThreatDetection,
			BrowserExtensionEnabled: res.Data.Features.BrowserExtension,
			GitCLIEnabled:           res.Data.Features.GitCLI,
		},
		Limits: app_config_domain.SubscriptionLimits{
			MaxVaults:  3,
			MaxUsers:   5,
			MaxDevices: 10,
			MaxShares:  200,
		},
	}

	utils.LogPretty("App - GetSubscriptionFromCloud - cfgs", &cfgs)

	err = a.AppConfigHandler.SubscriptionConfigRepository.Update(cfgs.Subscription.ID, cfgs.Subscription)
	if err != nil {
		a.Logger.Error("App - GetSubscriptionFromCloud - error: %v", err)
		return nil, err
	}

	utils.LogPretty("App - GetSubscriptionFromCloud - response", res.Data)
	return &res.Data, nil
}

// -----------------------------
// Link shares
// -----------------------------
type CreateLinkShareOutput struct {
	Data   *share_entry_domain.LinkShare `json:"data"`
	Status string                        `json:"status"`
	Error  string                        `json:"error"`
	Code   string                        `json:"code"`
}

func (a *App) CreateLinkShare(payload share_entry_application_dto.LinkShareCreateRequest, jwtToken string) (*CreateLinkShareOutput, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - CreateLinkShare - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - CreateLinkShare - payload", payload)
	output := CreateLinkShareOutput{}
	output.Data, err = a.LinkShareHandler.CreateLinkShare(context.Background(), claims.Email, payload)
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
	res.Data, err = a.LinkShareHandler.ListLinkSharesByMe(context.Background(), claims.Email)
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
	return a.LinkShareHandler.ListLinkSharesWithMe(context.Background(), claims.Email)
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
	Payload  share_entry_application_dto.CreateShareEntryPayload `json:"payload"`
	JwtToken string                                              `json:"jwtToken"`
}

func (a *App) CreateShare(input CreateShareInput) (*share_entry_domain.ShareEntry, error) {
	claims, err := a.RequireAuth(input.JwtToken)
	if err != nil {
		a.Logger.Error("App - CreateShare - error: %v", err)
		return nil, err
	}

	a.Logger.LogPretty("App - CreateShare - AttachmentCIDs", input.Payload.AttachmentCIDs)

	userOnboarding, err := a.OnBoardingHandler.FindUsersUseCase.FindByEmail(claims.Email)
	if err != nil {
		a.Logger.Error("App - CreateShare - error: %v", err)
		return nil, err
	}

	_, _, configs, err := a.GetAllConfigs(claims.UserID, claims.Email)
	if err != nil {
		a.Logger.Error("App - CreateShare - error: %v", err)
		return nil, err
	}

	subscription, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - CreateShare - error: %v", err)
		return nil, err
	}

	return a.CryptographicShareHandler.CreateShareEntry(
		context.Background(),
		input.Payload,
		claims.UserID,
		claims.Email,
		a.AppConfigHandler,
		a.config.ANCHORA_SECRET,
		a.Vault,
		a.Vault.TracecoreClient,
		userOnboarding.ID,
		*configs,
		subscription.UserID,
	)
}

// Cryptographic share by me
func (a *App) ListSharedEntries(jwtToken string) (*[]share_entry_domain.ShareEntry, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListSharedEntries - error: %v", err)
		return nil, fmt.Errorf("ListSharedEntries - auth failed: %w", err)
	}

	if a.ShareEntryHandler == nil {
		a.Logger.LogPretty("App - ListSharedEntries - ShareEntryHandler is nil", nil)
		return nil, errors.New("ShareEntryHandler is nil")
	}
	a.Logger.LogPretty("App - ListSharedEntries - claims", claims.Email)

	entries, err := a.CryptographicShareHandler.ListSharedEntries(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - ListSharedEntries - error: %v", err)
		return nil, err
	}

	return &entries, nil
}

// Cryptographic share with by me
func (a *App) ListReceivedShares(jwtToken string) (*[]share_entry_domain.ShareEntry, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListReceivedShares - error: %v", err)
		return nil, err
	}

	entries, err := a.CryptographicShareHandler.ListReceivedShares(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - ListReceivedShares - error: %v", err)
		return nil, err
	}

	return &entries, nil // Wails wants pointer
}
func (a *App) GetShareForAccept(jwt, shareID string) (*share_entry_domain.ShareAcceptData, error) {
	claims, err := a.RequireAuth(jwt)
	if err != nil {
		a.Logger.Error("App - GetShareForAccept - error: %v", err)
		return nil, err
	}

	return a.CryptographicShareHandler.CryptographicShareUseCase.GetShareForAccept(
		context.Background(),
		claims.UserID,
		shareID,
	)
}
func (a *App) AddReceiver(jwtToken string, inputReceiver share_entry_use_cases.AddReceiverInput) (*share_entry_use_cases.AddReceiverResult, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AddReceiver - error: %v", err)
		return nil, err
	}

	return a.CryptographicShareHandler.CryptographicShareUseCase.AddReceiver(context.Background(), claims.UserID, inputReceiver)
}

func (a *App) AddRecipient(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AddRecipient - error: %v", err)
		return nil, err
	}

	var addRecipRequest share_entry_application_dto.AddRecipientRequest
	if err := json.Unmarshal(raw, &addRecipRequest); err != nil {
		a.Logger.Error("App - AddRecipient - error: %v", err)
		return nil, err
	}
	return a.CryptographicShareHandler.CryptographicShareUseCase.AddRecipient(
		context.Background(),
		claims.UserID,
		addRecipRequest,
		a.AppConfigHandler,
		a.config.ANCHORA_SECRET,
	)
}

func (a *App) UpdateRecipient(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - UpdateRecipient - error: %v", err)
		return nil, err
	}

	var updateRecipRequest share_entry_application_dto.UpdateRecipientRequest
	if err := json.Unmarshal(raw, &updateRecipRequest); err != nil {
		a.Logger.Error("App - UpdateRecipient - error: %v", err)
		return nil, err
	}
	return a.CryptographicShareHandler.UpdateRecipient(context.Background(), claims.UserID, updateRecipRequest)
}

func (a *App) RevokeRecipient(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - RevokeRecipient - error: %v", err)
		return nil, err
	}

	var revokeRecipRequest share_entry_application_dto.UpdateRecipientRequest
	if err := json.Unmarshal(raw, &revokeRecipRequest); err != nil {
		a.Logger.Error("App - RevokeRecipient - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - RevokeRecipient - request: %v", revokeRecipRequest)
	return a.CryptographicShareHandler.RevokeRecipient(context.Background(), claims.UserID, revokeRecipRequest)
}
func (a *App) AcceptShare(jwtToken string, shareID string, intentID string) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - AcceptShare - error: %v", err)
		return nil, err
	}
	return a.CryptographicShareHandler.AcceptShare(context.Background(), shareID, intentID, claims.Email)
}
func (a *App) RejectShare(jwtToken string, shareID string, intentID string) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - RejectShare - error: %v", err)
		return nil, err
	}

	return a.CryptographicShareHandler.RejectShare(context.Background(), shareID, intentID, claims.Email)
}
func (a *App) RevokeShare(jwtToken string, raw json.RawMessage) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - RevokeShare - error: %v", err)
		return nil, err
	}

	var revokeShareRequest share_entry_application_dto.UpdateRecipientRequest
	if err := json.Unmarshal(raw, &revokeShareRequest); err != nil {
		a.Logger.Error("App - RevokeShare - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - RevokeShare - request: %v", revokeShareRequest)
	return a.CryptographicShareHandler.RevokeShare(context.Background(), claims.UserID, revokeShareRequest, a.AppConfigHandler)
}

func (a *App) ListPendingIntentSharesByMe(jwtToken string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListPendingIntentSharesByMe - error: %v", err)
		return nil, err
	}
	res, err := a.CryptographicShareHandler.ListPendingIntentSharesByMe(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - ListPendingIntentSharesByMe - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - ListPendingIntentSharesByMe - res", res)
	return res, nil
}
func (a *App) ListPendingIntentSharesWithMe(jwtToken string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("App - ListPendingIntentSharesWithMe - error: %v", err)
		return nil, err
	}
	res, err := a.CryptographicShareHandler.ListPendingIntentSharesWithMe(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - ListPendingIntentSharesWithMe - error: %v", err)
		return nil, err
	}
	a.Logger.LogPretty("App - ListPendingIntentSharesWithMe - res", res)
	return res, nil
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
	claims, err := a.RequireAuth(input.JwtToken)
	if err != nil {
		a.Logger.Error("❌ GenerateApiKey - Failed to authenticate user: %v", err)
		return nil, err
	}
	userID := claims.UserID
	a.Logger.Info("GenerateApiKey - user id %s:", userID)

	// -------------------------------------------------------------------------------------------------
	// Stellar - Create Stellar account keypair with no friendbot funding
	// -------------------------------------------------------------------------------------------------
	// account, err := a.StellarService.OnGenerateApiKey(input.Password)
	keyEnc := vault_infrastructure_crypto.NewKeyService()
	account, err := keyEnc.CreateAccount("password")
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

	response, err := a.BillingHandler.GetBillingHistoryByUserID(a.ctx, sub.UserID, limit)
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

	sub, err := a.SubscriptionHandler.GetUserSubscriptionByEmail(context.Background(), claims.Email)
	if err != nil {
		a.Logger.Error("App - GetStorageUsage - error: %v", err)
		return nil, err
	}

	return a.SubscriptionHandler.GetStorageUsage(
		a.ctx,
		sub.UserID,
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
		a.Logger.Error("APP - FetchUsers -failed to load all vault users", err)
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
// Websocket integration
// -----------------------------
func (a *App) ConnectToRealtime(user identity_domain.User) error {

	a.Logger.Info("ConnectToRealtime called")
	a.Logger.Info("ConnectToRealtime - user id: %s", user.ID)
	a.ctx = context.WithValue(a.ctx, "user", user)

	ws := realtime_client_infrastructure_websocket.NewClient(a.config.ANKHORA_WEBSOCKET_GATEWAY + "/ws/" + user.ID)

	handlers := realtime_client_handlers.RegisterHandlers(a.ctx)

	client := realtime_client_application_services.NewClient(handlers)

	go func() {
		err := ws.Run(a.ctx, func(msg shared_realtime.Message) error {
			return client.Handle(a.ctx, msg)
		})
		if err != nil {
			utils.LogPretty("APP - ConnectToRealtime - error: %v", err)
		}
	}()

	return nil
}

// -----------------------------
// Notifications Center
// -----------------------------
func (a *App) ListByUser(jwtToken string, limit int, offset int) ([]notification_center_domain.Notification, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	return a.NotificationCenterHandler.ListByUser(context.Background(), claims.UserID, limit, offset)
}
func (a *App) CountUnread(jwtToken string) (int64, error) {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return 0, err
	}

	return a.NotificationCenterHandler.CountUnread(context.Background(), claims.UserID)
}
func (a *App) MarkRead(jwtToken string, id string) error {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}

	return a.NotificationCenterHandler.MarkRead(context.Background(), id)
}
func (a *App) Archive(jwtToken string, id string) error {
	_, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	return a.NotificationCenterHandler.Archive(context.Background(), id)
}
func (a *App) MarkAllRead(jwtToken string) error {
	claims, err := a.RequireAuth(jwtToken)
	if err != nil {
		return err
	}
	return a.NotificationCenterHandler.MarkAllRead(context.Background(), claims.UserID)
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
		StellarNetwork:            os.Getenv("STELLAR_NETWORK"),
		StellarHorizonURL:         os.Getenv("STELLAR_HORIZON_URL"),
		StellarAssetCode:          os.Getenv("STELLAR_ASSET_CODE"),
		StellarAssetIssuer:        os.Getenv("STELLAR_ASSET_ISSUER"),
		IPFSClient:                os.Getenv("IPFS_CLIENT"),
		IPFSGateway:               os.Getenv("IPFS_GATEWAY"),
		IPFSNetwork:               os.Getenv("IPFS_NETWORK"),
		Branch:                    os.Getenv("BRANCH"),
		EncryptionPolicy:          os.Getenv("ENCRYPTION_POLICY"),
		TracecoreURL:              os.Getenv("TRACECORE_URL"),
		TracecoreToken:            os.Getenv("TRACECORE_TOKEN"),
		CloudBackURL:              os.Getenv("CLOUD_BACK_URL"),
		CloudFrontURL:             os.Getenv("CLOUD_FRONT_URL"),
		ANCHORA_SECRET:            os.Getenv("ANCHORA_SECRET"),
		KEYRING_PATH:              os.Getenv("KEYRING_PATH"),
		ANKHORA_WEBSOCKET_GATEWAY: os.Getenv("ANKHORA_WEBSOCKET_GATEWAY"),
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
	claims, err := a.RequireAuth(JwtToken)
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

// C3 Wails App Methods

func (a *App) CreateWorkspace(JwtToken string, vaultId string, name string, description string) (*tracecore_types.Workspace, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if a.WorkspaceHandler == nil {
		return nil, fmt.Errorf("workspace handler is not initialized")
	}
	return a.WorkspaceHandler.CreateWorkspace(a.ctx, claims.UserID, vaultId, name, description)
}

func (a *App) ListWorkspaces(JwtToken string, vaultId string) ([]tracecore_types.Workspace, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	a.Logger.Info("Listing workspaces for user: %v, vaultId: %v", claims.UserID, vaultId)

	// Restore cloud token before making the request
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		a.Logger.Error("Failed to restore cloud token: %v", err)
	}

	log.Printf("[CLOUD-TRACE] E: ListWorkspaces: token_present=%v token_length=%d",
		a.Vault.TracecoreClient.Token != "", len(a.Vault.TracecoreClient.Token))

	if a.WorkspaceHandler == nil {
		return nil, fmt.Errorf("workspace handler is not initialized")
	}

	res, err := a.WorkspaceHandler.ListWorkspaces(a.ctx, vaultId)
	if err != nil {
		a.Logger.Error("App.ListWorkspaces error: %v", err)
		return nil, err
	}
	utils.LogPretty("[Workspace] App.ListWorkspaces final result", res)
	return res, nil
}

func (a *App) CreateChannel(JwtToken string, workspaceID string, title string, templateID string, slots []channel_domain.Slot, assignments []channel_domain.Assignment, properties []channel_domain.ChannelProperty, policy map[string]interface{}, federation string) (*tracecore_types.ChannelDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.CreateChannel(a.ctx, claims.UserID, workspaceID, title, templateID, slots, assignments, properties, policy, federation)
}

func (a *App) ListChannels(JwtToken string, workspaceID string) ([]tracecore_types.ChannelDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.ListChannels(a.ctx, claims.UserID, workspaceID)
}

// GetChannel fetches a single Channel from the authoritative Cloud backend
// (GET /channels/{id}). The returned Channel is the Cloud-persisted aggregate;
// no local channel is ever fabricated.
func (a *App) GetChannel(JwtToken string, channelID string) (*tracecore_types.ChannelDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.GetChannel(a.ctx, claims.UserID, channelID)
}

// UpdateChannel updates an existing Channel through the authoritative Cloud
// backend (PUT /channels/{id}). The Cloud-persisted aggregate is returned; no
// local mutation is performed.
func (a *App) UpdateChannel(JwtToken string, channelID string, title string, slots []channel_domain.Slot, assignments []channel_domain.Assignment, properties []channel_domain.ChannelProperty, policy map[string]interface{}) (*tracecore_types.ChannelDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.UpdateChannel(a.ctx, claims.UserID, channelID, title, slots, assignments, properties, policy)
}

// DeleteChannel deletes a Channel through the authoritative Cloud backend
// (DELETE /channels/{id}). Cloud is the single source of truth for channel
// existence; a 2xx response is success and HTTP >=400 is surfaced verbatim.
func (a *App) DeleteChannel(JwtToken string, channelID string) error {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.DeleteChannel(a.ctx, claims.UserID, channelID)
}

func (a *App) ActivateChannel(JwtToken string, channelID string) (*tracecore_types.ChannelDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.ActivateChannel(a.ctx, claims.UserID, channelID)
}

// RevokeChannel revokes an active Channel through the authoritative Cloud
// backend. The Cloud response carries no Channel data; the frontend refreshes
// the workspace channel list to observe the revoked status.
func (a *App) RevokeChannel(JwtToken string, channelID string) error {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.RevokeChannel(a.ctx, claims.UserID, channelID)
}

// AddParticipant joins an external vault to a Channel through the authoritative
// Cloud backend (POST /channels/{id}/participants). slotID and role are
// optional; when supplied they are forwarded verbatim to the Cloud contract.
func (a *App) AddParticipant(JwtToken string, channelID string, vaultID string, publicKey string, direction string, slotID string, role string) (*tracecore_types.ChannelParticipantDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.AddParticipant(a.ctx, claims.UserID, channelID, vaultID, publicKey, direction, slotID, role)
}

// ListParticipants returns the vaults Cloud has persisted as participants for
// the Channel (GET /channels/{id}/participants). An empty result is valid.
func (a *App) ListParticipants(JwtToken string, channelID string) ([]tracecore_types.ChannelParticipantDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.ListParticipants(a.ctx, claims.UserID, channelID)
}

// InviteToChannel creates a channel invitation through the authoritative Cloud
// backend (POST /channels/{id}/invitations). The invitation carries no slot or
// role information; role semantics are a channel participant concern.
func (a *App) InviteToChannel(JwtToken string, channelID string, inviterVaultID string, inviteeVaultID string) (*tracecore_types.ChannelInvitationDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if err := a.RestoreCloudTokenForUser(claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to restore Cloud token: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.InviteToChannel(a.ctx, claims.UserID, channelID, inviterVaultID, inviteeVaultID)
}

// AcceptChannelInvitation accepts a pending channel invitation through the
// authoritative Cloud backend (POST /channels/invitations/{id}/accept). Cloud
// validates the acceptance and persists the resulting participant; the accept
// response carries the accepted Invitation, not the participant. Cloud is
// idempotent: accepting an already-accepted invitation returns the accepted
// invitation without a duplicate participant.
func (a *App) AcceptChannelInvitation(JwtToken string, invitationID string, inviteeVaultID string, inviteePublicKey string) (*tracecore_types.ChannelInvitationDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if a.ChannelHandler == nil {
		return nil, fmt.Errorf("channel handler is not initialized")
	}
	return a.ChannelHandler.AcceptChannelInvitation(a.ctx, claims.UserID, invitationID, inviteeVaultID, inviteePublicKey)
}

func (a *App) CreateThread(JwtToken string, channelID string, title string, subtitle string, assetType string) (*tracecore_types.ThreadDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if a.ThreadHandler == nil {
		return nil, fmt.Errorf("thread handler is not initialized")
	}
	return a.ThreadHandler.CreateThread(a.ctx, claims.UserID, channelID, title, subtitle, assetType)
}

func (a *App) ListThreads(JwtToken string, channelID string) ([]tracecore_types.ThreadDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if a.ThreadHandler == nil {
		return nil, fmt.Errorf("thread handler is not initialized")
	}
	return a.ThreadHandler.ListThreads(a.ctx, claims.UserID, channelID)
}

func (a *App) ListThreadEvents(JwtToken string, threadID string) ([]tracecore_types.ThreadEventDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if a.ThreadHandler == nil {
		return nil, fmt.Errorf("thread handler is not initialized")
	}
	return a.ThreadHandler.ListThreadEvents(a.ctx, claims.UserID, threadID)
}

func (a *App) AppendThreadEvent(JwtToken string, threadID string, eventType string, payloadJson string) (*tracecore_types.ThreadEventDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	var payload map[string]interface{}
	if payloadJson != "" {
		_ = json.Unmarshal([]byte(payloadJson), &payload)
	}
	if payload == nil {
		payload = make(map[string]interface{})
	}
	if a.ThreadHandler == nil {
		return nil, fmt.Errorf("thread handler is not initialized")
	}
	return a.ThreadHandler.AppendThreadEvent(a.ctx, claims.UserID, threadID, eventType, payload)
}

func (a *App) CreateCollaborativeShare(JwtToken string, threadID string, trustGroupID string, assetCID string, targetVaultID string, notes string) (*tracecore_types.ShareEntryRefDTO, error) {
	claims, err := a.RequireAuth(JwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	if a.CollaborationHandler == nil {
		return nil, fmt.Errorf("collaboration handler is not initialized")
	}
	return a.CollaborationHandler.CreateCollaborativeShare(a.ctx, claims.UserID, threadID, trustGroupID, assetCID, targetVaultID, notes)
}

// ConnectVault explicitly connects the local vault to Ankhora Cloud.
// This is a first-class vault lifecycle operation, NOT a hidden side effect
// of SignIn or ListWorkspaces. The Ledger should not discover a 403 and
// secretly register the vault.
//
// Precondition: User must be signed in with a valid Cloud bearer token.
// This method should be invoked from the Wails UI when the user clicks
// "Connect Vault" — it is a first-class vault lifecycle operation.
//
// The flow is:
//   1. POST /api/identity/challenge (vault_id)
//   2. Sign challenge locally with the vault's stellar signing key
//   3. POST /api/identity/register (challenge_id + signature + vault details)
//   4. Return Connected — /api/workspaces will then return 200 instead of 403.
//
// If the vault is already delegated, the Cloud protocol handles idempotency;
