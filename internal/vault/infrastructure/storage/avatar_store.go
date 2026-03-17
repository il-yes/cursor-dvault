// internal/vault/infrastructure/storage/avatar_store.go
package vaults_storage

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// AvatarStore manages user avatars in the vault.
type AvatarStore struct {
	Root string // root vault path
}

// NewAvatarStore creates a new AvatarStore for the given vault root.
func NewAvatarStore(root string) *AvatarStore {
	return &AvatarStore{Root: root}
}


// Save writes an avatar for a specific user.
// userID: unique identifier for the user
// data: avatar image bytes
// ext: file extension including dot (e.g., ".png", ".jpg")
func (s *AvatarStore) Save(userID string, data []byte) (string, error) {
	fmt.Println("AvatarStore - Save - userID", userID)
	dir := filepath.Join(s.Root, "avatars")		// "vault/<user_id>/<vault_name>/avatars"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	fmt.Println("AvatarStore - Save - dir", dir)

	contentType := http.DetectContentType(data)

	ext := ".bin"

	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpg"
	case "image/webp":
		ext = ".webp"
	case "image/gif":
		ext = ".gif"
	}

	filename := userID + ext
	fmt.Println("AvatarStore - Save - filename", filename)
	path := filepath.Join(dir, filename)
	fmt.Println("AvatarStore - Save - path", path)
	if err := os.WriteFile(path, data, 0644); err != nil {
		fmt.Println("AvatarStore - Save - error", err)
		return "", err
	}
	fmt.Println("AvatarStore - Save - path", path)

	return path, nil
}

// Load reads the avatar file for a user.
func (s *AvatarStore) Load(userID string, ext string) ([]byte, error) {
	path := filepath.Join(s.Root, "avatars", userID+ext)
	return os.ReadFile(path)
}

// Delete removes the avatar for a user.
func (s *AvatarStore) Delete(userID string, ext string) error {
	path := filepath.Join(s.Root, "avatars", userID+ext)
	return os.Remove(path)
}

// Exists checks if an avatar exists for a user.
func (s *AvatarStore) Exists(userID string, ext string) bool {
	path := filepath.Join(s.Root, "avatars", userID+ext)
	_, err := os.Stat(path)
	return err == nil
}

// Cleanup removes all avatars (optional utility, e.g., during full vault reset).
func (s *AvatarStore) Cleanup() error {
	dir := filepath.Join(s.Root, "avatars")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
			return err
		}
	}

	return nil
}


func (s *AvatarStore) LoadBase64(userID, ext string) (string, error) {
    path := filepath.Join(s.Root, "avatars", userID+ext)
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    contentType := http.DetectContentType(data)
    b64 := base64.StdEncoding.EncodeToString(data)
    return fmt.Sprintf("data:%s;base64,%s", contentType, b64), nil
}
