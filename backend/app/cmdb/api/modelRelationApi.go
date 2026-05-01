package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type modelRelation struct {
}

var ModelRelation = new(modelRelation)

// CreateModelRelation
//
//	@Description: 模型创建,所有的都是从源开始
//	@receiver m
//	@param c
func (m *modelRelation) CreateModelRelation(c *gin.Context) {
	var body models.ModelRelation
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 判断是否存在
	var modelRelationCount int64
	if err := database.DB.Model(&models.ModelRelation{}).
		Where(map[string]any{"source_id": body.SourceId, "target_id": body.TargetId}).
		Count(&modelRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if modelRelationCount > 0 {
		response.Fail(c, fmt.Sprintf("源模型和目标模型的关系已经存在-无需创建"))
		return
	}
	// 创建
	if err := database.DB.Model(&models.ModelRelation{}).Create(&body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "创建成功", nil)
}

// ListModelRelation
//
//	@Description: 模型关系展示
//	@receiver m
//	@param c
func (m *modelRelation) ListModelRelation(c *gin.Context) {
	sourceId := c.Query("source_id") // 源模型id

	// 正向+反向
	modelRelations := make([]models.ModelRelation, 0)
	if err := database.DB.Model(&models.ModelRelation{}).
		Where(map[string]any{"source_id": sourceId}).
		Or(map[string]any{"target_id": sourceId}).
		Scan(&modelRelations).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	results := make([]map[string]any, 0)
	for _, modelRelation := range modelRelations {
		result := make(map[string]any)
		result["id"] = modelRelation.ID
		result["source_id"] = modelRelation.SourceId
		result["target_id"] = modelRelation.TargetId
		result["type_id"] = modelRelation.TypeId
		result["source_display"], _ = modelRelation.SourceDisplay()
		result["target_display"], _ = modelRelation.TargetDisplay()
		result["type_display"], _ = modelRelation.TypeDisplay()
		results = append(results, result)
	}
	count := len(results)
	result := resp.CommonList{
		Count:   int64(count),
		Results: results,
	}
	response.Success(c, "查询成功", result)
}

// DeleteModelRelation
//
//	@Description: 删除模型关系
//	@receiver m
//	@param c
func (m *modelRelation) DeleteModelRelation(c *gin.Context) {
	id := c.Param("id")
	// 查询源模型和目标模型的ID
	var modelRelation models.ModelRelation
	if err := database.DB.Model(&models.ModelRelation{}).
		Where(map[string]any{"id": id}).
		Scan(&modelRelation).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询出错-%s", err.Error()))
		return
	}
	sourceModelId, targetModelId := modelRelation.SourceId, modelRelation.TargetId

	// 判断有无模型字段关联
	var modelFiledRelationCount int64
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": sourceModelId, "target_model_id": targetModelId}).
		Count(&modelFiledRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询出错-%s", err.Error()))
		return
	}
	if modelFiledRelationCount > 0 {
		response.Fail(c, fmt.Sprintf("有模型字段关联无法删除"))
		return
	}

	// 判断有无实例关联
	var instanceRelationCount int64
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{"source_model_id": sourceModelId, "target_model_id": targetModelId}).
		Count(&instanceRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询出错-%s", err.Error()))
		return
	}
	if instanceRelationCount > 0 {
		response.Fail(c, fmt.Sprintf("有实例关联无法删除"))
		return
	}

	// 删除
	if err := database.DB.Unscoped().Model(&models.ModelRelation{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelRelation{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
