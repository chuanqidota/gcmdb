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

	db := database.DB.Model(&models.ModelGroup{})

	if search != "" {
		db = db.Where("alias like ?", "%"+search+"%").
			Or("name like ?", "%"+search+"%").
			Or("description like ?", "%"+search+"%")
	}
	var results []models.ModelGroup
	if err := db.Limit(limit).Offset(offset).Scan(&results).Error;err!=nil{
		response.Fail(c,fmt.Sprintf("查询失败-%s",err.Error()))
		return
	}
	response.Success(c,"执行成功",results)
}

// RetrieveModelGroup
//
//	@Description: 获取具体的一个模型分组
//	@receiver m
//	@param c
func (m *modelGroup) RetrieveModelGroup(c *gin.Context) {

}

// UpdateModelGroup
//
//	@Description: 更新模型分组
//	@receiver m
//	@param c
func (m *modelGroup) UpdateModelGroup(c *gin.Context) {
}

// DeleteModelGroup
//
//	@Description: 删除模型分组
//	@receiver m
//	@param c
func (m *modelGroup) DeleteModelGroup(c *gin.Context) {

}
