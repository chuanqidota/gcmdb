package utils

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/openapi/resp"
	"gcmdb/pkg/database"
)

type model struct {
}

var Model = new(model)

// ModelAll
//
//	@Description: 查询所有模型
//	@receiver m
//	@return []models.Model
//	@return error
func (m *model) ModelAll() ([]models.Model, error) {
	_models := make([]models.Model, 0)
	if err := database.DB.Model(&models.Model{}).
		Scan(&_models).Error; err != nil {
		return nil, err
	}
	return _models, nil
}

// ModelSingle
//
//	@Description: 查询单个模型
//	@receiver m
//	@param id
//	@return *resp.ModelInfo
//	@return error
func (m *model) ModelSingle(id string) (*resp.ModelInfo, error) {
	var model models.Model
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"id": id}).
		Scan(&model).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	modelRelations := make([]models.ModelRelation, 0)
	if err := database.DB.Model(&models.ModelRelation{}).
		Where(map[string]any{"source_id": id}).
		Or(map[string]any{"target_id": id}).
		Scan(&modelRelations).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	modelFieldGroups := make([]models.ModelFieldGroup, 0)
	if err := database.DB.Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFieldGroups).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	modelFieldUniques := make([]models.ModelFieldUnique, 0)
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFieldUniques).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": id}).
		Or(map[string]any{"target_model_id": id}).
		Scan(&modelFieldRelations).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFields).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	result := &resp.ModelInfo{
		Model:              model,
		ModelRealtion:      modelRelations,
		ModelFieldGroup:    modelFieldGroups,
		ModelFieldUnique:   modelFieldUniques,
		ModelFieldRelation: modelFieldRelations,
		ModelField:         modelFields,
	}
	return result, nil
}
