package identity_usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	identity_eventbus "vault-app/internal/identity/application"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
)

//
// -------- FAKES --------
//

type fakeRepoHandler struct {
	user       *identity_domain.User
	saveCalled bool
	saveErr    error
}

func (f *fakeRepoHandler) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	return f.user, nil
}
func (f *fakeRepoHandler) FindByPublicKey(ctx context.Context, publicKey string) (*identity_domain.User, error) {
	return f.user, nil
}
func (f *fakeRepoHandler) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	return f.user, nil
}

func (f *fakeRepoHandler) Save(ctx context.Context, u *identity_domain.User) error {
	f.saveCalled = true
	f.user = u
	return f.saveErr
}
func (f *fakeRepoHandler) Update(ctx context.Context, u *identity_domain.User) error {
	f.user = u
	return nil
}

type fakeBusHandler struct {
	published bool
	event     identity_eventbus.UserLoggedIn
}

func (f *fakeBusHandler) PublishUserLoggedIn(ctx context.Context, e identity_eventbus.UserLoggedIn) error {
	f.published = true
	f.event = e
	return nil
}
func (f *fakeBusHandler) PublishUserRegistered(ctx context.Context, e identity_eventbus.UserRegistered) error {
	return nil
}
func (f *fakeBusHandler) SubscribeToUserRegistered(handler identity_eventbus.UserRegisteredHandler) error {
	return nil
}
func (f *fakeBusHandler) SubscribeToUserLoggedIn(handler identity_eventbus.UserLoggedInHandler) error {
	return nil
}
//
// -------- TESTS --------
//

func TestLoginUseCase_Execute(t *testing.T) {
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	existingUser := &identity_domain.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: string(hashedPwd),
	}

	idGen := func() string { return "new-id-456" }

	tests := []struct {
		name         string
		existingUser *identity_domain.User
		email        string
		password     string
		publicKey    string
		signedMsg    string
		signature    string
		expectErr    bool
		assertFunc   func(t *testing.T, user *identity_domain.User, repo *fakeRepoHandler, bus *fakeBusHandler)
	}{
		{
			name:         "successful login",
			existingUser: existingUser,
			email:        "test@example.com",
			password:     "secret",
			expectErr:    false,
			assertFunc: func(t *testing.T, user *identity_domain.User, repo *fakeRepoHandler, bus *fakeBusHandler) {
				assert.Equal(t, "user-123", user.ID)
				assert.True(t, repo.saveCalled)
				assert.True(t, bus.published)
				assert.Equal(t, user.ID, bus.event.UserID)
				assert.Equal(t, user.Email, bus.event.Email)
				assert.NotEmpty(t, user.LastConnectedAt)
			},
		},
		{
			name:         "invalid password",
			existingUser: existingUser,
			email:        "test@example.com",
			password:     "wrong",
			expectErr:    true,
		},
		{
			name:         "nonexistent email",
			existingUser: nil,
			email:        "missing@example.com",
			password:     "secret",
			expectErr:    true,
		},
		{
			name:         "public key login missing signature",
			existingUser: existingUser,
			email:        "test@example.com",
			publicKey:    "stellar-key",
			signedMsg:    "",
			signature:    "",
			expectErr:    true,
		},
		{
			name:         "public key login success (placeholder password)",
			existingUser: nil,
			email:        "new@example.com",
			publicKey:    "stellar-key",
			signedMsg:    "signed-msg",
			signature:    "sig",
			expectErr:    false,
			assertFunc: func(t *testing.T, user *identity_domain.User, repoHandler *fakeRepoHandler, bus *fakeBusHandler) {
				assert.Equal(t, "new-id-456", user.ID)
				assert.True(t, repoHandler.saveCalled)
				assert.True(t, bus.published)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoHandler := &fakeRepoHandler{user: tt.existingUser}
			busHandler := &fakeBusHandler{}
			uc := identity_usecase.NewLoginUseCase(repoHandler, busHandler, idGen)

			user, err := uc.Execute(context.Background(), tt.email, tt.password, tt.publicKey, tt.signedMsg, tt.signature)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.assertFunc != nil {
				tt.assertFunc(t, user, repoHandler, busHandler)
			}
		})
	}
}

//
// -------- TIMING ATTACK TEST (simplified) --------
//

func TestLoginUseCase_TimingAttack(t *testing.T) {
	repoHandler := &fakeRepoHandler{user: nil} // user does not exist
	busHandler := &fakeBusHandler{}
	uc := identity_usecase.NewLoginUseCase(repoHandler, busHandler, func() string { return "id" })

	start := time.Now()
	_, _ = uc.Execute(context.Background(), "missing@example.com", "wrongpass", "", "", "")
	duration1 := time.Since(start)

	start = time.Now()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	repoHandler.user = &identity_domain.User{Email: "exists@example.com", PasswordHash: string(hashed)}
	_, _ = uc.Execute(context.Background(), "exists@example.com", "wrongpass", "", "", "")
	duration2 := time.Since(start)

	// Duration should be roughly similar to prevent timing attacks
	diff := duration1 - duration2
	if diff < 0 {
		diff = -diff
	}
	assert.Less(t, diff.Seconds(), 0.01, "timing difference should be small")
}
