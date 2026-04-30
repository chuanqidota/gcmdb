package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": 1,
		"data": data,
	})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":  msg,
		"code": 0,
		"data": nil,
	})
}

func FailWithStatus(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{
		"msg":  msg,
		"code": 0,
		"data": nil,
	})
}
