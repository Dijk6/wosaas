package errno

// 框架通用错误码。业务错误码按模块继续追加。
var (
	OK                 = NewCodeError(0, "success")
	ServerErr          = NewCodeError(1000, "服务器内部错误")
	UnLoginErr         = NewCodeError(1001, "未登录！")
	SignatureErr       = NewCodeError(1002, "签名错误")
	RequestTimeoutErr  = NewCodeError(1003, "请求超时，请稍后重试")
	ServerParamsErr    = NewCodeError(1004, "参数错误")
	NotAllowErr        = NewCodeError(1005, "非授权操作")
	RequestTooOftenErr = NewCodeError(1006, "操作过于频繁，请稍后重试")
	NotFoundErr        = NewCodeError(1007, "资源不存在")

	// 登录态相关。
	SessionInvalidErr = NewCodeError(1100, "登录状态失效，请重新登录!")
	LoginFailedErr    = NewCodeError(1101, "账号或密码错误，请重新输入")

	// 后台权限相关。
	ManagePermissionErr = NewCodeError(1200, "请通过后台登录!")

	// 导出相关。
	ExportTaskCountErr    = NewCodeError(1300, "导出条数超过限制")
	ExportTaskSQLErr      = NewCodeError(1301, "导出 SQL 不合法")
	ExportTaskNotFoundErr = NewCodeError(1302, "导出任务不存在")
	ExportTaskDBErr       = NewCodeError(1303, "导出数据库未初始化")
	ExportTaskFileErr     = NewCodeError(1304, "导出文件不存在")
)
