package vault_infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
)


type FileLocalStore struct {
    BasePath string // directory to store vaults
}

func NewFileLocalStore(basePath string) *FileLocalStore {
    return &FileLocalStore{BasePath: basePath}
}

// helper to build the vault path
func (s *FileLocalStore) vaultPath(userID int64) string {
    return fmt.Sprintf("%s/vault_%d.enc", s.BasePath, userID)
}

// Load vault from disk
func (s *FileLocalStore) LoadVault(userID int64) ([]byte, error) {
    path := s.vaultPath(userID)
    data, err := os.ReadFile(path)
    if os.IsNotExist(err) {
        return nil, fmt.Errorf("vault not found for user %d", userID)
    }
    return data, err
}

// Save vault to disk
func (s *FileLocalStore) SaveVault(userID int64, data []byte) error {
    path := s.vaultPath(userID)
    // ensure directory exists
    if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
        return err
    }
    // atomic write
    tmp := path + ".tmp"
    if err := os.WriteFile(tmp, data, 0600); err != nil {
        return err
    }
    return os.Rename(tmp, path)
}

// Delete vault from disk
func (s *FileLocalStore) DeleteVault(userID int64) error {
    path := s.vaultPath(userID)
    if err := os.Remove(path); os.IsNotExist(err) {
        return nil
    } else {
        return err
    }
}
