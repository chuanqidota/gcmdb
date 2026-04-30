package utils

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/openapi/params"
	"gcmdb/pkg/database"
	"gorm.io/gorm"
)

type instanceRelation struct {
}

var InstanceRelation = new(instanceRelation)

// CreateInstanceRelation
//
//	@Description: 创建实例关系
//	@receiver ir
//	@param body
//	@return error
func (ir *instanceRelation) CreateInstanceRelation(body params.CreateInstanceRelation) error {
	SourceModel := body.SourceModel
	TargetModel := body.TargetModel
	SourceInstanceId := body.SourceInstanceId
	TargetInstanceIds := body.TargetInstanceIds
	// 模型别名对应的ID
	var _models []models.Model
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"alias": []string{SourceModel, TargetModel}}).
		Scan(&_models).Error; err != nil {
		return err
	}
	modelAlias2Id := make(map[string]uint)
	for _, item := range _models {
		modelAlias2Id[item.Alias] = item.ID
	}

	sourceModelId, ok := modelAlias2Id[SourceModel]
	if !ok {
		return fmt.Errorf("没有找到源模型:%s", SourceModel)
	}
	targetModelId, ok := modelAlias2Id[TargetModel]
	if !ok {
		return fmt.Errorf("没有找到目标模型:%s", TargetModel)
	}
	// 创建的数据
	createData := make([]models.InstanceRelation, 0)
	for _, targetInstanceId := range TargetInstanceIds {
		createData = append(createData, models.InstanceRelation{
			SourceModelId:    sourceModelId,
			TargetModelId:    targetModelId,
			SourceInstanceId: SourceInstanceId,
			TargetInstanceId: targetInstanceId,
		})
	}
	// 去重：先查询已存在的关系
	for _, data := range createData {
		var existing models.InstanceRelation
		result := database.DB.Model(&models.InstanceRelation{}).
			Where(map[string]any{
				"source_model_id": data.SourceModelId,
				"target_model_id": data.TargetModelId,
				"source_id":       data.SourceInstanceId,
				"target_id":       data.TargetInstanceId,
			}).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := database.DB.Model(&models.InstanceRelation{}).Create(&data).Error; err != nil {
				return err
			}
		} else if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// DeleteInstanceRelation
//
//	@Description: 删除实例关系
//	@receiver ir
//	@param body
//	@return error
func (ir *instanceRelation) DeleteInstanceRelation(body params.DeleteInstanceRelation) error {
	SourceModel := body.SourceModel
	TargetModel := body.TargetModel
	SourceInstanceId := body.SourceInstanceId
	TargetInstanceId := body.TargetInstanceId

	// 模型别名对应的ID
	var _models []models.Model
	if err := database.DB.Model(&models.Model{}).
		Where(map[string]any{"alias": []string{SourceModel, TargetModel}}).
		Scan(&_models).Error; err != nil {
		return err
	}
	modelAlias2Id := make(map[string]uint)
	for _, item := range _models {
		modelAlias2Id[item.Alias] = item.ID
	}

	sourceModelId, ok := modelAlias2Id[SourceModel]
	if !ok {
		return fmt.Errorf("没有找到源模型:%s", SourceModel)
	}
	targetModelId, ok := modelAlias2Id[TargetModel]
	if !ok {
		return fmt.Errorf("没有找到目标模型:%s", TargetModel)
	}

	if err := database.DB.Unscoped().Model(&models.InstanceRelation{}).
		Where(map[string]any{
			"source_model_id": sourceModelId,
			"target_model_id": targetModelId,
			"source_id":       SourceInstanceId,
			"target_id":       TargetInstanceId,
		}).Delete(&models.InstanceRelation{}).Error; err != nil {
		return err
	}
	return nil

}
