package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	cmdbUtils "gcmdb/app/cmdb/utils"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
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
	if !models.IsValidFieldType(body.Type) {
		response.Fail(c, fmt.Sprintf("不支持的字段类型:%s, 支持:string/number/bool/date/datetime/json", body.Type))
		return
	}
	// 校验模型存在
	var modelCount int64
	if err := database.DB.Model(&models.Model{}).Where(map[string]any{"id": body.ModelId}).Count(&modelCount).Error; err != nil || modelCount == 0 {
		response.Fail(c, "关联模型不存在")
		return
	}
	// 校验同模型下 alias 唯一
	var aliasCount int64
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"model_id": body.ModelId, "alias": body.Alias}).
		Count(&aliasCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if aliasCount > 0 {
		response.Fail(c, fmt.Sprintf("模型下已存在别名「%s」的字段", body.Alias))
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
		"is_required":    body.IsRequired,
		"is_indexed":     body.IsIndexed,
		"order":          body.Order,
	}
	if err := database.DB.Model(&models.ModelField{}).
		Create(&data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	// 同步为已有实例补充新字段默认值
	value := models.DefaultValueByType(body.Type)
	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = fmt.Sprintf("'%s'", v)
	default:
		valueStr = fmt.Sprintf("%v", v)
	}
	expr := gorm.Expr(fmt.Sprintf("JSON_SET(data, '$.%s', %s)", body.Alias, valueStr))
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"model_id": body.ModelId}).
		Update("data", expr).Error; err != nil {
		logger.Error(fmt.Sprintf("更新字段失败-%s", err.Error()))
	}
	// is_indexed 为 true 时，同步刷新已有实例的 indexed_values
	if body.IsIndexed {
		cmdbUtils.IndexedValues.SyncModelIndexedValues(body.ModelId)
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
		"updated_at":  time.Now(),
		"is_required": body.IsRequired,
		"is_indexed":  body.IsIndexed,
		"name":        body.Name,
		"order":       body.Order,
	}
	if body.FieldGroupId != nil {
		data["field_group_id"] = *body.FieldGroupId
	} else {
		data["field_group_id"] = nil
	}
	// 查询旧的 is_indexed 状态，判断是否需要同步
	var oldField models.ModelField
	database.DB.Model(&models.ModelField{}).Where(map[string]any{"id": id}).Select("model_id", "is_indexed").Scan(&oldField)

	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"id": id}).
		Updates(data).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			response.Fail(c, fmt.Sprintf("更新失败-模型下已存在名称「%s」的字段", body.Name))
		} else {
			response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		}
		return
	}
	// is_indexed 状态变更时，同步刷新该模型所有实例的 indexed_values
	if oldField.IsIndexed != body.IsIndexed {
		cmdbUtils.IndexedValues.SyncModelIndexedValues(oldField.ModelId)
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
	var modelFieldUniques []models.ModelFieldUnique
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": modelId}).
		Scan(&modelFieldUniques).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	for _, modelFieldUnique := range modelFieldUniques {
		fields := strings.Split(modelFieldUnique.Fields, ",")
		if slices.Contains(fields, alias) {
			response.Fail(c, fmt.Sprintf("该字段在模型唯一校验中被引用,无法删除"))
			return
		}
	}
	// 删除模型字段
	if err := database.DB.Unscoped().Model(&models.ModelField{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelField{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	// 同步删除实例中的字段数据
	expr := gorm.Expr(fmt.Sprintf("JSON_REMOVE(data, '$.%s')", alias))
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"model_id": modelId}).
		Update("data", expr).Error; err != nil {
		logger.Error(fmt.Sprintf("删除字段失败-%s", err.Error()))
	}

	response.Success(c, "执行成功", nil)
}
