package identity_usecase_test

import (
	"context"
	"errors"
	"testing"

	identity_app "vault-app/internal/identity/application"
	identity_domain "vault-app/internal/identity/domain"
	identity_uc "vault-app/internal/identity/application/usecase"
)

/* ------------------------------
   Fake In-Memory Repo
--------------------------------*/
type fakeRepo struct {
	usersByID    map[string]*identity_domain.User
	usersByEmail map[string]*identity_domain.User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		usersByID:    map[string]*identity_domain.User{},
		usersByEmail: map[string]*identity_domain.User{},
	}
}

func (r *fakeRepo) Save(ctx context.Context, u *identity_domain.User) error {
	r.usersByID[u.ID] = u
	if u.Email != "" {
		r.usersByEmail[u.Email] = u
	}
	return nil
}
func (r *fakeRepo) Update(ctx context.Context, u *identity_domain.User) error {
	return nil
}
func (r *fakeRepo) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	return r.usersByID[id], nil
}

func (r *fakeRepo) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	return r.usersByEmail[email], nil
}

func (r *fakeRepo) FindByPublicKey(ctx context.Context, publicKey string) (*identity_domain.User, error) {
	return r.usersByEmail[publicKey], nil
}

/* ------------------------------
   Fake Event Bus
--------------------------------*/
type fakeEventBus struct {
	published []identity_app.UserRegistered
	publishedLogins []identity_app.UserLoggedIn
}

func newFakeEventBus() *fakeEventBus {
	return &fakeEventBus{
		published: []identity_app.UserRegistered{},
		publishedLogins: []identity_app.UserLoggedIn{},
	}
}

func (b *fakeEventBus) PublishUserRegistered(ctx context.Context, e identity_app.UserRegistered) error {
	b.published = append(b.published, e)
	return nil
}

func (b *fakeEventBus) SubscribeToUserRegistered(handler identity_app.UserRegisteredHandler) error {
	return nil
}

func (b *fakeEventBus) PublishUserLoggedIn(ctx context.Context, e identity_app.UserLoggedIn) error {
	return nil
}

func (b *fakeEventBus) SubscribeToUserLoggedIn(handler identity_app.UserLoggedInHandler) error {
	return nil
}

func (b *fakeEventBus) PublishDeviceCreated(ctx context.Context, e identity_app.DeviceCreated) error {
	return nil
}
func (b *fakeEventBus) SubscribeToDeviceCreated(handler identity_app.DeviceCreatedHandler) error {
	return nil
}
func (b *fakeEventBus) PublishDeviceRevoked(ctx context.Context, e identity_app.DeviceRevoked) error {
	return nil
}
func (b *fakeEventBus) SubscribeToDeviceRevoked(handler identity_app.DeviceRevokedHandler) error {
	return nil
}


/* ------------------------------
   Fake ID Generator
--------------------------------*/
func fakeIDGen() string {
	return "test-user-id"
}

/* ------------------------------
   TESTS
--------------------------------*/

func TestRegisterStandardUser_Success(t *testing.T) {
	t.Log("➡️  START TestRegisterStandardUser_Success")
	ctx := context.Background()
	repo := newFakeRepo()
	bus := newFakeEventBus()

	t.Log("🔧 Creating use case")
	ucase := identity_uc.NewRegisterStandardUserUseCase(repo, bus, fakeIDGen)

	email := "test@example.com"
	hash := "hashed-password"

	t.Log("📩 Executing use case")
	u, err := ucase.Execute(ctx, email, hash)
	if err != nil {
		t.Fatalf("❌ expected no error, got %v", err)
	}

	t.Log("✅ Use case returned user:", u)

	// Validate saved user
	t.Log("🔍 Checking saved user fields")
	if u.ID != "test-user-id" {
		t.Errorf("❌ expected ID 'test-user-id', got %s", u.ID)
	}
	if u.Email != email {
		t.Errorf("❌ expected email %s, got %s", email, u.Email)
	}
	if !u.IsStandard() {
		t.Errorf("❌ expected standard user, got anonymous")
	}

	// Validate event
	t.Log("🔍 Checking published events")
	if len(bus.published) != 1 {
		t.Fatalf("❌ expected 1 user-registered event, got %d", len(bus.published))
	}

	ev := bus.published[0]
	t.Logf("📦 Event received: %+v", ev)

	if ev.UserID != "test-user-id" {
		t.Errorf("❌ expected event userID 'test-user-id', got %s", ev.UserID)
	}
	if ev.IsAnonymous {
		t.Errorf("❌ expected IsAnonymous=false in event")
	}
	if ev.OccurredAt == 0 {
		t.Errorf("❌ expected OccurredAt timestamp to be set")
	}

	t.Log("🎉 TestRegisterStandardUser_Success PASSED")
}

func TestRegisterStandardUser_DuplicateEmail(t *testing.T) {
	t.Log("➡️  START TestRegisterStandardUser_DuplicateEmail")

	ctx := context.Background()
	repo := newFakeRepo()
	bus := newFakeEventBus()
	ucase := identity_uc.NewRegisterStandardUserUseCase(repo, bus, fakeIDGen)

	t.Log("🔧 Creating existing user in repo")
	repo.Save(ctx, identity_domain.NewStandardUser("existing", "dup@example.com", "hash"))

	t.Log("📩 Executing use case with duplicate email")
	_, err := ucase.Execute(ctx, "dup@example.com", "hash")
	if !errors.Is(err, identity_domain.ErrUserExists) {
		t.Fatalf("❌ expected ErrUserExists, got %v", err)
	}

	t.Log("🔍 Verifying no events were published")
	if len(bus.published) != 0 {
		t.Fatalf("❌ expected 0 events, got %d", len(bus.published))
	}

	t.Log("🎉 TestRegisterStandardUser_DuplicateEmail PASSED")
}
