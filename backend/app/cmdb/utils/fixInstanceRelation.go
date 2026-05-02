package utils

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/config"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	fieldId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		fieldId2Alias[modelField.ID] = modelField.Alias
	}

	sourceField, ok := fieldId2Alias[sourceFieldId]
	if !ok {
		return fmt.Errorf("没有找到源字段ID对应的别名")
	}
	targetField, ok := fieldId2Alias[targetFieldId]
	if !ok {
		return fmt.Errorf("没有找到目标字段ID对应的别名")
	}

	sql := `SELECT
				instance1.model_id as source_model_id,
				instance2.model_id as target_model_id,
				instance1.id as source_id,
				instance2.id as target_id
			FROM
				instance instance1
			INNER JOIN
				instance instance2
			ON
				JSON_EXTRACT(instance1.data, CONCAT('$.', ?)) = JSON_EXTRACT(instance2.data, CONCAT('$.', ?))
			WHERE
				instance1.model_id = ?
			AND
				instance2.model_id = ?`

	instanceRelations := make([]models.InstanceRelation, 0)
	if err := database.DB.Raw(sql, sourceField, targetField, sourceModelId, targetModelId).Scan(&instanceRelations).Error; err != nil {
		return fmt.Errorf("查询失败-%s", err.Error())
	}
	if len(instanceRelations) == 0 {
		return nil
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
	fieldId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		fieldId2Alias[modelField.ID] = modelField.Alias
	}

	sourceField, ok := fieldId2Alias[sourceFieldId]
	if !ok {
		return fmt.Errorf("没有找到源字段ID对应的别名")
	}
	targetField, ok := fieldId2Alias[targetFieldId]
	if !ok {
		return fmt.Errorf("没有找到目标字段ID对应的别名")
	}

	sql := `SELECT
				instance1.model_id as source_model_id,
				instance2.model_id as target_model_id,
				instance1.id as source_id,
				instance2.id as target_id
			FROM
				instance instance1
			INNER JOIN
				instance instance2
			ON
				JSON_EXTRACT(instance1.data, CONCAT('$.', ?)) = JSON_EXTRACT(instance2.data, CONCAT('$.', ?))
			WHERE
				instance1.model_id = ?
			AND
				instance2.model_id = ?`

	instanceRelations := make([]models.InstanceRelation, 0)
	if err := database.DB.Raw(sql, sourceField, targetField, sourceModelId, targetModelId).Scan(&instanceRelations).Error; err != nil {
		return fmt.Errorf("查询失败-%s", err.Error())
	}
	if len(instanceRelations) == 0 {
		return nil
	}
	// 按源/目标模型ID删除匹配的实例关系
	if err := database.DB.Model(&models.InstanceRelation{}).
		Where("source_model_id = ? AND target_model_id = ?", sourceModelId, targetModelId).
		Delete(&models.InstanceRelation{}).Error; err != nil {
		return fmt.Errorf("删除失败-%s", err.Error())
	}

	return nil
}

