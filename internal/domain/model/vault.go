package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/contract"
)

// Vault represents an encrypted vault containing entries for an environment.
// Vault must be created using NewVault() to ensure proper initialization.
// Fields are private to enforce encapsulation and prevent invalid state.
type Vault struct {
	meta       Meta
	entries    map[string]Entry
	path       string
	passphrase string
}

// NewVault creates a new vault instance
func NewVault(env, fingerprint, salt string) (*Vault, error) {
	if env == "" {
		return nil, errors.New("environment cannot be empty")
	}
	if fingerprint == "" {
		return nil, errors.New("fingerprint cannot be empty")
	}
	if salt == "" {
		return nil, errors.New("salt cannot be empty")
	}

	vault := &Vault{
		meta: Meta{
			Env:         env,
			Salt:        salt,
			FingerPrint: fingerprint,
		},
		entries: make(map[string]Entry),
	}

	return vault, nil
}

// Path returns the vault file path
func (v *Vault) Path() string {
	return v.path
}

// SetPath sets the vault file path
func (v *Vault) SetPath(path string) {
	v.path = path
}

// Meta returns a copy of the vault metadata.
// Returns a copy to prevent external mutation of the vault's internal state.
func (v *Vault) Meta() Meta {
	return v.meta
}

// Salt returns the salt used for encryption key derivation.
func (v *Vault) Salt() string {
	return v.meta.Salt
}

// Env returns the environment name for this vault.
func (v *Vault) Env() string {
	return v.meta.Env
}

// FingerPrint returns the fingerprint hash of the passphrase.
func (v *Vault) FingerPrint() string {
	return v.meta.FingerPrint
}

// EntriesCount returns the number of entries in the vault.
// This is safer than exposing the entries map directly.
func (v *Vault) EntriesCount() int {
	return len(v.entries)
}

// Passphrase returns the vault passphrase.
// WARNING: This returns the passphrase stored in memory. For security, call ClearPassphrase()
// when the passphrase is no longer needed.
func (v *Vault) Passphrase() string {
	return v.passphrase
}

// ClearPassphrase clears the passphrase from memory for security.
// This removes the reference to the passphrase string.
// Note: Due to Go's garbage collector and string immutability, complete memory clearing
// cannot be guaranteed, but this removes the reference and allows GC to reclaim memory.
// Should be called when passphrase is no longer needed (e.g., after operations complete).
// SECURITY WARNING: Passphrases stored in Go strings cannot be reliably zeroed from memory.
// For maximum security, consider using byte slices and clearing them explicitly, or
// limiting the lifetime of Vault instances that hold passphrases.
func (v *Vault) ClearPassphrase() {
	v.passphrase = ""
}

// SetPassphrase validates the passphrase against the fingerprint and sets it if valid.
// This ensures the passphrase matches the vault's fingerprint before allowing access.
func (v *Vault) SetPassphrase(passphrase string, hashService contract.HashService) error {
	if hashService == nil {
		return errors.New("hash service is required")
	}
	if passphrase == "" {
		return errors.New("passphrase cannot be empty")
	}
	if v.meta.FingerPrint == "" {
		return errors.New("vault fingerprint is not set")
	}
	if err := hashService.Verify(v.meta.FingerPrint, passphrase); err != nil {
		return fmt.Errorf("passphrase does not match vault fingerprint: %w", err)
	}
	v.passphrase = passphrase
	return nil
}

// GetEntry retrieves an entry by key
func (v *Vault) GetEntry(key string) (Entry, error) {
	if key == "" {
		return Entry{}, errors.New("key cannot be empty")
	}
	entry, exists := v.entries[key]
	if !exists {
		return Entry{}, fmt.Errorf("key %q not found", key)
	}
	return entry, nil
}

