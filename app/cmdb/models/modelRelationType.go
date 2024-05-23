package models

// ModelRelationType
// @Description: 模型关系类型
type ModelRelationType struct {
	BaseModel
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
