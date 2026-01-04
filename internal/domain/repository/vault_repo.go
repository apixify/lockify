package repository

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

// VaultRepository provides operations for managing vault storage
type VaultRepository interface {
	// Create creates a new vault file
	Create(vctx *model.VaultContext, vault *model.Vault) error
	// Load loads a vault from the storage
	Load(vctx *model.VaultContext) (*model.Vault, error)
	// Save saves a vault to the storage
	Save(vctx *model.VaultContext, vault *model.Vault) error
	// Exists checks if a vault exists for an environment
	Exists(vctx *model.VaultContext) (bool, error)
}
