package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/resp"
	"gcmdb/app/cmdb/utils"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type modelFieldRelation struct {
}

var ModelFieldRelation = new(modelFieldRelation)

// CreateModelFieldRelation
//
//	@Description: 创建模型字段关联
//	@receiver m
//	@param c
func (m *modelFieldRelation) CreateModelFieldRelation(c *gin.Context) {
	var body models.ModelFieldRelation
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
		return
	}
	// 校验源模型存在
	var srcModelCount, tgtModelCount int64
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": body.SourceModelId}).Count(&srcModelCount)
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": body.TargetModelId}).Count(&tgtModelCount)
	if srcModelCount == 0 || tgtModelCount == 0 {
		response.Fail(c, "关联模型不存在")
		return
	}
	// 校验字段归属
	var srcField, tgtField models.ModelField
	if err := database.DB.Model(&models.ModelField{}).Where(map[string]any{"id": body.SourceFieldId, "model_id": body.SourceModelId}).Scan(&srcField).Error; err != nil || srcField.ID == 0 {
		response.Fail(c, "源字段不属于源模型")
		return
	}
	if err := database.DB.Model(&models.ModelField{}).Where(map[string]any{"id": body.TargetFieldId, "model_id": body.TargetModelId}).Scan(&tgtField).Error; err != nil || tgtField.ID == 0 {
		response.Fail(c, "目标字段不属于目标模型")
		return
	}
	// 校验模型关系存在
	var relCount int64
	database.DB.Model(&models.ModelRelation{}).
		Where("(source_id = ? AND target_id = ?) OR (source_id = ? AND target_id = ?)", body.SourceModelId, body.TargetModelId, body.TargetModelId, body.SourceModelId).
		Count(&relCount)
	if relCount == 0 {
		response.Fail(c, "两个模型之间未定义模型关系")
		return
	}
	var results models.ModelFieldRelation
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		FirstOrCreate(&results, body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	// 同步创建实例关系
	if err := utils.InstanceRelation.CreateModelFieldRelation(body.SourceModelId, body.TargetModelId, body.SourceFieldId, body.TargetFieldId); err != nil {
		logger.Error(fmt.Sprintf("创建实例关系失败-%s", err.Error()))
	}
	response.Success(c, "创建成功", nil)
}

// ListModelFieldRelation
//
//	@Description: 从源模型展示所有的模型字段关联
//	@receiver m
//	@param c
func (m *modelFieldRelation) ListModelFieldRelation(c *gin.Context) {
	sourceModelId := c.Param("source_model_id")
	var modelFieldRelations []models.ModelFieldRelation
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where("source_model_id = ? OR target_model_id = ?", sourceModelId, sourceModelId).
		Scan(&modelFieldRelations).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 批量收集所有需要的 model ID 和 field ID，避免 N+1 查询
	modelIds := make(map[uint]bool)
	fieldIds := make(map[uint]bool)
	for _, mfr := range modelFieldRelations {
		modelIds[mfr.SourceModelId] = true
		modelIds[mfr.TargetModelId] = true
		fieldIds[mfr.SourceFieldId] = true
		fieldIds[mfr.TargetFieldId] = true
	}

	// 批量查询 models
	modelIdSlice := make([]uint, 0, len(modelIds))
	for id := range modelIds {
		modelIdSlice = append(modelIdSlice, id)
	}
	modelsList := make([]models.Model, 0)
	if len(modelIdSlice) > 0 {
		database.DB.Model(&models.Model{}).Where("id IN ?", modelIdSlice).Scan(&modelsList)
	}
	modelMap := make(map[uint]models.Model, len(modelsList))
	for _, m := range modelsList {
		modelMap[m.ID] = m
	}

	// 批量查询 fields
	fieldIdSlice := make([]uint, 0, len(fieldIds))
	for id := range fieldIds {
		fieldIdSlice = append(fieldIdSlice, id)
	}
	fieldsList := make([]models.ModelField, 0)
	if len(fieldIdSlice) > 0 {
		database.DB.Model(&models.ModelField{}).Where("id IN ?", fieldIdSlice).Scan(&fieldsList)
	}
	fieldMap := make(map[uint]models.ModelField, len(fieldsList))
	for _, f := range fieldsList {
		fieldMap[f.ID] = f
	}

	// 组装结果
	results := make([]resp.ListModelFieldRelation, 0, len(modelFieldRelations))
	for _, mfr := range modelFieldRelations {
		srcModel := modelMap[mfr.SourceModelId]
		tgtModel := modelMap[mfr.TargetModelId]
		srcField := fieldMap[mfr.SourceFieldId]
		tgtField := fieldMap[mfr.TargetFieldId]
		results = append(results, resp.ListModelFieldRelation{
			ModelFieldRelation: mfr,
			SourceModelDisplay: fmt.Sprintf("%s(%s)", srcModel.Name, srcModel.Alias),
			TargetModelDisplay: fmt.Sprintf("%s(%s)", tgtModel.Name, tgtModel.Alias),
			SourceFieldDisplay: fmt.Sprintf("%s(%s)", srcField.Name, srcField.Alias),
			TargetFieldDisplay: fmt.Sprintf("%s(%s)", tgtField.Name, tgtField.Alias),
		})
	}
	response.Success(c, "查询成功", results)
}

// DeleteModelFieldRelation
//
//	@Description: 删除模型关系实例
//	@receiver m
//	@param c
func (m *modelFieldRelation) DeleteModelFieldRelation(c *gin.Context) {
	id := c.Param("id")
	// 查询模型字段关系实例
	var modelFieldRelation models.ModelFieldRelation
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"id": id}).
		Scan(&modelFieldRelation).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}

	sourceModelId := modelFieldRelation.SourceModelId
	targetModelId := modelFieldRelation.TargetModelId
	sourceFieldId := modelFieldRelation.SourceFieldId
	targetFieldId := modelFieldRelation.TargetFieldId
	// 删除模型字段关系
	if err := database.DB.Unscoped().Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelFieldRelation{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	// 同步删除实例关系
	if err := utils.InstanceRelation.DeleteModelFieldRelation(sourceModelId, targetModelId, sourceFieldId, targetFieldId); err != nil {
		logger.Error(fmt.Sprintf("删除实例关系失败-%s", err.Error()))
	}
	response.Success(c, "删除成功", nil)
}
