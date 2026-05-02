package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement;comment:主键ID" json:"id"`
	Username     string    `gorm:"column:username;uniqueIndex;size:64;not null;comment:用户名" json:"username"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null;comment:密码哈希" json:"-"`
	Token        string    `gorm:"column:token;uniqueIndex;size:64;not null;comment:API Token" json:"token"`
	IsActive     bool      `gorm:"column:is_active;default:true;comment:是否启用" json:"is_active"`
	IsAdmin      bool      `gorm:"column:is_admin;default:false;comment:是否管理员" json:"is_admin"`
	CreatedAt    time.Time `gorm:"column:created_at;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;comment:更新时间" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
