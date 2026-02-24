package identity_commands

import (
	"context"
	"log"
	"time"
	utils "vault-app/internal/utils"
	auth_domain "vault-app/internal/auth/domain"
	"vault-app/internal/blockchain"
	identity_eventbus "vault-app/internal/identity/application"
	identity_domain "vault-app/internal/identity/domain"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	onboarding_domain "vault-app/internal/onboarding/domain"
	onboarding_persistence "vault-app/internal/onboarding/infrastructure/persistence"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// internal/identity/application/commands/login_command.go
type LoginCommand struct {
	Email         string
	Password      string
	PublicKey     string
	SignedMessage string
	Signature     string
}

type LoginResult struct {
	User      *identity_domain.User
	Tokens    auth_domain.TokenPairs
}
type ManagerInterface interface {
	Get(userID string) (*vault_session.Session, bool)
	Prepare(userID string) (*vault_session.Session, error)
	AttachVault(
		userID string,
		vault *vaults_domain.VaultPayload,
		runtime *vault_session.RuntimeContext,
		lastCID string,
	) *vault_session.Session
	MarkDirty(userID string)
	Close(userID string)
}
type TokenServiceInterface interface {
	GenerateTokenPair(user *auth_domain.JwtUser) (auth_domain.TokenPairs, error)
	Persist(tp auth_domain.TokenPairs) error
    SaveJwtToken(tokens auth_domain.TokenPairs) (*auth_domain.TokenPairs, error)
}

type LoginCommandHandler struct {
	onboardingRepo onboarding_domain.UserRepository
	userRepo       identity_domain.UserRepository
	NowUTC         func() string
}

func NewLoginCommandHandler(
	db *gorm.DB,
) *LoginCommandHandler {
	userRepo := identity_persistence.NewGormUserRepository(db)
	onboardingRepo := onboarding_persistence.NewGormUserRepository(db)	

	return &LoginCommandHandler{
		onboardingRepo: onboardingRepo,
		userRepo:       userRepo,
		NowUTC:         func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *LoginCommandHandler) Handle(
	cmd LoginCommand, tokenService TokenServiceInterface, eventBus identity_eventbus.EventBus,
) (*LoginResult, error) {

	utils.LogPretty("Login command", cmd)
	// 1. Resolve credentials
	creds := auth_domain.Credentials{
		Email:    cmd.Email,
		Password: cmd.Password,
		PublicKey: cmd.PublicKey,
		SignedMessage: cmd.SignedMessage,
		Signature: cmd.Signature,
	}

	// 2. Identity - Authenticate
	if cmd.PublicKey != "" {
		// 2.1. Resolve stellar credentials if publicKey login
		var err error
		email, err := h.resolveStellarPassword(cmd.PublicKey, cmd.SignedMessage, cmd.Signature)
		if err != nil {
			log.Println("❌ resolveStellarPassword - Error: ", err)
			return nil, identity_domain.ErrInvalidCredentials
		}
		if email == nil {
			log.Println("❌ LoginCommandHandler.  resolveStellarPassword - Authentication failed")
			return nil, identity_domain.ErrInvalidCredentials
		}
		creds.Email = *email
		utils.LogPretty("Authenticated with stellar", creds.Email)
	} else {
		// 2.2. Onboarding - Load user
		onboardingUser, err := h.onboardingRepo.FindByEmail(creds.Email)
		if err != nil || onboardingUser == nil {
			return nil, auth_domain.ErrInvalidCredentials
		}
		utils.LogPretty("LoginCommandHandler - Onboarded user", onboardingUser)
		// 2.3. Resolve password if publicKey login
		if err := bcrypt.CompareHashAndPassword(
			[]byte(onboardingUser.Password),
			[]byte(creds.Password),
		); err != nil {
			return nil, auth_domain.ErrInvalidCredentials
		}
		utils.LogPretty("LoginCommandHandler - Authenticated with bcrypt", "true")
	}
	
	// 3. Identity - Load identity user
	utils.LogPretty("LoginCommandHandler - Loading identity user by email", creds.Email)
	user, err := h.userRepo.FindByEmail(context.Background(), creds.Email)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("LoginCommandHandler - Loaded identity user", user)

	user.LastConnectedAt = time.Now()
	if err := h.userRepo.Update(context.Background(), user); err != nil {
		return nil, err
	}

	// 4. Auth - Generate tokens
	tokens, err := tokenService.GenerateTokenPair(user.ToJwtUser())
	if err != nil {
		utils.LogPretty("Failed to generate tokens", user)
		return nil, err
	}
	utils.LogPretty("LoginCommandHandler - Generated tokens", tokens)

	if _, err := tokenService.SaveJwtToken(tokens); err != nil {
		utils.LogPretty("Failed to persist tokens", tokens)
		return nil, err
	}
	utils.LogPretty("LoginCommandHandler - Persisted tokens", tokens)

	// 5. Publish event
	if eventBus != nil {
		_ = eventBus.PublishUserLoggedIn(context.Background(), identity_eventbus.UserLoggedIn{
			UserID:     user.ID,
			Email:      user.Email,
			OccurredAt: time.Now(),
		})
	}

	return &LoginResult{
		User:      user,
		Tokens:    tokens,
	}, nil
}


// resolveStellarPassword resolves a password from public key login
func (uc *LoginCommandHandler) resolveStellarPassword(pubKey string, signedMessage, signature string) (*string, error) {
	if pubKey == "" || signedMessage == "" || signature == "" {
		return nil, identity_domain.ErrInvalidCredentials
	}
	// TODO: verify signature against pubKey (Stellar logic)
	isVerified := blockchain.VerifySignature(pubKey, signedMessage, signature)
	if !isVerified {
		log.Println("❌ resolveStellarPassword - Signature verification failed")
		return 	nil, identity_domain.ErrInvalidCredentials
	}
	utils.LogPretty("resolveStellarPassword - Signature verified", isVerified)
	// Find user by pubKey
	user, err := uc.userRepo.FindByPublicKey(context.Background(), pubKey)
	if err != nil {
		return nil, identity_domain.ErrInvalidCredentials
	}

	return &user.Email, nil
}


