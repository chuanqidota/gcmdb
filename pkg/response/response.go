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
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": 0,
		"data": nil,
	})
}
