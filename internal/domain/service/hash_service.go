package service

// CryptoService provides cryptographic utility operations:
// - Passphrase hashing and verification (for fingerprints)
// - Salt generation
type HashService interface {
	// Hash creates a hash of the passphrase (for fingerprinting)
	Hash(passphrase string) (string, error)
	// Verify verifies if a passphrase matches the hash
	Verify(hashedPassphrase, passphrase string) error
	// GenerateSalt generates a random salt for encryption key derivation
	GenerateSalt(size int) (string, error)
}
