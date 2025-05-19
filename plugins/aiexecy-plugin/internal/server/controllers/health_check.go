package controllers

import (
	"github.com/gin-gonic/gin"
	"aiexec-plugin/internal/manifest"
	"aiexec-plugin/internal/types/app"
	"aiexec-plugin/internal/utils/routine"
)

func HealthCheck(app *app.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"pool_status": routine.FetchRoutineStatus(),
			"version":     manifest.VersionX,
			"build_time":  manifest.BuildTimeX,
			"platform":    app.Platform,
		})
	}
}
