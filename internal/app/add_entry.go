package app

import (
	"context"
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type AddEntryUc interface {
	Execute(context.Context, AddEntryDTO) error
}

type AddEntryUseCase struct {
	vaultService      service.VaultServiceInterface
	encryptionService service.EncryptionService
}

type AddEntryDTO struct {
	Env   string
	Key   string
	Value string
}

func NewAddEntryUseCase(
	vaultService service.VaultServiceInterface,
	encryptionService service.EncryptionService,
) AddEntryUc {
	return &AddEntryUseCase{vaultService, encryptionService}
}

func (useCase *AddEntryUseCase) Execute(ctx context.Context, dto AddEntryDTO) error {
	vault, err := useCase.vaultService.Open(ctx, dto.Env)
	if err != nil {
		return fmt.Errorf("failed to open vault for environment %s: %w", dto.Env, err)
	}

	encryptedValue, err := useCase.encryptionService.Encrypt([]byte(dto.Value), vault.Meta.Salt, vault.Passphrase())
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	err = vault.SetEntry(dto.Key, encryptedValue)
	if err != nil {
		return fmt.Errorf("failed to set entry: %w", err)
	}

	return useCase.vaultService.Save(ctx, vault)
}
