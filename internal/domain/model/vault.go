package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/contract"
)

// Vault represents an encrypted vault containing entries for an environment.
type Vault struct {
	Meta       Meta             `json:"meta"`
	Entries    map[string]Entry `json:"entries"`
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
		Meta: Meta{
			Env:         env,
			Salt:        salt,
			FingerPrint: fingerprint,
		},
		Entries: make(map[string]Entry),
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

// Passphrase returns the vault passphrase
func (v *Vault) Passphrase() string {
	return v.passphrase
}

// SetPassphrase validates the passphrase against the fingerprint and sets it if valid.
// This ensures the passphrase matches the vault's fingerprint before allowing access.
func (v *Vault) SetPassphrase(passphrase string, hashService contract.HashService) error {
	if passphrase == "" {
		return errors.New("passphrase cannot be empty")
	}
	if v.Meta.FingerPrint == "" {
		return errors.New("vault fingerprint is not set")
	}
	if err := hashService.Verify(v.Meta.FingerPrint, passphrase); err != nil {
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
	entry, exists := v.Entries[key]
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
	entry, exists := v.Entries[key]

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

	v.Entries[key] = entry
	return nil
}

// DeleteEntry removes an entry by key
func (v *Vault) DeleteEntry(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if _, exists := v.Entries[key]; !exists {
		return fmt.Errorf("key %q not found", key)
	}
	delete(v.Entries, key)
	return nil
}

// ListKeys returns all keys in the vault
func (v *Vault) ListKeys() []string {
	keys := make([]string, 0, len(v.Entries))
	for k := range v.Entries {
		keys = append(keys, k)
	}
	return keys
}

// ForEachEntry iterates over all entries in the vault, calling fn for each entry.
// This provides controlled access to entries without exposing the internal map.
func (v *Vault) ForEachEntry(fn func(key string, entry Entry) error) error {
	for key, entry := range v.Entries {
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

	currentSalt := v.Meta.Salt
	for key, entry := range v.Entries {
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
		v.Entries[key] = entry
	}

	newFingerprint, err := hashService.Hash(newPassphrase)
	if err != nil {
		return fmt.Errorf("failed to hash new passphrase: %w", err)
	}

	v.Meta.Salt = newSalt
	v.Meta.FingerPrint = newFingerprint
	v.passphrase = newPassphrase

	return nil
}
