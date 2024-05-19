package models

import (
	"gorm.io/gorm"
)

// Model
// @Description:模型
type Model struct {
	gorm.Model
	GroupId     uint   `gorm:"column:group_id;type:uint;not null;comment:分组id" json:"group_id"`
	Alias       string `gorm:"column:alias;type:string;size:255;unique;not null;comment:别名" json:"alias"`
	Name        string `gorm:"column:name;type:string;size:255;unique;not null;comment:名称" json:"name"`
	Description string `gorm:"column:description;type:text;comment:描述" json:"description"`
	IsUsable    bool   `gorm:"column:is_usable;type:boolean;not null;default:true;comment:是否可用" json:"is_usable"`
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
