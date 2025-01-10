package resp

import "gcmdb/app/cmdb/models"

type RetrieveModelGroupField struct {
	Group       models.ModelFieldGroup `json:"groups" label:"组"`
	ModelFields []models.ModelField    `json:"fields" label:"模型字段"`
}

type RetrieveModelResp struct {
	models.Model
	Groups []RetrieveModelGroupField `json:"groups" label:"组列表"`
}
