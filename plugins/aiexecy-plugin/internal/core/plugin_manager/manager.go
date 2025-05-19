package plugin_manager

import (
	"errors"
	"fmt"
	"os"

	"aiexec-plugin/internal/core/aiexec_invocation"
	"aiexec-plugin/internal/core/aiexec_invocation/real"
	"aiexec-plugin/internal/core/plugin_manager/debugging_runtime"
	"aiexec-plugin/internal/core/plugin_manager/media_transport"
	serverless "aiexec-plugin/internal/core/plugin_manager/serverless_connector"
	"aiexec-plugin/internal/db"
	"aiexec-plugin/internal/oss"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/internal/types/models"
	"aiexec-plugin/internal/utils/cache"
	"aiexec-plugin/internal/utils/cache/helper"
	"aiexec-plugin/internal/utils/lock"
	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/internal/utils/mapping"
	"aiexec-plugin/pkg/entities/plugin_entities"
	"aiexec-plugin/pkg/plugin_packager/decoder"
)

type PluginManager struct {
	m mapping.Map[string, plugin_entities.PluginLifetime]

	// max size of a plugin package
	maxPluginPackageSize int64

	// where the plugin finally running
	workingDirectory string

	// where the plugin finally installed but not running
	pluginStoragePath string

	// mediaBucket is used to manage media files like plugin icons, images, etc.
	mediaBucket *media_transport.MediaBucket

	// packageBucket is used to manage plugin packages, all the packages uploaded by users will be saved here
	packageBucket *media_transport.PackageBucket

	// installedBucket is used to manage installed plugins, all the installed plugins will be saved here
	installedBucket *media_transport.InstalledBucket

	// register plugin
	pluginRegisters []func(lifetime plugin_entities.PluginLifetime) error

	// localPluginLaunchingLock is a lock to launch local plugins
	localPluginLaunchingLock *lock.GranularityLock

	// backwardsInvocation is a handle to invoke aiexec
	backwardsInvocation aiexec_invocation.BackwardsInvocation

	// python interpreter path
	pythonInterpreterPath string

	// uv path
	uvPath string

	// python env init timeout
	pythonEnvInitTimeout int

	// proxy settings
	HttpProxy  string
	HttpsProxy string
	NoProxy    string

	// pip mirror url
	pipMirrorUrl string

	// pip prefer binary
	pipPreferBinary bool

	// pip verbose
	pipVerbose bool

	// pip extra args
	pipExtraArgs string

	// python compileall extra args
	pythonCompileAllExtraArgs string

	// remote plugin server
	remotePluginServer debugging_runtime.RemotePluginServerInterface

	// max launching lock to prevent too many plugins launching at the same time
	maxLaunchingLock chan bool

	// platform, local or serverless
	platform app.PlatformType

	// serverless connector launch timeout
	serverlessConnectorLaunchTimeout int

	// plugin stdio buffer size
	pluginStdioBufferSize    int
	pluginStdioMaxBufferSize int
}

var (
	manager *PluginManager
)

func InitGlobalManager(oss oss.OSS, configuration *app.Config) *PluginManager {
	manager = &PluginManager{
		maxPluginPackageSize: configuration.MaxPluginPackageSize,
		pluginStoragePath:    configuration.PluginInstalledPath,
		workingDirectory:     configuration.PluginWorkingPath,
		mediaBucket: media_transport.NewAssetsBucket(
			oss,
			configuration.PluginMediaCachePath,
			configuration.PluginMediaCacheSize,
		),
		packageBucket: media_transport.NewPackageBucket(
			oss,
			configuration.PluginPackageCachePath,
		),
		installedBucket: media_transport.NewInstalledBucket(
			oss,
			configuration.PluginInstalledPath,
		),
		localPluginLaunchingLock:         lock.NewGranularityLock(),
		maxLaunchingLock:                 make(chan bool, 2), // by default, we allow 2 plugins launching at the same time
		pythonInterpreterPath:            configuration.PythonInterpreterPath,
		uvPath:                           configuration.UvPath,
		pythonEnvInitTimeout:             configuration.PythonEnvInitTimeout,
		pythonCompileAllExtraArgs:        configuration.PythonCompileAllExtraArgs,
		platform:                         configuration.Platform,
		HttpProxy:                        configuration.HttpProxy,
		HttpsProxy:                       configuration.HttpsProxy,
		NoProxy:                          configuration.NoProxy,
		pipMirrorUrl:                     configuration.PipMirrorUrl,
		pipPreferBinary:                  *configuration.PipPreferBinary,
		pipVerbose:                       *configuration.PipVerbose,
		pipExtraArgs:                     configuration.PipExtraArgs,
		serverlessConnectorLaunchTimeout: configuration.AiexecPluginServerlessConnectorLaunchTimeout,
		pluginStdioBufferSize:            configuration.PluginStdioBufferSize,
		pluginStdioMaxBufferSize:         configuration.PluginStdioMaxBufferSize,
	}

	return manager
}

