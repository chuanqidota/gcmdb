package api

import (
	"fmt"
	"gcmdb/app/auth/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct{}

var User = new(user)

// ListUsers 管理员查看用户列表
func (u *user) ListUsers(c *gin.Context) {
	username, _ := c.Get("username")
	if username != "admin" {
		response.FailWithStatus(c, http.StatusForbidden, "仅管理员可操作")
		return
	}

	var users []models.User
	if err := database.DB.Order("id asc").Find(&users).Error; err != nil {
		response.Fail(c, "查询失败")
		return
	}

	response.Success(c, "查询成功", users)
}

// CreateUser 管理员创建用户
func (u *user) CreateUser(c *gin.Context) {
	username, _ := c.Get("username")
	if username != "admin" {
		response.FailWithStatus(c, http.StatusForbidden, "仅管理员可操作")
		return
	}

	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, "密码加密失败")
		return
	}

	user := models.User{
		Username:     body.Username,
		PasswordHash: string(hash),
		Token:        uuid.New().String(),
		IsActive:     true,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		response.Fail(c, "创建失败，用户名可能已存在")
		return
	}

	response.Success(c, "创建成功", gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"token":      user.Token,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
	})
}

// PatchUser 管理员更新用户状态
func (u *user) PatchUser(c *gin.Context) {
	username, _ := c.Get("username")
	if username != "admin" {
		response.FailWithStatus(c, http.StatusForbidden, "仅管理员可操作")
		return
	}

	id := c.Param("id")
	var body struct {
		IsActive *bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}

	updates := map[string]any{}
	if body.IsActive != nil {
		updates["is_active"] = *body.IsActive
	}
	if len(updates) == 0 {
		response.Fail(c, "无更新内容")
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Fail(c, "更新失败")
		return
	}

	response.Success(c, "更新成功", nil)
}
