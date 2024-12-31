package resp

import "gcmdb/app/cmdb/models"

type RetrieveModelGroup struct {
	models.ModelGroup
	Models []models.Model `json:"models"`
}
