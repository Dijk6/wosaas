package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"wosaas/library/errno"
	"wosaas/library/types/response"
)

// Timeout 为请求上下文设置超时。
//
// 注意：Go 无法强制终止正在执行的 handler，这里通过 context 通知业务层，
// 并在 handler 返回后兜底返回超时错误。
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeout <= 0 {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if ctx.Err() == context.DeadlineExceeded && !c.Writer.Written() {
			response.AbortJSON(c, errno.RequestTimeoutErr)
		}
	}
}
