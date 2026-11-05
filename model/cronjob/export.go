// Package cronjob 定时任务。
package cronjob

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"

	model "wosaas/bizdata/model"
	"wosaas/library/constant"
	"wosaas/library/filestore"
	"wosaas/library/resource"
	exportservice "wosaas/model/service/manage/exporttask"
)

const (
	exportWorkerLimit   = 100
	exportConcurrency   = 5
	exportRowsPerSheet  = 500_000
	exportTaskTimeout   = 10 * time.Minute
	exportFailReasonMax = 500
)

// ExportTaskControl 导出任务执行器。
type ExportTaskControl struct{}

// Run 扫描并执行导出任务。
func (e ExportTaskControl) Run() {
	ctx := context.Background()

	tasks, err := exportservice.ListPending(ctx, exportWorkerLimit)
	if err != nil {
		exportWarn("list pending export tasks failed", err)
		return
	}
	if len(tasks) == 0 {
		return
	}

	sem := make(chan struct{}, exportConcurrency)
	var wg sync.WaitGroup

	for _, task := range tasks {
		claimed, err := exportservice.Claim(ctx, task.TaskID)
		if err != nil {
			exportWarn("claim export task failed: "+task.TaskID, err)
			continue
		}
		if !claimed {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(t *model.ExportTask) {
			defer wg.Done()
			defer func() { <-sem }()

			taskCtx, cancel := context.WithTimeout(ctx, exportTaskTimeout)
			defer cancel()
			e.runTask(taskCtx, t)
		}(task)
	}

	wg.Wait()
}

func (e ExportTaskControl) runTask(ctx context.Context, task *model.ExportTask) {
	rows, err := execSearch(ctx, task)
	if err != nil {
		e.finishFail(ctx, task, "SQL 执行失败: "+err.Error())
		return
	}
	defer rows.Close()

	relativePath, count, size, err := writeExcel(ctx, rows, task)
	if err != nil {
		e.finishFail(ctx, task, "Excel 写入失败: "+err.Error())
		return
	}

	task.Status = constant.ExportStatusSuccess
	task.FileURL = filestore.URL(relativePath)
	task.SQLCount = count
	task.FileSize = size
	task.FinishTime = time.Now().Unix()
	task.FailReason = ""
	if err := exportservice.Update(ctx, task); err != nil {
		exportWarn("update export task success state failed: "+task.TaskID, err)
	}
}

func (e ExportTaskControl) finishFail(ctx context.Context, task *model.ExportTask, reason string) {
	task.Status = constant.ExportStatusFailed
	task.FailReason = truncateRunes(reason, exportFailReasonMax)
	task.FinishTime = time.Now().Unix()
	if err := exportservice.Update(ctx, task); err != nil {
		exportWarn("update export task failed state failed: "+task.TaskID, err)
	}
}

func execSearch(ctx context.Context, task *model.ExportTask) (*sql.Rows, error) {
	if task.DatabaseName == "" {
		return nil, fmt.Errorf("database name is empty")
	}
	db := resource.DB(task.DatabaseName)
	if db == nil {
		return nil, fmt.Errorf("database %s is not initialized", task.DatabaseName)
	}
	return db.WithContext(ctx).Raw(task.SQLContent).Rows()
}

func writeExcel(ctx context.Context, rows *sql.Rows, task *model.ExportTask) (relativePath string, count int64, size int64, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return "", 0, 0, err
	}

	headers := parseHeaders(task.ColumnsName, columns)
	if len(headers) != len(columns) {
		return "", 0, 0, fmt.Errorf("column count mismatch: sql=%d header=%d", len(columns), len(headers))
	}

	relativePath = filepath.Join(time.Now().Format("20060102"), task.TaskID+".xlsx")
	fullPath, err := filestore.FilePath(relativePath)
	if err != nil {
		return "", 0, 0, err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", 0, 0, err
	}

	file := excelize.NewFile()
	defer file.Close()

	sheetIndex := 1
	sheetRow := 1
	var writer *excelize.StreamWriter

	newSheet := func() error {
		name := fmt.Sprintf("Sheet%d", sheetIndex)
		if sheetIndex > 1 {
			if _, err := file.NewSheet(name); err != nil {
				return err
			}
		}
		var err error
		writer, err = file.NewStreamWriter(name)
		if err != nil {
			return err
		}
		headerValues := make([]any, len(headers))
		for i, h := range headers {
			headerValues[i] = h
		}
		return writer.SetRow("A1", headerValues)
	}

	if err := newSheet(); err != nil {
		return "", 0, 0, err
	}

	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(sql.NullString)
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return "", 0, 0, err
		}

		rowValues := make([]any, len(values))
		for i, value := range values {
			rowValues[i] = value.(*sql.NullString).String
		}

		sheetRow++
		cell, _ := excelize.CoordinatesToCellName(1, sheetRow)
		if err := writer.SetRow(cell, rowValues); err != nil {
			return "", 0, 0, err
		}
		count++

		if sheetRow >= exportRowsPerSheet {
			if err := writer.Flush(); err != nil {
				return "", 0, 0, err
			}
			sheetIndex++
			sheetRow = 1
			if err := newSheet(); err != nil {
				return "", 0, 0, err
			}
		}
	}
	if err := rows.Err(); err != nil {
		return "", 0, 0, err
	}

	if err := writer.Flush(); err != nil {
		return "", 0, 0, err
	}
	if err := file.SaveAs(fullPath); err != nil {
		return "", 0, 0, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return "", 0, 0, err
	}
	return relativePath, count, info.Size(), nil
}

func parseHeaders(columnsName string, sqlColumns []string) []string {
	if strings.TrimSpace(columnsName) == "" {
		return sqlColumns
	}

	parts := strings.Split(columnsName, ",")
	headers := make([]string, 0, len(parts))
	for _, part := range parts {
		headers = append(headers, strings.TrimSpace(part))
	}
	return headers
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}

func exportWarn(message string, err error) {
	if resource.LoggerService == nil {
		return
	}
	resource.LoggerService.WithError(err).Warn(message)
}
