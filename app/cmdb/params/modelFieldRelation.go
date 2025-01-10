package params

type UpdateModelFieldRelation struct {
	TargetModelId uint `json:"target_model_id" label:"目标模型ID"`
	SourceFieldId uint `json:"source_field_id" label:"源字段ID"`
	TargetFieldId uint `json:"target_field_id" label:"目的字段ID"`
}