// SetEntry adds or updates an entry
func (v *Vault) SetEntry(key, encryptedValue string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if encryptedValue == "" {
		return errors.New("encrypted value cannot be empty")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	entry, exists := v.entries[key]

	if exists {
		entry.Value = encryptedValue
		entry.UpdatedAt = now
	} else {
		entry = Entry{
			Value:     encryptedValue,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	v.entries[key] = entry
	return nil
}

// DeleteEntry removes an entry by key
func (v *Vault) DeleteEntry(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if _, exists := v.entries[key]; !exists {
		return fmt.Errorf("key %q not found", key)
	}
	delete(v.entries, key)
	return nil
}

// ListKeys returns all keys in the vault
func (v *Vault) ListKeys() []string {
	keys := make([]string, 0, len(v.entries))
	for k := range v.entries {
		keys = append(keys, k)
	}
	return keys
}

// ForEachEntry iterates over all entries in the vault, calling fn for each entry.
// This provides controlled access to entries without exposing the internal map.
func (v *Vault) ForEachEntry(fn func(key string, entry Entry) error) error {
	for key, entry := range v.entries {
		if err := fn(key, entry); err != nil {
			return err
		}
	}
	return nil
}

// RotatePassphrase rotates the vault's passphrase by re-encrypting all entries
// with the new passphrase and updating the vault's metadata.
func (v *Vault) RotatePassphrase(
	currentPassphrase, newPassphrase string,
	encryptionService contract.EncryptionService,
	hashService contract.HashService,
	saltSize int,
) error {
	if encryptionService == nil {
		return errors.New("encryption service is required")
	}
	if hashService == nil {
		return errors.New("hash service is required")
	}
	if currentPassphrase == "" {
		return errors.New("current passphrase cannot be empty")
	}
	if newPassphrase == "" {
		return errors.New("new passphrase cannot be empty")
	}
	if err := v.SetPassphrase(currentPassphrase, hashService); err != nil {
		return fmt.Errorf("invalid current passphrase: %w", err)
	}

	newSalt, err := hashService.GenerateSalt(saltSize)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	currentSalt := v.meta.Salt
	if err := v.rotateEntries(currentSalt, currentPassphrase, newSalt, newPassphrase, encryptionService); err != nil {
		return err
	}

	newFingerprint, err := hashService.Hash(newPassphrase)
	if err != nil {
		return fmt.Errorf("failed to hash new passphrase: %w", err)
	}

	v.meta.Salt = newSalt
	v.meta.FingerPrint = newFingerprint
	v.passphrase = newPassphrase

	return nil
}

// rotateEntries re-encrypts all entries with new salt and passphrase.
// This is a helper method extracted from RotatePassphrase for better organization.
func (v *Vault) rotateEntries(
	currentSalt, currentPassphrase string,
	newSalt, newPassphrase string,
	encryptionService contract.EncryptionService,
) error {
	for key, entry := range v.entries {
		// Decrypt with current passphrase
		decryptedValue, err := encryptionService.Decrypt(
			entry.Value,
			currentSalt,
			currentPassphrase,
		)
		if err != nil {
			return fmt.Errorf("failed to decrypt entry %q: %w", key, err)
		}

		encryptedValue, err := encryptionService.Encrypt(
			decryptedValue,
			newSalt,
			newPassphrase,
		)
		if err != nil {
			return fmt.Errorf("failed to encrypt entry %q: %w", key, err)
		}

		entry.Value = encryptedValue
		entry.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		v.entries[key] = entry
	}
	return nil
}

// vaultJSON is a helper struct for JSON marshaling/unmarshaling.
// This allows us to maintain compatibility with existing vault file format
// while keeping the internal fields private.
type vaultJSON struct {
	Meta    Meta             `json:"meta"`
	Entries map[string]Entry `json:"entries"`
}

// MarshalJSON implements json.Marshaler for Vault.
// This ensures private fields are properly serialized to maintain
// compatibility with existing vault file format.
func (v *Vault) MarshalJSON() ([]byte, error) {
	return json.Marshal(vaultJSON{
		Meta:    v.meta,
		Entries: v.entries,
	})
}

// UnmarshalJSON implements json.Unmarshaler for Vault.
// This ensures private fields are properly deserialized from
// existing vault file format.
func (v *Vault) UnmarshalJSON(data []byte) error {
	var vj vaultJSON
	if err := json.Unmarshal(data, &vj); err != nil {
		return err
	}

	v.meta = vj.Meta
	v.entries = vj.Entries
	if v.entries == nil {
		v.entries = make(map[string]Entry)
	}

	return nil
}
