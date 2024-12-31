package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
)

type modelFieldGroup struct {
}

var ModelFieldGroup = new(modelFieldGroup)

// CreateModelFieldGroup
//
//	@Description: 创建模型字段分组
//	@receiver m
//	@param c
func (m *modelFieldGroup) CreateModelFieldGroup(c *gin.Context) {
	var body models.ModelFieldGroup
	if err := database.DB.Model(&models.ModelFieldGroup{}).Create(&body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)
}

//
// RetrieveModelFieldGroup
//  @Description: 列出模型字段分组详情
//  @receiver m
//  @param c
//

func (m *modelFieldGroup) RetrieveModelFieldGroup(c *gin.Context) {
	modelId := c.Query("model_id")

	// 查询模型对应的模型字段分组
	modelFieldGroups := make([]models.ModelFieldGroup, 0)
	if err := database.DB.Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"model_id": modelId}).
		Order("order asc").
		Scan(&modelFieldGroups).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 额外补充字段信息
	results := make([]map[string]any, 0)
	for _, modelFieldGroup := range modelFieldGroups {
		result := map[string]any{
			"id":    modelFieldGroup.ID,
			"name":  modelFieldGroup.Name,
			"order": modelFieldGroup.Order,
		}
		modelFields, err := modelFieldGroup.GetModelFields()
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		result["fields"] = modelFields
		results = append(results, result)
	}
	response.Success(c, "执行成功", results)

}

// UpdateModelFieldGroup
//
//	@Description: 更新模型字段-只能更新名称
//	@receiver m
//	@param c
func (m *modelFieldGroup) UpdateModelFieldGroup(c *gin.Context) {
	id := c.Param("id")
	var name string
	if err := c.ShouldBindJSON(&name); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := database.DB.Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"id": id}).
		Updates(map[string]any{"name": name}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "更新成功", nil)
}

// DeleteModelFieldGroup
//
//	@Description: 删除模型字段分组
//	@receiver m
//	@param c
func (m *modelFieldGroup) DeleteModelFieldGroup(c *gin.Context) {
	id := c.Param("id")
	// 判断分组下有无字段
	var modelFieldCount int64
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"field_group_id": id}).
		Count(&modelFieldCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if modelFieldCount > 0 {
		response.Fail(c, fmt.Sprintf("该分组下存在字段，无法删除"))
		return
	}
	// 删除分组
	if err := database.DB.Unscoped().Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelFieldGroup{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
