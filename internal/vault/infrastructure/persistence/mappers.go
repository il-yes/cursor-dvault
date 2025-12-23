package vaults_persistence

import (
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

type SessionMapper struct {
    UserID   string
    Vault    *vaults_domain.VaultPayload `gorm:"-"`
    Dirty    bool
    LastCID  string

    LastSynced string
    LastUpdated string
    Runtime *vault_session.RuntimeContext `gorm:"-"`
}

func (sm *SessionMapper) TableName() string {
    return "vault_sessions"
}
func (sm *SessionMapper) ToDomain() *vault_session.Session {
    return &vault_session.Session{
        UserID:   sm.UserID,
        Vault:    sm.Vault,
        Dirty:    sm.Dirty,
        LastCID:  sm.LastCID,
        LastSynced: sm.LastSynced,
        LastUpdated: sm.LastUpdated,
        Runtime: sm.Runtime,
    }
}
func SessionDomainToMapper(session *vault_session.Session) *SessionMapper {
    return &SessionMapper{
        UserID:   session.UserID,
        Vault:    session.Vault,
        Dirty:    session.Dirty,
        LastCID:  session.LastCID,
        LastSynced: session.LastSynced,
        LastUpdated: session.LastUpdated,
        Runtime: session.Runtime,
    }
}
