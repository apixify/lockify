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
		for k, v := range vault.Entries {
			decryptedVal, err := useCase.encryptionService.Decrypt(
				v.Value,
				vault.Meta.Salt,
				vault.Passphrase(),
			)
			if err != nil {
				return fmt.Errorf("failed to decrypt value: %v", err)
			}
			useCase.logger.Output("%s=%s\n", k, decryptedVal)
		}
	} else {
		mappedEntries := make(map[string]string)
		for k, v := range vault.Entries {
			decryptedVal, err := useCase.encryptionService.Decrypt(
				v.Value,
				vault.Meta.Salt,
				vault.Passphrase(),
			)
			if err != nil {
				return fmt.Errorf("failed to decrypt value: %v", err)
			}
			mappedEntries[k] = string(decryptedVal)
		}

		data, err := json.MarshalIndent(mappedEntries, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal entries: %v", err)
		}
		useCase.logger.Output(string(data))
	}

	return nil
}
