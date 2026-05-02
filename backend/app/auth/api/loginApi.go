package api

import (
	"fmt"
	"gcmdb/app/auth/models"
	"gcmdb/app/auth/session"
	"gcmdb/config"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct{}

var Login = new(login)

type loginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录
func (l *login) Login(c *gin.Context) {
	var body loginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 查询用户
	var user models.User
	if err := database.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		response.FailWithStatus(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	if !user.IsActive {
		response.FailWithStatus(c, http.StatusUnauthorized, "用户已禁用")
		return
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		response.FailWithStatus(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	// 创建 session
	sid := session.Create(user.ID, user.Username)
	c.SetCookie(session.CookieName(), sid, config.Conf.Server.SessionMaxAge, "/", "", false, true)
	response.Success(c, "登录成功", gin.H{
		"user_id":  user.ID,
		"username": user.Username,
	})
}

// Logout 登出
func (l *login) Logout(c *gin.Context) {
	sid, _ := c.Cookie(session.CookieName())
	if sid != "" {
		session.Delete(sid)
	}
	c.SetCookie(session.CookieName(), "", -1, "/", "", false, true)
	response.Success(c, "登出成功", nil)
}

// ChangePassword 修改密码
func (l *login) ChangePassword(c *gin.Context) {
	var body struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.Fail(c, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.OldPassword)); err != nil {
		response.Fail(c, "当前密码错误")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, "密码加密失败")
		return
	}

	if err := database.DB.Model(&user).Update("password_hash", string(hash)).Error; err != nil {
		response.Fail(c, "修改失败")
		return
	}

	response.Success(c, "密码修改成功", nil)
}

// ResetPassword 管理员重置密码
func (l *login) ResetPassword(c *gin.Context) {
	var body struct {
		UserID      uint   `json:"user_id" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}

	username, _ := c.Get("username")
	if username != "admin" {
		response.FailWithStatus(c, http.StatusForbidden, "仅管理员可重置密码")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", body.UserID).First(&user).Error; err != nil {
		response.Fail(c, "用户不存在")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, "密码加密失败")
		return
	}

	if err := database.DB.Model(&user).Update("password_hash", string(hash)).Error; err != nil {
		response.Fail(c, "重置失败")
		return
	}

	response.Success(c, fmt.Sprintf("用户 %s 密码已重置", user.Username), nil)
}

// Me 当前用户信息
func (l *login) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.Fail(c, "用户不存在")
		return
	}
	response.Success(c, "查询成功", gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"token":     user.Token,
		"is_active": user.IsActive,
		"created_at": user.CreatedAt.Format(time.DateTime),
		"_username":  username,
	})
}
