package api

import (
	auditModels "gcmdb/app/audit/models"
	"gcmdb/app/cmdb/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type stats struct{}

var Stats = new(stats)

func (s *stats) GetStats(c *gin.Context) {
	var data resp.StatsResponse

	// 模型分组数
	database.DB.Table("model_group").Count(&data.ModelGroupCount)

	// 模型数
	database.DB.Table("model").Count(&data.ModelCount)

	// 实例总数 + 每个模型的实例数
	var modelCounts []resp.ModelInstanceCount
	database.DB.Table("instance i").
		Select(`i.model_id, m.name as model_name, m.alias as model_alias,
			COALESCE(mg.name, '') as group_name, COUNT(*) as instance_count`).
		Joins("JOIN model m ON m.id = i.model_id").
		Joins("LEFT JOIN model_group mg ON mg.id = m.group_id").
		Group("i.model_id, m.name, m.alias, mg.name").
		Order("instance_count DESC").
		Scan(&modelCounts)
	data.ModelInstanceCounts = modelCounts

	totalInstances := 0
	for _, mc := range modelCounts {
		totalInstances += mc.InstanceCount
	}
	data.InstanceCount = totalInstances

	// 关系总数
	database.DB.Table("instance_relation").Count(&data.RelationCount)

	// 最近 10 条审计日志
	var logs []auditModels.AuditLog
	database.DB.Order("id DESC").Limit(10).Find(&logs)
	recentLogs := make([]resp.AuditLogBrief, len(logs))
	for i, l := range logs {
		recentLogs[i] = resp.AuditLogBrief{
			ID:           l.ID,
			Action:       l.Action,
			ResourceType: l.ResourceType,
			ResourceID:   l.ResourceID,
			Path:         l.Path,
			Username:     l.Username,
			CreatedAt:    l.CreatedAt,
		}
	}
	data.RecentLogs = recentLogs

	response.Success(c, "查询成功", data)
}
