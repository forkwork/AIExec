package serverless_runtime

import (
	"net/http"

	"aiexec-plugin/internal/core/plugin_manager/basic_runtime"
	"aiexec-plugin/internal/utils/mapping"
	"aiexec-plugin/pkg/entities"
	"aiexec-plugin/pkg/entities/plugin_entities"
)

type AWSPluginRuntime struct {
	basic_runtime.BasicChecksum
	plugin_entities.PluginRuntime

	// access url for the lambda function
	LambdaURL  string
	LambdaName string

	// listeners mapping session id to the listener
	listeners mapping.Map[string, *entities.Broadcast[plugin_entities.SessionMessage]]

	client *http.Client
}
