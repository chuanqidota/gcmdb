package api

import "github.com/gin-gonic/gin"

type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

func (i *instanceRelation) CreateInstanceRelation(c *gin.Context) {
}

func (i *instanceRelation) ListInstanceRelation(c *gin.Context) {
}

func (i *instanceRelation) UpdateInstanceRelation(c *gin.Context) {
}

func (i *instanceRelation) DeleteInstanceRelation(c *gin.Context) {
}
