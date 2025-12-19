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

func TestRotatePassphraseUseCase_Execute_Success(t *testing.T) {
	env := "test"
	currentPassphrase := "old-passphrase"
	newPassphrase := "new-passphrase"
	currentSalt := "old-salt"
	newSalt := "new-salt"
	currentFingerprint := "old-fingerprint"
	newFingerprint := "new-fingerprint"

	vault, _ := model.NewVault(env, currentFingerprint, currentSalt)
	vault.SetEntry("key1", "encrypted-value-1")
	vault.SetEntry("key2", "encrypted-value-2")

	var savedVault *model.Vault
	decryptCallCount := 0
	encryptCallCount := 0

	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			savedVault = vault
			return nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			decryptCallCount++
			if encodedSalt != currentSalt {
				t.Errorf("Decrypt() called with salt %q, want %q", encodedSalt, currentSalt)
			}
			if passphrase != currentPassphrase {
				t.Errorf("Decrypt() called with passphrase %q, want %q", passphrase, currentPassphrase)
			}
			return []byte("decrypted-value"), nil
		},
		EncryptFunc: func(plaintext []byte, encodedSalt, passphrase string) (string, error) {
			encryptCallCount++
			if encodedSalt != newSalt {
				t.Errorf("Encrypt() called with salt %q, want %q", encodedSalt, newSalt)
			}
			if passphrase != newPassphrase {
				t.Errorf("Encrypt() called with passphrase %q, want %q", passphrase, newPassphrase)
			}
			return "new-encrypted-value", nil
		},
	}

	hashService := &test.MockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			if hashedPassphrase != currentFingerprint {
				t.Errorf("Verify() called with fingerprint %q, want %q", hashedPassphrase, currentFingerprint)
			}
			if passphrase != currentPassphrase {
				t.Errorf("Verify() called with passphrase %q, want %q", passphrase, currentPassphrase)
			}
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			if size != 16 {
				t.Errorf("GenerateSalt() called with size %d, want 16", size)
			}
			return newSalt, nil
		},
		HashFunc: func(passphrase string) (string, error) {
			if passphrase != newPassphrase {
				t.Errorf("Hash() called with passphrase %q, want %q", passphrase, newPassphrase)
			}
			return newFingerprint, nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, encryptionService, hashService)

	err := useCase.Execute(context.Background(), env, currentPassphrase, newPassphrase)
	assert.Nil(t, err, fmt.Sprintf("Execute() returned unexpected error: %v", err))

	// Verify vault was saved with new salt and fingerprint
	assert.NotNil(t, savedVault, "Execute() should call Save() with the vault, but Save() was not called")
	assert.Equal(t, newSalt, savedVault.Meta.Salt, fmt.Sprintf("Execute() should update salt to %q, got %q", newSalt, savedVault.Meta.Salt))
	assert.Equal(t, newFingerprint, savedVault.Meta.FingerPrint, fmt.Sprintf("Execute() should update fingerprint to %q, got %q", newFingerprint, savedVault.Meta.FingerPrint))

	// Verify all entries were re-encrypted
	assert.Equal(t, 2, decryptCallCount, fmt.Sprintf("Execute() should decrypt 2 entries, decrypted %d", decryptCallCount))
	assert.Equal(t, 2, encryptCallCount, fmt.Sprintf("Execute() should encrypt 2 entries, encrypted %d", encryptCallCount))

	// Verify entries have new encrypted values
	entry1, _ := savedVault.GetEntry("key1")
	entry2, _ := savedVault.GetEntry("key2")

	assert.Equal(t, "new-encrypted-value", entry1.Value, fmt.Sprintf("Execute() should update entry1 value, got %q", entry1.Value))
	assert.Equal(t, "new-encrypted-value", entry2.Value, fmt.Sprintf("Execute() should update entry2 value, got %q", entry2.Value))
}

func TestRotatePassphraseUseCase_Execute_LoadError(t *testing.T) {
	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return nil, errors.New("load error")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &test.MockEncryptionService{}, &test.MockHashService{})

	err := useCase.Execute(context.Background(), "test", "old", "new")
	assert.NotNil(t, err, "Execute() with load error expected error, got nil")
	assert.Contains(t, "failed to open vault for environment", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'failed to open vault for environment'", err.Error()))
}

func TestRotatePassphraseUseCase_Execute_VerifyError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
	}

	hashService := &test.MockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return errors.New("invalid passphrase")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &test.MockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "wrong", "new")
	assert.NotNil(t, err, "Execute() with invalid passphrase expected error, got nil")
	assert.Contains(t, "invalid credentials", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'invalid credentials'", err.Error()))
}

func TestRotatePassphraseUseCase_Execute_GenerateSaltError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
	}

	hashService := &test.MockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "", errors.New("salt error")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &test.MockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	assert.NotNil(t, err, "Execute() with salt error expected error, got nil")
	assert.Contains(t, "failed to generate salt", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'failed to generate salt'", err.Error()))
}

func TestRotatePassphraseUseCase_Execute_HashError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
	}

	hashService := &test.MockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "", errors.New("hash error")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &test.MockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	assert.NotNil(t, err, "Execute() with hash error expected error, got nil")
	assert.Contains(t, "failed to hash the fingerprint", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'failed to hash the fingerprint'", err.Error()))
}

func TestRotatePassphraseUseCase_Execute_DecryptError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vault.SetEntry("key1", "encrypted-value")
	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			v, _ := model.NewVault(env, "fingerprint", "salt")
			v.SetEntry("key1", "encrypted-value")
			return v, nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return nil, errors.New("decrypt error")
		},
	}

	hashService := &test.MockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "new-fingerprint", nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, encryptionService, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	assert.NotNil(t, err, "Execute() with decrypt error expected error, got nil")
	assert.Contains(t, "failed to decrypt key", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'failed to decrypt key'", err.Error()))
}

func TestRotatePassphraseUseCase_Execute_EncryptError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vault.SetEntry("key1", "encrypted-value")
	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			v, _ := model.NewVault(env, "fingerprint", "salt")
			v.SetEntry("key1", "encrypted-value")
			return v, nil
		},
	}

	encryptionService := &test.MockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return []byte("decrypted"), nil
		},
		EncryptFunc: func(plaintext []byte, encodedSalt, passphrase string) (string, error) {
			return "", errors.New("encrypt error")
		},
	}

	hashService := &test.MockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "new-fingerprint", nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, encryptionService, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	assert.NotNil(t, err, "Execute() with encrypt error expected error, got nil")
	assert.Contains(t, "failed to encrypt key", err.Error(), fmt.Sprintf("Execute() error = %q, want to contain 'failed to encrypt key'", err.Error()))
}

func TestRotatePassphraseUseCase_Execute_SaveError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &test.MockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			return errors.New("save error")
		},
	}

	hashService := &test.MockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "new-fingerprint", nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &test.MockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	assert.NotNil(t, err, "Execute() with save error expected error, got nil")
	assert.Equal(t, "save error", err.Error(), fmt.Sprintf("Execute() error = %q, want %q", err.Error(), "save error"))
}
