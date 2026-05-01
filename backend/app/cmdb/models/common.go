package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// CustomTime 自定义时间类型
type CustomTime time.Time

// MarshalJSON 自定义 JSON 序列化方法
func (t CustomTime) MarshalJSON() ([]byte, error) {
	// 格式化时间为 "2006-01-02 15:04:05"
	formattedTime := time.Time(t).Format(time.DateTime)
	return []byte(fmt.Sprintf(`"%s"`, formattedTime)), nil
}

// UnmarshalJSON 自定义 JSON 反序列化方法
func (t *CustomTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsedTime, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*t = CustomTime(parsedTime)
	return nil
}

// Value 实现 driver.Valuer 接口
func (t CustomTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner 接口
func (t *CustomTime) Scan(value interface{}) error {
	var tm time.Time
	if value != nil {
		switch v := value.(type) {
		case time.Time:
			tm = v
		case []byte:
			parsed, err := time.Parse(time.DateTime, string(v))
			if err != nil {
				return err
			}
			tm = parsed
		case string:
			parsed, err := time.Parse(time.DateTime, v)
			if err != nil {
				return err
			}
			tm = parsed
		default:
			return fmt.Errorf("cannot scan %T into CustomTime", value)
		}
	}
	*t = CustomTime(tm)
	return nil
}

// BaseModel
// @Description: 基础模型
type BaseModel struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	CreatedAt CustomTime     `gorm:"column:created_at; comment:创建时间" json:"created_at"`
	UpdatedAt CustomTime     `gorm:"column:updated_at; comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index; comment:删除时间" json:"-"`
}
