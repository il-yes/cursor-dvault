package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	utils "vault-app/internal"
	share_application "vault-app/internal/application/use_cases"
	"vault-app/internal/auth"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	share_domain "vault-app/internal/domain/shared"
	"vault-app/internal/driver"
	"vault-app/internal/handlers"
	"vault-app/internal/logger/logger"
	"vault-app/internal/registry"
	shared "vault-app/internal/shared/stellar"
	stellar_recovery_usecase "vault-app/internal/stellar_recovery/application/use_case"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
	stellar "vault-app/internal/stellar_recovery/infrastructure"
	"vault-app/internal/stellar_recovery/infrastructure/events"
	stellar_recovery_persistence "vault-app/internal/stellar_recovery/infrastructure/persistence"
	"vault-app/internal/stellar_recovery/infrastructure/token"
	stellar_recovery_ui_api "vault-app/internal/stellar_recovery/ui/api"
	"vault-app/internal/tracecore"

	// "vault-app/internal/logger/logger"
	"vault-app/internal/models"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
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
}

type App struct {
	config   config
	Logger   logger.Logger
	version  string
	DB       models.DBModel
	ctx      context.Context
	sessions map[int]*models.VaultSession
	NowUTC   func() string

	// Core handlers
	StellarRecoveryHandler    *stellar_recovery_ui_api.StellarRecoveryHandler
	ConnectWithStellarHandler *stellar_recovery_ui_api.StellarRecoveryHandler
	Vaults                    *handlers.VaultHandler
	Auth                      *handlers.AuthHandler
	EntryRegistry             *registry.EntryRegistry

	// New: Global state
	RuntimeContext *models.VaultRuntimeContext
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

	// Use auto-init from env
	appLogger := logger.NewFromEnv()
	appLogger.Info("🚀 Starting D-Vault initialization...")

	// Pick DSN
	dsn := cfg.db.dsn
	if dsn == "" {
		dsn = "sqlite3.db"
	}

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

	// Init services
	appLogger.Info("🔧 Initializing IPFS client...")
	ipfs := blockchain.NewIPFSClient("localhost:5001")
	appLogger.Info("✅ IPFS client initialized (connection will be tested on first use)")

	sessions := make(map[int]*models.VaultSession)

	appLogger.Info("🔧 Initializing Tracecore client...")
	tcClient := tracecore.NewTracecoreClient(os.Getenv("ANKHORA_URL"), os.Getenv("TRACECORE_TOKEN"))
	appLogger.Info("✅ Tracecore client initialized")

	reg := registry.NewRegistry(appLogger)
	reg.RegisterDefinitions([]registry.EntryDefinition{
		{
			Type:    "login",
			Factory: func() models.VaultEntry { return &models.LoginEntry{} },
			Handler: handlers.NewLoginHandler(*db, *ipfs, sessions, appLogger),
		},
		{
			Type:    "card",
			Factory: func() models.VaultEntry { return &models.CardEntry{} },
			Handler: handlers.NewCardHandler(*db, *ipfs, sessions, appLogger),
		},
		{
			Type:    "note",
			Factory: func() models.VaultEntry { return &models.NoteEntry{} },
			Handler: handlers.NewNoteHandler(*db, *ipfs, sessions, appLogger),
		},
		{
			Type:    "identity",
			Factory: func() models.VaultEntry { return &models.IdentityEntry{} },
			Handler: handlers.NewIdentityHandler(*db, *ipfs, sessions, appLogger),
		},
		{
			Type:    "sshkey",
			Factory: func() models.VaultEntry { return &models.SSHKeyEntry{} },
			Handler: handlers.NewSSHKeyHandler(*db, *ipfs, sessions, appLogger),
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	runtimeCtx := &models.VaultRuntimeContext{
		AppSettings: app_config.AppConfig{
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
		LoadedEntries:  []models.VaultEntry{},
	}
	vaults := handlers.NewVaultHandler(*db, ipfs, reg, sessions, appLogger, tcClient, *runtimeCtx)
	auth := handlers.NewAuthHandler(*db, vaults, ipfs, appLogger, tcClient, cfg.auth)

	// Recovery context implementations
	userRepo := stellar_recovery_persistence.NewGormUserRepository(db.DB)
	vaultRepo := stellar_recovery_persistence.NewGormVaultRepository(db.DB)
	subRepo := stellar_recovery_persistence.NewGormSubscriptionRepository(db.DB)
	verifier := stellar.NewStellarKeyAdapter()
	eventDisp := events.NewLocalDispatcher()
	tokenGen := token.NewSimpleTokenGen()

	checkUC := stellar_recovery_usecase.NewCheckKeyUseCase(userRepo, vaultRepo, subRepo, verifier)
	recoverUC := stellar_recovery_usecase.NewRecoverVaultUseCase(userRepo, vaultRepo, subRepo, verifier, eventDisp, tokenGen)
	importUC := stellar_recovery_usecase.NewImportKeyUseCase(userRepo, verifier)
	loginAdapter := shared.NewStellarLoginAdapter(db)

	connectUC := stellar_recovery_usecase.NewConnectWithStellarUseCase(
		loginAdapter,
		vaultRepo,
		subRepo,
	)

	stellarHandler := stellar_recovery_ui_api.NewStellarRecoveryHandler(checkUC, recoverUC, importUC, connectUC)

	// ⚡ Restore sessions asynchronously to speed up startup
	go func() {
		appLogger.Info("🔄 Restoring sessions in background...")
		storedSessions, err := db.GetAllSessions()
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

	// Start pending commit worker
	vaults.StartPendingCommitWorker(ctx, 2*time.Minute)

	elapsed := time.Since(startTime)
	appLogger.Info("✅ D-Vault initialized successfully in %v", elapsed)

	return &App{
		config:                    cfg,
		version:                   version,
		DB:                        *db,
		StellarRecoveryHandler:    stellarHandler,
		Vaults:                    vaults,
		Auth:                      auth,
		Logger:                    *appLogger,
		EntryRegistry:             reg,
		sessions:                  sessions,
		RuntimeContext:            runtimeCtx,
		cancel:                    cancel,
		NowUTC:                    func() string { return time.Now().Format(time.RFC3339) },
		ConnectWithStellarHandler: stellarHandler,
	}
}

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

// -----------------------------
// Connexion
// -----------------------------
func (a *App) SignIn(req handlers.LoginRequest) (*handlers.LoginResponse, error) {
	return a.Auth.Login(req)
}
func (a *App) SignInWithStellar(req handlers.LoginRequest) (*handlers.LoginResponse, error) {
	return a.Auth.Login(req)
}

func (a *App) SignUp(setup handlers.OnBoarding) (*handlers.OnBoardingResponse, error) {
	return a.Auth.OnBoarding(setup)
}
func (a *App) SignOut(userID int, cid string, password string) {
	a.Vaults.EndSession(userID)
	a.Auth.Logout(userID)
}
func (a *App) CheckSession(userID int) (string, error) {
	return a.Auth.RefreshToken(userID) // same logic you already wrote
}
func (a *App) CheckEmail(email string) (*handlers.CheckEmailResponse, error) {
	return a.Auth.CheckEmail(email)
}

// -----------------------------
// JWT Token
// -----------------------------
func (a *App) RefreshToken(userID int) (string, error) {
	token, err := a.Auth.RefreshToken(userID)
	if err != nil {
		return "", err
	}
	return token, nil
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
func (a *App) AddEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vaults.AddEntry(claims.UserID, entryType, raw)
}

func (a *App) EditEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vaults.UpdateEntry(claims.UserID, entryType, raw)
}

func (a *App) TrashEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vaults.TrashEntry(claims.UserID, entryType, raw)
}

func (a *App) RestoreEntry(entryType string, raw json.RawMessage, jwtToken string) (any, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vaults.RestoreEntry(claims.UserID, entryType, raw)
}

func (a *App) CreateFolder(name string, jwtToken string) (*models.VaultPayload, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vaults.CreateFolder(claims.UserID, name)
}
func (a *App) GetFoldersByVault(vaultCID string, jwtToken string) ([]models.Folder, error) {
	_, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vaults.GetFoldersByVault(vaultCID)
}
func (a *App) UpdateFolder(id int, newName string, isDraft bool, jwtToken string) (*models.Folder, error) {
	_, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}
	return a.Vaults.UpdateFolder(id, newName, isDraft)
}
func (a *App) DeleteFolder(userID int, id int, jwtToken string) (string, error) {
	_, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Folder deleted %d successfuly", id), a.Vaults.DeleteFolder(userID, id)
}

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
	userID := int(claims.UserID)
	return a.Vaults.CreateShareEntry(context.Background(), input.Payload, userID)
}

func (a *App) ListSharedEntries(jwtToken string) (*[]share_domain.ShareEntry, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, fmt.Errorf("ListSharedEntries - auth failed: %w", err)
	}
	utils.LogPretty("ListSharedEntries - claims", claims)

	entries, err := a.Vaults.ListSharedEntries(context.Background(), int(claims.UserID))
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

	entries, err := a.Vaults.ListReceivedShares(context.Background(), int(claims.UserID))
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
		int(claims.UserID),
		shareID,
	)
}
func (a *App) RejectShare(jwtToken string, shareID uint) (*share_application.RejectShareResult, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	return a.Vaults.RejectShare(context.Background(), int(claims.UserID), shareID)
}
func (a *App) AddReceiver(jwtToken string, payload share_application.AddReceiverInput) (*share_application.AddReceiverResult, error) {
	claims, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		a.Logger.Error("❌ Failed to authenticate user: %v", err)
		return nil, err
	}

	return a.Vaults.AddReceiver(context.Background(), int(claims.UserID), payload)
}

