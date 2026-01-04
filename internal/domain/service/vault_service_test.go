package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
)

// ============================================================================
// Helpers
// ============================================================================
func createTestVault(env string) *model.Vault {
	vault, _ := model.NewVault(env, "test-fingerprint", "test-salt")
	vault.SetEntry("test-entry", "test-value")
	return vault
}

func createVaultServiceWithMocks(
	repo *test.MockVaultRepository,
	passphrase *test.MockPassphraseService,
	hash *test.MockHashService,
) VaultServiceInterface {
	return NewVaultService(repo, passphrase, hash)
}

// ============================================================================
// Tests
// ============================================================================
func TestCreate_Success(t *testing.T) {
	vaultService := createVaultServiceWithMocks(
		&test.MockVaultRepository{},
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)
	vault, err := vaultService.Create(model.NewVaultContext(context.Background(), "test", false))
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
	repo := &test.MockVaultRepository{
		ExistsFunc: func(vctx *model.VaultContext) (bool, error) {
			return true, nil
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	_, err := vaultService.Create(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Create() with existing vault expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault already exists") {
		t.Errorf("Create() error = %q, want to contain 'vault already exists'", err.Error())
	}
}

func TestCreate_RepositoryExistsError(t *testing.T) {
	repo := &test.MockVaultRepository{
		ExistsFunc: func(vctx *model.VaultContext) (bool, error) {
			return false, errors.New("repository error")
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	_, err := vaultService.Create(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Create() with repository error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to check vault existence") {
		t.Errorf(
			"Create() error = %q, want to contain 'failed to check vault existence'",
			err.Error(),
		)
	}
}

func TestCreate_PassphraseGetError(t *testing.T) {
	passphrase := &test.MockPassphraseService{
		GetWithConfirmationFunc: func(vctx *model.VaultContext) (string, error) {
			return "", errors.New("passphrase error")
		},
	}
	vaultService := createVaultServiceWithMocks(
		&test.MockVaultRepository{},
		passphrase,
		&test.MockHashService{},
	)

	_, err := vaultService.Create(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Create() with passphrase error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get passphrase") {
		t.Errorf("Create() error = %q, want to contain 'failed to get passphrase'", err.Error())
	}
}

func TestCreate_HashError(t *testing.T) {
	hash := &test.MockHashService{
		HashFunc: func(passphrase string) (string, error) {
			return "", errors.New("hash error")
		},
	}
	vaultService := createVaultServiceWithMocks(
		&test.MockVaultRepository{},
		&test.MockPassphraseService{},
		hash,
	)

	_, err := vaultService.Create(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Create() with hash error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to hash passphrase") {
		t.Errorf("Create() error = %q, want to contain 'failed to hash passphrase'", err.Error())
	}
}

func TestCreate_GenerateSaltError(t *testing.T) {
	hash := &test.MockHashService{
		GenerateSaltFunc: func(length int) (string, error) {
			return "", errors.New("salt error")
		},
	}
	vaultService := createVaultServiceWithMocks(
		&test.MockVaultRepository{},
		&test.MockPassphraseService{},
		hash,
	)

	_, err := vaultService.Create(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Create() with salt error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to generate salt") {
		t.Errorf("Create() error = %q, want to contain 'failed to generate salt'", err.Error())
	}
}

func TestCreate_RepositoryCreateError(t *testing.T) {
	repo := &test.MockVaultRepository{
		CreateFunc: func(vctx *model.VaultContext, vault *model.Vault) error {
			return errors.New("create error")
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	_, err := vaultService.Create(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Create() with repository create error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to save vault") {
		t.Errorf("Create() error = %q, want to contain 'failed to save vault'", err.Error())
	}
}

func TestOpen_Success(t *testing.T) {
	testVault := createTestVault("test")
	repo := &test.MockVaultRepository{
		ExistsFunc: func(vctx *model.VaultContext) (bool, error) {
			return true, nil
		},
		LoadFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			return testVault, nil
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	vault, err := vaultService.Open(model.NewVaultContext(context.Background(), "test", false))
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
	repo := &test.MockVaultRepository{
		ExistsFunc: func(vctx *model.VaultContext) (bool, error) {
			return false, nil
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	_, err := vaultService.Open(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Open() with non-existent vault expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault for env") ||
		!strings.Contains(err.Error(), "does not exist") {
		t.Errorf(
			"Open() error = %q, want to contain 'vault for env' and 'does not exist'",
			err.Error(),
		)
	}
}

func TestOpen_PassphraseGetError(t *testing.T) {
	repo := &test.MockVaultRepository{
		ExistsFunc: func(vctx *model.VaultContext) (bool, error) {
			return true, nil
		},
	}
	passphrase := &test.MockPassphraseService{
		GetFunc: func(vctx *model.VaultContext) (string, error) {
			return "", errors.New("passphrase error")
		},
	}
	vaultService := createVaultServiceWithMocks(repo, passphrase, &test.MockHashService{})

	_, err := vaultService.Open(model.NewVaultContext(context.Background(), "test", false))
	if err == nil {
		t.Fatal("Open() with passphrase error expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to retrieve passphrase") {
		t.Errorf("Open() error = %q, want to contain 'failed to retrieve passphrase'", err.Error())
	}
}

func TestOpen_RepositoryLoadError(t *testing.T) {
	repo := &test.MockVaultRepository{
		ExistsFunc: func(vctx *model.VaultContext) (bool, error) {
			return true, nil
		},
		LoadFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			return nil, errors.New("load error")
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	_, err := vaultService.Open(model.NewVaultContext(context.Background(), "test", false))
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
	repo := &test.MockVaultRepository{
		ExistsFunc: func(vctx *model.VaultContext) (bool, error) {
			return true, nil
		},
		LoadFunc: func(vctx *model.VaultContext) (*model.Vault, error) {
			return testVault, nil
		},
	}
	passphrase := &test.MockPassphraseService{
		ValidateFunc: func(vctx *model.VaultContext, vault *model.Vault, passphrase string) error {
			return errors.New("invalid passphrase")
		},
		ClearFunc: func(vctx *model.VaultContext) error {
			clearCalled = true
			return nil
		},
	}
	vaultService := createVaultServiceWithMocks(repo, passphrase, &test.MockHashService{})

	_, err := vaultService.Open(model.NewVaultContext(context.Background(), "test", false))
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
	repo := &test.MockVaultRepository{
		SaveFunc: func(vctx *model.VaultContext, vault *model.Vault) error {
			saveCalled = true
			if vault == nil {
				return errors.New("vault is nil")
			}
			return nil
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	err := vaultService.Save(model.NewVaultContext(context.Background(), "test", false), vault)
	if err != nil {
		t.Fatalf("Save() returned unexpected error: %v", err)
	}
	if !saveCalled {
		t.Error("Save() should call repository.Save(), but it didn't")
	}
}

func TestSave_RepositoryError(t *testing.T) {
	vault := createTestVault("test")
	repo := &test.MockVaultRepository{
		SaveFunc: func(vctx *model.VaultContext, vault *model.Vault) error {
			return errors.New("save error")
		},
	}
	vaultService := createVaultServiceWithMocks(
		repo,
		&test.MockPassphraseService{},
		&test.MockHashService{},
	)

	err := vaultService.Save(model.NewVaultContext(context.Background(), "test", false), vault)
	if err == nil {
		t.Fatal("Save() with repository error expected error, got nil")
	}
	if err.Error() != "save error" {
		t.Errorf("Save() error = %q, want %q", err.Error(), "save error")
	}
}
