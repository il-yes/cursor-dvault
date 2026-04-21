package vaults_domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	// "vault-app/internal/models"
	"vault-app/internal/utils"

	"github.com/google/uuid"
)

const DefaultVaultName = "Default Vault"

// ==============================================================================
// Vault
// ==============================================================================
type Vault struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Version string `json:"version" gorm:"column:version"`

	Name               string `json:"name" gorm:"column:name"`
	Type               string `json:"type" gorm:"column:type"`
	UserID             string `json:"user_id" gorm:"column:user_id"`
	UserSubscriptionID string `json:"user_subscription_id" gorm:"column:user_subscription_id"`
	CID                string `json:"cid" gorm:"column:cid"`               // ✅ Explicitly map this!
	CreatedAt          string `json:"created_at" gorm:"column:created_at"` // change to time.Time later
	UpdatedAt          string `json:"updated_at" gorm:"column:updated_at"` // change to time.Time later
	VaultMeta          VaultMeta

	// new
	KeyVersion int
	Keyring    []WrappedKey

	Avatar string `json:"avatar" gorm:"column:avatar,omitempty"`
	TxHash string `json:"tx_hash" gorm:"column:tx_hash,omitempty"`
}

func NewVault(userID string, name string) *Vault {
	if name == "" {
		name = DefaultVaultName
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
func (v *Vault) AttachUserSubscriptionID(userSubscriptionID string) {
	v.UserSubscriptionID = userSubscriptionID
	v.UpdatedAt = time.Now().Format(time.RFC3339)
}
func (v *Vault) BuildInitialPayload(version string) *VaultPayload {
	return InitEmptyVaultPayload(v.Name, version)
}
func (v *Vault) GetVaultPath() string {
	return filepath.Join("vault", v.UserID, v.Name)
}
func (v *Vault) GetVaultAttachmentPath() string {
	return filepath.Join("vault")
}
func (v *Vault) AttachAvatar(avatar string) {
	v.Avatar = avatar
	v.UpdatedAt = time.Now().Format(time.RFC3339)
}

type VaultContentInterface interface {
	GetFolder(folderID string) (Folder, Entries)
	GetEntriesByFolder(folderID string) Entries
	MoveEntriesToUnsorted(folderID string) Entries
}

// ==============================================================================
// Folder
// ==============================================================================
type Folder struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"varchar(100)"`
	CreatedAt string `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt string `json:"updated_at" gorm:"varchar(100)"`
	IsDraft   bool   `json:"is_draft"`
	// VaultCID  string `json:"vault_cid"`
}

// ==============================================================================
// BaseVaultContent
// ==============================================================================
type BaseVaultContent struct {
	Folders   []Folder `json:"folders"`
	Entries   Entries  `json:"entries"`
	CreatedAt string   `json:"created_at" gorm:"-"`
	UpdatedAt string   `json:"updated_at" gorm:"-"`
}

// ==============================================================================
// VaultPayload
// ==============================================================================
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
		vp.Version = "1.0.0"
	}
	if name != "" {
		vp.Name = name
	} else {
		vp.Name = "New Vault"
	}
	vp.InitBaseVaultContent()

	return &vp
}
func (v *VaultPayload) InitFolders() {
	v.Folders = []Folder{}
}
func (v *VaultPayload) InitEntries() {
	v.Entries = Entries{
		Login:    []LoginEntry{},
		Card:     []CardEntry{},
		Identity: []IdentityEntry{},
		Note:     []NoteEntry{},
		SSHKey:   []SSHKeyEntry{},
	}
}
func (v *VaultPayload) InitBaseVaultContent() {
	v.InitFolders()
	v.InitEntries()
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
func (v *VaultPayload) GetEntriesByType(entryType string) []VaultEntry {
	var results []VaultEntry
	switch entryType {
	case "login":
		for _, e := range v.Entries.Login {
			results = append(results, e)
		}
	case "identity":
		for _, e := range v.Entries.Identity {
			results = append(results, e)
		}
	case "note":
		for _, e := range v.Entries.Note {
			results = append(results, e)
		}
	case "card":
		for _, e := range v.Entries.Card {
			results = append(results, e)
		}
	case "sshkey":
		for _, e := range v.Entries.SSHKey {
			results = append(results, e)
		}
	}
	return results
}
func (v *VaultPayload) GetEntry(entryType string, entryName string) VaultEntry {
	entries := v.GetEntriesByType(entryType)

	for _, entry := range entries {
		if entry.GetName() == entryName {
			return entry
		}
	}
	return nil
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
func (v *VaultPayload) GetContentBytes() ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
func (v *VaultPayload) Normalize() {
	if v.Folders == nil {
		v.Folders = []Folder{}
	}

	// Entries struct itself always exists, but slices may not
	if v.Entries.Login == nil {
		v.Entries.Login = []LoginEntry{}
	}
	if v.Entries.Card == nil {
		v.Entries.Card = []CardEntry{}
	}
	if v.Entries.Identity == nil {
		v.Entries.Identity = []IdentityEntry{}
	}
	if v.Entries.Note == nil {
		v.Entries.Note = []NoteEntry{}
	}
	if v.Entries.SSHKey == nil {
		v.Entries.SSHKey = []SSHKeyEntry{}
	}
}

func (s *VaultPayload) ToBytes() []byte {
	raw, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	return raw
}


// AddEntryAttachment adds a new attachment to the entry with entryID.
// If the entry is not found, returns an error.
func (v *VaultPayload) AddEntryAttachment(
    entryID string,
    att Attachment,
) error {
    find := func(entries interface{}) bool {
        switch xs := entries.(type) {
        case []LoginEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    xs[i].BaseEntry.Attachments = append(xs[i].BaseEntry.Attachments, att)
                    return true
                }
            }
        case []CardEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    xs[i].BaseEntry.Attachments = append(xs[i].BaseEntry.Attachments, att)
                    return true
                }
            }
        case []IdentityEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    xs[i].BaseEntry.Attachments = append(xs[i].BaseEntry.Attachments, att)
                    return true
                }
            }
        case []NoteEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    xs[i].BaseEntry.Attachments = append(xs[i].BaseEntry.Attachments, att)
                    return true
                }
            }
        case []SSHKeyEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    xs[i].BaseEntry.Attachments = append(xs[i].BaseEntry.Attachments, att)
                    return true
                }
            }
        }
        return false
    }

    found := false
    found = find(v.Entries.Login)
    if !found {
        found = find(v.Entries.Card)
    }
    if !found {
        found = find(v.Entries.Identity)
    }
    if !found {
        found = find(v.Entries.Note)
    }
    if !found {
        found = find(v.Entries.SSHKey)
    }

    if !found {
        return errors.New("entry not found")
    }

    return nil
}
func (v *VaultPayload) GetEntryAttachments(entryID string) []Attachment {
    find := func(es interface{}) []Attachment {
        switch xs := es.(type) {
        case []LoginEntry:
            for _, e := range xs {
                if e.ID == entryID {
                    return e.Attachments
                }
            }
        case []CardEntry:
            for _, e := range xs {
                if e.ID == entryID {
                    return e.Attachments
                }
            }
        case []IdentityEntry:
            for _, e := range xs {
                if e.ID == entryID {
                    return e.Attachments
                }
            }
        case []NoteEntry:
            for _, e := range xs {
                if e.ID == entryID {
                    return e.Attachments
                }
            }
        case []SSHKeyEntry:
            for _, e := range xs {
                if e.ID == entryID {
                    return e.Attachments
                }
            }
        }
        return nil
    }

    if a := find(v.Entries.Login); a != nil {
        return a
    }
    if a := find(v.Entries.Card); a != nil {
        return a
    }
    if a := find(v.Entries.Identity); a != nil {
        return a
    }
    if a := find(v.Entries.Note); a != nil {
        return a
    }
    if a := find(v.Entries.SSHKey); a != nil {
        return a
    }

    return nil
}

