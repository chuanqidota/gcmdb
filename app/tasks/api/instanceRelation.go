package api

import (
	"fmt"
	"gcmdb/app/cmdb/utils"
	"gcmdb/app/tasks/params"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

// SyncInstanceRelation
//
//	@Description: 同步实例关系
//	@receiver ir
//	@param c
func (ir *instanceRelation) SyncInstanceRelation(c *gin.Context) {
	var body params.SyncInstanceRelation
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	modelId := body.ModelId
	// 同步指定源模型对应的实例关系
	go func() {
		if err := utils.InstanceRelation.SyncSourceModelInstanceRelation(modelId); err != nil {
			fmt.Printf("同步指定模型:%+v,同步实例关系成功", modelId)
		}
	}()
	response.Success(c, "异步执行中", nil)
}
