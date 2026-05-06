package resp

type StatsResponse struct {
	ModelGroupCount     int64                `json:"model_group_count"`
	ModelCount          int64                `json:"model_count"`
	InstanceCount       int                  `json:"instance_count"`
	RelationCount       int64                `json:"relation_count"`
	ModelInstanceCounts []ModelInstanceCount `json:"model_instance_counts"`
	GroupSummaries      []GroupSummary       `json:"group_summaries"`
}

type ModelInstanceCount struct {
	ModelID       uint   `json:"model_id"`
	ModelName     string `json:"model_name"`
	ModelAlias    string `json:"model_alias"`
	GroupName     string `json:"group_name"`
	InstanceCount int    `json:"instance_count"`
}

type GroupSummary struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ModelCount    int    `json:"model_count"`
	InstanceCount int    `json:"instance_count"`
}
