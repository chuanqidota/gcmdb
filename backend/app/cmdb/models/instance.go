package models

import (
	"gorm.io/datatypes"
)

// Instance
// @Description: 实例
type Instance struct {
	BaseModel
	ModelId    uint           `gorm:"column:model_id;type:uint;not null;comment:模型id" json:"model_id" binding:"required"`
	ModelAlias string         `gorm:"column:model_alias;type:string;size:255;not null;comment:模型别名" json:"model_alias" binding:"required"`
	Data       datatypes.JSON `gorm:"column:data;type:json;comment:数据" json:"data" binding:"required"`
}

// TableName
//
//	@Description: 实例>表
//	@receiver Instance
//	@return string
func (Instance) TableName() string {
	return "instance"
}
