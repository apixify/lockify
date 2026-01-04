package service

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

// PassphraseService manages passphrase retrieval and caching
type PassphraseService interface {
	// Get retrieves a passphrase from environment variable, cache, or user input
	Get(vctx *model.VaultContext) (string, error)
	// GetWithConfirmation retrieves a passphrase with confirmation for new vaults
	GetWithConfirmation(vctx *model.VaultContext) (string, error)
	// Cache caches a passphrase for an environment
	Cache(vctx *model.VaultContext, passphrase string) error
	// Clear clears a cached passphrase for an environment
	Clear(vctx *model.VaultContext) error
	// ClearAll clears all cached passphrases
	ClearAll(vctx *model.VaultContext) error
	// Validate validates a passphrase against a vault's fingerprint
	Validate(vctx *model.VaultContext, vault *model.Vault, passphrase string) error
}
