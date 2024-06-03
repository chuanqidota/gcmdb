package models

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel
// @Description: 基础模型
type BaseModel struct {
	ID        uint           `gorm:"column:id; primary_key; AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
