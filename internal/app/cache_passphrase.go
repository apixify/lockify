package app

import (
	"context"
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/repository"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// CachePassphraseUc defines the interface for caching a passphrase.
type CachePassphraseUc interface {
	Execute(ctx context.Context, env string, passphrase string) error
}

// CachePassphraseUseCase implements the use case for caching a passphrase.
type CachePassphraseUseCase struct {
	vaultRepo         repository.VaultRepository
	passphraseService service.PassphraseService
	hashService       service.HashService
}

// NewCachePassphraseUseCase creates a new CachePassphraseUseCase instance.
func NewCachePassphraseUseCase(
	vaultRepo repository.VaultRepository,
	passphraseService service.PassphraseService,
	hashService service.HashService,
) CachePassphraseUc {
	return &CachePassphraseUseCase{vaultRepo, passphraseService, hashService}
}

// Execute caches a passphrase after validating it against the vault.
func (useCase *CachePassphraseUseCase) Execute(
	ctx context.Context,
	env string,
	passphrase string,
) error {
	exists, err := useCase.vaultRepo.Exists(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to check vault existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("vault for environment %q does not exist", env)
	}

	vault, err := useCase.vaultRepo.Load(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to load vault: %w", err)
	}

	if err := useCase.passphraseService.Validate(ctx, vault, passphrase); err != nil {
		return fmt.Errorf("invalid passphrase: %w", err)
	}

	if err := useCase.passphraseService.Cache(ctx, env, passphrase); err != nil {
		return fmt.Errorf("failed to cache passphrase: %w", err)
	}

	return nil
}
