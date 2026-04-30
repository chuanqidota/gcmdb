package models

import (
	"gcmdb/pkg/database"
)

// Model
// @Description:模型
type Model struct {
	BaseModel
	GroupId     uint   `gorm:"column:group_id;type:uint;not null;comment:分组id" json:"group_id" binding:"required"`
	Alias       string `gorm:"column:alias;type:string;size:255;unique;not null;comment:别名" json:"alias" binding:"required"`
	Name        string `gorm:"column:name;type:string;size:255;unique;not null;comment:名称" json:"name" binding:"required"`
	Description string `gorm:"column:description;type:text;comment:描述" json:"description"`
	IsUsable    bool   `gorm:"column:is_usable;type:boolean;not null;default:true;comment:是否可用" json:"is_usable" binding:"required"`
	Icon        string `gorm:"column:icon;type:string;size:255;comment:图标" json:"icon"`
	Order       uint   `gorm:"column:order;type:uint;default:0;comment:排序" json:"order"`
}

// TableName
//
//	@Description:模型>表
//	@receiver Model
//	@return string
func (Model) TableName() string {
	return "model"
}

// GetModelFieldGroups
//
//	@Description: 获取模型对应的字段分组
//	@receiver m
//	@return []ModelFieldGroup
//	@return error
func (m *Model) GetModelFieldGroups() ([]ModelFieldGroup, error) {
	modelFieldGroups := make([]ModelFieldGroup, 0)
	if err := database.DB.Model(&ModelFieldGroup{}).
		Where(map[string]any{"model_id": m.ID}).
		Order("order asc").
		Scan(&modelFieldGroups).Error; err != nil {
		return nil, err
	}
	return modelFieldGroups, nil
}

// GetModelFields
//
//	@Description: 获取模型对应的字段
//	@receiver m
//	@return []ModelField
//	@return error
func (m *Model) GetModelFields() ([]ModelField, error) {
	modelFields := make([]ModelField, 0)
	if err := database.DB.Model(&ModelField{}).
		Where(map[string]any{"model_id": m.ID}).
		Order("order asc").
		Scan(&modelFields).Error; err != nil {
		return nil, err
	}
	return modelFields, nil
}
