package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type modelGroup struct {
}

var ModelGroup = new(modelGroup)

// CreateModelGroup
//
//	@Description: 创建模型分组
//	@receiver m
//	@param c
func (m *modelGroup) CreateModelGroup(c *gin.Context) {
	var group models.ModelGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	if err := database.DB.Model(&models.ModelGroup{}).Create(&group).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "创建成功", nil)
}

// ListModelGroup
//
//	@Description: 模型分组列表
//	@receiver m
//	@param c
func (m *modelGroup) ListModelGroup(c *gin.Context) {
	var _params params.CommonQuery
	if err := c.ShouldBindQuery(&_params); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	search := _params.Search
	limit := _params.Limit
	offset := _params.Offset
	if limit == 0 {
		limit = 10
	}
	db := database.DB.Model(&models.ModelGroup{})

	if search != "" {
		db = db.Where("alias like ?", "%"+search+"%").
			Or("name like ?", "%"+search+"%").
			Or("description like ?", "%"+search+"%")
	}
	var results []models.ModelGroup
	if err := db.Limit(limit).Offset(offset).Scan(&results).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", results)
}

// RetrieveModelGroup
//
//	@Description: 获取具体的一个模型分组
//	@receiver m
//	@param c
func (m *modelGroup) RetrieveModelGroup(c *gin.Context) {
	groupId := c.Param("id")
	var result params.RetrieveModelGroup
	if err := database.DB.Model(&models.ModelGroup{}).
		Where("id = ?", groupId).
		Scan(&result.ModelGroup).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if err := database.DB.Model(&models.Model{}).
		Where("group_id = ?", groupId).
		Scan(&result.Models).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if result.Models == nil {
		result.Models = make([]models.Model, 0)
	}
	response.Success(c, "执行成功", result)

}

// PatchModelGroup
//
//	@Description: 更新模型分组
//	@receiver m
//	@param c
func (m *modelGroup) PatchModelGroup(c *gin.Context) {
	groupId := c.Param("id")
	var body params.ModelGroupBody

	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	if err := database.DB.Model(&models.ModelGroup{}).Where("id = ?", groupId).Updates(&body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "更新成功", nil)
}

// DeleteModelGroup
//
//	@Description: 删除模型分组
//	@receiver m
//	@param c
func (m *modelGroup) DeleteModelGroup(c *gin.Context) {
	groupId := c.Param("id")
	var count int64
	if err := database.DB.Model(&models.Model{}).Where("group_id = ?", groupId).Count(&count).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if count > 0 {
		response.Fail(c, fmt.Sprintf("该分组下存在模型，无法删除"))
		return
	}

	if err := database.DB.Model(&models.ModelGroup{}).Where("id = ?", groupId).Delete(&models.ModelGroup{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)

}
