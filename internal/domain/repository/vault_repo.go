package repository

import (
	"context"

	"github.com/apixify/lockify/internal/domain/model"
)

// VaultRepository provides operations for managing vault storage
type VaultRepository interface {
	// Create creates a new vault file
	Create(ctx context.Context, vault *model.Vault) error
	// Load loads a vault from the storage
	Load(ctx context.Context, env string) (*model.Vault, error)
	// Save saves a vault to the storage
	Save(ctx context.Context, vault *model.Vault) error
	// Exists checks if a vault exists for an environment
	Exists(ctx context.Context, env string) (bool, error)
}
