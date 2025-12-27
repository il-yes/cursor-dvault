package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	utils "vault-app/internal"
	"vault-app/internal/auth"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	onboarding_domain "vault-app/internal/onboarding/domain"
	"vault-app/internal/registry"
	"vault-app/internal/tracecore"

	// "os"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB                       models.DBModel
	Vaults                   *VaultHandler
	IPFS                     *blockchain.IPFSClient
	NowUTC                   func() string
	logger                   logger.Logger
	TracecoreClient          *tracecore.TracecoreClient
	auth                     auth.Auth
	UserOnboardingRepository onboarding_domain.UserRepository
}

func NewAuthHandler(db models.DBModel, vaults *VaultHandler, ipfs *blockchain.IPFSClient, logger *logger.Logger, tc *tracecore.TracecoreClient,
	auth auth.Auth, userOnboardingRepository onboarding_domain.UserRepository) *AuthHandler {
	return &AuthHandler{
		DB:                       db,
		Vaults:                   vaults,
		IPFS:                     ipfs,
		NowUTC:                   func() string { return time.Now().Format(time.RFC3339) },
		logger:                   *logger,
		TracecoreClient:          tc,
		auth:                     auth,
		UserOnboardingRepository: userOnboardingRepository,
	}
}

// -----------------------------
// Sign In
// -----------------------------
type LoginRequest struct {
	Email         string `json:"email,omitempty"`
	Password      string `json:"password,omitempty"`
	PublicKey     string `json:"publicKey,omitempty"`     // optional
	PrivateKey    string `json:"privateKey,omitempty"`    // optional
	SignedMessage string `json:"signedMessage,omitempty"` // optional
	Signature     string `json:"signature,omitempty"`     // optional
}

type LoginResponse struct {
	User                models.User                 `json:"User"`
	Vault               *models.VaultPayload         `json:"Vault"`
	Tokens              *auth.TokenPairs            `json:"Tokens"`
	CloudToken          string                      `json:"cloud_token"`
	VaultRuntimeContext *models.VaultRuntimeContext `json:"vault_runtime_context"`
	LastCID             string                      `json:"last_cid"`
	Dirty               bool                        `json:"dirty"`
	SessionID           string                      `json:"session_id"`
	// SharedEntries      []share_domain.ShareEntry `json:"share_domain_entries"`
}

