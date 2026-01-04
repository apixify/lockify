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
	expectedVault, _ := model.NewVault(envTest, fingerprintTest, saltTest)

	vaultService := &test.MockVaultService{
		CreateFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			if vctx.Env != envTest {
				t.Errorf("Create() called with env %q, want %q", vctx.Env, envTest)
			}
			return expectedVault, nil
		},
	}

	useCase := NewInitializeVaultUseCase(vaultService)

	vault, err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false))

	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.NotNil(t, vault, "Execute() should return a vault, but got nil")
	assert.Equal(
		t,
		envTest,
		vault.Env(),
		fmt.Sprintf("Execute() returned vault with env %q, want %q", vault.Env(), envTest),
	)
	assert.Equal(
		t,
		fingerprintTest,
		vault.FingerPrint(),
		fmt.Sprintf(
			"Execute() returned vault with fingerprint %q, want %q",
			vault.FingerPrint(),
			fingerprintTest,
		),
	)
	assert.Equal(
		t,
		saltTest,
		vault.Salt(),
		fmt.Sprintf("Execute() returned vault with salt %q, want %q", vault.Salt(), saltTest),
	)
}

func TestInitializeVaultUseCase_Execute_CreateError(t *testing.T) {
	expectedError := errors.New("vault already exists")

	vaultService := &test.MockVaultService{
		CreateFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			return nil, expectedError
		},
	}

	useCase := NewInitializeVaultUseCase(vaultService)

	vault, err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false))

	assert.NotNil(t, err, "Execute() should return error, got nil")
	assert.Nil(t, vault, fmt.Sprintf("Execute() should return nil vault on error, got %v", vault))
	assert.Contains(
		t,
		expectedError.Error(),
		err.Error(),
		fmt.Sprintf("Execute() error = %q, want to contain %q", err.Error(), expectedError.Error()),
	)
}
