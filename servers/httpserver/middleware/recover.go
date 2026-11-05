package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"wosaas/library/errno"
	"wosaas/library/resource"
	"wosaas/library/types/response"
)

func logger() *logrus.Logger {
	if resource.LoggerService != nil {
		return resource.LoggerService
	}
	return logrus.StandardLogger()
}

// Recovery 捕获 handler panic，返回统一错误响应。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger().WithFields(logrus.Fields{
					"trace_id": traceID(c),
					"uri":      c.Request.URL.Path,
					"panic":    fmt.Sprint(rec),
					"stack":    string(debug.Stack()),
				}).Error("panic recovered")

				if !c.Writer.Written() {
					response.AbortJSON(c, errno.ServerErr)
					return
				}
				c.Abort()
			}
		}()

		c.Next()
	}
}
