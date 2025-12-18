package onboarding_usecase

import (
	"context"
	"fmt"
	"time"
	utils "vault-app/internal"
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

type UserServiceInterface interface {
	Create(user *onboarding_domain.User) (*onboarding_domain.User, error)
}

type CreateAccountUseCase struct {
	StellarService StellarServiceInterface
	UserService    UserServiceInterface
	Bus            onboarding_application_events.OnboardingEventBus
    Logger         *logger.Logger
}

func NewCreateAccountUseCase(stellarService StellarServiceInterface, userService UserServiceInterface, eventBus onboarding_application_events.OnboardingEventBus, logger *logger.Logger) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		StellarService: stellarService,
		UserService:    userService,    
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
func (a *CreateAccountUseCase) Execute(req AccountCreationRequest) (*AccountCreationResponse, error) {
	if req.IsAnonymous {
		// Create Stellar account for anonymous or user account (included encrypted password)
		pub, secret, txID, err := a.StellarService.CreateKeypair()
		if err != nil {
			return nil, err
		}
				
		// Create user with Stellar key as identifier
		user := &onboarding_domain.User{
			IsAnonymous:      true,
			StellarPublicKey: pub,
			CreatedAt:        time.Now(),
		}
        utils.LogPretty("user", user)

		createdUser, err := a.UserService.Create(user)
		if err != nil {
			return nil, err
		}

		// Fire creation event
		accountCreatedEvent := onboarding_application_events.AccountCreatedEvent{
			UserID:           user.ID,
			StellarPublicKey: pub,
			OccurredAt:       time.Now(),
		}
        utils.LogPretty("accountCreatedEvent", accountCreatedEvent)
        fmt.Println("accountCreatedEvent", accountCreatedEvent)
        // a.Logger.Info("accountCreatedEvent", accountCreatedEvent)   

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
		// -----------------------------
	// 2. Hash password & create user
	// -----------------------------
	utils.LogPretty("Plain password", req.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to hash password: %w", err)
	}

	// Standard email/password account
	user := &onboarding_domain.User{
        ID: uuid.New().String(),
		Email:     req.Email,
        Password: string(hashedPassword),
		CreatedAt: time.Now(),
	}

	createdUser, err := a.UserService.Create(user)
	if err != nil {
		return nil, err
	}
	fmt.Println("createdUser", createdUser)

	// Fire creation event
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
