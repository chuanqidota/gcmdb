package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
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
	if err := database.DB.Model(&models.Model{}).Create(&model).Error; err != nil {
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
	db := database.DB.Model(&models.Model{})
	if search != "" {
		db = db.Where("alias like ?", "%"+search+"%").
			Or("name like ?", "%"+search+"%").
			Or("description like ?", "%"+search+"%")
	}
	var results []models.Model
	if err := db.Limit(limit).Offset(offset).Scan(&results).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", results)

}

// RetrieveModel
//
//	@Description: 模型详情
//	@receiver m
//	@param c
func (m *model) RetrieveModel(c *gin.Context) {
	modelId := c.Param("id")
	var result models.Model
	if err := database.DB.Model(&models.Model{}).Where("id = ?", modelId).Scan(&result).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", result)
}

func (m *model) UpdateModel(c *gin.Context) {

}

// DeleteModel
//
//	@Description: 删除模型
//	@receiver m
//	@param c
func (m *model) DeleteModel(c *gin.Context) {
	modelId := c.Param("id")
	if err := database.DB.Model(&models.Model{}).Where("id = ?", modelId).Delete(&models.Model{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)

}
