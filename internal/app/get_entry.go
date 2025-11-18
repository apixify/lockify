package app

import (
	"context"

	"github.com/apixify/lockify/internal/domain/service"
)

type GetEntryUseCase struct {
	vaultService      service.VaultService
	encryptionService service.EncryptionService
}

func NewGetEntryUseCase(vaultService service.VaultService, encryptionService service.EncryptionService) GetEntryUseCase {
	return GetEntryUseCase{vaultService, encryptionService}
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
