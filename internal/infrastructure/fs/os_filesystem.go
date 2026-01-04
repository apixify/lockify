package fs

import (
	"os"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/storage"
)

// OSFileSystem implements FileSystem using the OS filesystem
type OSFileSystem struct{}

// NewOSFileSystem creates a new OS filesystem implementation
func NewOSFileSystem() storage.FileSystem {
	return &OSFileSystem{}
}

// MkdirAll creates a directory and all parent directories
func (f *OSFileSystem) MkdirAll(path string, perm uint32) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// WriteFile writes data to a file
func (f *OSFileSystem) WriteFile(path string, data []byte, perm uint32) error {
	return os.WriteFile(path, data, os.FileMode(perm))
}

// ReadFile reads a file
func (f *OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Stat returns file information
func (f *OSFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
