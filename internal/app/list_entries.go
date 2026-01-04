package app

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ListEntriesUc defines the interface for listing entries in the vault.
type ListEntriesUc interface {
	Execute(vctx *model.VaultContext) ([]string, error)
}

// ListEntriesUseCase implements the use case for listing entries in the vault.
type ListEntriesUseCase struct {
	vaultService service.VaultServiceInterface
}

// NewListEntriesUseCase creates a new ListEntriesUseCase instance.
func NewListEntriesUseCase(vaultService service.VaultServiceInterface) ListEntriesUc {
	return &ListEntriesUseCase{vaultService}
}

// Execute lists all entry keys in the vault for the specified environment.
func (useCase *ListEntriesUseCase) Execute(vctx *model.VaultContext) ([]string, error) {
	vault, err := useCase.vaultService.Open(vctx)
	if err != nil {
		return nil, err
	}

	return vault.ListKeys(), nil
}
