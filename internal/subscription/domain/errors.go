package subscription_domain

import "errors"

var (
    ErrInvalidSubscription = errors.New("invalid subscription")
    ErrSubscriptionNotFound           = errors.New("subscription not found")
)
    
var (
	ErrUserNotFound     = errors.New("User: not found")
	ErrUserInvalidState = errors.New("User: invalid state")
	ErrUserValidation   = errors.New("User: validation failed")
)