func (ah *AuthHandler) Login(credentials LoginRequest) (*LoginResponse, error) {
	var user *models.User
	var userOnboarding *onboarding_domain.User
	var err error
	var cloudLoginResponse *tracecore.LoginResponse
	fmt.Println("🔥 SIGNIN ROUTE HIT AT:", time.Now())

	// -----------------------------
	// 1. Identify login method
	// -----------------------------
	if credentials.PublicKey != "" && credentials.SignedMessage != "" {
		ah.logger.Info("🔑 Stellar login request: %s", credentials.SignedMessage)

		user, _, err = ah.DB.GetUserByPublicKey(credentials.PublicKey)
		if err != nil || user == nil {
			return nil, fmt.Errorf("❌ user not found for public key %s: %w", credentials.PublicKey, err)
		}

		if !blockchain.VerifySignature(credentials.PublicKey, credentials.SignedMessage, credentials.Signature) {
			return nil, fmt.Errorf("❌ stellar signature verification failed: %w", err)
		}

		// 🔓 Recover the plain password from stored encrypted data
		userCfg, err := ah.DB.GetUserConfigByUserID(user.ID)
		if err != nil {
			return nil, fmt.Errorf("❌ failed to load user config: %w", err)
		}
		plainPassword, err := blockchain.DecryptPasswordWithStellar(
			userCfg.StellarAccount.EncNonce,
			userCfg.StellarAccount.EncPassword,
			userCfg.StellarAccount.PrivateKey, // frontend must provide this
		)
		if err != nil {
			return nil, fmt.Errorf("❌ failed to recover password: %w", err)
		}
		println("plainPassword", plainPassword)
		credentials.Password = plainPassword

	} else {
		ah.logger.Info("📧 Email login request: %s", credentials.Email)
		// authenticate to Ankhora.io
		// cloudLoginResponse, err = ah.TracecoreClient.Login(context.Background(), tracecore.LoginRequest{
		// 	Email:         credentials.Email,
		// 	Password:      credentials.Password,
		// 	PublicKey:     credentials.PublicKey,
		// 	SignedMessage: credentials.SignedMessage,
		// 	Signature:     credentials.Signature,
		// })
		// if err != nil {
		// 	return nil, fmt.Errorf("❌ failed to authenticate with Ankhora: %w", err)
		// }
		// utils.LogPretty("cloud login respoonse", cloudLoginResponse)

		userOnboarding, err = ah.UserOnboardingRepository.FindByEmail(credentials.Email)
		if err != nil || userOnboarding == nil {
			if err != nil {
				ah.logger.Error("Monitor - Error retrieving user onboarding: %v", err)
			} else {
				ah.logger.Error("Monitor - User onboarding is nil for email: %s", credentials.Email)
			}
			return nil, fmt.Errorf("❌ User not found: %w", err)
		}

		// Log the retrieved user onboarding
		utils.LogPretty("userOnboarding", userOnboarding)

		// Defensive check before bcrypt
		if userOnboarding.Password == "" {
			ah.logger.Error("Monitor - Retrieved user has empty password: %v", userOnboarding)
			return nil, fmt.Errorf("❌ empty password for user: %s", credentials.Email)
		}

		// Compare hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(userOnboarding.Password), []byte(credentials.Password)); err != nil {
			ah.logger.Error("Monitor - Invalid credentials for user: %s", credentials.Email)
			return nil, fmt.Errorf("❌ invalid credentials: %w", err)
		}

		ah.logger.Info("✅ Password verified for user: %s", credentials.Email)

	}
	ah.logger.Info("✅ User found: %s", userOnboarding.Email)

	// -----------------------------
	// 2. Always update last connection
	// -----------------------------
	user, errUser := ah.DB.GetUserByEmail(userOnboarding.Email)
	if errUser != nil {
		return nil, fmt.Errorf("❌ User not found: %w", errUser)
	}
	utils.LogPretty("connexion - user", user)
	user.LastConnectedAt = time.Now().UTC()
	ah.logger.Info("connexion - user last connected at", user.LastConnectedAt)

	// 9. create a jwt user
	u := auth.JwtUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	utils.LogPretty("Auth user", u)

	// generate tokens
	tokens, err := ah.auth.GenerateTokenPair(&u)
	ah.logger.Info("Auth tokens", tokens)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to generate token for user %d - %w", user.ID, err)
	}
	utils.LogPretty("Auth tokens", tokens)
	//  Save tokens to DB (for persistence across restarts)
	savedtoken, err := ah.DB.SaveJwtToken(tokens)
	fmt.Println("savedtoken", savedtoken)
	if err != nil {
		ah.logger.Error("❌ failed to persist tokens: %v", err)
		return nil, fmt.Errorf("❌ failed to persist tokens: %w", err)
	}
	ah.logger.Info("saved token: ", savedtoken.Token)

	// I. LOAD VAULT
	// -----------------------------
	// 3. Try to reuse existing session
	// -----------------------------
	if existingSession, ok := ah.Vaults.Sessions[user.ID]; ok {
		fmt.Println("existingSession", existingSession)
		if existingSession.Dirty {
			ah.Vaults.MarkDirty(user.ID)
		}

		ah.logger.Info("♻️ Reusing in-memory session for user %s", user.ID)

		existingSession.VaultRuntimeContext.SessionSecrets["dvault_jwt"] = tokens.Token

		if cloudLoginResponse != nil && cloudLoginResponse.AuthenticationToken.Token != "" {
			existingSession.VaultRuntimeContext.SessionSecrets["cloud_jwt"] = cloudLoginResponse.AuthenticationToken.Token
			ah.TracecoreClient.Token = cloudLoginResponse.Token
			ah.logger.Info("cloud token set to: ", ah.TracecoreClient.Token)
		}

		return &LoginResponse{
			User:                *user,
			Vault:               existingSession.Vault,
			Tokens:              &tokens,
			CloudToken:          cloudLoginResponse.AuthenticationToken.Token,
			VaultRuntimeContext: &existingSession.VaultRuntimeContext,
			LastCID:             existingSession.LastCID,
			Dirty:               existingSession.Dirty,
		}, nil
	}
	ah.logger.Info("No existing session found for user %s", user.ID)
	// -----------------------------
	// 4. Try to load from DB
	// -----------------------------
	storedSession, err := ah.DB.LoadSession(user.ID)
	if err == nil && storedSession != nil {
		storedSession = RehydrateSession(storedSession)
		ah.Vaults.Sessions[user.ID] = storedSession

		// re-queue any pending commits
		if len(storedSession.PendingCommits) > 0 {
			for _, commit := range storedSession.PendingCommits {
				if err := ah.Vaults.QueuePendingCommits(user.ID, commit); err != nil {
					ah.logger.Error("❌ Failed to queue commit for user %s: %v", user.ID, err)
				}
			}
		}

		if storedSession.Dirty {
			ah.Vaults.MarkDirty(user.ID)
		}

		storedSession.VaultRuntimeContext.SessionSecrets["dvault_jwt"] = tokens.Token

		if cloudLoginResponse != nil && cloudLoginResponse.AuthenticationToken.Token != "" {
			storedSession.VaultRuntimeContext.SessionSecrets["cloud_jwt"] = cloudLoginResponse.AuthenticationToken.Token
			ah.TracecoreClient.Token = cloudLoginResponse.Token
			ah.logger.Info("cloud token set to: ", ah.TracecoreClient.Token)
		}
		ah.logger.Info("🔄 Restored session for user %s from DB", user.ID)

		return &LoginResponse{
			User:                *user,
			Vault:               storedSession.Vault,
			Tokens:              &tokens,
			CloudToken:          cloudLoginResponse.AuthenticationToken.Token,
			VaultRuntimeContext: &storedSession.VaultRuntimeContext,
			LastCID:             storedSession.LastCID,
			Dirty:               storedSession.Dirty,
		}, nil
	}
	ah.logger.Info("after No stored session found for user %s", user.ID)

	// -----------------------------
	// 5. Fresh login → fetch vault
	// -----------------------------
	vaultMeta, err := ah.DB.GetLatestVaultCIDByUserID(user.ID)
	// If the user exists but has NO VAULT → minimal onboarding
	if errors.Is(err, gorm.ErrRecordNotFound) || vaultMeta == nil {
		ah.logger.Warn("⚠️ User %s has no vault, performing minimal onboarding...", user.ID)
		vaultPayload, _ := ah.OnboardMissingVault(user, credentials.Password)
		utils.LogPretty("✅ Minimal vault created", vaultPayload)
	}

	rawVault, err := ah.IPFS.GetData(vaultMeta.CID)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to fetch vault from IPFS: %w", err)
	}
	if rawVault == nil || len(rawVault) == 0 {
		return nil, fmt.Errorf("❌ empty vault data for CID %s", vaultMeta.CID)
	}

	decrypted, err := blockchain.Decrypt(rawVault, credentials.Password)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to decrypt vault: %w", err)
	}
	if len(decrypted) == 0 {
		return nil, fmt.Errorf("❌ vault decryption returned empty result")
	}

	vaultPayload := ParseVaultPayload(decrypted)

	// -----------------------------
	// 6. Load App & User Config
	// -----------------------------
	appCfg, _ := ah.DB.GetAppConfigByUserID(user.ID)
	userCfg, _ := ah.DB.GetUserConfigByUserID(user.ID)

	// If either config missing → minimal config onboarding
	if appCfg == nil || userCfg == nil {
		ah.logger.Warn("⚠️ Missing configs for user %s — creating minimal config...", user.ID)
		appCfg, userCfg, err = ah.OnboardMissingConfig(user, credentials.Password)
		if err != nil {
			return nil, fmt.Errorf("❌ failed to onboard missing config: %w", err)
		}
	}

	// -----------------------------
	// 7. Create runtime context
	// -----------------------------
	runtimeCtx := &models.VaultRuntimeContext{
		AppSettings:    *appCfg,
		CurrentUser:    *userCfg,
		SessionSecrets: make(map[string]string),
		WorkingBranch:  "main",
		LoadedEntries:  []models.VaultEntry{},
	}

	// -----------------------------
	// 8. Start new session
	// -----------------------------
	ah.Vaults.StartSession(user.ID, vaultPayload, "main", runtimeCtx)
	ah.logger.Info("✅ Vault session started for user %s", user.ID)

	// -----------------------------
	// 10. Return login response
	// -----------------------------
	return &LoginResponse{
		User:                *user,
		Vault:               &vaultPayload,
		Tokens:              &tokens,
		CloudToken:          cloudLoginResponse.Token,
		VaultRuntimeContext: runtimeCtx,
		LastCID:             vaultMeta.CID,
		Dirty:               false,
	}, nil
}

