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

func TestInitializeVaultUseCase_Execute_Success(t *testing.T) {
	env := "test"
	fingerprint := "test-fingerprint"
	salt := "test-salt"

	expectedVault, _ := model.NewVault(env, fingerprint, salt)

	vaultService := &test.MockVaultService{
		CreateFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			if env != "test" {
				t.Errorf("Create() called with env %q, want %q", env, "test")
			}
			return expectedVault, nil
		},
	}

	useCase := NewInitializeVaultUseCase(vaultService)

	vault, err := useCase.Execute(context.Background(), env)

	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.NotNil(t, vault, "Execute() should return a vault, but got nil")
	assert.Equal(t, env, vault.Meta.Env, fmt.Sprintf("Execute() returned vault with env %q, want %q", vault.Meta.Env, env))
	assert.Equal(t, fingerprint, vault.Meta.FingerPrint, fmt.Sprintf("Execute() returned vault with fingerprint %q, want %q", vault.Meta.FingerPrint, fingerprint))
	assert.Equal(t, salt, vault.Meta.Salt, fmt.Sprintf("Execute() returned vault with salt %q, want %q", vault.Meta.Salt, salt))
}

func TestInitializeVaultUseCase_Execute_CreateError(t *testing.T) {
	env := "test"
	expectedError := errors.New("vault already exists")

	vaultService := &test.MockVaultService{
		CreateFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return nil, expectedError
		},
	}

	useCase := NewInitializeVaultUseCase(vaultService)

	vault, err := useCase.Execute(context.Background(), env)

	assert.NotNil(t, err, "Execute() should return error, got nil")
	assert.Nil(t, vault, fmt.Sprintf("Execute() should return nil vault on error, got %v", vault))
	assert.Contains(t, expectedError.Error(), err.Error(), fmt.Sprintf("Execute() error = %q, want to contain %q", err.Error(), expectedError.Error()))
}
