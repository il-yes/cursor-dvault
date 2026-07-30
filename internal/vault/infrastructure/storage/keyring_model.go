package vaults_storage

type WrappedKeyring struct {
	Type       string `json:"type"` // "password" | "stellar"
	Ciphertext []byte `json:"ciphertext"`
}


type StoredKeyring struct {
	VaultID string `json:"vault_id"`
	Ciphertext []byte `json:"ciphertext"`

	Wrappers []WrappedKeyring `json:"wrappers"`
	// encrypted blob of VaultKeyring (domain)

	// how it was encrypted
	KDF        string `json:"kdf"` // "password" | "stellar"
	Version    int    `json:"version"`
}