func (ah *AuthHandler) OnboardMissingVault(user *models.User, password string) (*models.VaultPayload, error) {
	vaultName := fmt.Sprintf("%s-vault", user.Username)

	// Create empty vault
	vaultPayload := models.VaultPayload{
		Version: "1.0.0",
		Name:    vaultName,
		BaseVaultContent: models.BaseVaultContent{
			Folders:   []models.Folder{},
			Entries:   models.Entries{},
			CreatedAt: ah.NowUTC(),
			UpdatedAt: ah.NowUTC(),
		},
	}
	utils.LogPretty("✅ OnboardMissingVault - vaultPayload", vaultPayload)

	vaultBytes, _ := json.MarshalIndent(vaultPayload, "", "  ")
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return nil, fmt.Errorf("❌ minimal vault encryption failed: %w", err)
	}

	cid, err := ah.IPFS.AddData(encrypted)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to save minimal vault: %w", err)
	}
	savedVault, err := ah.DB.SaveVaultCID(models.VaultCID{
		Name:      vaultName,
		Type:      "vault",
		UserID:    user.ID,
		CID:       cid,
		CreatedAt: ah.NowUTC(),
		UpdatedAt: ah.NowUTC(),
	})
	if err != nil {
		return nil, fmt.Errorf("❌ failed to persist minimal vault metadata: %w", err)
	}

	ah.logger.Info("✅ Minimal vault created for user %s, CID=%s", user.ID, savedVault.CID)

	return &vaultPayload, nil

	// return &LoginResponse{
	// 	User:                *user,
	// 	Vault:               vaultPayload,
	// 	Tokens:              nil, // or generate JWT here if needed
	// 	CloudToken:          "",
	// 	VaultRuntimeContext: runtimeCtx,
	// 	LastCID:             savedVault.CID,
	// 	Dirty:               false,
	// }, nil
}
func (ah *AuthHandler) OnboardMissingConfig(user *models.User, password string) (*app_config.AppConfig, *app_config.UserConfig, error) {
	ah.logger.Info("🛠 Creating minimal config for user %d", user.ID)

	// ---- AppConfig ----
	appCfg := app_config.AppConfig{
		UserID:           user.ID,
		Branch:           "main",
		TracecoreEnabled: false,
		EncryptionPolicy: "AES-256-GCM",
		VaultSettings: app_config.VaultConfig{
			MaxEntries:       1000,
			AutoSyncEnabled:  false,
			EncryptionScheme: "AES-256-GCM",
		},
		Blockchain: app_config.BlockchainConfig{
			Stellar: app_config.StellarConfig{
				Network:    "testnet",
				HorizonURL: "https://horizon-testnet.stellar.org",
				Fee:        100,
			},
		},
	}

	// ---- UserConfig ----
	userCfg := app_config.UserConfig{
		ID:             user.ID,
		Role:           "user",
		Signature:      "",
		SharingRules:   []app_config.SharingRule{},
		StellarAccount: app_config.StellarAccountConfig{}, // optional
	}

	if err := ah.SaveConfigurations(appCfg, userCfg); err != nil {
		return nil, nil, fmt.Errorf("❌ failed to save minimal configs: %w", err)
	}

	ah.logger.Info("✅ Minimal configs created for user %s", user.ID)

	return &appCfg, &userCfg, nil
}

