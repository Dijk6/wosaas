// Package example 后台管理示例接口。
package example

import (
	"github.com/gin-gonic/gin"

	"wosaas/library/errno"
	"wosaas/library/resource"
	"wosaas/library/types/response"
	"wosaas/servers/httpserver/middleware"
)

// HelloReq 示例请求参数。
type HelloReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}

// Hello 是 /manage/example/hello 的示例接口。
//
// 调用示例：
//
//	curl "http://127.0.0.1:8600/manage/example/hello?name=wosaas"
//	curl -X POST http://127.0.0.1:8600/manage/example/hello \
//	  -H 'Content-Type: application/json' \
//	  -d '{"name":"wosaas"}'
func Hello(c *gin.Context) {
	var req HelloReq
	if err := c.ShouldBind(&req); err != nil {
		resource.LoggerService.WithError(err).Warn("bind manage example params failed")
		response.JSON(c, nil, errno.ServerParamsErr)
		return
	}

	data := gin.H{
		"message": "hello " + req.Name,
	}

	// 如果带了登录态，可以把用户信息一起返回，演示中间件上下文用法。
	if info, ok := middleware.GetUserInfo(c); ok {
		data["user_id"] = info.UserID
	}

	response.OK(c, data)
}
