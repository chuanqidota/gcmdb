package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/app/cmdb/utils"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"strconv"
	"time"

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
	// 验证参数
	verifyData, err := utils.Verify.VerifyCreateInstance(body.ModelId, body.Data)
	if err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 创建实例
	instance := models.Instance{
		ModelId:    body.ModelId,
		ModelAlias: body.ModelAlias,
		Data:       verifyData,
	}
	if err := database.DB.Model(&models.Instance{}).
		Create(&instance).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建实例失败-%s", err.Error()))
		return
	}
	// 异步创建实例关联
	go func() {
		if err := utils.InstanceRelation.CreateInstance(body.ModelId, instance.ID); err != nil {
			fmt.Printf("创建实例关联失败-%s", err.Error())
		}
	}()
	response.Success(c, "创建成功", nil)

}

// ListInstance
//
//	@Description: 查询实例
//	@receiver i
//	@param c
func (i *instance) ListInstance(c *gin.Context) {
	modelId := c.Param("model_id")
	var query params.ListInstance
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	_modelId, err := strconv.Atoi(modelId)
	if err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	resp, err := utils.SimpleSearchInstance(uint(_modelId), query)
	if err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "查询成功", resp)

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

// UpdateInstance
//
//	@Description: 更新实例
//	@receiver i
//	@param c
func (i *instance) UpdateInstance(c *gin.Context) {
	id := c.Param("id")
	var body params.UpdateInstance
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 参数转化
	_id, err := strconv.Atoi(id)
	if err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 验证参数
	verifyData, err := utils.Verify.VerifyUpdateInstance(uint(_id), body.Data)
	if err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	// 更新实例
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id}).
		Updates(map[string]any{
			"updated_at": time.Now(),
			"data":       verifyData,
		}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)

}

// DeleteInstance
//
//	@Description: 删除实例
//	@receiver i
//	@param c
func (i *instance) DeleteInstance(c *gin.Context) {
	id := c.Param("id")
	_id, err := strconv.Atoi(id)
	if err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
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
	_ = _id
}
