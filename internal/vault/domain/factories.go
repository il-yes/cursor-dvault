package vaults_domain

import (
	"github.com/google/uuid"
)

func NewAttachment(fileCID string, nodeCID string, hash string, name string, size int64, ext string) *Attachment {
	return &Attachment{
		ID:      uuid.New().String(),
		FileCID: fileCID,
		NodeCID: nodeCID,
		Hash:    hash,
		Name:    name,
		Size:    size,
		Ext:     ext,
		IsDirty: true,
		RecipientCIDs: map[string]string{},
	}
}

func NewVaultKeyring(uID string) *VaultKeyring {
	return &VaultKeyring{
		UserID: uID,
		Keys:   make([]EncryptedKey, 0),
	}
}
