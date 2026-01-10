// internal/vault/application/session/session.go
package vault_session

import (
	"encoding/json"
	"vault-app/internal/tracecore"
	vaults_domain "vault-app/internal/vault/domain"

)

type SessionState string

const (
	SessionPrepared  SessionState = "prepared"
	SessionVaultOpen SessionState = "vault_open"
)

type Session struct {
	UserID      string `gorm:"uniqueIndex"`
	Vault       []byte `json:"vault_blob,omitempty"`
	LastCID     string
	LastSynced  string
	LastUpdated string
	Runtime     *RuntimeContext `json:"-" gorm:"-"`
	Dirty       bool
	PendingCommits []tracecore.CommitEnvelope `json:"pending_commits,omitempty" gorm:"-"`
}

func InitNewSession(userID string) *Session {
	return &Session{
		UserID:      userID,
		Vault:       nil,
		LastCID:     "",
		LastSynced:  "",
		LastUpdated: "",
		Runtime:     nil,
		Dirty:       false,
		PendingCommits: []tracecore.CommitEnvelope{},
	}
}

func DecodeSessionVault(blob []byte) (*vaults_domain.VaultPayload, error) {
	if len(blob) == 0 {
		return nil, nil
	}

	var payload vaults_domain.VaultPayload
	if err := json.Unmarshal(blob, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

type SessionRepository interface {
	CreateSession(session *Session) error
	GetSession(sessionID string) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSession(sessionID string) error
	SaveSession(userID string, session *Session) error
	GetLatestByUserID(userID string) (*Session, error)
}


