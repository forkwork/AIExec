package middleware

import (
	"github.com/gin-gonic/gin"
	"aiexec-sandbox/internal/static"
)

func Auth() gin.HandlerFunc {
	config := static.GetAiexecSandboxGlobalConfigurations()
	return func(c *gin.Context) {
		if config.App.Key != c.GetHeader("X-Api-Key") {
			c.AbortWithStatus(401)
			return
		}
	}
}
