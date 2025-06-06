package server

import (
	"errors"
	"aiexec-plugin/internal/utils/cache"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/db"
	"aiexec-plugin/internal/service"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/internal/types/exception"
	"aiexec-plugin/internal/types/models"
	"aiexec-plugin/internal/utils/log"
	"aiexec-plugin/pkg/entities/endpoint_entities"
	"aiexec-plugin/pkg/entities/plugin_entities"
)

// AiexecPlugin supports register and use endpoint to improve the plugin's functionality
// you can use it to do some magics, looking forward to your imagination, Ciallo～(∠·ω< )⌒
// - Yeuoly

// EndpointHandler is a function type that can be used to handle endpoint requests
type EndpointHandler func(ctx *gin.Context, hookId string, maxExecutionTime time.Duration, path string)

func (app *App) Endpoint(config *app.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		hookId := c.Param("hook_id")
		path := c.Param("path")

		// set X-Original-Host
		if c.Request.Header.Get(endpoint_entities.HeaderXOriginalHost) == "" {
			c.Request.Header.Set(endpoint_entities.HeaderXOriginalHost, c.Request.Host)
		}

		if app.endpointHandler != nil {
			app.endpointHandler(c, hookId, time.Duration(config.PluginMaxExecutionTimeout)*time.Second, path)
		} else {
			app.EndpointHandler(c, hookId, time.Duration(config.PluginMaxExecutionTimeout)*time.Second, path)
		}
	}
}

func (app *App) EndpointHandler(ctx *gin.Context, hookId string, maxExecutionTime time.Duration, path string) {
	endpointCacheKey := strings.Join(
		[]string{
			"hook_id",
			hookId,
		},
		":",
	)
	endpoint, err := cache.AutoGetWithGetter[models.Endpoint](
		endpointCacheKey,
		func() (*models.Endpoint, error) {
			v, err := db.GetOne[models.Endpoint](
				db.Equal("hook_id", hookId),
			)
			return &v, err
		})
	if err == db.ErrDatabaseNotFound {
		ctx.JSON(404, exception.BadRequestError(errors.New("endpoint not found")).ToResponse())
		return
	}

	if err != nil {
		log.Error("get endpoint error %v", err)
		ctx.JSON(500, exception.InternalServerError(errors.New("internal server error")).ToResponse())
		return
	}

	// get plugin installation
	pluginInstallationCacheKey := strings.Join(
		[]string{
			"plugin_id",
			endpoint.PluginID,
			"tenant_id",
			endpoint.TenantID,
		},
		":",
	)
	pluginInstallation, err := cache.AutoGetWithGetter[models.PluginInstallation](
		pluginInstallationCacheKey,
		func() (*models.PluginInstallation, error) {
			v, err := db.GetOne[models.PluginInstallation](
				db.Equal("plugin_id", endpoint.PluginID),
				db.Equal("tenant_id", endpoint.TenantID),
			)
			return &v, err
		},
	)
	if err != nil {
		ctx.JSON(404, exception.BadRequestError(errors.New("plugin installation not found")).ToResponse())
		return
	}

	pluginUniqueIdentifier, err := plugin_entities.NewPluginUniqueIdentifier(
		pluginInstallation.PluginUniqueIdentifier,
	)
	if err != nil {
		ctx.JSON(400, exception.UniqueIdentifierError(
			errors.New("invalid plugin unique identifier"),
		).ToResponse())
		return
	}

	// check if plugin exists in current node
	if ok, originalError := app.cluster.IsPluginOnCurrentNode(pluginUniqueIdentifier); !ok {
		app.redirectPluginInvokeByPluginIdentifier(ctx, pluginUniqueIdentifier, originalError)
	} else {
		service.Endpoint(ctx, endpoint, pluginInstallation, maxExecutionTime, path)
	}
}
