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

const (
	envTest            = "test"
	keyTest            = "test-key"
	valueTest          = "test-value"
	fingerprintTest    = "test-fingerprint"
	saltTest           = "test-salt"
	passphraseTest     = "test-passphrase"
	encryptedValueTest = "encrypted-test-value"
)

func TestAddEntryUseCase_Execute_Success(t *testing.T) {
	testVault, _ := model.NewVault(envTest, fingerprintTest, saltTest)
	_ = test.SetPassphraseForTest(testVault, passphraseTest)

	var savedVault *model.Vault

	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			vault, _ := model.NewVault(vctx.Env, fingerprintTest, saltTest)
			_ = test.SetPassphraseForTest(vault, passphraseTest)
			return vault, nil
		},
		SaveFunc: func(vctx *model.VaultContext, vault *model.Vault) error {
			savedVault = vault
			return nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		EncryptFunc: func(plaintext []byte, encodedSalt, pwd string) (string, error) {
			if string(plaintext) != valueTest {
				t.Errorf(
					"Encrypt() called with plaintext %q, want %q",
					string(plaintext),
					valueTest,
				)
			}
			if encodedSalt != saltTest {
				t.Errorf("Encrypt() called with salt %q, want %q", encodedSalt, saltTest)
			}
			if pwd != passphraseTest {
				t.Errorf("Encrypt() called with passphrase %q, want %q", pwd, passphraseTest)
			}
			return encryptedValueTest, nil
		},
	}

	useCase := NewAddEntryUseCase(vaultService, encryptionService)

	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), AddEntryDTO{
		Env:   envTest,
		Key:   keyTest,
		Value: valueTest,
	})

	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))
	assert.NotNil(
		t,
		savedVault,
		"Execute() should call Save() with the vault, but Save() was not called",
	)

	entry, err := savedVault.GetEntry(keyTest)
	assert.Nil(
		t,
		err,
		fmt.Sprintf(
			"Execute() should add entry with key %q, but GetEntry() failed: %v",
			keyTest,
			err,
		),
	)
	assert.Equal(
		t,
		entry.Value,
		encryptedValueTest,
		fmt.Sprintf(
			"Execute() added entry with value %q, want %q",
			entry.Value,
			encryptedValueTest,
		),
	)
}

func TestAddEntryUseCase_Execute_VaultOpenError(t *testing.T) {
	vaultService := &test.MockVaultService{
		OpenFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			return nil, errors.New("open vault error")
		},
	}

	useCase := NewAddEntryUseCase(vaultService, &test.MockEncryptionService{})
	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), AddEntryDTO{
		Env:   envTest,
		Key:   keyTest,
		Value: valueTest,
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
	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), AddEntryDTO{
		Env:   envTest,
		Key:   keyTest,
		Value: valueTest,
	})

	assert.NotNil(t, err, "Execute() should return encryption error, got nil")
	assert.Contains(
		t,
		"failed to encrypt value",
		err.Error(),
		fmt.Sprintf("Execute() error = %q, want to contain 'failed to encrypt value'", err.Error()),
	)
}

func TestAddEntryUseCase_Execute_SaveError(t *testing.T) {
	useCase := NewAddEntryUseCase(&test.MockVaultService{
		SaveFunc: func(vctx *model.VaultContext, vault *model.Vault) error {
			return errors.New("save failed")
		},
	}, &test.MockEncryptionService{})

	err := useCase.Execute(model.NewVaultContext(context.Background(), envTest, false), AddEntryDTO{
		Env:   envTest,
		Key:   keyTest,
		Value: valueTest,
	})

	assert.NotNil(t, err, "Execute() should return save error, got nil")
	assert.Contains(
		t,
		"save failed",
		err.Error(),
		fmt.Sprintf("Execute() error = %q, want to contain 'save failed'", err.Error()),
	)
}
