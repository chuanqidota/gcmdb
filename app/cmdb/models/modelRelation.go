package models

// ModelRelation
// @Description:模型关系
type ModelRelation struct {
	BaseModel
	SourceId    uint   `gorm:"column:source_id;type:uint;not null;comment:源模型id" json:"source_id"`
	TargetId    uint   `gorm:"column:target_id;type:uint;not null;comment:目标模型id" json:"target_id"`
	TypeId      uint   `gorm:"column:type_id;type:uint;not null;comment:模型关系类型id" json:"type_id"`
	Description string `gorm:"column:description;type:text;comment:描述" json:"description"`
}

// TableName
//
//	@Description: 模型关系>表
//	@receiver ModelRelation
//	@return string
func (ModelRelation) TableName() string {
	return "model_relation"
}
