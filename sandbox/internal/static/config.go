package static

import (
	"os"
	"strconv"
	"strings"

	"aiexec-sandbox/internal/types"
	"aiexec-sandbox/internal/utils/log"
	"gopkg.in/yaml.v3"
)

var aiexecSandboxGlobalConfigurations types.AiexecSandboxGlobalConfigurations

func InitConfig(path string) error {
	aiexecSandboxGlobalConfigurations = types.AiexecSandboxGlobalConfigurations{}

	// read config file
	configFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer configFile.Close()

	// parse config file
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&aiexecSandboxGlobalConfigurations)
	if err != nil {
		return err
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err == nil {
		aiexecSandboxGlobalConfigurations.App.Debug = debug
	}

	max_workers := os.Getenv("MAX_WORKERS")
	if max_workers != "" {
		aiexecSandboxGlobalConfigurations.MaxWorkers, _ = strconv.Atoi(max_workers)
	}

	max_requests := os.Getenv("MAX_REQUESTS")
	if max_requests != "" {
		aiexecSandboxGlobalConfigurations.MaxRequests, _ = strconv.Atoi(max_requests)
	}

	port := os.Getenv("SANDBOX_PORT")
	if port != "" {
		aiexecSandboxGlobalConfigurations.App.Port, _ = strconv.Atoi(port)
	}

	timeout := os.Getenv("WORKER_TIMEOUT")
	if timeout != "" {
		aiexecSandboxGlobalConfigurations.WorkerTimeout, _ = strconv.Atoi(timeout)
	}

	api_key := os.Getenv("API_KEY")
	if api_key != "" {
		aiexecSandboxGlobalConfigurations.App.Key = api_key
	}

	python_path := os.Getenv("PYTHON_PATH")
	if python_path != "" {
		aiexecSandboxGlobalConfigurations.PythonPath = python_path
	}

	if aiexecSandboxGlobalConfigurations.PythonPath == "" {
		aiexecSandboxGlobalConfigurations.PythonPath = "/usr/local/bin/python3"
	}

	python_lib_path := os.Getenv("PYTHON_LIB_PATH")
	if python_lib_path != "" {
		aiexecSandboxGlobalConfigurations.PythonLibPaths = strings.Split(python_lib_path, ",")
	}

	if len(aiexecSandboxGlobalConfigurations.PythonLibPaths) == 0 {
		aiexecSandboxGlobalConfigurations.PythonLibPaths = DEFAULT_PYTHON_LIB_REQUIREMENTS
	}

	python_pip_mirror_url := os.Getenv("PIP_MIRROR_URL")
	if python_pip_mirror_url != "" {
		aiexecSandboxGlobalConfigurations.PythonPipMirrorURL = python_pip_mirror_url
	}

	python_deps_update_interval := os.Getenv("PYTHON_DEPS_UPDATE_INTERVAL")
	if python_deps_update_interval != "" {
		aiexecSandboxGlobalConfigurations.PythonDepsUpdateInterval = python_deps_update_interval
	}

	// if not set "PythonDepsUpdateInterval", update python dependencies every 30 minutes to keep the sandbox up-to-date
	if aiexecSandboxGlobalConfigurations.PythonDepsUpdateInterval == "" {
		aiexecSandboxGlobalConfigurations.PythonDepsUpdateInterval = "30m"
	}

	nodejs_path := os.Getenv("NODEJS_PATH")
	if nodejs_path != "" {
		aiexecSandboxGlobalConfigurations.NodejsPath = nodejs_path
	}

	if aiexecSandboxGlobalConfigurations.NodejsPath == "" {
		aiexecSandboxGlobalConfigurations.NodejsPath = "/usr/local/bin/node"
	}

	enable_network := os.Getenv("ENABLE_NETWORK")
	if enable_network != "" {
		aiexecSandboxGlobalConfigurations.EnableNetwork, _ = strconv.ParseBool(enable_network)
	}

	enable_preload := os.Getenv("ENABLE_PRELOAD")
	if enable_preload != "" {
		aiexecSandboxGlobalConfigurations.EnablePreload, _ = strconv.ParseBool(enable_preload)
	}

	allowed_syscalls := os.Getenv("ALLOWED_SYSCALLS")
	if allowed_syscalls != "" {
		strs := strings.Split(allowed_syscalls, ",")
		ary := make([]int, len(strs))
		for i := range ary {
			ary[i], err = strconv.Atoi(strs[i])
			if err != nil {
				return err
			}
		}
		aiexecSandboxGlobalConfigurations.AllowedSyscalls = ary
	}

	if aiexecSandboxGlobalConfigurations.EnableNetwork {
		log.Info("network has been enabled")
		socks5_proxy := os.Getenv("SOCKS5_PROXY")
		if socks5_proxy != "" {
			aiexecSandboxGlobalConfigurations.Proxy.Socks5 = socks5_proxy
		}

		if aiexecSandboxGlobalConfigurations.Proxy.Socks5 != "" {
			log.Info("using socks5 proxy: %s", aiexecSandboxGlobalConfigurations.Proxy.Socks5)
		}

		https_proxy := os.Getenv("HTTPS_PROXY")
		if https_proxy != "" {
			aiexecSandboxGlobalConfigurations.Proxy.Https = https_proxy
		}

		if aiexecSandboxGlobalConfigurations.Proxy.Https != "" {
			log.Info("using https proxy: %s", aiexecSandboxGlobalConfigurations.Proxy.Https)
		}

		http_proxy := os.Getenv("HTTP_PROXY")
		if http_proxy != "" {
			aiexecSandboxGlobalConfigurations.Proxy.Http = http_proxy
		}

		if aiexecSandboxGlobalConfigurations.Proxy.Http != "" {
			log.Info("using http proxy: %s", aiexecSandboxGlobalConfigurations.Proxy.Http)
		}
	}
	return nil
}

// avoid global modification, use value copy instead
func GetAiexecSandboxGlobalConfigurations() types.AiexecSandboxGlobalConfigurations {
	return aiexecSandboxGlobalConfigurations
}

type RunnerDependencies struct {
	PythonRequirements string
}

var runnerDependencies RunnerDependencies

func GetRunnerDependencies() RunnerDependencies {
	return runnerDependencies
}

func SetupRunnerDependencies() error {
	file, err := os.ReadFile("dependencies/python-requirements.txt")
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}

	runnerDependencies.PythonRequirements = string(file)

	return nil
}
