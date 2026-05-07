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

// instanceRelation 实例关系维护工具
// 核心逻辑：当模型定义了字段关联（ModelFieldRelation）时，自动同步实例级别的关系（InstanceRelation）
// 例如：模型 A 的 host_ip 字段关联模型 B 的 ip 字段，则当实例 A.host_ip = "1.1.1.1" 时
// 自动查找 B 中 ip = "1.1.1.1" 的实例并建立关系
type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

// findMatchingRelations 查询通过字段值匹配的实例关系
// baseOnly=true: 只查询两端模型匹配的关系
// baseOnly=false: 额外限定 instanceId 必须参与
func (f *instanceRelation) findMatchingRelations(sourceField, targetField string, sourceModelId, targetModelId uint, instanceId uint, baseOnly bool) ([]models.InstanceRelation, error) {
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
	args := []any{sourceField, targetField, sourceModelId, targetModelId}
	if !baseOnly {
		sql += "\n\t\t\tAND\n\t\t\t\t(instance1.id = ? OR instance2.id = ?)"
		args = append(args, instanceId, instanceId)
	}
	instanceRelations := make([]models.InstanceRelation, 0)
	if err := database.DB.Raw(sql, args...).Scan(&instanceRelations).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	return instanceRelations, nil
}

// CreateModelFieldRelation
//
//	@Description: 创建模型字段关联-维护实例关系
//	@receiver f
func (f *instanceRelation) CreateModelFieldRelation(sourceModelId, targetModelId, sourceFieldId, targetFieldId uint) error {
	sourceField, targetField, err := f.lookupFieldAliases(sourceFieldId, targetFieldId)
	if err != nil {
		return err
	}

	instanceRelations, err := f.findMatchingRelations(sourceField, targetField, sourceModelId, targetModelId, 0, true)
	if err != nil {
		return err
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
	sourceField, targetField, err := f.lookupFieldAliases(sourceFieldId, targetFieldId)
	if err != nil {
		return err
	}

	instanceRelations, err := f.findMatchingRelations(sourceField, targetField, sourceModelId, targetModelId, 0, true)
	if err != nil {
		return err
	}
	if len(instanceRelations) == 0 {
		return nil
	}
	// 只删除通过该字段对匹配的特定实例关系，不影响其他字段关联的实例关系
	if err := f.deleteMatchingRelations(instanceRelations); err != nil {
		return err
	}
	return nil
}

// deleteMatchingRelations 根据复合唯一键删除 instance_relation 记录
func (f *instanceRelation) deleteMatchingRelations(relations []models.InstanceRelation) error {
	if len(relations) == 0 {
		return nil
	}
	for _, ir := range relations {
		if err := database.DB.Unscoped().Model(&models.InstanceRelation{}).
			Where("source_model_id = ? AND target_model_id = ? AND source_id = ? AND target_id = ?",
				ir.SourceModelId, ir.TargetModelId, ir.SourceInstanceId, ir.TargetInstanceId).
			Delete(&models.InstanceRelation{}).Error; err != nil {
			return fmt.Errorf("删除失败-%s", err.Error())
		}
	}
	return nil
}

// lookupFieldAliases 根据字段ID查询字段别名
func (f *instanceRelation) lookupFieldAliases(sourceFieldId, targetFieldId uint) (string, string, error) {
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where("id in ?", []uint{sourceFieldId, targetFieldId}).
		Scan(&modelFields).Error; err != nil {
		return "", "", fmt.Errorf("查询出错-%s", err.Error())
	}
	fieldId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		fieldId2Alias[modelField.ID] = modelField.Alias
	}
	sourceField, ok := fieldId2Alias[sourceFieldId]
	if !ok {
		return "", "", fmt.Errorf("没有找到源字段ID对应的别名")
	}
	targetField, ok := fieldId2Alias[targetFieldId]
	if !ok {
		return "", "", fmt.Errorf("没有找到目标字段ID对应的别名")
	}
	return sourceField, targetField, nil
}

// CreateInstance
//
//	@Description: 创建实例-维护实例关系
//	@receiver f
func (f *instanceRelation) CreateInstance(ModelId uint, instanceId uint) error {
	if !config.Conf.Server.AutoRelation {
		return nil
	}
	modelFieldRelations, fieldId2Alias, err := f.loadRelatedFieldRelations(ModelId)
	if err != nil {
		return err
	}
	if len(modelFieldRelations) == 0 {
		return nil
	}

	for _, mfr := range modelFieldRelations {
		sourceField, ok := fieldId2Alias[mfr.SourceFieldId]
		if !ok {
			continue
		}
		targetField, ok := fieldId2Alias[mfr.TargetFieldId]
		if !ok {
			continue
		}

		instanceRelations, err := f.findMatchingRelations(sourceField, targetField, mfr.SourceModelId, mfr.TargetModelId, instanceId, false)
		if err != nil {
			return err
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

// loadRelatedFieldRelations 加载模型字段关联及相关字段别名
func (f *instanceRelation) loadRelatedFieldRelations(modelId uint) ([]models.ModelFieldRelation, map[uint]string, error) {
	modelFieldRelations := make([]models.ModelFieldRelation, 0)
	if err := database.DB.Model(&models.ModelFieldRelation{}).
		Where(map[string]any{"source_model_id": modelId}).
		Scan(&modelFieldRelations).Error; err != nil {
		return nil, nil, fmt.Errorf("查询模型字段关联出错")
	}
	if len(modelFieldRelations) == 0 {
		return nil, nil, nil
	}

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
		return nil, nil, fmt.Errorf("查询出错-%s", err.Error())
	}

	fieldId2Alias := make(map[uint]string)
	for _, modelField := range modelFields {
		fieldId2Alias[modelField.ID] = modelField.Alias
	}
	return modelFieldRelations, fieldId2Alias, nil
}

// SyncSourceModelInstanceRelation
//
//	@Description: 同步指定源模型的实例关系
//	@receiver f
func (f *instanceRelation) SyncSourceModelInstanceRelation(modelId uint) error {
	modelFieldRelations, fieldId2Alias, err := f.loadRelatedFieldRelations(modelId)
	if err != nil {
		return err
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Model(&models.InstanceRelation{}).
			Where(map[string]any{"source_model_id": modelId}).
			Delete(&models.InstanceRelation{}).Error; err != nil {
			return fmt.Errorf("删除实例关联出错")
		}

		for _, mfr := range modelFieldRelations {
			sourceField, ok := fieldId2Alias[mfr.SourceFieldId]
			if !ok {
				continue
			}
			targetField, ok := fieldId2Alias[mfr.TargetFieldId]
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
			if err := tx.Raw(sql, sourceField, targetField, mfr.SourceModelId, mfr.TargetModelId).Scan(&instanceRelations).Error; err != nil {
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
