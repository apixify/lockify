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
	entries := map[string]string{
		keyTest: valueTest,
	}

	var savedVault *model.Vault
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			vault, _ := model.NewVault(envTest, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(vault, passphraseTest)
			return vault, nil
		},
		SaveFunc: func(vctx *model.VaultContext, vault *model.Vault) error {
			savedVault = vault
			return nil
		},
	}

	importService := &test.MockImportService{
		FromJSONFunc: func(r io.Reader) (map[string]string, error) {
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

	imported, skipped, err := useCase.Execute(
		model.NewVaultContext(context.Background(), envTest, false),
		value.JSON,
		reader,
		false,
	)
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

	entry, err := savedVault.GetEntry(keyTest)
	if err != nil {
		t.Fatalf("entry not found in vault: %v", err)
	}

	if entry.Value != "encrypted-"+valueTest {
		t.Errorf("want encrypted value: %q, got: %q", "encrypted-"+valueTest, entry.Value)
	}
}

func TestImportEnvUseCase_Execute_Dotenv(t *testing.T) {
	entries := map[string]string{
		keyTest: valueTest,
	}

	var savedVault *model.Vault
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			vault, _ := model.NewVault(envTest, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(vault, passphraseTest)
			return vault, nil
		},
		SaveFunc: func(vctx *model.VaultContext, vault *model.Vault) error {
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

	imported, skipped, err := useCase.Execute(
		model.NewVaultContext(context.Background(), envTest, false),
		value.DotEnv,
		reader,
		false,
	)

	assert.Nil(t, err, fmt.Sprintf("unexpected error: %v", err))
	assert.Equal(t, 1, imported)
	assert.Equal(t, 0, skipped)
	assert.NotNil(t, savedVault, "vault was not saved")

	entry, err := savedVault.GetEntry(keyTest)
	assert.Nil(t, err, fmt.Sprintf("entry not found in vault: %v", err))
	assert.Equal(
		t,
		encryptedValueTest,
		entry.Value,
		fmt.Sprintf("want encrypted value: %q, got: %q", encryptedValueTest, entry.Value),
	)
}
