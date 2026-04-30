package models

import (
	"fmt"
	"gcmdb/pkg/database"
)

// InstanceRelation
// @Description: 实例关系
type InstanceRelation struct {
	BaseModel
	SourceModelId    uint `gorm:"column:source_model_id;type:uint;not null;comment:源模型id" json:"source_model_id" binding:"required"`
	TargetModelId    uint `gorm:"column:target_model_id;type:uint;not null;comment:目标模型id" json:"target_model_id" binding:"required"`
	SourceInstanceId uint `gorm:"column:source_id;type:uint;not null;comment:源实例id" json:"source_instance_id" binding:"required"`
	TargetInstanceId uint `gorm:"column:target_id;type:uint;not null;comment:目标实例id" json:"target_instance_id" binding:"required"`
}

// TableName
//
//	@Description: 实例关系>表
//	@receiver InstanceRelation
//	@return string
func (InstanceRelation) TableName() string {
	return "instance_relation"
}

// GetSourceModel
//
//	@Description: 获取源模型
//	@receiver ir
//	@return string
//	@return error
func (ir *InstanceRelation) GetSourceModel() (string, error) {
	var model Model
	if err := database.DB.Model(&Model{}).
		Where(map[string]any{"id": ir.SourceModelId}).
		First(&model).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", model.Name, model.Alias), nil
}

// GetTargetModel
//
//	@Description: 获取目标模型
//	@receiver ir
//	@return string
//	@return error
func (ir *InstanceRelation) GetTargetModel() (string, error) {
	var model Model
	if err := database.DB.Model(&Model{}).
		Where(map[string]any{"id": ir.TargetModelId}).
		First(&model).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", model.Name, model.Alias), nil
}