func RehydrateSession(s *models.VaultSession) *models.VaultSession {
	if s.VaultRuntimeContext.SessionSecrets == nil {
		s.VaultRuntimeContext.SessionSecrets = make(map[string]string)
	}
	if s.VaultRuntimeContext.WorkingBranch == "" {
		s.VaultRuntimeContext.WorkingBranch = "main"
	}
	if s.VaultRuntimeContext.LoadedEntries == nil {
		s.VaultRuntimeContext.LoadedEntries = []models.VaultEntry{}
	}
	return s
}

func ParseVaultPayload(decrypted []byte) models.VaultPayload {
	var vault models.VaultPayload

	err := json.Unmarshal(decrypted, &vault)
	if err != nil {
		fmt.Println("❌ Failed to parse vault JSON:", err)
		// Fallback: return empty vault
		vault = models.VaultPayload{}
	}
	return vault
}

func (ah *AuthHandler) RequestChallenge(req blockchain.ChallengeRequest) (blockchain.ChallengeResponse, error) {
	challenge := blockchain.GenerateChallenge(req.PublicKey)

	response := blockchain.ChallengeResponse{
		Challenge: challenge,
		ExpiresAt: time.Now().Add(5 * time.Minute).Format(time.RFC3339),
	}

	return response, nil
}
func (ah *AuthHandler) AuthVerify(req *blockchain.SignatureVerification) (string, error) {
	expectedChallenge, err := blockchain.ChallengeStore[req.PublicKey]
	if !err || expectedChallenge != req.Challenge {
		return "", fmt.Errorf("❌ Invalid or expired challenge: %w", err)
	}
	if blockchain.VerifySignature(req.PublicKey, req.Challenge, req.Signature) {
		return fmt.Sprintf("✅ Signature verified. Login successful."), nil
	} else {
		return "", fmt.Errorf("❌ Signature verification failed: %w", err)
	}
}

