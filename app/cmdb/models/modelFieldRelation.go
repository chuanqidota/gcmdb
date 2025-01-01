package models

import (
	"fmt"
	"gcmdb/pkg/database"
)

// ModelFieldRelation
// @Description: 模型字段关联
type ModelFieldRelation struct {
	BaseModel
	SourceModelId uint `gorm:"column:source_model_id;type:uint;not null;comment:源模型id" json:"source_model_id"`
	TargetModelId uint `gorm:"column:target_model_id;type:uint;not null;comment:目标模型id" json:"target_model_id"`
	SourceFieldId uint `gorm:"column:source_field_id;type:uint;not null;comment:源字段id" json:"source_field_id"`
	TargetFieldId uint `gorm:"column:target_field_id;type:uint;not null;comment:目标字段id" json:"target_field_id"`
}

// TableName
//
//	@Description: 模型字段关联>表
//	@receiver ModelFieldRelation
//	@return string
func (ModelFieldRelation) TableName() string {
	return "model_field_relation"
}

// GetSourceModelId
//
//	@Description: 获取源模型信息
//	@receiver mfr
//	@return string
//	@return error
func (mfr *ModelFieldRelation) GetSourceModelId() (string, error) {
	var result Model
	if err := database.DB.Model(&Model{}).
		Where(map[string]any{"id": mfr.SourceModelId}).
		Scan(&result).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", result.Name, result.Alias), nil
}

// GetTargetModelId
//
//	@Description: 获取目标模型信息
//	@receiver mfr
//	@return string
//	@return error
func (mfr *ModelFieldRelation) GetTargetModelId() (string, error) {
	var result Model
	if err := database.DB.Model(&Model{}).
		Where(map[string]any{"id": mfr.TargetModelId}).
		Scan(&result).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", result.Name, result.Alias), nil
}

// GetSourceFieldId
//
//	@Description: 获取源字段信息
//	@receiver mfr
//	@return string
//	@return error
func (mfr *ModelFieldRelation) GetSourceFieldId() (string, error) {
	var result ModelField
	if err := database.DB.Model(&ModelField{}).
		Where(map[string]any{"id": mfr.SourceFieldId}).
		Scan(&result).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", result.Name, result.Alias), nil
}

// GetTargetFieldId
//
//	@Description: 获取目的字段信息
//	@receiver mfr
//	@return string
//	@return error
func (mfr *ModelFieldRelation) GetTargetFieldId() (string, error) {
	var result ModelField
	if err := database.DB.Model(&ModelField{}).
		Where(map[string]any{"id": mfr.TargetFieldId}).
		Scan(&result).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", result.Name, result.Alias), nil
}
