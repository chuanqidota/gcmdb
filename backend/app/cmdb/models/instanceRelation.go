package models

// InstanceRelation
// @Description: 实例关系
type InstanceRelation struct {
	BaseModel
	SourceModelId    uint `gorm:"column:source_model_id;type:uint;not null;uniqueIndex:idx_instance_relation_unique;comment:源模型id" json:"source_model_id" binding:"required"`
	TargetModelId    uint `gorm:"column:target_model_id;type:uint;not null;uniqueIndex:idx_instance_relation_unique;comment:目标模型id" json:"target_model_id" binding:"required"`
	SourceInstanceId uint `gorm:"column:source_id;type:uint;not null;uniqueIndex:idx_instance_relation_unique;index;comment:源实例id" json:"source_instance_id" binding:"required"`
	TargetInstanceId uint `gorm:"column:target_id;type:uint;not null;uniqueIndex:idx_instance_relation_unique;index;comment:目标实例id" json:"target_instance_id" binding:"required"`
}

// TableName
//
//	@Description: 实例关系>表
//	@receiver InstanceRelation
//	@return string
func (InstanceRelation) TableName() string {
	return "instance_relation"
}
