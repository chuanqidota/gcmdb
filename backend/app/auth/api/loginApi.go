package api

import (
	"fmt"
	"gcmdb/app/auth/models"
	"gcmdb/app/auth/session"
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
	c.SetCookie(session.CookieName(), sid, 86400*7, "/", "", false, true)
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
