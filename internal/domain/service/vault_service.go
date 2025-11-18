package service

import (
	"context"
	"fmt"

	"github.com/apixify/lockify/internal/domain/model"
	"github.com/apixify/lockify/internal/domain/repository"
)

type VaultService struct {
	vaultRepo         repository.VaultRepository
	passphraseService PassphraseService
	hashService       HashService
}

func NewVaultService(
	vaultRepo repository.VaultRepository,
	passphraseService PassphraseService,
	hashService HashService,
) VaultService {
	return VaultService{vaultRepo, passphraseService, hashService}
}

func (vs *VaultService) Create(ctx context.Context, env string) (*model.Vault, error) {
	exists, err := vs.vaultRepo.Exists(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("failed to check vault existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("vault already exists for environment %q", env)
	}

	passphrase, err := vs.passphraseService.Get(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("failed to get passphrase: %w", err)
	}

	fingerprint, err := vs.hashService.Hash(passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to hash passphrase: %w", err)
	}

	salt, err := vs.hashService.GenerateSalt(16)
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

func (vs *VaultService) Open(ctx context.Context, env string) (*model.Vault, error) {
	if exists, _ := vs.vaultRepo.Exists(ctx, env); !exists {
		return nil, fmt.Errorf("vault for env %s does not exist", env)
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
		vs.passphraseService.Clear(ctx, env)
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	vault.SetPassphrase(passphrase)

	return vault, nil
}

func (vs *VaultService) Save(ctx context.Context, vault *model.Vault) error {
	return vs.vaultRepo.Save(ctx, vault)
}
