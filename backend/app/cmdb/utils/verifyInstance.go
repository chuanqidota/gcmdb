package utils

import (
	"encoding/json"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gorm.io/datatypes"
	"regexp"
	"slices"
	"strings"
)

var dateRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// verifyUniqueness 验证实例数据的唯一性约束
// excludeId: 排除的实例ID（更新时传入当前实例ID，创建时传入0）
func verifyUniqueness(modelId uint, dataMap map[string]any, excludeId uint) error {
	modelFieldsUniques := make([]models.ModelFieldUnique, 0)
	if err := database.DB.Model(&models.ModelFieldUnique{}).
		Where(map[string]any{"model_id": modelId}).
		Scan(&modelFieldsUniques).Error; err != nil {
		return fmt.Errorf("查询失败-%s", err.Error())
	}
	for _, modelFieldUnique := range modelFieldsUniques {
		fields := strings.Split(modelFieldUnique.Fields, ",")
		conditions := make([]string, 0)
		_params := make([]any, len(fields))
		for i, _field := range fields {
			conditions = append(conditions, fmt.Sprintf("data->'$.%s' = ?", _field))
			_params[i] = fmt.Sprintf("%v", dataMap[_field])
		}
		query := strings.Join(conditions, " AND ")
		q := database.DB.Model(&models.Instance{}).
			Where(map[string]any{"model_id": modelId}).
			Where(query, _params...)
		if excludeId > 0 {
			q = q.Not(map[string]any{"id": excludeId})
		}
		var count int64
		if err := q.Count(&count).Error; err != nil {
			return fmt.Errorf("查询失败-%s", err.Error())
		}
		if count > 0 {
			return fmt.Errorf("字段%s值重复", modelFieldUnique.Fields)
		}
	}
	return nil
}
var datetimeRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)

// validateFieldValue 校验字段值是否匹配字段类型
func validateFieldValue(field models.ModelField, value any) error {
	if value == nil {
		return nil // nil 值由 required 检查处理
	}
	switch field.Type {
	case "number":
		switch value.(type) {
		case float64, float32, int, int64, uint, uint64:
			// JSON number 解码为 float64
		default:
			return fmt.Errorf("字段%s类型应为number", field.Alias)
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("字段%s类型应为bool", field.Alias)
		}
	case "date":
		s, ok := value.(string)
		if !ok || !dateRegexp.MatchString(s) {
			return fmt.Errorf("字段%s格式应为YYYY-MM-DD", field.Alias)
		}
	case "datetime":
		s, ok := value.(string)
		if !ok || !datetimeRegexp.MatchString(s) {
			return fmt.Errorf("字段%s格式应为YYYY-MM-DD HH:mm:ss", field.Alias)
		}
	case "json":
		// JSON 类型接受 object 和 array
		switch value.(type) {
		case map[string]any, []any:
			// ok
		default:
			return fmt.Errorf("字段%s类型应为JSON对象或数组", field.Alias)
		}
	case "enum":
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("字段%s类型应为字符串", field.Alias)
		}
		if field.Options != "" {
			var options []string
			if err := json.Unmarshal([]byte(field.Options), &options); err == nil {
				if len(options) > 0 {
					found := false
					for _, opt := range options {
						if opt == s {
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("字段%s的值不在可选列表中", field.Alias)
					}
				}
			}
		}
	}
	return nil
}

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
		Where(map[string]any{"model_id": modelId}).
		Scan(&modelFields).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	// 需要的字段
	requiredFields := make([]string, 0)
	fieldAlias := make([]string, 0)
	for _, modelField := range modelFields {
		if modelField.IsRequired {
			requiredFields = append(requiredFields, modelField.Alias)
		}
		fieldAlias = append(fieldAlias, modelField.Alias)
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
		if !slices.Contains(fieldAlias, key) {
			delete(dataMap, key)
		}
	}

	// 补充字段
	fieldMap := make(map[string]models.ModelField)
	for _, modelField := range modelFields {
		fieldMap[modelField.Alias] = modelField
		if _, ok := dataMap[modelField.Alias]; !ok {
			defaultVal := models.DefaultValueByType(modelField.Type)
			if defaultVal != nil {
				dataMap[modelField.Alias] = defaultVal
			}
		}
	}

	// 校验字段类型
	for key, value := range dataMap {
		if f, ok := fieldMap[key]; ok {
			if err := validateFieldValue(f, value); err != nil {
				return nil, err
			}
		}
	}

	// 验证唯一性
	if err := verifyUniqueness(modelId, dataMap, 0); err != nil {
		return nil, err
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

	// 校验字段类型
	modelFields := make([]models.ModelField, 0)
	if err := database.DB.Model(&models.ModelField{}).
		Where(map[string]any{"model_id": instance.ModelId}).
		Scan(&modelFields).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	fieldMap := make(map[string]models.ModelField)
	for _, f := range modelFields {
		fieldMap[f.Alias] = f
	}
	for key, value := range dataMap {
		if f, ok := fieldMap[key]; ok {
			if err := validateFieldValue(f, value); err != nil {
				return nil, err
			}
		}
	}

	// 验证唯一性
	if err := verifyUniqueness(instance.ModelId, dataMap, id); err != nil {
		return nil, err
	}
	// 响应数据
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("json数据异常-%s", err.Error())
	}
	return jsonData, nil

}
