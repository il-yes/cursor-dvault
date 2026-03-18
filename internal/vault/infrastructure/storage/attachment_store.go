package vaults_storage

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
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

	path, err := s.path(hash)
	if err != nil {
		s.logger.Error("AttachmentStore - Save - Failed to get path: %v", err)
		return "", err
	}

	if _, err := os.Stat(path); err == nil {
		return hash, nil
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		s.logger.Error("AttachmentStore - Save - Failed to write file: %v", err)
		return "", err
	}
	s.logger.Info("AttachmentStore - Save - File saved: %s", path)

	return hash, nil
}

func (s *AttachmentStore) LoadBase64(hash string) (string, error) {
	path, err := s.path(hash)
	if err != nil {
		s.logger.Error("AttachmentStore - LoadBase64 - Failed to get path: %v", err)
		return "", err
	}

    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }

    b64 := base64.StdEncoding.EncodeToString(data)

    return fmt.Sprintf("data:application/octet-stream;base64,%s", b64), nil
}

func (s *AttachmentStore) Load(hash string) ([]byte, error) {
	path, err := s.path(hash)
	if err != nil {
		s.logger.Error("AttachmentStore - Load - Failed to get path: %v", err)
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		s.logger.Error("AttachmentStore - Load - Failed to read file: %v", err)
		return nil, err
	}

	return data, nil
}

func (s *AttachmentStore) Exists(hash string) bool {
	path, err := s.path(hash)
	if err != nil {
		s.logger.Error("AttachmentStore - Exists - Failed to get path: %v", err)
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}

func (s *AttachmentStore) Delete(hash string) error {
	path, err := s.path(hash)
	if err != nil {
		s.logger.Error("AttachmentStore - Delete - Failed to get path: %v", err)
		return err
	}

	err = os.Remove(path)
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
func (s *AttachmentStore) path(hash string) (string, error) {
	if len(hash) < 2 {
		s.logger.Error("AttachmentStore - path - Invalid hash: %s", hash)
		return "", fmt.Errorf("invalid hash")
	}
	return filepath.Join(s.Root, "attachments", "sha256", hash[:2], hash+".enc"), nil
}