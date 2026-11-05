package manage

// OrderListReq 订单列表查询条件。
type OrderListReq struct {
	Page        int    `form:"page"`
	Size        int    `form:"size"`
	Keyword     string `form:"keyword"`
	OrderStatus *uint  `form:"order_status"`
	StartTime   int64  `form:"start_time"`
	EndTime     int64  `form:"end_time"`
}

// OrderExportReq 订单导出条件。
type OrderExportReq struct {
	Keyword     string `json:"keyword" form:"keyword"`
	OrderStatus *uint  `json:"order_status" form:"order_status"`
	StartTime   int64  `json:"start_time" form:"start_time"`
	EndTime     int64  `json:"end_time" form:"end_time"`
}
