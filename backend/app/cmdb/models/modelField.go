package models

import (
	"time"
)

var allowedFieldTypes = map[string]bool{
	"string": true, "number": true, "bool": true,
	"date": true, "datetime": true, "json": true,
}

func DefaultValueByType(fieldType string) any {
	switch fieldType {
	case "string":
		return ""
	case "number":
		return 0
	case "bool":
		return false
	case "date":
		return time.Now().Format(time.DateOnly)
	case "datetime":
		return time.Now().Format(time.DateTime)
	case "json":
		return map[string]any{}
	default:
		return ""
	}
}

func IsValidFieldType(t string) bool {
	return allowedFieldTypes[t]
}

// ModelField
// @Description:模型字段>模型和模型字段的别名及名称联合唯一
type ModelField struct {
	BaseModel
	ModelId      uint   `gorm:"column:model_id;type:uint;not null;uniqueIndex:idx_alias_model_id;uniqueIndex:idx_name_model_id;comment:模型id" json:"model_id" binding:"required"`
	FieldGroupId uint   `gorm:"column:field_group_id;type:uint;comment:字段分组id" json:"field_group_id" binding:"required"`
	Alias        string `gorm:"column:alias;type:string;size:255;not null;uniqueIndex:idx_alias_model_id;comment:别名" json:"alias" binding:"required"`
	Name         string `gorm:"column:name;type:string;size:255;not null;uniqueIndex:idx_name_model_id;comment:名称" json:"name" binding:"required"`
	Type         string `gorm:"column:type;type:string;size:255;not null;comment:类型" json:"type" binding:"required"`
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
