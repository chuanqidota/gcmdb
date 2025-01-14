package utils

import (
	"errors"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/utils"
	"gcmdb/app/openapi/params"
	"gcmdb/app/openapi/resp"
	"gcmdb/pkg/database"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type instance struct {
}

var Instance = new(instance)

// CreateInstance
//
//	@Description: 创建实例
//	@receiver i
//	@param modelAlias 模型英文名
//	@param data
//	@return error
func (i *instance) CreateInstance(modelAlias string, data datatypes.JSON) error {
	var modelId uint
	if err := database.DB.Model(&models.Model{}).
		Select("id").
		Where(map[string]any{"alias": modelAlias}).
		Scan(&modelId).Error; err != nil {
		return fmt.Errorf("查询模型失败：%s", err.Error())
	}
	verifyData, err := utils.Verify.VerifyCreateInstance(modelId, data)
	if err != nil {
		return err
	}
	// 创建实例数据
	createInstance := models.Instance{
		ModelId:    modelId,
		ModelAlias: modelAlias,
		Data:       verifyData,
	}

	if err := database.DB.Model(&models.Instance{}).
		Create(&createInstance).Error; err != nil {
		return fmt.Errorf("创建实例失败:%s", err.Error())
	}
	// 异步创建实例关系
	go func() {
		if err := utils.InstanceRelation.CreateInstance(modelId, createInstance.ID); err != nil {
			fmt.Printf("创建实例关联失败-%s", err.Error())
		}
	}()
	return nil

}

// UpdateInstance
//
//	@Description: 更新实例
//	@receiver i
//	@param id 实例ID
//	@param data 更新数据
//	@return error
func (i *instance) UpdateInstance(id uint, data datatypes.JSON) error {
	verifyData, err := utils.Verify.VerifyUpdateInstance(id, data)
	if err != nil {
		return fmt.Errorf("参数校验失败:%s", err.Error())
	}
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id}).
		Update("data", verifyData).Error; err != nil {
		return fmt.Errorf("更新失败-%s", err.Error())
	}
	return nil

}

// DeleteInstance
//
//	@Description: 删除实例
//	@receiver i
//	@param modelAlias 模型英文名
//	@param id
//	@return error
func (i *instance) DeleteInstance(modelAlias string, id uint) error {
	// 删除实例
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id, "model_alias": modelAlias}).
		Delete(&models.Instance{}).Error; err != nil {
		return fmt.Errorf("删除失败-%s", err.Error())
	}

	// 删除实例关联
	go func() {
		if err := utils.InstanceRelation.DeleteInstance(id); err != nil {
			fmt.Printf("删除实例关联失败-%s", err.Error())
		}
	}()
	return nil
}

// MulDeleteInstance
//
//	@Description: 批量删除实例
//	@receiver i
//	@param modelAlias 模型英文名
//	@param ids 实例ID列表
//	@return error
func (i *instance) MulDeleteInstance(modelAlias string, ids []uint) error {
	// 删除实例
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": ids, "model_alias": modelAlias}).
		Delete(&models.Instance{}).Error; err != nil {
		return fmt.Errorf("删除失败-%s", err.Error())
	}

	// 删除实例关联
	go func() {
		if err := utils.InstanceRelation.MulDeleteInstance(ids); err != nil {
			fmt.Printf("删除实例关联失败-%s", err.Error())
		}
	}()
	return nil
}

// DirectSearch
//
//	@Description: 直接搜索sql
//	@receiver i
//	@param uuid
//	@return error
func (i *instance) DirectSearch(uuid string) (any, error) {
	var searchDirectSql models.SearchDirectSql
	if err := database.DB.Model(&models.SearchDirectSql{}).
		Where(map[string]any{"uuid": uuid}).
		First(&searchDirectSql).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询失败:%s",err.Error())
	}
	sql := searchDirectSql.Sql
	if sql == "" {
		return nil, fmt.Errorf("sql语句为空")
	}
	var result interface{}
	if err := database.DB.Raw(sql).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	return result, nil
}

// FulltextInstance
//
//	@Description: 全文搜索
//	@receiver i
//	@param body
//	@return int64
//	@return []resp.FulltextInstance
//	@return error
func (i *instance) FulltextInstance(body params.FulltextInstance) (int64, []resp.FulltextInstance, error) {
	search := body.Search
	limit := body.Limit
	offset := body.Offset
	if limit == 0 {
		limit = 10
	}
	modelAlias := body.ModelAlias
	modelAliasList := strings.Split(modelAlias, ",")

	sql := fmt.Sprintf(`
			SELECT
			    i.id,
			    i.model_id,
			    i.model_alias,
			    m.name as model_name
			    cast(i.data as JSON) data
			FROM
				instance i
			LEFT JOIN 
				model m
			ON
				i.model_id = m.id
			WHERE
			    m.is_usable = true 
			AND
				JSON_SEARCH(i.data,"one","%%%s%%")
    `, search)
	if modelAliasList != nil {
		for _, modelAlias := range modelAliasList {
			sql += fmt.Sprintf(`
				AND
					m.model_alias = '%s'
			`, modelAlias)
		}
	}
	sql += fmt.Sprintf(`
				LIMIT
					%d
				OFFSET
                    %d
			`, limit, offset)
	var result []resp.FulltextInstance
	qs := database.DB.Raw(sql).Scan(&result)
	if qs.Error != nil {
		return 0, nil, fmt.Errorf("查询失败-%s", qs.Error.Error())
	}
	count := qs.RowsAffected
	return count, result, nil
}

