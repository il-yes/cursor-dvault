package identity_commands_test

import (
	"context"
	"errors"
	"testing"

	auth_domain "vault-app/internal/auth/domain"
	identity_eventbus "vault-app/internal/identity/application"
	identity_commands "vault-app/internal/identity/application/commands"
	identity_domain "vault-app/internal/identity/domain"
	onboarding_domain "vault-app/internal/onboarding/domain"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

//
// -------- FAKES --------
//

type fakeOnboardingRepo struct {
	user *onboarding_domain.User
	err  error
}

func (f *fakeOnboardingRepo) FindByEmail(email string) (*onboarding_domain.User, error) {
	return f.user, f.err
}

func (f *fakeOnboardingRepo) GetByID(id string) (*onboarding_domain.User, error) {
	return nil, nil
}

func (f *fakeOnboardingRepo) Create(u *onboarding_domain.User) (*onboarding_domain.User, error) {
	return nil, nil
}

func (f *fakeOnboardingRepo) Delete(id string) error {
	return nil
}
func (f *fakeOnboardingRepo) List() ([]onboarding_domain.User, error) {
	return nil, nil
}
func (f *fakeOnboardingRepo) Update(u *onboarding_domain.User) error {
	return nil
}
func (f *fakeOnboardingRepo) FinndByEmail(email string) (*onboarding_domain.User, error) {
	return nil, nil
}

type fakeIdentityRepo struct {
	user       *identity_domain.User
	findErr    error
	saveCalled bool
	saveErr    error
}

func (f *fakeIdentityRepo) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	return f.user, f.findErr
}

func (f *fakeIdentityRepo) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	return nil, nil
}

func (f *fakeIdentityRepo) Save(ctx context.Context, u *identity_domain.User) error {
	f.saveCalled = true
	f.user = u
	return f.saveErr
}
func (f *fakeIdentityRepo) Update(ctx context.Context, u *identity_domain.User) error {
	return nil
}
type fakeTokenService struct {
	generateErr error
	persistErr  error
	generated   bool
	persisted   bool
}

func (f *fakeTokenService) GenerateTokenPair(user *auth_domain.JwtUser) (auth_domain.TokenPairs, error) {
	if f.generateErr != nil {
		return auth_domain.TokenPairs{}, f.generateErr
	}
	f.generated = true
	return auth_domain.TokenPairs{
		Token:        "access",
		RefreshToken: "refresh",
		UserID:       user.ID,
	}, nil
}

func (f *fakeTokenService) Persist(tp auth_domain.TokenPairs) error {
	if f.persistErr != nil {
		return f.persistErr
	}
	f.persisted = true
	return nil
}
func (f *fakeTokenService) SaveJwtToken(tokens auth_domain.TokenPairs) (*auth_domain.TokenPairs, error) {
	return &tokens, nil
}	

type fakeSessionManager struct {
	sessionID string
	prepared  bool
}

func (f *fakeSessionManager) Get(userID string) (*vault_session.Session, bool) {
	return nil, false
}

func (f *fakeSessionManager) Prepare(userID string) string {
	f.prepared = true
	return f.sessionID
}

func (f *fakeSessionManager) AttachVault(
	userID string,
	vault *vaults_domain.VaultPayload,
	runtime *vault_session.RuntimeContext,
	lastCID string,
) *vault_session.Session {
	return nil
}

func (f *fakeSessionManager) MarkDirty(userID string) {
}

func (f *fakeSessionManager) Close(userID string) {
}

/* ------------------------------
   Fake EventBus
--------------------------------*/
type fakeEventBus struct {
	eventBus identity_eventbus.EventBus
	publishedLogins []identity_eventbus.UserLoggedIn
}

func newFakeEventBus() *fakeEventBus {
	return &fakeEventBus{
		publishedLogins: []identity_eventbus.UserLoggedIn{},
	}
}

func (f *fakeEventBus) PublishUserLoggedIn(ctx context.Context, e identity_eventbus.UserLoggedIn) error {
	f.publishedLogins = append(f.publishedLogins, e)
	return nil
}

func (f *fakeEventBus) SubscribeToUserLoggedIn(handler identity_eventbus.UserLoggedInHandler) error {
	return nil
}

func (f *fakeEventBus) PublishUserRegistered(ctx context.Context, e identity_eventbus.UserRegistered) error {
	return nil
}

func (f *fakeEventBus) SubscribeToUserRegistered(handler identity_eventbus.UserRegisteredHandler) error {
	return nil
}
func (f *fakeEventBus) SetEventBus(bus identity_eventbus.EventBus) {
	f.eventBus = bus
}	
//
// -------- TESTS --------
//

