package service

import (
	"context"
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/config"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/repository"
)

// VaultServiceInterface defines the interface for vault operations.
type VaultServiceInterface interface {
	Open(ctx context.Context, env string) (*model.Vault, error)
	Save(ctx context.Context, vault *model.Vault) error
	Create(ctx context.Context, env string, shouldCache bool) (*model.Vault, error)
}

// VaultService implements vault operations including create, open, and save.
type VaultService struct {
	vaultRepo         repository.VaultRepository
	passphraseService PassphraseService
	hashService       HashService
}

// NewVaultService creates a new VaultService instance.
func NewVaultService(
	vaultRepo repository.VaultRepository,
	passphraseService PassphraseService,
	hashService HashService,
) *VaultService {
	return &VaultService{vaultRepo, passphraseService, hashService}
}

// Create creates a new vault for the specified environment with cache preference.
func (vs *VaultService) Create(
	ctx context.Context,
	env string,
	shouldCache bool,
) (*model.Vault, error) {
	exists, err := vs.vaultRepo.Exists(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("failed to check vault existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("vault already exists for environment %q", env)
	}

	passphrase, err := vs.passphraseService.GetWithConfirmation(ctx, env, shouldCache)
	if err != nil {
		return nil, fmt.Errorf("failed to get passphrase: %w", err)
	}

	fingerprint, err := vs.hashService.Hash(passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to hash passphrase: %w", err)
	}

	salt, err := vs.hashService.GenerateSalt(config.DefaultSaltSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	vault, err := model.NewVault(env, fingerprint, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault: %w", err)
	}

	if err := vs.vaultRepo.Create(ctx, vault); err != nil {
		return nil, fmt.Errorf("failed to save vault: %w", err)
	}

	return vault, nil
}

// Open opens an existing vault for the specified environment.
func (vs *VaultService) Open(ctx context.Context, env string) (*model.Vault, error) {
	if exists, err := vs.vaultRepo.Exists(ctx, env); !exists || err != nil {
		return nil, fmt.Errorf("vault for env %s does not exist %w", env, err)
	}

	passphrase, err := vs.passphraseService.Get(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve passphrase: %w", err)
	}

	vault, err := vs.vaultRepo.Load(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("failed to open vault for environment %s: %w", env, err)
	}

	if err = vs.passphraseService.Validate(ctx, vault, passphrase); err != nil {
		if err = vs.passphraseService.Clear(ctx, env); err != nil {
			return nil, fmt.Errorf("failed to clear passphrase: %w", err)
		}
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	vault.SetPassphrase(passphrase)

	return vault, nil
}

// Save saves the vault to persistent storage.
func (vs *VaultService) Save(ctx context.Context, vault *model.Vault) error {
	return vs.vaultRepo.Save(ctx, vault)
}
