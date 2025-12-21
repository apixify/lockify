package config

const (
	// DefaultArgonTime is the default time parameter for Argon2 key derivation.
	DefaultArgonTime uint32 = 3
	// DefaultArgonMemoryKB is the default memory parameter in KB for Argon2 key derivation.
	DefaultArgonMemoryKB uint32 = 64
	// DefaultArgonThreads is the default number of threads for Argon2 key derivation.
	DefaultArgonThreads uint8 = 4
	// DefaultKeyLength is the default key length in bytes for encryption.
	DefaultKeyLength uint32 = 32
	// DefaultNonceSize is the default nonce size in bytes for AES-GCM encryption.
	DefaultNonceSize int = 12
	// DefaultSaltSize is the default salt size in bytes for key derivation.
	DefaultSaltSize int = 16
	// bytesPerKB is the number of bytes in a kilobyte.
	bytesPerKB uint32 = 1024
)

// EncryptionConfig holds cryptographic configuration
type EncryptionConfig struct {
	ArgonTime    uint32
	ArgonMemory  uint32
	ArgonThreads uint8
	KeyLength    uint32
	NonceSize    int
}

// DefaultEncryptionConfig returns default cryptographic settings.
func DefaultEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		ArgonTime:    DefaultArgonTime,
		ArgonMemory:  DefaultArgonMemoryKB * bytesPerKB,
		ArgonThreads: DefaultArgonThreads,
		KeyLength:    DefaultKeyLength,
		NonceSize:    DefaultNonceSize,
	}
}

// VaultConfig holds vault-related configuration
type VaultConfig struct {
	BaseDir       string
	FileMode      uint32
	DirMode       uint32
	DefaultEnv    string
	PassphraseEnv string
}

// DefaultVaultConfig returns default vault configuration
func DefaultVaultConfig() VaultConfig {
	return VaultConfig{
		BaseDir:       ".lockify",
		FileMode:      0600, // rw-------
		DirMode:       0700, // rwx------
		DefaultEnv:    "local",
		PassphraseEnv: "LOCKIFY_PASSPHRASE",
	}
}

// GetVaultPath returns the path to a vault file for an environment
func (c VaultConfig) GetVaultPath(env string) string {
	if c.BaseDir == "" {
		return env + ".vault.enc"
	}
	return c.BaseDir + "/" + env + ".vault.enc"
}
