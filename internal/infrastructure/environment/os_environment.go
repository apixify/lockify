package environment

import (
	"os"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// OSEnvironmentProvider implements EnvironmentProvider using the OS environment.
type OSEnvironmentProvider struct{}

// NewOSEnvironmentProvider creates a new OS environment provider.
func NewOSEnvironmentProvider() service.EnvironmentProvider {
	return &OSEnvironmentProvider{}
}

// GetPassphrase retrieves a passphrase from an environment variable.
func (p *OSEnvironmentProvider) GetPassphrase(envVar string) string {
	return os.Getenv(envVar)
}
