package model

// ExportTask 导出任务表。
type ExportTask struct {
	BaseModel

	TaskID          string `gorm:"column:task_id;size:64;uniqueIndex" json:"task_id"`
	SpID            string `gorm:"column:sp_id;size:64;index" json:"sp_id"`
	CpID            int64  `gorm:"column:cp_id" json:"cp_id"`
	SrcPage         string `gorm:"column:src_page;size:64;index" json:"src_page"`
	SQLContent      string `gorm:"column:sql_content;type:text" json:"-"`
	SQLType         int8   `gorm:"column:sql_type" json:"sql_type"`
	SourceTable     string `gorm:"column:table_name;size:128" json:"table_name"`
	ColumnsName     string `gorm:"column:columns_name;type:text" json:"columns_name"`
	DatabaseName    string `gorm:"column:database_name;size:64" json:"database_name"`
	Status          int8   `gorm:"column:status;index" json:"status"`
	FileSize        int64  `gorm:"column:file_size" json:"file_size"`
	SQLCount        int64  `gorm:"column:sql_count" json:"sql_count"`
	FileURL         string `gorm:"column:file_url;size:512" json:"file_url"`
	FailReason      string `gorm:"column:fail_reason;size:512" json:"fail_reason"`
	CreatorUserID   string `gorm:"column:creator_user_id;size:64" json:"creator_user_id"`
	CreatorUserName string `gorm:"column:creator_user_name;size:64" json:"creator_user_name"`
	ExtraInfo       string `gorm:"column:extra_info;type:text" json:"extra_info"`
	FinishTime      int64  `gorm:"column:finish_time" json:"finish_time"`
}

// TableName 指定导出任务表名。
func (ExportTask) TableName() string {
	return "export_task"
}
