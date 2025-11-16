package model

import (
	"errors"
	"fmt"
	"time"
)

type Vault struct {
	Meta    Meta             `json:"meta"`
	Entries map[string]Entry `json:"entries"`
	path    string
}

// Path returns the vault file path
func (v *Vault) Path() string {
	return v.path
}

// SetPath sets the vault file path
func (v *Vault) SetPath(path string) {
	v.path = path
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
