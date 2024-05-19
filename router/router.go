package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	router := gin.Default()

	gcmdb := router.Group("gcmdb")
	v1 := gcmdb.Group("v1")

	{
		fmt.Println(v1)

	}

	return router
}
