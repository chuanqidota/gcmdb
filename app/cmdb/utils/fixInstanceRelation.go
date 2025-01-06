package utils

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
)

// 维护实例关系
type fix struct {
}

var Fix = new(fix)

// CreateModelFieldRelation
//
//	@Description: 创建模型字段关联-维护实例关系
//	@receiver f
func (f *fix) CreateModelFieldRelation(sourceModelId, targetModelId, sourceFieldId, targetFieldId uint) error{
	modelFields := make([]models.ModelField,0)
	if err:=database.DB.Model(&models.ModelField{}).
		Where("id in ?",[]uint{sourceFieldId,targetFieldId}).
		Scan(&modelFields).Error;err!=nil{
			return fmt.Errorf("查询出错-%s",err.Error())
		}
	filedId2Alias := make(map[uint]string)
	for _,modelField := range modelFields {
		filedId2Alias[modelField.ID] = modelField.Alias
	}

	sourceFiled,ok := filedId2Alias[sourceFieldId]
	if !ok{
		return fmt.Errorf("没有找到源字段ID对应的别名")
	}
	targetField,ok := filedId2Alias[targetFieldId]
	if !ok{
		return fmt.Errorf("没有找到目标字段ID对应的别名")
	}
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
							instance2.model_id = %+v`, sourceFiled, targetField, sourceModelId, targetModelId)

	instanceRelations := make([]models.InstanceRelation, 0)
	if err := database.DB.Raw(sql).Scan(&instanceRelations).Error; err != nil {
		return fmt.Errorf("查询失败-%s", err.Error())
	}

	if err := database.DB.Model(&models.InstanceRelation{}).
		CreateInBatches(instanceRelations, 100).Error; err != nil {
		return fmt.Errorf("创建失败-%s",err.Error())
	}
	return nil
}

// DeleteModelFieldRelation
//
//	@Description: 删除模型字段关联-维护实例关系
//	@receiver f
func (f *fix) DeleteModelFieldRelation(sourceModelId, targetModelId, sourceFieldId, targetFieldId uint) error {
	modelFields := make([]models.ModelField,0)
	if err:=database.DB.Model(&models.ModelField{}).
		Where("id in ?",[]uint{sourceFieldId,targetFieldId}).
		Scan(&modelFields).Error;err!=nil{
			return fmt.Errorf("查询出错-%s",err.Error())
		}
	filedId2Alias := make(map[uint]string)
	for _,modelField := range modelFields {
		filedId2Alias[modelField.ID] = modelField.Alias
	}

	sourceFiled,ok := filedId2Alias[sourceFieldId]
	if !ok{
		return fmt.Errorf("没有找到源字段ID对应的别名")
	}
	targetField,ok := filedId2Alias[targetFieldId]
	if !ok{
		return fmt.Errorf("没有找到目标字段ID对应的别名")
	}
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
							instance2.model_id = %+v`, sourceFiled, targetField, sourceModelId, targetModelId)

	instanceRelations := make([]models.InstanceRelation, 0)
	if err := database.DB.Raw(sql).Scan(&instanceRelations).Error; err != nil {
		return fmt.Errorf("查询失败-%s", err.Error())
	}
	if err := database.DB.Delete(&instanceRelations).Error; err != nil {
		return fmt.Errorf("删除失败-%s", err.Error())
	}

	return nil
}

// CreateInstance
//
//	@Description: 创建实例-维护实例关系
//	@receiver f
func (f *fix) CreateInstance() {

}

// DeleteInstance
//
//	@Description: 删除实例-维护实例关系
//	@receiver f
func (f *fix) DeleteInstance() {

}

// SyncSourceModelInstanceRelation
//
//	@Description: 同步指定源模型的实例关系
//	@receiver f
func (f *fix) SyncSourceModelInstanceRelation() {

}