func Manager() *PluginManager {
	return manager
}

func (p *PluginManager) Get(
	identity plugin_entities.PluginUniqueIdentifier,
) (plugin_entities.PluginLifetime, error) {
	if identity.RemoteLike() || p.platform == app.PLATFORM_LOCAL {
		// check if it's a debugging plugin or a local plugin
		if v, ok := p.m.Load(identity.String()); ok {
			return v, nil
		}
		return nil, errors.New("plugin not found")
	} else {
		// otherwise, use serverless runtime instead
		pluginSessionInterface, err := p.getServerlessPluginRuntime(identity)
		if err != nil {
			return nil, err
		}

		return pluginSessionInterface, nil
	}
}

func (p *PluginManager) GetAsset(id string) ([]byte, error) {
	return p.mediaBucket.Get(id)
}

func (p *PluginManager) Launch(configuration *app.Config) {
	log.Info("start plugin manager daemon...")

	// init redis client
	if err := cache.InitRedisClient(
		fmt.Sprintf("%s:%d", configuration.RedisHost, configuration.RedisPort),
		configuration.RedisUser,
		configuration.RedisPass,
		configuration.RedisUseSsl,
		configuration.RedisDB,
	); err != nil {
		log.Panic("init redis client failed: %s", err.Error())
	}

	invocation, err := real.NewAiexecInvocationDaemon(
		real.NewAiexecInvocationDaemonPayload{
			BaseUrl:      configuration.AiexecInnerApiURL,
			CallingKey:   configuration.AiexecInnerApiKey,
			WriteTimeout: configuration.AiexecInvocationWriteTimeout,
			ReadTimeout:  configuration.AiexecInvocationReadTimeout,
		},
	)
	if err != nil {
		log.Panic("init aiexec invocation daemon failed: %s", err.Error())
	}
	p.backwardsInvocation = invocation

	// start local watcher
	if configuration.Platform == app.PLATFORM_LOCAL {
		p.startLocalWatcher()
	}

	// launch serverless connector
	if configuration.Platform == app.PLATFORM_SERVERLESS {
		serverless.Init(configuration)
	}

	// start remote watcher
	p.startRemoteWatcher(configuration)
}

func (p *PluginManager) BackwardsInvocation() aiexec_invocation.BackwardsInvocation {
	return p.backwardsInvocation
}

func (p *PluginManager) SavePackage(plugin_unique_identifier plugin_entities.PluginUniqueIdentifier, pkg []byte, thirdPartySignatureVerificationConfig *decoder.ThirdPartySignatureVerificationConfig) (
	*plugin_entities.PluginDeclaration, error,
) {
	// try to decode the package
	packageDecoder, err := decoder.NewZipPluginDecoderWithThirdPartySignatureVerificationConfig(pkg, thirdPartySignatureVerificationConfig)
	if err != nil {
		return nil, err
	}

	// get the declaration
	declaration, err := packageDecoder.Manifest()
	if err != nil {
		return nil, err
	}

	// get the assets
	assets, err := packageDecoder.Assets()
	if err != nil {
		return nil, err
	}

	// remap the assets
	_, err = p.mediaBucket.RemapAssets(&declaration, assets)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to remap assets"))
	}

	uniqueIdentifier, err := packageDecoder.UniqueIdentity()
	if err != nil {
		return nil, err
	}

	// save to storage
	err = p.packageBucket.Save(plugin_unique_identifier.String(), pkg)
	if err != nil {
		return nil, err
	}

	// create plugin if not exists
	if _, err := db.GetOne[models.PluginDeclaration](
		db.Equal("plugin_unique_identifier", uniqueIdentifier.String()),
	); err == db.ErrDatabaseNotFound {
		err = db.Create(&models.PluginDeclaration{
			PluginUniqueIdentifier: uniqueIdentifier.String(),
			PluginID:               uniqueIdentifier.PluginID(),
			Declaration:            declaration,
		})
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &declaration, nil
}

func (p *PluginManager) GetPackage(
	plugin_unique_identifier plugin_entities.PluginUniqueIdentifier,
) ([]byte, error) {
	file, err := p.packageBucket.Get(plugin_unique_identifier.String())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("plugin package not found, please upload it firstly")
		}
		return nil, err
	}

	return file, nil
}

func (p *PluginManager) GetDeclaration(
	plugin_unique_identifier plugin_entities.PluginUniqueIdentifier,
	tenant_id string,
	runtime_type plugin_entities.PluginRuntimeType,
) (
	*plugin_entities.PluginDeclaration, error,
) {
	return helper.CombinedGetPluginDeclaration(
		plugin_unique_identifier, runtime_type,
	)
}
