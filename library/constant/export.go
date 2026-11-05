package constant

// 导出任务状态。
const (
	ExportStatusReady     int8 = 1 // 待执行
	ExportStatusExecuting int8 = 2 // 执行中
	ExportStatusSuccess   int8 = 3 // 成功
	ExportStatusFailed    int8 = 4 // 失败
	ExportStatusRetry     int8 = 5 // 待重试
)

// 导出 SQL 类型。
const (
	ExportSQLTypeMySQL int8 = 1
)

// ExportStatusName 返回状态文案。
func ExportStatusName(status int8) string {
	switch status {
	case ExportStatusReady:
		return "待执行"
	case ExportStatusExecuting:
		return "执行中"
	case ExportStatusSuccess:
		return "导出成功"
	case ExportStatusFailed:
		return "导出失败"
	case ExportStatusRetry:
		return "待重试"
	default:
		return "未知"
	}
}
