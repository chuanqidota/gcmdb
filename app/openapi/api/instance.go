package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/utils"
	"gcmdb/app/openapi/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type instance struct {
}

var Instance = new(instance)

func (i *instance) CreateInstance(c *gin.Context) {
	var body params.CreateInstance
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
		return
	}

	modelAlias := body.ModelAlias
	data := body.Data

	// 根据模型别称获取模型id
	var modelId uint
	if err := database.DB.Model(&models.Model{}).
		Select("id").
		Where(map[string]any{"alias": modelAlias}).
		Scan(&modelId).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询出错-%s", err.Error()))
		return
	}
	// 校验数据
	verifyData, err := utils.Verify.CreateInstance(modelId, data)
	if err != nil {
		response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
		return
	}
	// 创建实例数据
	createInstance := models.Instance{
		ModelId:    modelId,
		ModelAlias: modelAlias,
		Data:       verifyData,
	}

	if err := database.DB.Model(&models.Instance{}).
		Create(&createInstance).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建实例失败-%s", err.Error()))
		return
	}
	// 异步创建实例关系
	go utils.InstanceRelation.CreateInstance(modelId, createInstance.ID)
	response.Success(c, "执行成功", nil)
}


func(i *instance)DeleteInstance(c *gin.Context){
	
}