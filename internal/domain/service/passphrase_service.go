package service

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

// PassphraseService manages passphrase retrieval and caching
type PassphraseService interface {
	// Get retrieves a passphrase from environment variable, cache, or user input
	Get(ctx context.Context, env string) (string, error)
	// GetWithConfirmation retrieves a passphrase with confirmation for new vaults
	GetWithConfirmation(ctx context.Context, env string, shouldCache bool) (string, error)
	// Cache caches a passphrase for an environment
	Cache(ctx context.Context, env string, passphrase string) error
	// Clear clears a cached passphrase for an environment
	Clear(ctx context.Context, env string) error
	// ClearAll clears all cached passphrases
	ClearAll(ctx context.Context) error
	// Validate validates a passphrase against a vault's fingerprint
	Validate(ctx context.Context, vault *model.Vault, passphrase string) error
}
