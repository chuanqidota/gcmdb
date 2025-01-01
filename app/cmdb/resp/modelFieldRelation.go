package resp

import "gcmdb/app/cmdb/models"

type ListModelFieldRelation struct {
	models.ModelFieldRelation
	SourceModelDisplay string `json:"source_model_display"`
	TargetModelDisplay string `json:"target_model_display"`
	SourceFieldDisplay string `json:"source_field_display"`
	TargetFieldDisplay string `json:"target_field_display"`
}
