package models

import "time"

// AuditLog
// @Description: 审计日志
type AuditLog struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	ResourceType string   `gorm:"column:resource_type;size:50;not null;index;comment:资源类型" json:"resource_type"`
	ResourceID  string    `gorm:"column:resource_id;size:50;index;comment:资源ID" json:"resource_id"`
	Action      string    `gorm:"column:action;size:20;not null;comment:操作类型" json:"action"`
	Method      string    `gorm:"column:method;size:10;not null;comment:HTTP方法" json:"method"`
	Path        string    `gorm:"column:path;size:500;not null;comment:请求路径" json:"path"`
	RequestBody string    `gorm:"column:request_body;type:text;comment:请求体" json:"request_body"`
	SourceIP    string    `gorm:"column:source_ip;size:50;comment:来源IP" json:"source_ip"`
	UserAgent   string    `gorm:"column:user_agent;size:500;comment:User-Agent" json:"user_agent"`
	StatusCode  int       `gorm:"column:status_code;comment:响应状态码" json:"status_code"`
	CreatedAt   time.Time `gorm:"column:created_at;index;comment:创建时间" json:"created_at"`
}

// TableName
//
//	@Description: 审计日志>表
//	@receiver AuditLog
//	@return string
func (AuditLog) TableName() string {
	return "audit_log"
}
