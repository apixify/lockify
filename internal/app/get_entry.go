package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type GetEntryUc interface {
	Execute(ctx context.Context, env, key string) (string, error)
}

type GetEntryUseCase struct {
	vaultService      service.VaultServiceInterface
	encryptionService service.EncryptionService
}

func NewGetEntryUseCase(vaultService service.VaultServiceInterface, encryptionService service.EncryptionService) GetEntryUc {
	return &GetEntryUseCase{vaultService, encryptionService}
}

func (useCase *GetEntryUseCase) Execute(ctx context.Context, env, key string) (string, error) {
	vault, err := useCase.vaultService.Open(ctx, env)
	if err != nil {
		return "", err
	}

	entry, err := vault.GetEntry(key)
	if err != nil {
		return "", err
	}

	value, err := useCase.encryptionService.Decrypt(entry.Value, vault.Meta.Salt, vault.Passphrase())
	if err != nil {
		return "", err
	}

	return string(value), nil
}
