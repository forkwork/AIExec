package persistence

import (
	"aiexec-plugin/internal/oss"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/internal/utils/log"
)

var (
	persistence *Persistence
)

func InitPersistence(oss oss.OSS, config *app.Config) {
	persistence = &Persistence{
		storage:        NewWrapper(oss, config.PersistenceStoragePath),
		maxStorageSize: config.PersistenceStorageMaxSize,
	}

	log.Info("Persistence initialized")
}

func GetPersistence() *Persistence {
	return persistence
}
