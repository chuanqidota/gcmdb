package models

import (
	"fmt"
	"gcmdb/pkg/database"
)

// ModelRelation
// @Description:模型关系
type ModelRelation struct {
	BaseModel
	SourceId    uint   `gorm:"column:source_id;type:uint;not null;comment:源模型id" json:"source_id" binding:"required"`
	TargetId    uint   `gorm:"column:target_id;type:uint;not null;comment:目标模型id" json:"target_id" binding:"required"`
	TypeId      uint   `gorm:"column:type_id;type:uint;not null;comment:模型关系类型id" json:"type_id" binding:"required"`
	Description string `gorm:"column:description;type:text;comment:描述" json:"description"`
}

// TableName
//
//	@Description: 模型关系>表
//	@receiver ModelRelation
//	@return string
func (ModelRelation) TableName() string {
	return "model_relation"
}

// SourceDisplay
//
//	@Description: 源模型展示
//	@receiver mr
//	@return string
//	@return error
func (mr *ModelRelation) SourceDisplay() (string, error) {
	var _model Model
	if err := database.DB.Model(&Model{}).
		Where(map[string]any{"id": mr.SourceId}).
		Scan(&_model).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", _model.Name, _model.Alias), nil
}

// TargetDisplay
//
//	@Description: 目的模型展示
//	@receiver mr
//	@return string
//	@return error
func (mr *ModelRelation) TargetDisplay() (string, error) {
	var _model Model
	if err := database.DB.Model(&Model{}).
		Where(map[string]any{"id": mr.TargetId}).
		Scan(&_model).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", _model.Name, _model.Alias), nil
}

// TypeDisplay
//
//	@Description: 模型关系类型展示
//	@receiver mr
//	@return string
//	@return error
func (mr *ModelRelation) TypeDisplay() (string, error) {
	var _modelRelationType ModelRelationType
	if err := database.DB.Model(&ModelRelationType{}).
		Where(map[string]any{"id": mr.TypeId}).
		Scan(&_modelRelationType).Error; err != nil {
		return "", err
	}
	return _modelRelationType.Name, nil
}
