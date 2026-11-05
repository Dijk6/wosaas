package manage

// ExportTaskAddReq 创建导出任务。
type ExportTaskAddReq struct {
	SrcPage      string `json:"src_page" form:"src_page" binding:"required"`
	SQLContent   string `json:"sql_content" form:"sql_content" binding:"required"`
	SQLCount     int64  `json:"sql_count" form:"sql_count"`
	TableName    string `json:"table_name" form:"table_name"`
	ColumnsName  string `json:"columns_name" form:"columns_name" binding:"required"`
	DatabaseName string `json:"database_name" form:"database_name" binding:"required"`
}

// ExportTaskListReq 查询导出任务列表。
type ExportTaskListReq struct {
	SrcPage string `json:"src_page" form:"src_page"`
	Page    int    `json:"page" form:"page"`
	Size    int    `json:"size" form:"size"`
}
