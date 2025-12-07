package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

// ============================================================================
// Shared Mocks for Use Case Tests
// ============================================================================

// mockVaultService mocks the VaultService for testing.
type mockVaultService struct {
	OpenFunc   func(ctx context.Context, env string) (*model.Vault, error)
	SaveFunc   func(ctx context.Context, vault *model.Vault) error
	CreateFunc func(ctx context.Context, env string) (*model.Vault, error)
}

func (m *mockVaultService) Open(ctx context.Context, env string) (*model.Vault, error) {
	if m.OpenFunc != nil {
		return m.OpenFunc(ctx, env)
	}
	vault, _ := model.NewVault(env, "test-fingerprint", "test-salt")
	vault.SetPassphrase("test-passphrase")
	return vault, nil
}

func (m *mockVaultService) Save(ctx context.Context, vault *model.Vault) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, vault)
	}
	return nil
}

func (m *mockVaultService) Create(ctx context.Context, env string) (*model.Vault, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, env)
	}
	vault, _ := model.NewVault(env, "test-fingerprint", "test-salt")
	return vault, nil
}

// mockEncryptionService mocks the EncryptionService for testing.
type mockEncryptionService struct {
	EncryptFunc func(plaintext []byte, encodedSalt, passphrase string) (string, error)
	DecryptFunc func(ciphertext, encodedSalt, passphrase string) ([]byte, error)
}

func (m *mockEncryptionService) Encrypt(plaintext []byte, encodedSalt, passphrase string) (string, error) {
	if m.EncryptFunc != nil {
		return m.EncryptFunc(plaintext, encodedSalt, passphrase)
	}

	return "encrypted-value", nil
}

func (m *mockEncryptionService) Decrypt(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
	if m.DecryptFunc != nil {
		return m.DecryptFunc(ciphertext, encodedSalt, passphrase)
	}

	return []byte("decrypted-value"), nil
}