func (i *instance) SearchInstance(body params.SearchInstance) {
	model := body.Model
	fields := body.Fields
	// children := body.Children
	__condition := body.Condition
	limit := __condition.Limit
	offset := __condition.Offset
	order := __condition.Order
	where := __condition.Where

	// 指定查询
	query := database.DB.Model(&models.Instance{}).Where(map[string]any{"model_alias": model})

	selectFields := ""
	for i, field := range fields {
		if i > 0 {
			selectFields += ", "
		}
		selectFields += fmt.Sprintf("JSON_EXTRACT(data, '$.%s') as %s", field, field)
	}
	if selectFields != "" {
		query = query.Select(selectFields)
	}

	var actionMap = map[string]string{
		"eq":         "=",
		"ne":         "!=",
		"contains":   "LIKE",
		"startswith": "LIKE",
		"endswith":   "LIKE",
		"gt":         ">",
		"ge":         ">=",
		"lt":         "<",
		"le":         "<=",
		"in":         "IN",
	}


	innerConditionSql := func(action string, value any) string {
		var conditions []string
	
		switch action {
		case "search":
			// 模糊查询
			conditions = append(conditions, fmt.Sprintf("JSON_SEARCH(data, 'one', '%s') IS NOT NULL", "%"+value.(string)+"%"))
	
		case "eq", "ne", "gt", "lt", "le", "ge":
			sqlAction := actionMap[action]
			for field, val := range value.(map[string]interface{}) {
				conditions = append(conditions, fmt.Sprintf("JSON_EXTRACT(data, '$.%s') %s '%v'", field, sqlAction, val))
			}
	
		case "in":
			sqlAction := actionMap[action]
			for field, val := range value.(map[string]interface{}) {
				// 将值转换为逗号分隔的字符串
				values := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(val)), "','"), "[]")
				conditions = append(conditions, fmt.Sprintf("JSON_EXTRACT(data, '$.%s') %s ('%s')", field, sqlAction, values))
			}
	
		case "startswith":
			for field, val := range value.(map[string]interface{}) {
				conditions = append(conditions, fmt.Sprintf("JSON_EXTRACT(data, '$.%s') LIKE '%s%%'", field, val.(string)))
			}
	
		case "endswith":
			for field, val := range value.(map[string]interface{}) {
				conditions = append(conditions, fmt.Sprintf("JSON_EXTRACT(data, '$.%s') LIKE '%%%s'", field, val.(string)))
			}
	
		default:
			return ""
		}
	
		return strings.Join(conditions, " AND ")

	}
	
	querySql := " 1=1 "
	for _, cond := range where {
		for action, value := range cond {
			switch action {
			case "search":
				// 模糊查询
				query = query.Where("JSON_SEARCH(data, 'one', ?) IS NOT NULL", "%"+value.(string)+"%")
			case "eq", "ne", "gt", "lt", "le", "in", "startswith", "endswith":
				querySql = innerConditionSql(action, value)
				query = query.Where(querySql)
			case "or":
				orConditions := value.([]map[string]interface{})
				// 外部是OR
				for _, cond := range orConditions {
					// 内部是AND
					andQuery := make([]string,0)
					for or_action,or_value := range cond {
						_query := innerConditionSql(or_action,or_value)
						andQuery = append(andQuery,_query)
					}
					query.Or(strings.Join(andQuery, " AND "))
				}
			}
		}
	
		// 排序
		for _, item := range order {
			query = query.Order(item)
		}

		// 分页
		query = query.Offset(offset).Limit(limit)

	}
}

func (i *instance) TargetInstance(sourceId int64, targetModel string) ([]models.Instance, error) {
	var modelId int64
	if err := database.DB.Model(&models.Model{}).
		Select("id").
		Where(map[string]any{"alias": targetModel}).
		Scan(&modelId).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())
	}
	query := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"source_id": sourceId})
	if modelId != 0 {
		query = query.Where(map[string]any{"target_model_id": modelId})
	}
	targetIds := make([]int64, 0)
	if err := query.Select("target_id").Scan(&targetIds).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	instances := make([]models.Instance, 0)
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": targetIds}).
		Scan(&instances).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())

	}
	return instances, nil
}

func (i *instance) SourceInstance(targetId int64, SourceModel string) ([]models.Instance, error) {
	var modelId int64
	if err := database.DB.Model(&models.Model{}).
		Select("id").
		Where(map[string]any{"alias": SourceModel}).
		Scan(&modelId).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())
	}
	query := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"target_id": targetId})
	if modelId != 0 {
		query = query.Where(map[string]any{"source_model_id": modelId})
	}
	sourceIds := make([]int64, 0)
	if err := query.Select("source_id").Scan(&sourceIds).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	instances := make([]models.Instance, 0)
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": sourceIds}).
		Scan(&instances).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())

	}
	return instances, nil
}
