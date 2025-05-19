package service

import (
	"errors"

	"aiexec-sandbox/internal/core/runner/types"
	"aiexec-sandbox/internal/static"
)

var (
	ErrNetworkDisabled = errors.New("network is disabled, please enable it in the configuration")
)

func checkOptions(options *types.RunnerOptions) error {
	configuration := static.GetAiexecSandboxGlobalConfigurations()

	if options.EnableNetwork && !configuration.EnableNetwork {
		return ErrNetworkDisabled
	}

	return nil
}
