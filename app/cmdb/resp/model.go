package resp

import "gcmdb/app/cmdb/models"

type RetrieveModelGroupField struct {
	Groups      models.ModelFieldGroup
	ModelFields []models.ModelField `json:"fields"`
}

type RetrieveModelResp struct {
	models.Model
	Groups []RetrieveModelGroupField `json:"groups"`
}
