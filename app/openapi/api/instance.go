package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	cmdbUtils "gcmdb/app/cmdb/utils"
	"gcmdb/app/openapi/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"strconv"
	"strings"

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
	verifyData, err := cmdbUtils.Verify.CreateInstance(modelId, data)
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
	go cmdbUtils.InstanceRelation.CreateInstance(modelId, createInstance.ID)
	response.Success(c, "执行成功", nil)
}

func (i *instance) DeleteInstanceById(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id}).
		Delete(&models.Instance{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	_id, err := strconv.Atoi(id)
	if err != nil {
		response.Fail(c, fmt.Sprintf("参数转换失败-%s", err.Error()))
		return
	}
	go cmdbUtils.InstanceRelation.DeleteInstance(uint(_id))
	response.Success(c, "执行成功", nil)
}

func(i *instance)DeleteInstanceByField(c *gin.Context){
	var body map[string]any
	if err:=c.ShouldBindJSON(&body);err!=nil{
		response.Fail(c,fmt.Sprintf("参数校验失败-%s",err.Error()))
		return
	}
	conditions := make([]string, 0)
	params := make([]any, 0)
	for key, value := range body {
		conditions = append(conditions, fmt.Sprintf("data->'$.%s' = ?", key))
		params = append(params,fmt.Sprintf("%%%s%%",value))
	}
	query := strings.Join(conditions, " AND ")
	if err:=database.DB.Model(&models.Instance{}).
	Where(query,params...).
	Delete(&models.Instance{}).Error;err!=nil{
		response.Fail(c,fmt.Sprintf("删除失败-%s",err.Error()))
		return
	}


}
