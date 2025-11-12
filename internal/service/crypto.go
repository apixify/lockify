package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	keyLen       = 32
	nonceSize    = 12
)

type CryptoService struct {
	aesgcm cipher.AEAD
}

func NewCryptoService(encodedSalt string, passphrase string) (*CryptoService, error) {
	salt, _ := base64.StdEncoding.DecodeString(encodedSalt)
	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &CryptoService{aesgcm: aesgcm}, nil
}

func GenerateSalt(n int) (string, error) {
	if n <= 0 {
		n = 16
	}
	s := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, s); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}

func (crypto *CryptoService) EncryptValue(plaintext []byte) (string, error) {
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := crypto.aesgcm.Seal(nil, nonce, plaintext, nil)

	out := append(nonce, ciphertext...)
	enc := base64.StdEncoding.EncodeToString(out)

	return enc, nil
}

func (crypto *CryptoService) DecryptValue(encoded string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(raw) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := raw[:nonceSize]
	ciphertext := raw[nonceSize:]

	plaintext, err := crypto.aesgcm.Open(nil, nonce, ciphertext, nil)
	value := string(plaintext)

	if err != nil {
		return "", err
	}

	return value, nil
}

func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(passphrase), salt, argonTime, argonMemory, argonThreads, keyLen)
}
