package vaults_domain

type VaultStorage interface {
	GetData(cid string) ([]byte, error)
}

type CryptoService interface {
	Decrypt(data []byte, password string) ([]byte, error)
}
