package resp

type ListInstanceRelation struct {
	ModelId      uint   `json:"modelId" label:"模型ID"`
	ModelDisplay string `json:"modelDisplay" label:"模型展示"`
	Count        int64  `json:"count" label:"数量"`
	Instances    []uint `json:"instances" label:"实例列表"`
}
