// Package price 价格服务示例。
//
// 该文件演示 controller -> service -> (后续 dao) 的分层方式。
package price

import (
	"context"
	"math"

	request "wosaas/bizdata/request/openapi"
	response "wosaas/bizdata/response/openapi"
)

const (
	earthRadiusMeter = 6371000.0
	basePrice        = 8.0
	pricePerKm       = 2.2
	averageSpeedKmH  = 30.0
)

// GetEstimate 计算预估价格。
func GetEstimate(ctx context.Context, req *request.EstimateReq) (*response.EstimateResp, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	distance := haversineMeter(req.StartLng, req.StartLat, req.EndLng, req.EndLat)
	durationMinute := int64(math.Ceil(distance / 1000 / averageSpeedKmH * 60))
	amount := math.Round((basePrice+distance/1000*pricePerKm)*100) / 100

	return &response.EstimateResp{
		DistanceMeter:   int64(math.Round(distance)),
		DurationMinute:  durationMinute,
		EstimatedAmount: amount,
		Currency:        "CNY",
	}, nil
}

func haversineMeter(lng1, lat1, lng2, lat2 float64) float64 {
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLng := (lng2 - lng1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMeter * c
}
