package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
)

type model struct {
}

var Model = new(model)

// GetAllModels
//
//	@Description: 获取所有模型的信息
//	@receiver m
//	@param c
func (m *model) GetAllModels(c *gin.Context) {
	_models := make([]models.Model, 0)
	if err := database.DB.Model(&models.Model{}).
		Scan(&_models).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", _models)
}
