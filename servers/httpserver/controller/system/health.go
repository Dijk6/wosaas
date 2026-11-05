package system

import (
	"time"

	"github.com/gin-gonic/gin"

	"wosaas/config"
	"wosaas/library/types/response"
)

// Health 健康检查。
func Health(c *gin.Context) {
	appName := ""
	if app := config.App(); app != nil {
		appName = app.Base.AppName
	}
	response.JSON(c, gin.H{
		"status": "ok",
		"app":    appName,
		"time":   time.Now().Format(time.RFC3339),
	}, nil)
}

// Ping 探活。
func Ping(c *gin.Context) {
	response.JSON(c, "pong", nil)
}
