package security

import (
	"fmt"
	"os"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// PassphraseService implements service.PassphraseService
type PassphraseService struct {
	cache               service.Cache
	cryptoUtil          service.HashService
	prompt              service.PromptService
	environmentProvider service.EnvironmentProvider
	envVar              string
}

// NewPassphraseService creates a new passphrase service
func NewPassphraseService(
	cache service.Cache,
	cryptoUtil service.HashService,
	prompt service.PromptService,
	environmentProvider service.EnvironmentProvider,
	envVar string,
) service.PassphraseService {
	if envVar == "" {
		envVar = "LOCKIFY_PASSPHRASE"
	}

	return &PassphraseService{
		cache:               cache,
		cryptoUtil:          cryptoUtil,
		prompt:              prompt,
		environmentProvider: environmentProvider,
		envVar:              envVar,
	}
}

// Get retrieves a passphrase from environment variable, keyring cache, or user input
func (s *PassphraseService) Get(vctx *model.VaultContext) (string, error) {
	if vctx.Env == "" {
		return "", fmt.Errorf("environment cannot be empty")
	}

	if passphrase := s.environmentProvider.GetPassphrase(s.envVar); passphrase != "" {
		return passphrase, nil
	}

	key := s.getKeyringKey(vctx.Env)
	passphrase, err := s.cache.Get(key)
	if err == nil && passphrase != "" {
		return passphrase, nil
	}

	return s.getFromUser(vctx)
}

// Clear clears a cached passphrase for an environment
func (s *PassphraseService) Clear(vctx *model.VaultContext) error {
	if vctx.Env == "" {
		return fmt.Errorf("environment cannot be empty")
	}
	key := s.getKeyringKey(vctx.Env)

	return s.cache.Delete(key)
}

// ClearAll clears all cached passphrases
func (s *PassphraseService) ClearAll(vctx *model.VaultContext) error {
	return s.cache.DeleteAll()
}

// Validate validates a passphrase against a vault's fingerprint
func (s *PassphraseService) Validate(
	vctx *model.VaultContext,
	vault *model.Vault,
	passphrase string,
) error {
	if vault == nil {
		return fmt.Errorf("vault cannot be nil")
	}
	if vault.FingerPrint() == "" {
		return fmt.Errorf("fingerprint cannot be empty")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	return s.cryptoUtil.Verify(vault.FingerPrint(), passphrase)
}

// GetWithConfirmation prompts the user for a passphrase with confirmation (for new vaults)
func (s *PassphraseService) GetWithConfirmation(vctx *model.VaultContext) (string, error) {
	if vctx.Env == "" {
		return "", fmt.Errorf("environment cannot be empty")
	}

	passphrase, err := s.prompt.GetPassphraseInput(
		fmt.Sprintf("Enter passphrase for environment %q:", vctx.Env),
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

	if !vctx.ShouldCache {
		shouldCacheInteractive, err := s.prompt.GetConfirmation(
			"Cache passphrase in system keyring?",
			false,
		)
		vctx.ShouldCache = err == nil && shouldCacheInteractive
	}

	if vctx.ShouldCache {
		if err := s.Cache(vctx, passphrase); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to cache passphrase: %v\n", err)
		}
	}

	return passphrase, nil
}

// Cache caches a passphrase for an environment
func (s *PassphraseService) Cache(vctx *model.VaultContext, passphrase string) error {
	if vctx.Env == "" {
		return fmt.Errorf("environment cannot be empty")
	}
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	key := s.getKeyringKey(vctx.Env)
	return s.cache.Set(key, passphrase)
}

// getFromUser prompts the user for a passphrase (for existing vaults)
func (s *PassphraseService) getFromUser(vctx *model.VaultContext) (string, error) {
	passphrase, err := s.prompt.GetPassphraseInput(
		fmt.Sprintf("Enter passphrase for environment %q:", vctx.Env),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase: %w", err)
	}

	if passphrase == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	shouldCache, err := s.prompt.GetConfirmation("Cache passphrase in system keyring?", false)
	if err == nil && shouldCache {
		if err := s.Cache(vctx, passphrase); err != nil {
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
