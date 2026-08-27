package vaults_persistence

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"vault-app/internal/blockchain"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
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

func createSessionWriteTrace(op string, userID string, session *vault_session.Session) map[string]interface{} {
	if session == nil {
		return map[string]interface{}{
			"operation": op,
			"userID":    userID,
			"session":   "nil",
		}
	}

	vaultPresent := session.Vault != nil
	vaultByteLen := 0
	attCount := 0
	noteCount := 0
	entriesByType := map[string]int{}
	entriesWithAttCIDs := make(map[string][]string)
	attNodeCIDs := []string{}

	if vaultPresent {
		vaultByteLen = len(session.Vault)
		parsed := vaults_domain.ParseVaultPayload(session.Vault)

		attCount = len(parsed.Attachments) + len(parsed.Personal.Attachments)
		for _, att := range parsed.Attachments {
			if att.NodeCID != "" {
				attNodeCIDs = append(attNodeCIDs, att.NodeCID)
			}
		}
		for _, att := range parsed.Personal.Attachments {
			if att.NodeCID != "" {
				attNodeCIDs = append(attNodeCIDs, att.NodeCID)
			}
		}

		entriesByType["Login"] = len(parsed.Entries.Login) + len(parsed.Personal.Entries.Login)
		entriesByType["Card"] = len(parsed.Entries.Card) + len(parsed.Personal.Entries.Card)
		entriesByType["Identity"] = len(parsed.Entries.Identity) + len(parsed.Personal.Entries.Identity)
		entriesByType["Note"] = len(parsed.Entries.Note) + len(parsed.Personal.Entries.Note)
		entriesByType["SSHKey"] = len(parsed.Entries.SSHKey) + len(parsed.Personal.Entries.SSHKey)
		noteCount = entriesByType["Note"]

		checkNotes := func(entries []vaults_domain.NoteEntry) {
			for _, e := range entries {
				if len(e.AttachmentCIDs) > 0 {
					entriesWithAttCIDs["note:"+e.ID] = e.AttachmentCIDs
				}
			}
		}
		checkLogins := func(entries []vaults_domain.LoginEntry) {
			for _, e := range entries {
				if len(e.AttachmentCIDs) > 0 {
					entriesWithAttCIDs["login:"+e.ID] = e.AttachmentCIDs
				}
			}
		}

		checkNotes(parsed.Entries.Note)
		checkNotes(parsed.Personal.Entries.Note)
		checkLogins(parsed.Entries.Login)
		checkLogins(parsed.Personal.Entries.Login)
	}

	return map[string]interface{}{
		"operation":                 op,
		"userID":                    userID,
		"lastCID":                   session.LastCID,
		"vaultPresent":              vaultPresent,
		"vaultByteLength":           vaultByteLen,
		"attachmentCount":           attCount,
		"noteCount":                 noteCount,
		"entriesByType":             entriesByType,
		"entriesWithAttachmentCIDs": entriesWithAttCIDs,
		"attachmentNodeCIDs":        attNodeCIDs,
	}
}

func (r *GormSessionRepository) SaveSession(userID string, session *vault_session.Session) error {
	log.Printf("SESSION WRITE TRACE (BEFORE SaveSession): %+v", createSessionWriteTrace("SaveSession", userID, session))

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
	err = r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(&mapper).Error

	if err == nil {
		log.Printf("SESSION WRITE TRACE (AFTER SaveSession SUCCESS): userID=%s lastCID=%s", userID, session.LastCID)
	} else {
		log.Printf("SESSION WRITE TRACE (AFTER SaveSession ERROR): userID=%s err=%v", userID, err)
	}
	return err
}

func (r *GormSessionRepository) GetSession(userID string) (*vault_session.Session, error) {
	var mapper SessionMapper
	err := r.db.Where("user_id = ?", userID).First(&mapper).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	if session.Vault != nil {
		h := sha256.Sum256(session.Vault)
		log.Printf("[BYTE IDENTITY TRACE 1 - GetSession DB READ] userID=%s vaultBytesLen=%d vaultSHA256=%x", userID, len(session.Vault), h[:])
	}
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
func (r *GormSessionRepository) GetEntries(s vault_session.Session) (*vaults_domain.Entries, error) {
	vp, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return nil, err
	}
	return &vp.Entries, nil
}
func (r *GormSessionRepository) GetFolders(s vault_session.Session) ([]vaults_domain.Folder, error) {
	vp, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return nil, err
	}
	return vp.Folders, nil
}
func (r *GormSessionRepository) GetAttachements(s vault_session.Session) ([]vaults_domain.Attachment, error) {
	vp, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return nil, err
	}
	return vp.Attachments, nil
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
		session, err := r.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to restore session for user %s: %w", r.UserID, err)
		}
		runtimePresent := session.Runtime != nil
		secretsPresent := runtimePresent && session.Runtime.SessionSecrets != nil
		cloudJWTLen := 0
		if secretsPresent {
			if v, ok := session.Runtime.SessionSecrets["cloud_jwt"]; ok {
				cloudJWTLen = len(v)
			}
		}
		log.Printf("[TOKEN-TRACE 1] FindAll user=%s runtime_present=%v secrets_present=%v cloud_jwt_length=%d", r.UserID, runtimePresent, secretsPresent, cloudJWTLen)
		sessions[r.UserID] = session
	}

	return sessions, nil
}
