package models

// ModelFieldGroup
// @Description:模型字段分组>模型和字段分组名称联合唯一
type ModelFieldGroup struct {
	BaseModel
	ModelId uint   `gorm:"column:model_id;type:uint;not null;uniqueIndex:idx_name_model_id;comment:模型id" json:"model_id"`
	Name    string `gorm:"column:name;type:string;size:255;unique;not null;uniqueIndex:idx_name_model_id;comment:名称" json:"name"`
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
