package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/openapi/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type model struct {
}

var Model = new(model)

// GetAllModels
//
//	@Description: 获取所有模型的信息
//	@receiver m
//	@param c
func (m *model) GetAllModels(c *gin.Context) {
	_models := make([]models.Model, 0)
	if err := database.DB.Model(&models.Model{}).
		Scan(&_models).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", _models)
}

func (m *model) ModelInfo(c *gin.Context) {
	id := c.Param("id")
	var model models.Model
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"id": id}).
		Scan(&model).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	modelRelations := make([]models.ModelRelation, 0)
	if err := database.DB.Model(&models.ModelRelation{}).
		Where(map[string]any{"source_id": id}).
		Or(map[string]any{"target_id": id}).
		Scan(&modelRelations).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}

	modelFieldGroups := make([]models.ModelFieldGroup, 0)
	if err := database.DB.Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFieldGroups).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}

	modelFieldUniques := make([]models.ModelFieldUnique, 0)
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFieldUniques).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": id}).
		Or(map[string]any{"target_model_id": id}).
		Scan(&modelFieldRelations).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}

	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFields).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	result := resp.ModelInfo{
		Model:              model,
		ModelRealtion:      modelRelations,
		ModelFieldGroup:    modelFieldGroups,
		ModelFieldUnique:   modelFieldUniques,
		ModelFieldRelation: modelFieldRelations,
		ModelField:         modelFields,
	}
	response.Success(c, "执行成功", result)

}
