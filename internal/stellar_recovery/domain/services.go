package stellar_recovery_domain

type StellarKeyValidator interface {
	ParseSecret(secret string) (publicKey string, err error)
	VerifyOwnership(publicKey string, secret string) error
}

type EventLogger interface {
	LogVaultRecovered(userID, vaultID, publicKey string) error
}

type TokenGenerator interface {
	NewSessionToken(userID string) (string, error)
}

type EventDispatcher interface {
	Dispatch(event any) error
}

type StellarKeyVerifier interface {
	ParseSecret(secret string) (publicKey string, err error)
	VerifyOwnership(secretKey string, expectedPublicKey string) error
}

