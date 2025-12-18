package identity_usecase

import (
	"context"
	identity_eventbus "vault-app/internal/identity/application"
	identity_domain "vault-app/internal/identity/domain"
)


// RegisterStandardUserUseCase registers a non-anonymous user
type RegisterStandardUserUseCase struct {
	repo  identity_domain.UserRepository
	bus   identity_eventbus.EventBus
	idGen IDGen
}

func NewRegisterStandardUserUseCase(repo identity_domain.UserRepository, bus identity_eventbus.EventBus, idGen IDGen) *RegisterStandardUserUseCase {
	return &RegisterStandardUserUseCase{repo: repo, bus: bus, idGen: idGen}
}

func (uc *RegisterStandardUserUseCase) Execute(ctx context.Context, email, passwordHash string) (*identity_domain.User, error) {
	// check existing
	if existing, _ := uc.repo.FindByEmail(ctx, email); existing != nil {
		return nil, identity_domain.ErrUserExists
	}
	id := uc.idGen()
	u := identity_domain.NewStandardUser(id, email, passwordHash)
	if err := uc.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	// publish event
	domainEvent := identity_domain.NewUserRegistered(u)
	if uc.bus != nil {
		// Convert domain event to application event
		appEvent := identity_eventbus.UserRegistered{
			UserID:      domainEvent.UserID,
			IsAnonymous: domainEvent.IsAnonymous,
			OccurredAt:  domainEvent.OccurredAt,
		}
		_ = uc.bus.PublishUserRegistered(ctx, appEvent)
	}
	return u, nil
}

// RegisterAnonymousUserUseCase registers an anonymous user with stellar key
type RegisterAnonymousUserUseCase struct {
	repo  identity_domain.UserRepository
	bus   identity_eventbus.EventBus
	idGen IDGen
}

func NewRegisterAnonymousUserUseCase(repo identity_domain.UserRepository, bus identity_eventbus.EventBus, idGen IDGen) *RegisterAnonymousUserUseCase {
	return &RegisterAnonymousUserUseCase{repo: repo, bus: bus, idGen: idGen}
}

func (uc *RegisterAnonymousUserUseCase) Execute(ctx context.Context, stellarPublicKey string) (*identity_domain.User, error) {
	id := uc.idGen()
	u := identity_domain.NewAnonymousUser(id, stellarPublicKey)
	if err := uc.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	if uc.bus != nil {
		// Convert domain event to application event
		domainEvent := identity_domain.NewUserRegistered(u)
		appEvent := identity_eventbus.UserRegistered{
			UserID:      domainEvent.UserID,
			IsAnonymous: domainEvent.IsAnonymous,
			OccurredAt:  domainEvent.OccurredAt,
		}
		_ = uc.bus.PublishUserRegistered(ctx, appEvent)
	}
	return u, nil
}
