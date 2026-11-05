package constant

// 上下文键。
const (
	ContextKeyTraceID  = "traceId"
	ContextKeyUserInfo = "userinfo"
	ContextKeyClientIP = "ClientIP"
	ContextKeyURI      = "URI"
	ContextKeyExt      = "__ext"

	HeaderTraceID = "X-Request-Id"
	HeaderRealIP  = "X-Real-Ip"
)
