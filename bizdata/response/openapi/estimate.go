package openapi

// EstimateResp 运力预估响应。
type EstimateResp struct {
	DistanceMeter   int64   `json:"distance_meter"`
	DurationMinute  int64   `json:"duration_minute"`
	EstimatedAmount float64 `json:"estimated_amount"`
	Currency        string  `json:"currency"`
}