// UpdateEntryAttachment applies the given update function to the attachment
// inside the entry with entryID whose Attachment.ID == attachmentID.
// If the entry or attachment is not found, it returns ErrNotFound.
func (v *VaultPayload) UpdateEntryAttachment(
    entryID string,
    attachmentID string,
    updateFn func(*Attachment) error,
) error {
    findAndUpdate := func(entries interface{}) error {
        switch xs := entries.(type) {
        case []LoginEntry:
            for i := range xs {
                if xs[i].ID == entryID {
                    for j := range xs[i].Attachments {
                        if xs[i].Attachments[j].ID == attachmentID {
                            return updateFn(&xs[i].Attachments[j])
                        }
                    }
                }
            }
        case []CardEntry:
            for i := range xs {
                if xs[i].ID == entryID {
                    for j := range xs[i].Attachments {
                        if xs[i].Attachments[j].ID == attachmentID {
                            return updateFn(&xs[i].Attachments[j])
                        }
                    }
                }
            }
        case []IdentityEntry:
            for i := range xs {
                if xs[i].ID == entryID {
                    for j := range xs[i].Attachments {
                        if xs[i].Attachments[j].ID == attachmentID {
                            return updateFn(&xs[i].Attachments[j])
                        }
                    }
                }
            }
        case []NoteEntry:
            for i := range xs {
                if xs[i].ID == entryID {
                    for j := range xs[i].Attachments {
                        if xs[i].Attachments[j].ID == attachmentID {
                            return updateFn(&xs[i].Attachments[j])
                        }
                    }
                }
            }
        case []SSHKeyEntry:
            for i := range xs {
                if xs[i].ID == entryID {
                    for j := range xs[i].Attachments {
                        if xs[i].Attachments[j].ID == attachmentID {
                            return updateFn(&xs[i].Attachments[j])
                        }
                    }
                }
            }
        }
        return errors.New("entry not found")
    }

    err := findAndUpdate(v.Entries.Login)
    if err == nil {
        return nil
    }
    err = findAndUpdate(v.Entries.Card)
    if err == nil {
        return nil
    }
    err = findAndUpdate(v.Entries.Identity)
    if err == nil {
        return nil
    }
    err = findAndUpdate(v.Entries.Note)
    if err == nil {
        return nil
    }
    err = findAndUpdate(v.Entries.SSHKey)
    if err == nil {
        return nil
    }

    // reuse a sentinel error
    return errors.New("entry or attachment not found")
}

