package app

import (
	"context"
	"fmt"
	"io"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type ImportEnvUc interface {
	Execute(ctx context.Context, env string, format value.FileFormat, r io.Reader, overwrite bool) (int, int, error)
}

type ImportEnvUseCase struct {
	vaultService      service.VaultServiceInterface
	importService     service.ImportService
	encryptionService service.EncryptionService
	logger            domain.Logger
}

func NewImportEnvUseCase(
	vaultService service.VaultServiceInterface,
	importService service.ImportService,
	encryptionService service.EncryptionService,
	logger domain.Logger,
) ImportEnvUc {
	return &ImportEnvUseCase{vaultService, importService, encryptionService, logger}
}

func (useCase *ImportEnvUseCase) Execute(ctx context.Context, env string, format value.FileFormat, r io.Reader, overwrite bool) (int, int, error) {
	imported := 0
	skipped := 0
	vault, err := useCase.vaultService.Open(ctx, env)
	if err != nil {
		return imported, skipped, fmt.Errorf("couln't open vault for env %s: %w", env, err)
	}

	var entries map[string]string
	switch format {
	case value.Json:
		entries, err = useCase.importService.FromJson(r)
	case value.DotEnv:
		entries, err = useCase.importService.FromDotEnv(r)
	default:
		return imported, skipped, fmt.Errorf("unsupported format: %q", format)
	}

	if err != nil {
		return imported, skipped, fmt.Errorf("failed to parse file: %w", err)
	}

	if len(entries) == 0 {
		return imported, skipped, fmt.Errorf("no entries found in file")
	}

	for key, value := range entries {
		_, err := vault.GetEntry(key)
		if err == nil && !overwrite {
			useCase.logger.Warning("Skipping existing key %q (use --overwrite to replace)", key)
			skipped++
			continue
		}

		encryptedValue, err := useCase.encryptionService.Encrypt([]byte(value), vault.Meta.Salt, vault.Passphrase())
		if err != nil {
			return imported, skipped, fmt.Errorf("failed to encrypt value: %w", err)
		}

		if err := vault.SetEntry(key, encryptedValue); err != nil {
			return imported, skipped, fmt.Errorf("failed to import key %q: %w", key, err)
		}
		imported++
	}

	if imported > 0 {
		useCase.vaultService.Save(ctx, vault)
	}

	return imported, skipped, nil
}
