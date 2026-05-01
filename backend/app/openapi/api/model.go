package api

import (
	"fmt"
	"gcmdb/app/openapi/utils"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type model struct {
}

var Model = new(model)

// ModelRange
//
//	@Description: 查询所有模型信息和指定模型信息
//	@receiver m
//	@param c
func (m *model) ModelRange(c *gin.Context) {
	_range := c.Param("range")
	switch _range {

	case "all": // 查询所有模型信息
		modelAll, err := utils.Model.ModelAll()
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "查询成功", modelAll)

	case "single": // 查询指定模型信息
		id := c.Query("id")
		modelSingle, err := utils.Model.ModelSingle(id)
		if err != nil {
			response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
			return
		}
		response.Success(c, "查询成功", modelSingle)

	default:
		response.Fail(c, fmt.Sprintf("参数:%+v不在[all,single]范围内", _range))
	}

}
