package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

func TestClearEnvCachedPassphraseUseCase_Execute_Success(t *testing.T) {
	clearCalled := false
	var clearedEnv string

	passphraseService := &test.MockPassphraseService{
		ClearFunc: func(vctx *model.VaultContext) error {
			clearCalled = true
			clearedEnv = vctx.Env
			return nil
		},
	}

	useCase := NewClearEnvCachedPassphraseUseCase(passphraseService)
	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false))

	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.True(t, clearCalled, "Execute() should call Clear(), but it didn't")
	assert.Equal(
		t,
		clearedEnv,
		envTest,
		fmt.Sprintf("Execute() called Clear() with env %q, want %q", clearedEnv, envTest),
	)
}

func TestClearEnvCachedPassphraseUseCase_Execute_Error(t *testing.T) {
	passphraseService := &test.MockPassphraseService{
		ClearFunc: func(vctx *model.VaultContext) error {
			return errors.New("clear error")
		},
	}

	useCase := NewClearEnvCachedPassphraseUseCase(passphraseService)

	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false))
	assert.NotNil(t, err, "Execute() with Clear error expected error, got nil")
	assert.Equal(
		t,
		err.Error(),
		"clear error",
		fmt.Sprintf("Execute() error = %q, want %q", err.Error(), "clear error"),
	)
}
