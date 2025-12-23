package identity_domain

import "time"

// Domain event
type UserRegistered struct {
	UserID    string
	IsAnonymous bool
	OccurredAt int64
}

func NewUserRegistered(u *User) UserRegistered {
	return UserRegistered{UserID: u.ID, IsAnonymous: u.IsAnonymous, OccurredAt: time.Now().UnixNano()}
}
