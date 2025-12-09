package models

import (
	// "database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"vault-app/internal/auth"
	app_config "vault-app/internal/config"
	share_domain "vault-app/internal/domain/shared"
	share_infrastructure "vault-app/internal/infrastructure/share"
	"vault-app/internal/tracecore"

	"gorm.io/gorm"
)

type DBModel struct {
	DB *gorm.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// New models return a model type with database connection pool
func NewModels(db *gorm.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Folder{},
		&VaultCID{},
		&VaultContent{},
		&LoginEntry{},
		&CardEntry{},
		&IdentityEntry{},
		&NoteEntry{},
		&SSHKeyEntry{},

		// App & User Configs
		&app_config.AppConfig{},
		&app_config.CommitRule{},
		&app_config.UserConfig{},
		&app_config.SharingRule{},
		&app_config.SharingConfig{}, // if used for advanced sharing
		&UserSession{},
		&auth.TokenPairs{},

		// Sharing
		&share_infrastructure.ShareEntryModel{},
		&share_infrastructure.RecipientModel{},
		&share_domain.AuditLog{},	
		
	)
}

type VaultEntry interface {
	GetId() string
	GetTypeName() string
	GetName() string
}

type User struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	Username        string    `gorm:"column:username" json:"username"`
	Email           string    `gorm:"column:email" json:"email"`
	Password        string    `gorm:"column:password" json:"password"`
	Role            string    `gorm:"column:role" json:"role"`
	CreatedAt       time.Time `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"varchar(100)"`
	LastConnectedAt time.Time `json:"last_connected_at" gorm:"last_connected_at"`
}
type UserDTO struct {
	ID              int    `json:"id"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	LastConnectedAt string `json:"last_connected_at"`
}

func toUserDTO(u User) UserDTO {
	return UserDTO{
		ID:              u.ID,
		Email:           u.Email,
		Role:            u.Role,
		CreatedAt:       u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       u.UpdatedAt.Format(time.RFC3339),
		LastConnectedAt: u.LastConnectedAt.Format(time.RFC3339),
	}
}

type Folder struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"varchar(100)"`
	CreatedAt string `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt string `json:"updated_at" gorm:"varchar(100)"`
	IsDraft   bool   `json:"is_draft"`
	// VaultCID  string `json:"vault_cid"`
}

type VaultItemType string

const (
	Credentials VaultItemType = "login"    // Email/password
	Card        VaultItemType = "card"     // Payment card
	Identity    VaultItemType = "identity" // Name, address, etc.
	Note        VaultItemType = "note"     // Secure notes
	SSHKey      VaultItemType = "sshkey"   // SSH private key
)

// Vault of the user (could be versionned)
type VaultCID struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"column:name"`
	Type      string `json:"type" gorm:"column:type"`
	UserID    int    `json:"user_id" gorm:"column:user_id"`
	CID       string `json:"cid" gorm:"column:cid"` // ✅ Explicitly map this!
	TxHash    string `json:"tx_hash" gorm:"column:tx_hash"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
}

type Entries struct {
	Login    []LoginEntry    `json:"login"`
	Card     []CardEntry     `json:"card"`
	Identity []IdentityEntry `json:"identity"`
	Note     []NoteEntry     `json:"note"`
	SSHKey   []SSHKeyEntry   `json:"sshkey"`
}
type BaseVaultContent struct {
	Folders   []Folder `json:"folders"`
	Entries   Entries  `json:"entries"`
	CreatedAt string   `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt string   `json:"updated_at" gorm:"varchar(100)"`
}

// For vault initialization (signup)
type VaultPayload struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	BaseVaultContent
}

