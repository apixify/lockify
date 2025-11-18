package app

import (
	"context"

	"github.com/apixify/lockify/internal/domain/service"
)

type ClearEnvCachedPassphraseUseCase struct {
	passphraseService service.PassphraseService
}

func NewClearEnvCachedPassphraseUseCase(passphraseService service.PassphraseService) ClearEnvCachedPassphraseUseCase {
	return ClearEnvCachedPassphraseUseCase{passphraseService}
}

func (useCase *ClearEnvCachedPassphraseUseCase) Execute(ctx context.Context, env string) error {
	return useCase.passphraseService.Clear(ctx, env)
}
