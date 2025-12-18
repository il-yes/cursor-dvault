// internal/vault/application/session/session.go
package vault_session

import (
	"sync"
	"vault-app/internal/tracecore"
	vaults_domain "vault-app/internal/vault/domain"

	"gorm.io/gorm"
)


type Session struct {
    UserID   string
    Vault    *vaults_domain.VaultPayload
    Dirty    bool
    LastCID  string

    LastSynced string
    LastUpdated string
    Runtime *RuntimeContext
}


type SessionRepository interface {
	CreateSession(session *Session) error
	GetSession(sessionID string) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSession(sessionID string) error
	SaveSession(userID string, session *Session) error


	GetLatestByUserID(userID string) (*Session, error)
}


// VaultSession holds the decrypted vault during an active session
type VaultSession struct {
	UserID              string           `json:"user_id"`
	Vault               *vaults_domain.VaultPayload // Decrypted vault format
	LastCID             string
	Dirty               bool
	LastSynced          string
	LastUpdated         string
	Mutex               sync.Mutex                 `json:"-"`
	VaultRuntimeContext RuntimeContext `json:"vault_runtime_context"`
	PendingCommits      []tracecore.CommitEnvelope `json:"pending_commits,omitempty"`
}
type UserSession struct {
	UserID      string            `gorm:"primaryKey;column:user_id" json:"user_id"`
	SessionData string         `gorm:"type:json" json:"session_data"` // Marshaled VaultSession
	UpdatedAt   string         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
} 

func (UserSession) TableName() string {
	return "user_sessions"
}
 