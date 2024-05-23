package params

import "gcmdb/app/cmdb/models"

type CommonQuery struct {
	Search string `form:"search"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

type RetrieveModelGroup struct {
	models.ModelGroup
	Models []models.Model `json:"models"`
}
