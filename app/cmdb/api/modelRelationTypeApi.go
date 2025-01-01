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

type modelRelationType struct {
}

var ModelRelationType = new(modelRelationType)

// CreateModelRelationType
//
//	@Description: 创建模型管理类型
//	@receiver m
//	@param c
func (m *modelRelationType) CreateModelRelationType(c *gin.Context) {
	var body models.ModelRelationType
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	if err := database.DB.Model(&models.ModelRelationType{}).Create(&body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "创建成功", nil)
}

// ListModelRelationType
//
//	@Description: 获取模型管理类型列表
//	@receiver m
//	@param c
func (m *modelRelationType) ListModelRelationType(c *gin.Context) {
	var query params.CommonQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 参数
	search := query.Search
	limit, offset := query.Limit, query.Offset
	if limit == 0 {
		limit = 10
	}
	// 查询
	db := database.DB.Model(&models.ModelRelationType{})
	if search != "" {
		db.Where("name like ?", "%"+search+"%").
			Or("s2t like ?", "%"+search+"%").
			Or("t2s like ?", "%"+search+"%")
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 分页
	var modelRelationTypes []models.ModelRelationType
	if err := db.Limit(limit).Offset(offset).Scan(&modelRelationTypes).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 响应
	results := params.CommonList{
		Count:   count,
		Results: modelRelationTypes,
	}
	response.Success(c, "执行成功", results)
}

// UpdateModelRelationType
//
//	@Description: 更新模型管理类型
//	@receiver m
//	@param c
func (m *modelRelationType) UpdateModelRelationType(c *gin.Context) {
	relationId := c.Param("id")
	var body models.ModelRelationType
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"updated_at": time.Now(),
		"name":       body.Name,
		"s2t":        body.S2T,
		"t2s":        body.T2S,
	}
	if err := database.DB.Model(&models.ModelRelationType{}).
		Where(map[string]any{"id": relationId}).
		Updates(data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "更新成功", nil)
}

// DeleteModelRelationType
//
//	@Description: 删除模型管理类型
//	@receiver m
//	@param c
func (m *modelRelationType) DeleteModelRelationType(c *gin.Context) {
	typeId := c.Param("id")
	// 判断是否存在模型关系
	var count int64
	if err := database.DB.Model(&models.ModelRelation{}).
		Where(map[string]any{"type_id": typeId}).
		Count(&count).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if count > 0 {
		response.Fail(c, fmt.Sprintf("该类型存在关联关系，无法删除"))
		return
	}
	// 删除
	if err := database.DB.Unscoped().Model(&models.ModelRelationType{}).
		Where(map[string]any{"id": typeId}).
		Delete(&models.ModelRelationType{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
