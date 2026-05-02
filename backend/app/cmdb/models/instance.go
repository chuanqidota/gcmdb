package models

import (
	"gorm.io/datatypes"
)

// Instance
// @Description: 实例
type Instance struct {
	BaseModel
	ModelId uint           `gorm:"column:model_id;type:uint;not null;index;comment:模型id" json:"model_id" binding:"required"`
	Data    datatypes.JSON `gorm:"column:data;type:json;comment:数据" json:"data" binding:"required"`
	Version int            `gorm:"column:version;default:1;comment:乐观锁版本号" json:"version"`
}

// TableName
//
//	@Description: 实例>表
//	@receiver Instance
//	@return string
func (Instance) TableName() string {
	return "instance"
}
