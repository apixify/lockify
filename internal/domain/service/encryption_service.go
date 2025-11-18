package service

// EncryptionService provides encryption and decryption operations for vault entries
type EncryptionService interface {
	// Encrypt encrypts plaintext and returns base64-encoded ciphertext
	Encrypt(plaintext []byte, encodedSalt, passphrase string) (string, error)
	// Decrypt decrypts base64-encoded ciphertext and returns plaintext
	Decrypt(ciphertext, encodedSalt, passphrase string) ([]byte, error)
}
