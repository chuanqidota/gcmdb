package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/app/cmdb/resp"
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
	id := c.Query("id")
	// 获取字段分组信息
	var modelFieldGroup models.ModelFieldGroup
	if err := database.DB.Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"id": id}).
		Scan(&modelFieldGroup).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 获取字段分组对应的字段
	fields, err := modelFieldGroup.GetModelFields()
	if err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 响应
	results := resp.RetrieveModelFieldGroup{
		ModelFieldGroup: modelFieldGroup,
		Fields:          fields,
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
	var body params.UpdateModelFieldGroup
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := database.DB.Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"id": id}).
		Updates(map[string]any{"name": body.Name}).Error; err != nil {
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
