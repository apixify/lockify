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
	return security.NewPassphraseService(
		getCacheService(),
		getHashService(),
		BuildPromptService(),
		vaultConfig.PassphraseEnv,
	)
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
	return fs.NewImportService()
}

// GetLogger returns the logger instance.
func GetLogger() domain.Logger {
	return log
}

// BuildAddEntry creates and returns an AddEntry use case.
func BuildAddEntry() app.AddEntryUc {
	return app.NewAddEntryUseCase(getVaultService(), getEncryptionService())
}

// BuildPromptService creates and returns a prompt service instance.
func BuildPromptService() service.PromptService {
	return prompt.NewService()
}

// BuildClearCachedPassphrase creates and returns a ClearCachedPassphrase use case.
func BuildClearCachedPassphrase() app.ClearCachedPassphraseUc {
	return app.NewClearCachedPassphraseUseCase(getPassphraseService())
}

// BuildClearEnvCachedPassphrase creates and returns a ClearEnvCachedPassphrase use case.
func BuildClearEnvCachedPassphrase() app.ClearEnvCachedPassphraseUseCase {
	return app.NewClearEnvCachedPassphraseUseCase(getPassphraseService())
}

// BuildDeleteEntry creates and returns a DeleteEntry use case.
func BuildDeleteEntry() app.DeleteEntryUc {
	return app.NewDeleteEntryUseCase(getVaultService())
}

// BuildExportEnv creates and returns an ExportEnv use case.
func BuildExportEnv() app.ExportEnvUc {
	return app.NewExportEnvUseCase(getVaultService(), getEncryptionService(), GetLogger())
}

// BuildGetEntry creates and returns a GetEntry use case.
func BuildGetEntry() app.GetEntryUc {
	return app.NewGetEntryUseCase(getVaultService(), getEncryptionService())
}

// BuildInitializeVault creates and returns an InitializeVault use case.
func BuildInitializeVault() app.InitUc {
	return app.NewInitializeVaultUseCase(getVaultService())
}

// BuildListEntries creates and returns a ListEntries use case.
func BuildListEntries() app.ListEntriesUc {
	return app.NewListEntriesUseCase(getVaultService())
}

// BuildRotatePassphrase creates and returns a RotatePassphrase use case.
func BuildRotatePassphrase() app.RotatePassphraseUc {
	return app.NewRotatePassphraseUseCase(
		getVaultRepository(),
		getEncryptionService(),
		getHashService(),
	)
}

// BuildImportEnv creates and returns an ImportEnv use case.
func BuildImportEnv() app.ImportEnvUc {
	return app.NewImportEnvUseCase(
		getVaultService(),
		getImportService(),
		getEncryptionService(),
		GetLogger(),
	)
}

// BuildCachePassphrase creates and returns a CachePassphrase use case.
func BuildCachePassphrase() app.CachePassphraseUc {
	return app.NewCachePassphraseUseCase(
		getVaultRepository(),
		getPassphraseService(),
		getHashService(),
	)
}
