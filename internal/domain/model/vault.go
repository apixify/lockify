package model

import (
	"errors"
	"fmt"
	"time"
)

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

// SetPassphrase sets the vault passphrase
func (v *Vault) SetPassphrase(passphrase string) {
	v.passphrase = passphrase
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
func (v *Vault) SetEntry(key string, encryptedValue string) error {
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

	if v.Entries == nil {
		v.Entries = make(map[string]Entry)
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
