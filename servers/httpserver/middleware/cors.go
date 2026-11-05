package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wosaas/config"
)

// CORS 跨域中间件。生产环境请维护精确白名单。
func CORS() gin.HandlerFunc {
	whiteList := []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		allowed := false
		for _, v := range whiteList {
			if v == origin {
				allowed = true
				break
			}
		}

		// 本地调试阶段允许任意来源，线上必须走白名单。
		if !allowed {
			if app := config.App(); app != nil && app.Base.RunMode == "debug" {
				allowed = true
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Token,X-Request-Id")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			if allowed {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}

		c.Next()
	}
}
