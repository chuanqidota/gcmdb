package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"strings"

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
	// 校验现有实例数据是否满足唯一性约束
	if body.Fields != "" {
		fields := strings.Split(body.Fields, ",")
		selectParts := make([]string, 0, len(fields))
		groupParts := make([]string, 0, len(fields))
		for i, field := range fields {
			field = strings.TrimSpace(field)
			selectParts = append(selectParts, fmt.Sprintf("JSON_EXTRACT(data, '$.%s') as v%d", field, i))
			groupParts = append(groupParts, fmt.Sprintf("v%d", i))
		}
		checkSQL := fmt.Sprintf("SELECT 1 FROM (SELECT %s, COUNT(*) as cnt FROM instance WHERE model_id = ? GROUP BY %s HAVING cnt > 1) t LIMIT 1",
			strings.Join(selectParts, ", "), strings.Join(groupParts, ", "))
		var dupCount int64
		database.DB.Raw(checkSQL, body.ModelId).Count(&dupCount)
		if dupCount > 0 {
			response.Fail(c, fmt.Sprintf("现有实例中存在「%s」字段组合重复的数据，无法创建唯一约束", body.Fields))
			return
		}
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
