package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/apixify/lockify/internal/domain/service"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHashService implements service.HashService
// providing hashing, verification, and salt generation capabilities
type BcryptHashService struct{}

// NewBcryptHashService creates a new hash service
func NewBcryptHashService() service.HashService {
	return &BcryptHashService{}
}

// Hash creates a hash of the passphrase (for fingerprinting)
func (service *BcryptHashService) Hash(passphrase string) (string, error) {
	if passphrase == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash: %w", err)
	}
	return string(hash), nil
}

// Verify verifies if a passphrase matches the hash
func (service *BcryptHashService) Verify(hashedPassphrase, passphrase string) error {
	if hashedPassphrase == "" {
		return fmt.Errorf("hashed passphrase cannot be empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassphrase), []byte(passphrase)); err != nil {
		return fmt.Errorf("invalid passphrase: %w", err)
	}

	return nil
}

// GenerateSalt generates a random salt for encryption key derivation
func (service *BcryptHashService) GenerateSalt(n int) (string, error) {
	if n <= 0 {
		n = 16
	}
	s := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, s); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}
