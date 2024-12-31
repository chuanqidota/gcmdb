package params

import "gcmdb/app/cmdb/models"

type RetrieveModelGroup struct {
	models.ModelGroup
	Models []models.Model `json:"models"`
}

type PatchModelGroupBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