type CheckEmailResponse struct {
	Status      string   `json:"status"`
	AuthMethods []string `json:"auth_methods,omitempty"`
}

func (ah *AuthHandler) CheckEmail(email string) (*CheckEmailResponse, error) {

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

// -----------------------------
// Sign Up
// -----------------------------
type OnBoardingResponse struct {
	Vault  models.VaultPayload
	User   models.User
	Tokens auth.TokenPairs
}

// type OnBoarding struct {
// 	UserID    string `json:"user_id"`    // e.g. "alice@example.com"
// 	UserAlias string `json:"user_alias"` // e.g. "Alice"
// 	Role      string `json:"role"`       // e.g. "user"

// 	EnableTracecore bool   `json:"enable_tracecore"` // optional override
// 	RepoTemplate    string `json:"repo_template"`    // e.g. "contract"

// 	VaultPath        string `json:"vault_path"`        // e.g. "/Users/alice/.dvault"
// 	EncryptionPolicy string `json:"encryption_policy"` // e.g. "AES-256-GCM"

// 	BranchingModel string   `json:"branching_model"` // optional override
// 	CommitRules    []string `json:"commit_rules"`    // optional override
// 	Actors         []string `json:"actors"`          // optional override

// 	FederatedProviders []string                 `json:"federated_providers"` // e.g. ["keycloak"]
// 	SharingDefaults    []app_config.SharingRule `json:"sharing_defaults"`    // optional override

// 	AppVersion string `json:"app_version"` // e.g. "1.0.0"`

// 	// Required internally
// 	Password  string `json:"password"`   // not exposed in payload; comes from secret input
// 	VaultName string `json:"vault_name"` // internal label — derive from alias if needed
// }

// frontend sends only the minimum
type OnBoarding struct {
	UserID    string `json:"user_id"`    // required (email)
	UserAlias string `json:"user_alias"` // required (display name)
	Password  string `json:"password"`   // required (secret)
	VaultName string `json:"vault_name"` // required (derive from alias if missing)

	// optional (defaults if not provided)
	Role               string   `json:"role"`
	RepoTemplate       string   `json:"repo_template"`
	EncryptionPolicy   string   `json:"encryption_policy"`
	FederatedProviders []string `json:"federated_providers"`
}

func (ah *AuthHandler) OnBoarding(setup OnBoarding) (*OnBoardingResponse, error) {
	// -----------------------------
	// 1. Validate & defaults
	// -----------------------------
	if setup.UserID == "" || setup.UserAlias == "" || setup.Password == "" {
		return nil, fmt.Errorf("❌ missing required fields (email, alias, password)")
	}
	if setup.VaultName == "" {
		setup.VaultName = setup.UserAlias + "-vault"
	}
	if setup.Role == "" {
		setup.Role = "user"
	}
	if setup.EncryptionPolicy == "" {
		setup.EncryptionPolicy = "AES-256-GCM"
	}

	// -----------------------------
	// 2. Hash password & create user
	// -----------------------------
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(setup.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to hash password: %w", err)
	}

	user := &models.User{
		Email:    setup.UserID,
		Username: setup.UserAlias,
		Password: string(hashedPassword),
		Role:     setup.Role,
	}
	user, err = ah.DB.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to create user: %w", err)
	}

	// -----------------------------
	// 3. Init empty vault
	// -----------------------------
	vaultPayload := models.VaultPayload{
		Version: "1.0.0",
		Name:    setup.VaultName,
		BaseVaultContent: models.BaseVaultContent{
			Folders:   []models.Folder{},
			Entries:   models.Entries{},
			CreatedAt: ah.NowUTC(),
			UpdatedAt: ah.NowUTC(),
		},
	}

	vaultBytes, _ := json.MarshalIndent(vaultPayload, "", "  ")
	encrypted, err := blockchain.Encrypt(vaultBytes, setup.Password)
	if err != nil {
		return nil, fmt.Errorf("❌ vault encryption failed: %w", err)
	}

	cid, err := ah.IPFS.AddData(encrypted)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to add vault to IPFS: %w", err)
	}
	vaultMeta := models.VaultCID{
		Name:      setup.VaultName,
		Type:      "vault",
		UserID:    user.ID,
		CID:       cid,
		CreatedAt: ah.NowUTC(),
		UpdatedAt: ah.NowUTC(),
	}
	savedVault, err := ah.DB.SaveVaultCID(vaultMeta)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to persist vault metadata: %w", err)
	}
	ah.logger.Info("Saved vault in db - CID: %s", savedVault.CID)

	// -----------------------------
	// 4. (Optional) Bootstrap Tracecore Repo & Load Template
	// -----------------------------
	// Use default template if none specified
	if setup.RepoTemplate == "" {
		setup.RepoTemplate = "personal" // default template
	}

	template, ok := registry.VaultTemplates[setup.RepoTemplate]
	if !ok {
		// Fallback to personal template if invalid template specified
		ah.logger.Warn("⚠️ Invalid template '%s', using 'personal' as fallback", setup.RepoTemplate)
		template = registry.VaultTemplates["personal"]
	}

	var repoIDStr string
	if template.Tracecore.Enabled && ah.TracecoreClient != nil {
		repoID, err := template.ShouldBootstrapTracecoreRepo(ah.TracecoreClient)
		if err == nil && repoID != nil {
			repoIDStr = *repoID
		}
	}
	ah.logger.Info("Repo ID: %s", repoIDStr)

	// -----------------------------
	// 5. (Optional) Stellar account
	// -----------------------------

	var stellarAccount *app_config.StellarAccountConfig
	if template.StellarAccount.Enabled {
		pub, secret, txID, err := blockchain.CreateStellarAccount()
		if err != nil {
			ah.logger.Warn("⚠️ Stellar account creation failed: %v", err)
		} else {
			nonce, ct, _ := blockchain.EncryptPasswordWithStellar(setup.Password, secret)
			stellarAccount = &app_config.StellarAccountConfig{
				PublicKey:   pub,
				PrivateKey:  secret, // TODO: encrypt before saving
				EncNonce:    nonce,
				EncPassword: ct,
			}
			ah.logger.Info("✅ Stellar account created: %s -  tx:", stellarAccount.PublicKey, txID)
		}
	}
	/*
		var stellarAccount *app_config.StellarAccountConfig
		if template.StellarAccount.Enabled {
			pub, secret, txID, err := blockchain.CreateStellarAccount()
			if err != nil {
				ah.logger.Warn("⚠️ Stellar account creation failed: %v", err)
			} else {
				// Encrypt the user password with Stellar secret
				salt, nonce, ct, err := blockchain.EncryptPasswordWithStellarSecure(setup.Password, secret)
				if err != nil {
					ah.logger.Warn("⚠️ Failed to encrypt password with Stellar secret: %v", err)
				}

				// TODO: encrypt the Stellar private key before storing (server-side master key or KMS)
				anchorSecret := os.Getenv("ANCHOR_SECRET")
				if anchorSecret == "" {
					ah.logger.Warn("⚠️ Anchor secret not found")
					return nil, errors.New("anchor secret not found")
				}
				encryptedSecret, err := blockchain.Encrypt([]byte(secret), anchorSecret)
				if err != nil {
					ah.logger.Warn("⚠️ Failed to encrypt Stellar private key: %v", err)
				}

				stellarAccount = &app_config.StellarAccountConfig{
					PublicKey:   pub,
					PrivateKey:  string(encryptedSecret),
					EncSalt:     salt,
					EncNonce:    nonce,
					EncPassword: ct,
				}
				ah.logger.Info("✅ Stellar account created: %s - txID: %s", stellarAccount.PublicKey, txID)
			}
		}
	*/

	// -----------------------------
	// 6. JWT token generation
	// -----------------------------
	u := auth.JwtUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	tokens, err := ah.auth.GenerateTokenPair(&u)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to generate token for user %d: %w", user.ID, err)
	}

	_, err = ah.DB.SaveJwtToken(tokens) // fixed arg
	if err != nil {
		return nil, fmt.Errorf("❌ failed to persist tokens: %w", err)
	}

	// 7. Build AppConfig
	appCfg := app_config.AppConfig{
		UserID:           user.ID,
		RepoID:           repoIDStr,
		Branch:           "main",
		TracecoreEnabled: template.Tracecore.Enabled,
		EncryptionPolicy: setup.EncryptionPolicy,
		VaultSettings: app_config.VaultConfig{
			MaxEntries:       1000,
			AutoSyncEnabled:  false,
			EncryptionScheme: setup.EncryptionPolicy,
		},
		Blockchain: app_config.BlockchainConfig{
			Stellar: app_config.StellarConfig{
				Network:    "testnet",
				HorizonURL: "https://horizon-testnet.stellar.org",
				Fee:        100,
			},
		},
	}

	// 8. Build UserConfig
	var sharingRules []app_config.SharingRule
	userCfg := app_config.UserConfig{
		ID:           user.ID,
		Role:         "user",
		Signature:    "", // will be generated later
		SharingRules: sharingRules,
	}

	// Only set StellarAccount if it was created
	if stellarAccount != nil {
		userCfg.StellarAccount = *stellarAccount
	}
	// 9. Save appconfigurations
	errCfg := ah.SaveConfigurations(appCfg, userCfg)
	if errCfg != nil {
		return nil, fmt.Errorf("❌ Failed to save app & user config: %w", errCfg)
	}

	// -----------------------------
	// 10. Response (lightweight)
	// -----------------------------
	res := &OnBoardingResponse{
		User:   *user,
		Vault:  vaultPayload,
		Tokens: tokens,
	}
	return res, nil
}

