package params

import (
	"gorm.io/datatypes"
)

type CreateInstance struct {
	ModelAlias string         `json:"model_alias" label:"模型别名"`
	Data       datatypes.JSON `json:"data" label:"实例数据"`
}

type UpdateInstance struct {
	Id   uint           `json:"id" label:"实例ID"`
	Data datatypes.JSON `json:"data" label:"更新数据"`
}

type DeleteInstance struct {
	ModelAlias string `json:"model_alias" label:"模型英文名"`
	Id         uint   `json:"id" label:"实例ID"`
}

type MulDeleteInstance struct {
	ModelAlias string `json:"model_alias" label:"模型英文名"`
	Ids        []uint `json:"ids" label:"实例ID"`
}

type DirectSearch struct {
	Uuid string `json:"uuid" label:"uuid"`
}
