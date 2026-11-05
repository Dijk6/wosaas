package saas

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	model "wosaas/bizdata/model"
	"wosaas/library/constant"
)

// ExportTaskFilter 导出任务查询条件。
type ExportTaskFilter struct {
	SrcPage       string
	SpID          string
	CreatorUserID string
}

// ExportTaskDao 导出任务数据访问对象。
type ExportTaskDao struct {
	db *gorm.DB
}

// NewExportTaskDao 创建导出任务 DAO。
func NewExportTaskDao(db *gorm.DB) *ExportTaskDao {
	return &ExportTaskDao{db: db}
}

// Create 创建导出任务。
func (d *ExportTaskDao) Create(ctx context.Context, task *model.ExportTask) error {
	return d.db.WithContext(ctx).Create(task).Error
}

// GetByTaskID 按任务 ID 查询。
func (d *ExportTaskDao) GetByTaskID(ctx context.Context, taskID string) (*model.ExportTask, error) {
	var task model.ExportTask
	err := d.db.WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List 分页查询导出任务。
func (d *ExportTaskDao) List(ctx context.Context, filter ExportTaskFilter, offset, limit int) ([]*model.ExportTask, int64, error) {
	query := d.db.WithContext(ctx).Model(&model.ExportTask{})
	if filter.SrcPage != "" {
		query = query.Where("src_page = ?", filter.SrcPage)
	}
	if filter.SpID != "" {
		query = query.Where("sp_id = ?", filter.SpID)
	}
	if filter.CreatorUserID != "" {
		query = query.Where("creator_user_id = ?", filter.CreatorUserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*model.ExportTask
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListPending 查询待执行/待重试任务。
func (d *ExportTaskDao) ListPending(ctx context.Context, limit int) ([]*model.ExportTask, error) {
	var list []*model.ExportTask
	err := d.db.WithContext(ctx).
		Where("status IN ?", []int8{constant.ExportStatusReady, constant.ExportStatusRetry}).
		Where("sql_type = ?", constant.ExportSQLTypeMySQL).
		Order("id ASC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

// Claim 通过条件更新抢占任务，保证多实例下只有一个执行者。
func (d *ExportTaskDao) Claim(ctx context.Context, taskID string) (bool, error) {
	result := d.db.WithContext(ctx).
		Model(&model.ExportTask{}).
		Where("task_id = ? AND status IN ?", taskID, []int8{constant.ExportStatusReady, constant.ExportStatusRetry}).
		Updates(map[string]any{
			"status":     constant.ExportStatusExecuting,
			"updated_at": time.Now(),
		})
	return result.RowsAffected == 1, result.Error
}

// Update 更新任务。
func (d *ExportTaskDao) Update(ctx context.Context, task *model.ExportTask) error {
	return d.db.WithContext(ctx).Save(task).Error
}
