// Package routers 后台管理路由。
package routers

import (
	"github.com/gin-gonic/gin"

	"wosaas/servers/httpserver/controller/manage/example"
	exporttaskctrl "wosaas/servers/httpserver/controller/manage/exporttask"
)

// RegisterManage 注册后台管理接口。
func RegisterManage(r *gin.RouterGroup) {
	r.GET("/example/hello", example.Hello)
	r.POST("/example/hello", example.Hello)

	// 导出任务
	r.POST("/export/task/add", exporttaskctrl.Add)
	r.GET("/export/task/list", exporttaskctrl.List)
	r.POST("/export/task/run", exporttaskctrl.Run)
	r.GET("/export/download/:date/:file", exporttaskctrl.Download)
}
