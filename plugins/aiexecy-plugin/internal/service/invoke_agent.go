package service

import (
	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/core/plugin_daemon"
	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/core/session_manager"
	"aiexec-plugin/internal/utils/stream"
	"aiexec-plugin/pkg/entities/agent_entities"
	"aiexec-plugin/pkg/entities/plugin_entities"
	"aiexec-plugin/pkg/entities/requests"
)

func InvokeAgentStrategy(
	r *plugin_entities.InvokePluginRequest[requests.RequestInvokeAgentStrategy],
	ctx *gin.Context,
	max_timeout_seconds int,
) {
	baseSSEWithSession(
		func(session *session_manager.Session) (*stream.Stream[agent_entities.AgentStrategyResponseChunk], error) {
			return plugin_daemon.InvokeAgentStrategy(session, &r.Data)
		},
		access_types.PLUGIN_ACCESS_TYPE_AGENT_STRATEGY,
		access_types.PLUGIN_ACCESS_ACTION_INVOKE_AGENT_STRATEGY,
		r,
		ctx,
		max_timeout_seconds,
	)
}