func (v *VaultPayload) GetFolder(folderID string) (Folder, Entries) {
	var folder Folder
	for _, f := range v.Folders {
		if fmt.Sprint(f.ID) == folderID {
			folder = f
			break
		}
	}
	return folder, v.GetEntriesByFolder(folderID)
}
func (v *VaultPayload) GetEntriesByFolder(folderID string) Entries {
	filtered := Entries{}

	for _, e := range v.Entries.Login {
		if e.FolderID == folderID {
			filtered.Login = append(filtered.Login, e)
		}
	}
	for _, e := range v.Entries.Card {
		if e.FolderID == folderID {
			filtered.Card = append(filtered.Card, e)
		}
	}
	for _, e := range v.Entries.Identity {
		if e.FolderID == folderID {
			filtered.Identity = append(filtered.Identity, e)
		}
	}
	for _, e := range v.Entries.Note {
		if e.FolderID == folderID {
			filtered.Note = append(filtered.Note, e)
		}
	}
	for _, e := range v.Entries.SSHKey {
		if e.FolderID == folderID {
			filtered.SSHKey = append(filtered.SSHKey, e)
		}
	}

	return filtered
}
func (v *VaultPayload) MoveEntriesToUnsorted(folderID string) Entries {
	moved := Entries{}

	// Login
	for i, e := range v.Entries.Login {
		if e.FolderID == folderID {
			v.Entries.Login[i].FolderID = ""
			moved.Login = append(moved.Login, v.Entries.Login[i])
		}
	}
	// Card
	for i, e := range v.Entries.Card {
		if e.FolderID == folderID {
			v.Entries.Card[i].FolderID = ""
			moved.Card = append(moved.Card, v.Entries.Card[i])
		}
	}
	// Identity
	for i, e := range v.Entries.Identity {
		if e.FolderID == folderID {
			v.Entries.Identity[i].FolderID = ""
			moved.Identity = append(moved.Identity, v.Entries.Identity[i])
		}
	}
	// Note
	for i, e := range v.Entries.Note {
		if e.FolderID == folderID {
			v.Entries.Note[i].FolderID = ""
			moved.Note = append(moved.Note, v.Entries.Note[i])
		}
	}
	// SSHKey
	for i, e := range v.Entries.SSHKey {
		if e.FolderID == folderID {
			v.Entries.SSHKey[i].FolderID = ""
			moved.SSHKey = append(moved.SSHKey, v.Entries.SSHKey[i])
		}
	}

	return moved
}

// undraft entries to render vault payload for syncing
// func (v *VaultPayload) Sync() VaultPayload {

// }

