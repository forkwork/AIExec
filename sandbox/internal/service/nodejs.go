package service

import (
	"time"

	"aiexec-sandbox/internal/core/runner/nodejs"
	runner_types "aiexec-sandbox/internal/core/runner/types"
	"aiexec-sandbox/internal/static"
	"aiexec-sandbox/internal/types"
)

func RunNodeJsCode(code string, preload string, options *runner_types.RunnerOptions) *types.AiexecSandboxResponse {
	if err := checkOptions(options); err != nil {
		return types.ErrorResponse(-400, err.Error())
	}

	
	if !static.GetAiexecSandboxGlobalConfigurations().EnablePreload {
	    preload = ""
	}
	
	timeout := time.Duration(
		static.GetAiexecSandboxGlobalConfigurations().WorkerTimeout * int(time.Second),
	)

	runner := nodejs.NodeJsRunner{}
	stdout, stderr, done, err := runner.Run(code, timeout, nil, preload, options)
	if err != nil {
		return types.ErrorResponse(-500, err.Error())
	}

	stdout_str := ""
	stderr_str := ""

	defer close(done)
	defer close(stdout)
	defer close(stderr)

	for {
		select {
		case <-done:
			return types.SuccessResponse(&RunCodeResponse{
				Stdout: stdout_str,
				Stderr: stderr_str,
			})
		case out := <-stdout:
			stdout_str += string(out)
		case err := <-stderr:
			stderr_str += string(err)
		}
	}
}
