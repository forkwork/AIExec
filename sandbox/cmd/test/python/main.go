package main

import (
	"fmt"

	"aiexec-sandbox/internal/core/runner/python"
	"aiexec-sandbox/internal/core/runner/types"
	"aiexec-sandbox/internal/service"
	"aiexec-sandbox/internal/static"
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
