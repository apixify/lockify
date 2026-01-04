package app

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// DeleteEntryUc defines the interface for deleting entries from the vault.
type DeleteEntryUc interface {
	Execute(vctx *model.VaultContext, key string) error
}

// DeleteEntryUseCase implements the use case for deleting entries from the vault.
type DeleteEntryUseCase struct {
	vaultService service.VaultServiceInterface
}

// NewDeleteEntryUseCase creates a new DeleteEntryUseCase instance.
func NewDeleteEntryUseCase(vaultService service.VaultServiceInterface) DeleteEntryUc {
	return &DeleteEntryUseCase{vaultService}
}

// Execute deletes an entry from the vault for the specified environment and key.
func (useCase *DeleteEntryUseCase) Execute(vctx *model.VaultContext, key string) error {
	vault, err := useCase.vaultService.Open(vctx)
	if err != nil {
		return err
	}

	if err = vault.DeleteEntry(key); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return useCase.vaultService.Save(vctx, vault)
}
