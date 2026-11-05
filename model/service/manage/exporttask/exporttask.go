// Package exporttask 导出任务服务。
package exporttask

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	model "wosaas/bizdata/model"
	request "wosaas/bizdata/request/manage"
	response "wosaas/bizdata/response/manage"
	"wosaas/library/constant"
	"wosaas/library/errno"
	dataexport "wosaas/model/data/exporttask"
)

const maxExportRows = 1_000_000

var (
	trailingLimit  = regexp.MustCompile(`(?is)\blimit\s+\d+(\s*,\s*\d+)?\s*$`)
	trailingOffset = regexp.MustCompile(`(?is)\boffset\s+\d+\s*$`)
)

// Add 创建导出任务。
func Add(ctx context.Context, req *request.ExportTaskAddReq, creatorUserID, creatorUserName string) (string, error) {
	if req.SQLCount > maxExportRows {
		return "", errno.ExportTaskCountErr
	}
	if err := ValidateSQL(req.SQLContent); err != nil {
		return "", errno.ExportTaskSQLErr
	}
	if req.DatabaseName == "" || req.ColumnsName == "" {
		return "", errno.ServerParamsErr
	}

	taskID := fmt.Sprintf("export_%s_%s",
		time.Now().Format("20060102150405"),
		strings.ReplaceAll(uuid.NewString(), "-", "")[:8],
	)

	task := &model.ExportTask{
		TaskID:          taskID,
		SQLContent:      req.SQLContent,
		SQLType:         constant.ExportSQLTypeMySQL,
		SourceTable:     req.TableName,
		ColumnsName:     req.ColumnsName,
		DatabaseName:    req.DatabaseName,
		SrcPage:         req.SrcPage,
		Status:          constant.ExportStatusReady,
		SQLCount:        req.SQLCount,
		CreatorUserID:   creatorUserID,
		CreatorUserName: creatorUserName,
	}
	if err := dataexport.Add(ctx, task); err != nil {
		return "", errno.ServerErr
	}
	return taskID, nil
}

// List 查询导出任务列表。
func List(ctx context.Context, req *request.ExportTaskListReq) (*response.ExportTaskListResp, error) {
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

	list, total, err := dataexport.List(ctx, dataexport.Filter{
		SrcPage: req.SrcPage,
	}, page, size)
	if err != nil {
		return nil, errno.ServerErr
	}

	result := &response.ExportTaskListResp{
		List:  make([]*response.ExportTaskItem, 0, len(list)),
		Total: total,
	}
	for _, task := range list {
		item := &response.ExportTaskItem{
			TaskID:          task.TaskID,
			SrcPage:         task.SrcPage,
			Status:          task.Status,
			StatusName:      constant.ExportStatusName(task.Status),
			SQLCount:        task.SQLCount,
			FileSize:        task.FileSize,
			FileURL:         task.FileURL,
			FailReason:      task.FailReason,
			CreatorUserID:   task.CreatorUserID,
			CreatorUserName: task.CreatorUserName,
			CreateTime:      task.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if task.FinishTime > 0 {
			item.FinishTime = time.Unix(task.FinishTime, 0).Format("2006-01-02 15:04:05")
		}
		result.List = append(result.List, item)
	}
	return result, nil
}

// ListPending 查询待执行任务，供 worker 使用。
func ListPending(ctx context.Context, limit int) ([]*model.ExportTask, error) {
	return dataexport.ListPending(ctx, limit)
}

// Claim 抢占任务。
func Claim(ctx context.Context, taskID string) (bool, error) {
	return dataexport.Claim(ctx, taskID)
}

// Update 更新任务。
func Update(ctx context.Context, task *model.ExportTask) error {
	return dataexport.Update(ctx, task)
}

// ValidateSQL 校验导出 SQL。
//
// 只允许单条 SELECT/WITH 查询，禁止分页和危险语句。
func ValidateSQL(sql string) error {
	s := strings.TrimSpace(sql)
	if s == "" {
		return errors.New("sql is empty")
	}
	lower := strings.ToLower(s)
	if !strings.HasPrefix(lower, "select") && !strings.HasPrefix(lower, "with") {
		return errors.New("only select/with sql is allowed")
	}
	if strings.Contains(s, ";") {
		return errors.New("multiple statements are not allowed")
	}
	if trailingLimit.MatchString(s) || trailingOffset.MatchString(s) {
		return errors.New("top-level limit/offset is not allowed")
	}
	if strings.Contains(lower, "into outfile") ||
		strings.Contains(lower, "for update") ||
		strings.Contains(lower, "into dumpfile") {
		return errors.New("dangerous sql is not allowed")
	}
	return nil
}
