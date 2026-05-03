package utils

import (
	"encoding/json"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	"strings"

	"gorm.io/gorm"
)

type indexedValues struct{}

var IndexedValues = new(indexedValues)

// BuildIndexedValues
//
//	@Description: 从实例 data 中提取索引字段的值，构建逗号分隔的索引字符串（首尾带逗号，便于 LIKE 查询）
//	@param modelId 模型ID
//	@param data 实例数据 JSON
//	@return string 索引值，如 ",value1,value2,"
func (v *indexedValues) BuildIndexedValues(modelId uint, data []byte) string {
	// 查询该模型下标记为索引的字段
	indexedFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where("model_id = ? AND is_indexed = ?", modelId, true).
		Select("alias").
		Scan(&indexedFields).Error; err != nil {
		logger.Error(fmt.Sprintf("查询索引字段失败-%s", err.Error()))
		return ""
	}
	if len(indexedFields) == 0 {
		return ""
	}

	// 解析 data JSON
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return ""
	}

	// 提取索引字段值
	values := make([]string, 0, len(indexedFields))
	for _, field := range indexedFields {
		if val, ok := dataMap[field.Alias]; ok && val != nil {
			values = append(values, fmt.Sprintf("%v", val))
		}
	}

	return "," + strings.Join(values, ",") + ","
}

// SyncModelIndexedValues
//
//	@Description: 刷新指定模型下所有实例的 indexed_values（字段 is_indexed 变更时调用）
//	@param modelId 模型ID
func (v *indexedValues) SyncModelIndexedValues(modelId uint) {
	// 查询索引字段列表（只查一次）
	indexedFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where("model_id = ? AND is_indexed = ?", modelId, true).
		Select("alias").
		Scan(&indexedFields).Error; err != nil {
		logger.Error(fmt.Sprintf("查询索引字段失败-%s", err.Error()))
		return
	}

	// 无索引字段时，清空所有实例的 indexed_values
	if len(indexedFields) == 0 {
		if err := database.DB.Model(&models.Instance{}).
			Where("model_id = ?", modelId).
			Update("indexed_values", "").Error; err != nil {
			logger.Error(fmt.Sprintf("清空索引值失败-%s", err.Error()))
		}
		return
	}

	// 查询该模型所有实例（只查一次）
	instances := make([]models.Instance, 0)
	if err := database.DB.Model(&models.Instance{}).
		Where("model_id = ?", modelId).
		Select("id", "data").
		Scan(&instances).Error; err != nil {
		logger.Error(fmt.Sprintf("查询实例失败-%s", err.Error()))
		return
	}
	if len(instances) == 0 {
		return
	}

	// 批量更新：CASE WHEN 一次 SQL 搞定
	caseSQL := "CASE id "
	args := make([]interface{}, 0)
	ids := make([]uint, 0, len(instances))
	for _, inst := range instances {
		indexedVal := v.buildIndexedValuesFromFields(indexedFields, inst.Data)
		caseSQL += "WHEN ? THEN ? "
		args = append(args, inst.ID, indexedVal)
		ids = append(ids, inst.ID)
	}
	caseSQL += "END"
	args = append(args, ids)

	if err := database.DB.Model(&models.Instance{}).
		Where("id IN ?", ids).
		Update("indexed_values", gorm.Expr(caseSQL, args...)).Error; err != nil {
		logger.Error(fmt.Sprintf("批量更新索引值失败-%s", err.Error()))
	}
}

// buildIndexedValuesFromFields 接收已查好的字段列表，不再重复查库
func (v *indexedValues) buildIndexedValuesFromFields(indexedFields []models.ModelField, data []byte) string {
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return ""
	}
	values := make([]string, 0, len(indexedFields))
	for _, field := range indexedFields {
		if val, ok := dataMap[field.Alias]; ok && val != nil {
			values = append(values, fmt.Sprintf("%v", val))
		}
	}
	return "," + strings.Join(values, ",") + ","
}
