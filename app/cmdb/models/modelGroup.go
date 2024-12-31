package models

import "gcmdb/pkg/database"

// ModelGroup
// @Description: 模型分组
type ModelGroup struct {
	BaseModel
	Alias       string `gorm:"column:alias;type:string;size:255;unique;not null;comment:别名" json:"alias"`
	Name        string `gorm:"column:name;type:string;size:255;unique;not null;comment:名称" json:"name"`
	Description string `gorm:"column:description;type:text;comment:描述" json:"description"`
	Order       uint   `gorm:"column:order;type:uint;default:0;comment:排序" json:"order"`
}

// TableName
//
//	@Description: 模型分组>表名
//	@receiver ModelGroup
//	@return string
func (ModelGroup) TableName() string {
	return "model_group"
}

// GetModels
//
//	@Description: 获取组对应的模型
//	@receiver mg
//	@return []Model
//	@return error
func (mg *ModelGroup) GetModels() ([]Model, error) {
	models := make([]Model, 0)
	if err := database.DB.Model(&Model{}).
		Where("group_id = ?", mg.ID).
		Order("order asc").
		Scan(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}
