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
//	@Description: 查询所有模型（含分组名称）
//	@receiver m
//	@return []resp.ModelListItem
//	@return error
func (m *model) ModelAll() ([]resp.ModelListItem, error) {
	var results []resp.ModelListItem
	err := database.DB.Table("model md").
		Select("md.id, md.group_id, mg.name as group_name, md.alias, md.name, md.description, md.is_usable, md.icon, md.`order`").
		Joins("LEFT JOIN model_group mg ON mg.id = md.group_id").
		Where("md.deleted_at IS NULL").
		Order("md.`order` asc").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// ModelGroupAll
//
//	@Description: 查询所有模型分组
//	@receiver m
//	@return []models.ModelGroup
//	@return error
func (m *model) ModelGroupAll() ([]models.ModelGroup, error) {
	groups := make([]models.ModelGroup, 0)
	if err := database.DB.Model(&models.ModelGroup{}).Order("`order` asc").Scan(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// ModelRelationTypeAll
//
//	@Description: 查询所有模型关系类型
//	@receiver m
//	@return []models.ModelRelationType
//	@return error
func (m *model) ModelRelationTypeAll() ([]models.ModelRelationType, error) {
	types := make([]models.ModelRelationType, 0)
	if err := database.DB.Model(&models.ModelRelationType{}).Scan(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

// ModelRelationAll
//
//	@Description: 查询所有模型关系
//	@receiver m
//	@return []models.ModelRelation
//	@return error
func (m *model) ModelRelationAll() ([]models.ModelRelation, error) {
	relations := make([]models.ModelRelation, 0)
	if err := database.DB.Model(&models.ModelRelation{}).Scan(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}

// ModelSingle
//
//	@Description: 查询单个模型（含关联名称）
//	@receiver m
//	@param alias
//	@return *resp.ModelInfo
//	@return error
func (m *model) ModelSingle(alias string) (*resp.ModelInfo, error) {
	// 模型
	var model models.Model
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"alias": alias}).
		Scan(&model).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	id := model.ID

	// 模型关系（含名称）
	modelRelations := make([]resp.ModelRelationInfo, 0)
	err := database.DB.Table("model_relation mr").
		Select(`mr.id, mr.source_id, ms.name as source_name, mr.target_id, mt.name as target_name,
			mr.type_id, mrt.name as type_name, mrt.s2t, mrt.t2s, mr.description`).
		Joins("LEFT JOIN model ms ON ms.id = mr.source_id").
		Joins("LEFT JOIN model mt ON mt.id = mr.target_id").
		Joins("LEFT JOIN model_relation_type mrt ON mrt.id = mr.type_id").
		Where("mr.source_id = ? OR mr.target_id = ?", id, id).
		Scan(&modelRelations).Error
	if err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	// 模型字段分组
	modelFieldGroups := make([]models.ModelFieldGroup, 0)
	if err := database.DB.Model(&models.ModelFieldGroup{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFieldGroups).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	// 模型唯一字段
	modelFieldUniques := make([]models.ModelFieldUnique, 0)
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFieldUniques).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	// 模型字段关联
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": id}).
		Or(map[string]any{"target_model_id": id}).
		Scan(&modelFieldRelations).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	// 模型字段
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"model_id": id}).
		Scan(&modelFields).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	result := &resp.ModelInfo{
		Model:              model,
		ModelRelation:      modelRelations,
		ModelFieldGroup:    modelFieldGroups,
		ModelFieldUnique:   modelFieldUniques,
		ModelFieldRelation: modelFieldRelations,
		ModelField:         modelFields,
	}
	return result, nil
}
