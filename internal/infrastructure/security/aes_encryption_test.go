package security

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/config"
)

// createTestEncryptionService creates a test encryption service with default config
func createTestEncryptionService(t *testing.T) *AESEncryptionService {
	t.Helper()
	return NewAESEncryptionService(config.DefaultEncryptionConfig()).(*AESEncryptionService)
}

// createTestSalt creates a base64-encoded test salt
func createTestSalt(t *testing.T) string {
	t.Helper()
	return base64.StdEncoding.EncodeToString([]byte("test salt"))
}

func TestEncrypt_Success(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"
	plaintext := []byte("test plaintext")

	ciphertext, err := encryptionService.Encrypt(plaintext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() returned unexpected error: %v", err)
	}

	if ciphertext == "" {
		t.Error("Encrypt() returned empty ciphertext")
	}

	// Verify ciphertext is valid base64
	_, err = base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		t.Errorf("Encrypt() returned invalid base64 ciphertext: %v", err)
	}
}

func TestEncrypt_SamePlaintextProducesDifferentCiphertexts(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"
	plaintext := []byte("test plaintext")

	ciphertext1, err := encryptionService.Encrypt(plaintext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() first call returned unexpected error: %v", err)
	}
	ciphertext2, err := encryptionService.Encrypt(plaintext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() second call returned unexpected error: %v", err)
	}

	if ciphertext1 == ciphertext2 {
		t.Error("Encrypt() produced identical ciphertexts for same plaintext (nonce should be random)")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"
	plaintext := []byte("test plaintext")

	ciphertext, err := encryptionService.Encrypt(plaintext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() returned unexpected error: %v", err)
	}

	decrypted, err := encryptionService.Decrypt(ciphertext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Decrypt() returned unexpected error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypt() returned %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"
	plaintext := []byte("")

	ciphertext, err := encryptionService.Encrypt(plaintext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() with empty plaintext returned unexpected error: %v", err)
	}

	decrypted, err := encryptionService.Decrypt(ciphertext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Decrypt() returned unexpected error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypt() returned %q, want empty", decrypted)
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"
	plaintext := []byte("test plaintext")

	ciphertext, err := encryptionService.Encrypt(plaintext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() returned unexpected error: %v", err)
	}

	wrongPassphrase := "wrong passphrase"
	_, err = encryptionService.Decrypt(ciphertext, encodedSalt, wrongPassphrase)
	if err == nil {
		t.Error("Decrypt() with wrong passphrase expected error, got nil")
	}
	if !strings.Contains(err.Error(), "decryption failed") {
		t.Errorf("Decrypt() with wrong passphrase returned unexpected error: %v", err)
	}
}

func TestDecrypt_WrongSalt(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"
	plaintext := []byte("test plaintext")

	ciphertext, err := encryptionService.Encrypt(plaintext, encodedSalt, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() returned unexpected error: %v", err)
	}

	wrongSalt := base64.StdEncoding.EncodeToString([]byte("wrong salt"))
	_, err = encryptionService.Decrypt(ciphertext, wrongSalt, passphrase)
	if err == nil {
		t.Error("Decrypt() with wrong salt expected error, got nil")
	}
	if !strings.Contains(err.Error(), "decryption failed") {
		t.Errorf("Decrypt() with wrong salt returned unexpected error: %v", err)
	}
}

func TestDecrypt_EmptyCiphertext(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"

	_, err := encryptionService.Decrypt("", encodedSalt, passphrase)
	if err == nil {
		t.Error("Decrypt() with empty ciphertext expected error, got nil")
	}
	if !strings.Contains(err.Error(), "ciphertext cannot be empty") {
		t.Errorf("Decrypt() with empty ciphertext returned unexpected error: %v", err)
	}
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"

	_, err := encryptionService.Decrypt("invalid", encodedSalt, passphrase)
	if err == nil {
		t.Error("Decrypt() with invalid ciphertext expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid ciphertext encoding") {
		t.Errorf("Decrypt() with invalid ciphertext returned unexpected error: %v", err)
	}
}

func TestEncrypt_NilPlaintext(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"

	_, err := encryptionService.Encrypt(nil, encodedSalt, passphrase)
	if err == nil {
		t.Error("Encrypt() with nil plaintext expected error, got nil")
	}
	if !strings.Contains(err.Error(), "plaintext cannot be nil") {
		t.Errorf("Encrypt() with nil plaintext returned unexpected error: %v", err)
	}
}

func TestEncrypt_EmptySalt(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	passphrase := "test passphrase"
	plaintext := []byte("test plaintext")

	_, err := encryptionService.Encrypt(plaintext, "", passphrase)
	if err == nil {
		t.Error("Encrypt() with empty salt expected error, got nil")
	}
	if !strings.Contains(err.Error(), "salt cannot be empty") {
		t.Errorf("Encrypt() with empty salt returned unexpected error: %v", err)
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	plaintext := []byte("test plaintext")

	_, err := encryptionService.Encrypt(plaintext, encodedSalt, "")
	if err == nil {
		t.Error("Encrypt() with empty passphrase expected error, got nil")
	}
	if !strings.Contains(err.Error(), "passphrase cannot be empty") {
		t.Errorf("Encrypt() with empty passphrase returned unexpected error: %v", err)
	}
}

func TestDecrypt_CiphertextTooShort(t *testing.T) {
	encryptionService := createTestEncryptionService(t)
	encodedSalt := createTestSalt(t)
	passphrase := "test passphrase"

	shortCiphertext := base64.StdEncoding.EncodeToString([]byte("short"))

	_, err := encryptionService.Decrypt(shortCiphertext, encodedSalt, passphrase)
	if err == nil {
		t.Error("Decrypt() with too short ciphertext expected error, got nil")
	}
	if !strings.Contains(err.Error(), "ciphertext too short") {
		t.Errorf("Decrypt() with too short ciphertext returned unexpected error: %v", err)
	}
}
