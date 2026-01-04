package service

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/config"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/repository"
)

// VaultServiceInterface defines the interface for vault operations.
type VaultServiceInterface interface {
	Open(vctx *model.VaultContext) (*model.Vault, error)
	Save(vctx *model.VaultContext, vault *model.Vault) error
	Create(vctx *model.VaultContext) (*model.Vault, error)
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
func (vs *VaultService) Create(vctx *model.VaultContext) (*model.Vault, error) {
	exists, err := vs.vaultRepo.Exists(vctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check vault existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("vault already exists for environment %q", vctx.Env)
	}

	passphrase, err := vs.passphraseService.GetWithConfirmation(vctx)
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

	vault, err := model.NewVault(vctx.Env, fingerprint, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault: %w", err)
	}

	if err := vs.vaultRepo.Create(vctx, vault); err != nil {
		return nil, fmt.Errorf("failed to save vault: %w", err)
	}

	return vault, nil
}

// Open opens an existing vault for the specified environment.
func (vs *VaultService) Open(vctx *model.VaultContext) (*model.Vault, error) {
	if exists, err := vs.vaultRepo.Exists(vctx); !exists || err != nil {
		return nil, fmt.Errorf("vault for env %s does not exist %w", vctx.Env, err)
	}

	passphrase, err := vs.passphraseService.Get(vctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve passphrase: %w", err)
	}

	vault, err := vs.vaultRepo.Load(vctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open vault for environment %s: %w", vctx.Env, err)
	}

	if err = vs.passphraseService.Validate(vctx, vault, passphrase); err != nil {
		if err = vs.passphraseService.Clear(vctx); err != nil {
			return nil, fmt.Errorf("failed to clear passphrase: %w", err)
		}
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	vault.SetPassphrase(passphrase)

	return vault, nil
}

// Save saves the vault to persistent storage.
func (vs *VaultService) Save(vctx *model.VaultContext, vault *model.Vault) error {
	return vs.vaultRepo.Save(vctx, vault)
}
