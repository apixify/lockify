package app

import (
	"fmt"
	"io"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ImportEnvUc defines the interface for importing entries into the vault.
type ImportEnvUc interface {
	Execute(
		vctx *model.VaultContext,
		format value.FileFormat,
		r io.Reader,
		overwrite bool,
	) (int, int, error)
}

// ImportEnvUseCase implements the use case for importing entries into the vault.
type ImportEnvUseCase struct {
	vaultService      service.VaultServiceInterface
	importService     service.ImportService
	encryptionService service.EncryptionService
	logger            domain.Logger
}

// NewImportEnvUseCase creates a new ImportEnvUseCase instance.
func NewImportEnvUseCase(
	vaultService service.VaultServiceInterface,
	importService service.ImportService,
	encryptionService service.EncryptionService,
	logger domain.Logger,
) ImportEnvUc {
	return &ImportEnvUseCase{vaultService, importService, encryptionService, logger}
}

// Execute imports entries from a reader into the vault.
func (uc *ImportEnvUseCase) Execute(
	vctx *model.VaultContext,
	format value.FileFormat,
	r io.Reader,
	overwrite bool,
) (imported, skipped int, err error) {
	vault, err := uc.vaultService.Open(vctx)
	if err != nil {
		return 0, 0, fmt.Errorf("couln't open vault for env %s: %w", vctx.Env, err)
	}

	var entries map[string]string
	switch format {
	case value.JSON:
		entries, err = uc.importService.FromJSON(r)
	case value.DotEnv:
		entries, err = uc.importService.FromDotEnv(r)
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
			uc.logger.Warning("Skipping existing key %q (use --overwrite to replace)", key)
			skipped++
			continue
		}

		encryptedValue, err := uc.encryptionService.Encrypt(
			[]byte(value),
			vault.Meta.Salt,
			vault.Passphrase(),
		)
		if err != nil {
			return imported, skipped, fmt.Errorf("failed to encrypt value: %w", err)
		}

		if err := vault.SetEntry(key, encryptedValue); err != nil {
			return imported, skipped, fmt.Errorf("failed to import key %q: %w", key, err)
		}
		imported++
	}

	if imported > 0 {
		err = uc.vaultService.Save(vctx, vault)
		if err != nil {
			return imported, skipped, fmt.Errorf("failed to save vault: %w", err)
		}
	}

	return imported, skipped, nil
}
