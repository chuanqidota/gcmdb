package resp

import "gcmdb/app/cmdb/models"

type RetrieveModelFieldGroup struct {
	models.ModelFieldGroup
	Fields []models.ModelField `json:"fields"`
}
