package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

// ============================================================================
// Mocks
// ============================================================================
type mockVaultRepository struct {
	ExistsFunc func(ctx context.Context, env string) (bool, error)
	CreateFunc func(ctx context.Context, vault *model.Vault) error
	LoadFunc   func(ctx context.Context, env string) (*model.Vault, error)
	SaveFunc   func(ctx context.Context, vault *model.Vault) error
}

func (m *mockVaultRepository) Exists(ctx context.Context, env string) (bool, error) {
	if m.ExistsFunc == nil {
		return false, nil
	}
	return m.ExistsFunc(ctx, env)
}

func (m *mockVaultRepository) Create(ctx context.Context, vault *model.Vault) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, vault)
}

func (m *mockVaultRepository) Load(ctx context.Context, env string) (*model.Vault, error) {
	if m.LoadFunc == nil {
		return nil, nil
	}
	return m.LoadFunc(ctx, env)
}

func (m *mockVaultRepository) Save(ctx context.Context, vault *model.Vault) error {
	if m.SaveFunc == nil {
		return nil
	}
	return m.SaveFunc(ctx, vault)
}

// ============================================================================
// Services
// ============================================================================
type mockPassphraseService struct {
	GetFunc      func(ctx context.Context, env string) (string, error)
	ValidateFunc func(ctx context.Context, vault *model.Vault, passphrase string) error
	ClearFunc    func(ctx context.Context, env string) error
	ClearAllFunc func(ctx context.Context) error
}

func (m *mockPassphraseService) Get(ctx context.Context, env string) (string, error) {
	if m.GetFunc == nil {
		return "test-passphrase", nil
	}
	return m.GetFunc(ctx, env)
}

func (m *mockPassphraseService) Validate(ctx context.Context, vault *model.Vault, passphrase string) error {
	if m.ValidateFunc == nil {
		return nil
	}
	return m.ValidateFunc(ctx, vault, passphrase)
}

func (m *mockPassphraseService) Clear(ctx context.Context, env string) error {
	if m.ClearFunc == nil {
		return nil
	}
	return m.ClearFunc(ctx, env)
}

func (m *mockPassphraseService) ClearAll(ctx context.Context) error {
	if m.ClearAllFunc == nil {
		return nil
	}
	return m.ClearAllFunc(ctx)
}

type mockHashService struct {
	HashFunc         func(passphrase string) (string, error)
	GenerateSaltFunc func(length int) (string, error)
	VerifyFunc       func(fingerprint, passphrase string) error
}

func (m *mockHashService) Hash(passphrase string) (string, error) {
	if m.HashFunc == nil {
		return "test-fingerprint", nil
	}
	return m.HashFunc(passphrase)
}

func (m *mockHashService) GenerateSalt(length int) (string, error) {
	if m.GenerateSaltFunc == nil {
		return "test-salt", nil
	}

	return m.GenerateSaltFunc(length)
}

func (m *mockHashService) Verify(hashedPassphrase, passphrase string) error {
	if m.VerifyFunc == nil {
		return nil
	}
	return m.VerifyFunc(hashedPassphrase, passphrase)
}

// ============================================================================
// Helpers
// ============================================================================
func createTestVault(env string) *model.Vault {
	vault, _ := model.NewVault(env, "test-fingerprint", "test-salt")
	vault.SetEntry("test-entry", "test-value")
	return vault
}

func createVaultServiceWithMocks(
	repo *mockVaultRepository,
	passphrase *mockPassphraseService,
	hash *mockHashService,
) VaultService {
	return NewVaultService(repo, passphrase, hash)
}

