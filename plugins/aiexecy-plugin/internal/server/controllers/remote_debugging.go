package controllers

import (
	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/service"
	"aiexec-plugin/pkg/entities/requests"
)

func GetRemoteDebuggingKey(c *gin.Context) {
	BindRequest(
		c, func(request requests.RequestGetRemoteDebuggingKey) {
			c.JSON(200, service.GetRemoteDebuggingKey(request.TenantID))
		},
	)
}
