package app

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

func TestExportEnvUseCase_Execute_Json(t *testing.T) {
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			vault, _ := model.NewVault(vctx.Env, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(vault, passphraseTest)
			vault.SetEntry(keyTest, valueTest)
			return vault, nil
		},
	}
	loggerService := &test.MockLogger{}
	encryptionService := &test.MockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return []byte(valueTest), nil
		},
	}

	useCase := NewExportEnvUseCase(vaultService, encryptionService, loggerService)

	useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), "json")

	var got map[string]string
	json.Unmarshal([]byte(loggerService.OutputLogs[0]), &got)
	assert.Equal(
		t,
		got[keyTest],
		valueTest,
		fmt.Sprintf("want: %q, got: %q", valueTest, got[keyTest]),
	)
}

func TestExportEnvUseCase_Execute_Dotenv(t *testing.T) {
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			vault, _ := model.NewVault(vctx.Env, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(vault, passphraseTest)
			vault.SetEntry(keyTest, valueTest)
			return vault, nil
		},
	}
	loggerService := &test.MockLogger{}
	encryptionService := &test.MockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return []byte(valueTest), nil
		},
	}

	useCase := NewExportEnvUseCase(vaultService, encryptionService, loggerService)

	useCase.Execute(
		model.NewVaultContext(context.Background(), envTest, false),
		value.DotEnv,
	)

	want := fmt.Sprintf("%s=%s\n", keyTest, valueTest)
	got := loggerService.OutputLogs[0]
	assert.Equal(t, want, got)
}
