package app

import (
	"context"
	"fmt"

	"github.com/apixify/lockify/internal/domain/service"
)

type DeleteEntryUseCase struct {
	vaultService service.VaultService
}

func NewDeleteEntryUseCase(vaultService service.VaultService) DeleteEntryUseCase {
	return DeleteEntryUseCase{vaultService}
}

func (useCase *DeleteEntryUseCase) Execute(ctx context.Context, env, key string) error {
	vault, err := useCase.vaultService.Open(ctx, env)
	if err != nil {
		return err
	}

	if err = vault.DeleteEntry(key); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return useCase.vaultService.Save(ctx, vault)
}
