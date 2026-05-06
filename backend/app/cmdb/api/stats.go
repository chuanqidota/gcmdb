package api

import (
	"gcmdb/app/cmdb/models"
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
	if err := database.DB.Table("model_group").Count(&data.ModelGroupCount).Error; err != nil {
		response.Fail(c, "查询失败-模型分组数")
		return
	}

	// 模型数
	if err := database.DB.Table("model").Count(&data.ModelCount).Error; err != nil {
		response.Fail(c, "查询失败-模型数")
		return
	}

	// 实例总数 + 每个模型的实例数
	var modelCounts []resp.ModelInstanceCount
	if err := database.DB.Table("instance i").
		Select(`i.model_id, m.name as model_name, m.alias as model_alias,
			COALESCE(mg.name, '') as group_name, COUNT(*) as instance_count`).
		Joins("JOIN model m ON m.id = i.model_id").
		Joins("LEFT JOIN model_group mg ON mg.id = m.group_id").
		Group("i.model_id, m.name, m.alias, mg.name").
		Order("instance_count DESC").
		Scan(&modelCounts).Error; err != nil {
		response.Fail(c, "查询失败-实例统计")
		return
	}
	data.ModelInstanceCounts = modelCounts

	totalInstances := 0
	for _, mc := range modelCounts {
		totalInstances += mc.InstanceCount
	}
	data.InstanceCount = totalInstances

	// 关系总数
	if err := database.DB.Table("instance_relation").Count(&data.RelationCount).Error; err != nil {
		response.Fail(c, "查询失败-关系总数")
		return
	}

	// 模型分组概览（附带每组的模型数和实例数）
	instanceCountByGroup := map[string]int{}
	modelCountByGroup := map[string]int{}
	for _, mc := range modelCounts {
		g := mc.GroupName
		if g == "" {
			g = "未分组"
		}
		instanceCountByGroup[g] += mc.InstanceCount
		modelCountByGroup[g]++
	}

	var groups []models.ModelGroup
	database.DB.Find(&groups)
	groupSummaries := make([]resp.GroupSummary, len(groups))
	for i, g := range groups {
		groupSummaries[i] = resp.GroupSummary{
			ID:            g.ID,
			Name:          g.Name,
			Description:   g.Description,
			ModelCount:    modelCountByGroup[g.Name],
			InstanceCount: instanceCountByGroup[g.Name],
		}
	}
	data.GroupSummaries = groupSummaries

	response.Success(c, "查询成功", data)
}
