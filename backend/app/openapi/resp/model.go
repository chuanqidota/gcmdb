package resp

import (
	"gcmdb/app/cmdb/models"
)

type ModelListItem struct {
	ID          uint   `json:"id"`
	GroupId     uint   `json:"group_id"`
	GroupName   string `json:"group_name"`
	Alias       string `json:"alias"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsUsable    bool   `json:"is_usable"`
	Icon        string `json:"icon"`
	Order       uint   `json:"order"`
}

type ModelRelationInfo struct {
	ID          uint   `json:"id"`
	SourceId    uint   `json:"source_id"`
	SourceName  string `json:"source_name"`
	TargetId    uint   `json:"target_id"`
	TargetName  string `json:"target_name"`
	TypeId      uint   `json:"type_id"`
	TypeName    string `json:"type_name"`
	S2T         string `json:"s2t"`
	T2S         string `json:"t2s"`
	Description string `json:"description"`
}

type ModelInfo struct {
	Model              models.Model                `json:"model"`
	ModelRelation      []ModelRelationInfo         `json:"model_relations"`
	ModelFieldGroup    []models.ModelFieldGroup    `json:"model_field_groups"`
	ModelFieldUnique   []models.ModelFieldUnique   `json:"model_field_uniques"`
	ModelFieldRelation []models.ModelFieldRelation `json:"model_field_relations"`
	ModelField         []models.ModelField         `json:"model_fields"`
}
