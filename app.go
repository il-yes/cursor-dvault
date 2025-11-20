package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"vault-app/internal/auth"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	"vault-app/internal/driver"
	"vault-app/internal/handlers"
	"vault-app/internal/logger/logger"
	"vault-app/internal/registry"
	"vault-app/internal/tracecore"

	// "vault-app/internal/logger/logger"
	"vault-app/internal/models"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
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

	// Core handlers
	Vaults        *handlers.VaultHandler
	Auth          *handlers.AuthHandler
	EntryRegistry *registry.EntryRegistry

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

	ctx, cancel := context.WithCancel(context.Background())

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
		config:         cfg,
		version:        version,
		DB:             *db,
		Vaults:         vaults,
		Auth:           auth,
		Logger:         *appLogger,
		EntryRegistry:  reg,
		sessions:       sessions,
		RuntimeContext: runtimeCtx,
		cancel:         cancel,
	}
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
	return a.Vaults.SyncVault(claims.UserID, password)
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
