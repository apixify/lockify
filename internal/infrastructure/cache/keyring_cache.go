package cache

import (
	"github.com/apixify/lockify/internal/domain/service"
	"github.com/zalando/go-keyring"
)

// OSKeyring implements Cache using the OS keyring
type OSKeyring struct {
	service string
}

// NewOSKeyring creates a new OS keyring implementation
func NewOSKeyring(service string) service.Cache {
	return &OSKeyring{service: service}
}

// Set stores a value in the keyring
func (osKeyRing *OSKeyring) Set(key, value string) error {
	return keyring.Set(osKeyRing.service, key, value)
}

// Get retrieves a value from the keyring
func (osKeyRing *OSKeyring) Get(key string) (string, error) {
	return keyring.Get(osKeyRing.service, key)
}

// Delete removes a value from the keyring
func (osKeyRing *OSKeyring) Delete(key string) error {
	return keyring.Delete(osKeyRing.service, key)
}

// DeleteAll removes all values for a service from the keyring
func (osKeyRing *OSKeyring) DeleteAll() error {
	return keyring.DeleteAll(osKeyRing.service)
}
