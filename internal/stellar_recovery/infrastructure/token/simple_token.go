package token


import "fmt"

type SimpleTokenGen struct{}

func NewSimpleTokenGen() *SimpleTokenGen { return &SimpleTokenGen{} }

func (t *SimpleTokenGen) NewSessionToken(userID string) (string, error) {
	// Replace with JWT or other real implementation
	return fmt.Sprintf("session_%s", userID), nil
}
