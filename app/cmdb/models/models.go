package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ModelGroup
// @Description: 模型分组
type ModelGroup struct {
	gorm.Model
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

// ModelFieldGroup
// @Description:模型字段分组>模型和字段分组名称联合唯一
type ModelFieldGroup struct {
	gorm.Model
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

// ModelField
// @Description:模型字段>模型和模型字段的别名及名称联合唯一
type ModelField struct {
	gorm.Model
	ModelId      uint   `gorm:"column:model_id;type:uint;not null;uniqueIndex:idx_alias_model_id;uniqueIndex:idx_name_model_id;comment:模型id" json:"model_id"`
	FieldGroupId uint   `gorm:"column:model_id;type:uint;comment:字段分组id" json:"field_group_id"`
	Alias        string `gorm:"column:alias;type:string;size:255;not null;uniqueIndex:idx_alias_model_id;comment:别名" json:"alias"`
	Name         string `gorm:"column:name;type:string;size:255;not null;uniqueIndex:idx_name_model_id;comment:名称" json:"name"`
	Type         string `gorm:"column:type;type:string;size:255;not null;comment:类型" json:"type"` //ENUM('string','number','bool')待定
	IsUnique     bool   `gorm:"column:is_unique;type:boolean;not null;default:false;comment:是否唯一" json:"is_unique"`
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

// ModelRelation
// @Description:模型关系
type ModelRelation struct {
	gorm.Model
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

// ModelRelationType
// @Description: 模型关系类型
type ModelRelationType struct {
	gorm.Model
	Name string `gorm:"column:name;type:string;size:255;unique;not null;comment:名称" json:"name"`
	S2T  string `gorm:"column:s2t;type:string;size:255;not null;comment:源到目标的描述" json:"s2t"`
	T2S  string `gorm:"column:t2s;type:string;size:255;not null;comment:目标到源的描述" json:"t2s"`
}

// TableName
//
//	@Description: 模型关系类型>表
//	@receiver ModelRelationType
//	@return string
func (ModelRelationType) TableName() string {
	return "model_relation_type"
}

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

// InstanceRelation
// @Description: 实例关系
type InstanceRelation struct {
	gorm.Model
	SourceModelId    uint `gorm:"column:source_model_id;type:uint;not null;comment:源模型id" json:"source_model_id"`
	TargetModelId    uint `gorm:"column:target_model_id;type:uint;not null;comment:目标模型id" json:"target_model_id"`
	SourceInstanceId uint `gorm:"column:source_id;type:uint;not null;comment:源实例id" json:"source_instance_id"`
	TargetInstanceId uint `gorm:"column:target_id;type:uint;not null;comment:目标实例id" json:"target_instance_id"`
}

// TableName
//
//	@Description: 实例关系>表
//	@receiver InstanceRelation
//	@return string
func (InstanceRelation) TableName() string {
	return "instance_relation"
}
