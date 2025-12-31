package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// InitUc defines the interface for initializing a new vault.
type InitUc interface {
	Execute(context.Context, string, bool) (*model.Vault, error)
}

// InitializeVaultUseCase implements the use case for initializing a new vault.
type InitializeVaultUseCase struct {
	vaultService service.VaultServiceInterface
}

// NewInitializeVaultUseCase creates a new InitializeVaultUseCase instance.
func NewInitializeVaultUseCase(vaultService service.VaultServiceInterface) InitUc {
	return &InitializeVaultUseCase{vaultService}
}

// Execute initializes a new vault for the specified environment with cache preference.
func (useCase *InitializeVaultUseCase) Execute(
	ctx context.Context,
	env string,
	shouldCache bool,
) (*model.Vault, error) {
	return useCase.vaultService.Create(ctx, env, shouldCache)
}
