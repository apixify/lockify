package security

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// PassphraseService implements service.PassphraseService
type PassphraseService struct {
	cache      service.Cache
	cryptoUtil service.HashService
	envVar     string
}

// NewPassphraseService creates a new passphrase service
func NewPassphraseService(
	cache service.Cache,
	cryptoUtil service.HashService,
	envVar string,
) service.PassphraseService {
	if envVar == "" {
		envVar = "LOCKIFY_PASSPHRASE"
	}

	return &PassphraseService{cache, cryptoUtil, envVar}
}

// Get retrieves a passphrase from environment variable, keyring cache, or user input
func (service *PassphraseService) Get(ctx context.Context, env string) (string, error) {
	if env == "" {
		return "", fmt.Errorf("environment cannot be empty")
	}

	if passphrase := os.Getenv(service.envVar); passphrase != "" {
		return passphrase, nil
	}

	key := service.getKeyringKey(env)
	passphrase, err := service.cache.Get(key)
	if err == nil && passphrase != "" {
		return passphrase, nil
	}

	return service.getFromUser(ctx, env)
}

// Clear clears a cached passphrase for an environment
func (service *PassphraseService) Clear(ctx context.Context, env string) error {
	if env == "" {
		return fmt.Errorf("environment cannot be empty")
	}
	key := service.getKeyringKey(env)

	return service.cache.Delete(key)
}

// ClearAll clears all cached passphrases
func (service *PassphraseService) ClearAll(ctx context.Context) error {
	return service.cache.DeleteAll()
}

// Validate validates a passphrase against a vault's fingerprint
func (service *PassphraseService) Validate(
	ctx context.Context,
	vault *model.Vault,
	passphrase string,
) error {
	if vault == nil {
		return fmt.Errorf("vault cannot be nil")
	}
	if vault.Meta.FingerPrint == "" {
		return fmt.Errorf("fingerprint cannot be empty")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	return service.cryptoUtil.Verify(vault.Meta.FingerPrint, passphrase)
}

// getFromUser prompts the user for a passphrase
func (s *PassphraseService) getFromUser(ctx context.Context, env string) (string, error) {
	var passphrase string
	prompt := &survey.Password{
		Message: fmt.Sprintf("Enter passphrase for environment %q:", env),
	}

	if err := survey.AskOne(prompt, &passphrase); err != nil {
		return "", fmt.Errorf("failed to get passphrase: %w", err)
	}

	if passphrase == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	// Cache passphrase in keyring (best effort, ignore errors)
	key := s.getKeyringKey(env)
	//nolint:errcheck // We don't want to return an error here
	s.cache.Set(key, passphrase)

	return passphrase, nil
}

// getKeyringKey returns the keyring key for an environment
func (service *PassphraseService) getKeyringKey(env string) string {
	return fmt.Sprintf("env:%s", env)
}
