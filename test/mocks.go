package test

import (
	"context"
	"fmt"
	"io"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

// MockPromptService mocks the PromptService for testing.
type MockPromptService struct {
	GetUserInputFunc       func(isSecret bool) (key, value string, err error)
	GetPassphraseInputFunc func(message string) (string, error)
}

// GetUserInputForKeyAndValue mocks the GetUserInputForKeyAndValue method.
func (m *MockPromptService) GetUserInputForKeyAndValue(
	isSecret bool,
) (key, value string, err error) {
	if m.GetUserInputFunc != nil {
		return m.GetUserInputFunc(isSecret)
	}

	return "test_key", "test_value", nil
}

// GetPassphraseInput mocks the GetPassphraseInput method.
func (m *MockPromptService) GetPassphraseInput(message string) (string, error) {
	if m.GetPassphraseInputFunc != nil {
		return m.GetPassphraseInputFunc(message)
	}
	return "test_passphrase", nil
}

// MockVaultService mocks the VaultService for testing.
type MockVaultService struct {
	OpenFunc   func(ctx context.Context, env string) (*model.Vault, error)
	SaveFunc   func(ctx context.Context, vault *model.Vault) error
	CreateFunc func(ctx context.Context, env string) (*model.Vault, error)
}

// Open mocks the Open method.
func (m *MockVaultService) Open(ctx context.Context, env string) (*model.Vault, error) {
	if m.OpenFunc != nil {
		return m.OpenFunc(ctx, env)
	}
	vault, _ := model.NewVault(env, "test-fingerprint", "test-salt")
	vault.SetPassphrase("test-passphrase")
	return vault, nil
}

// Save mocks the Save method.
func (m *MockVaultService) Save(ctx context.Context, vault *model.Vault) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, vault)
	}
	return nil
}

// Create mocks the Create method.
func (m *MockVaultService) Create(ctx context.Context, env string) (*model.Vault, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, env)
	}
	vault, _ := model.NewVault(env, "test-fingerprint", "test-salt")
	return vault, nil
}

// MockEncryptionService mocks the EncryptionService for testing.
type MockEncryptionService struct {
	EncryptFunc func(plaintext []byte, encodedSalt, passphrase string) (string, error)
	DecryptFunc func(ciphertext, encodedSalt, passphrase string) ([]byte, error)
}

// Encrypt mocks the Encrypt method.
func (m *MockEncryptionService) Encrypt(
	plaintext []byte,
	encodedSalt, passphrase string,
) (string, error) {
	if m.EncryptFunc != nil {
		return m.EncryptFunc(plaintext, encodedSalt, passphrase)
	}

	return "encrypted-value", nil
}

// Decrypt mocks the Decrypt method.
func (m *MockEncryptionService) Decrypt(
	ciphertext, encodedSalt, passphrase string,
) ([]byte, error) {
	if m.DecryptFunc != nil {
		return m.DecryptFunc(ciphertext, encodedSalt, passphrase)
	}

	return []byte("decrypted-value"), nil
}

