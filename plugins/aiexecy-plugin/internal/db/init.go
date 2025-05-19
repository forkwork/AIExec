package db

import (
	"aiexec-plugin/internal/db/mysql"
	"aiexec-plugin/internal/db/pg"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/internal/types/models"
	"aiexec-plugin/internal/utils/log"
)

func autoMigrate() error {
	err := AiexecPluginDB.AutoMigrate(
		models.Plugin{},
		models.PluginInstallation{},
		models.PluginDeclaration{},
		models.Endpoint{},
		models.ServerlessRuntime{},
		models.ToolInstallation{},
		models.AIModelInstallation{},
		models.InstallTask{},
		models.TenantStorage{},
		models.AgentStrategyInstallation{},
	)

	if err != nil {
		return err
	}

	// check if "declaration" column exists in Plugin/ServerlessRuntime/ToolInstallation/AIModelInstallation/AgentStrategyInstallation
	// drop the "declaration" column not null constraint if exists
	ignoreDeclarationColumn := func(table string) error {
		if AiexecPluginDB.Migrator().HasColumn(table, "declaration") {
			// remove NOT NULL constraint on declaration column
			if err := AiexecPluginDB.Exec("ALTER TABLE " + table + " ALTER COLUMN declaration DROP NOT NULL").Error; err != nil {
				return err
			}
		}
		return nil
	}

	tables := []string{
		"plugins",
		"serverless_runtimes",
		"tool_installations",
		"ai_model_installations",
		"agent_strategy_installations",
	}

	for _, table := range tables {
		if err := ignoreDeclarationColumn(table); err != nil {
			return err
		}
	}

	return nil
}

func Init(config *app.Config) {
	var err error
	if config.DBType == "postgresql" {
		AiexecPluginDB, err = pg.InitPluginDB(
			config.DBHost,
			int(config.DBPort),
			config.DBDatabase,
			config.DBDefaultDatabase,
			config.DBUsername,
			config.DBPassword,
			config.DBSslMode,
			config.DBMaxIdleConns,
			config.DBMaxOpenConns,
			config.DBConnMaxLifetime,
		)
	} else if config.DBType == "mysql" {
		AiexecPluginDB, err = mysql.InitPluginDB(
			config.DBHost,
			int(config.DBPort),
			config.DBDatabase,
			config.DBDefaultDatabase,
			config.DBUsername,
			config.DBPassword,
			config.DBSslMode,
			config.DBMaxIdleConns,
			config.DBMaxOpenConns,
			config.DBConnMaxLifetime,
		)
	} else {
		log.Panic("unsupported database type: %v", config.DBType)
	}

	if err != nil {
		log.Panic("failed to init aiexec plugin db: %v", err)
	}

	err = autoMigrate()
	if err != nil {
		log.Panic("failed to auto migrate: %v", err)
	}

	log.Info("aiexec plugin db initialized")
}

func Close() {
	db, err := AiexecPluginDB.DB()
	if err != nil {
		log.Error("failed to close aiexec plugin db: %v", err)
		return
	}

	err = db.Close()
	if err != nil {
		log.Error("failed to close aiexec plugin db: %v", err)
		return
	}

	log.Info("aiexec plugin db closed")
}
