// Package response 统一 HTTP JSON 响应。
package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"wosaas/library/constant"
	"wosaas/library/errno"
	"wosaas/library/log"
	"wosaas/library/resource"
)

// JSONOutput 通用响应。
type JSONOutput struct {
	ErrNO     int64  `json:"err_no"`
	ErrMsg    string `json:"err_msg"`
	RequestID string `json:"request_id"`
	Data      any    `json:"data"`
}

// JSONOutputNew 新格式响应（部分内部接口使用 code/message）。
type JSONOutputNew struct {
	Code      int64  `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Data      any    `json:"data"`
}

// JSONOutputDianshang 电商侧响应格式。
type JSONOutputDianshang struct {
	ErrNO     int64  `json:"errno"`
	RequestID string `json:"request_id"`
	Data      any    `json:"data"`
}

func traceID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Get(constant.ContextKeyTraceID); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	if v, ok := c.Request.Context().Value(constant.ContextKeyTraceID).(string); ok {
		return v
	}
	return ""
}

func logger() *logrus.Logger {
	if resource.LoggerService != nil {
		return resource.LoggerService
	}
	return logrus.StandardLogger()
}

// JSON 返回通用 JSON 响应。
func JSON(c *gin.Context, data any, err error) {
	code := errno.Cause(err)
	if code != errno.OK {
		logger().WithField("origin-errmsg", strings.ReplaceAll(err.Error(), "\n", ";")).Warn("response error")
	}
	log.AddExtNotice(c.Request.Context(), logger())

	c.JSON(http.StatusOK, JSONOutput{
		ErrNO:     code.ErrCode(),
		ErrMsg:    code.Error(),
		RequestID: traceID(c),
		Data:      data,
	})
}

// JSONOpenApi 返回 OpenAPI JSON 响应。
func JSONOpenApi(c *gin.Context, data any, err error) {
	JSON(c, data, err)
}

// JSONDianshang 返回电商侧 JSON 响应。
func JSONDianshang(c *gin.Context, data any, err error) {
	code := errno.Cause(err)
	if code != errno.OK {
		logger().WithField("origin-errmsg", strings.ReplaceAll(err.Error(), "\n", ";")).Warn("response error")
	}
	c.JSON(http.StatusOK, JSONOutputDianshang{
		ErrNO:     code.ErrCode(),
		RequestID: traceID(c),
		Data:      data,
	})
}

// JSONNew 返回新格式 JSON 响应。
func JSONNew(c *gin.Context, data any, err error) {
	code := errno.Cause(err)
	c.JSON(http.StatusOK, JSONOutputNew{
		Code:      code.ErrCode(),
		Message:   code.Error(),
		RequestID: traceID(c),
		Data:      data,
	})
}

// AbortJSON 返回错误并终止后续 handler。
func AbortJSON(c *gin.Context, err error) {
	JSON(c, nil, err)
	c.Abort()
}

// OK 返回成功响应。
func OK(c *gin.Context, data any) {
	JSON(c, data, errno.OK)
}
