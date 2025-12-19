package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

func TestClearCachedPassphraseUseCase_Execute_Success(t *testing.T) {
	clearAllCalled := false
	passphraseService := &test.MockPassphraseService{
		ClearAllFunc: func(ctx context.Context) error {
			clearAllCalled = true
			return nil
		},
	}

	useCase := NewClearCachedPassphraseUseCase(passphraseService)
	err := useCase.Execute(context.Background())

	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.True(t, clearAllCalled, "Execute() should call ClearAll(), but it didn't")
}

func TestClearCachedPassphraseUseCase_Execute_Error(t *testing.T) {
	passphraseService := &test.MockPassphraseService{
		ClearAllFunc: func(ctx context.Context) error {
			return errors.New("clear all error")
		},
	}

	useCase := NewClearCachedPassphraseUseCase(passphraseService)

	err := useCase.Execute(context.Background())
	assert.NotNil(t, err, "Execute() with ClearAll error expected error, got nil")
	assert.Equal(t, err.Error(), "clear all error", fmt.Sprintf("Execute() error = %q, want %q", err.Error(), "clear all error"))
}
