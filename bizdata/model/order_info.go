package model

// OrderInfo 订单信息表。
type OrderInfo struct {
	ID                  uint64 `gorm:"column:id;primaryKey" json:"id"`
	OrderNo             string `gorm:"column:order_no" json:"order_no"`
	ChannelOrderNo      string `gorm:"column:channel_order_no" json:"channel_order_no"`
	OrderStatus         uint   `gorm:"column:order_status" json:"order_status"`
	SpID                string `gorm:"column:sp_id" json:"sp_id"`
	CpID                int64  `gorm:"column:cp_id" json:"cp_id"`
	AgentID             int64  `gorm:"column:agent_id" json:"agent_id"`
	FleetID             int64  `gorm:"column:fleet_id" json:"fleet_id"`
	ChannelPartnerID    int    `gorm:"column:channel_partner_id" json:"channel_partner_id"`
	ChannelUserID       string `gorm:"column:channel_user_id" json:"channel_user_id"`
	Mobile              string `gorm:"column:mobile" json:"mobile"`
	CreateMobile        string `gorm:"column:create_mobile" json:"create_mobile"`
	TotalPrice          uint   `gorm:"column:total_price" json:"total_price"`
	PayPrice            uint   `gorm:"column:pay_price" json:"pay_price"`
	ChannelPayPrice     uint   `gorm:"column:channel_pay_price" json:"channel_pay_price"`
	DriverTotalPrice    uint   `gorm:"column:driver_total_price" json:"driver_total_price"`
	StartPlace          string `gorm:"column:start_place" json:"start_place"`
	StartAddress        string `gorm:"column:start_address" json:"start_address"`
	StartCityID         int    `gorm:"column:start_city_id" json:"start_city_id"`
	EndPlace            string `gorm:"column:end_place" json:"end_place"`
	EndAddress          string `gorm:"column:end_address" json:"end_address"`
	EndCityID           int    `gorm:"column:end_city_id" json:"end_city_id"`
	StartGeo            string `gorm:"column:start_geo" json:"start_geo"`
	EndGeo              string `gorm:"column:end_geo" json:"end_geo"`
	DriverClientVersion string `gorm:"column:driver_client_version" json:"driver_client_version"`
	DriverDevice        int    `gorm:"column:driver_device" json:"driver_device"`
	DriverCuid          string `gorm:"column:driver_cuid" json:"driver_cuid"`
	DriverUID           string `gorm:"column:driver_uid" json:"driver_uid"`
	DriverMobile        string `gorm:"column:driver_mobile" json:"driver_mobile"`
	VehicleUID          string `gorm:"column:vehicle_uid" json:"vehicle_uid"`
	VehicleNumber       string `gorm:"column:vehicle_number" json:"vehicle_number"`
	CarType             uint   `gorm:"column:car_type" json:"car_type"`
	CancelType          uint   `gorm:"column:cancel_type" json:"cancel_type"`
	PrePay              uint   `gorm:"column:pre_pay" json:"pre_pay"`
	Status              uint   `gorm:"column:status" json:"status"`
	StartTime           int64  `gorm:"column:start_time" json:"start_time"`
	ConfirmTime         int64  `gorm:"column:confirm_time" json:"confirm_time"`
	PickupTime          int64  `gorm:"column:pickup_time" json:"pickup_time"`
	ArriveTime          int64  `gorm:"column:arrive_time" json:"arrive_time"`
	ReceiveTime         int64  `gorm:"column:receive_time" json:"receive_time"`
	DeliveryTime        int64  `gorm:"column:delivery_time" json:"delivery_time"`
	BillTime            int64  `gorm:"column:bill_time" json:"bill_time"`
	PayTime             int64  `gorm:"column:pay_time" json:"pay_time"`
	CreateTime          int64  `gorm:"column:create_time" json:"create_time"`
	UpdateTime          int64  `gorm:"column:update_time" json:"update_time"`
}

// TableName 指定订单表名。
func (OrderInfo) TableName() string {
	return "order_info"
}
