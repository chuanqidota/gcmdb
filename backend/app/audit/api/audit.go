package api

import (
	"fmt"
	auditModels "gcmdb/app/audit/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type audit struct{}

var Audit = new(audit)

// ListAuditLog
//
//	@Description: 查询审计日志列表
//	@receiver a
//	@param c
func (a *audit) ListAuditLog(c *gin.Context) {
	var query struct {
		ResourceType string `form:"resource_type"`
		ResourceID   string `form:"resource_id"`
		Action       string `form:"action"`
		Username     string `form:"username"`
		Keyword      string `form:"keyword"`
		StartTime    string `form:"start_time"`
		EndTime      string `form:"end_time"`
		Limit        int    `form:"limit"`
		Offset       int    `form:"offset"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	db := database.DB.Model(&auditModels.AuditLog{})
	if query.ResourceType != "" {
		db = db.Where("resource_type = ?", query.ResourceType)
	}
	if query.ResourceID != "" {
		db = db.Where("resource_id = ?", query.ResourceID)
	}
	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}
	if query.Username != "" {
		db = db.Where("username = ?", query.Username)
	}
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where("path LIKE ? OR request_body LIKE ? OR resource_id LIKE ?", like, like, like)
	}
	if query.StartTime != "" {
		db = db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime != "" {
		db = db.Where("created_at <= ?", query.EndTime)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}

	var logs []auditModels.AuditLog
	if err := db.Order("id desc").Limit(query.Limit).Offset(query.Offset).Scan(&logs).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}

	response.Success(c, "查询成功", map[string]any{
		"count":   count,
		"results": logs,
	})
}

// RetrieveAuditLog
//
//	@Description: 查询某资源的审计历史
//	@receiver a
//	@param c
func (a *audit) RetrieveAuditLog(c *gin.Context) {
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")

	var logs []auditModels.AuditLog
	if err := database.DB.Model(&auditModels.AuditLog{}).
		Where(map[string]any{"resource_type": resourceType, "resource_id": resourceID}).
		Order("id desc").
		Scan(&logs).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}

	response.Success(c, "查询成功", logs)
}
