package tracecore

import (
	"errors"
	"fmt"
	"net/http"
)

// Standardized Cloud & Identity Errors
var (
	ErrCloudUnauthorized      = errors.New("cloud authentication required (401)")
	ErrVaultForbidden         = errors.New("vault access forbidden or delegation required (403)")
	ErrResourceNotFound       = errors.New("resource not found (404)")
	ErrChallengeExpired       = errors.New("vault challenge expired")
	ErrInvalidSignature       = errors.New("invalid cryptographic signature for challenge")
	ErrRegistrationBadRequest = errors.New("malformed vault registration request (400)")
	ErrCloudServerError       = errors.New("cloud server error (5xx)")
)

// MapHTTPStatusToError converts an HTTP status code and response body into a typed error
func MapHTTPStatusToError(statusCode int, body string) error {
	switch statusCode {
	case http.StatusBadRequest:
		if body != "" {
			return fmt.Errorf("%w: %s", ErrRegistrationBadRequest, body)
		}
		return ErrRegistrationBadRequest
	case http.StatusUnauthorized:
		return ErrCloudUnauthorized
	case http.StatusForbidden:
		return ErrVaultForbidden
	case http.StatusNotFound:
		return ErrResourceNotFound
	default:
		if statusCode >= 500 {
			return fmt.Errorf("%w (status %d): %s", ErrCloudServerError, statusCode, body)
		}
		return fmt.Errorf("cloud returned status %d: %s", statusCode, body)
	}
}
