package params

type UpdateModelFieldRelation struct {
	TargetModelId uint `json:"target_model_id"`
	SourceFieldId uint `json:"source_field_id"`
	TargetFieldId uint `json:"target_field_id"`
}
