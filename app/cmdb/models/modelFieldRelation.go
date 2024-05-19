package models

import "gorm.io/gorm"

// ModelFieldRelation
// @Description: 模型字段关联
type ModelFieldRelation struct {
	gorm.Model
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
