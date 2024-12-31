package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
	"time"
)

type modelField struct {
}

var ModelField = new(modelField)

// CreateModelField
//
//	@Description: 创建模型字段
//	@receiver m
//	@param c
func (m *modelField) CreateModelField(c *gin.Context) {
	var body models.ModelField
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
		"model_id":       body.ModelId,
		"field_group_id": body.FieldGroupId,
		"alias":          body.Alias,
		"name":           body.Name,
		"type":           body.Type,
		"is_unique":      body.IsUnique,
		"is_required":    body.IsRequired,
		"order":          body.Order,
	}
	if err := database.DB.Model(&models.ModelField{}).
		Create(&data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)
}

// RetrieveModelField
//
//	@Description: 根据模型id查询模型字段
//	@receiver m
//	@param c
func (m *modelField) RetrieveModelField(c *gin.Context) {
	modelId := c.Param("model_id")
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"model_id": modelId}).
		Scan(&modelFields).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", modelFields)
}

// UpdateModelField
//
//	@Description: 更新模型字段
//	@receiver m
//	@param c
func (m *modelField) UpdateModelField(c *gin.Context) {
	id := c.Param("id")
	var body params.UpdateModelFieldBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"updated_at":     time.Now(),
		"field_group_id": body.FieldGroupId,
		"is_required":    body.IsRequired,
		"is_unique":      body.IsUnique,
		"name":           body.Name,
		"order":          body.Order,
	}
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"id": id}).
		Updates(data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "更新成功", nil)
}

func (m *modelField) DeleteModelField(c *gin.Context) {
	id := c.Param("id")
	// 判断字段关联中时候在使用
	var modelFieldRelationCount int64
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_field_id": id}).
		Or(map[string]any{"target_field_id": id}).
		Count(&modelFieldRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if modelFieldRelationCount > 0 {
		response.Fail(c, fmt.Sprintf("该字段被引用，无法删除"))
		return
	}
	// 删除模型字段
	if err := database.DB.Unscoped().Model(&models.ModelField{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelField{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)
}
