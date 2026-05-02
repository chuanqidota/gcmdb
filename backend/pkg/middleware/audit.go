package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	auditModels "gcmdb/app/audit/models"
	cmdbModels "gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// resourceMap URL路径片段 → 资源类型映射
var resourceMap = map[string]string{
	"models-group":           "model_group",
	"models-field-group":     "model_field_group",
	"models-field-unique":    "model_field_unique",
	"models-field-relation":  "model_field_relation",
	"models-field":           "model_field",
	"models-relation-type":   "model_relation_type",
	"models-relation":        "model_relation",
	"models":                 "model",
	"instance":               "instance",
	"instance-relation":      "instance_relation",
	"instance-relation-by-keys": "instance_relation",
	"search-direct-sql":      "search_direct_sql",
	"sync-instance-relation": "task",
}

// actionMap HTTP方法 → 操作类型映射
var actionMap = map[string]string{
	"POST":   "create",
	"PUT":    "update",
	"PATCH":  "update",
	"DELETE": "delete",
}

// modelRegistry 资源类型 → GORM 模型工厂，用于查询变更前数据
var modelRegistry = map[string]func() any{
	"model_group":          func() any { return &cmdbModels.ModelGroup{} },
	"model":                func() any { return &cmdbModels.Model{} },
	"model_field_group":    func() any { return &cmdbModels.ModelFieldGroup{} },
	"model_field":          func() any { return &cmdbModels.ModelField{} },
	"model_field_unique":   func() any { return &cmdbModels.ModelFieldUnique{} },
	"model_field_relation": func() any { return &cmdbModels.ModelFieldRelation{} },
	"model_relation":       func() any { return &cmdbModels.ModelRelation{} },
	"model_relation_type":  func() any { return &cmdbModels.ModelRelationType{} },
	"instance":             func() any { return &cmdbModels.Instance{} },
	"instance_relation":    func() any { return &cmdbModels.InstanceRelation{} },
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
				if len(requestBody) > 10000 {
					requestBody = requestBody[:10000] + "...(truncated)"
				}
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		// 解析资源信息
		resourceType, resourceID := parseResource(c)

		// update/delete 操作：在执行前查询旧数据
		var beforeData string
		if action == "update" || action == "delete" {
			beforeData = queryBeforeData(c, resourceType, resourceID)
		}

		// 执行后续 handler
		c.Next()

		// 获取操作人
		username := c.GetString("username")

		// 确定 afterData：create/update 使用请求体
		var afterData string
		if action == "create" || action == "update" {
			afterData = requestBody
		}

		// handler 执行完毕，异步写入审计日志
		statusCode := c.Writer.Status()
		sourceIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		go func() {
			auditLog := auditModels.AuditLog{
				ResourceType: resourceType,
				ResourceID:   resourceID,
				Action:       action,
				Method:       method,
				Path:         c.Request.URL.Path,
				RequestBody:  requestBody,
				Username:     username,
				BeforeData:   beforeData,
				AfterData:    afterData,
				SourceIP:     sourceIP,
				UserAgent:    userAgent,
				StatusCode:   statusCode,
			}
			if err := database.DB.Create(&auditLog).Error; err != nil {
				logger.Error(fmt.Sprintf("审计日志写入失败: %s", err.Error()))
			}
		}()
	}
}

// queryBeforeData 查询变更前的数据并序列化为 JSON
func queryBeforeData(c *gin.Context, resourceType, resourceID string) string {
	factory, ok := modelRegistry[resourceType]
	if !ok {
		return ""
	}

	model := factory()

	// 特殊处理：instance-relation-by-keys 通过复合键查询
	if strings.Contains(c.Request.URL.Path, "instance-relation-by-keys") {
		sourceModelId := c.Query("source_model_id")
		targetModelId := c.Query("target_model_id")
		sourceId := c.Query("source_id")
		targetId := c.Query("target_id")
		if sourceModelId == "" || targetModelId == "" || sourceId == "" || targetId == "" {
			return ""
		}
		err := database.DB.Where(map[string]any{
			"source_model_id": sourceModelId,
			"target_model_id": targetModelId,
			"source_id":       sourceId,
			"target_id":       targetId,
		}).First(model).Error
		if err != nil {
			return ""
		}
	} else if resourceID != "" {
		// 按 ID 查询
		if err := database.DB.Where("id = ?", resourceID).First(model).Error; err != nil {
			return ""
		}
	} else {
		return ""
	}

	data, err := json.Marshal(model)
	if err != nil {
		return ""
	}
	return string(data)
}

// parseResource 从请求路径和参数中解析资源类型和资源ID
func parseResource(c *gin.Context) (resourceType, resourceID string) {
	path := c.Request.URL.Path

	// 遍历 resourceMap 匹配路径片段（优先匹配长路径）
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
