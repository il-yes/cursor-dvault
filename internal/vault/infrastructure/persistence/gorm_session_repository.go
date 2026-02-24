package vaults_persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	utils "vault-app/internal/utils"
	"vault-app/internal/blockchain"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormSessionRepository struct {
	db *gorm.DB
}

func NewGormSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{db: db}
}

func (r *GormSessionRepository) CreateSession(session *vault_session.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	mapper := &SessionMapper{
		UserID: session.UserID,
		Vault:  data,
	}
	return r.db.Create(mapper).Error
}

func (r *GormSessionRepository) SaveSession(userID string, session *vault_session.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	cryptoS := blockchain.CryptoService{}
	encryptedData, err := cryptoS.Encrypt(data, "password")
	if err != nil {
		return err
	}

	mapper := SessionMapper{
		UserID:      session.UserID,
		Vault:       encryptedData,
		LastCID:     session.LastCID,
		LastSynced:  session.LastSynced,
		LastUpdated: time.Now().Format(time.RFC3339),
	}
	_, err = vault_session.DecodeSessionVault(session.Vault)
	if err != nil {
		return err
	}

	// Upsert: if exists, update; else insert
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(&mapper).Error
}

func (r *GormSessionRepository) GetSession(userID string) (*vault_session.Session, error) {
	var mapper SessionMapper
	err := r.db.Where("user_id = ?", userID).First(&mapper).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogPretty("GormSessionRepository - GetSession - mapper not found", mapper)
			// create a default session if not found
			// session := vault_session.InitNewSession(userID)
			// if err := r.SaveSession(userID, session); err != nil {
			//     return nil, fmt.Errorf("failed to create default session: %w", err)
			// }
			// return session, nil
		}
		return nil, err
	}

	cryptoS := blockchain.CryptoService{}
	decrypted, err := cryptoS.Decrypt(mapper.Vault, "password")
	if err != nil {
		return nil, err
	}

	var session vault_session.Session
	if err := json.Unmarshal(decrypted, &session); err != nil {
		return nil, err
	}
	utils.LogPretty("GormSessionRepository - GetSession - mapper", mapper)
	utils.LogPretty("GormSessionRepository - GetSession - session", session)
	utils.LogPretty("GormSessionRepository - GetSession - decrypted", vaults_domain.ParseVaultPayload(decrypted))
	return &session, nil
}

func (r *GormSessionRepository) DeleteSession(sessionID string) error {
	return r.db.Delete(&SessionMapper{}, sessionID).Error
}

func (r *GormSessionRepository) GetLatestByUserID(userID string) (*vault_session.Session, error) {
	var session SessionMapper
	if err := r.db.Last(&session, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return session.ToDomain()
}
func (r *GormSessionRepository) UpdateSession(session *vault_session.Session) error {
	return r.SaveSession(session.UserID, session)
}

type SessionDBModel struct {
	db *gorm.DB
}

func NewSessionDBModel(db *gorm.DB) *SessionDBModel {
	return &SessionDBModel{db: db}
}
func (db *SessionDBModel) FindAll() (map[string]*vault_session.Session, error) {
	var records []SessionMapper
	if err := db.db.Find(&records).Error; err != nil {
		return nil, err
	}

	sessions := make(map[string]*vault_session.Session)

	for _, r := range records {
		session := &vault_session.Session{
			UserID:      r.UserID,
			Vault:       r.Vault, // ✅ raw encrypted bytes
			LastCID:     r.LastCID,
			LastSynced:  r.LastSynced,
			LastUpdated: r.LastUpdated,
			Dirty:       r.Dirty,
		}

		sessions[r.UserID] = session
	}

	return sessions, nil
}

