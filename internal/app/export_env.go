package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apixify/lockify/internal/domain"
	"github.com/apixify/lockify/internal/domain/model/value"
	"github.com/apixify/lockify/internal/domain/service"
)

type ExportEnvUseCase struct {
	vaultService      service.VaultService
	encryptionService service.EncryptionService
	logger            domain.Logger
}

func NewExportEnvUseCase(
	vaultService service.VaultService,
	encryptionService service.EncryptionService,
	logger domain.Logger,
) ExportEnvUseCase {
	return ExportEnvUseCase{vaultService, encryptionService, logger}
}

func (useCase *ExportEnvUseCase) Execute(ctx context.Context, env string, exportFormat value.FileFormat) error {
	if !exportFormat.IsValid() {
		return fmt.Errorf("format must be either %s or %s. %s is given", value.Json, value.DotEnv, exportFormat)
	}

	vault, err := useCase.vaultService.Open(ctx, env)
	if err != nil {
		return err
	}

	if exportFormat.IsDotEnv() {
		for k, v := range vault.Entries {
			decryptedVal, _ := useCase.encryptionService.Decrypt(v.Value, vault.Meta.Salt, vault.Passphrase())
			useCase.logger.Output("%s=%s\n", k, decryptedVal)
		}
	} else {
		mappedEntries := make(map[string]string)
		for k, v := range vault.Entries {
			decryptedVal, _ := useCase.encryptionService.Decrypt(v.Value, vault.Meta.Salt, vault.Passphrase())
			mappedEntries[k] = string(decryptedVal)
		}

		data, _ := json.MarshalIndent(mappedEntries, "", "  ")
		useCase.logger.Output(string(data))
	}

	return nil
}
