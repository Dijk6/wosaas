package exporttask

import (
	"github.com/gin-gonic/gin"

	request "wosaas/bizdata/request/manage"
	"wosaas/library/errno"
	"wosaas/library/types/response"
	exportservice "wosaas/model/service/manage/exporttask"
)

// List 查询导出任务列表。
func List(c *gin.Context) {
	var req request.ExportTaskListReq
	if err := c.ShouldBind(&req); err != nil {
		response.JSON(c, nil, errno.ServerParamsErr)
		return
	}

	result, err := exportservice.List(c.Request.Context(), &req)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}
	response.OK(c, result)
}