func TestLoginCommandHandler_Handle(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	baseOnboardingUser := &onboarding_domain.User{
		Email:    "test@example.com",
		Password: string(hashed),
	}

	baseIdentityUser := &identity_domain.User{
		ID:    "user-123",	
		Email: "test@example.com",
	}

	tests := []struct {
		name       string
		setup      func() (*identity_commands.LoginCommandHandler, *fakeIdentityRepo, *fakeTokenService, *fakeSessionManager)
		command    identity_commands.LoginCommand
		expectErr  bool
		assertions func(t *testing.T, res *identity_commands.LoginResult, repo *fakeIdentityRepo, ts *fakeTokenService, sm *fakeSessionManager)
	}{
		{
			name: "successful login",
			setup: func() (*identity_commands.LoginCommandHandler, *fakeIdentityRepo, *fakeTokenService, *fakeSessionManager) {
				onboarding := &fakeOnboardingRepo{user: baseOnboardingUser}
				repo := &fakeIdentityRepo{user: baseIdentityUser}
				session := &fakeSessionManager{sessionID: "session-1"}
				ts := &fakeTokenService{}
				eventBus := newFakeEventBus()
				eventBus.SetEventBus(eventBus)
				handler := identity_commands.NewLoginCommandHandler(onboarding, repo, ts, session, eventBus)
				handler.NowUTC = func() string { return "now" }

				return handler, repo, ts, session
			},
			command: identity_commands.LoginCommand{
				Email:    "test@example.com",
				Password: "secret",
			},
			expectErr: false,
			assertions: func(t *testing.T, res *identity_commands.LoginResult, repo *fakeIdentityRepo, ts *fakeTokenService, sm *fakeSessionManager) {
				assert.NotNil(t, res)
				assert.Equal(t, "user-123", res.User.ID)
				assert.Equal(t, "access", res.Tokens.Token)
				assert.True(t, repo.saveCalled)
				assert.Equal(t, "now", res.User.LastConnectedAt)
				assert.True(t, ts.generated)
				assert.True(t, ts.persisted)
				assert.True(t, sm.prepared)
			},
		},
		{
			name: "invalid email",
			setup: func() (*identity_commands.LoginCommandHandler, *fakeIdentityRepo, *fakeTokenService, *fakeSessionManager) {
				onboarding := &fakeOnboardingRepo{user: nil}
				repo := &fakeIdentityRepo{}
				session := &fakeSessionManager{}
				ts := &fakeTokenService{}
				eventBus := newFakeEventBus()
				eventBus.SetEventBus(eventBus)
				handler := identity_commands.NewLoginCommandHandler(onboarding, repo, ts, session, eventBus)
				return handler, repo, ts, session
			},
			command: identity_commands.LoginCommand{
				Email:    "missing@example.com",
				Password: "secret",
			},
			expectErr: true,
		},
		{
			name: "invalid password",
			setup: func() (*identity_commands.LoginCommandHandler, *fakeIdentityRepo, *fakeTokenService, *fakeSessionManager) {
				onboarding := &fakeOnboardingRepo{user: baseOnboardingUser}
				repo := &fakeIdentityRepo{user: baseIdentityUser}
				ts := &fakeTokenService{}
				session := &fakeSessionManager{}
				eventBus := newFakeEventBus()
				eventBus.SetEventBus(eventBus)
				handler := identity_commands.NewLoginCommandHandler(onboarding, repo, ts, session, eventBus)
				return handler, repo, ts, session
			},
			command: identity_commands.LoginCommand{
				Email:    "test@example.com",
				Password: "wrong",
			},
			expectErr: true,
		},
		{
			name: "identity user not found",
			setup: func() (*identity_commands.LoginCommandHandler, *fakeIdentityRepo, *fakeTokenService, *fakeSessionManager) {
				onboarding := &fakeOnboardingRepo{user: baseOnboardingUser}
				repo := &fakeIdentityRepo{user: nil, findErr: errors.New("not found")}
				session := &fakeSessionManager{}
				ts := &fakeTokenService{}
				eventBus := newFakeEventBus()
				eventBus.SetEventBus(eventBus)
				handler := identity_commands.NewLoginCommandHandler(onboarding, repo, ts, session, eventBus)
				return handler, repo, ts, session
			},
			command: identity_commands.LoginCommand{
				Email:    "test@example.com",
				Password: "secret",
			},
			expectErr: true,
		},
		{
			name: "token generation fails",
			setup: func() (*identity_commands.LoginCommandHandler, *fakeIdentityRepo, *fakeTokenService, *fakeSessionManager) {
				onboarding := &fakeOnboardingRepo{user: baseOnboardingUser}
				repo := &fakeIdentityRepo{user: baseIdentityUser}
				session := &fakeSessionManager{}
				ts := &fakeTokenService{generateErr: errors.New("boom")}
				eventBus := newFakeEventBus()
				eventBus.SetEventBus(eventBus)
				handler := identity_commands.NewLoginCommandHandler(onboarding, repo, ts, session, eventBus)
				return handler, repo, ts, session
			},
			command: identity_commands.LoginCommand{
				Email:    "test@example.com",
				Password: "secret",
			},
			expectErr: true,
		},
		{
			name: "token persist fails",
			setup: func() (*identity_commands.LoginCommandHandler, *fakeIdentityRepo, *fakeTokenService, *fakeSessionManager) {
				onboarding := &fakeOnboardingRepo{user: baseOnboardingUser}
				repo := &fakeIdentityRepo{user: baseIdentityUser}
				session := &fakeSessionManager{}
				ts := &fakeTokenService{persistErr: errors.New("boom")}
				eventBus := newFakeEventBus()
				eventBus.SetEventBus(eventBus)
				handler := identity_commands.NewLoginCommandHandler(onboarding, repo, ts, session, eventBus)
				return handler, repo, ts, session
			},
			command: identity_commands.LoginCommand{
				Email:    "test@example.com",
				Password: "secret",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, repo, ts, sm := tt.setup()

			res, err := handler.Handle(tt.command)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.assertions != nil {
				tt.assertions(t, res, repo, ts, sm)
			}
		})
	}
}
