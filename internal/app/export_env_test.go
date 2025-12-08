package app

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

func TestExportEnvUseCase_Execute_Json(t *testing.T) {
	env := "test"
	key := "test-key"
	value := "test-value"
	vaultService := &mockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			vault, _ := model.NewVault(env, "test-fingerprint", "salt")
			vault.SetPassphrase("passphrase")
			vault.SetEntry(key, value)
			return vault, nil
		},
	}
	loggerService := &mockLogger{}
	encryptionService := &mockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return []byte(value), nil
		},
	}

	useCase := NewExportEnvUseCase(vaultService, encryptionService, loggerService)

	useCase.Execute(context.Background(), env, "json")

	var got map[string]string
	json.Unmarshal([]byte(loggerService.OutputLogs[0]), &got)

	if got[key] != value {
		t.Errorf("want: %q, got: %q", value, got[key])
	}
}

func TestExportEnvUseCase_Execute_Dotenv(t *testing.T) {
	env := "test"
	key := "test-key"
	value := "test-value"
	vaultService := &mockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			vault, _ := model.NewVault(env, "test-fingerprint", "salt")
			vault.SetPassphrase("passphrase")
			vault.SetEntry(key, value)
			return vault, nil
		},
	}
	loggerService := &mockLogger{}
	encryptionService := &mockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return []byte(value), nil
		},
	}

	useCase := NewExportEnvUseCase(vaultService, encryptionService, loggerService)

	useCase.Execute(context.Background(), env, "dotenv")

	want := fmt.Sprintf("%s=%s\n", key, value)
	got := loggerService.OutputLogs[0]

	if got != want {
		t.Errorf("want: %q, got: %q", want, got)
	}
}
