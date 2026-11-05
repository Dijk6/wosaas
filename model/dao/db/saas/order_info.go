package saas

import (
	"context"

	"gorm.io/gorm"

	model "wosaas/bizdata/model"
)

// OrderInfoFilter 订单查询条件。
type OrderInfoFilter struct {
	Keyword     string
	OrderStatus *uint
	StartTime   int64
	EndTime     int64
}

// OrderInfoDao 订单数据访问对象。
type OrderInfoDao struct {
	db *gorm.DB
}

// NewOrderInfoDao 创建订单 DAO。
func NewOrderInfoDao(db *gorm.DB) *OrderInfoDao {
	return &OrderInfoDao{db: db}
}

func (d *OrderInfoDao) query(ctx context.Context, filter OrderInfoFilter) *gorm.DB {
	query := d.db.WithContext(ctx).Model(&model.OrderInfo{}).Where("status = ?", 1)
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where(
			"(order_no LIKE ? OR mobile LIKE ? OR driver_uid LIKE ? OR vehicle_number LIKE ?)",
			like, like, like, like,
		)
	}
	if filter.OrderStatus != nil {
		query = query.Where("order_status = ?", *filter.OrderStatus)
	}
	if filter.StartTime > 0 {
		query = query.Where("create_time >= ?", filter.StartTime)
	}
	if filter.EndTime > 0 {
		query = query.Where("create_time <= ?", filter.EndTime)
	}
	return query
}

// Count 统计订单数量。
func (d *OrderInfoDao) Count(ctx context.Context, filter OrderInfoFilter) (int64, error) {
	var total int64
	err := d.query(ctx, filter).Count(&total).Error
	return total, err
}

// List 分页查询订单。
func (d *OrderInfoDao) List(ctx context.Context, filter OrderInfoFilter, offset, limit int) ([]*model.OrderInfo, int64, error) {
	var total int64
	if err := d.query(ctx, filter).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*model.OrderInfo
	err := d.query(ctx, filter).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&list).Error
	return list, total, err
}
