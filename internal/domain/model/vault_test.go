package model

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

const (
	testEnv         = "test"
	testFingerprint = "test-fingerprint"
	testPassphrase  = "test-passphrase"
	testSalt        = "test-salt"
	testKey         = "test-key"
	testValue       = "test-value"
)

func createTestVault(t *testing.T) *Vault {
	t.Helper()
	vault, err := NewVault(testEnv, testFingerprint, testSalt)
	if err != nil {
		t.Fatalf("failed to create test vault: %v", err)
	}
	if vault == nil {
		t.Fatal("NewVault() should return a vault, got nil")
	}

	return vault
}

func TestNewVault(t *testing.T) {
	vault := createTestVault(t)

	if vault.Env() != testEnv {
		t.Errorf("expected env %q, got %q", testEnv, vault.Env())
	}
	if vault.FingerPrint() != testFingerprint {
		t.Errorf("expected fingerprint %q, got %q", testFingerprint, vault.FingerPrint())
	}
	if vault.Salt() != testSalt {
		t.Errorf("expected salt %q, got %q", testSalt, vault.Salt())
	}
	if vault.EntriesCount() != 0 {
		t.Errorf("expected 0 entries, got %d", vault.EntriesCount())
	}
}

func TestSetEntry(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry(testKey, testValue)
	if vault.EntriesCount() != 1 {
		t.Errorf("expected 1 entry, got %d", vault.EntriesCount())
	}
	entry, err := vault.GetEntry(testKey)
	if err != nil {
		t.Fatalf("failed to get entry: %v", err)
	}
	if entry.Value != testValue {
		t.Errorf("expected value %q, got %q", testValue, entry.Value)
	}
	if entry.CreatedAt == "" {
		t.Errorf("expected created at, got empty")
	}
}

func TestGetEntry(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry(testKey, testValue)
	if vault.EntriesCount() != 1 {
		t.Errorf("expected 1 entry, got %d", vault.EntriesCount())
	}

	entry, err := vault.GetEntry(testKey)
	if err != nil {
		t.Fatalf("failed to get entry: %v", err)
	}
	if entry.Value != testValue {
		t.Errorf("expected value %q, got %q", testValue, entry.Value)
	}
	if entry.CreatedAt == "" {
		t.Errorf("expected created at, got empty")
	}
	if entry.UpdatedAt == "" {
		t.Errorf("expected updated at, got empty")
	}
}

func TestSetEntryUpdateKey(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry(testKey, testValue)
	if vault.EntriesCount() != 1 {
		t.Errorf("expected 1 entry, got %d", vault.EntriesCount())
	}
	// Get the first entry to capture its CreatedAt timestamp
	firstEntry, err := vault.GetEntry(testKey)
	if err != nil {
		t.Fatalf("failed to get entry: %v", err)
	}
	firstCreatedAt := firstEntry.CreatedAt
	if firstCreatedAt == "" {
		t.Fatalf("expected CreatedAt to be set")
	}

	// Update the entry - CreatedAt should remain the same, UpdatedAt should change
	testValue2 := "test-value-2"
	vault.SetEntry(testKey, testValue2)
	entry, err := vault.GetEntry(testKey)
	if err != nil {
		t.Fatalf("failed to get entry: %v", err)
	}
	if entry.Value != testValue2 {
		t.Errorf("expected value %q, got %q", testValue2, entry.Value)
	}
	if entry.CreatedAt != firstCreatedAt {
		t.Errorf("expected CreatedAt to remain %q, got %q", firstCreatedAt, entry.CreatedAt)
	}
	if entry.UpdatedAt == "" {
		t.Errorf("expected UpdatedAt to be set, got empty")
	}
	if entry.Value != testValue2 {
		t.Errorf("expected value %q, got %q", testValue2, entry.Value)
	}
}

func TestGetNonExistentEntry(t *testing.T) {
	vault := createTestVault(t)

	_, err := vault.GetEntry("non_existent_key")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key \"non_existent_key\" not found" {
		t.Errorf("expected error %q, got %q", "key \"non_existent_key\" not found", err.Error())
	}
}

func TestGetEntryWithEmptyKey(t *testing.T) {
	vault := createTestVault(t)

	_, err := vault.GetEntry("")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key cannot be empty" {
		t.Errorf("expected error %q, got %q", "key cannot be empty", err.Error())
	}
}

