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
	env := "test"
	key := "test-key"
	value := "test-value"
	salt := "test-salt"
	passphrase := "test-passphrase"

	var savedVault *model.Vault
	vaultService := &test.MockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			savedVault, _ = model.NewVault(env, "test-fingerprint", salt)
			savedVault.SetPassphrase(passphrase)
			savedVault.SetEntry(key, base64.StdEncoding.EncodeToString([]byte(value)))
			return savedVault, nil
		},
	}

	useCase := NewDeleteEntryUseCase(vaultService)

	err := useCase.Execute(context.Background(), env, key)
	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))

	_, err = savedVault.GetEntry(key)
	assert.NotNil(t, err, "Execute() did not delete the key successfully")
}

func TestDeleteEntryUseCase_Execute_EntryNotFound(t *testing.T) {
	env := "test"
	key := "test-key"
	salt := "test-salt"
	passphrase := "test-passphrase"

	vaultService := &test.MockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			savedVault, _ := model.NewVault(env, "test-fingerprint", salt)
			savedVault.SetPassphrase(passphrase)
			return savedVault, nil
		},
	}

	useCase := NewDeleteEntryUseCase(vaultService)

	err := useCase.Execute(context.Background(), env, key)
	assert.NotNil(t, err, "Execute() should return non-existence error, got nil")
	assert.Contains(t, fmt.Sprintf("key %q not found", key), err.Error(), fmt.Sprintf("Execute() error = %q, want to contain '%s'", err.Error(), fmt.Sprintf("key %q not found", key)))
}
