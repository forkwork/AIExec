package integrationtests_test

import (
	"aiexec-sandbox/internal/core/runner/python"
	"aiexec-sandbox/internal/static"
	"aiexec-sandbox/internal/utils/log"
)

func init() {
	static.InitConfig("conf/config.yaml")

	err := python.PreparePythonDependenciesEnv()
	if err != nil {
		log.Panic("failed to initialize python dependencies sandbox: %v", err)
	}
}
