package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"wosaas/library/constant"
	"wosaas/library/errno"
	"wosaas/library/sessionlib"
	"wosaas/library/types"
	"wosaas/library/types/response"
)

// Auth 校验登录态，并将用户信息放入上下文。
//
// Token 来源优先级：Authorization: Bearer xxx -> X-Token -> token 参数/Cookie。
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := tokenFromRequest(c)
		if token == "" {
			response.AbortJSON(c, errno.SessionInvalidErr)
			return
		}

		claims, err := sessionlib.Default.Parse(token)
		if err != nil {
			response.AbortJSON(c, errno.SessionInvalidErr)
			return
		}

		setUserInfo(c, types.UserInfo{
			UserID:   claims.UserID,
			UserName: claims.UserName,
			Mobile:   claims.Mobile,
			Role:     claims.Role,
			SpID:     claims.SpID,
			CpID:     claims.CpID,
			FleetID:  claims.FleetID,
			CityID:   claims.CityID,
		})
		c.Next()
	}
}

// OptionalAuth 有 token 则解析，无 token 也放行。
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := tokenFromRequest(c); token != "" {
			if claims, err := sessionlib.Default.Parse(token); err == nil {
				setUserInfo(c, types.UserInfo{
					UserID:   claims.UserID,
					UserName: claims.UserName,
					Mobile:   claims.Mobile,
					Role:     claims.Role,
					SpID:     claims.SpID,
					CpID:     claims.CpID,
					FleetID:  claims.FleetID,
					CityID:   claims.CityID,
				})
			}
		}
		c.Next()
	}
}

// GetUserInfo 从上下文读取用户信息。
func GetUserInfo(c *gin.Context) (types.UserInfo, bool) {
	if v, ok := c.Get(constant.ContextKeyUserInfo); ok {
		if info, ok := v.(types.UserInfo); ok {
			return info, true
		}
	}
	return types.UserInfo{}, false
}

func setUserInfo(c *gin.Context, info types.UserInfo) {
	ctx := context.WithValue(c.Request.Context(), constant.ContextKeyUserInfo, info)
	c.Request = c.Request.WithContext(ctx)
	c.Set(constant.ContextKeyUserInfo, info)
}

func tokenFromRequest(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	if token := c.GetHeader("X-Token"); token != "" {
		return token
	}
	if token := c.Query("token"); token != "" {
		return token
	}
	if cookie, err := c.Cookie("WOSAAS_SESSION"); err == nil && cookie != "" {
		return cookie
	}
	return ""
}
