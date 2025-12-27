package identity_commands

import (
	"context"
	"time"
	utils "vault-app/internal"
	auth_domain "vault-app/internal/auth/domain"
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
	}
	utils.LogPretty("Resolved credentials", creds)
	if cmd.PublicKey != "" {
		plain, err := h.resolveStellarPassword(cmd)
		if err != nil {
			return nil, err
		}
		creds.Password = plain
	}
	utils.LogPretty("Resolved credentials", creds)
	// 2. Onboarding - Authenticate
	onboardingUser, err := h.onboardingRepo.FindByEmail(creds.Email)
	if err != nil || onboardingUser == nil {
		return nil, auth_domain.ErrInvalidCredentials
	}
	utils.LogPretty("Authenticated user", onboardingUser)

	if err := bcrypt.CompareHashAndPassword(
		[]byte(onboardingUser.Password),
		[]byte(creds.Password),
	); err != nil {
		return nil, auth_domain.ErrInvalidCredentials
	}
	utils.LogPretty("Authenticated bcrypt checked", "true")
	// 3. Identity - Load identity user
	utils.LogPretty("Loading identity user by email", creds.Email)
	user, err := h.userRepo.FindByEmail(context.Background(), creds.Email)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("Loaded identity user", user)

	user.LastConnectedAt = time.Now()
	if err := h.userRepo.Update(context.Background(), user); err != nil {
		return nil, err
	}

	// 5. Auth - Generate tokens
	tokens, err := tokenService.GenerateTokenPair(user.ToJwtUser())
	if err != nil {
		utils.LogPretty("Failed to generate tokens", user)
		return nil, err
	}
	utils.LogPretty("Generated tokens", tokens)

	if _, err := tokenService.SaveJwtToken(tokens); err != nil {
		utils.LogPretty("Failed to persist tokens", tokens)
		return nil, err
	}
	utils.LogPretty("Persisted tokens", tokens)

	// 7. Publish event
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

func (h *LoginCommandHandler) resolveStellarPassword(cmd LoginCommand) (string, error) {
	return "", nil
}

