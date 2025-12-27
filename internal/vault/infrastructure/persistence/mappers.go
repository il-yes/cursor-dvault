package vaults_persistence

import (
	"encoding/json"
	"fmt"
	"log"
	"vault-app/internal/blockchain"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)


type VaultMapper struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"column:name"`
	Type      string `json:"type" gorm:"column:type"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`
	CID       string `json:"cid" gorm:"column:cid"` // ✅ Explicitly map this!
	TxHash    string `json:"tx_hash" gorm:"column:tx_hash"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
}

func (vm *VaultMapper) TableName() string {
	return "vault_vaults"
}
func (vm *VaultMapper) ToDomain() *vaults_domain.Vault {
	return &vaults_domain.Vault{
		ID:        vm.ID,
		Name:      vm.Name,
		Type:      vm.Type,
		UserID:    vm.UserID,
		CID:       vm.CID,
		TxHash:    vm.TxHash,
		CreatedAt: vm.CreatedAt,
		UpdatedAt: vm.UpdatedAt,
	}
}
func VaultDomainToMapper(vault *vaults_domain.Vault) *VaultMapper {
	return &VaultMapper{
		ID:        vault.ID,
		Name:      vault.Name,
		Type:      vault.Type,
		UserID:    vault.UserID,
		CID:       vault.CID,
		TxHash:    vault.TxHash,
		CreatedAt: vault.CreatedAt,
		UpdatedAt: vault.UpdatedAt,
	}
}

type SessionMapper0 struct {
    UserID   string `json:"user_id"`
    Vault    string `json:"vault" gorm:"type:longtext"`
    Dirty    bool   `json:"dirty"`
    LastCID  string `json:"last_cid"`

    LastSynced  string `json:"last_synced"`
    LastUpdated string `json:"last_updated"`
}
type SessionMapper struct {
    UserID      string `gorm:"uniqueIndex"`
    Vault       []byte `gorm:"type:bytea"` // store encrypted VaultPayload + other session info
	LastCID     string
	LastSynced  string
	LastUpdated string

	Dirty       bool
}

func (sm *SessionMapper) TableName() string {
    return "vault_sessions"
}
func (m *SessionMapper) ToDomain() (*vault_session.Session, error) {
    cryptoS := blockchain.CryptoService{}
    decrypted, err := cryptoS.Decrypt(m.Vault, "password")
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt session for user %s: %w", m.UserID, err)
    }

    session := &vault_session.Session{}
    if err := json.Unmarshal(decrypted, session); err != nil {
        return nil, fmt.Errorf("failed to decode session for user %s: %w", m.UserID, err)
    }
    return session, nil
}



func SessionDomainToMapper(session *vault_session.Session) *SessionMapper {
    // Marshal the full session domain into JSON for storage
    data, err := json.Marshal(session)
    if err != nil {
        // If you want, you can panic here or handle the error upstream
        log.Printf("SessionDomainToMapper: failed to marshal session: %v", err)
        data = []byte("{}") // fallback to empty JSON
    }

    return &SessionMapper{
        UserID:      session.UserID,
        Vault: data, // store the full session JSON
    }
}



