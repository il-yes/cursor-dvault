package onboarding_usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"
	"vault-app/internal/blockchain"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_domain "vault-app/internal/onboarding/domain"
	"vault-app/internal/utils"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StellarServiceInterface interface {
	CreateAccount(plainPassword string) (*blockchain.CreateAccountRes, error)
	CreateKeypair() (string, string, string, error)
}
type KeyringServiceInterface interface {
    SaveHybrid(
        kr *vaults_domain.VaultKeyring,
        userID string,
        password string,
        stellarSecret string,
    ) error
}
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// OnBoardingUserRepository
type UserServiceInterface interface {
	Create(user *onboarding_domain.User) (*onboarding_domain.User, error)
	FindByEmail(email string) (*onboarding_domain.User, error)
}

type CreateAccountUseCase struct {
	StellarService StellarServiceInterface
	UserRepo       UserServiceInterface
	Bus            onboarding_application_events.OnboardingEventBus
	Logger         Logger
	// CryptoService    blockchain.CryptoService

	KeyringService KeyringServiceInterface
	KeyEncryption  vaults_domain.KeyEncryption
}

func NewCreateAccountUseCase(
	stellarService StellarServiceInterface,
	userRepo UserServiceInterface,
	eventBus onboarding_application_events.OnboardingEventBus,
	logger Logger,
	keyringService KeyringServiceInterface,
	keyEncryption vaults_domain.KeyEncryption,
) *CreateAccountUseCase {

	return &CreateAccountUseCase{
		StellarService: stellarService,
		UserRepo:       userRepo,
		Bus:            eventBus,
		Logger:         logger,
		KeyringService: keyringService,
		KeyEncryption:  keyEncryption,
	}
}

// Step 4: Account Creation
type AccountCreationRequest struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	IsAnonymous bool   `json:"is_anonymous"`
	StellarKey  string `json:"stellar_key,omitempty"` // For anonymous accounts
}

type AccountCreationResponse struct {
	UserID     string `json:"user_id"`
	StellarKey string `json:"stellar_key,omitempty"` // Generated for anonymous
	SecretKey  string `json:"secret_key,omitempty"`  // CRITICAL: User must save this
	TxID       string `json:"tx_id,omitempty"`
}

