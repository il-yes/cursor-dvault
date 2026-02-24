package identity_usecase

import (
	"context"
	"log"
	"time"

	"vault-app/internal/blockchain"
	identity_eventbus "vault-app/internal/identity/application"
	identity_domain "vault-app/internal/identity/domain"

	"golang.org/x/crypto/bcrypt"
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

	if publicKey != "" {
		// 2️⃣ Resolve password if publicKey login
		var err error
		isAuthenticated, err := uc.resolveStellarPassword(publicKey, signedMessage, signature, email)
		if err != nil {
			log.Println("❌ resolveStellarPassword - Error: ", err)
			return nil, identity_domain.ErrInvalidCredentials
		}
		if !isAuthenticated {
			log.Println("❌ resolveStellarPassword - Authentication failed")
			return nil, identity_domain.ErrInvalidCredentials
		}
	} else {
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
func (uc *LoginUseCase) resolveStellarPassword(pubKey string, signedMessage, signature string, email string) (bool, error) {
	if pubKey == "" || signedMessage == "" || signature == "" {
		return false, identity_domain.ErrInvalidCredentials
	}
	// TODO: verify signature against pubKey (Stellar logic)
	isVerified := blockchain.VerifySignature(pubKey, signedMessage, signature)
	if !isVerified {
		log.Println("❌ resolveStellarPassword - Signature verification failed")
		return false, identity_domain.ErrInvalidCredentials
	}
	// Find user by pubKey
	user, err := uc.repo.FindByPublicKey(context.Background(), pubKey)
	if err != nil {
		return false, identity_domain.ErrInvalidCredentials
	}
	return user.Email == email, nil
}
