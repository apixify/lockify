package app

import (
	"fmt"

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
	saltSize          int
}

// NewRotatePassphraseUseCase creates a new RotatePassphraseUseCase instance.
func NewRotatePassphraseUseCase(
	vaultRepo repository.VaultRepository,
	encryptionService service.EncryptionService,
	hashService service.HashService,
	saltSize int,
) RotatePassphraseUc {
	return &RotatePassphraseUseCase{
		vaultRepo:         vaultRepo,
		encryptionService: encryptionService,
		hashService:       hashService,
		saltSize:          saltSize,
	}
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

	if err := vault.RotatePassphrase(
		currentPassphrase,
		newPassphrase,
		useCase.encryptionService,
		useCase.hashService,
		useCase.saltSize,
	); err != nil {
		return fmt.Errorf("failed to rotate passphrase: %w", err)
	}

	return useCase.vaultRepo.Save(vctx, vault)
}
