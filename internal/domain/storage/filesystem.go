package storage

// FileSystem abstracts file operations for testing
type FileSystem interface {
	// MkdirAll creates a directory and all parent directories
	MkdirAll(path string, perm uint32) error
	// WriteFile writes data to a file
	WriteFile(path string, data []byte, perm uint32) error
	// ReadFile reads a file
	ReadFile(path string) ([]byte, error)
	// Stat returns file information
	Stat(path string) (FileInfo, error)
}

// FileInfo represents file metadata
type FileInfo interface {
	IsDir() bool
	Mode() uint32
}