// DeleteEntryAttachment removes the attachment with attachmentID from the entry
// with entryID.
// If entry or attachment is not found, returns nil (no error).
func (v *VaultPayload) DeleteEntryAttachment(
    entryID string,
    attachmentID string,
) error {
    findAndDelete := func(entries interface{}) bool {
        switch xs := entries.(type) {
        case []LoginEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    for j, att := range xs[i].BaseEntry.Attachments {
                        if att.ID == attachmentID {
                            xs[i].BaseEntry.Attachments = append(
                                xs[i].BaseEntry.Attachments[:j],
                                xs[i].BaseEntry.Attachments[j+1:]...,
                            )
                            return true
                        }
                    }
                }
            }
        case []CardEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    for j, att := range xs[i].BaseEntry.Attachments {
                        if att.ID == attachmentID {
                            xs[i].BaseEntry.Attachments = append(
                                xs[i].BaseEntry.Attachments[:j],
                                xs[i].BaseEntry.Attachments[j+1:]...,
                            )
                            return true
                        }
                    }
                }
            }
        case []IdentityEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    for j, att := range xs[i].BaseEntry.Attachments {
                        if att.ID == attachmentID {
                            xs[i].BaseEntry.Attachments = append(
                                xs[i].BaseEntry.Attachments[:j],
                                xs[i].BaseEntry.Attachments[j+1:]...,
                            )
                            return true
                        }
                    }
                }
            }
        case []NoteEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    for j, att := range xs[i].BaseEntry.Attachments {
                        if att.ID == attachmentID {
                            xs[i].BaseEntry.Attachments = append(
                                xs[i].BaseEntry.Attachments[:j],
                                xs[i].BaseEntry.Attachments[j+1:]...,
                            )
                            return true
                        }
                    }
                }
            }
        case []SSHKeyEntry:
            for i := range xs {
                if xs[i].BaseEntry.ID == entryID {
                    for j, att := range xs[i].BaseEntry.Attachments {
                        if att.ID == attachmentID {
                            xs[i].BaseEntry.Attachments = append(
                                xs[i].BaseEntry.Attachments[:j],
                                xs[i].BaseEntry.Attachments[j+1:]...,
                            )
                            return true
                        }
                    }
                }
            }
        }
        return false
    }

    found := false
    found = findAndDelete(v.Entries.Login)
    if !found {
        found = findAndDelete(v.Entries.Card)
    }
    if !found {
        found = findAndDelete(v.Entries.Identity)
    }
    if !found {
        found = findAndDelete(v.Entries.Note)
    }
    if !found {
        found = findAndDelete(v.Entries.SSHKey)
    }

    if !found {
        return errors.New("entry or attachment not found")
    }

    return nil
}
// ==============================================================================
// VaultEntry Interface
// ==============================================================================
type VaultEntry interface {
	GetId() string
	GetTypeName() string
	GetName() string
}

