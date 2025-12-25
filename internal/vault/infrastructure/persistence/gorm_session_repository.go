package vaults_persistence

import (
	"encoding/json"
	"fmt"
	utils "vault-app/internal"
	vault_session "vault-app/internal/vault/application/session"

	"gorm.io/gorm"
)


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
	if err := r.db.Where("user_id = ?", sessionID).First(&session).Error; err != nil {
		return nil, err
	}
	return session.ToDomain(), nil
}

func (r *GormSessionRepository) SaveSession(userID string, session *vault_session.Session) error {
	utils.LogPretty("GormSessionRepository - SaveSession", session)
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	r.db = r.db.Debug()


	userSession := vault_session.UserSession{
		UserID:      userID,
		SessionData: string(data),
	}
	return r.db.Save(&userSession).Error
}

func (r *GormSessionRepository) DeleteSession(sessionID string) error {
	return r.db.Delete(&SessionMapper{}, sessionID).Error
}

func (r *GormSessionRepository) GetLatestByUserID(userID string) (*vault_session.Session, error) {
	var session SessionMapper
	if err := r.db.Last(&session, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return session.ToDomain(), nil
}
func (r *GormSessionRepository) UpdateSession(session *vault_session.Session) error {
	return r.db.Save(SessionDomainToMapper(session)).Error
}	


type DBModel struct {
	db *gorm.DB
}

func NewDBModel(db *gorm.DB) *DBModel {
	return &DBModel{db: db}
}	

func (db *DBModel) SaveSessionV1(userID string, session *vault_session.Session) error {
	utils.LogPretty("DBModel - SaveSessionV1", session)
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	userSession := vault_session.UserSession{
		UserID:      userID,
		SessionData: string(data),
	}
	return db.db.Save(&userSession).Error
}

func (db *DBModel) LoadSessionV1(userID string) (*vault_session.Session, error) {
	var userSession vault_session.UserSession
	if err := db.db.Model(&vault_session.UserSession{}).Where("user_id = ?", userID).First(&userSession).Error; err != nil {
		return nil, err
	}

	var session vault_session.Session
	if err := json.Unmarshal([]byte(userSession.SessionData), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &session, nil
}

func (db *DBModel) GetAllSessionsV1() (map[string]*vault_session.Session, error) {
	var userSessions []vault_session.UserSession
	if err := db.db.Model(&vault_session.UserSession{}).Find(&userSessions).Error; err != nil {
		return nil, err
	}

	sessionMap := make(map[string]*vault_session.Session)
	for _, s := range userSessions {
		var session vault_session.Session
		s.UserID = fmt.Sprintf("%s", s.UserID)
		
		if	 err := json.Unmarshal([]byte(s.SessionData), &session); err != nil {
			return nil, fmt.Errorf("failed to decode session for user %s: %w", s.UserID, err)
		}
		sessionMap[s.UserID] = &session
	}
	
	return sessionMap, nil
}
