package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	go func(fieldType, field string) {
		value, ok := models.DefaultValueByType[fieldType]
		if !ok {
			value = ""
		}
		expr := fmt.Sprintf("JSON_SET(data,'$.%+v',%+v)", field, value)
		if err := database.DB.Model(&models.Instance{}).
			Where(map[string]any{"model_id": body.ModelId}).
			Update("data", gorm.Expr(expr)).Error; err != nil {
			fmt.Printf("更新字段失败-%s", err.Error())
		}
	}(body.Type, body.Alias)

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

// DeleteModelField
//
//	@Description: 删除模型字段
//	@receiver m
//	@param c
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
	// 获取modelId
	var modelField models.ModelField
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"id": id}).
		Scan(&modelField).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	modelId := modelField.ModelId
	alias := modelField.Alias
	// 判断字段是否在模型字段唯一字段中
	var modelFieldUnique models.ModelFieldUnique
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": modelId}).
		Scan(&modelFieldUnique).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	fields := strings.Split(modelFieldUnique.Fields, ",")
	if slices.Contains(fields, alias) {
		response.Fail(c, fmt.Sprintf("该字段在模型唯一校验中被引用,无法删除"))
		return
	}
	// 删除模型字段
	if err := database.DB.Unscoped().Model(&models.ModelField{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelField{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	// 删除实例里面的字段
	go func(alias string, modelId uint) {
		expr := fmt.Sprintf("JSON_REMOVE(data,'$.%+v')", alias)
		if err := database.DB.Model(&models.Instance{}).
			Where(map[string]any{"model_id": modelId}).
			Update("data", gorm.Expr(expr)).Error; err != nil {
			fmt.Printf("删除失败-%s", err.Error())
		}
	}(alias, modelId)

	response.Success(c, "执行成功", nil)
}
