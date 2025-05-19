package server

import (
	"context"

	"github.com/getsentry/sentry-go"
	"aiexec-plugin/internal/cluster"
	"aiexec-plugin/internal/core/persistence"
	"aiexec-plugin/internal/core/plugin_manager"
	"aiexec-plugin/internal/db"
	"aiexec-plugin/internal/oss"
	"aiexec-plugin/internal/oss/aliyun"
	"aiexec-plugin/internal/oss/azure"
	"aiexec-plugin/internal/oss/gcs"
	"aiexec-plugin/internal/oss/local"
	"aiexec-plugin/internal/oss/s3"
	"aiexec-plugin/internal/oss/tencent_cos"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/internal/utils/routine"
)

func initOSS(config *app.Config) oss.OSS {
	// init storage
	ctx := context.TODO()
	var storage oss.OSS
	var err error
	switch config.PluginStorageType {
	case oss.OSS_TYPE_S3:
		storage, err = s3.NewS3Storage(
			config.S3UseAwsManagedIam,
			config.S3Endpoint,
			config.S3UsePathStyle,
			config.AWSAccessKey,
			config.AWSSecretKey,
			config.PluginStorageOSSBucket,
			config.AWSRegion,
		)
	case oss.OSS_TYPE_LOCAL:
		storage = local.NewLocalStorage(config.PluginStorageLocalRoot)
	case oss.OSS_TYPE_TENCENT_COS:
		storage, err = tencent_cos.NewTencentCOSStorage(
			config.TencentCOSSecretId,
			config.TencentCOSSecretKey,
			config.TencentCOSRegion,
			config.PluginStorageOSSBucket,
		)
	case oss.OSS_TYPE_AZURE_BLOB:
		storage, err = azure.NewAzureBlobStorage(
			config.AzureBlobStorageContainerName,
			config.AzureBlobStorageConnectionString,
		)
	case oss.OSS_TYPE_GCS:
		storage, err = gcs.NewGCSStorage(ctx, config.PluginStorageOSSBucket)
	case oss.OSS_TYPE_ALIYUN_OSS:
		storage, err = aliyun.NewAliyunOSSStorage(
			config.AliyunOSSRegion,
			config.AliyunOSSEndpoint,
			config.AliyunOSSAccessKeyID,
			config.AliyunOSSAccessKeySecret,
			config.AliyunOSSAuthVersion,
			config.AliyunOSSPath,
			config.PluginStorageOSSBucket,
		)
	default:
		log.Panic("Invalid plugin storage type: %s", config.PluginStorageType)
	}

	if err != nil {
		log.Panic("Failed to create storage: %s", err)
	}

	return storage
}

func (app *App) Run(config *app.Config) {
	// init routine pool
	if config.SentryEnabled {
		routine.InitPool(config.RoutinePoolSize, sentry.ClientOptions{
			Dsn:              config.SentryDSN,
			AttachStacktrace: config.SentryAttachStacktrace,
			TracesSampleRate: config.SentryTracesSampleRate,
			SampleRate:       config.SentrySampleRate,
			EnableTracing:    config.SentryTracingEnabled,
		})
	} else {
		routine.InitPool(config.RoutinePoolSize)
	}

	// init db
	db.Init(config)

	// init oss
	oss := initOSS(config)

	// create manager
	manager := plugin_manager.InitGlobalManager(oss, config)

	// create cluster
	app.cluster = cluster.NewCluster(config, manager)

	// register plugin lifetime event
	manager.AddPluginRegisterHandler(app.cluster.RegisterPlugin)

	// init manager
	manager.Launch(config)

	// init persistence
	persistence.InitPersistence(oss, config)

	// launch cluster
	app.cluster.Launch()

	// start http server
	app.server(config)

	// block
	select {}
}
