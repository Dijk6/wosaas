package order

import (
	"github.com/gin-gonic/gin"

	request "wosaas/bizdata/request/manage"
	"wosaas/library/errno"
	"wosaas/library/types/response"
	orderservice "wosaas/model/service/manage/order"
)

// List 查询订单列表。
func List(c *gin.Context) {
	var req request.OrderListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.JSON(c, nil, errno.ServerParamsErr)
		return
	}

	result, err := orderservice.List(c.Request.Context(), &req)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}
	response.OK(c, result)
}
