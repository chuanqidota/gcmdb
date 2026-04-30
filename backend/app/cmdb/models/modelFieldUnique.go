package models

// ModelFieldUnique
// @Description: 模型唯一字段
type ModelFieldUnique struct {
	BaseModel
	Fields      string `gorm:"column:fields; type:string; size:255; uniqueIndex:idx_fields_model_id; comment:字段" json:"fields" binding:"required"` // 唯一字段,联合唯一用英文逗号分隔                                                 // 字段分组ID
	ModelId     uint   `gorm:"column:model_id; type:uint; uniqueIndex:idx_fields_model_id; comment:模型" json:"model_id" binding:"required"`         // 对应的模型ID
	Description string `gorm:"column:description;type:text;comment:描述" json:"description"`                                                         // 描述
}

func (ModelFieldUnique) TableName() string {
	return "model_field_unique"
}
