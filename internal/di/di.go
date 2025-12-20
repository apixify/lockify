package di

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/config"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/repository"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/storage"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/cache"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/fs"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/logger"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/prompt"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/security"
)

var (
	vaultConfig      = config.DefaultVaultConfig()
	encryptionConfig = config.DefaultEncryptionConfig()
	log              = logger.New()
)

func getHashService() service.HashService {
	return security.NewBcryptHashService()
}

func getCacheService() service.Cache {
	return cache.NewOSKeyring("lockify")
}

func getPassphraseService() service.PassphraseService {
	return security.NewPassphraseService(getCacheService(), getHashService(), vaultConfig.PassphraseEnv)
}

func getEncryptionService() service.EncryptionService {
	return security.NewAESEncryptionService(encryptionConfig)
}

func getFileSystemStorage() storage.FileSystem {
	return fs.NewOSFileSystem()
}

func getVaultRepository() repository.VaultRepository {
	return fs.NewFileVaultRepository(getFileSystemStorage(), vaultConfig)
}

func getVaultService() service.VaultServiceInterface {
	return service.NewVaultService(getVaultRepository(), getPassphraseService(), getHashService())
}

func getImportService() service.ImportService {
	return fs.NewFsImportService()
}

func GetLogger() domain.Logger {
	return log
}

func BuildAddEntry() app.AddEntryUc {
	return app.NewAddEntryUseCase(getVaultService(), getEncryptionService())
}

func BuildPromptService() service.PromptService {
	return prompt.NewPromptService()
}

func BuildClearCachedPassphrase() app.ClearCachedPassphraseUc {
	return app.NewClearCachedPassphraseUseCase(getPassphraseService())
}

func BuildClearEnvCachedPassphrase() app.ClearEnvCachedPassphraseUseCase {
	return app.NewClearEnvCachedPassphraseUseCase(getPassphraseService())
}

func BuildDeleteEntry() app.DeleteEntryUc {
	return app.NewDeleteEntryUseCase(getVaultService())
}

func BuildExportEnv() app.ExportEnvUc {
	return app.NewExportEnvUseCase(getVaultService(), getEncryptionService(), GetLogger())
}

func BuildGetEntry() app.GetEntryUc {
	return app.NewGetEntryUseCase(getVaultService(), getEncryptionService())
}

func BuildInitializeVault() app.InitUc {
	return app.NewInitializeVaultUseCase(getVaultService())
}

func BuildListEntries() app.ListEntriesUc {
	return app.NewListEntriesUseCase(getVaultService())
}

func BuildRotatePassphrase() app.RotatePassphraseUc {
	return app.NewRotatePassphraseUseCase(getVaultRepository(), getEncryptionService(), getHashService())
}

func BuildImportEnv() app.ImportEnvUc {
	return app.NewImportEnvUseCase(getVaultService(), getImportService(), getEncryptionService(), GetLogger())
}
