package fs

import (
	"os"

	"github.com/apixify/lockify/internal/domain/storage"
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
func (f *OSFileSystem) Stat(path string) (storage.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &fileInfo{info}, nil
}

// fileInfo wraps os.FileInfo
type fileInfo struct {
	info os.FileInfo
}

func (f *fileInfo) IsDir() bool {
	return f.info.IsDir()
}

func (f *fileInfo) Mode() uint32 {
	return uint32(f.info.Mode())
}