// MockLogger mocks the MockLogger for testing.
type MockLogger struct {
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

// Info mocks the Info method.
func (l *MockLogger) Info(format string, args ...interface{}) {
	l.InfoLogs = append(l.InfoLogs, fmt.Sprintf(format, args...))
	if l.InfoFunc == nil {
		return
	}

	l.InfoFunc(format, args...)
}

func (l *MockLogger) Error(format string, args ...interface{}) {
	l.ErrorLogs = append(l.ErrorLogs, fmt.Sprintf(format, args...))
	if l.ErrorFunc == nil {
		return
	}

	l.ErrorFunc(format, args...)
}

// Warning mocks the Warning method.
func (l *MockLogger) Warning(format string, args ...interface{}) {
	l.WarningLogs = append(l.WarningLogs, fmt.Sprintf(format, args...))
	if l.WarningFunc == nil {
		return
	}

	l.WarningFunc(format, args...)
}

// Success mocks the Success method.
func (l *MockLogger) Success(format string, args ...interface{}) {
	l.SuccessLogs = append(l.SuccessLogs, fmt.Sprintf(format, args...))
	if l.SuccessFunc == nil {
		return
	}

	l.SuccessFunc(format, args...)
}

// Progress mocks the Progress method.
func (l *MockLogger) Progress(format string, args ...interface{}) {
	l.ProgressLogs = append(l.ProgressLogs, fmt.Sprintf(format, args...))
	if l.ProgressFunc == nil {
		return
	}

	l.ProgressFunc(format, args...)
}

// Output mocks the Output method.
func (l *MockLogger) Output(format string, args ...interface{}) {
	l.OutputLogs = append(l.OutputLogs, fmt.Sprintf(format, args...))
	if l.OutputFunc == nil {
		return
	}

	l.OutputFunc(format, args...)
}

// MockImportService mocks the ImportService for testing.
type MockImportService struct {
	FromJSONFunc   func(r io.Reader) (map[string]string, error)
	FromDotEnvFunc func(r io.Reader) (map[string]string, error)
}

// FromJSON mocks the FromJSON method.
func (m *MockImportService) FromJSON(r io.Reader) (map[string]string, error) {
	if m.FromJSONFunc != nil {
		return m.FromJSONFunc(r)
	}
	return make(map[string]string), nil
}

// FromDotEnv mocks the FromDotEnv method.
func (m *MockImportService) FromDotEnv(r io.Reader) (map[string]string, error) {
	if m.FromDotEnvFunc != nil {
		return m.FromDotEnvFunc(r)
	}
	return make(map[string]string), nil
}

// MockVaultRepository mocks the VaultRepository for testing.
type MockVaultRepository struct {
	CreateFunc func(ctx context.Context, vault *model.Vault) error
	LoadFunc   func(ctx context.Context, env string) (*model.Vault, error)
	SaveFunc   func(ctx context.Context, vault *model.Vault) error
	ExistsFunc func(ctx context.Context, env string) (bool, error)
}

// Create mocks the Create method.
func (m *MockVaultRepository) Create(ctx context.Context, vault *model.Vault) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, vault)
	}
	return nil
}

// Load mocks the Load method.
func (m *MockVaultRepository) Load(ctx context.Context, env string) (*model.Vault, error) {
	if m.LoadFunc != nil {
		return m.LoadFunc(ctx, env)
	}
	vault, _ := model.NewVault(env, "test-fingerprint", "test-salt")
	return vault, nil
}

// Save mocks the Save method.
func (m *MockVaultRepository) Save(ctx context.Context, vault *model.Vault) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, vault)
	}
	return nil
}

// Exists mocks the Exists method.
func (m *MockVaultRepository) Exists(ctx context.Context, env string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, env)
	}
	return false, nil
}

// MockHashService mocks the HashService for testing.
type MockHashService struct {
	HashFunc         func(passphrase string) (string, error)
	VerifyFunc       func(hashedPassphrase, passphrase string) error
	GenerateSaltFunc func(size int) (string, error)
}

// Hash mocks the Hash method.
func (m *MockHashService) Hash(passphrase string) (string, error) {
	if m.HashFunc != nil {
		return m.HashFunc(passphrase)
	}
	return "test-fingerprint", nil
}

// Verify mocks the Verify method.
func (m *MockHashService) Verify(hashedPassphrase, passphrase string) error {
	if m.VerifyFunc != nil {
		return m.VerifyFunc(hashedPassphrase, passphrase)
	}
	return nil
}

// GenerateSalt mocks the GenerateSalt method.
func (m *MockHashService) GenerateSalt(size int) (string, error) {
	if m.GenerateSaltFunc != nil {
		return m.GenerateSaltFunc(size)
	}
	return "test-salt", nil
}

// MockPassphraseService mocks the PassphraseService for testing.
type MockPassphraseService struct {
	GetFunc      func(ctx context.Context, env string) (string, error)
	ClearFunc    func(ctx context.Context, env string) error
	ClearAllFunc func(ctx context.Context) error
	ValidateFunc func(ctx context.Context, vault *model.Vault, passphrase string) error
}

// Get mocks the Get method.
func (m *MockPassphraseService) Get(ctx context.Context, env string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, env)
	}
	return "test-passphrase", nil
}

// Clear mocks the Clear method.
func (m *MockPassphraseService) Clear(ctx context.Context, env string) error {
	if m.ClearFunc != nil {
		return m.ClearFunc(ctx, env)
	}
	return nil
}

// ClearAll mocks the ClearAll method.
func (m *MockPassphraseService) ClearAll(ctx context.Context) error {
	if m.ClearAllFunc != nil {
		return m.ClearAllFunc(ctx)
	}
	return nil
}

// Validate mocks the Validate method.
func (m *MockPassphraseService) Validate(
	ctx context.Context,
	vault *model.Vault,
	passphrase string,
) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(ctx, vault, passphrase)
	}
	return nil
}
