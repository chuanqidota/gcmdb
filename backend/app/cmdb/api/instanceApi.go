package api

import (
	"errors"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/app/cmdb/utils"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	"gcmdb/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	// 校验模型存在
	var modelCount int64
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": body.ModelId}).Count(&modelCount)
	if modelCount == 0 {
		response.Fail(c, "关联模型不存在")
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
		ModelId: body.ModelId,
		Data:    verifyData,
	}
	if err := database.DB.Model(&models.Instance{}).
		Create(&instance).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建实例失败-%s", err.Error()))
		return
	}
	// 异步创建实例关联
	go func() {
		if err := utils.InstanceRelation.CreateInstance(body.ModelId, instance.ID); err != nil {
			logger.Error(fmt.Sprintf("创建实例关联失败-%s", err.Error()))
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
		First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, "实例不存在")
			return
		}
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "查询成功", instance)
}

// UpdateInstance
//
//	@Description: 更新实例（乐观锁）
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
	// 乐观锁更新：WHERE 加 version 条件，更新后 version+1
	version := body.Version
	result := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id, "version": version}).
		Updates(map[string]any{
			"updated_at": time.Now(),
			"data":       verifyData,
			"version":    gorm.Expr("version + 1"),
		})
	if result.Error != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", result.Error.Error()))
		return
	}
	if result.RowsAffected == 0 {
		response.FailWithStatus(c, 409, "数据已被其他人修改，请刷新后重试")
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
	// 删除实例关系和实例（事务）
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Model(&models.InstanceRelation{}).
			Where(map[string]any{"source_id": id}).
			Or(map[string]any{"target_id": id}).
			Delete(&models.InstanceRelation{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Model(&models.Instance{}).
			Where(map[string]any{"id": id}).
			Delete(&models.Instance{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
