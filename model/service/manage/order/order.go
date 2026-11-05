// Package order 后台订单服务。
package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	request "wosaas/bizdata/request/manage"
	response "wosaas/bizdata/response/manage"
	"wosaas/library/errno"
	dataorder "wosaas/model/data/orderinfo"
	exportservice "wosaas/model/service/manage/exporttask"
)

const (
	orderDBName       = "saas"
	orderTableName    = "order_info"
	orderExportSrc    = "order_list"
	orderExportHeader = "订单号,订单状态,乘客手机号,司机ID,车牌号,起点,终点,订单金额(元),实付金额(元),创建时间"
)

type orderStatusOption struct {
	Code uint
	Name string
}

var orderStatusOptions = []orderStatusOption{
	{0, "已创建"},
	{8, "接单超时"},
	{9, "接单前取消"},
	{10, "已接单"},
	{15, "司机出发"},
	{19, "到达前取消"},
	{20, "司机到达"},
	{29, "行程开始前取消"},
	{30, "行程开始"},
	{38, "行程结束前取消"},
	{39, "行程结束"},
	{40, "确认账单"},
	{49, "支付前取消"},
	{50, "已支付"},
}

// StatusName 返回订单状态文案。
func StatusName(status uint) string {
	for _, item := range orderStatusOptions {
		if item.Code == status {
			return item.Name
		}
	}
	return "未知"
}

// List 分页查询订单。
func List(ctx context.Context, req *request.OrderListReq) (*response.OrderListResp, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}
	if size > 200 {
		size = 200
	}

	filter := dataorder.Filter{
		Keyword:     strings.TrimSpace(req.Keyword),
		OrderStatus: req.OrderStatus,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	list, total, err := dataorder.List(ctx, filter, page, size)
	if err != nil {
		return nil, errno.ServerErr
	}

	result := &response.OrderListResp{
		List:  make([]*response.OrderItem, 0, len(list)),
		Total: total,
		Page:  page,
		Size:  size,
	}
	for _, item := range list {
		result.List = append(result.List, &response.OrderItem{
			ID:              item.ID,
			OrderNo:         item.OrderNo,
			OrderStatus:     item.OrderStatus,
			OrderStatusName: StatusName(item.OrderStatus),
			SpID:            item.SpID,
			CpID:            item.CpID,
			Mobile:          item.Mobile,
			DriverUID:       item.DriverUID,
			VehicleNumber:   item.VehicleNumber,
			StartPlace:      item.StartPlace,
			EndPlace:        item.EndPlace,
			TotalPrice:      item.TotalPrice,
			PayPrice:        item.PayPrice,
			CreateTime:      formatUnix(item.CreateTime),
		})
	}
	return result, nil
}

// Export 创建订单导出任务。
func Export(ctx context.Context, req *request.OrderExportReq, creatorUserID, creatorUserName string) (string, error) {
	filter := dataorder.Filter{
		Keyword:     strings.TrimSpace(req.Keyword),
		OrderStatus: req.OrderStatus,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	count, err := dataorder.Count(ctx, filter)
	if err != nil {
		return "", errno.ServerErr
	}

	return exportservice.Add(ctx, &request.ExportTaskAddReq{
		SrcPage:      orderExportSrc,
		SQLContent:   buildExportSQL(filter),
		SQLCount:     count,
		TableName:    orderTableName,
		ColumnsName:  orderExportHeader,
		DatabaseName: orderDBName,
	}, creatorUserID, creatorUserName)
}

func buildExportSQL(filter dataorder.Filter) string {
	var b strings.Builder
	b.WriteString("SELECT order_no, ")
	b.WriteString(orderStatusCase())
	b.WriteString(", mobile, driver_uid, vehicle_number, start_place, end_place, ")
	b.WriteString("ROUND(total_price / 100, 2) AS total_price, ROUND(pay_price / 100, 2) AS pay_price, ")
	b.WriteString("FROM_UNIXTIME(create_time, '%Y-%m-%d %H:%i:%s') AS create_time ")
	b.WriteString("FROM order_info WHERE status = 1")

	if filter.Keyword != "" {
		like := sqlQuote("%" + filter.Keyword + "%")
		fmt.Fprintf(&b, " AND (order_no LIKE %s OR mobile LIKE %s OR driver_uid LIKE %s OR vehicle_number LIKE %s)",
			like, like, like, like)
	}
	if filter.OrderStatus != nil {
		fmt.Fprintf(&b, " AND order_status = %d", *filter.OrderStatus)
	}
	if filter.StartTime > 0 {
		fmt.Fprintf(&b, " AND create_time >= %d", filter.StartTime)
	}
	if filter.EndTime > 0 {
		fmt.Fprintf(&b, " AND create_time <= %d", filter.EndTime)
	}
	b.WriteString(" ORDER BY id DESC")
	return b.String()
}

func orderStatusCase() string {
	var b strings.Builder
	b.WriteString("CASE order_status")
	for _, item := range orderStatusOptions {
		fmt.Fprintf(&b, " WHEN %d THEN '%s'", item.Code, item.Name)
	}
	b.WriteString(" ELSE '未知' END AS order_status")
	return b.String()
}

func sqlQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `''`)
	return "'" + s + "'"
}

func formatUnix(ts int64) string {
	if ts <= 0 {
		return ""
	}
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
