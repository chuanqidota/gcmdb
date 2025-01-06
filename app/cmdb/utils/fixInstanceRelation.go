package utils

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
)

// 维护实例关系
type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

// CreateModelFieldRelation
//
//	@Description: 创建模型字段关联-维护实例关系
//	@receiver f
func (f *instanceRelation) CreateModelFieldRelation(sourceModelId, targetModelId, sourceFieldId, targetFieldId uint) error {
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where("id in ?", []uint{sourceFieldId, targetFieldId}).
		Scan(&modelFields).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}
	filedId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		filedId2Alias[modelField.ID] = modelField.Alias
	}

	sourceFiled, ok := filedId2Alias[sourceFieldId]
	if !ok {
		return fmt.Errorf("没有找到源字段ID对应的别名")
	}
	targetField, ok := filedId2Alias[targetFieldId]
	if !ok {
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
		return fmt.Errorf("创建失败-%s", err.Error())
	}
	return nil
}

// DeleteModelFieldRelation
//
//	@Description: 删除模型字段关联-维护实例关系
//	@receiver f
func (f *instanceRelation) DeleteModelFieldRelation(sourceModelId, targetModelId, sourceFieldId, targetFieldId uint) error {
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where("id in ?", []uint{sourceFieldId, targetFieldId}).
		Scan(&modelFields).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}
	filedId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		filedId2Alias[modelField.ID] = modelField.Alias
	}

	sourceFiled, ok := filedId2Alias[sourceFieldId]
	if !ok {
		return fmt.Errorf("没有找到源字段ID对应的别名")
	}
	targetField, ok := filedId2Alias[targetFieldId]
	if !ok {
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
func (f *instanceRelation) CreateInstance(ModelId uint, instanceId uint) error {
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": ModelId}).
		Scan(&modelFieldRelations).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}

	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Scan(&modelFields).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}

	filedId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		filedId2Alias[modelField.ID] = modelField.Alias
	}

	for _, modelFieldRelation := range modelFieldRelations {
		sourceFiled, ok := filedId2Alias[modelFieldRelation.SourceFieldId]
		if !ok {
			return fmt.Errorf("没有找到源字段ID对应的别名")
		}
		targetField, ok := filedId2Alias[modelFieldRelation.TargetFieldId]
		if !ok {
			return fmt.Errorf("没有找到目标字段ID对应的别名")
		}
		sourceModelId := modelFieldRelation.SourceModelId
		targetModelId := modelFieldRelation.TargetModelId
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
										instance2.model_id = %+v
									AND 
										(source_id=%+v OR target_id =%+v)`, sourceFiled, targetField, sourceModelId, targetModelId, instanceId, instanceId)
		instanceRelations := make([]models.InstanceRelation, 0)
		if err := database.DB.Raw(sql).Scan(&instanceRelations).Error; err != nil {
			return fmt.Errorf("查询失败-%s", err.Error())
		}
		if err := database.DB.Model(&models.InstanceRelation{}).
			CreateInBatches(instanceRelations, 100).Error; err != nil {
			return fmt.Errorf("创建失败-%s", err.Error())
		}
	}
	return nil
}

// DeleteInstance
//
//	@Description: 删除实例-维护实例关系
//	@receiver f
func (f *instanceRelation) DeleteInstance(id uint) error {
	if err := database.DB.Unscoped().Model(&models.InstanceRelation{}).
		Where(map[string]any{"source_id": id}).
		Or(map[string]any{"target_id": id}).
		Delete(&models.InstanceRelation{}).Error; err != nil {
		return fmt.Errorf("删除错误-%s", err.Error())
	}
	return nil
}

// SyncSourceModelInstanceRelation
//
//	@Description: 同步指定源模型的实例关系
//	@receiver f
func (f *instanceRelation) SyncSourceModelInstanceRelation(modelId uint) error {

	if err := database.DB.Model(&models.InstanceRelation{}).
		Where(map[string]any{"source_model_id": modelId}).
		Delete(&models.InstanceRelation{}).Error; err != nil {
		return fmt.Errorf("删除实例关联出错")
	}
	// 查询模型字段关联
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": modelId}).
		Scan(&modelFieldRelations).Error; err != nil {
		return fmt.Errorf("查询模型字段关联出错")
	}
	// 查询所有模型字段
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Scan(&modelFields).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}
	// 构建字段id和别名的对应关系
	filedId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		filedId2Alias[modelField.ID] = modelField.Alias
	}
	// 遍历模型字段关联，查询实例关联
	for _, modelFieldRelation := range modelFieldRelations {
		sourceModelId := modelFieldRelation.SourceModelId
		targetModelId := modelFieldRelation.TargetModelId
		sourceFiledId := modelFieldRelation.SourceFieldId
		targetFieldId := modelFieldRelation.TargetFieldId

		sourceFiled, ok := filedId2Alias[sourceFiledId]
		if !ok {
			return fmt.Errorf("没有找到源字段ID对应的别名")
		}
		targetField, ok := filedId2Alias[targetFieldId]
		if !ok {
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
			return fmt.Errorf("插入失败-%s", err.Error())
		}
	}
	return nil
}
