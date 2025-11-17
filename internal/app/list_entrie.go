package app

import (
	"context"

	"github.com/apixify/lockify/internal/domain/service"
)

type ListEntriesUseCase struct {
	vaultService service.VaultService
}

func NewListEntriesUseCase(vaultService service.VaultService) ListEntriesUseCase {
	return ListEntriesUseCase{vaultService}
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
