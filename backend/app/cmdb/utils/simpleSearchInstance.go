package utils

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/app/cmdb/resp"
	"gcmdb/pkg/database"
	"strings"

	"gorm.io/gorm"
)

// SimpleSearchInstance
//
//	@Description: 页面指定条件的搜索
//	@param modelId
//	@param listInstance
//	@return *resp.CommonList
//	@return error
func SimpleSearchInstance(modelId uint, listInstance params.ListInstance) (*resp.CommonList, error) {
	// 参数
	search := listInstance.Search
	field := listInstance.Field
	value := listInstance.Value
	compare := listInstance.Compare
	limit := listInstance.Limit
	offset := listInstance.Offset
	if limit == 0 {
		limit = 10
	}

	if compare != "" {
		if !params.ValidatorCompare(compare) {
			return nil, fmt.Errorf("参数错误-compare")
		}
	}

	// 校验字段名安全性
	validFieldName := func(name string) bool {
		for _, c := range name {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
				return false
			}
		}
		return name != ""
	}

	// 匿名函数-拼接条件
	whereClause := func(fields []string) (string, []interface{}) {
		var conditions []string
		_params := make([]interface{}, len(fields))
		for i, _field := range fields {
			conditions = append(conditions, fmt.Sprintf("data->'$.%s' LIKE ?", _field))
			_params[i] = fmt.Sprintf("%%%s%%", search)
		}
		query := strings.Join(conditions, " OR ")
		return query, _params
	}

	baseQuery := func() *gorm.DB {
		return database.DB.Model(&models.Instance{}).Where(map[string]any{"model_id": modelId})
	}

	results := make([]models.Instance, 0)
	if search == "" || (field == "" && value == "") {
		// 第一种情况,无任何条件
		var count int64
		if err := baseQuery().Count(&count).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		if err := baseQuery().Limit(limit).Offset(offset).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		return &resp.CommonList{
			Count:   count,
			Results: results,
		}, nil
	} else if search != "" && (field == "" || value == "") {
		// 第二种情况,只有search

		// 获取模型字段
		modelFields := make([]models.ModelField, 0)
		if err := database.DB.Model(&models.ModelField{}).
			Where(map[string]any{"model_id": modelId}).
			Scan(&modelFields).Error; err != nil {
			return nil, err
		}
		fields := make([]string, 0)
		for _, modelField := range modelFields {
			fields = append(fields, modelField.Alias)
		}
		// 拼接条件
		query, _params := whereClause(fields)
		db := baseQuery().Where(query, _params...)

		var count int64
		if err := db.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		if err := db.Limit(limit).Offset(offset).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		return &resp.CommonList{
			Count:   count,
			Results: results,
		}, nil
	} else if search == "" && (field != "" && value != "" && compare != "") {
		// 第三种情况,只有field和value
		if !validFieldName(field) {
			return nil, fmt.Errorf("非法字段名:%s", field)
		}
		condition := fmt.Sprintf("data->'$.%s' %s ?", field, compare)
		db := baseQuery().Where(condition, value)

		var count int64
		if err := db.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		if err := db.Limit(limit).Offset(offset).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		return &resp.CommonList{
			Count:   count,
			Results: results,
		}, nil
	} else if search != "" && (field != "" && value != "" && compare != "") {
		// 第四种情况,search和field和value
		if !validFieldName(field) {
			return nil, fmt.Errorf("非法字段名:%s", field)
		}
		// 先精准搜索
		condition := fmt.Sprintf("data->'$.%s' %s ?", field, compare)
		db := baseQuery().Where(condition, value)
		// 再模糊搜索
		query, _params := whereClause([]string{field})
		db = db.Where(query, _params...)

		var count int64
		if err := db.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		if err := db.Limit(limit).Offset(offset).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		return &resp.CommonList{
			Count:   count,
			Results: results,
		}, nil
	}
	return nil, fmt.Errorf("搜索不支持")
}
