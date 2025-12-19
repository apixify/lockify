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

func TestListEntriesUseCase_Execute_Success(t *testing.T) {
	env := "test"
	key1 := "test-key-1"
	key2 := "test-key-2"
	value := "test-value"
	salt := "test-salt"
	passphrase := "test-passphrase"

	vaultService := &test.MockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			savedVault, _ := model.NewVault(env, "test-fingerprint", salt)
			savedVault.SetPassphrase(passphrase)
			savedVault.SetEntry(key1, base64.StdEncoding.EncodeToString([]byte(value)))
			savedVault.SetEntry(key2, base64.StdEncoding.EncodeToString([]byte(value)))
			return savedVault, nil
		},
	}

	useCase := NewListEntriesUseCase(vaultService)

	allKeys, err := useCase.Execute(context.Background(), env)
	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.Count(t, 2, allKeys, fmt.Sprintf("length of keys error, want: 2, got: %v", len(allKeys)))
	assert.Contains(t, key1, allKeys, fmt.Sprintf("keys should contain %v", key1))
	assert.Contains(t, key2, allKeys, fmt.Sprintf("keys should contain %v", key2))
}