// ==============================================================================
// EntryType
// ==============================================================================

type EntryType string

const (
	EntryLogin    EntryType = "Login"
	EntryCard     EntryType = "Card"
	EntryIdentity EntryType = "Identity"
	EntryNote     EntryType = "Note"
	EntrySSHKey   EntryType = "SSHKey"
)

// ==============================================================================
// BaseEntry
// ==============================================================================
type BaseEntry struct {
	ID              string       `json:"id"`
	CID             string       `json:"cid,omitempty" gorm:"column:cid"`
	EntryName       string       `json:"entry_name"`
	FolderID        string       `json:"folder_id"`
	Type            EntryType    `json:"type"`
	AdditionnalNote string       `json:"additionnal_note,omitempty"`
	CustomFields    JSONMap      `json:"custom_fields,omitempty" gorm:"type:jsonb"`
	Trashed         bool         `json:"trashed"`
	IsDraft         bool         `json:"is_draft"`
	IsDirty         bool         `json:"is_dirty"`
	IsFavorite      bool         `json:"is_favorite"`
	CreatedAt       string       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       string       `json:"updated_at" gorm:"autoUpdateTime"`
	Attachments     []Attachment `json:"attachments,omitempty" gorm:"foreignKey:EntryID"`
	AttachmentCIDs  []string     `json:"attachmentCIDs,omitempty"`
}
func ParseAndUpdateCIDs(attachments []Attachment, cids []string, atts[]string) bool {
	modified := false
	for i := range attachments {
		for j, newAttID := range atts {
			if attachments[i].ID == newAttID {
				attachments[i].HashShare = cids[j]
				modified = true
			}
		}
	}
	return modified
}

type LoginEntry struct {
	BaseEntry
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Website  string `json:"web_site,omitempty"`
}

func (e *LoginEntry) AddAttachments(attachments []Attachment) *LoginEntry {
	e.Attachments = append(e.Attachments, attachments...)
	e.UpdatedAt = time.Now().Format(time.RFC3339)
	return e
}
func (e *LoginEntry) OnShareCreated(cids []string, atts []string) bool {
	return	ParseAndUpdateCIDs(e.BaseEntry.Attachments, cids, atts)
}

type CardEntry struct {
	BaseEntry
	Owner      string `json:"owner"`
	Number     string `json:"number"`
	Expiration string `json:"expiration"`
	CVC        string `json:"cvc"`
}

func (e *CardEntry) AddAttachments(attachments []Attachment) *CardEntry {
	e.Attachments = append(e.Attachments, attachments...)
	e.UpdatedAt = time.Now().Format(time.RFC3339)
	return e
}

