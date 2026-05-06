package resp

import "time"

type StatsResponse struct {
	ModelGroupCount     int64                 `json:"model_group_count"`
	ModelCount          int64                 `json:"model_count"`
	InstanceCount       int                   `json:"instance_count"`
	RelationCount       int64                 `json:"relation_count"`
	ModelInstanceCounts []ModelInstanceCount  `json:"model_instance_counts"`
	RecentLogs          []AuditLogBrief       `json:"recent_logs"`
}

type ModelInstanceCount struct {
	ModelID       uint   `json:"model_id"`
	ModelName     string `json:"model_name"`
	ModelAlias    string `json:"model_alias"`
	GroupName     string `json:"group_name"`
	InstanceCount int    `json:"instance_count"`
}

type AuditLogBrief struct {
	ID           uint      `json:"id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Path         string    `json:"path"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
}
