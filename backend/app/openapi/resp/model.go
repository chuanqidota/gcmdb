package resp

import (
	"gcmdb/app/cmdb/models"
)

type ModelInfo struct {
	Model              models.Model                `json:"model" label:"模型"`
	ModelRelation      []models.ModelRelation      `json:"model_relations" label:"模型关联"`
	ModelFieldGroup    []models.ModelFieldGroup    `json:"model_field_groups" label:"模型字段分组"`
	ModelFieldUnique   []models.ModelFieldUnique   `json:"model_field_uniques" label:"模型字段唯一约束"`
	ModelFieldRelation []models.ModelFieldRelation `json:"model_field_relations" label:"模型字段关联"`
	ModelField         []models.ModelField         `json:"model_fields" label:"模型字段"`
}
