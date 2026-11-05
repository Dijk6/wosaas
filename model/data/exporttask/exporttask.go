// Package exporttask 导出任务数据层。
package exporttask

import (
	"context"

	model "wosaas/bizdata/model"
	"wosaas/model/dao/db"
	"wosaas/model/dao/db/saas"
)

const dbName = "saas"

// Filter 查询条件。
type Filter struct {
	SrcPage       string
	SpID          string
	CreatorUserID string
}

func dao() (*saas.ExportTaskDao, error) {
	gdb, err := db.GetDB(dbName)
	if err != nil {
		return nil, err
	}
	return saas.NewExportTaskDao(gdb), nil
}

// Add 创建任务。
func Add(ctx context.Context, task *model.ExportTask) error {
	d, err := dao()
	if err != nil {
		return err
	}
	return d.Create(ctx, task)
}

// GetByTaskID 查询任务。
func GetByTaskID(ctx context.Context, taskID string) (*model.ExportTask, error) {
	d, err := dao()
	if err != nil {
		return nil, err
	}
	return d.GetByTaskID(ctx, taskID)
}

// List 分页查询任务。
func List(ctx context.Context, filter Filter, page, size int) ([]*model.ExportTask, int64, error) {
	d, err := dao()
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	return d.List(ctx, saas.ExportTaskFilter{
		SrcPage:       filter.SrcPage,
		SpID:          filter.SpID,
		CreatorUserID: filter.CreatorUserID,
	}, offset, size)
}

// ListPending 查询待执行任务。
func ListPending(ctx context.Context, limit int) ([]*model.ExportTask, error) {
	d, err := dao()
	if err != nil {
		return nil, err
	}
	return d.ListPending(ctx, limit)
}

// Claim 抢占任务。
func Claim(ctx context.Context, taskID string) (bool, error) {
	d, err := dao()
	if err != nil {
		return false, err
	}
	return d.Claim(ctx, taskID)
}

// Update 更新任务。
func Update(ctx context.Context, task *model.ExportTask) error {
	d, err := dao()
	if err != nil {
		return err
	}
	return d.Update(ctx, task)
}
