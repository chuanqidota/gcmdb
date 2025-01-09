package api

import (
	"fmt"
	"gcmdb/pkg/response"

	"gcmdb/app/openapi/params"
	"gcmdb/app/openapi/utils"

	"github.com/gin-gonic/gin"
)

type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

func (ir *instanceRelation) InstanceRelationAction(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "create":
		var createBody params.CreateInstanceRelation
		if err := c.ShouldBindJSON(&createBody); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if err := utils.InstanceRelation.CreateInstanceRelation(createBody); err != nil {
			response.Fail(c, fmt.Sprintf("创建关联失败-%s", err.Error()))
			return
		}
	case "delete":
		var deleteBody params.DeleteInstanceRelation
		if err := c.ShouldBindJSON(&deleteBody); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if err:=utils.InstanceRelation.DeleteInstanceRelation(deleteBody);err!=nil{
			response.Fail(c,fmt.Sprintf("删除关联失败-%s",err.Error()))
			return
		}

	default:
		response.Fail(c, fmt.Sprintf("参数错误,action不在create和delete范围内:%s", action))
	}
	response.Success(c,"执行成功",nil)
}
