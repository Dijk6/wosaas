package order

import (
	"github.com/gin-gonic/gin"

	request "wosaas/bizdata/request/manage"
	"wosaas/library/errno"
	"wosaas/library/types/response"
	orderservice "wosaas/model/service/manage/order"
	"wosaas/servers/httpserver/middleware"
)

// Export 创建订单导出任务。
func Export(c *gin.Context) {
	var req request.OrderExportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, nil, errno.ServerParamsErr)
		return
	}

	info, _ := middleware.GetUserInfo(c)
	taskID, err := orderservice.Export(c.Request.Context(), &req, info.UserID, info.UserName)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}
	response.OK(c, gin.H{"task_id": taskID})
}