// FlushAllSessions persists and clears all active sessions.
func (a *App) FlushAllSessions() {
	a.Vaults.SessionsMu.Lock()
	defer a.Vaults.SessionsMu.Unlock()

	if len(a.Vaults.Sessions) == 0 {
		a.Logger.Info("No sessions to flush")
		return
	}

	a.Logger.Info("💾 Flushing %d active sessions...", len(a.Vaults.Sessions))

	for userID, session := range a.Vaults.Sessions {
		if err := a.DB.SaveSession(userID, session); err != nil {
			a.Logger.Error("❌ Failed to flush session for user %d: %v", userID, err)
		} else {
			a.Logger.Info("✅ Session flushed for user %d", userID)
		}
		delete(a.Vaults.Sessions, userID)
	}

	a.Logger.Info("✨ All sessions flushed and cleared")
}

func (a *App) IsVaultDirty(jwtToken string) (bool, error) {
	_, err := a.Auth.RequireAuth(jwtToken)
	if err != nil {
		return false, err
	}
	return a.Vaults.IsVaultDirty(), nil
}

func (a *App) FetchUsers() ([]models.UserDTO, error) {
	return a.DB.FindUsers()
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
	}
}
func (a *App) startup(ctx context.Context) {
	fmt.Println("App has started.")
	a.ctx = ctx
}
