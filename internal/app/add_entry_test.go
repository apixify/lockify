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

func TestAddEntryUseCase_Execute_Success(t *testing.T) {
	env := "test"
	key := "test-key"
	value := "test-value"
	encryptedValue := "encrypted-test-value"
	salt := "test-salt"
	passphrase := "test-passphrase"

	testVault, _ := model.NewVault(env, "test-fingerprint", salt)
	testVault.SetPassphrase(passphrase)

	var savedVault *model.Vault

	vaultService := &test.MockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			vault, _ := model.NewVault(env, "test-fingerprint", salt)
			vault.SetPassphrase(passphrase)
			return vault, nil
		},
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			savedVault = vault
			return nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		EncryptFunc: func(plaintext []byte, encodedSalt, pwd string) (string, error) {
			if string(plaintext) != value {
				t.Errorf("Encrypt() called with plaintext %q, want %q", string(plaintext), value)
			}
			if encodedSalt != salt {
				t.Errorf("Encrypt() called with salt %q, want %q", encodedSalt, salt)
			}
			if pwd != passphrase {
				t.Errorf("Encrypt() called with passphrase %q, want %q", pwd, passphrase)
			}
			return encryptedValue, nil
		},
	}

	useCase := NewAddEntryUseCase(vaultService, encryptionService)

	err := useCase.Execute(context.Background(), AddEntryDTO{
		Env:   env,
		Key:   key,
		Value: value,
	})

	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.NotNil(t, savedVault, "Execute() should call Save() with the vault, but Save() was not called")

	entry, err := savedVault.GetEntry(key)
	assert.Nil(t, err, fmt.Sprintf("Execute() should add entry with key %q, but GetEntry() failed: %v", key, err))
	assert.Equal(t, entry.Value, encryptedValue, fmt.Sprintf("Execute() added entry with value %q, want %q", entry.Value, encryptedValue))
}

func TestAddEntryUseCase_Execute_VaultOpenError(t *testing.T) {
	vaultService := &test.MockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return nil, errors.New("open vault error")
		},
	}

	useCase := NewAddEntryUseCase(vaultService, &test.MockEncryptionService{})
	err := useCase.Execute(context.Background(), AddEntryDTO{
		Env:   "test",
		Key:   "test-key",
		Value: "test-value",
	})

	assert.NotNil(t, err, "Execute() should return vault open error, got nil")
	assert.Contains(t, "failed to open vault for environment", err.Error())
}

func TestAddEntryUseCase_Execute_EncryptionError(t *testing.T) {
	encryptionService := &test.MockEncryptionService{
		EncryptFunc: func(plaintext []byte, encodedSalt, passphrase string) (string, error) {
			return "", errors.New("encryption failed")
		},
	}
	useCase := NewAddEntryUseCase(&test.MockVaultService{}, encryptionService)
	err := useCase.Execute(context.Background(), AddEntryDTO{
		Env:   "test",
		Key:   "test-key",
		Value: "test-value",
	})

	assert.NotNil(t, err, "Execute() should return encryption error, got nil")
	assert.Contains(t, "failed to encrypt value", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'failed to encrypt value'", err.Error()))
}

func TestAddEntryUseCase_Execute_SaveError(t *testing.T) {
	useCase := NewAddEntryUseCase(&test.MockVaultService{
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			return errors.New("save failed")
		},
	}, &test.MockEncryptionService{})

	err := useCase.Execute(context.Background(), AddEntryDTO{
		Env:   "test",
		Key:   "test-key",
		Value: "test-value",
	})

	assert.NotNil(t, err, "Execute() should return save error, got nil")
	assert.Contains(t, "save failed", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'save failed'", err.Error()))
}
