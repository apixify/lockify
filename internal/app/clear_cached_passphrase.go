package app

import (
	"context"

	"github.com/apixify/lockify/internal/domain/service"
)

type ClearCachedPassphraseUseCase struct {
	passphraseService service.PassphraseService
}

func NewClearCachedPassphraseUseCase(passphraseService service.PassphraseService) ClearCachedPassphraseUseCase {
	return ClearCachedPassphraseUseCase{passphraseService}
}

func (useCase *ClearCachedPassphraseUseCase) Execute(ctx context.Context) error {
	return useCase.passphraseService.ClearAll(ctx)
}
