package vault_infrastructure_security

import (
	"errors"
	"os"
	"path/filepath"
)

type FileSystem interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
	ReadFile(name string) ([]byte, error)
}

type OSFileSystem struct{}

func (OSFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func (OSFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}



type FailingFS struct{}

func (FailingFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return errors.New("disk full")
}

func (FailingFS) ReadFile(name string) ([]byte, error) {
	return nil, errors.New("fail")
}