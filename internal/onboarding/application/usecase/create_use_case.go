package onboarding_usecase

import (
	"context"
	"fmt"
	"time"
	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_domain "vault-app/internal/onboarding/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StellarServiceInterface interface {
	CreateAccount(plainPassword string) (*blockchain.CreateAccountRes, error)
	CreateKeypair() (string, string, string, error)
}
// OnBoardingUserRepository
type UserServiceInterface interface {
	Create(user *onboarding_domain.User) (*onboarding_domain.User, error)
	FindByEmail(email string) (*onboarding_domain.User, error)
}

type CreateAccountUseCase struct {
	StellarService StellarServiceInterface
	UserRepo    UserServiceInterface
	Bus            onboarding_application_events.OnboardingEventBus
    Logger         *logger.Logger
}

func NewCreateAccountUseCase(stellarService StellarServiceInterface, userRepo UserServiceInterface, eventBus onboarding_application_events.OnboardingEventBus, logger *logger.Logger) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		StellarService: stellarService,
		UserRepo:    userRepo,    
		Bus:            eventBus,
	    Logger:         logger,
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
    TxID string `json:"tx_id,omitempty"`
}



// CreateAccount handles account creation (Step 4)
// TODO: handle non-anonymous accounts with their own stellar keypair
func (a *CreateAccountUseCase) Execute(req AccountCreationRequest) (*AccountCreationResponse, error) {
	// I. ---------- Check existing user ----------
	if existingUser, _ := a.UserRepo.FindByEmail(req.Email); existingUser != nil {
		return nil, onboarding_domain.ErrUserExists
	}

	// II. ---------- Anonymous Case ----------
	if req.IsAnonymous {
		// 1. ---------- Create Stellar account for anonymous or user account (included encrypted password) ----------
		pub, secret, txID, err := a.StellarService.CreateKeypair()
		if err != nil {
			return nil, err
		}

		// 2. ---------- Create user Onboarding with Stellar key as identifier ----------
		user := &onboarding_domain.User{
			IsAnonymous:      true,
			StellarPublicKey: pub,
			CreatedAt:        time.Now(),
		}
        // utils.LogPretty("user", user)

		createdUser, err := a.UserRepo.Create(user)
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

		if err := a.Bus.PublishCreated(context.Background(), accountCreatedEvent); err != nil {
			return nil, err
		}

		return &AccountCreationResponse{
			UserID:     createdUser.ID,
			StellarKey: pub,
			SecretKey:  secret, // MUST be saved by user, just for dev
            TxID: txID,
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
        ID: uuid.New().String(),
		Email:     req.Email,
		IsAnonymous: false,
        Password: string(hashedPassword),
		CreatedAt: time.Now(),
	}

	createdUser, err := a.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	a.Logger.Info("createdUser: %v", createdUser)

	// 3. ---------- Fire Onboarding creation event ----------
	accountCreatedEvent := onboarding_application_events.AccountCreatedEvent{
		UserID:     user.ID,
		OccurredAt: time.Now(),
	}

	if err := a.Bus.PublishCreated(context.Background(), accountCreatedEvent); err != nil {
		return nil, err
	}

	return &AccountCreationResponse{
		UserID: user.ID,
	}, nil
}
