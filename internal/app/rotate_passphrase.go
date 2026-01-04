package app

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/config"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/repository"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// RotatePassphraseUc defines the interface for rotating vault passphrases.
type RotatePassphraseUc interface {
	Execute(vctx *model.VaultContext, currentPassphrase, newPassphrase string) error
}

// RotatePassphraseUseCase implements the use case for rotating vault passphrases.
type RotatePassphraseUseCase struct {
	vaultRepo         repository.VaultRepository
	encryptionService service.EncryptionService
	hashService       service.HashService
}

// NewRotatePassphraseUseCase creates a new RotatePassphraseUseCase instance.
func NewRotatePassphraseUseCase(
	vaultRepo repository.VaultRepository,
	encryptionService service.EncryptionService,
	hashService service.HashService,
) RotatePassphraseUc {
	return &RotatePassphraseUseCase{vaultRepo, encryptionService, hashService}
}

// Execute rotates the passphrase for a vault by re-encrypting all entries with the new passphrase.
func (useCase *RotatePassphraseUseCase) Execute(
	vctx *model.VaultContext,
	currentPassphrase, newPassphrase string,
) error {
	vault, err := useCase.vaultRepo.Load(vctx)
	if err != nil {
		return fmt.Errorf("failed to open vault for environment %s: %w", vctx.Env, err)
	}

	if err = useCase.hashService.Verify(vault.Meta.FingerPrint, currentPassphrase); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	currentSalt := vault.Meta.Salt
	newSalt, err := useCase.hashService.GenerateSalt(config.DefaultSaltSize)
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
		decryptedValue, err := useCase.encryptionService.Decrypt(
			entry.Value,
			currentSalt,
			currentPassphrase,
		)
		if err != nil {
			return fmt.Errorf("failed to decrypt key %s: %w", key, err)
		}

		encryptedValue, err := useCase.encryptionService.Encrypt(
			decryptedValue,
			newSalt,
			newPassphrase,
		)
		if err != nil {
			return fmt.Errorf("failed to encrypt key %s: %w", key, err)
		}

		entry.Value = encryptedValue
		vault.Entries[key] = entry
	}

	return useCase.vaultRepo.Save(vctx, vault)
}
