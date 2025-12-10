package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type ClearCachedPassphraseUc interface {
	Execute(context.Context) error
}

type ClearCachedPassphraseUseCase struct {
	passphraseService service.PassphraseService
}

func NewClearCachedPassphraseUseCase(passphraseService service.PassphraseService) ClearCachedPassphraseUc {
	return &ClearCachedPassphraseUseCase{passphraseService}
}

func (useCase *ClearCachedPassphraseUseCase) Execute(ctx context.Context) error {
	return useCase.passphraseService.ClearAll(ctx)
}
