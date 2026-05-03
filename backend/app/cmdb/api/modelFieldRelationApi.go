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
	results := make([]resp.ListModelFieldRelation, 0)
	for _, modelFieldRelation := range modelFieldRelations {
		SourceModelDisplay, _ := modelFieldRelation.GetSourceModelId()
		TargetModelDisplay, _ := modelFieldRelation.GetTargetModelId()
		SourceFieldDisplay, _ := modelFieldRelation.GetSourceFieldId()
		TargetFieldDisplay, _ := modelFieldRelation.GetTargetFieldId()
		results = append(results, resp.ListModelFieldRelation{
			ModelFieldRelation: modelFieldRelation,
			SourceModelDisplay: SourceModelDisplay,
			TargetModelDisplay: TargetModelDisplay,
			SourceFieldDisplay: SourceFieldDisplay,
			TargetFieldDisplay: TargetFieldDisplay,
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
