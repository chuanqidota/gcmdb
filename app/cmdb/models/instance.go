package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Instance
// @Description: 实例
type Instance struct {
	gorm.Model
	ModelId    uint           `gorm:"column:model_id;type:uint;not null;comment:模型id" json:"model_id"`
	ModelAlias string         `gorm:"column:model_alias;type:string;size:255;not null;comment:模型别名" json:"model_alias"`
	Data       datatypes.JSON `gorm:"column:data;type:json;comment:数据" json:"data"`
}

// TableName
//
//	@Description: 实例>表
//	@receiver Instance
//	@return string
func (Instance) TableName() string {
	return "instance"
}
