package exporttask

import (
	"github.com/gin-gonic/gin"

	request "wosaas/bizdata/request/manage"
	"wosaas/library/errno"
	"wosaas/library/types/response"
	exportservice "wosaas/model/service/manage/exporttask"
	"wosaas/servers/httpserver/middleware"
)

// Add 创建导出任务。
func Add(c *gin.Context) {
	var req request.ExportTaskAddReq
	if err := c.ShouldBind(&req); err != nil {
		response.JSON(c, nil, errno.ServerParamsErr)
		return
	}

	info, _ := middleware.GetUserInfo(c)
	taskID, err := exportservice.Add(c.Request.Context(), &req, info.UserID, info.UserName)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}

	response.OK(c, gin.H{"task_id": taskID})
}
