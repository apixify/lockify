package security

import (
	"context"
	"fmt"
	"os"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// PassphraseService implements service.PassphraseService
type PassphraseService struct {
	cache      service.Cache
	cryptoUtil service.HashService
	prompt     service.PromptService
	envVar     string
}

// NewPassphraseService creates a new passphrase service
func NewPassphraseService(
	cache service.Cache,
	cryptoUtil service.HashService,
	prompt service.PromptService,
	envVar string,
) service.PassphraseService {
	if envVar == "" {
		envVar = "LOCKIFY_PASSPHRASE"
	}

	return &PassphraseService{cache, cryptoUtil, prompt, envVar}
}

// Get retrieves a passphrase from environment variable, keyring cache, or user input
func (s *PassphraseService) Get(ctx context.Context, env string) (string, error) {
	if env == "" {
		return "", fmt.Errorf("environment cannot be empty")
	}

	if passphrase := os.Getenv(s.envVar); passphrase != "" {
		return passphrase, nil
	}

	key := s.getKeyringKey(env)
	passphrase, err := s.cache.Get(key)
	if err == nil && passphrase != "" {
		return passphrase, nil
	}

	return s.getFromUser(ctx, env)
}

// Clear clears a cached passphrase for an environment
func (s *PassphraseService) Clear(ctx context.Context, env string) error {
	if env == "" {
		return fmt.Errorf("environment cannot be empty")
	}
	key := s.getKeyringKey(env)

	return s.cache.Delete(key)
}

// ClearAll clears all cached passphrases
func (s *PassphraseService) ClearAll(ctx context.Context) error {
	return s.cache.DeleteAll()
}

// Validate validates a passphrase against a vault's fingerprint
func (s *PassphraseService) Validate(
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

	return s.cryptoUtil.Verify(vault.Meta.FingerPrint, passphrase)
}

// GetWithConfirmation prompts the user for a passphrase with confirmation (for new vaults)
func (s *PassphraseService) GetWithConfirmation(
	ctx context.Context,
	env string,
	shouldCache bool,
) (string, error) {
	if env == "" {
		return "", fmt.Errorf("environment cannot be empty")
	}

	passphrase, err := s.prompt.GetPassphraseInput(
		fmt.Sprintf("Enter passphrase for environment %q:", env),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase: %w", err)
	}
	if passphrase == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	confirmation, err := s.prompt.GetPassphraseInput("Confirm passphrase:")
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase confirmation: %w", err)
	}

	if passphrase != confirmation {
		return "", fmt.Errorf("passphrases do not match")
	}

	if !shouldCache {
		shouldCacheInteractive, err := s.prompt.GetConfirmation(
			"Cache passphrase in system keyring?",
			false,
		)
		shouldCache = err == nil && shouldCacheInteractive
	}

	if shouldCache {
		if err := s.Cache(ctx, env, passphrase); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to cache passphrase: %v\n", err)
		}
	}

	return passphrase, nil
}

// Cache caches a passphrase for an environment
func (s *PassphraseService) Cache(ctx context.Context, env, passphrase string) error {
	if env == "" {
		return fmt.Errorf("environment cannot be empty")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	key := s.getKeyringKey(env)
	return s.cache.Set(key, passphrase)
}

// getFromUser prompts the user for a passphrase (for existing vaults)
func (s *PassphraseService) getFromUser(ctx context.Context, env string) (string, error) {
	passphrase, err := s.prompt.GetPassphraseInput(
		fmt.Sprintf("Enter passphrase for environment %q:", env),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase: %w", err)
	}

	if passphrase == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	shouldCache, err := s.prompt.GetConfirmation("Cache passphrase in system keyring?", false)
	if err == nil && shouldCache {
		if err := s.Cache(ctx, env, passphrase); err != nil {
			// Best effort, don't fail if caching fails
			fmt.Fprintf(os.Stderr, "Warning: failed to cache passphrase: %v\n", err)
		}
	}

	return passphrase, nil
}

// getKeyringKey returns the keyring key for an environment
func (s *PassphraseService) getKeyringKey(env string) string {
	return fmt.Sprintf("env:%s", env)
}
