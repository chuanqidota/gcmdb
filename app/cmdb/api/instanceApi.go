package api

import "github.com/gin-gonic/gin"

type instance struct {
}

var Instance = new(instance)

func (i *instance) CreateInstance(c *gin.Context) {
}

func (i *instance) ListInstance(c *gin.Context) {

}

func (i *instance) RetrieveInstance(c *gin.Context) {

}

func (i *instance) UpdateInstance(c *gin.Context) {

}

func (i *instance) DeleteInstance(c *gin.Context) {

}
