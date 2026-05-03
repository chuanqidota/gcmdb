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
	search := listInstance.Search
	field := listInstance.Field
	value := listInstance.Value
	compare := listInstance.Compare
	limit := listInstance.Limit
	offset := listInstance.Offset
	if limit == 0 {
		limit = 10
	}

	if compare != "" && !params.ValidatorCompare(compare) {
		return nil, fmt.Errorf("参数错误-compare")
	}

	baseQuery := func() *gorm.DB {
		return database.DB.Model(&models.Instance{}).Where(map[string]any{"model_id": modelId})
	}

	db := baseQuery()

	// like 比较自动加 % 通配符
	compareValue := value
	if compare == "like" {
		compareValue = "%" + value + "%"
	}

	// 搜索条件：在所有字段中模糊匹配
	addSearchCondition := func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}
		modelFields := make([]models.ModelField, 0)
		if err := database.DB.Model(&models.ModelField{}).
			Where(map[string]any{"model_id": modelId}).
			Scan(&modelFields).Error; err != nil {
			return db
		}
		conditions := make([]string, 0, len(modelFields))
		_params := make([]interface{}, 0, len(modelFields))
		for _, mf := range modelFields {
			conditions = append(conditions, fmt.Sprintf("data->>'$.%s' LIKE ?", mf.Alias))
			_params = append(_params, "%"+search+"%")
		}
		return db.Where(strings.Join(conditions, " OR "), _params...)
	}

	// 检查字段是否已索引
	isFieldIndexed := func(fieldName string) bool {
		var count int64
		database.DB.Model(&models.ModelField{}).
			Where("model_id = ? AND alias = ? AND is_indexed = ?", modelId, fieldName, true).
			Count(&count)
		return count > 0
	}

	switch {
	// 只有 field+value+compare
	case search == "" && field != "" && value != "" && compare != "":
		if !ValidFieldName(field) {
			return nil, fmt.Errorf("非法字段名:%s", field)
		}
		// 精确匹配且字段已索引时，使用 indexed_values 命中索引
		if compare == "=" && isFieldIndexed(field) {
			db = db.Where("FIND_IN_SET(?, indexed_values) > 0", value)
		} else {
			db = db.Where(fmt.Sprintf("data->>'$.%s' %s ?", field, compare), compareValue)
		}

	// search + field+value+compare
	case search != "" && field != "" && value != "" && compare != "":
		if !ValidFieldName(field) {
			return nil, fmt.Errorf("非法字段名:%s", field)
		}
		if compare == "=" && isFieldIndexed(field) {
			db = db.Where("FIND_IN_SET(?, indexed_values) > 0", value)
		} else {
			db = db.Where(fmt.Sprintf("data->>'$.%s' %s ?", field, compare), compareValue)
		}
		db = addSearchCondition(db)

	// 只有 search（无 field/value）
	case search != "" && (field == "" || value == ""):
		db = addSearchCondition(db)

	// 无任何条件：search=="" 且无 field/value（分支1）
	default:
		// 不加额外条件，列出全部
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())
	}

	results := make([]models.Instance, 0)
	if err := db.Limit(limit).Offset(offset).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())
	}

	return &resp.CommonList{
		Count:   count,
		Results: results,
	}, nil
}
