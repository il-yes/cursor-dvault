package onboarding_domain

import "errors"



var (
    ErrUserExists =     errors.New("user_exists")
	ErrUserNotFound =     errors.New("user_not_found")
)	