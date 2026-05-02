package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
)

type modelFieldUnique struct {
}

var ModelFieldUnique = new(modelFieldUnique)

// CreateModelFieldUnique
//
//	@Description: 创建模型唯一性字段,fields用英文逗号分割
//	@receiver mfu
//	@param c
func (mfu *modelFieldUnique) CreateModelFieldUnique(c *gin.Context) {
	var body models.ModelFieldUnique
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 校验模型存在
	var modelCount int64
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": body.ModelId}).Count(&modelCount)
	if modelCount == 0 {
		response.Fail(c, "关联模型不存在")
		return
	}
	// 校验 fields 中的字段存在
	if body.Fields == "" {
		response.Fail(c, "字段组合不能为空")
		return
	}
	// 判断是否实例
	var instanceCount int64
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"model_id": body.ModelId}).
		Count(&instanceCount).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	if instanceCount > 0 {
		response.Fail(c, fmt.Sprintf("该模型存在实例，无法创建唯一字段"))
		return
	}

	var result models.ModelFieldUnique
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		FirstOrCreate(&result, body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "创建成功", nil)
}

// RetrieveModelFieldUnique
//
//	@Description: 查询指定模型的唯一性字段
//	@receiver mfu
//	@param c
func (mfu *modelFieldUnique) RetrieveModelFieldUnique(c *gin.Context) {
	modelId := c.Param("model_id")
	var modelFieldUniques []models.ModelFieldUnique
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": modelId}).
		Scan(&modelFieldUniques).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", modelFieldUniques)
}

// DeleteModelFieldUnique
//
//	@Description: 删除唯一性校验
//	@receiver mfu
//	@param c
func (mfu *modelFieldUnique) DeleteModelFieldUnique(c *gin.Context) {
	id := c.Param("id")
	// 校验记录存在
	var count int64
	database.DB.Model(&models.ModelFieldUnique{}).Where(map[string]any{"id": id}).Count(&count)
	if count == 0 {
		response.Fail(c, "记录不存在")
		return
	}
	if err := database.DB.Unscoped().Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelFieldUnique{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
