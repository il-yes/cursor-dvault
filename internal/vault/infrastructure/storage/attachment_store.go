package vaults_storage

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"vault-app/internal/logger/logger"
)

type AttachmentStore struct {
	Root string
	logger logger.Logger
}

func NewAttachmentStore(root string) *AttachmentStore {
	return &AttachmentStore{Root: root, logger: *logger.NewFromEnv()}
}

func HashFile(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *AttachmentStore) Save(data []byte) (string, error) {

	hash := HashFile(data)

	dir := filepath.Join(s.Root, "attachments", "sha256", hash[:2])
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		s.logger.Error("AttachmentStore - Save - Failed to create directory: %v", err)
		return "", err
	}

	path := s.path(hash)

	if _, err := os.Stat(path); err == nil {
		return hash, nil
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		s.logger.Error("AttachmentStore - Save - Failed to write file: %v", err)
		return "", err
	}

	return hash, nil
}

func (s *AttachmentStore) Load(hash string) ([]byte, error) {
	path := s.path(hash)

	data, err := os.ReadFile(path)
	if err != nil {
		s.logger.Error("AttachmentStore - Load - Failed to read file: %v", err)
		return nil, err
	}

	return data, nil
}

func (s *AttachmentStore) Exists(hash string) bool {
	path := s.path(hash)

	_, err := os.Stat(path)
	return err == nil
}

func (s *AttachmentStore) Delete(hash string) error {
	path := s.path(hash)

	err := os.Remove(path)
	if err != nil {
		s.logger.Error("AttachmentStore - Delete - Failed to remove file: %v", err)
		return err
	}

	return nil
}

func (s *AttachmentStore) Cleanup() error {
	// 1️⃣ collect hashes referenced in DB
	// 2️⃣ scan attachment directory
	// 3️⃣ delete unused files
	return nil
}
func (s *AttachmentStore) path(hash string) string {
	return filepath.Join(s.Root, "attachments", "sha256", hash[:2], hash+".enc")
}