package resp

import "gcmdb/app/cmdb/models"

type ListModelFieldRelation struct {
	models.ModelFieldRelation
	SourceModelDisplay string `json:"source_model_display" label:"源模型展示"`
	TargetModelDisplay string `json:"target_model_display" label:"目标模型展示"`
	SourceFieldDisplay string `json:"source_field_display" label:"源字段展示"`
	TargetFieldDisplay string `json:"target_field_display" label:"目的字段展示"`
}
