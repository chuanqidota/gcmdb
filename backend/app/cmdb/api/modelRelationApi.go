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
	// 禁止自关联
	if body.SourceId == body.TargetId {
		response.Fail(c, "源模型和目标模型不能相同")
		return
	}
	// 校验源模型和目标模型存在
	var srcCount, tgtCount int64
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": body.SourceId}).Count(&srcCount)
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": body.TargetId}).Count(&tgtCount)
	if srcCount == 0 || tgtCount == 0 {
		response.Fail(c, "关联模型不存在")
		return
	}
	// 校验关系类型存在
	var typeCount int64
	database.DB.Model(&models.ModelRelationType{}).Where(map[string]any{"id": body.TypeId}).Count(&typeCount)
	if typeCount == 0 {
		response.Fail(c, "关系类型不存在")
		return
	}
	// 判断是否存在（双向去重）
	var modelRelationCount int64
	if err := database.DB.Model(&models.ModelRelation{}).
		Where("(source_id = ? AND target_id = ?) OR (source_id = ? AND target_id = ?)",
			body.SourceId, body.TargetId, body.TargetId, body.SourceId).
		Count(&modelRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if modelRelationCount > 0 {
		response.Fail(c, "源模型和目标模型的关系已存在（含反向）")
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

	// 批量收集所有需要的 model ID 和 type ID
	modelIds := make(map[uint]bool)
	typeIds := make(map[uint]bool)
	for _, mr := range modelRelations {
		modelIds[mr.SourceId] = true
		modelIds[mr.TargetId] = true
		typeIds[mr.TypeId] = true
	}

	// 批量查询 models
	modelIdSlice := make([]uint, 0, len(modelIds))
	for id := range modelIds {
		modelIdSlice = append(modelIdSlice, id)
	}
	models_ := make([]models.Model, 0)
	if len(modelIdSlice) > 0 {
		database.DB.Model(&models.Model{}).Where("id IN ?", modelIdSlice).Scan(&models_)
	}
	modelMap := make(map[uint]models.Model, len(models_))
	for _, m := range models_ {
		modelMap[m.ID] = m
	}

	// 批量查询 relation types
	typeIdSlice := make([]uint, 0, len(typeIds))
	for id := range typeIds {
		typeIdSlice = append(typeIdSlice, id)
	}
	relationTypes := make([]models.ModelRelationType, 0)
	if len(typeIdSlice) > 0 {
		database.DB.Model(&models.ModelRelationType{}).Where("id IN ?", typeIdSlice).Scan(&relationTypes)
	}
	typeMap := make(map[uint]models.ModelRelationType, len(relationTypes))
	for _, t := range relationTypes {
		typeMap[t.ID] = t
	}

	results := make([]map[string]any, 0)
	for _, modelRelation := range modelRelations {
		result := make(map[string]any)
		result["id"] = modelRelation.ID
		result["source_id"] = modelRelation.SourceId
		result["target_id"] = modelRelation.TargetId
		result["type_id"] = modelRelation.TypeId
		if src, ok := modelMap[modelRelation.SourceId]; ok {
			result["source_display"] = fmt.Sprintf("%s(%s)", src.Name, src.Alias)
		}
		if tgt, ok := modelMap[modelRelation.TargetId]; ok {
			result["target_display"] = fmt.Sprintf("%s(%s)", tgt.Name, tgt.Alias)
		}
		if tp, ok := typeMap[modelRelation.TypeId]; ok {
			result["type_display"] = tp.Name
		}
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
