package vault_infrastructure_security

import (
	"crypto/rand"
	"encoding/hex"
)

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}