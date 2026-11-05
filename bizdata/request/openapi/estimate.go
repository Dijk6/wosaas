package openapi

// EstimateReq 运力预估请求。
type EstimateReq struct {
	StartLng float64 `json:"start_lng" form:"start_lng" binding:"required"`
	StartLat float64 `json:"start_lat" form:"start_lat" binding:"required"`
	EndLng   float64 `json:"end_lng" form:"end_lng" binding:"required"`
	EndLat   float64 `json:"end_lat" form:"end_lat" binding:"required"`
	CityID   int     `json:"city_id" form:"city_id"`
}