func convertSharingRules(defaults []registry.SharingRule) []app_config.SharingRule {
	var rules []app_config.SharingRule
	for _, r := range defaults {
		rules = append(rules, app_config.SharingRule{
			EntryType: r.Role,
			Targets:   []string{r.Group},
			Encrypted: true,
		})
	}
	return rules
}
func convertSharingRules2(defaults []app_config.SharingRule) []app_config.SharingRule {
	var rules []app_config.SharingRule
	for _, r := range defaults {
		rules = append(rules, app_config.SharingRule{
			EntryType: r.EntryType,
			Targets:   r.Targets,
			Encrypted: r.Encrypted,
		})
	}
	return rules
}

func (ah *AuthHandler) SaveConfigurations(appCfg app_config.AppConfig, userCfg app_config.UserConfig) error {
	savedAppCfg, err := ah.DB.SaveAppConfig(appCfg)
	if err != nil {
		ah.logger.Error("❌ Failed to save App Config saved: %d", savedAppCfg.ID)
		return err
	}
	savedUserCfg, err := ah.DB.SaveUserConfig(userCfg)
	if err != nil {
		ah.logger.Error("❌ Failed to save User Config: %d", savedUserCfg.ID)
		return err
	}
	ah.logger.Info("App Config saved: %s & User Config saved: %s", savedAppCfg.ID, savedUserCfg.ID)
	return nil
}
func (ah *AuthHandler) LinkStellarKey(user *app_config.UserConfig, stellarSecret string, plainPassword string) error {
	nonce, ct, err := blockchain.EncryptPasswordWithStellar(plainPassword, stellarSecret)
	if err != nil {
		return err
	}
	user.StellarAccount = app_config.StellarAccountConfig{
		EncPassword: ct,
		EncNonce:    nonce,
	}

	_, err = ah.DB.SaveUserConfig(*user)
	if err != nil {
		ah.logger.Error("❌ Failed to encrypted password with stellar key")
		return err
	}
	// utils.LogPretty("userCfg", userCfg)
	return nil
}

