package app

import (
	"context"
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type DeleteEntryUc interface {
	Execute(ctx context.Context, env, key string) error
}

type DeleteEntryUseCase struct {
	vaultService service.VaultServiceInterface
}

func NewDeleteEntryUseCase(vaultService service.VaultServiceInterface) DeleteEntryUc {
	return &DeleteEntryUseCase{vaultService}
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
