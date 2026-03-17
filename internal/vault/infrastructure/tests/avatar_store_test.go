package vaults_storage_tests

import (
	"os"
	"path/filepath"
	"testing"

	vaults_storage "vault-app/internal/vault/infrastructure/storage"
)

func TestAvatarStore_SaveAndLoad(t *testing.T) {
	root := t.TempDir()

	store := vaults_storage.NewAvatarStore(root)

	userID := "user123"
	ext := ".png"
	data := []byte("avatar-image-bytes")

	// Save avatar
	relPath, err := store.Save(userID, data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	expectedRel := filepath.Join("avatars", userID+ext)
	if relPath != expectedRel {
		t.Fatalf("unexpected relative path: got %s want %s", relPath, expectedRel)
	}

	// Verify file exists
	fullPath := filepath.Join(root, relPath)
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("avatar file not created: %v", err)
	}

	// Load avatar
	loaded, err := store.Load(userID, ext)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if string(loaded) != string(data) {
		t.Fatalf("loaded data mismatch")
	}
}

func TestAvatarStore_Exists(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAvatarStore(root)

	userID := "user123"
	ext := ".png"
	data := []byte("avatar")

	// Should not exist initially
	if store.Exists(userID, ext) {
		t.Fatalf("avatar should not exist yet")
	}

	// Save avatar
	_, err := store.Save(userID, data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if !store.Exists(userID, ext) {
		t.Fatalf("avatar should exist after save")
	}
}

func TestAvatarStore_Delete(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAvatarStore(root)

	userID := "user123"
	ext := ".png"
	data := []byte("avatar")

	_, err := store.Save(userID, data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	err = store.Delete(userID, ext)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if store.Exists(userID, ext) {
		t.Fatalf("avatar should not exist after delete")
	}
}

func TestAvatarStore_Cleanup(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAvatarStore(root)

	// Create multiple avatars
	store.Save("user1", []byte("a"))
	store.Save("user2", []byte("b"))

	err := store.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	avatarDir := filepath.Join(root, "avatars")

	files, err := os.ReadDir(avatarDir)
	if err != nil {
		t.Fatalf("failed to read avatar directory: %v", err)
	}

	if len(files) != 0 {
		t.Fatalf("cleanup did not remove avatars")
	}
}

func TestAvatarStore_Overwrite(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAvatarStore(root)

	userID := "user123"
	ext := ".png"

	store.Save(userID, []byte("old"))
	store.Save(userID, []byte("new"))

	data, _ := store.Load(userID, ext)

	if string(data) != "new" {
		t.Fatalf("avatar was not overwritten")
	}
}