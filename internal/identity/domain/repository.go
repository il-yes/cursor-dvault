package identity_domain	

import "context"

// Repository interface for identity
type UserRepository interface {
	Save(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
}
