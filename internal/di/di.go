package di

import (
	"github.com/apixify/lockify/internal/app"
	"github.com/apixify/lockify/internal/config"
	"github.com/apixify/lockify/internal/domain"
	"github.com/apixify/lockify/internal/domain/repository"
	"github.com/apixify/lockify/internal/domain/service"
	"github.com/apixify/lockify/internal/domain/storage"
	"github.com/apixify/lockify/internal/infrastructure/cache"
	"github.com/apixify/lockify/internal/infrastructure/crypto"
	"github.com/apixify/lockify/internal/infrastructure/fs"
	"github.com/apixify/lockify/internal/infrastructure/logger"
)

var (
	vaultConfig      = config.DefaultVaultConfig()
	encryptionConfig = config.DefaultEncryptionConfig()
	log              = logger.New()
)

func getHashService() service.HashService {
	return crypto.NewBcryptHashService()
}

func getCacheService() service.Cache {
	return cache.NewOSKeyring("lockify")
}

func getPassphraseService() service.PassphraseService {
	return crypto.NewPassphraseService(getCacheService(), getHashService(), vaultConfig.PassphraseEnv)
}

func getEncryptionService() service.EncryptionService {
	return crypto.NewAESEncryptionService(encryptionConfig)
}

func getFileSystemStorage() storage.FileSystem {
	return fs.NewOSFileSystem()
}

func getVaultRepository() repository.VaultRepository {
	return fs.NewFileVaultRepository(getFileSystemStorage(), vaultConfig)
}

func getVaultService() service.VaultService {
	return service.NewVaultService(getVaultRepository(), getPassphraseService(), getHashService())
}

func GetLogger() domain.Logger {
	return log
}

func BuildAddEntry() app.AddEntryUseCase {
	return app.NewAddEntryUseCase(getVaultService(), getEncryptionService())
}

func BuildClearCachedPassphrase() app.ClearCachedPassphraseUseCase {
	return app.NewClearCachedPassphraseUseCase(getPassphraseService())
}

func BuildDeleteEntry() app.DeleteEntryUseCase {
	return app.NewDeleteEntryUseCase(getVaultService())
}

func BuildExportEnv() app.ExportEnvUseCase {
	return app.NewExportEnvUseCase(getVaultService(), getEncryptionService(), GetLogger())
}

func BuildGetEntry() app.GetEntryUseCase {
	return app.NewGetEntryUseCase(getVaultService(), getEncryptionService())
}

func BuildInitializeVault() app.InitializeVaultUseCase {
	return app.NewInitializeVaultUseCase(getVaultService())
}

func BuildListEntries() app.ListEntriesUseCase {
	return app.NewListEntriesUseCase(getVaultService())
}
