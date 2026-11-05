package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wosaas/library/constant"
	"wosaas/library/log"
	"wosaas/library/resource"
)

// Trace 生成/透传 trace id，并写入请求上下文与响应头。
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceIDValue := c.GetHeader(constant.HeaderTraceID)
		if traceIDValue == "" {
			if resource.Tracer != nil {
				traceIDValue = resource.Tracer.NewTraceID().String()
			} else {
				traceIDValue = uuid.NewString()
			}
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, constant.ContextKeyTraceID, traceIDValue)
		ctx = context.WithValue(ctx, constant.ContextKeyClientIP, c.ClientIP())
		ctx = context.WithValue(ctx, constant.ContextKeyURI, c.Request.URL.Path)
		ctx = log.ContextWithExt(ctx)

		c.Request = c.Request.WithContext(ctx)
		c.Set(constant.ContextKeyTraceID, traceIDValue)
		c.Set(constant.ContextKeyClientIP, c.ClientIP())
		c.Header(constant.HeaderTraceID, traceIDValue)

		c.Next()
	}
}

func traceID(c *gin.Context) string {
	if v, ok := c.Get(constant.ContextKeyTraceID); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	if c.Request != nil {
		if s, ok := c.Request.Context().Value(constant.ContextKeyTraceID).(string); ok {
			return s
		}
	}
	return ""
}
