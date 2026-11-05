package routers

import (
	"github.com/gin-gonic/gin"

	"wosaas/library/types/response"
	"wosaas/servers/httpserver/controller/system"
	"wosaas/servers/httpserver/middleware"
)

// RegisterDriver 注册司机端示例接口。
func RegisterDriver(r *gin.RouterGroup) {
	r.GET("/ping", system.Ping)

	// 可选登录态示例：带合法 Token 时返回用户信息，否则返回未登录。
	r.GET("/profile", func(c *gin.Context) {
		info, ok := middleware.GetUserInfo(c)
		if !ok {
			response.JSON(c, gin.H{"login": false}, nil)
			return
		}
		response.JSON(c, gin.H{"login": true, "user": info}, nil)
	})
}
