package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

func TestClearEnvCachedPassphraseUseCase_Execute_Success(t *testing.T) {
	env := "test"
	clearCalled := false
	var clearedEnv string

	passphraseService := &test.MockPassphraseService{
		ClearFunc: func(ctx context.Context, env string) error {
			clearCalled = true
			clearedEnv = env
			return nil
		},
	}

	useCase := NewClearEnvCachedPassphraseUseCase(passphraseService)
	err := useCase.Execute(context.Background(), env)

	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.True(t, clearCalled, "Execute() should call Clear(), but it didn't")
	assert.Equal(t, clearedEnv, env, fmt.Sprintf("Execute() called Clear() with env %q, want %q", clearedEnv, env))
}

func TestClearEnvCachedPassphraseUseCase_Execute_Error(t *testing.T) {
	passphraseService := &test.MockPassphraseService{
		ClearFunc: func(ctx context.Context, env string) error {
			return errors.New("clear error")
		},
	}

	useCase := NewClearEnvCachedPassphraseUseCase(passphraseService)

	err := useCase.Execute(context.Background(), "test")
	assert.NotNil(t, err, "Execute() with Clear error expected error, got nil")
	assert.Equal(t, err.Error(), "clear error", fmt.Sprintf("Execute() error = %q, want %q", err.Error(), "clear error"))
}
