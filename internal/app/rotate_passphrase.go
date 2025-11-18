package app

import (
	"context"
	"fmt"

	"github.com/apixify/lockify/internal/domain/repository"
	"github.com/apixify/lockify/internal/domain/service"
)

type RotatePassphraseUseCase struct {
	vaultRepo         repository.VaultRepository
	encryptionService service.EncryptionService
	hashService       service.HashService
}

func NewRotatePassphraseUseCase(
	vaultRepo repository.VaultRepository,
	encryptionService service.EncryptionService,
	hashService service.HashService,
) RotatePassphraseUseCase {
	return RotatePassphraseUseCase{vaultRepo, encryptionService, hashService}
}

func (useCase *RotatePassphraseUseCase) Execute(ctx context.Context, env, currentPassphrase, newPassphrase string) error {
	vault, err := useCase.vaultRepo.Load(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
	}

	if err = useCase.hashService.Verify(vault.Meta.FingerPrint, currentPassphrase); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	currentSalt := vault.Meta.Salt
	newSalt, err := useCase.hashService.GenerateSalt(16)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	vault.Meta.Salt = newSalt
	vault.Meta.FingerPrint, err = useCase.hashService.Hash(newPassphrase)
	if err != nil {
		return fmt.Errorf("failed to hash the fingerprint")
	}

	for key := range vault.Entries {
		entry := vault.Entries[key]
		decryptedValue, err := useCase.encryptionService.Decrypt(entry.Value, currentSalt, currentPassphrase)
		if err != nil {
			return fmt.Errorf("failed to decrypt key %s: %w", key, err)
		}

		encryptedValue, err := useCase.encryptionService.Encrypt(decryptedValue, newSalt, newPassphrase)
		if err != nil {
			return fmt.Errorf("failed to encrypt key %s: %w", key, err)
		}

		entry.Value = encryptedValue
		vault.Entries[key] = entry
	}

	return useCase.vaultRepo.Save(ctx, vault)
}
