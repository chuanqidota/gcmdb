package utils

import (
	"errors"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/utils"
	"gcmdb/app/openapi/params"
	"gcmdb/app/openapi/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
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
		First(&modelId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("未找到模型:%s", modelAlias)
		}
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
			logger.Error(fmt.Sprintf("创建实例关联失败-%s", err.Error()))
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
			logger.Error(fmt.Sprintf("删除实例关联失败-%s", err.Error()))
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
			logger.Error(fmt.Sprintf("删除实例关联失败-%s", err.Error()))
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
		First(&searchDirectSql).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未找到查询:%s", uuid)
		}
		return nil, fmt.Errorf("查询失败:%s", err.Error())
	}
	sql := searchDirectSql.Sql
	if sql == "" {
		return nil, fmt.Errorf("sql语句为空")
	}
	upperSQL := strings.ToUpper(strings.TrimSpace(sql))
	if !strings.HasPrefix(upperSQL, "SELECT") {
		return nil, fmt.Errorf("仅支持SELECT查询语句")
	}
	var result interface{}
	if err := database.DB.Raw(sql).Limit(1000).Scan(&result).Error; err != nil {
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

	query := database.DB.Table("instance i").
		Select("i.id, i.model_id, i.model_alias, m.name as model_name, i.data").
		Joins("LEFT JOIN model m ON i.model_id = m.id").
		Where("m.is_usable = ?", true).
		Where("JSON_SEARCH(i.data, 'one', ?) IS NOT NULL", "%"+search+"%")

	if modelAlias != "" {
		modelAliasList := strings.Split(modelAlias, ",")
		query = query.Where("m.alias IN ?", modelAliasList)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, nil, fmt.Errorf("查询失败-%s", err.Error())
	}

	var result []resp.FulltextInstance
	if err := query.Limit(limit).Offset(offset).Scan(&result).Error; err != nil {
		return 0, nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	return count, result, nil
}

// SearchInstance
//
//	@Description: 搜索实例信息
//	@receiver i
//	@param body
//	@return []models.Instance
//	@return error
func (i *instance) SearchInstance(body params.SearchInstance) (int64, any, error) {
	limit := body.Condition.Limit
	offset := body.Condition.Offset
	if limit == 0 {
		limit = 10
	}

	validateFields := func(fields map[string]interface{}) error {
		for f := range fields {
			if !utils.ValidFieldName(f) {
				return fmt.Errorf("非法字段名:%s", f)
			}
		}
		return nil
	}
	jsonExtract := func(field, expr string, args ...any) (string, []any) {
		return fmt.Sprintf("data->>'$.%s' "+expr, field), args
	}

	query := database.DB.Model(&models.Instance{}).Where(map[string]any{"model_alias": body.Model})

	// 指定查询字段
	if len(body.Fields) > 0 {
		baseCols := map[string]bool{"id": true, "model_id": true, "model_alias": true}
		parts := []string{"id", "model_id", "model_alias"}
		for _, f := range body.Fields {
			if !utils.ValidFieldName(f) {
				return 0, nil, fmt.Errorf("非法字段名:%s", f)
			}
			if baseCols[f] {
				continue
			}
			parts = append(parts, fmt.Sprintf("data->>'$.%s' as `%s`", f, f))
		}
		query = query.Select(strings.Join(parts, ", "))
	}

	// where 条件
	compares := map[string]string{
		"eq": "=", "ne": "!=", "gt": ">", "ge": ">=",
		"lt": "<", "le": "<=",
	}
	likes := map[string]func(string) string{
		"startswith": func(s string) string { return s + "%" },
		"endswith":   func(s string) string { return "%" + s },
		"contains":   func(s string) string { return "%" + s + "%" },
	}

	for _, cond := range body.Condition.Where {
		for action, value := range cond {
			switch {
			case action == "search":
				if s, ok := value.(string); ok {
					query = query.Where("JSON_SEARCH(data, 'one', ?) IS NOT NULL", "%"+s+"%")
				}
			case action == "or":
				ors, ok := value.([]interface{})
				if !ok {
					continue
				}
				for _, orCond := range ors {
					condMap, ok := orCond.(map[string]interface{})
					if !ok {
						continue
					}
					orQ := database.DB.Model(&models.Instance{}).Where("model_alias = ?", body.Model)
					for orAct, orVal := range condMap {
						op, ok := compares[orAct]
						if !ok {
							continue
						}
						if mv, ok := orVal.(map[string]interface{}); ok {
							for f, v := range mv {
								if !utils.ValidFieldName(f) {
									return 0, nil, fmt.Errorf("非法字段名:%s", f)
								}
								sql, args := jsonExtract(f, op+" ?", v)
								orQ = orQ.Where(sql, args...)
							}
						}
					}
					query = query.Or(orQ)
				}
			case compares[action] != "":
				if mv, ok := value.(map[string]interface{}); ok {
					if err := validateFields(mv); err != nil {
						return 0, nil, err
					}
					for f, v := range mv {
						sql, args := jsonExtract(f, compares[action]+" ?", v)
						query = query.Where(sql, args...)
					}
				}
			case action == "in":
				if mv, ok := value.(map[string]interface{}); ok {
					if err := validateFields(mv); err != nil {
						return 0, nil, err
					}
					for f, v := range mv {
						sql, args := jsonExtract(f, "IN ?", v)
						query = query.Where(sql, args...)
					}
				}
			case likes[action] != nil:
				if mv, ok := value.(map[string]interface{}); ok {
					if err := validateFields(mv); err != nil {
						return 0, nil, err
					}
					wrap := likes[action]
					for f, v := range mv {
						if s, ok := v.(string); ok {
							sql, args := jsonExtract(f, "LIKE ?", wrap(s))
							query = query.Where(sql, args...)
						}
					}
				}
			}
		}
	}

	// 排序
	allowedOrders := map[string]bool{
		"id": true, "created_at": true, "updated_at": true, "model_id": true,
	}
	for _, item := range body.Condition.Order {
		if item == "" {
			continue
		}
		col, dir := item, "ASC"
		if strings.HasPrefix(item, "-") {
			col, dir = item[1:], "DESC"
		} else if parts := strings.Fields(item); len(parts) > 1 {
			dir = strings.ToUpper(parts[1])
			if dir != "ASC" && dir != "DESC" {
				return 0, nil, fmt.Errorf("非法排序方向:%s", dir)
			}
		}
		if !allowedOrders[col] {
			return 0, nil, fmt.Errorf("不允许按字段排序:%s", col)
		}
		query = query.Order(col + " " + dir)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, nil, fmt.Errorf("查询出错-%s", err.Error())
	}

	if len(body.Fields) > 0 {
		results := make([]map[string]interface{}, 0)
		if err := query.Offset(offset).Limit(limit).Scan(&results).Error; err != nil {
			return 0, nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		return count, results, nil
	}

	instances := make([]models.Instance, 0)
	if err := query.Offset(offset).Limit(limit).Scan(&instances).Error; err != nil {
		return 0, nil, fmt.Errorf("查询出错-%s", err.Error())
	}
	return count, instances, nil
}

// TargetInstance
//
//	@Description: 根据源目标实例 查询 目标实例信息
//	@receiver i
//	@param sourceId 源实例id
//	@param targetModel 目标模型别名
//	@return []models.Instance
//	@return error
func (i *instance) TargetInstance(sourceId uint, targetModel string) ([]models.Instance, error) {
	return i.queryRelatedInstance(sourceId, targetModel, "target")
}

// SourceInstance
//
//	@Description: 根据目标实例 查询 源实例信息
//	@receiver i
//	@param targetId 目标实例id
//	@param SourceModel 源模型别名
//	@return []models.Instance
//	@return error
func (i *instance) SourceInstance(targetId uint, SourceModel string) ([]models.Instance, error) {
	return i.queryRelatedInstance(targetId, SourceModel, "source")
}

// queryRelatedInstance 通用关联实例查询, direction: "source" 或 "target"
func (i *instance) queryRelatedInstance(id uint, modelAlias, direction string) ([]models.Instance, error) {
	instances := make([]models.Instance, 0)
	// direction="target": 查目标实例 (join on target_id, where source_id)
	// direction="source": 查源实例 (join on source_id, where target_id)
	joinCol := "ir.target_id"
	whereCol := "ir.source_id"
	modelCol := "ir.target_model_id"
	if direction == "source" {
		joinCol = "ir.source_id"
		whereCol = "ir.target_id"
		modelCol = "ir.source_model_id"
	}

	query := database.DB.Table("instance i").
		Select("i.*").
		Joins("INNER JOIN instance_relation ir ON i.id = "+joinCol).
		Where(whereCol+" = ?", id)

	if modelAlias != "" {
		var modelId uint
		if err := database.DB.Model(&models.Model{}).
			Select("id").
			Where(map[string]any{"alias": modelAlias}).
			First(&modelId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("未找到模型:%s", modelAlias)
			}
			return nil, fmt.Errorf("查询出错-%s", err.Error())
		}
		query = query.Where(modelCol+" = ?", modelId)
	}

	if err := query.Scan(&instances).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())
	}
	return instances, nil
}
