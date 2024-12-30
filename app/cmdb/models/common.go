package models

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// CustomTime 自定义时间类型
type CustomTime time.Time

// MarshalJSON 自定义 JSON 序列化方法
func (t CustomTime) MarshalJSON() ([]byte, error) {
	// 格式化时间为 "2006-01-02 15:04:05"
	formattedTime := time.Time(t).Format("2006-01-02 15:04:05")
	return []byte(fmt.Sprintf(`"%s"`, formattedTime)), nil
}

// Value 实现 driver.Valuer 接口
func (t CustomTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner 接口
func (t *CustomTime) Scan(value interface{}) error {
	var tm time.Time
	if value != nil {
		tm = value.(time.Time)
	}
	*t = CustomTime(tm)
	return nil
}

// BaseModel
// @Description: 基础模型
type BaseModel struct {
	ID        uint           `gorm:"column:id; primary_key; AUTO_INCREMENT" json:"id"`
	CreatedAt CustomTime     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt CustomTime     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
