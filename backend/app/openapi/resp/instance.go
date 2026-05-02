package resp

import "gorm.io/datatypes"

type FulltextInstance struct {
	ID        uint           `json:"id" label:"实例ID"`
	ModelId   uint           `json:"model_id" label:"模型ID"`
	ModelName string         `json:"model_name" label:"模型名称"`
	Data      datatypes.JSON `json:"data" label:"实例数据"`
}

type DetailInstance struct {
	ID        uint           `json:"id"`
	ModelId   uint           `json:"model_id"`
	ModelName string         `json:"model_name"`
	ModelAlias string        `json:"model_alias"`
	Data      datatypes.JSON `json:"data"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

type TopologyRelation struct {
	InstanceId   uint           `json:"instance_id"`
	ModelId      uint           `json:"model_id"`
	ModelName    string         `json:"model_name"`
	ModelAlias   string         `json:"model_alias"`
	RelationType string         `json:"relation_type"`
	Data         datatypes.JSON `json:"data"`
}

type TopologyInfo struct {
	Instance DetailInstance      `json:"instance"`
	Upstream []TopologyRelation  `json:"upstream"`
	Downstream []TopologyRelation `json:"downstream"`
}
