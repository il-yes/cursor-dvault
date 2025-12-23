package identity_usecase

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	identity_eventbus "vault-app/internal/identity/application"
	identity_domain "vault-app/internal/identity/domain"
)

// IDGen defines a function that generates unique IDs
type IDGen func() string

// LoginUseCase handles user login
type LoginUseCase struct {
	repo  identity_domain.UserRepository
	bus   identity_eventbus.EventBus
	idGen IDGen
}

// NewLoginUseCase constructs a LoginUseCase
func NewLoginUseCase(repo identity_domain.UserRepository, bus identity_eventbus.EventBus, idGen IDGen) *LoginUseCase {
	return &LoginUseCase{
		repo:  repo,
		bus:   bus,
		idGen: idGen,
	}
}

// Execute performs login, returns user on success
func (uc *LoginUseCase) Execute(
	ctx context.Context,
	email, password, publicKey, signedMessage, signature string,
) (*identity_domain.User, error) {
	// 1️⃣ Lookup user by email
	user, _ := uc.repo.FindByEmail(ctx, email)

	// 2️⃣ Resolve password if publicKey login
	if publicKey != "" {
		var err error
		password, err = uc.resolveStellarPassword(publicKey, signedMessage, signature)
		if err != nil {
			return nil, identity_domain.ErrInvalidCredentials
		}
	}

	// 3️⃣ Safe password comparison to prevent timing attacks
	storedHash := ""
	if user != nil {
		storedHash = user.PasswordHash
	}
	if err := uc.safePasswordCompare(storedHash, password); err != nil {
		return nil, identity_domain.ErrInvalidCredentials
	}

	// 4️⃣ If user does not exist, create new one (optional behavior)
	if user == nil {
		id := uc.idGen()
		user = identity_domain.NewStandardUser(id, email, password)
		if err := uc.repo.Save(ctx, user); err != nil {
			return nil, err
		}
	}

	// 5️⃣ Update last connected timestamp
	user.LastConnectedAt = time.Now()
	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	// 6️⃣ Publish UserLoggedIn event
	if uc.bus != nil {
		event := identity_eventbus.UserLoggedIn{
			UserID:    user.ID,
			Email:     user.Email,
			OccurredAt: time.Now(),
		}
		_ = uc.bus.PublishUserLoggedIn(ctx, event)
	}

	return user, nil
}

// safePasswordCompare ensures constant-time comparison to prevent timing attacks
func (uc *LoginUseCase) safePasswordCompare(storedHash, password string) error {
	if storedHash == "" {
		// Compare against dummy hash
		dummy, _ := bcrypt.GenerateFromPassword([]byte("dummy"), bcrypt.DefaultCost)
		storedHash = string(dummy)
	}
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
}

// resolveStellarPassword resolves a password from public key login
func (uc *LoginUseCase) resolveStellarPassword(pubKey, signedMessage, signature string) (string, error) {
	if pubKey == "" || signedMessage == "" || signature == "" {
		return "", identity_domain.ErrInvalidCredentials
	}
	// TODO: verify signature against pubKey (Stellar logic)
	// For now, return a derived placeholder password
	return "derived-password", nil
}
