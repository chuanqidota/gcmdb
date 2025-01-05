package resp

type ListInstanceRelation struct {
	ModelId      uint   `json:"modelId"`
	ModelDisplay string `json:"modelDisplay"`
	Count        int64  `json:"count"`
	Instances    []uint `json:"instances"`
}
