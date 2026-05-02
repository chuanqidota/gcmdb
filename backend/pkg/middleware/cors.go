package middleware

import (
	"gcmdb/config"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := config.Conf.CORS.AllowedOrigins

		if origin != "" && len(allowedOrigins) > 0 {
			if slices.Contains(allowedOrigins, origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else if origin != "" {
			// 未配置白名单时允许所有（向后兼容）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
