package models

import "gorm.io/gorm"

// ModelField
// @Description:模型字段>模型和模型字段的别名及名称联合唯一
type ModelField struct {
	gorm.Model
	ModelId      uint   `gorm:"column:model_id;type:uint;not null;uniqueIndex:idx_alias_model_id;uniqueIndex:idx_name_model_id;comment:模型id" json:"model_id"`
	FieldGroupId uint   `gorm:"column:model_id;type:uint;comment:字段分组id" json:"field_group_id"`
	Alias        string `gorm:"column:alias;type:string;size:255;not null;uniqueIndex:idx_alias_model_id;comment:别名" json:"alias"`
	Name         string `gorm:"column:name;type:string;size:255;not null;uniqueIndex:idx_name_model_id;comment:名称" json:"name"`
	Type         string `gorm:"column:type;type:string;size:255;not null;comment:类型" json:"type"` //ENUM('string','number','bool')待定
	IsUnique     bool   `gorm:"column:is_unique;type:boolean;not null;default:false;comment:是否唯一" json:"is_unique"`
	IsRequired   bool   `gorm:"column:is_required;type:boolean;not null;default:false;comment:是否必填" json:"is_required"`
	Order        uint   `gorm:"column:order;type:uint;default:0;comment:排序" json:"order"`
}

// TableName
//
//	@Description: 模型字段>表
//	@receiver ModelField
//	@return string
func (ModelField) TableName() string {
	return "model_field"
}
