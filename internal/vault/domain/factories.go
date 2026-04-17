package vaults_domain

import (

	"github.com/google/uuid"
)

func NewAttachment(entryID string, hash string, name string, size int64) *Attachment {
	return &Attachment{
		ID:        uuid.New().String(),
		EntryID:   entryID,
		Hash:      hash,
		Name:      name,
		Size:      size,
	}
}

func NewVaultKeyring(uID string) *VaultKeyring {
	return &VaultKeyring{
		UserID: uID,
		Keys:    make([]EncryptedKey, 0),
	}
}



