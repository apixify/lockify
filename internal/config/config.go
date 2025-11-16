package config

// EncryptionConfig holds cryptographic configuration
type EncryptionConfig struct {
	ArgonTime    uint32
	ArgonMemory  uint32
	ArgonThreads uint8
	KeyLength    uint32
	NonceSize    int
}

// DefaultCryptoConfig returns cryptographic settings
func DefaultEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		ArgonTime:    3,
		ArgonMemory:  64 * 1024,
		ArgonThreads: 4,
		KeyLength:    32,
		NonceSize:    12,
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