// CreateAccount handles account creation (Step 4)
// TODO: handle non-anonymous accounts with their own stellar keypair
func (a *CreateAccountUseCase) Execute(req AccountCreationRequest) (*AccountCreationResponse, error) {
	if a.KeyEncryption == nil {
		panic("KeyEncryption is nil")
	}
	if a.KeyringService == nil {
		panic("KeyringService is nil")
	}
	if a.UserRepo == nil {
		panic("UserRepo is nil")
	}
	if a.StellarService == nil {
		panic("StellarService is nil")
	}
	// I. ---------- Check existing user ----------
	if existingUser, _ := a.UserRepo.FindByEmail(req.Email); existingUser != nil {
		a.Logger.Error("CreateAccountUseCase - FindByEmail - Failed to find user: %v", existingUser)
		return nil, onboarding_domain.ErrUserExists
	}

	// II. ---------- Anonymous Case ----------
	if req.IsAnonymous {
		// 1. ---------- Create Stellar account for anonymous or user account (included encrypted password) ----------
		pub, secret, txID, err := a.StellarService.CreateKeypair()
		if err != nil {
			a.Logger.Error("CreateAccountUseCase - CreateKeypair - Failed to create stellar account: %v", err)
			return nil, err
		}
		utils.LogPretty("CreateAccountUseCase - CreateKeypair - pub", pub)
		utils.LogPretty("CreateAccountUseCase - CreateKeypair - secret", secret)
		

		// 2. ---------- Create user Onboarding with Stellar key as identifier ----------
		user := &onboarding_domain.User{
			IsAnonymous:      true,
			StellarPublicKey: pub,
			CreatedAt:        time.Now(),
		}
		// utils.LogPretty("user", user)

		createdUser, err := a.UserRepo.Create(user)
		if createdUser == nil {
			return nil, errors.New("createdUser is nil")
		}
		if err != nil {
			return nil, err
		}

		// 3. ---------- Fire Onboarding creation event ----------
		accountCreatedEvent := onboarding_application_events.AccountCreatedEvent{
			UserID:           user.ID,
			StellarPublicKey: pub,
			OccurredAt:       time.Now(),
		}
		a.Logger.Info("accountCreatedEvent", accountCreatedEvent)

		// 3. ---------- Create Vault key for user ----------
		// 1. Generate VaultKey (DEK)
		vaultKey := make([]byte, 32)
		if _, err := rand.Read(vaultKey); err != nil {
			return nil, fmt.Errorf("failed to generate vault key: %w", err)
		}
		utils.LogPretty("CreateAccountUseCase - CreateKeypair - vaultKey", vaultKey)

		// 2. Create empty keyring
		kr := &vaults_domain.VaultKeyring{
			UserID:    createdUser.ID,
			Keys:      []vaults_domain.EncryptedKey{},
			Wrappers:  []vaults_domain.WrappedKey{},
			UpdatedAt: time.Now().Unix(),
		}

		// 3. Store raw VaultKey as EncryptedKey (internal representation)
		kr.Keys = append(kr.Keys, vaults_domain.EncryptedKey{
			ID:         uuid.New().String(),
			Type:       vaults_domain.KeyTypeVault,
			Version:    1,
			Ciphertext: vaultKey,
			CreatedAt:  time.Now().Unix(),
		})
		// STELLAR WRAP
		if createdUser.StellarPublicKey != "" {
			enc, err := a.KeyEncryption.WrapKeyWithStellar(vaultKey, secret)
			if err != nil {
				return nil, err
			}

			kr.AddWrapper(vaults_domain.WrappedKey{
				ID:        uuid.New().String(),
				Type:      "stellar",
				Data:      enc,
				CreatedAt: time.Now().Format(time.RFC3339),
			})
		}

		if err := a.KeyringService.SaveHybrid(kr, createdUser.ID, "", secret); err != nil {
			return nil, fmt.Errorf("CreateAccountUseCase - Execute - failed to save keyring: %w", err)
		}
		utils.LogPretty("CreateAccountUseCase - Execute - kr", kr)

		if err := a.Bus.PublishCreated(context.Background(), accountCreatedEvent); err != nil {
			return nil, err
		}

		return &AccountCreationResponse{
			UserID:     createdUser.ID,
			StellarKey: pub,
			SecretKey:  secret, // MUST be saved by user, just for dev
			TxID:       txID,
		}, nil
	}

	// III. ---------- Standard Case ----------
	// 2. ---------- Create user Onboarding  ----------
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to hash password: %w", err)
	}

	// Standard email/password account
	user := &onboarding_domain.User{
		ID:          uuid.New().String(),
		Email:       req.Email,
		IsAnonymous: false,
		Password:    string(hashedPassword),
		CreatedAt:   time.Now(),
	}

	createdUser, err := a.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	a.Logger.Info("createdUser: %v", createdUser)

	// 3. ---------- Create Vault key for user ----------
	// 1. Generate VaultKey (DEK)
	vaultKey := make([]byte, 32)
	if _, err := rand.Read(vaultKey); err != nil {
		utils.LogPretty("CreateAccountUseCase - CreateKeypair - failed to genarate vaultKey", err)
		return nil, fmt.Errorf("failed to generate vault key: %w", err)
	}
	utils.LogPretty("CreateAccountUseCase - CreateKeypair - vaultKey", vaultKey)

	// 2. Create empty keyring
	kr := &vaults_domain.VaultKeyring{
		UserID:    createdUser.ID, // user onboarding
		Keys:      []vaults_domain.EncryptedKey{},
		Wrappers:  []vaults_domain.WrappedKey{},
		UpdatedAt: time.Now().Unix(),
	}
	utils.LogPretty("CreateAccountUseCase - CreateKeypair - kr", kr)

	// 3. Store raw VaultKey as EncryptedKey (internal representation)
	kr.Keys = append(kr.Keys, vaults_domain.EncryptedKey{
		ID:         uuid.New().String(),
		Type:       vaults_domain.KeyTypeVault,
		Version:    1,
		Ciphertext: vaultKey,
		CreatedAt:  time.Now().Unix(),
	})
	utils.LogPretty("CreateAccountUseCase - CreateKeypair - kr.Keys", kr.Keys)

	// save keyring WITH USER ONBOARDING
	if err := a.KeyringService.SaveHybrid(kr, createdUser.ID, req.Password, ""); err != nil {
		utils.LogPretty("CreateAccountUseCase - CreateKeypair - failed to save keyring", err)
		return nil, fmt.Errorf("CreateAccountUseCase - Execute - failed to save keyring: %w", err)
	}

	// 4. ---------- Fire Onboarding creation event ----------
	accountCreatedEvent := onboarding_application_events.AccountCreatedEvent{
		UserID:     user.ID,
		OccurredAt: time.Now(),
	}
	utils.LogPretty("CreateAccountUseCase - CreateKeypair - accountCreatedEvent", accountCreatedEvent)

	if err := a.Bus.PublishCreated(context.Background(), accountCreatedEvent); err != nil {
		utils.LogPretty("CreateAccountUseCase - CreateKeypair - failed to publish event", err)
		return nil, err
	}

	return &AccountCreationResponse{
		UserID: user.ID,
	}, nil
}
