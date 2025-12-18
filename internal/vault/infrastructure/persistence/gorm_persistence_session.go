package vaults_persistence

import (
	"encoding/json"
	"fmt"
	vault_session "vault-app/internal/vault/application/session"

	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model
	UserID      string `json:"user_id"`
	SessionData string `json:"session_data"`
}
type GormSessionRepository struct {
	db *gorm.DB
}

func NewGormSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{db: db}
}

func (r *GormSessionRepository) CreateSession(session *vault_session.Session) error {
	sessionMapper := SessionDomainToMapper(session)
	return r.db.Create(&sessionMapper).Error
}

func (r *GormSessionRepository) GetSession(sessionID string) (*vault_session.Session, error) {
	var session SessionMapper
	if err := r.db.First(&session, sessionID).Error; err != nil {
		return nil, err
	}
	return session.ToDomain(), nil
}

func (r *GormSessionRepository) SaveSession(userID string, session *vault_session.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	userSession := UserSession{
		UserID:      userID,
		SessionData: string(data),
	}
	return r.db.Save(&userSession).Error
}

func (r *GormSessionRepository) DeleteSession(sessionID string) error {
	return r.db.Delete(&UserSession{}, sessionID).Error
}

func (r *GormSessionRepository) GetLatestByUserID(userID string) (*vault_session.Session, error) {
	var session SessionMapper
	if err := r.db.Last(&session, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return session.ToDomain(), nil
}
func (r *GormSessionRepository) UpdateSession(session *vault_session.Session) error {
	return r.db.Save(session).Error
}	
