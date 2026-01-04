package app

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// AddEntryUc defines the interface for adding entries to the vault.
type AddEntryUc interface {
	Execute(vctx *model.VaultContext, dto AddEntryDTO) error
}

// AddEntryUseCase implements the use case for adding entries to the vault.
type AddEntryUseCase struct {
	vaultService      service.VaultServiceInterface
	encryptionService service.EncryptionService
}

// AddEntryDTO contains the data needed to add an entry to the vault.
type AddEntryDTO struct {
	Env   string
	Key   string
	Value string
}

// NewAddEntryUseCase creates a new AddEntryUseCase instance.
func NewAddEntryUseCase(
	vaultService service.VaultServiceInterface,
	encryptionService service.EncryptionService,
) AddEntryUc {
	return &AddEntryUseCase{vaultService, encryptionService}
}

// Execute adds or updates an entry in the vault.
func (useCase *AddEntryUseCase) Execute(vctx *model.VaultContext, dto AddEntryDTO) error {
	vault, err := useCase.vaultService.Open(vctx)
	if err != nil {
		return fmt.Errorf("failed to open vault for environment %s: %w", dto.Env, err)
	}

	encryptedValue, err := useCase.encryptionService.Encrypt(
		[]byte(dto.Value),
		vault.Meta.Salt,
		vault.Passphrase(),
	)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	err = vault.SetEntry(dto.Key, encryptedValue)
	if err != nil {
		return fmt.Errorf("failed to set entry: %w", err)
	}

	return useCase.vaultService.Save(vctx, vault)
}
