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
	i.listInstanceRelation(c, "source_id", sourceId, func(rel models.InstanceRelation) (uint, uint) {
		return rel.TargetModelId, rel.TargetInstanceId
	})
}

// TargetSourceInstanceRelation
//
//	@Description: 目标到源的实例关系
//	@receiver i
//	@param c
func (i *instanceRelation) TargetSourceInstanceRelation(c *gin.Context) {
	targetId := c.Param("target_id")
	i.listInstanceRelation(c, "target_id", targetId, func(rel models.InstanceRelation) (uint, uint) {
		return rel.SourceModelId, rel.SourceInstanceId
	})
}

// listInstanceRelation 通用实例关系查询
func (i *instanceRelation) listInstanceRelation(c *gin.Context, whereField, whereValue string, extract func(models.InstanceRelation) (uint, uint)) {
	modelInstances := make([]models.InstanceRelation, 0)
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{whereField: whereValue}).
		Scan(&modelInstances).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	// 查询模型ID对应的展示名称
	var _models []models.Model
	if err := database.DB.Model(&models.Model{}).Scan(&_models).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	modelDisplayMap := make(map[uint]string)
	for _, _model := range _models {
		modelDisplayMap[_model.ID] = _model.Name
	}
	// 按模型分组
	grouped := make(map[uint][]uint)
	for _, rel := range modelInstances {
		modelId, instanceId := extract(rel)
		grouped[modelId] = append(grouped[modelId], instanceId)
	}
	// 整理结果
	results := make([]resp.ListInstanceRelation, 0)
	for modelId, instances := range grouped {
		modelDisplay, ok := modelDisplayMap[modelId]
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

// DeleteInstanceRelationByKeys
//
//	@Description: 通过源/目标模型和实例ID删除实例关系
//	@receiver i
//	@param c
func (i *instanceRelation) DeleteInstanceRelationByKeys(c *gin.Context) {
	sourceModelId := c.Query("source_model_id")
	targetModelId := c.Query("target_model_id")
	sourceId := c.Query("source_id")
	targetId := c.Query("target_id")
	if sourceModelId == "" || targetModelId == "" || sourceId == "" || targetId == "" {
		response.Fail(c, "参数不完整，需要 source_model_id, target_model_id, source_id, target_id")
		return
	}
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{
			"source_model_id": sourceModelId,
			"target_model_id": targetModelId,
			"source_id":       sourceId,
			"target_id":       targetId,
		}).
		Delete(&models.InstanceRelation{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
