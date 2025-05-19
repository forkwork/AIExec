package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/core/plugin_daemon"
	"aiexec-plugin/internal/core/plugin_daemon/access_types"
	"aiexec-plugin/internal/core/session_manager"
	"aiexec-plugin/internal/utils/stream"
	"aiexec-plugin/pkg/entities/oauth_entities"
	"aiexec-plugin/pkg/entities/plugin_entities"
	"aiexec-plugin/pkg/entities/requests"
)

func OAuthGetAuthorizationURL(
	r *plugin_entities.InvokePluginRequest[requests.RequestOAuthGetAuthorizationURL],
	ctx *gin.Context,
	maxExecutionTimeout time.Duration,
) {
	baseSSEWithSession(
		func(session *session_manager.Session) (*stream.Stream[oauth_entities.OAuthGetAuthorizationURLResult], error) {
			return plugin_daemon.OAuthGetAuthorizationURL(session, &r.Data)
		},
		access_types.PLUGIN_ACCESS_TYPE_OAUTH,
		access_types.PLUGIN_ACCESS_ACTION_GET_AUTHORIZATION_URL,
		r,
		ctx,
		int(maxExecutionTimeout.Seconds()),
	)
}

func OAuthGetCredentials(
	r *plugin_entities.InvokePluginRequest[requests.RequestOAuthGetCredentials],
	ctx *gin.Context,
	maxExecutionTimeout time.Duration,
) {
	baseSSEWithSession(
		func(session *session_manager.Session) (*stream.Stream[oauth_entities.OAuthGetCredentialsResult], error) {
			return plugin_daemon.OAuthGetCredentials(session, &r.Data)
		},
		access_types.PLUGIN_ACCESS_TYPE_OAUTH,
		access_types.PLUGIN_ACCESS_ACTION_GET_CREDENTIALS,
		r,
		ctx,
		int(maxExecutionTimeout.Seconds()),
	)
}
