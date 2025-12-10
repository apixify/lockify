package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type ListEntriesUc interface {
	Execute(ctx context.Context, env string) ([]string, error)
}
type ListEntriesUseCase struct {
	vaultService service.VaultServiceInterface
}

func NewListEntriesUseCase(vaultService service.VaultServiceInterface) ListEntriesUc {
	return &ListEntriesUseCase{vaultService}
}

func (useCase *ListEntriesUseCase) Execute(ctx context.Context, env string) ([]string, error) {
	vault, err := useCase.vaultService.Open(ctx, env)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(vault.Entries))
	for k := range vault.Entries {
		keys = append(keys, k)
	}

	return keys, nil
}
