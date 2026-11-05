package manage

// ExportTaskItem 导出任务列表项。
type ExportTaskItem struct {
	TaskID          string `json:"task_id"`
	SrcPage         string `json:"src_page"`
	Status          int8   `json:"status"`
	StatusName      string `json:"status_name"`
	SQLCount        int64  `json:"sql_count"`
	FileSize        int64  `json:"file_size"`
	FileURL         string `json:"file_url"`
	FailReason      string `json:"fail_reason,omitempty"`
	CreatorUserID   string `json:"creator_user_id"`
	CreatorUserName string `json:"creator_user_name"`
	CreateTime      string `json:"create_time"`
	FinishTime      string `json:"finish_time,omitempty"`
}

// ExportTaskListResp 导出任务列表响应。
type ExportTaskListResp struct {
	List  []*ExportTaskItem `json:"list"`
	Total int64             `json:"total"`
}
