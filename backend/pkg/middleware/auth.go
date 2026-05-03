package middleware

import (
	"gcmdb/app/auth/models"
	"gcmdb/app/auth/session"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 三种认证中间件的适用场景：
// - SessionAuthMiddleware: 前端页面接口（/v1/cmdb/*），基于 cookie 中的 session ID
// - OpenAPIAuthMiddleware: 对外开放接口（/openapi/*），支持 session 或 Bearer token 双模式
// - CORSMiddleware: 跨域请求处理，放行 OPTIONS 预检请求

// SessionAuthMiddleware 页面接口 session 认证
func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, _ := c.Cookie(session.CookieName())
		if sid == "" {
			response.FailWithStatus(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		data := session.Get(sid)
		if data == nil {
			response.FailWithStatus(c, http.StatusUnauthorized, "会话已过期")
			c.Abort()
			return
		}
		c.Set("user_id", data.UserID)
		c.Set("username", data.Username)
		c.Set("is_admin", data.IsAdmin)
		c.Next()
	}
}

// TokenAuthMiddleware OpenAPI 接口 token 认证
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			response.FailWithStatus(c, http.StatusUnauthorized, "缺少 Authorization header")
			c.Abort()
			return
		}
		var user models.User
		if err := database.DB.Where("token = ? AND is_active = ?", token, true).First(&user).Error; err != nil {
			response.FailWithStatus(c, http.StatusUnauthorized, "无效的 Token")
			c.Abort()
			return
		}
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Next()
	}
}

// OpenAPIAuthMiddleware OpenAPI 接口双重认证：先尝试 session，再尝试 token
func OpenAPIAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 先尝试 session
		if sid, _ := c.Cookie(session.CookieName()); sid != "" {
			if data := session.Get(sid); data != nil {
				c.Set("user_id", data.UserID)
				c.Set("username", data.Username)
				c.Set("is_admin", data.IsAdmin)
				c.Next()
				return
			}
		}
		// 2. 再尝试 token
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			response.FailWithStatus(c, http.StatusUnauthorized, "缺少认证信息")
			c.Abort()
			return
		}
		var user models.User
		if err := database.DB.Where("token = ? AND is_active = ?", token, true).First(&user).Error; err != nil {
			response.FailWithStatus(c, http.StatusUnauthorized, "无效的 Token")
			c.Abort()
			return
		}
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("is_admin", user.IsAdmin)
		c.Next()
	}
}
