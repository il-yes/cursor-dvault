// internal/vault/application/session/session.go
package vault_session

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
	"vault-app/internal/models"
	"vault-app/internal/tracecore"
	vaults_domain "vault-app/internal/vault/domain"

	"gorm.io/gorm"
)

type SessionState string

const (
	SessionPrepared  SessionState = "prepared"
	SessionVaultOpen SessionState = "vault_open"
)

type SessionV0 struct {
	ID      string
	UserID  string
	Vault   *vaults_domain.VaultPayload
	Dirty   bool
	LastCID string

	LastSynced     string
	LastUpdated    string
	Runtime        *RuntimeContext
	PendingCommits []tracecore.CommitEnvelope `json:"pending_commits,omitempty"`
}
type Session struct {
	UserID      string `gorm:"uniqueIndex"`
	Vault       []byte `json:"vault_blob,omitempty"`
	LastCID     string
	LastSynced  string
	LastUpdated string
	Runtime     *RuntimeContext `json:"-"`
	Dirty       bool
	PendingCommits []tracecore.CommitEnvelope `json:"pending_commits,omitempty"`
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

// VaultSession holds the decrypted vault during an active session
type VaultSession struct {
	UserID              string                     `json:"user_id"`
	Vault               vaults_domain.VaultPayload // Decrypted vault format
	LastCID             string
	Dirty               bool
	LastSynced          string
	LastUpdated         string
	Mutex               sync.Mutex                 `json:"-"`
	VaultRuntimeContext RuntimeContext             `json:"vault_runtime_context"`
	PendingCommits      []tracecore.CommitEnvelope `json:"pending_commits,omitempty"`
}

func (v *VaultSession) ToFormerModel() *models.VaultSession {
	return &models.VaultSession{
		UserID:              v.UserID,
		Vault:               v.Vault.ToFormerVaultPayload(),
		LastCID:             v.LastCID,
		Dirty:               v.Dirty,
		LastSynced:          v.LastSynced,
		LastUpdated:         v.LastUpdated,
		VaultRuntimeContext: *v.VaultRuntimeContext.ToFormerRuntimeContext(),
		PendingCommits:      v.PendingCommits,
	}
}

type UserSession struct {
	UserID      string `gorm:"primaryKey;column:user_id" json:"user_id"`
	SessionData string `gorm:"type:json" json:"session_data"` // Marshaled VaultSession
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

// --------- Refactoring ---------

// --- NEW STRUCT ---
type UserJSON struct {
	UserID    string `gorm:"primaryKey"`
	Data      []byte `gorm:"type:bytea"` // raw JSON storage
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserJSON) TableName() string {
	return "user_json_sessions"
}

// --- REPOSITORY ---
type JSONSessionRepository struct {
	db *gorm.DB
}

func NewJSONSessionRepository(db *gorm.DB) *JSONSessionRepository {
	return &JSONSessionRepository{db: db}
}

// Save or update session (upsert)
func (r *JSONSessionRepository) SaveSession(userID string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var session UserJSON
	err = r.db.First(&session, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// create new
		session = UserJSON{
			UserID:    userID,
			Data:      data,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return r.db.Create(&session).Error
	} else if err != nil {
		return err
	}

	// update existing
	session.Data = data
	session.UpdatedAt = time.Now()
	return r.db.Save(&session).Error
}

// Get session by userID
func (r *JSONSessionRepository) GetSession(userID string, out interface{}) error {
	var session UserJSON
	if err := r.db.First(&session, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return json.Unmarshal(session.Data, out)
}

// Delete session
func (r *JSONSessionRepository) DeleteSession(userID string) error {
	return r.db.Delete(&UserJSON{}, "user_id = ?", userID).Error
}
