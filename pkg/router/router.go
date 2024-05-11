package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	router := gin.Default()

	task := router.Group("task")

	{
		fmt.Println(task)

	}

	return router
}
