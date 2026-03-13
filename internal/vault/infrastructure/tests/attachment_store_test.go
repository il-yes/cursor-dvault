package vaults_storage_tests

import (
	"os"
	"path/filepath"
	"testing"
	vaults_storage "vault-app/internal/vault/infrastructure/storage"
)

func TestAttachmentStore_SaveAndLoad(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAttachmentStore(root)

	data := []byte("hello attachment")

	hash, err := store.Save(data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	expectedHash := vaults_storage.HashFile(data)

	if hash != expectedHash {
		t.Fatalf("hash mismatch: got %s want %s", hash, expectedHash)
	}

	// Check file exists
	path := filepath.Join(root, "attachments", "sha256", hash[:2], hash+".enc")

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("attachment file not created: %v", err)
	}

	// Load file
	loaded, err := store.Load(hash)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if string(loaded) != string(data) {
		t.Fatalf("loaded data mismatch")
	}
}

func TestAttachmentStore_Deduplication(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAttachmentStore(root)

	data := []byte("duplicate file")

	hash1, err := store.Save(data)
	if err != nil {
		t.Fatalf("first save failed: %v", err)
	}

	hash2, err := store.Save(data)
	if err != nil {
		t.Fatalf("second save failed: %v", err)
	}

	if hash1 != hash2 {
		t.Fatalf("hash should be identical for same content")
	}

	path := filepath.Join(root, "attachments", "sha256", hash1[:2], hash1+".enc")

	files, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("failed reading directory: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("deduplication failed, expected 1 file got %d", len(files))
	}
}

func TestAttachmentStore_Exists(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAttachmentStore(root)

	data := []byte("exists test")

	hash := vaults_storage.HashFile(data)

	if store.Exists(hash) {
		t.Fatalf("file should not exist yet")
	}

	_, err := store.Save(data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if !store.Exists(hash) {
		t.Fatalf("file should exist after save")
	}
}

func TestAttachmentStore_Delete(t *testing.T) {
	root := t.TempDir()
	store := vaults_storage.NewAttachmentStore(root)

	data := []byte("delete test")

	hash, err := store.Save(data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	err = store.Delete(hash)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if store.Exists(hash) {
		t.Fatalf("file should not exist after delete")
	}
}

func TestAttachmentStore_HashDeterministic(t *testing.T) {
	data := []byte("same content")

	hash1 := vaults_storage.HashFile(data)
	hash2 := vaults_storage.HashFile(data)

	if hash1 != hash2 {
		t.Fatalf("hash should be deterministic")
	}
}