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
	key1 := "test-key-1"
	key2 := "test-key-2"

	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			savedVault, _ := model.NewVault(envTest, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(savedVault, passphraseTest)
			savedVault.SetEntry(key1, base64.StdEncoding.EncodeToString([]byte(valueTest)))
			savedVault.SetEntry(key2, base64.StdEncoding.EncodeToString([]byte(valueTest)))
			return savedVault, nil
		},
	}

	useCase := NewListEntriesUseCase(vaultService)

	allKeys, err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false))
	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.Count(t, 2, allKeys, fmt.Sprintf("length of keys error, want: 2, got: %v", len(allKeys)))
	assert.Contains(t, key1, allKeys, fmt.Sprintf("keys should contain %v", key1))
	assert.Contains(t, key2, allKeys, fmt.Sprintf("keys should contain %v", key2))
}
