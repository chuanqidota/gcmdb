package models

import (
	"gcmdb/pkg/database"
)

// ModelFieldGroup
// @Description:模型字段分组>模型和字段分组名称联合唯一
type ModelFieldGroup struct {
	BaseModel
	ModelId uint   `gorm:"column:model_id;type:uint;not null;uniqueIndex:idx_name_model_id;comment:模型id" json:"model_id"`
	Name    string `gorm:"column:name;type:string;size:255;not null;uniqueIndex:idx_name_model_id;comment:名称" json:"name"`
	Order   uint   `gorm:"column:order;type:uint;default:0;comment:排序" json:"order"`
}

// TableName
//
//	@Description:模型字段分组>表
//	@receiver ModelFieldGroup
//	@return string
func (ModelFieldGroup) TableName() string {
	return "model_field_group"
}

// Indexes 定义复合唯一索引

// GetModelFields
//
//	@Description: 根据模型字段分组id获取模型字段
//	@receiver mfg
//	@return []ModelField
//	@return error
func (mfg *ModelFieldGroup) GetModelFields() ([]ModelField, error) {
	modelFields := make([]ModelField, 0)
	if err := database.DB.Model(&ModelField{}).
		Where(map[string]any{"field_group_id": mfg.ID}).
		Order("order asc").
		Scan(&modelFields).Error; err != nil {
		return nil, err
	}
	return modelFields, nil
}
