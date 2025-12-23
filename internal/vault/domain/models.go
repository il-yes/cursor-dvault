package vaults_domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"vault-app/internal/models"

	"github.com/google/uuid"
)

// former VaultCID
type Vault struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"column:name"`
	Type      string `json:"type" gorm:"column:type"`
	UserID    string `json:"user_id" gorm:"column:user_id"`
	CID       string `json:"cid" gorm:"column:cid"` // ✅ Explicitly map this!
	TxHash    string `json:"tx_hash" gorm:"column:tx_hash,omitempty"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"` // change to time.Time later
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"` // change to time.Time later
}

func NewVault(userID string, name string) *Vault {
	if name == "" {
		name = "New Vault"
	}

	return &Vault{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      "default",
		UserID:    userID,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}
func (v *Vault) AttachCID(cid string) {
    v.CID = cid
    v.UpdatedAt = time.Now().Format(time.RFC3339)
}
func (v *Vault) AttachTxHash(txHash string) {
    v.TxHash = txHash
    v.UpdatedAt = time.Now().Format(time.RFC3339)
}	
func (v *Vault) BuildInitialPayload(version string) *VaultPayload {
	return InitEmptyVaultPayload(v.Name, version)
}


type VaultContentInterface interface {
	GetFolder(folderID string) (Folder, Entries)
	GetEntriesByFolder(folderID string) Entries
	MoveEntriesToUnsorted(folderID string) Entries
}


type Folder struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"varchar(100)"`
	CreatedAt string `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt string `json:"updated_at" gorm:"varchar(100)"`
	IsDraft   bool   `json:"is_draft"`
	// VaultCID  string `json:"vault_cid"`
}
type BaseVaultContent struct {
	Folders   []Folder `json:"folders"`
	Entries   Entries  `json:"entries"`
	CreatedAt string   `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt string   `json:"updated_at" gorm:"varchar(100)"`
}

// For vault content
type VaultPayload struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	BaseVaultContent
}
func InitEmptyVaultPayload(name string, version string) *VaultPayload {
	var vp VaultPayload
	if version != "" {
		vp.Version = version
	} else {
		vp.Version = "1.0"
	}
	if name != "" {
		vp.Name = name
	} else {
		vp.Name = "New Vault"
	}
	vp.InitBaseVaultContent()	

	return &vp
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
func (v *VaultPayload) InitFolders() {	
	v.Folders = []Folder{}
}
func (v *VaultPayload) InitEntries() {
	v.Entries = Entries{}
}
func (v *VaultPayload) InitBaseVaultContent() {
	v.InitFolders()
	v.InitEntries()
}
func (v *VaultPayload) GetContentBytes() ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
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
	CreatedAt string `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt string `json:"updated_at" gorm:"autoUpdateTime"`
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

type Entries struct {
	Login    []LoginEntry    `json:"login"`
	Card     []CardEntry     `json:"card"`
	Identity []IdentityEntry `json:"identity"`
	Note     []NoteEntry     `json:"note"`
	SSHKey   []SSHKeyEntry   `json:"sshkey"`
}

type VaultEntry struct {
	ID        string    `json:"id"`
	EntryName string    `json:"entry_name"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}
// Utilities
func ParseVaultPayload(decrypted []byte) VaultPayload {
	var vault VaultPayload

	err := json.Unmarshal(decrypted, &vault)
	if err != nil {
		fmt.Println("❌ Failed to parse vault JSON:", err)
		// Fallback: return empty vault
		vault = VaultPayload{}
	}
	return vault
}

func (f *Folder) ToFormerFolder() models.Folder {
	return models.Folder{
		ID:        f.ID,
		Name:      f.Name,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		IsDraft:   f.IsDraft,
	}
}
func (e *Entries) ToFormerEntries() models.Entries {
	login := make([]models.LoginEntry, len(e.Login))
	for i, v := range e.Login {
		login[i] = v.ToFormerLoginEntry()
	}
	card := make([]models.CardEntry, len(e.Card))
	for i, v := range e.Card {
		card[i] = v.ToFormerCardEntry()
	}
	identity := make([]models.IdentityEntry, len(e.Identity))
	for i, v := range e.Identity {
		identity[i] = v.ToFormerIdentityEntry()
	}
	note := make([]models.NoteEntry, len(e.Note))
	for i, v := range e.Note {
		note[i] = v.ToFormerNoteEntry()
	}
	sshKey := make([]models.SSHKeyEntry, len(e.SSHKey))
	for i, v := range e.SSHKey {
		sshKey[i] = v.ToFormerSSHKeyEntry()
	}

	return models.Entries{
		Login:    login,
		Card:     card,
		Identity: identity,
		Note:     note,
		SSHKey:   sshKey,
	}
}
func (e *LoginEntry) ToFormerLoginEntry() models.LoginEntry {

	return models.LoginEntry{
		BaseEntry: models.BaseEntry{
			ID:              e.ID,
			EntryName:       e.EntryName,
			FolderID:        e.FolderID,
			Type:            models.EntryType(e.Type),
			AdditionnalNote: e.AdditionnalNote,
			CustomFields:    models.JSONMap(e.CustomFields),
			Trashed:         e.Trashed,
			IsDraft:         e.IsDraft,
			IsFavorite:      e.IsFavorite,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		},
		UserName: e.UserName,
		Password: e.Password,
		Website:  e.Website,
	}
}
func (v *VaultPayload) ToFormerVaultPayload() *models.VaultPayload {
	folders := make([]models.Folder, len(v.Folders))
	for i, f := range v.Folders {
		folders[i] = f.ToFormerFolder()
	}

	return &models.VaultPayload{
		Version: v.Version,
		Name:    v.Name,
		BaseVaultContent: models.BaseVaultContent{
			Folders: folders,
			Entries: v.Entries.ToFormerEntries(),
		},
	}
}
func (e *CardEntry) ToFormerCardEntry() models.CardEntry {
	return models.CardEntry{
		BaseEntry: models.BaseEntry{
			ID:              e.ID,
			EntryName:       e.EntryName,
			FolderID:        e.FolderID,
			Type:            models.EntryType(e.Type),
			AdditionnalNote: e.AdditionnalNote,
			CustomFields:    models.JSONMap(e.CustomFields),
			Trashed:         e.Trashed,
			IsDraft:         e.IsDraft,
			IsFavorite:      e.IsFavorite,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		},
		Owner:      e.Owner,
		Number:     e.Number,
		Expiration: e.Expiration,
		CVC:        e.CVC,
	}
}
func (e *IdentityEntry) ToFormerIdentityEntry() models.IdentityEntry {
	return models.IdentityEntry{
		BaseEntry: models.BaseEntry{
			ID:              e.ID,
			EntryName:       e.EntryName,
			FolderID:        e.FolderID,
			Type:            models.EntryType(e.Type),
			AdditionnalNote: e.AdditionnalNote,
			CustomFields:    models.JSONMap(e.CustomFields),
			Trashed:         e.Trashed,
			IsDraft:         e.IsDraft,
			IsFavorite:      e.IsFavorite,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		},
		Genre:                e.Genre,
		FirstName:            e.FirstName,
		SecondFirstName:      e.SecondFirstName,
		LastName:             e.LastName,
		Username:             e.Username,
		Company:              e.Company,
		SocialSecurityNumber: e.SocialSecurityNumber,
		IDNumber:             e.IDNumber,
		DriverLicense:        e.DriverLicense,
		Mail:                 e.Mail,
		Telephone:            e.Telephone,
		AddressOne:           e.AddressOne,
		AddressTwo:           e.AddressTwo,
		AddressThree:         e.AddressThree,
		City:                 e.City,
		State:                e.State,
		PostalCode:           e.PostalCode,
		Country:              e.Country,
	}
}
func (e *NoteEntry) ToFormerNoteEntry() models.NoteEntry {
	return models.NoteEntry{
		BaseEntry: models.BaseEntry{
			ID:              e.ID,
			EntryName:       e.EntryName,
			FolderID:        e.FolderID,
			Type:            models.EntryType(e.Type),
			AdditionnalNote: e.AdditionnalNote,
			CustomFields:    models.JSONMap(e.CustomFields),
			Trashed:         e.Trashed,
			IsDraft:         e.IsDraft,
			IsFavorite:      e.IsFavorite,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		},
	}
}
func (e *SSHKeyEntry) ToFormerSSHKeyEntry() models.SSHKeyEntry {
	return models.SSHKeyEntry{
		BaseEntry: models.BaseEntry{
			ID:              e.ID,
			EntryName:       e.EntryName,
			FolderID:        e.FolderID,
			Type:            models.EntryType(e.Type),
			AdditionnalNote: e.AdditionnalNote,
			CustomFields:    models.JSONMap(e.CustomFields),
			Trashed:         e.Trashed,
			IsDraft:         e.IsDraft,
			IsFavorite:      e.IsFavorite,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		},
		PrivateKey:   e.PrivateKey,
		PublicKey:    e.PublicKey,
		EFingerprint: e.EFingerprint,
	}
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
