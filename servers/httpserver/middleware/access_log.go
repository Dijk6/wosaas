package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AccessLog 打印访问日志。
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		fields := logrus.Fields{
			"trace_id":  traceID(c),
			"method":    c.Request.Method,
			"uri":       c.Request.URL.Path,
			"status":    status,
			"cost_ms":   time.Since(start).Milliseconds(),
			"client_ip": c.ClientIP(),
		}
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			fields["query"] = rawQuery
		}
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		entry := logger().WithFields(fields)
		switch {
		case status >= 500:
			entry.Error("access")
		case status >= 400:
			entry.Warn("access")
		default:
			entry.Info("access")
		}
	}
}