// For draft mode in temp storage
type VaultContent struct {
	ID        string `json:"id"`
	UserID    int    `json:"user_id"`
	CID       string `json:"cid"`
	IsDraft   bool   `json:"is_draft"`
	CreatedAt string `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt string `json:"updated_at" gorm:"varchar(100)"`
}

// Save VaultCID
func (m *DBModel) InsertCID(vault VaultCID) error {

	if err := m.DB.Save(&vault).Error; err != nil {
		return fmt.Errorf("failed to insert new vaultCID: %w", err)
	}
	return nil
}

type VaultSaveResult struct {
	CID    string `json:"cid"`
	TxHash string `json:"txHash"`
	// You can add time if you format it
	Timestamp string `json:"timestamp,omitempty"` // time.Now().Format(time.RFC3339),
}

// TODO: userId argument missing in DecryptVault()
func (m *DBModel) GetLastVaultCID() (*VaultCID, error) {
	var record VaultCID

	if err := m.DB.First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}
func (m *DBModel) GetLastVaultByCID(cid string) (*VaultCID, error) {
	var record VaultCID

	if err := m.DB.First(&record, "cid = ?", cid).Error; err != nil {
		return nil, err
	}
	return &record, nil
}
func (m *DBModel) GetLatestVaultCIDByUserID(id int) (*VaultCID, error) {
	var record VaultCID

	if err := m.DB.Order("created_at DESC").First(&record, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (m *DBModel) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := m.DB.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type EntryType string

const (
	EntryLogin    EntryType = "Login"
	EntryCard     EntryType = "Card"
	EntryIdentity EntryType = "Identity"
	EntryNote     EntryType = "Note"
	EntrySSHKey   EntryType = "SSHKey"
)

type BaseEntry struct {
	ID              string    `json:"id"`
	EntryName       string    `json:"entry_name"`
	FolderID        string    `json:"folder_id"`
	Type            EntryType `json:"type"`
	AdditionnalNote string    `json:"additionnal_note,omitempty"`
	CustomFields    JSONMap   `json:"custom_fields,omitempty" gorm:"type:jsonb"`
	Trashed         bool      `json:"trashed"`
	IsDraft         bool      `json:"is_draft"`
	IsFavorite      bool      `json:"is_favorite"`
	// or type:text if jsonb unsupported
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type JSONMap map[string]string

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONMap value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type LoginEntry struct {
	BaseEntry
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Website  string `json:"web_site,omitempty"`
}

type CardEntry struct {
	BaseEntry
	Owner      string `json:"owner"`
	Number     string `json:"number"`
	Expiration string `json:"expiration"`
	CVC        string `json:"cvc"`
}

type IdentityEntry struct {
	BaseEntry
	Genre                string `json:"genre,omitempty"`
	FirstName            string `json:"firstname,omitempty"`
	SecondFirstName      string `json:"second_firstname,omitempty"`
	LastName             string `json:"lastname,omitempty"`
	Username             string `json:"username,omitempty"`
	Company              string `json:"company,omitempty"`
	SocialSecurityNumber string `json:"social_security_number,omitempty"`
	IDNumber             string `json:"ID_number,omitempty"`
	DriverLicense        string `json:"driver_license,omitempty"`
	Mail                 string `json:"mail,omitempty"`
	Telephone            string `json:"telephone,omitempty"`
	AddressOne           string `json:"address_one,omitempty"`
	AddressTwo           string `json:"address_two,omitempty"`
	AddressThree         string `json:"address_three,omitempty"`
	City                 string `json:"city,omitempty"`
	State                string `json:"state,omitempty"`
	PostalCode           string `json:"postal_code,omitempty"`
	Country              string `json:"country,omitempty"`
}

type NoteEntry struct {
	BaseEntry
}

type SSHKeyEntry struct {
	BaseEntry
	PrivateKey   string `json:"private_key"`
	PublicKey    string `json:"public_key"`
	EFingerprint string `json:"e_fingerprint"`
}

func (e LoginEntry) GetId() string          { return e.ID }
func (e LoginEntry) GetTypeName() string    { return "login" }
func (e LoginEntry) GetName() string        { return e.EntryName }
func (e CardEntry) GetId() string           { return e.ID }
func (e CardEntry) GetTypeName() string     { return "card" }
func (e CardEntry) GetName() string         { return e.EntryName }
func (e IdentityEntry) GetId() string       { return e.ID }
func (e IdentityEntry) GetTypeName() string { return "identity" }
func (e IdentityEntry) GetName() string     { return e.EntryName }
func (e NoteEntry) GetId() string           { return e.ID }
func (e NoteEntry) GetTypeName() string     { return "note" }
func (e NoteEntry) GetName() string         { return e.EntryName }
func (e SSHKeyEntry) GetId() string         { return e.ID }
func (e SSHKeyEntry) GetTypeName() string   { return "sshkey" }
func (e SSHKeyEntry) GetName() string       { return e.EntryName }

func (db *DBModel) SaveVaultCID(vault VaultCID) (*VaultCID, error) {
	result := db.DB.Create(&vault)
	if result.Error != nil {
		return nil, fmt.Errorf("❌ failed to insert new VaultCID: %w", result.Error)
	}
	return &vault, nil
}

func (m *DBModel) DeleteDraftVaultByUserIDAndCID(vault VaultContent) error {
	if err := m.DB.Delete(&vault, "id = ?", vault.ID); err != nil {
		return fmt.Errorf("❌ failed to delete drafts: %w", err)
	}
	return nil
}
func (m *DBModel) CreateUser(user *User) (*User, error) {
	if err := m.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to save user: %w", err)
	}
	return user, nil
}
func (m *DBModel) TouchLastConnected(userID int) (*User, error) {
	now := time.Now().UTC()

	// IMPORTANT: explicitly set the model/table
	tx := m.DB.Model(&User{}).
		Where("id = ?", userID).
		Update("last_connected_at", now)

	if tx.Error != nil {
		return nil, fmt.Errorf("❌ failed to update last connection: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("❌ no rows updated for user_id=%d", userID)
	}

	var u User
	if err := m.DB.First(&u, userID).Error; err != nil {
		return nil, fmt.Errorf("❌ re-fetch updated user failed: %w", err)
	}
	return &u, nil
}

func (m *DBModel) GetDraftVaultContentByUserIDAndCID(userID int, cid string, isDraft bool) (*VaultContent, error) {
	var vault VaultContent
	if err := m.DB.Where("user_id = ? AND cid = ? AND is_draft = ?", userID, cid, isDraft).
		First(&vault).Error; err != nil {
		return nil, err
	}
	return &vault, nil
}

func (m *DBModel) UpsertVaultContent(vault *VaultContent) error {
	return m.DB.Save(&vault).Error
}

// VaultSession holds the decrypted vault during an active session
type VaultSession struct {
	UserID              int           `json:"user_id"`
	Vault               *VaultPayload // Decrypted vault format
	LastCID             string
	Dirty               bool
	LastSynced          string
	LastUpdated         string
	Mutex               sync.Mutex                 `json:"-"`
	VaultRuntimeContext VaultRuntimeContext        `json:"vault_runtime_context"`
	PendingCommits      []tracecore.CommitEnvelope `json:"pending_commits,omitempty"`
}
type UserSession struct {
	UserID      int            `gorm:"primaryKey;column:user_id" json:"user_id"`
	SessionData string         `gorm:"type:json" json:"session_data"` // Marshaled VaultSession
	UpdatedAt   string         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// 🧠 This lives in memory, and is injected into handlers (vault_handler.go, etc.)
type VaultRuntimeContext struct {
	CurrentUser    app_config.UserConfig
	AppSettings    app_config.AppConfig
	SessionSecrets map[string]string
	WorkingBranch  string
	LoadedEntries  []VaultEntry
}

func (ctx *VaultRuntimeContext) IsMultiActorMode() bool {
	return len(ctx.AppSettings.Actors) > 1
}

func (ctx *VaultRuntimeContext) CurrentUserID() string {
	return ctx.CurrentUser.ID
}
func (m *DBModel) FindUsers() ([]UserDTO, error) {
	var users []User
	if err := m.DB.Order("last_connected_at desc").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to find users accounts: %w", err)
	}
	// Convert to DTOs
	userDTOs := make([]UserDTO, len(users))
	for i, u := range users {
		userDTOs[i] = toUserDTO(u)
	}

	return userDTOs, nil
}
func (m *DBModel) FindUserById(id int) (*User, error) {
	var user User
	if err := m.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to find user with id: %d - %w", id, err)
	}

	return &user, nil
}

func (m *DBModel) SaveAppConfig(ac app_config.AppConfig) (*app_config.AppConfig, error) {
	err := m.DB.Create(&ac).Error
	if err != nil {
		return nil, fmt.Errorf("❌ failed to save app config: %w", err)
	}
	return &ac, nil
}
func (m *DBModel) SaveUserConfig(uc app_config.UserConfig) (*app_config.UserConfig, error) {
	err := m.DB.Create(&uc).Error
	if err != nil {
		return nil, fmt.Errorf("❌ failed to save app config: %w", err)
	}
	return &uc, nil
}
func (m *DBModel) GetAppConfigByUserID(userID int) (*app_config.AppConfig, error) {
	var appCfg app_config.AppConfig
	if err := m.DB.Find(&appCfg, "user_id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to find app config from user Id: %d", userID)
	}
	return &appCfg, nil
}
func (m *DBModel) GetUserConfigByUserID(userID int) (*app_config.UserConfig, error) {
	var userCfg app_config.UserConfig
	if err := m.DB.First(&userCfg, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to find user config from user Id: %d", userID)
	}
	return &userCfg, nil
}
func (m *DBModel) GetUserByPublicKey(pubKey string) (*User, *app_config.UserConfig, error) {
	var userCfg app_config.UserConfig
	if err := m.DB.First(&userCfg, "stellar_public_key = ?", pubKey).Error; err != nil {
		return nil, nil, fmt.Errorf("❌ failed to find user config from public key: %s", pubKey)
	}

	var user User
	if err := m.DB.First(&user, "id = ?", userCfg.ID).Error; err != nil {
		return nil, nil, fmt.Errorf("❌ failed to find user from userCfg id: %s", userCfg.ID)
	}

	// utils.LogPretty("user from db", user)
	return &user, &userCfg, nil
}

func (db *DBModel) SaveSession(userID int, session *VaultSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	userSession := UserSession{
		UserID:      userID,
		SessionData: string(data),
	}
	return db.DB.Save(&userSession).Error
}

func (db *DBModel) LoadSession(userID int) (*VaultSession, error) {
	var userSession UserSession
	if err := db.DB.First(&userSession, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	var session VaultSession
	if err := json.Unmarshal([]byte(userSession.SessionData), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &session, nil
}

func (db *DBModel) GetAllSessions() (map[int]*VaultSession, error) {
	var sessions []UserSession
	if err := db.DB.Find(&sessions).Error; err != nil {
		return nil, err
	}

	sessionMap := make(map[int]*VaultSession)
	for _, s := range sessions {
		var vaultSession VaultSession
		if err := json.Unmarshal([]byte(s.SessionData), &vaultSession); err != nil {
			return nil, fmt.Errorf("failed to decode session for user %d: %w", s.UserID, err)
		}
		sessionMap[s.UserID] = &vaultSession
	}

	return sessionMap, nil
}

// internal/models/db_jwt.go

// SaveJwtToken saves the refresh token for a user.
// Accepts the whole TokenPairs so callers that already have both tokens can pass them through.
//
//	func (m *DBModel) SaveJwtToken(userID int, tokens auth.TokenPairs) error {
//	    // Persist only the refresh token (access tokens are short-lived; no need to store them long-term)
//	    return m.DB.Model(&User{}).
//	        Where("id = ?", userID).
//	        Update("refresh_token", tokens.RefreshToken).Error
//	}
func (m *DBModel) SaveJwtToken(tokens auth.TokenPairs) (*auth.TokenPairs, error) {
	// Persist only the refresh token (access tokens are short-lived; no need to store them long-term)
	if err := m.DB.Model(&auth.TokenPairs{}).
		Create(&tokens).Error; err != nil {
		return nil, err
	}
	return &tokens, nil
}
func (m *DBModel) UpdateJwtToken(userID int, tokens auth.TokenPairs) (*auth.TokenPairs, error) {
	// Persist only the refresh token (access tokens are short-lived; no need to store them long-term)
	if err := m.DB.Model(&auth.TokenPairs{}).
		Where("user_id = ?", userID).
		Update("refresh_token", &tokens.RefreshToken).Error; err != nil {
		return nil, err
	}
	return &tokens, nil
}

// DeleteJwtToken removes the refresh token (logout)
func (m *DBModel) DeleteJwtToken(userID int) error {
	return m.DB.Delete(&auth.TokenPairs{}, "user_id = ?", userID).Error
}

func (m *DBModel) GetJwtTokenByUserId(userId int) (*auth.TokenPairs, error) {
	var token auth.TokenPairs
	err := m.DB.Find(&token, "user_id = ?", userId).Error
	if err != nil {
		return nil, fmt.Errorf("❌ failed to find token for user: %d", userId)
	}
	return &token, nil
}

func (m *DBModel) CreateFolder(f Folder) (*Folder, error) {
	if err := m.DB.Create(&f).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to create folder: %w", err)
	}
	return &f, nil
}
func (m *DBModel) SaveFolder(f Folder) (*Folder, error) {
	if err := m.DB.Save(&f).Error; err != nil {
		return nil, fmt.Errorf("❌ failed to save folder: %w", err)
	}
	return &f, nil
}
func (m *DBModel) GetFoldersByVault(vaultCID string) ([]Folder, error) {
	var folders []Folder
	if err := m.DB.Where("vault_cid = ?", vaultCID).Find(&folders).Error; err != nil {
		return nil, fmt.Errorf("❌ server internal error: %w", err)
	}
	return folders, nil
}
func (m *DBModel) GetFolderById(id int) (*Folder, error) {
	var folder Folder
	if err := m.DB.First(&folder, id).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}
func (m *DBModel) DeleteFolder(id int) error {
	if err := m.DB.Delete(&Folder{}, id).Error; err != nil {
		return err
	}
	return nil
}
