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
func (osKeyRing *OSKeyring) Set(service, key, value string) error {
	return keyring.Set(service, key, value)
}

// Get retrieves a value from the keyring
func (osKeyRing *OSKeyring) Get(service, key string) (string, error) {
	return keyring.Get(service, key)
}

// Delete removes a value from the keyring
func (osKeyRing *OSKeyring) Delete(service, key string) error {
	return keyring.Delete(service, key)
}

// DeleteAll removes all values for a service from the keyring
func (osKeyRing *OSKeyring) DeleteAll(service string) error {
	return keyring.DeleteAll(service)
}
