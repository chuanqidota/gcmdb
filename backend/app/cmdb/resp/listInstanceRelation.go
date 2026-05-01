package resp

type ListInstanceRelation struct {
	ModelId      uint   `json:"model_id" label:"模型ID"`
	ModelDisplay string `json:"model_display" label:"模型展示"`
	Count        int64  `json:"count" label:"数量"`
	Instances    []uint `json:"instances" label:"实例列表"`
}