func (e *CardEntry) OnShareCreated(cids []string, atts []string) bool {
	return	ParseAndUpdateCIDs(e.BaseEntry.Attachments, cids, atts)
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

func (e *IdentityEntry) AddAttachments(attachments []Attachment) *IdentityEntry {
	e.Attachments = append(e.Attachments, attachments...)
	e.UpdatedAt = time.Now().Format(time.RFC3339)
	return e
}
func (e *IdentityEntry) OnShareCreated(cids []string, atts []string) bool {
	return	ParseAndUpdateCIDs(e.BaseEntry.Attachments, cids, atts)
}

type NoteEntry struct {
	BaseEntry
}

func (e *NoteEntry) AddAttachments(attachments []Attachment) *NoteEntry {
	e.Attachments = append(e.Attachments, attachments...)
	e.UpdatedAt = time.Now().Format(time.RFC3339)
	return e
}
func (e *NoteEntry) OnShareCreated(cids []string, atts []string) bool {
	return	ParseAndUpdateCIDs(e.BaseEntry.Attachments, cids, atts)
}

type SSHKeyEntry struct {
	BaseEntry
	PrivateKey   string `json:"private_key"`
	PublicKey    string `json:"public_key"`
	EFingerprint string `json:"e_fingerprint"`
}

func (e *SSHKeyEntry) AddAttachments(attachments []Attachment) *SSHKeyEntry {
	e.Attachments = append(e.Attachments, attachments...)
	e.UpdatedAt = time.Now().Format(time.RFC3339)
	return e
}
func (e *SSHKeyEntry) OnShareCreated(cids []string, atts []string) bool {
	return	ParseAndUpdateCIDs(e.BaseEntry.Attachments, cids, atts)
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

type EntryInterface interface {
	GetBase() *BaseEntry
}

func (e *LoginEntry) GetBase() *BaseEntry    { return &e.BaseEntry }
func (e *CardEntry) GetBase() *BaseEntry     { return &e.BaseEntry }
func (e *IdentityEntry) GetBase() *BaseEntry { return &e.BaseEntry }
func (e *NoteEntry) GetBase() *BaseEntry     { return &e.BaseEntry }
func (e *SSHKeyEntry) GetBase() *BaseEntry   { return &e.BaseEntry }

// ==============================================================================
// Utilities
// ==============================================================================
func ParseVaultPayload(decrypted []byte) VaultPayload {
	var vault VaultPayload
	utils.LogPretty("vault_domain - ParseVaultPayload - decrypted", decrypted)

	err := json.Unmarshal(decrypted, &vault)
	if err != nil {
		utils.LogPretty("vault_domain - ParseVaultPayload - Failed to parse vault JSON:", err)
		// Fallback: return empty vault
		vault = VaultPayload{}
	}
	// utils.LogPretty("vault_domain - ParseVaultPayload - Vault Payload", vault)
	return vault
}

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Vault models - JSONMap scan - failed to unmarshal JSONMap value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// ==============================================================================
// Attachments
// ==============================================================================
type Attachment struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	EntryID      string    `json:"entry_id"`
	Hash         string    `json:"hash"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	CID          string    `json:"cid,omitempty" gorm:"column:cid"`
	Storage      string    `json:"storage,omitempty"` // "local" | "cloud" | "ipfs";
	Ext          string    `json:"ext,omitempty" gorm:"column:ext"`
	DownloadedAt time.Time `json:"downloaded_at,omitempty" gorm:"column:downloaded_at"`
	DownloadedTo string    `json:"downloaded_to,omitempty" gorm:"column:downloaded_to"`
	HashLocal    string    `json:"hash_local" gorm:"column:hash_local"`
	HashShare    string    `json:"hash_share" gorm:"column:hash_share"`
}

// ==============================================================================
// VaultMeta
// ==============================================================================
type VaultMeta struct {
	Name       string `json:"name" gorm:"column:name"`
	UserID     string `json:"user_id" gorm:"column:user_id"`
	CreatedAt  string `json:"created_at" gorm:"column:created_at"` // change to time.Time later
	UpdatedAt  string `json:"updated_at" gorm:"column:updated_at"` // change to time.Time later"Dirty": true,
	LastSynced string `json:"last_synced" gorm:"column:updated_at"`
}

// ==============================================================================
// VaultNode
// ==============================================================================
type VaultNode struct {
	Type        string
	Version     string
	Folders     Link
	Entries     Link
	Index       Link
	Attachments Link `json:"attachments"`
}

func (v *VaultNode) ParseVaultNode(decrypted []byte) VaultNode {
	var vault VaultNode
	utils.LogPretty("vault_domain - ParseVaultPayload - decrypted", decrypted)

	err := json.Unmarshal(decrypted, &vault)
	if err != nil {
		utils.LogPretty("vault_domain - ParseVaultPayload - Failed to parse vault JSON:", err)
		// Fallback: return empty vault
		vault = VaultNode{}
	}
	utils.LogPretty("", vault)
	return vault
}

// func (v *VaultNode) Normalize() {
// 	if v.Folders == nil {
// 		v.Folders = []Folder{}
// 	}

// 	// Entries struct itself always exists, but slices may not
// 	if v.Entries.Login == nil {
// 		v.Entries.Login = []LoginEntry{}
// 	}
// 	if v.Entries.Card == nil {
// 		v.Entries.Card = []CardEntry{}
// 	}
// 	if v.Entries.Identity == nil {
// 		v.Entries.Identity = []IdentityEntry{}
// 	}
// 	if v.Entries.Note == nil {
// 		v.Entries.Note = []NoteEntry{}
// 	}
// 	if v.Entries.SSHKey == nil {
// 		v.Entries.SSHKey = []SSHKeyEntry{}
// 	}
// }

func (s *VaultNode) ToBytes() []byte {
	raw, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	return raw
}

type Link struct {
	CID string `json:"/"`
}
type EntriesRoot struct {
	Items []Link `json:"items"`
}
type FoldersRoot struct {
	Items []Link `json:"items"`
}
type AttachementsRoot struct {
	Items []Link `json:"items"`
}
type Index struct {
	ByType   map[string][]Link `json:"byType"`
	ByFolder map[string][]Link `json:"byFolder"`
}
type WrappedKey struct {
	ID        string // uuid
	Type      string // "password", "stellar", "device"
	Data      []byte // encrypted vault key
	CreatedAt string
}

// ==============================================================================
// Keyring
// ==============================================================================
type KeyType string

const (
	KeyTypeEntry      KeyType = "entry"
	KeyTypeAttachment KeyType = "attachment"
	KeyTypeIndex      KeyType = "index"
	KeyTypeVault      KeyType = "vault"
)
const (
	KeyWrapperPassword string = "password"
)

type VaultKey struct {
	Key       []byte
	Version   int
	CreatedAt time.Time
}

type EncryptedKey struct {
	ID      string
	Type    KeyType
	Version int

	// encrypted DEK
	Ciphertext []byte
	Nonce      []byte

	CreatedAt int64
}

type VaultKeyring struct {
	UserID    string
	VaultID   string
	Keys      []EncryptedKey
	Wrappers  []WrappedKey
	UpdatedAt int64
}

func (kr *VaultKeyring) AddWrapper(w WrappedKey) {
	kr.Wrappers = append(kr.Wrappers, w)
}
func (kr *VaultKeyring) GetLatestKey(t KeyType) *EncryptedKey {
	var latest *EncryptedKey

	for i := range kr.Keys {
		k := &kr.Keys[i]

		if k.Type == t {
			if latest == nil || k.Version > latest.Version {
				latest = k
			}
		}
	}

	return latest
}

/*
Yes — you are **100% moving in the right direction**, but you’re now at the point where you must clearly separate **3 encryption layers** or your architecture will become inconsistent.

I’ll clarify everything and then answer your direct question precisely.

---

# 🧠 1. You now have 3 different cryptographic responsibilities

## 🔐 (A) Vault Keyring (master unlock layer)

> “How do I unlock the vault?”

✔ password / stellar / hybrid
✔ handled by `KeyringService`
✔ outputs: **Vault Master Key (DEK or KEK)**

---

## 🔐 (B) Vault Crypto (data encryption layer)

> “How do I encrypt vault content?”

✔ entries
✔ folders
✔ attachments
✔ index
✔ IPFS payload

👉 This uses the **Vault Key (DEK)**

---

## 🔐 (C) Key Encryption (wrapping layer)

> “How do I protect the vault key itself?”

✔ password wraps vault key
✔ stellar wraps vault key
✔ hybrid storage

👉 This is your `KeyEncryption`

---

# ⚠️ 2. Your confusion (very normal here)

You wrote:

> “should I use VaultCrypto instead of CryptoService?”

👉 YES — BUT ONLY for **vault data encryption**

---

# ✅ 3. Correct mapping (VERY IMPORTANT)

## ✔ WRONG (your current issue)

```go
h.CryptoService.Encrypt(vaultBytes, cmd.Password)
```

❌ wrong because:

* password ≠ vault encryption key
* this bypasses keyring model
* breaks rotation + hybrid unlock

---

## ✔ CORRECT (DDD design)

### Step 1 — unlock vault key

```go
vaultKeyring, _ := keyringService.LoadHybrid(password, stellarSecret)
vaultKey := vaultKeyring.GetLatestKey(KeyTypeEntry)
```

---

### Step 2 — use vault key for encryption

```go
cipher := h.vaultCrypto.AESEncrypt(vaultBytes, vaultKey.Ciphertext)
```

---

# 🚀 4. Answer to your direct question

> “should I use cipher := h.crypto.AESEncrypt(entryData, dek)?”

## ✔ YES — BUT WITH ONE CONDITION:

### 👉 ONLY if `dek` = Vault Key (from keyring)

NOT:

* password ❌
* stellar secret ❌
* derived key ❌

BUT:

* VaultKeyring key ✔

---

# 🧱 5. Correct architecture flow (BOOTSTRAP + OPEN)

## 🟢 A. Vault creation (onboarding)

```text
password / stellar
        ↓
KeyringService generates VaultKey
        ↓
Wrap VaultKey (password + stellar)
        ↓
store keyring (file / db / ipfs)
```

---

## 🟡 B. Vault initialization (IPFS upload)

```text
unlock keyring
   ↓
get VaultKey (DEK)
   ↓
encrypt vault payload (entries, folders, index)
   ↓
upload to IPFS
   ↓
store CID in Vault model
```

---

## 🔵 C. Vault open

```text
unlock keyring
   ↓
get VaultKey
   ↓
fetch IPFS CID
   ↓
decrypt vault payload using VaultKey
   ↓
hydrate entries/folders/index
```

---

## 🔴 D. Vault close

```text
wipe VaultKey from memory
clear decrypted state
```

---

# ⚠️ 6. Your real correction in `CreateIPFSPayload`

## ❌ WRONG

```go
h.CryptoService.Encrypt(vaultBytes, cmd.Password)
```

---

## ✅ CORRECT (DDD version)

You must inject:

```go
KeyringService
VaultCrypto
```

---

## ✔ FINAL VERSION

```go
keyring, err := h.KeyringService.LoadHybrid(cmd.Password, cmd.StellarSecret)
if err != nil {
	return nil, err
}

vaultKey := keyring.GetLatestKey(vaults_domain.KeyTypeEntry)

encrypted, err := h.VaultCrypto.AESEncrypt(vaultBytes, vaultKey.Ciphertext)
if err != nil {
	return nil, err
}
```

---

# 🧠 7. Key insight (VERY IMPORTANT)

## Password is NEVER used for vault encryption

Password is ONLY for:

✔ unlocking keyring
✔ deriving/unwrapping vault key

---

## Vault key is ALWAYS used for:

✔ entries encryption
✔ folders encryption
✔ attachments encryption
✔ IPFS payload encryption
✔ index encryption

---

# 🚀 8. What you just built (this is big)

You now have the exact architecture used by:

* 🔐 Bitwarden (vault DEK model)
* 🔐 1Password (key hierarchy model)
* 🔐 Proton Pass (similar hybrid model)

---

# If you want next step (high value)

I can now help you implement:

## 🔥 “Vault Session Key Cache”

This gives you:

* instant unlock (no re-derivation)
* memory-only DEK
* auto-lock timer
* secure wipe on close
* multi-device sync ready

This is the final missing piece before your system feels “enterprise-grade”.

Just say **next**.

*/
