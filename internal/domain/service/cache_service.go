package service

type Cache interface {
	// Set stores a value in cache
	Set(key, value string) error
	// Get retrieves a value from cache
	Get(key string) (string, error)
	// Delete removes a value from cache
	Delete(key string) error
	// DeleteAll removes all values for a service from cache
	DeleteAll() error
}
