package app

import (
	"context"
	"fmt"

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

type mockLogger struct {
	InfoLogs     []string
	ErrorLogs    []string
	SuccessLogs  []string
	WarningLogs  []string
	ProgressLogs []string
	OutputLogs   []string
	InfoFunc     func(format string, args ...interface{})
	ErrorFunc    func(format string, args ...interface{})
	WarningFunc  func(format string, args ...interface{})
	SuccessFunc  func(format string, args ...interface{})
	ProgressFunc func(format string, args ...interface{})
	OutputFunc   func(format string, args ...interface{})
}

func (l *mockLogger) Info(format string, args ...interface{}) {
	l.InfoLogs = append(l.InfoLogs, fmt.Sprintf(format, args...))
	if l.InfoFunc == nil {
		return
	}

	l.InfoFunc(format, args...)
}

func (l *mockLogger) Error(format string, args ...interface{}) {
	l.ErrorLogs = append(l.ErrorLogs, fmt.Sprintf(format, args...))
	if l.ErrorFunc == nil {
		return
	}

	l.ErrorFunc(format, args...)
}

func (l *mockLogger) Warning(format string, args ...interface{}) {
	l.WarningLogs = append(l.WarningLogs, fmt.Sprintf(format, args...))
	if l.WarningFunc == nil {
		return
	}

	l.WarningFunc(format, args...)
}

func (l *mockLogger) Success(format string, args ...interface{}) {
	l.SuccessLogs = append(l.SuccessLogs, fmt.Sprintf(format, args...))
	if l.SuccessFunc == nil {
		return
	}

	l.SuccessFunc(format, args...)
}

func (l *mockLogger) Progress(format string, args ...interface{}) {
	l.ProgressLogs = append(l.ProgressLogs, fmt.Sprintf(format, args...))
	if l.ProgressFunc == nil {
		return
	}

	l.ProgressFunc(format, args...)
}

func (l *mockLogger) Output(format string, args ...interface{}) {
	l.OutputLogs = append(l.OutputLogs, fmt.Sprintf(format, args...))
	if l.OutputFunc == nil {
		return
	}

	l.OutputFunc(format, args...)
}
