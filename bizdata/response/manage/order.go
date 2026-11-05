package manage

// OrderItem 订单列表项。
type OrderItem struct {
	ID              uint64 `json:"id"`
	OrderNo         string `json:"order_no"`
	OrderStatus     uint   `json:"order_status"`
	OrderStatusName string `json:"order_status_name"`
	SpID            string `json:"sp_id"`
	CpID            int64  `json:"cp_id"`
	Mobile          string `json:"mobile"`
	DriverUID       string `json:"driver_uid"`
	VehicleNumber   string `json:"vehicle_number"`
	StartPlace      string `json:"start_place"`
	EndPlace        string `json:"end_place"`
	TotalPrice      uint   `json:"total_price"`
	PayPrice        uint   `json:"pay_price"`
	CreateTime      string `json:"create_time"`
}

// OrderListResp 订单列表响应。
type OrderListResp struct {
	List  []*OrderItem `json:"list"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}
