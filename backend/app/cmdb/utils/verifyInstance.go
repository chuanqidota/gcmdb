package utils

import (
	"encoding/json"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gorm.io/datatypes"
	"slices"
	"strings"
)

// 校验实例能否被调整
type verify struct {
}

var Verify = new(verify)

// VerifyCreateInstance
//
//	@Description: 校验创建实例(验证必须的字段、删除多余字段、验证唯一性)
//	@receiver v
//	@param modelId 模型id
//	@param data 实例数据
//	@return datatypes.JSON
//	@return error
func (v *verify) VerifyCreateInstance(modelId uint, data datatypes.JSON) (datatypes.JSON, error) {
	// 校验字段是否必须有
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Or(map[string]any{"model_id": modelId}).
		Scan(&modelFields).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	// 需要的字段
	requiredFields := make([]string, 0)
	filedAlias := make([]string, 0)
	for _, modelField := range modelFields {
		if modelField.IsRequired {
			requiredFields = append(requiredFields, modelField.Alias)
		}
		filedAlias = append(filedAlias, modelField.Alias)
	}

	// 验证字段是否必须
	var dataMap map[string]any
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, fmt.Errorf("json解析失败-%s", err.Error())
	}
	for _, requiredField := range requiredFields {
		if _, ok := dataMap[requiredField]; !ok {
			return nil, fmt.Errorf("字段%s必须", requiredField)
		}
	}

	// 删除多余字段
	for key := range dataMap {
		if !slices.Contains(filedAlias, key) {
			delete(dataMap, key)
		}
	}

	// 补充字段
	for _, modelField := range modelFields {
		if _, ok := dataMap[modelField.Alias]; !ok {
			if models.DefaultValueByType[modelField.Type] != nil {
				dataMap[modelField.Alias] = models.DefaultValueByType[modelField.Type]
			}
		}
	}

	// 验证唯一性
	modelFieldsUniques := make([]models.ModelFieldUnique, 0)
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": modelId}).
		Scan(&modelFieldsUniques).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	db := database.DB.Model(&models.Instance{}).Where(map[string]any{"model_id": modelId})
	for _, modelFieldUnique := range modelFieldsUniques {
		fields := strings.Split(modelFieldUnique.Fields, ",") // 获取字段
		conditions := make([]string, 0)
		_params := make([]any, len(fields))
		for i, _field := range fields {
			conditions = append(conditions, fmt.Sprintf("data->'$.%s' = ?", _field))
			_params[i] = fmt.Sprintf("%%%s%%", dataMap[_field])
		}
		query := strings.Join(conditions, " AND ")
		if err := db.Where(query, _params...).
			First(&models.Instance{}).Error; err == nil {
			return nil, fmt.Errorf("字段%s值重复", modelFieldUnique.Fields)
		}
	}
	// 响应数据
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("json数据异常-%s", err.Error())
	}
	return jsonData, nil

}

// VerifyUpdateInstance
//
//	@Description: 校验更新实例(删除多余字段、验证唯一性)
//	@receiver v
//	@param modelId 模型id
//	@param data 实例数据
//	@return datatypes.JSON
//	@return error
func (v *verify) VerifyUpdateInstance(id uint, data datatypes.JSON) (datatypes.JSON, error) {
	// 查询实例
	var instance models.Instance
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id}).Scan(&instance).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	// 删除和补充字段
	var dataMap map[string]any
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, fmt.Errorf("json解析失败-%s", err.Error())
	}
	var instanceData map[string]any
	if err := json.Unmarshal(instance.Data, &instanceData); err != nil {
		return nil, fmt.Errorf("json解析失败-%s", err.Error())
	}
	// 删除多余字段
	for key := range dataMap {
		_, ok := instanceData[key]
		if !ok {
			delete(dataMap, key)
		}
	}
	// 补充字段值
	for key, value := range instanceData {
		if _, ok := dataMap[key]; !ok {
			dataMap[key] = value
		}
	}

	// 验证唯一性
	modelFieldsUniques := make([]models.ModelFieldUnique, 0)
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": instance.ModelId}).
		Scan(&modelFieldsUniques).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	db := database.DB.Model(&models.Instance{}).Where(map[string]any{"model_id": instance.ModelId})
	for _, modelFieldUnique := range modelFieldsUniques {
		fields := strings.Split(modelFieldUnique.Fields, ",") // 获取字段
		conditions := make([]string, 0)
		_params := make([]any, len(fields))
		for i, _field := range fields {
			conditions = append(conditions, fmt.Sprintf("data->'$.%s' = ?", _field))
			_params[i] = fmt.Sprintf("%%%s%%", dataMap[_field])
		}
		query := strings.Join(conditions, " AND ")
		if err := db.Where(query, _params...).
			Not(map[string]any{"id": id}).
			First(&models.Instance{}).Error; err == nil {
			return nil, fmt.Errorf("字段%s值重复", modelFieldUnique.Fields)
		}
	}
	// 响应数据
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("json数据异常-%s", err.Error())
	}
	return jsonData, nil

}