// ============================================================================
// Tests
// ============================================================================
func TestCreate_Success(t *testing.T) {
	vaultService := createVaultServiceWithMocks(&mockVaultRepository{}, &mockPassphraseService{}, &mockHashService{})
	vault, err := vaultService.Create(context.Background(), "test")
	if err != nil {
		t.Fatalf("Create() returned unexpected error: %v", err)
	}
	if vault == nil {
		t.Fatal("Create() returned nil vault")
	}
	if vault.Meta.Env != "test" {
		t.Errorf("Create() vault.Meta.Env = %q, want %q", vault.Meta.Env, "test")
	}
	if vault.Meta.FingerPrint == "" {
		t.Error("Create() vault.Meta.FingerPrint is empty")
	}
	if vault.Meta.Salt == "" {
		t.Error("Create() vault.Meta.Salt is empty")
	}
	if len(vault.Entries) != 0 {
		t.Errorf("Create() vault.Entries length = %d, want 0", len(vault.Entries))
	}
}

func TestCreate_VaultAlreadyExists(t *testing.T) {
	repo := &mockVaultRepository{
		ExistsFunc: func(ctx context.Context, env string) (bool, error) {
			return true, nil
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	_, err := vaultService.Create(context.Background(), "test")
	if err == nil {
		t.Fatal("Create() with existing vault expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault already exists") {
		t.Errorf("Create() error = %q, want to contain 'vault already exists'", err.Error())
	}
}

func TestCreate_RepositoryExistsError(t *testing.T) {
	repo := &mockVaultRepository{
		ExistsFunc: func(ctx context.Context, env string) (bool, error) {
			return false, errors.New("repository error")
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	_, err := vaultService.Create(context.Background(), "test")
	if err == nil {
		t.Fatal("Create() with repository error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to check vault existence") {
		t.Errorf("Create() error = %q, want to contain 'failed to check vault existence'", err.Error())
	}
}

func TestCreate_PassphraseGetError(t *testing.T) {
	passphrase := &mockPassphraseService{
		GetFunc: func(ctx context.Context, env string) (string, error) {
			return "", errors.New("passphrase error")
		},
	}
	vaultService := createVaultServiceWithMocks(&mockVaultRepository{}, passphrase, &mockHashService{})

	_, err := vaultService.Create(context.Background(), "test")
	if err == nil {
		t.Fatal("Create() with passphrase error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get passphrase") {
		t.Errorf("Create() error = %q, want to contain 'failed to get passphrase'", err.Error())
	}
}

func TestCreate_HashError(t *testing.T) {
	hash := &mockHashService{
		HashFunc: func(passphrase string) (string, error) {
			return "", errors.New("hash error")
		},
	}
	vaultService := createVaultServiceWithMocks(&mockVaultRepository{}, &mockPassphraseService{}, hash)

	_, err := vaultService.Create(context.Background(), "test")
	if err == nil {
		t.Fatal("Create() with hash error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to hash passphrase") {
		t.Errorf("Create() error = %q, want to contain 'failed to hash passphrase'", err.Error())
	}
}

func TestCreate_GenerateSaltError(t *testing.T) {
	hash := &mockHashService{
		GenerateSaltFunc: func(length int) (string, error) {
			return "", errors.New("salt error")
		},
	}
	vaultService := createVaultServiceWithMocks(&mockVaultRepository{}, &mockPassphraseService{}, hash)

	_, err := vaultService.Create(context.Background(), "test")
	if err == nil {
		t.Fatal("Create() with salt error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to generate salt") {
		t.Errorf("Create() error = %q, want to contain 'failed to generate salt'", err.Error())
	}
}

func TestCreate_RepositoryCreateError(t *testing.T) {
	repo := &mockVaultRepository{
		CreateFunc: func(ctx context.Context, vault *model.Vault) error {
			return errors.New("create error")
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	_, err := vaultService.Create(context.Background(), "test")
	if err == nil {
		t.Fatal("Create() with repository create error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to save vault") {
		t.Errorf("Create() error = %q, want to contain 'failed to save vault'", err.Error())
	}
}

func TestOpen_Success(t *testing.T) {
	testVault := createTestVault("test")
	repo := &mockVaultRepository{
		ExistsFunc: func(ctx context.Context, env string) (bool, error) {
			return true, nil
		},
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return testVault, nil
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	vault, err := vaultService.Open(context.Background(), "test")
	if err != nil {
		t.Fatalf("Open() returned unexpected error: %v", err)
	}
	if vault == nil {
		t.Fatal("Open() returned nil vault")
	}
	if vault.Meta.Env != "test" {
		t.Errorf("Open() vault.Meta.Env = %q, want %q", vault.Meta.Env, "test")
	}
	if vault.Passphrase() != "test-passphrase" {
		t.Errorf("Open() vault.Passphrase() = %q, want %q", vault.Passphrase(), "test-passphrase")
	}
}

func TestOpen_VaultDoesNotExist(t *testing.T) {
	repo := &mockVaultRepository{
		ExistsFunc: func(ctx context.Context, env string) (bool, error) {
			return false, nil
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	_, err := vaultService.Open(context.Background(), "test")
	if err == nil {
		t.Fatal("Open() with non-existent vault expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault for env") || !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Open() error = %q, want to contain 'vault for env' and 'does not exist'", err.Error())
	}
}

func TestOpen_PassphraseGetError(t *testing.T) {
	repo := &mockVaultRepository{
		ExistsFunc: func(ctx context.Context, env string) (bool, error) {
			return true, nil
		},
	}
	passphrase := &mockPassphraseService{
		GetFunc: func(ctx context.Context, env string) (string, error) {
			return "", errors.New("passphrase error")
		},
	}
	vaultService := createVaultServiceWithMocks(repo, passphrase, &mockHashService{})

	_, err := vaultService.Open(context.Background(), "test")
	if err == nil {
		t.Fatal("Open() with passphrase error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to retrieve passphrase") {
		t.Errorf("Open() error = %q, want to contain 'failed to retrieve passphrase'", err.Error())
	}
}

func TestOpen_RepositoryLoadError(t *testing.T) {
	repo := &mockVaultRepository{
		ExistsFunc: func(ctx context.Context, env string) (bool, error) {
			return true, nil
		},
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return nil, errors.New("load error")
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	_, err := vaultService.Open(context.Background(), "test")
	if err == nil {
		t.Fatal("Open() with load error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to open vault") {
		t.Errorf("Open() error = %q, want to contain 'failed to open vault'", err.Error())
	}
}

func TestOpen_InvalidPassphrase(t *testing.T) {
	testVault := createTestVault("test")
	clearCalled := false
	repo := &mockVaultRepository{
		ExistsFunc: func(ctx context.Context, env string) (bool, error) {
			return true, nil
		},
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return testVault, nil
		},
	}
	passphrase := &mockPassphraseService{
		ValidateFunc: func(ctx context.Context, vault *model.Vault, passphrase string) error {
			return errors.New("invalid passphrase")
		},
		ClearFunc: func(ctx context.Context, env string) error {
			clearCalled = true
			return nil
		},
	}
	vaultService := createVaultServiceWithMocks(repo, passphrase, &mockHashService{})

	_, err := vaultService.Open(context.Background(), "test")
	if err == nil {
		t.Fatal("Open() with invalid passphrase expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Errorf("Open() error = %q, want to contain 'invalid credentials'", err.Error())
	}
	if !clearCalled {
		t.Error("Open() with invalid passphrase should call Clear(), but it didn't")
	}
}

func TestSave_Success(t *testing.T) {
	vault := createTestVault("test")
	saveCalled := false
	repo := &mockVaultRepository{
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			saveCalled = true
			if vault == nil {
				return errors.New("vault is nil")
			}
			return nil
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	err := vaultService.Save(context.Background(), vault)
	if err != nil {
		t.Fatalf("Save() returned unexpected error: %v", err)
	}
	if !saveCalled {
		t.Error("Save() should call repository.Save(), but it didn't")
	}
}

func TestSave_RepositoryError(t *testing.T) {
	vault := createTestVault("test")
	repo := &mockVaultRepository{
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			return errors.New("save error")
		},
	}
	vaultService := createVaultServiceWithMocks(repo, &mockPassphraseService{}, &mockHashService{})

	err := vaultService.Save(context.Background(), vault)
	if err == nil {
		t.Fatal("Save() with repository error expected error, got nil")
	}
	if err.Error() != "save error" {
		t.Errorf("Save() error = %q, want %q", err.Error(), "save error")
	}
}
