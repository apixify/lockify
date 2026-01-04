package app

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// GetEntryUc defines the interface for retrieving entries from the vault.
type GetEntryUc interface {
	Execute(vctx *model.VaultContext, key string) (string, error)
}

// GetEntryUseCase implements the use case for retrieving entries from the vault.
type GetEntryUseCase struct {
	vaultService      service.VaultServiceInterface
	encryptionService service.EncryptionService
}

// NewGetEntryUseCase creates a new GetEntryUseCase instance.
func NewGetEntryUseCase(
	vaultService service.VaultServiceInterface,
	encryptionService service.EncryptionService,
) GetEntryUc {
	return &GetEntryUseCase{vaultService, encryptionService}
}

// Execute retrieves and decrypts an entry from the vault.
func (useCase *GetEntryUseCase) Execute(vctx *model.VaultContext, key string) (string, error) {
	vault, err := useCase.vaultService.Open(vctx)
	if err != nil {
		return "", err
	}

	entry, err := vault.GetEntry(key)
	if err != nil {
		return "", err
	}

	value, err := useCase.encryptionService.Decrypt(
		entry.Value,
		vault.Salt(),
		vault.Passphrase(),
	)
	if err != nil {
		return "", err
	}

	return string(value), nil
}
