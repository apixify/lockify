package app

import (
	"encoding/json"
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ExportEnvUc defines the interface for exporting vault entries.
type ExportEnvUc interface {
	Execute(vctx *model.VaultContext, exportFormat value.FileFormat) error
}

// ExportEnvUseCase implements the use case for exporting vault entries in various formats.
type ExportEnvUseCase struct {
	vaultService      service.VaultServiceInterface
	encryptionService service.EncryptionService
	logger            domain.Logger
}

// NewExportEnvUseCase creates a new ExportEnvUseCase instance.
func NewExportEnvUseCase(
	vaultService service.VaultServiceInterface,
	encryptionService service.EncryptionService,
	logger domain.Logger,
) ExportEnvUc {
	return &ExportEnvUseCase{vaultService, encryptionService, logger}
}

// Execute exports all entries from the vault in the specified format.
func (useCase *ExportEnvUseCase) Execute(
	vctx *model.VaultContext,
	exportFormat value.FileFormat,
) error {
	vault, err := useCase.vaultService.Open(vctx)
	if err != nil {
		return err
	}

	if exportFormat.IsDotEnv() {
		err := vault.ForEachEntry(func(key string, entry model.Entry) error {
			decryptedVal, err := useCase.encryptionService.Decrypt(
				entry.Value,
				vault.Salt(),
				vault.Passphrase(),
			)
			if err != nil {
				return fmt.Errorf("failed to decrypt value for key %q: %w", key, err)
			}
			useCase.logger.Output("%s=%s\n", key, decryptedVal)
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		mappedEntries := make(map[string]string)
		err := vault.ForEachEntry(func(key string, entry model.Entry) error {
			decryptedVal, err := useCase.encryptionService.Decrypt(
				entry.Value,
				vault.Salt(),
				vault.Passphrase(),
			)
			if err != nil {
				return fmt.Errorf("failed to decrypt value for key %q: %w", key, err)
			}
			mappedEntries[key] = string(decryptedVal)
			return nil
		})
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(mappedEntries, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal entries: %v", err)
		}
		useCase.logger.Output(string(data))
	}

	return nil
}
