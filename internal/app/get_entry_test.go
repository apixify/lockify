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

func TestGetEntryUseCase_Execute_Success(t *testing.T) {
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			savedVault, _ := model.NewVault(envTest, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(savedVault, passphraseTest)
			savedVault.SetEntry(keyTest, base64.StdEncoding.EncodeToString([]byte(valueTest)))
			return savedVault, nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			decodedValue, _ := base64.StdEncoding.DecodeString(ciphertext)
			return []byte(decodedValue), nil
		},
	}

	useCase := NewGetEntryUseCase(vaultService, encryptionService)

	valueRetrieved, err := useCase.Execute(
		model.NewVaultContext(context.Background(), envTest, false),
		keyTest,
	)
	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.Equal(
		t,
		valueTest,
		valueRetrieved,
		fmt.Sprintf("Execute() got %s, want %s", valueRetrieved, valueTest),
	)
}

func TestGetEntryUseCase_Execute_EntryNotFound(t *testing.T) {
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			savedVault, _ := model.NewVault(envTest, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(savedVault, passphraseTest)
			return savedVault, nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			decodedValue, _ := base64.StdEncoding.DecodeString(ciphertext)
			return []byte(decodedValue), nil
		},
	}

	useCase := NewGetEntryUseCase(vaultService, encryptionService)

	_, err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), keyTest)
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
