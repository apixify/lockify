package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ahmed-abdelgawad92/lockify/internal/config"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/repository"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/storage"
)

// FileVaultRepository implements VaultRepository using the filesystem
type FileVaultRepository struct {
	fs     storage.FileSystem
	config config.VaultConfig
}

// NewFileVaultRepository creates a new file-based vault repository
func NewFileVaultRepository(fs storage.FileSystem, config config.VaultConfig) repository.VaultRepository {
	return &FileVaultRepository{fs, config}
}

// Create creates a new vault file
func (repo *FileVaultRepository) Create(ctx context.Context, vault *model.Vault) error {
	if vault == nil {
		return fmt.Errorf("vault cannot be nil")
	}

	vaultPath := repo.config.GetVaultPath(vault.Meta.Env)
	vault.SetPath(vaultPath)

	if err := repo.fs.MkdirAll(repo.config.BaseDir, repo.config.DirMode); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	exists, err := repo.Exists(ctx, vault.Meta.Env)
	if err != nil {
		return fmt.Errorf("failed to check vault existence: %w", err)
	}
	if exists {
		return fmt.Errorf("vault already exists for environment %q", vault.Meta.Env)
	}

	return repo.Save(ctx, vault)
}

// Load loads a vault from the filesystem
func (repo *FileVaultRepository) Load(ctx context.Context, env string) (*model.Vault, error) {
	if env == "" {
		return nil, fmt.Errorf("environment cannot be empty")
	}

	vaultPath := repo.config.GetVaultPath(env)

	data, err := repo.fs.ReadFile(vaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("vault not found for environment %q: %w", env, err)
		}
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	var vault model.Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vault: %w", err)
	}

	vault.SetPath(vaultPath)

	if vault.Meta.Env != env {
		return nil, fmt.Errorf("vault environment mismatch: expected %q, got %q", env, vault.Meta.Env)
	}

	return &vault, nil
}

// Save saves a vault to the filesystem
func (repo *FileVaultRepository) Save(ctx context.Context, vault *model.Vault) error {
	if vault == nil {
		return fmt.Errorf("vault cannot be nil")
	}

	vaultPath := vault.Path()
	if vaultPath == "" {
		vaultPath = repo.config.GetVaultPath(vault.Meta.Env)
		vault.SetPath(vaultPath)
	}

	// Ensure base directory exists
	dir := filepath.Dir(vaultPath)
	if dir != "." && dir != "" {
		if err := repo.fs.MkdirAll(dir, repo.config.DirMode); err != nil {
			return fmt.Errorf("failed to create vault directory: %w", err)
		}
	}

	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	if err := repo.fs.WriteFile(vaultPath, data, repo.config.FileMode); err != nil {
		return fmt.Errorf("failed to write vault file: %w", err)
	}

	return nil
}

// Exists checks if a vault exists for an environment
func (repo *FileVaultRepository) Exists(ctx context.Context, env string) (bool, error) {
	if env == "" {
		return false, fmt.Errorf("environment cannot be empty")
	}

	vaultPath := repo.config.GetVaultPath(env)
	_, err := repo.fs.Stat(vaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check vault existence: %w", err)
	}
	return true, nil
}
