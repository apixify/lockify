package app

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ClearEnvCachedPassphraseUseCase implements the use case for clearing cached passphrase for a specific environment.
type ClearEnvCachedPassphraseUseCase struct {
	passphraseService service.PassphraseService
}

// NewClearEnvCachedPassphraseUseCase creates a new ClearEnvCachedPassphraseUseCase instance.
func NewClearEnvCachedPassphraseUseCase(
	passphraseService service.PassphraseService,
) ClearEnvCachedPassphraseUseCase {
	return ClearEnvCachedPassphraseUseCase{passphraseService}
}

// Execute clears the cached passphrase for the specified environment.
func (useCase *ClearEnvCachedPassphraseUseCase) Execute(vctx *model.VaultContext) error {
	return useCase.passphraseService.Clear(vctx)
}
