package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/service"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/pkg/entities/plugin_entities"
	"aiexec-plugin/pkg/entities/requests"
)

func OAuthGetAuthorizationURL(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestOAuthGetAuthorizationURL]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(ipr request) {
				service.OAuthGetAuthorizationURL(
					&ipr,
					c,
					time.Duration(config.PluginMaxExecutionTimeout)*time.Second,
				)
			},
		)
	}
}

func OAuthGetCredentials(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestOAuthGetCredentials]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(c, func(ipr request) {
			service.OAuthGetCredentials(
				&ipr,
				c,
				time.Duration(config.PluginMaxExecutionTimeout)*time.Second,
			)
		})
	}
}
