package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/tasks/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
)

type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

// SyncInstanceRelation
//
//	@Description: 同步实例关系
//	@receiver ir
//	@param c
func (ir *instanceRelation) SyncInstanceRelation(c *gin.Context) {
	var body params.SyncInstanceRelation
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	modelId := body.ModelId
	go func(modelId uint) {
		// 只删除源模型的实例关联
		if err := database.DB.Model(&models.InstanceRelation{}).
			Where(map[string]any{"source_model_id": modelId}).
			Delete(&models.InstanceRelation{}).Error; err != nil {
			fmt.Printf("删除实例关联出错")
			return
		}
		// 查询模型字段关联
		modelFieldRelations := make([]models.ModelFieldRelation, 0)
		if err := database.DB.Model(&models.ModelFieldRelation{}).
			Where(map[string]any{"source_model_id": modelId}).
			Scan(&modelFieldRelations).Error; err != nil {
			fmt.Printf("查询模型字段关联出错")
			return
		}
		// 遍历模型字段关联，查询实例关联
		for _, modelFieldRelation := range modelFieldRelations {
			sourceModel := modelFieldRelation.SourceModelId
			targetModel := modelFieldRelation.TargetModelId
			sourceFiled := modelFieldRelation.SourceFieldId
			targetField := modelFieldRelation.TargetFieldId
			sql := fmt.Sprintf(`SELECT 
							instance1.model_id as source_model_id,
							instance2.model_id as target_model_id,
							instance1.id as source_id,
							instance2.id as target_id 
						FROM 
							instance instance1
						INNER JOIN 
							instance instance2 
						ON 
							data->'$.%+v'=data->'$.%+v'
						WHERE 
							instance1.model_id = %+v 
						AND 
							instance2.model_id = %+v`, sourceFiled, targetField, sourceModel, targetModel)
			instanceRelations := make([]models.InstanceRelation, 0)
			if err := database.DB.Raw(sql).Scan(&instanceRelations).Error; err != nil {
				response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
				return
			}
			if err := database.DB.Model(&models.InstanceRelation{}).
				CreateInBatches(instanceRelations, 100).Error; err != nil {
				response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
				return
			}
		}
	}(modelId)
	response.Success(c, "异步执行中", nil)
}
