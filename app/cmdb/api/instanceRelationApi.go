package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
)

type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

// CreateInstanceRelation
//
//	@Description: 创建实例关系
//	@receiver i
//	@param c
func (i *instanceRelation) CreateInstanceRelation(c *gin.Context) {
	var body models.InstanceRelation
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
		return
	}
	var result models.InstanceRelation
	if err := database.DB.Model(&models.InstanceRelation{}).
		FirstOrCreate(&result, body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "创建成功", nil)
}

func (i *instanceRelation) ListInstanceRelation(c *gin.Context) {
}

// DeleteInstanceRelation
//
//	@Description: 删除实例关系
//	@receiver i
//	@param c
func (i *instanceRelation) DeleteInstanceRelation(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{"id": id}).
		Delete(&models.InstanceRelation{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
