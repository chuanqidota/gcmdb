package utils

import (
	"encoding/json"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gorm.io/datatypes"
	"strings"
)

type verify struct {
}

var Verify = new(verify)

// CreateInstqnce
//
//	@Description: 校验创建实例(验证必须的字段、删除多余字段、验证唯一性)
//	@receiver v
//	@param modelId 模型id
//	@param data 实例数据
//	@return datatypes.JSON
//	@return error
func (v *verify) CreateInstqnce(modelId uint, data datatypes.JSON) (datatypes.JSON, error) {
	// 校验字段是否必须有
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Or(map[string]any{"model_id": modelId}).
		Scan(&modelFields).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	// 需要的字段
	requiredFields := make([]string, 0)
	for _, modelField := range modelFields {
		if modelField.IsRequired {
			requiredFields = append(requiredFields, modelField.Alias)
		}
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
	// 多余的字段剔除
	for _, modelField := range modelFields {
		if _, ok := dataMap[modelField.Alias]; !ok {
			delete(dataMap, modelField.Alias)
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
		for _, field := range fields {
			if err := db.Where(fmt.Sprintf("JSON_EXTRACT(data,'$.%s') = ?", field), dataMap[field]).
				First(&models.Instance{}).Error; err == nil {
				return nil, fmt.Errorf("字段%s值重复", field)
			}
		}
	}
	// 响应数据
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("json数据异常-%s", err.Error())
	}
	return jsonData, nil

}

// UpdateInstance
//
//	@Description: 校验更新实例(删除多余字段、验证唯一性)
//	@receiver v
//	@param modelId 模型id
//	@param data 实例数据
//	@return datatypes.JSON
//	@return error
func (v *verify) UpdateInstance(modelId uint, data datatypes.JSON) (datatypes.JSON, error) {
	// 校验字段是否必须有
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Or(map[string]any{"model_id": modelId}).
		Scan(&modelFields).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	// 多余的字段剔除
	var dataMap map[string]any
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, fmt.Errorf("json解析失败-%s", err.Error())
	}
	for _, modelField := range modelFields {
		if _, ok := dataMap[modelField.Alias]; !ok {
			delete(dataMap, modelField.Alias)
		}
	}
	if len(dataMap) == 0 {
		return nil, fmt.Errorf("没有需要更新的字段在模型中")
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
		for _, field := range fields {
			if err := db.Where(fmt.Sprintf("JSON_EXTRACT(data,'$.%s') = ?", field), dataMap[field]).
				First(&models.Instance{}).Error; err == nil {
				return nil, fmt.Errorf("字段%s值重复", field)
			}
		}
	}
	// 响应数据
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("json数据异常-%s", err.Error())
	}
	return jsonData, nil

}
