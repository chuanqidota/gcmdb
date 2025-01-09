package params

type CreateInstanceRelation struct {
	SourceModel       string `json:"source_model" binding:"required" label:"源模型英文名"`
	TargetModel       string `json:"target_model" binding:"required" label:"目标模型英文名"`
	SourceInstanceId  uint   `json:"source_id" binding:"required" label:"源实例id"`
	TargetInstanceIds []uint `json:"target_ids" binding:"required" label:"目标实例数组"`
}

type DeleteInstanceRelation struct {
	SourceModel      string `json:"source_model" binding:"required" label:"源模型英文名"`
	TargetModel      string `json:"target_model" binding:"required" label:"目标模型英文名"`
	SourceInstanceId uint   `json:"source_id" binding:"required" label:"源实例id"`
	TargetInstanceId uint   `json:"target_id" binding:"required" label:"目标实例id"`
}
