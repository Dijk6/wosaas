package routers

import (
	"github.com/gin-gonic/gin"

	"wosaas/servers/httpserver/controller/openapi/price"
)

// RegisterOpenApi 注册 OpenAPI 路由。
func RegisterOpenApi(r *gin.RouterGroup) {
	partner := r.Group("/partner")

	// 运力预估（示例）。
	partner.POST("/price/estimate", price.Estimate)

	// 后续按模块继续补充：
	//   /price/modestimate  修改目的地预估
	//   /price/ruleinfo     计价规则
	//   /order/create       下单
	//   /order/cancel       取消
	//   /order/confirm      决策通知
	//   /order/detail       订单详情
	//   /pay/notify         支付通知
	//   /invoice/*          发票
	//   /issues/*           工单
}
