package service

import (
	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/core/plugin_daemon/backwards_invocation/transaction"
)

func HandleAWSPluginTransaction(handler *transaction.AWSTransactionHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get session id from the context
		sessionId := c.Request.Header.Get("Aiexec-Plugin-Session-ID")

		handler.Handle(c, sessionId)
	}
}
