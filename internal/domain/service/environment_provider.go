package service

// EnvironmentProvider provides access to environment variables.
type EnvironmentProvider interface {
	// GetPassphrase retrieves a passphrase from an environment variable.
	// Returns empty string if the variable is not set.
	GetPassphrase(envVar string) string
}
