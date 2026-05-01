package middleware

import (
	"bytes"
	"fmt"
	auditModels "gcmdb/app/audit/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// resourceMap URL路径片段 → 资源类型映射
var resourceMap = map[string]string{
	"models-group":          "model_group",
	"models-field-group":    "model_field_group",
	"models-field-unique":   "model_field_unique",
	"models-field-relation": "model_field_relation",
	"models-field":          "model_field",
	"models-relation-type":  "model_relation_type",
	"models-relation":       "model_relation",
	"models":                "model",
	"instance":              "instance",
	"instance-relation":     "instance_relation",
	"search-direct-sql":     "search_direct_sql",
	"sync-instance-relation": "task",
}

// actionMap HTTP方法 → 操作类型映射
var actionMap = map[string]string{
	"POST":   "create",
	"PUT":    "update",
	"PATCH":  "update",
	"DELETE": "delete",
}

// AuditMiddleware 审计中间件，拦截写操作并记录到 audit_log 表
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		action, isWrite := actionMap[method]

		// 特殊处理: search-direct-sql 的 GET /:id 是执行SQL操作，也需要审计
		if !isWrite && method == "GET" && strings.Contains(c.Request.URL.Path, "search-direct-sql/") {
			action = "execute"
			isWrite = true
		}

		if !isWrite {
			c.Next()
			return
		}

		// 读取请求体并恢复
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) > 0 {
				requestBody = string(bodyBytes)
				// 截断过长的请求体
				if len(requestBody) > 10000 {
					requestBody = requestBody[:10000] + "...(truncated)"
				}
				// 恢复请求体供 handler 使用
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		// 执行后续 handler
		c.Next()

		// handler 执行完毕，异步写入审计日志
		resourceType, resourceID := parseResource(c)
		statusCode := c.Writer.Status()
		sourceIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		go func() {
			auditLog := auditModels.AuditLog{
				ResourceType: resourceType,
				ResourceID:  resourceID,
				Action:      action,
				Method:      method,
				Path:        c.Request.URL.Path,
				RequestBody: requestBody,
				SourceIP:    sourceIP,
				UserAgent:   userAgent,
				StatusCode:  statusCode,
			}
			if err := database.DB.Create(&auditLog).Error; err != nil {
				logger.Error(fmt.Sprintf("审计日志写入失败: %s", err.Error()))
			}
		}()
	}
}

// parseResource 从请求路径和参数中解析资源类型和资源ID
func parseResource(c *gin.Context) (resourceType, resourceID string) {
	path := c.Request.URL.Path

	// 遍历 resourceMap 匹配路径片段
	for fragment, resType := range resourceMap {
		if strings.Contains(path, fragment) {
			resourceType = resType
			break
		}
	}

	// 提取资源 ID：优先从 URL 参数获取
	resourceID = c.Param("id")
	if resourceID == "" {
		resourceID = c.Param("model_id")
	}
	if resourceID == "" {
		resourceID = c.Param("source_id")
	}
	if resourceID == "" {
		resourceID = c.Param("target_id")
	}
	if resourceID == "" {
		resourceID = c.Param("source_model_id")
	}

	return resourceType, resourceID
}
