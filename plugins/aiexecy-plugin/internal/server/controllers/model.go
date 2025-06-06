package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/service"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/pkg/entities/plugin_entities"
	"aiexec-plugin/pkg/entities/requests"
)

func InvokeLLM(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestInvokeLLM]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.InvokeLLM(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func InvokeTextEmbedding(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestInvokeTextEmbedding]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.InvokeTextEmbedding(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func InvokeRerank(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestInvokeRerank]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.InvokeRerank(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func InvokeTTS(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestInvokeTTS]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.InvokeTTS(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func InvokeSpeech2Text(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestInvokeSpeech2Text]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.InvokeSpeech2Text(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func InvokeModeration(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestInvokeModeration]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.InvokeModeration(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func ValidateProviderCredentials(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestValidateProviderCredentials]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.ValidateProviderCredentials(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func ValidateModelCredentials(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestValidateModelCredentials]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.ValidateModelCredentials(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func GetTTSModelVoices(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestGetTTSModelVoices]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.GetTTSModelVoices(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func GetTextEmbeddingNumTokens(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestGetTextEmbeddingNumTokens]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.GetTextEmbeddingNumTokens(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func GetLLMNumTokens(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestGetLLMNumTokens]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.GetLLMNumTokens(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func GetAIModelSchema(config *app.Config) gin.HandlerFunc {
	type request = plugin_entities.InvokePluginRequest[requests.RequestGetAIModelSchema]

	return func(c *gin.Context) {
		BindPluginDispatchRequest(
			c,
			func(itr request) {
				service.GetAIModelSchema(&itr, c, config.PluginMaxExecutionTimeout)
			},
		)
	}
}

func ListModels(c *gin.Context) {
	BindRequest(c, func(request struct {
		TenantID string `uri:"tenant_id" validate:"required"`
		Page     int    `form:"page" validate:"required,min=1"`
		PageSize int    `form:"page_size" validate:"required,min=1,max=256"`
	}) {
		c.JSON(http.StatusOK, service.ListModels(request.TenantID, request.Page, request.PageSize))
	})
}
