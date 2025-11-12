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
	"vault-app/internal/registry"
	"vault-app/internal/tracecore"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB              models.DBModel
	Vaults          *VaultHandler
	IPFS            *blockchain.IPFSClient
	NowUTC          func() string
	logger          logger.Logger
	TracecoreClient *tracecore.TracecoreClient
	auth            auth.Auth
}

func NewAuthHandler(db models.DBModel, vaults *VaultHandler, ipfs *blockchain.IPFSClient, logger *logger.Logger, tc *tracecore.TracecoreClient,
	auth auth.Auth) *AuthHandler {
	return &AuthHandler{
		DB:              db,
		Vaults:          vaults,
		IPFS:            ipfs,
		NowUTC:          func() string { return time.Now().Format(time.RFC3339) },
		logger:          *logger,
		TracecoreClient: tc,
		auth:            auth,
	}
}

// -----------------------------
// Sign In
// -----------------------------
type LoginRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	PublicKey     string `json:"publicKey,omitempty"`     // optional
	SignedMessage string `json:"signedMessage,omitempty"` // optional
	Signature     string `json:"signature,omitempty"`     // optional
}
type LoginResponse struct {
	User   models.User
	Vault  models.VaultPayload
	Tokens auth.TokenPairs
}

func (ah *AuthHandler) Login(credentials LoginRequest) (*LoginResponse, error) {
	var user *models.User
	var err error

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
		credentials.Password = plainPassword

	} else {
		ah.logger.Info("📧 Email login request: %s", credentials.Email)

		user, err = ah.DB.GetUserByEmail(credentials.Email)
		if err != nil || user == nil {
			return nil, fmt.Errorf("❌ User not found: %w", err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			return nil, fmt.Errorf("❌ invalid credentials: %w", err)
		}
	}
	ah.logger.Info("✅ User found: %s", user.Email)

	// -----------------------------
	// 2. Always update last connection
	// -----------------------------
	user.LastConnectedAt = time.Now().UTC()
	if saved, errSave := ah.DB.TouchLastConnected(user.ID); errSave != nil {
		ah.logger.Error("❌ failed to update last connection: %v", errSave)
	} else {
		utils.LogPretty("last connected user", saved)
	}

	// 9. create a jwt user
	u := auth.JwtUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	// generate tokens
	tokens, err := ah.auth.GenerateTokenPair(&u)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to generate token for user %d - %w", user.ID, err)
	}
	//  Save tokens to DB (for persistence across restarts)
	savedtoken, err := ah.DB.SaveJwtToken(tokens)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to persist tokens: %w", err)
	}
	ah.logger.Info("saved token: ", savedtoken.Token)

	// -----------------------------
	// 3. Try to reuse existing session
	// -----------------------------
	if existingSession, ok := ah.Vaults.Sessions[user.ID]; ok {
		if existingSession.Dirty {
			ah.Vaults.MarkDirty(user.ID)
		}
		ah.logger.Info("♻️ Reusing in-memory session for user %d", user.ID)
	utils.LogPretty("session", existingSession)
		
		return &LoginResponse{
			User:   *user,
			Vault:  *existingSession.Vault,
			Tokens: tokens,
		}, nil
	}

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
					ah.logger.Error("❌ Failed to queue commit for user %d: %v", user.ID, err)
				}
			}
		}

		if storedSession.Dirty {
			ah.Vaults.MarkDirty(user.ID)
		}
		ah.logger.Info("🔄 Restored session for user %d from DB", user.ID)
		
		return &LoginResponse{
			User:   *user,
			Vault:  *storedSession.Vault,
			Tokens: tokens,
		}, nil
	}

	// -----------------------------
	// 5. Fresh login → fetch vault
	// -----------------------------
	vaultMeta, err := ah.DB.GetLatestVaultCIDByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("❌ vault metadata lookup failed: %w", err)
	}
	if vaultMeta == nil {
		return nil, fmt.Errorf("❌ no vault metadata found for user %d", user.ID)
	}
	utils.LogPretty("✅ vaultMeta", vaultMeta)

	rawVault, err := ah.IPFS.GetData(vaultMeta.CID)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to fetch vault from IPFS: %w", err)
	}
	if rawVault == nil || len(rawVault) == 0 {
		return nil, fmt.Errorf("❌ empty vault data for CID %s", vaultMeta.CID)
	}
	utils.LogPretty("✅ rawVault", rawVault)

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
	appCfg, err := ah.DB.GetAppConfigByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to load app config: %w", err)
	}
	if appCfg == nil {
		return nil, fmt.Errorf("❌ app config missing for user %d", user.ID)
	}
	userCfg, err := ah.DB.GetUserConfigByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to load user config: %w", err)
	}
	if userCfg == nil {
		return nil, fmt.Errorf("❌ user config missing for user %d", user.ID)
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
	ah.logger.Info("✅ Vault session started for user %d", user.ID)

	// -----------------------------
	// 10. Return login response
	// -----------------------------
	return &LoginResponse{
		User:   *user,
		Vault:  vaultPayload,
		Tokens: tokens,
	}, nil
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
	// 4. (Optional) Bootstrap Tracecore Repo
	// -----------------------------
	template, ok := registry.VaultTemplates[setup.RepoTemplate]
	var repoIDStr string
	if setup.RepoTemplate != "" && ah.TracecoreClient != nil {
		if ok {
			repoID, err := template.ShouldBootstrapTracecoreRepo(ah.TracecoreClient)
			if err == nil && repoID != nil {
				repoIDStr = *repoID
			}
		}
	}
	ah.logger.Info("Repo ID: %d", repoIDStr)

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
		ID:             strconv.Itoa(user.ID),
		Role:           "to-define",
		Signature:      "", // will be generated later
		StellarAccount: *stellarAccount,
		SharingRules:   sharingRules,
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

	userCfg, err := ah.DB.SaveUserConfig(*user)
	if err != nil {
		ah.logger.Error("❌ Failed to encrypted password with stellar key")
		return err
	}
	utils.LogPretty("userCfg", userCfg)
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
func (ah *AuthHandler) RefreshToken(userID int) (string, error) {
	// 1) Get refresh token from DB
	tokenPair, err := ah.DB.GetJwtTokenByUserId(userID)
	if err != nil || tokenPair.RefreshToken == "" {
		return "", errors.New("no active session")
	}
	utils.LogPretty("tokenpair", tokenPair)

	// 2) Validate refresh token string
	claims, err := ah.auth.VerifyToken(tokenPair.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}
	utils.LogPretty("claims", claims)

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return "", errors.New("refresh token expired; please login again")
	}

	if claims.Subject != fmt.Sprint(userID) {
		return "", errors.New("refresh token does not belong to this user")
	}

	// 3) Load user
	user, err := ah.DB.FindUserById(userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	u := auth.JwtUser{ID: user.ID, Username: user.Username, Email: user.Email}

	// 4) Generate new token pair
	newTokens, err := ah.auth.GenerateTokenPair(&u)
	if err != nil {
		return "", fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// 5) Save new refresh token
	_, errToken := ah.DB.UpdateJwtToken(user.ID, newTokens)
	if errToken != nil {
		return "", fmt.Errorf("failed to save new refresh token: %w", errToken)
	}

	// 6) Return only new access token
	return newTokens.Token, nil
}

// Logout removes refresh token so session is invalidated.
func (ah *AuthHandler) Logout(userID int) error {
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
	fmt.Println("UserID:", claims.UserID)     // int
	fmt.Println("Username:", claims.Username) // string
	fmt.Println("Email:", claims.Email)       // string
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
