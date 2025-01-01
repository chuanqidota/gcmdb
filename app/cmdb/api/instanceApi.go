package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
)

type instance struct {
}

var Instance = new(instance)

// CreateInstance
//
//	@Description: 没有做唯一性判断的常见
//	@receiver i
//	@param c
func (i *instance) CreateInstance(c *gin.Context) {
	var body models.Instance
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"model_id":    body.ModelId,
		"model_alias": body.ModelAlias,
		"data":        body.Data,
	}
	if err := database.DB.Model(&models.Instance{}).
		Create(data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建实例失败-%s", err.Error()))
		return
	}
	response.Success(c, "创建成功", nil)

}

func (i *instance) ListInstance(c *gin.Context) {

}

// RetrieveInstance
//
//	@Description: 查询实例详情
//	@receiver i
//	@param c
func (i *instance) RetrieveInstance(c *gin.Context) {
	id := c.Param("id")
	var instance models.Instance
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id}).
		Scan(&instance).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "查询成功", instance)
}

func (i *instance) UpdateInstance(c *gin.Context) {

}

func (i *instance) DeleteInstance(c *gin.Context) {
	id := c.Param("id")
	// 删除实例关系
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{"source_id": id}).
		Or(map[string]any{"target_id": id}).
		Delete(&models.InstanceRelation{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除实例关系失败-%s", err.Error()))
		return
	}
	// 删除实例
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id}).
		Delete(&models.Instance{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除实例失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
