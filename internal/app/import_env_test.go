package app

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

func TestImportEnvUseCase_Execute_Json(t *testing.T) {
	env := "test"
	key := "test-key"
	testValue := "test-value"

	entries := map[string]string{
		key: testValue,
	}

	var savedVault *model.Vault
	vaultService := &test.MockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			vault, _ := model.NewVault(env, "test-fingerprint", "salt")
			vault.SetPassphrase("passphrase")
			return vault, nil
		},
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			savedVault = vault
			return nil
		},
	}

	importService := &test.MockImportService{
		FromJsonFunc: func(r io.Reader) (map[string]string, error) {
			return entries, nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		EncryptFunc: func(plaintext []byte, encodedSalt, passphrase string) (string, error) {
			return "encrypted-" + string(plaintext), nil
		},
	}

	loggerService := &test.MockLogger{}

	useCase := NewImportEnvUseCase(vaultService, importService, encryptionService, loggerService)

	jsonInput := `{"test-key": "test-value"}`
	reader := strings.NewReader(jsonInput)

	imported, skipped, err := useCase.Execute(context.Background(), env, value.Json, reader, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if imported != 1 {
		t.Errorf("want imported: 1, got: %d", imported)
	}

	if skipped != 0 {
		t.Errorf("want skipped: 0, got: %d", skipped)
	}

	if savedVault == nil {
		t.Fatal("vault was not saved")
	}

	entry, err := savedVault.GetEntry(key)
	if err != nil {
		t.Fatalf("entry not found in vault: %v", err)
	}

	if entry.Value != "encrypted-"+testValue {
		t.Errorf("want encrypted value: %q, got: %q", "encrypted-"+testValue, entry.Value)
	}
}

func TestImportEnvUseCase_Execute_Dotenv(t *testing.T) {
	env := "test"
	key := "test-key"
	testValue := "test-value"
	encryptedValue := "encrypted-test-value"

	entries := map[string]string{
		key: testValue,
	}

	var savedVault *model.Vault
	vaultService := &test.MockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			vault, _ := model.NewVault(env, "test-fingerprint", "salt")
			vault.SetPassphrase("passphrase")
			return vault, nil
		},
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			savedVault = vault
			return nil
		},
	}

	importService := &test.MockImportService{
		FromDotEnvFunc: func(r io.Reader) (map[string]string, error) {
			return entries, nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		EncryptFunc: func(plaintext []byte, encodedSalt, passphrase string) (string, error) {
			return "encrypted-" + string(plaintext), nil
		},
	}

	loggerService := &test.MockLogger{}

	useCase := NewImportEnvUseCase(vaultService, importService, encryptionService, loggerService)

	dotenvInput := "test-key=test-value"
	reader := strings.NewReader(dotenvInput)

	imported, skipped, err := useCase.Execute(context.Background(), env, value.DotEnv, reader, false)

	assert.Nil(t, err, fmt.Sprintf("unexpected error: %v", err))
	assert.Equal(t, 1, imported)
	assert.Equal(t, 0, skipped)
	assert.NotNil(t, savedVault, "vault was not saved")

	entry, err := savedVault.GetEntry(key)
	assert.Nil(t, err, fmt.Sprintf("entry not found in vault: %v", err))
	assert.Equal(t, encryptedValue, entry.Value, fmt.Sprintf("want encrypted value: %q, got: %q", encryptedValue, entry.Value))
}
