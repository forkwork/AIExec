package service

import (
	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/core/plugin_daemon"
	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/core/session_manager"
	"aiexec-plugin/internal/utils/stream"
	"aiexec-plugin/pkg/entities/plugin_entities"
	"aiexec-plugin/pkg/entities/requests"
	"aiexec-plugin/pkg/entities/tool_entities"
)

func InvokeTool(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeTool],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	baseSSEWithSession(
		func(session *session_manager.Session) (*stream.Stream[tool_entities.ToolResponseChunk], error) {
			return plugin_daemon.InvokeTool(session, &r.Data)
		},
		access_types.PLUGIN_ACCESS_TYPE_TOOL,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_TOOL,
		r,
		ctx,
		max_timeout_seconds,
	)
}

func ValidateToolCredentials(
	r *plugin_entities.InvokePluginRequest[requests.RequestValidateToolCredentials],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	baseSSEWithSession(
		func(session *session_manager.Session) (*stream.Stream[tool_entities.ValidateCredentialsResult], error) {
			return plugin_daemon.ValidateToolCredentials(session, &r.Data)
		},
		access_types.PLUGIN_ACCESS_TYPE_TOOL,
		access_types.PLUGIN_ACCESS_ACTION_VALIDATE_TOOL_CREDENTIALS,
		r,
		ctx,
		max_timeout_seconds,
	)
}

func GetToolRuntimeParameters(
	r *plugin_entities.InvokePluginRequest[requests.RequestGetToolRuntimeParameters],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	baseSSEWithSession(
		func(session *session_manager.Session) (*stream.Stream[tool_entities.GetToolRuntimeParametersResponse], error) {
			return plugin_daemon.GetToolRuntimeParameters(session, &r.Data)
		},
		access_types.PLUGIN_ACCESS_TYPE_TOOL,
		access_types.PLUGIN_ACCESS_ACTION_GET_TOOL_RUNTIME_PARAMETERS,
		r,
		ctx,
		max_timeout_seconds,
	)
}
