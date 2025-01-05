package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
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

// SourceTargetInstanceRelation
//
//	@Description: 源到目标的实例关系
//	@receiver i
//	@param c
func (i *instanceRelation) SourceTargetInstanceRelation(c *gin.Context) {
	sourceId := c.Param("source_id")
	modelInstances := make([]models.InstanceRelation, 0)
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{"source_id": sourceId}).
		Scan(&modelInstances).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 匿名函数，模型ID对应的展示名称
	modelIdToModelDisplay := func() map[uint]string {
		var _models []models.Model
		if err := database.DB.Model(&models.Model{}).
			Scan(&_models).Error; err != nil {
			return nil
		}
		results := make(map[uint]string)
		for _, _model := range _models {
			results[_model.ID] = _model.Name
		}
		return results
	}
	// 处理数据
	_results := make(map[uint][]uint)
	for _, modelInstance := range modelInstances {
		_results[modelInstance.TargetModelId] = append(_results[modelInstance.TargetModelId], modelInstance.TargetInstanceId)
	}
	// 整理数据
	results := make([]resp.ListInstanceRelation, 0)
	modelIdToModelDisplayResult := modelIdToModelDisplay()
	for modelId, instances := range _results {
		modelDisplay, ok := modelIdToModelDisplayResult[modelId]
		if !ok {
			modelDisplay = strconv.Itoa(int(modelId))
		}
		results = append(results, resp.ListInstanceRelation{
			ModelId:      modelId,
			ModelDisplay: modelDisplay,
			Count:        int64(len(instances)),
			Instances:    instances,
		})
	}
	response.Success(c, "查询成功", results)
}

// TargetSourceInstanceRelation
//
//	@Description: 目标到源的实例关系
//	@receiver i
//	@param c
func (i *instanceRelation) TargetSourceInstanceRelation(c *gin.Context) {
	targetId := c.Param("source_id")
	modelInstances := make([]models.InstanceRelation, 0)
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{"target_id": targetId}).
		Scan(&modelInstances).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 匿名函数，模型ID对应的展示名称
	modelIdToModelDisplay := func() map[uint]string {
		var _models []models.Model
		if err := database.DB.Model(&models.Model{}).
			Scan(&_models).Error; err != nil {
			return nil
		}
		results := make(map[uint]string)
		for _, _model := range _models {
			results[_model.ID] = _model.Name
		}
		return results
	}
	// 处理数据
	_results := make(map[uint][]uint)
	for _, modelInstance := range modelInstances {
		_results[modelInstance.SourceModelId] = append(_results[modelInstance.SourceModelId], modelInstance.SourceInstanceId)
	}
	// 整理数据
	results := make([]resp.ListInstanceRelation, 0)
	modelIdToModelDisplayResult := modelIdToModelDisplay()
	for modelId, instances := range _results {
		modelDisplay, ok := modelIdToModelDisplayResult[modelId]
		if !ok {
			modelDisplay = strconv.Itoa(int(modelId))
		}
		results = append(results, resp.ListInstanceRelation{
			ModelId:      modelId,
			ModelDisplay: modelDisplay,
			Count:        int64(len(instances)),
			Instances:    instances,
		})
	}
	response.Success(c, "查询成功", results)

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