// -----------------------------
// JWT Token
// -----------------------------
// VerifyToken checks a JWT (access or refresh).
func (ah *AuthHandler) VerifyToken(token string) (*auth.Claims, error) {
	utils.LogPretty("token in app.go", token)
	return ah.auth.VerifyToken(token)
}

// RefreshToken:
// - fetch refresh from DB using userID
// - validate refresh (string only!)
// - issue new pair
// - persist new refresh (rotation)
// - return ONLY new access token to the frontend
// AuthHandler.go
func (ah *AuthHandler) RefreshToken(userID string) (*auth.TokenPairs, error) {
	utils.LogPretty("Auth - RefreshToken - userID", userID)
	// 1) Get refresh token from DB
	tokenPair, err := ah.DB.GetJwtTokenByUserId(userID)
	if err != nil || tokenPair.RefreshToken == "" {
		return nil, errors.New("no active session")
	}
	utils.LogPretty("RefreshToken - tokenpair", tokenPair)

	// 2) Validate refresh token string
	claims, err := ah.auth.VerifyToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	// utils.LogPretty("claims", claims)

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("refresh token expired; please login again")
	}

	if claims.Subject != fmt.Sprint(userID) {
		return nil, errors.New("refresh token does not belong to this user")
	}

	// 3) Load user
	user, err := ah.DB.FindUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("RefreshToken - user not found: %w", err)
	}

	u := auth.JwtUser{ID: user.ID, Username: user.Username, Email: user.Email}

	// 4) Generate new token pair
	newTokens, err := ah.auth.GenerateTokenPair(&u)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// 5) Save new refresh token
	id, err := strconv.Atoi(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert userID to string: %w", err)
	}
	_, errToken := ah.DB.UpdateJwtToken(id, newTokens)
	if errToken != nil {
		return nil, fmt.Errorf("failed to save new refresh token: %w", errToken)
	}

	// 6) Return only new access token
	return &newTokens, nil
}

// Logout removes refresh token so session is invalidated.
func (ah *AuthHandler) Logout(userID string) error {
	return ah.DB.DeleteJwtToken(userID)
}

// -----------------------------
// Middleware
// -----------------------------
func (ah *AuthHandler) RequireAuth(jwtToken string) (*auth.Claims, error) {
	claims, err := ah.VerifyToken(jwtToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	utils.LogPretty("claims", claims)
	return claims, nil
}

func (ah *AuthHandler) GetProfile(jwtToken string) (*models.User, error) {
	claims, err := ah.RequireAuth(jwtToken)
	if err != nil {
		return nil, err
	}

	user, err := ah.DB.FindUserById(claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}



	