package price

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"

	request "wosaas/bizdata/request/openapi"
	"wosaas/library/errno"
	"wosaas/library/resource"
	"wosaas/library/types/response"
	serviceprice "wosaas/model/service/openapi/price"
)

// Estimate 运力预估接口。
func Estimate(c *gin.Context) {
	params := new(request.EstimateReq)
	if err := c.ShouldBind(params); err != nil {
		resource.LoggerService.WithError(err).Warn("bind estimate params failed")
		response.JSONOpenApi(c, nil, errno.ServerParamsErr)
		return
	}

	data, err := serviceprice.GetEstimate(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			response.JSONOpenApi(c, nil, errno.RequestTimeoutErr)
			return
		}
		resource.LoggerService.WithError(err).Warn("get estimate failed")
		response.JSONOpenApi(c, nil, errno.ServerErr)
		return
	}

	response.JSONOpenApi(c, data, errno.OK)
}
