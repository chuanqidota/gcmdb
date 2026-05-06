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

// InstanceQuery 实例查询 handler（仅 GET）
// 与 InstanceWrite 分离，避免 isGet 双模式分支带来的代码复杂度
func (i *instance) InstanceQuery(c *gin.Context) {
	action := c.Param("action")

	switch action {

	case "direct": // 直接sql
		uuid := c.Query("uuid")
		if uuid == "" {
			response.Fail(c, "参数uuid不能为空")
			return
		}
		result, err := utils.Instance.DirectSearch(uuid)
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", result)

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

	case "detail": // 查询单个实例详情
		idStr := c.Query("id")
		if idStr == "" {
			response.Fail(c, "参数id不能为空")
			return
		}
		result, err := utils.Instance.DetailInstance(idStr)
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "查询成功", result)

	case "topology": // 查询实例拓扑关系
		idStr := c.Query("id")
		if idStr == "" {
			response.Fail(c, "参数id不能为空")
			return
		}
		modelAlias := c.Query("model")
		result, err := utils.Instance.TopologyInstance(idStr, modelAlias)
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "查询成功", result)

	default:
		response.Fail(c, fmt.Sprintf("路径参数不对,不可以为:%+v", action))
	}
}

// InstanceWrite
//
//	@Description: 实例写入（POST）
//	@receiver i
//	@param c
func (i *instance) InstanceWrite(c *gin.Context) {
	action := c.Param("action")

	switch action {

	case "fulltext": // 全文检索（POST）
		var body params.FulltextInstance
		if err := c.ShouldBindJSON(&body); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if body.Search == "" {
			response.Fail(c, "参数search不能为空")
			return
		}
		if count, result, err := utils.Instance.FulltextInstance(body); err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		} else {
			response.Success(c, "执行成功", map[string]any{
				"count":  count,
				"result": result,
			})
		}

	case "search": // 结构化搜索（POST）
		var body params.SearchInstance
		if err := c.ShouldBindJSON(&body); err != nil {
			response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
			return
		}
		if body.Model == "" {
			response.Fail(c, "参数model不能为空")
			return
		}
		count, results, err := utils.Instance.SearchInstance(body)
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "执行成功", map[string]any{
			"count":   count,
			"results": results,
		})

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

	default:
		response.Fail(c, fmt.Sprintf("路径参数不对,不可以为:%+v", action))
	}
}
