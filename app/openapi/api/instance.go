package api

import (
	"fmt"
	"gcmdb/app/openapi/params"
	"gcmdb/app/openapi/utils"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type instance struct {
}

var Instance = new(instance)

// InstanceAction
//
//	@Description: 实例的创建，更新，删除，批量删除
//	@receiver i
//	@param c
func (i *instance) InstanceAction(c *gin.Context) {
	action := c.Param("action")
	switch action {
	case "create": // 创建
		var createBody params.CreateInstance
		if err := c.ShouldBindJSON(&createBody); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if err := utils.Instance.CreateInstance(createBody.ModelAlias, createBody.Data); err != nil {
			response.Fail(c, fmt.Sprintf("创建实例失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", nil)

	case "update": // 更新
		var updateBody params.UpdateInstance
		if err := c.ShouldBindJSON(&updateBody); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if err := utils.Instance.UpdateInstance(updateBody.Id, updateBody.Data); err != nil {
			response.Fail(c, fmt.Sprintf("更新实例失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", nil)

	case "delete": // 删除
		var deleteBody params.DeleteInstance
		if err := c.ShouldBindJSON(&deleteBody); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if err := utils.Instance.DeleteInstance(deleteBody.ModelAlias, deleteBody.Id); err != nil {
			response.Fail(c, fmt.Sprintf("删除实例失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", nil)

	case "mul_delete": // 批量删除
		var mulDeleteBody params.MulDeleteInstance
		if err := c.ShouldBindJSON(&mulDeleteBody); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if err := utils.Instance.MulDeleteInstance(mulDeleteBody.ModelAlias, mulDeleteBody.Ids); err != nil {
			response.Fail(c, fmt.Sprintf("删除实例失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", nil)

	case "direct": // 直接sql
		var uuid params.DirectSearch
		if err := c.ShouldBindJSON(&uuid); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		result, err := utils.Instance.DirectSearch(uuid.Uuid)
		if err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", result)

	case "fulltext": // 全文检索
		var body params.FulltextInstance
		if err := c.ShouldBindJSON(&body); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if count, result, err := utils.Instance.FulltextInstance(body); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		} else {
			response.Success(c, "执行成功", map[string]any{
				"count":  count,
				"result": result,
			})
		}

	case "search": // 搜索
		var body params.SearchInstance
		if err := c.ShouldBindJSON(&body); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}

	case "target": // 通过源查询目标模型
		var body params.SourceTargetInstance
		if err := c.ShouldBindQuery(&body); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		result, err := utils.Instance.TargetInstance(body.Id, body.Model)
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", result)

	case "source": // 通过目标去查源模型
		var body params.SourceTargetInstance
		if err := c.ShouldBindQuery(&body); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		result, err := utils.Instance.SourceInstance(body.Id, body.Model)
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", result)

	default:
		response.Fail(c, fmt.Sprintf("路径参数不对,不可以为:%+v", action))
	}

}
