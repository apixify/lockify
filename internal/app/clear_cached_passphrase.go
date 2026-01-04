package app

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ClearCachedPassphraseUc defines the interface for clearing all cached passphrases.
type ClearCachedPassphraseUc interface {
	Execute(vctx *model.VaultContext) error
}

// ClearCachedPassphraseUseCase implements the use case for clearing all cached passphrases.
type ClearCachedPassphraseUseCase struct {
	passphraseService service.PassphraseService
}

// NewClearCachedPassphraseUseCase creates a new ClearCachedPassphraseUseCase instance.
func NewClearCachedPassphraseUseCase(
	passphraseService service.PassphraseService,
) ClearCachedPassphraseUc {
	return &ClearCachedPassphraseUseCase{passphraseService}
}

// Execute clears all cached passphrases from the system keyring.
func (useCase *ClearCachedPassphraseUseCase) Execute(vctx *model.VaultContext) error {
	return useCase.passphraseService.ClearAll(vctx)
}
