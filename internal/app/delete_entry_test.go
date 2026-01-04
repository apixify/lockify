package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

func TestDeleteEntryUseCase_Execute_Success(t *testing.T) {
	var savedVault *model.Vault
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			savedVault, _ = model.NewVault(vctx.Env, fingerprintTest, saltTest)
			savedVault.SetPassphrase(passphraseTest)
			savedVault.SetEntry(keyTest, base64.StdEncoding.EncodeToString([]byte(valueTest)))
			return savedVault, nil
		},
	}

	useCase := NewDeleteEntryUseCase(vaultService)

	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), keyTest)
	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))

	_, err = savedVault.GetEntry(keyTest)
	assert.NotNil(t, err, "Execute() did not delete the key successfully")
}

func TestDeleteEntryUseCase_Execute_EntryNotFound(t *testing.T) {
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			savedVault, _ := model.NewVault(vctx.Env, fingerprintTest, saltTest)
			savedVault.SetPassphrase(passphraseTest)
			return savedVault, nil
		},
	}

	useCase := NewDeleteEntryUseCase(vaultService)

	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), keyTest)
	assert.NotNil(t, err, "Execute() should return non-existence error, got nil")
	assert.Contains(
		t,
		fmt.Sprintf("key %q not found", keyTest),
		err.Error(),
		fmt.Sprintf(
			"Execute() error = %q, want to contain '%s'",
			err.Error(),
			fmt.Sprintf("key %q not found", keyTest),
		),
	)
}