func TestDeleteEntry(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry(testKey, testValue)
	if vault.EntriesCount() != 1 {
		t.Errorf("expected 1 entry, got %d", vault.EntriesCount())
	}

	vault.DeleteEntry(testKey)
	if vault.EntriesCount() != 0 {
		t.Errorf("expected 0 entries, got %d", vault.EntriesCount())
	}

	_, err := vault.GetEntry(testKey)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != fmt.Sprintf("key %q not found", testKey) {
		t.Errorf("expected error %q, got %q", fmt.Sprintf("key %q not found", testKey), err.Error())
	}
}

func TestDeleteNonExistentEntry(t *testing.T) {
	vault := createTestVault(t)

	err := vault.DeleteEntry("non_existent_key")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key \"non_existent_key\" not found" {
		t.Errorf("expected error %q, got %q", "key \"non_existent_key\" not found", err.Error())
	}
}

func TestListKeys(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry("test_key", "test_value")
	if vault.EntriesCount() != 1 {
		t.Errorf("expected 1 entry, got %d", vault.EntriesCount())
	}

	keys := vault.ListKeys()
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
	if keys[0] != "test_key" {
		t.Errorf("expected key %q, got %q", "test_key", keys[0])
	}
}

func TestListKeysEmpty(t *testing.T) {
	vault := createTestVault(t)

	keys := vault.ListKeys()
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestListKeysMultiple(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry(testKey, testValue)
	vault.SetEntry("test_key2", "test_value2")

	keys := vault.ListKeys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
	if !slices.Contains(keys, testKey) {
		t.Errorf("expected keys to contain %q", testKey)
	}
	if !slices.Contains(keys, "test_key2") {
		t.Errorf("expected keys to contain %q", "test_key2")
	}
}

func TestErrorSetEntryWithEmptyKey(t *testing.T) {
	vault := createTestVault(t)

	err := vault.SetEntry("", testValue)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key cannot be empty" {
		t.Errorf("expected error %q, got %q", "key cannot be empty", err.Error())
	}
}

func TestSetEntryWithEmptyValue(t *testing.T) {
	vault := createTestVault(t)

	err := vault.SetEntry(testKey, "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "encrypted value cannot be empty" {
		t.Errorf("expected error %q, got %q", "encrypted value cannot be empty", err.Error())
	}
}

func TestSetPassphrase(t *testing.T) {
	vault := createTestVault(t)

	// Use SetPassphrase with a mock hash service that always succeeds
	hashService := &mockHashService{}
	if err := vault.SetPassphrase(testPassphrase, hashService); err != nil {
		t.Fatalf("SetPassphrase() error = %v, want nil", err)
	}
	if vault.Passphrase() != testPassphrase {
		t.Errorf("expected passphrase %q, got %q", testPassphrase, vault.Passphrase())
	}
}

func TestSetPath(t *testing.T) {
	vault := createTestVault(t)

	vault.SetPath("test_path")
	if vault.Path() != "test_path" {
		t.Errorf("expected path %q, got %q", "test_path", vault.Path())
	}
}

func TestNewVault_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		fingerprint string
		salt        string
		wantErr     string
	}{
		{
			name:        "empty env",
			env:         "",
			fingerprint: "test",
			salt:        "test",
			wantErr:     "environment cannot be empty",
		},
		{
			name:        "empty fingerprint",
			env:         "test",
			fingerprint: "",
			salt:        "test",
			wantErr:     "fingerprint cannot be empty",
		},
		{
			name:        "empty salt",
			env:         "test",
			fingerprint: "test",
			salt:        "",
			wantErr:     "salt cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVault(tt.env, tt.fingerprint, tt.salt)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error to contain %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// mockHashService is a simple mock for testing that always succeeds verification
type mockHashService struct{}

func (m *mockHashService) Hash(passphrase string) (string, error) {
	return "hashed-" + passphrase, nil
}

func (m *mockHashService) Verify(hashedPassphrase, passphrase string) error {
	return nil // Always succeeds for testing
}

func (m *mockHashService) GenerateSalt(size int) (string, error) {
	return "test-salt", nil
}
