package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/app/cmdb/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
	"time"
)

type model struct {
}

var Model = new(model)

// CreateModel
//
//	@Description:创建模型
//	@receiver m
//	@param c
func (m *model) CreateModel(c *gin.Context) {
	var model models.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
		"group_id":    model.GroupId,
		"alias":       model.Alias,
		"name":        model.Name,
		"description": model.Description,
		"is_usable":   model.IsUsable,
		"icon":        model.Icon,
		"order":       model.Order,
	}
	if err := database.DB.Model(&models.Model{}).
		Create(data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)

}

// ListModel
//
//	@Description: 列出模型
//	@receiver m
//	@param c
func (m *model) ListModel(c *gin.Context) {
	// 参数校验
	var query params.CommonQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	search := query.Search
	limit, offset := query.Limit, query.Offset
	if limit == 0 {
		limit = 10
	}
	// 查询
	db := database.DB.Model(&models.Model{})
	if search != "" {
		db.Where("alias like ?", "%"+search+"%").
			Or("name like ?", "%"+search+"%").
			Or("description like ?", "%"+search+"%")
	}
	// 分页数量
	var count int64
	if err := db.Count(&count).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	//
	var _models []models.Model
	if err := db.Order("order asc").Limit(limit).Offset(offset).Scan(&_models).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 结果响应
	data := params.CommonList{
		Count:   count,
		Results: _models,
	}

	response.Success(c, "执行成功", data)

}

// RetrieveModel
//
//	@Description: 模型详情
//	@receiver m
//	@param c
func (m *model) RetrieveModel(c *gin.Context) {
	modelId := c.Param("id")
	var _model models.Model
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"id": modelId}).
		Scan(&_model).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 模型对应的字段分组
	groups, err := _model.GetModelFieldGroups()
	if err != nil {
		response.Fail(c, fmt.Sprintf("获取字段组失败-%s", err.Error()))
		return
	}
	// 字段分组对应的字段
	groupResults := make([]resp.RetrieveModelGroupField, 0)
	for _, group := range groups {
		modelFields, err := group.GetModelFields()
		if err != nil {
			response.Fail(c, fmt.Sprintf("获取字段失败-%s", err.Error()))
			return
		}
		item := resp.RetrieveModelGroupField{
			Groups:      group,
			ModelFields: modelFields,
		}
		groupResults = append(groupResults, item)
	}
	// 响应
	results := resp.RetrieveModelResp{
		Model:  _model,
		Groups: groupResults,
	}
	response.Success(c, "执行成功", results)
}

// PatchModelGroupId
//
//	@Description: 更新模型分组
//	@receiver m
//	@param c
func (m *model) PatchModelGroupId(c *gin.Context) {
	var body params.PatchModelGroupIdBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"updated_at": time.Now(),
		"group_id":   body.GroupId,
	}
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"id": body.ModelId}).
		Updates(data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "更新成功", nil)
}

// UpdateModel
//
//	@Description: 更新模型信息
//	@receiver m
//	@param c
func (m *model) UpdateModel(c *gin.Context) {
	var body params.PatchModelBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"updated_at":  time.Now(),
		"name":        body.Name,
		"description": body.Description,
		"is_usable":   body.IsUsable,
		"icon":        body.Icon,
		"order":       body.Order,
	}
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"id": body.ID}).
		Updates(data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "更新成功", nil)
}

// DeleteModel
//
//	@Description: 删除模型
//	@receiver m
//	@param c
func (m *model) DeleteModel(c *gin.Context) {
	modelId := c.Param("id")
	var (
		modelRelationCount    int64
		fieldRelationCount    int64
		instanceCount         int64
		instanceRelationCount int64
	)
	// 判断模型关联
	if err := database.DB.Model(&models.ModelRelation{}).
		Where(map[string]any{"source_id": modelId}).
		Or(map[string]any{"target_id": modelId}).
		Count(&modelRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if modelRelationCount > 0 {
		response.Fail(c, fmt.Sprintf("该模型存在关联关系，无法删除"))
		return
	}
	// 判断字段关联
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": modelId}).
		Or(map[string]any{"target_model_id": modelId}).
		Count(&fieldRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if fieldRelationCount > 0 {
		response.Fail(c, fmt.Sprintf("该模型存在字段关联，无法删除"))
		return
	}
	// 判断实例关联
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{"source_model_id": modelId}).
		Or(map[string]any{"target_model_id": modelId}).
		Count(&instanceRelationCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if instanceRelationCount > 0 {
		response.Fail(c, fmt.Sprintf("该模型存在实例关联，无法删除"))
		return
	}
	// 判断实例
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"model_id": modelId}).
		Count(&instanceCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if instanceCount > 0 {
		response.Fail(c, fmt.Sprintf("该模型存在实例，无法删除"))
		return
	}
	// 删除模型字段
	if err := database.DB.Unscoped().Model(&models.ModelField{}).
		Where(map[string]any{"model_id": modelId}).
		Delete(&models.ModelField{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除模型字段失败-%s", err.Error()))
		return
	}
	// 删除模型字段分组
	if err := database.DB.Unscoped().Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"model_id": modelId}).
		Delete(&models.ModelFieldGroup{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除字段分组失败-%s", err.Error()))
		return
	}
	// 删除模型
	if err := database.DB.Unscoped().Model(&models.Model{}).
		Where(map[string]any{"id": modelId}).
		Delete(&models.Model{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除模型失败-%s", err.Error()))
		return
	}

	response.Success(c, "执行成功", nil)
}
