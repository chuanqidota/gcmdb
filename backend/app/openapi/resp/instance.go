package resp

import "gorm.io/datatypes"

// InstanceItem 实例统一响应结构
type InstanceItem struct {
	ID         uint           `json:"id"`
	ModelId    uint           `json:"model_id"`
	ModelName  string         `json:"model_name"`
	ModelAlias string         `json:"model_alias"`
	Data       datatypes.JSON `json:"data"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}

// TopologyRelation 拓扑关系中的关联实例
type TopologyRelation struct {
	InstanceId   uint           `json:"instance_id"`
	ModelId      uint           `json:"model_id"`
	ModelName    string         `json:"model_name"`
	ModelAlias   string         `json:"model_alias"`
	RelationType string         `json:"relation_type"`
	Data         datatypes.JSON `json:"data"`
}

// TopologyInfo 实例拓扑信息（自身 + 上下游）
type TopologyInfo struct {
	Instance   InstanceItem        `json:"instance"`
	Upstream   []TopologyRelation  `json:"upstream"`
	Downstream []TopologyRelation  `json:"downstream"`
}