// CreateInstance
//
//	@Description: 创建实例-维护实例关系
//	@receiver f
func (f *instanceRelation) CreateInstance(ModelId uint, instanceId uint) error {
	if !config.Conf.Server.AutoRelation {
		return nil
	}
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": ModelId}).
		Scan(&modelFieldRelations).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}
	if len(modelFieldRelations) == 0 {
		return nil
	}

	// 只查询相关模型的字段，而非全表扫描
	relatedModelIds := make(map[uint]bool)
	for _, mfr := range modelFieldRelations {
		relatedModelIds[mfr.SourceModelId] = true
		relatedModelIds[mfr.TargetModelId] = true
	}
	modelIdSlice := make([]uint, 0, len(relatedModelIds))
	for id := range relatedModelIds {
		modelIdSlice = append(modelIdSlice, id)
	}

	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where("model_id in ?", modelIdSlice).
		Scan(&modelFields).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}

	fieldId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		fieldId2Alias[modelField.ID] = modelField.Alias
	}

	for _, modelFieldRelation := range modelFieldRelations {
		sourceField, ok := fieldId2Alias[modelFieldRelation.SourceFieldId]
		if !ok {
			continue
		}
		targetField, ok := fieldId2Alias[modelFieldRelation.TargetFieldId]
		if !ok {
			continue
		}
		sourceModelId := modelFieldRelation.SourceModelId
		targetModelId := modelFieldRelation.TargetModelId

		sql := `SELECT
					instance1.model_id as source_model_id,
					instance2.model_id as target_model_id,
					instance1.id as source_id,
					instance2.id as target_id
				FROM
					instance instance1
				INNER JOIN
					instance instance2
				ON
					JSON_EXTRACT(instance1.data, CONCAT('$.', ?)) = JSON_EXTRACT(instance2.data, CONCAT('$.', ?))
				WHERE
					instance1.model_id = ?
				AND
					instance2.model_id = ?
				AND
					(instance1.id = ? OR instance2.id = ?)`

		instanceRelations := make([]models.InstanceRelation, 0)
		if err := database.DB.Raw(sql, sourceField, targetField, sourceModelId, targetModelId, instanceId, instanceId).Scan(&instanceRelations).Error; err != nil {
			return fmt.Errorf("查询失败-%s", err.Error())
		}
		if len(instanceRelations) == 0 {
			continue
		}
		if err := database.DB.Model(&models.InstanceRelation{}).
			Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(instanceRelations, 100).Error; err != nil {
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
	if !config.Conf.Server.AutoRelation {
		return nil
	}
	if err := database.DB.Unscoped().Model(&models.InstanceRelation{}).
		Where(map[string]any{"source_id": id}).
		Or(map[string]any{"target_id": id}).
		Delete(&models.InstanceRelation{}).Error; err != nil {
		return fmt.Errorf("删除错误-%s", err.Error())
	}
	return nil
}

// MulDeleteInstance
//
//	@Description: 批量删除实例关系
//	@receiver f
//	@param ids
//	@return error
func (f *instanceRelation) MulDeleteInstance(ids []uint) error {
	if !config.Conf.Server.AutoRelation {
		return nil
	}
	if err := database.DB.Unscoped().Model(&models.InstanceRelation{}).
		Where("source_id in ?", ids).
		Or("target_id in ?", ids).
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
	// 查询模型字段关联
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": modelId}).
		Scan(&modelFieldRelations).Error; err != nil {
		return fmt.Errorf("查询模型字段关联出错")
	}

	// 只查询相关模型的字段
	relatedModelIds := make(map[uint]bool)
	for _, mfr := range modelFieldRelations {
		relatedModelIds[mfr.SourceModelId] = true
		relatedModelIds[mfr.TargetModelId] = true
	}
	modelIdSlice := make([]uint, 0, len(relatedModelIds))
	for id := range relatedModelIds {
		modelIdSlice = append(modelIdSlice, id)
	}

	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where("model_id in ?", modelIdSlice).
		Scan(&modelFields).Error; err != nil {
		return fmt.Errorf("查询出错-%s", err.Error())
	}
	// 构建字段id和别名的对应关系
	fieldId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		fieldId2Alias[modelField.ID] = modelField.Alias
	}

	// 使用事务：先删除再重建
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.InstanceRelation{}).
			Where(map[string]any{"source_model_id": modelId}).
			Delete(&models.InstanceRelation{}).Error; err != nil {
			return fmt.Errorf("删除实例关联出错")
		}

		// 遍历模型字段关联，查询实例关联
		for _, modelFieldRelation := range modelFieldRelations {
			sourceModelId := modelFieldRelation.SourceModelId
			targetModelId := modelFieldRelation.TargetModelId
			sourceFieldId := modelFieldRelation.SourceFieldId
			targetFieldId := modelFieldRelation.TargetFieldId

			sourceField, ok := fieldId2Alias[sourceFieldId]
			if !ok {
				continue
			}
			targetField, ok := fieldId2Alias[targetFieldId]
			if !ok {
				continue
			}

			sql := `SELECT
						instance1.model_id as source_model_id,
						instance2.model_id as target_model_id,
						instance1.id as source_id,
						instance2.id as target_id
					FROM
						instance instance1
					INNER JOIN
						instance instance2
					ON
						JSON_EXTRACT(instance1.data, CONCAT('$.', ?)) = JSON_EXTRACT(instance2.data, CONCAT('$.', ?))
					WHERE
						instance1.model_id = ?
					AND
						instance2.model_id = ?`

			instanceRelations := make([]models.InstanceRelation, 0)
			if err := tx.Raw(sql, sourceField, targetField, sourceModelId, targetModelId).Scan(&instanceRelations).Error; err != nil {
				logger.Error(fmt.Sprintf("查询失败-%s", err.Error()))
				continue
			}
			if len(instanceRelations) == 0 {
				continue
			}
			if err := tx.Model(&models.InstanceRelation{}).
				Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(instanceRelations, 100).Error; err != nil {
				logger.Error(fmt.Sprintf("插入失败-%s", err.Error()))
				continue
			}
		}
		return nil
	})
}
