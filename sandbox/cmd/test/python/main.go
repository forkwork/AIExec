package main

import (
	"fmt"

	"github.com/khulnasoft/aiexec/sandbox/internal/core/runner/python"
	"github.com/khulnasoft/aiexec/sandbox/internal/core/runner/types"
	"github.com/khulnasoft/aiexec/sandbox/internal/service"
	"github.com/khulnasoft/aiexec/sandbox/internal/static"
)

func main() {
	static.InitConfig("conf/config.yaml")
	python.PreparePythonDependenciesEnv()
	resp := service.RunPython3Code(`import json;print(json.dumps({"hello": "world"}))`,
		``,
		&types.RunnerOptions{
			EnableNetwork: true,
		})

	fmt.Println(resp.Data)
}
