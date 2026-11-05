package httpserver

import (
	"time"

	"github.com/gin-gonic/gin"

	"wosaas/config"
	"wosaas/servers/httpserver/controller/system"
	"wosaas/servers/httpserver/middleware"
	"wosaas/servers/httpserver/routers"
)

// NewRouter 构建 HTTP 路由。
//
// 两个分组示例：
//
//	/openapi  公开接口示例
//	/driver   带跨域与可选登录态的分组示例
func NewRouter(cfg config.ServerItem) *gin.Engine {
	setGinMode()

	r := gin.New()

	// 顺序：panic 恢复 -> trace id -> 访问日志 -> 超时控制。
	r.Use(
		middleware.Recovery(),
		middleware.Trace(),
		middleware.AccessLog(),
		middleware.Timeout(time.Duration(cfg.WriteTimeout)*time.Millisecond),
	)

	// 基础探活接口。
	r.GET("/health", system.Health)
	r.GET("/ping", system.Ping)

	routers.RegisterOpenApi(r.Group("/openapi"))
	routers.RegisterDriver(r.Group("/driver", middleware.CORS(), middleware.OptionalAuth()))
	routers.RegisterManage(r.Group("/manage", middleware.OptionalAuth()))

	return r
}

func setGinMode() {
	app := config.App()
	if app == nil {
		return
	}
	switch app.Base.RunMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
