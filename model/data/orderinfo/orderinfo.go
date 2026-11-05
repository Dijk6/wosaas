// Package orderinfo 订单数据层。
package orderinfo

import (
	"context"

	model "wosaas/bizdata/model"
	"wosaas/model/dao/db"
	"wosaas/model/dao/db/saas"
)

const dbName = "saas"

// Filter 订单查询条件。
type Filter struct {
	Keyword     string
	OrderStatus *uint
	StartTime   int64
	EndTime     int64
}

func dao() (*saas.OrderInfoDao, error) {
	gdb, err := db.GetDB(dbName)
	if err != nil {
		return nil, err
	}
	return saas.NewOrderInfoDao(gdb), nil
}

func toDAOFilter(filter Filter) saas.OrderInfoFilter {
	return saas.OrderInfoFilter{
		Keyword:     filter.Keyword,
		OrderStatus: filter.OrderStatus,
		StartTime:   filter.StartTime,
		EndTime:     filter.EndTime,
	}
}

// Count 统计订单数量。
func Count(ctx context.Context, filter Filter) (int64, error) {
	d, err := dao()
	if err != nil {
		return 0, err
	}
	return d.Count(ctx, toDAOFilter(filter))
}

// List 分页查询订单。
func List(ctx context.Context, filter Filter, page, size int) ([]*model.OrderInfo, int64, error) {
	d, err := dao()
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	return d.List(ctx, toDAOFilter(filter), offset, size)
}
