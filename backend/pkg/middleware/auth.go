package middleware

import (
	"gcmdb/config"
	"gcmdb/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := config.Conf.Server.ApiKey
		if apiKey == "" {
			c.Next()
			return
		}
		token := c.GetHeader("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token != apiKey {
			response.FailWithStatus(c, http.StatusUnauthorized, "未授权访问")
			c.Abort()
			return
		}
		c.Next()
	}
}
