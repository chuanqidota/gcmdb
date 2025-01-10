package resp

import (
	"gcmdb/app/cmdb/models"
)

type ModelInfo struct {
	Model              models.Model                `json:"model"`
	ModelRelation      []models.ModelRelation      `json:"model_relations"`
	ModelFieldGroup    []models.ModelFieldGroup    `json:"model_field_groups"`
	ModelFieldUnique   []models.ModelFieldUnique   `json:"model_field_uniques"`
	ModelFieldRelation []models.ModelFieldRelation `json:"model_field_relations"`
	ModelField         []models.ModelField         `json:"model_fields"`
}
