package security

import (
	"encoding/base64"
	"testing"
)

func TestHash_Success(t *testing.T) {
	hashService := NewBcryptHashService()
	hash, err := hashService.Hash("test")
	if err != nil {
		t.Fatalf("Hash() returned unexpected error: %v", err)
	}
	if hash == "" {
		t.Error("Hash() returned empty hash")
	}
}

func TestVerify_CorrectPassphrase(t *testing.T) {
	hashService := NewBcryptHashService()
	hash, err := hashService.Hash("test")
	if err != nil {
		t.Fatalf("Hash() returned unexpected error: %v", err)
	}
	if err := hashService.Verify(hash, "test"); err != nil {
		t.Errorf("Verify() with correct passphrase returned error: %v", err)
	}
}

func TestVerify_WrongPassphrase(t *testing.T) {
	hashService := NewBcryptHashService()
	hash, err := hashService.Hash("test")
	if err != nil {
		t.Fatalf("Hash() returned unexpected error: %v", err)
	}
	if err := hashService.Verify(hash, "wrong"); err == nil {
		t.Error("Verify() with wrong passphrase expected error, got nil")
	}
}

func TestVerify_EmptyHashedPassphrase(t *testing.T) {
	hashService := NewBcryptHashService()
	if err := hashService.Verify("", "test"); err == nil {
		t.Error("Verify() with empty hashed passphrase expected error, got nil")
	}
}

func TestVerify_InvalidHashedPassphrase(t *testing.T) {
	hashService := NewBcryptHashService()
	if err := hashService.Verify("invalid", "test"); err == nil {
		t.Error("Verify() with invalid hashed passphrase expected error, got nil")
	}
}

func TestGenerateSalt_Success(t *testing.T) {
	hashService := NewBcryptHashService()
	saltSet := make(map[string]bool, 20)
	for i := 0; i < 20; i++ {
		salt, err := hashService.GenerateSalt(16)
		if err != nil {
			t.Fatalf("GenerateSalt() returned unexpected error: %v", err)
		}
		if salt == "" {
			t.Error("GenerateSalt() returned empty salt")
		}
		if saltSet[salt] {
			t.Errorf("GenerateSalt() returned duplicate salt: %q", salt)
		}
		saltSet[salt] = true
	}

	if len(saltSet) != 20 {
		t.Errorf("GenerateSalt() produced %d unique salts, want 20", len(saltSet))
	}
}

func TestGenerateSalt_ZeroSize(t *testing.T) {
	hashService := NewBcryptHashService()
	salt, err := hashService.GenerateSalt(0)
	if err != nil {
		t.Fatalf("GenerateSalt(0) returned unexpected error: %v", err)
	}
	if salt == "" {
		t.Error("GenerateSalt(0) returned empty salt")
	}
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		t.Fatalf("failed to decode salt as base64: %v", err)
	}

	if len(saltBytes) != 16 {
		t.Errorf("GenerateSalt(0) produced salt of %d bytes, want 16", len(saltBytes))
	}
}

func TestGenerateSalt_NegativeSize(t *testing.T) {
	hashService := NewBcryptHashService()
	salt, err := hashService.GenerateSalt(-1)
	if err != nil {
		t.Fatalf("GenerateSalt(-1) returned unexpected error: %v", err)
	}
	if salt == "" {
		t.Error("GenerateSalt(-1) returned empty salt")
	}
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		t.Fatalf("failed to decode salt as base64: %v", err)
	}

	if len(saltBytes) != 16 {
		t.Errorf("GenerateSalt(-1) produced salt of %d bytes, want 16", len(saltBytes))
	}
}

func TestHash_ProducesDifferentHashes(t *testing.T) {
	hashService := NewBcryptHashService()
	hash1, err := hashService.Hash("test")
	if err != nil {
		t.Fatalf("Hash() first call returned unexpected error: %v", err)
	}
	hash2, err := hashService.Hash("test")
	if err != nil {
		t.Fatalf("Hash() second call returned unexpected error: %v", err)
	}
	if hash1 == hash2 {
		t.Error("Hash() produced identical hashes for same passphrase (bcrypt should salt internally)")
	}

	// Verify both hashes work correctly
	if err := hashService.Verify(hash1, "test"); err != nil {
		t.Errorf("first hash failed to verify: %v", err)
	}
	if err := hashService.Verify(hash2, "test"); err != nil {
		t.Errorf("second hash failed to verify: %v", err)
	}
}
