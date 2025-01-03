package api

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"

	"github.com/gin-gonic/gin"
)

type modelFieldRelation struct {
}

var ModelFieldRelation = new(modelFieldRelation)

// CreateModelFieldRelation
//
//	@Description: 创建模型字段关联
//	@receiver m
//	@param c
func (m *modelFieldRelation) CreateModelFieldRelation(c *gin.Context) {
	var body models.ModelFieldRelation
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
		return
	}
	var results models.ModelFieldRelation
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		FirstOrCreate(&results, body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	sourceModel := body.SourceModelId
	targetModel := body.TargetModelId
	sourceFiled := body.SourceFieldId
	targetField := body.TargetFieldId

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
							data->'$.%s'=data->'$.%s'
						WHERE 
							instance1.model_id = %+v 
						AND 
							instance2.model_id = %+v`, sourceFiled, targetField, sourceModel, targetModel)
	instanceRelations := make([]models.InstanceRelation,0)
	if err:=database.DB.Raw(sql).Scan(&instanceRelations).Error;err!=nil{
		response.Fail(c,fmt.Sprintf("查询失败-%s",err.Error()))
		return
	}
	if err:=database.DB.Model(&models.InstanceRelation{}).
		CreateInBatches(instanceRelations,100).Error;err!=nil{
			response.Fail(c,fmt.Sprintf("创建失败-%s",err.Error()))
			return
		}
	response.Success(c, "创建成功", nil)
}

// ListModelFieldRelation
//
//	@Description: 从源模型展示所有的模型字段关联
//	@receiver m
//	@param c
func (m *modelFieldRelation) ListModelFieldRelation(c *gin.Context) {
	sourceModelId := c.Param("source_model_id")
	var modelFieldRelations []models.ModelFieldRelation
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": sourceModelId}).
		Scan(&modelFieldRelations).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	results := make([]resp.ListModelFieldRelation, 0)
	for _, modelFieldRelation := range modelFieldRelations {
		SourceModelDisplay, _ := modelFieldRelation.GetSourceModelId()
		TargetModelDisplay, _ := modelFieldRelation.GetTargetModelId()
		SourceFieldDisplay, _ := modelFieldRelation.GetSourceFieldId()
		TargetFieldDisplay, _ := modelFieldRelation.GetTargetFieldId()
		results = append(results, resp.ListModelFieldRelation{
			ModelFieldRelation: modelFieldRelation,
			SourceModelDisplay: SourceModelDisplay,
			TargetModelDisplay: TargetModelDisplay,
			SourceFieldDisplay: SourceFieldDisplay,
			TargetFieldDisplay: TargetFieldDisplay,
		})
	}
	response.Success(c, "查询成功", results)
}

// DeleteModelFieldRelation
//
//	@Description: 删除模型关系实例
//	@receiver m
//	@param c
func (m *modelFieldRelation) DeleteModelFieldRelation(c *gin.Context) {
	id := c.Param("id")
	// 查询模型字段关系实例
	var modelFieldRelation models.ModelFieldRelation
	if err:=database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"id":id}).
		Scan(&modelFieldRelation).Error;err!=nil{
			response.Fail(c,fmt.Sprintf("查询失败-%s",err.Error()))
			return
		}
	
	sourceModel := modelFieldRelation.SourceModelId
	targetModel := modelFieldRelation.TargetModelId
	sourceFiled := modelFieldRelation.SourceFieldId
	targetField := modelFieldRelation.TargetFieldId

	// 查询需要删除那些实例关系
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
							data->'$.%s'=data->'$.%s'
						WHERE 
							instance1.model_id = %+v 
						AND 
							instance2.model_id = %+v`, sourceFiled, targetField, sourceModel, targetModel)
	
	instanceRelations := make([]models.InstanceRelation,0)
	if err:=database.DB.Raw(sql).Scan(&instanceRelations).Error;err!=nil{
		response.Fail(c,fmt.Sprintf("查询失败-%s",err.Error()))
		return
	}
	// 删除实例关系
	if err:=database.DB.Delete(&instanceRelations).Error;err!=nil{
		response.Fail(c,fmt.Sprintf("查询失败-%s",err.Error()))
		return
	}
	// 删除模型字段关系
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"id": id}).
		Delete(&models.ModelFieldRelation{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("参数校验失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)
}
