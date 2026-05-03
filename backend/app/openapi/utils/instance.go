package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/utils"
	"gcmdb/app/openapi/params"
	"gcmdb/app/openapi/resp"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	pkgUtils "gcmdb/pkg/utils"
	"strconv"
	"strings"
	"time"

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
		ModelId: modelId,
		Data:    verifyData,
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
	// 通过 model_alias 查找 model_id
	var modelId uint
	if err := database.DB.Model(&models.Model{}).Select("id").Where("alias = ?", modelAlias).First(&modelId).Error; err != nil {
		return fmt.Errorf("未找到模型:%s", modelAlias)
	}
	// 删除实例
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": id, "model_id": modelId}).
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
	// 通过 model_alias 查找 model_id
	var modelId uint
	if err := database.DB.Model(&models.Model{}).Select("id").Where("alias = ?", modelAlias).First(&modelId).Error; err != nil {
		return fmt.Errorf("未找到模型:%s", modelAlias)
	}
	// 删除实例
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": ids, "model_id": modelId}).
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
	if err := pkgUtils.ValidateSelectSQL(sql); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := database.DB.Raw(sql).Limit(1000).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	// MySQL JSON 列扫描后可能为 []uint8 或 string，需反序列化为 map/slice
	for _, row := range result {
		for k, v := range row {
			switch val := v.(type) {
			case []uint8:
				var parsed interface{}
				if json.Unmarshal(val, &parsed) == nil {
					row[k] = parsed
				}
			case string:
				if len(val) > 0 && (val[0] == '{' || val[0] == '[') {
					var parsed interface{}
					if json.Unmarshal([]byte(val), &parsed) == nil {
						row[k] = parsed
					}
				}
			}
		}
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
		Select("i.id, i.model_id, m.name as model_name, m.alias as model_alias, i.data, i.created_at, i.updated_at").
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

	type fulltextRow struct {
		ID         uint
		ModelId    uint
		ModelName  string
		ModelAlias string
		Data       datatypes.JSON
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	var rows []fulltextRow
	if err := query.Limit(limit).Offset(offset).Scan(&rows).Error; err != nil {
		return 0, nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	result := make([]resp.FulltextInstance, 0, len(rows))
	for _, r := range rows {
		result = append(result, resp.FulltextInstance{
			ID:         r.ID,
			ModelId:    r.ModelId,
			ModelName:  r.ModelName,
			ModelAlias: r.ModelAlias,
			Data:       r.Data,
			CreatedAt:  r.CreatedAt.Format(time.DateTime),
			UpdatedAt:  r.UpdatedAt.Format(time.DateTime),
		})
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

	// 通过 model_alias 查找 model_id
	var modelId uint
	if err := database.DB.Model(&models.Model{}).Select("id").Where("alias = ?", body.Model).First(&modelId).Error; err != nil {
		return 0, nil, fmt.Errorf("未找到模型:%s", body.Model)
	}
	query := database.DB.Model(&models.Instance{}).Where(map[string]any{"model_id": modelId})

	// 指定查询字段
	if len(body.Fields) > 0 {
		baseCols := map[string]bool{"id": true, "model_id": true}
		parts := []string{"id", "model_id"}
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
					orQ := database.DB.Model(&models.Instance{}).Where("model_id = ?", modelId)
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

	type searchRow struct {
		ID        uint
		ModelId   uint
		Data      datatypes.JSON
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	var rows []searchRow
	if err := query.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return 0, nil, fmt.Errorf("查询出错-%s", err.Error())
	}
	// 查询模型名称
	var m models.Model
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": modelId}).First(&m)
	items := make([]resp.SearchInstanceItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, resp.SearchInstanceItem{
			ID:         r.ID,
			ModelId:    r.ModelId,
			ModelName:  m.Name,
			ModelAlias: m.Alias,
			Data:       r.Data,
			CreatedAt:  r.CreatedAt.Format(time.DateTime),
			UpdatedAt:  r.UpdatedAt.Format(time.DateTime),
		})
	}
	return count, items, nil
}

// TargetInstance
//
//	@Description: 根据源目标实例 查询 目标实例信息
//	@receiver i
//	@param sourceId 源实例id
//	@param targetModel 目标模型别名
//	@return []resp.RelatedInstance
//	@return error
func (i *instance) TargetInstance(sourceId uint, targetModel string) ([]resp.RelatedInstance, error) {
	return i.queryRelatedInstance(sourceId, targetModel, "target")
}

// SourceInstance
//
//	@Description: 根据目标实例 查询 源实例信息
//	@receiver i
//	@param targetId 目标实例id
//	@param SourceModel 源模型别名
//	@return []resp.RelatedInstance
//	@return error
func (i *instance) SourceInstance(targetId uint, SourceModel string) ([]resp.RelatedInstance, error) {
	return i.queryRelatedInstance(targetId, SourceModel, "source")
}

// queryRelatedInstance 通用关联实例查询, direction: "source" 或 "target"
func (i *instance) queryRelatedInstance(id uint, modelAlias, direction string) ([]resp.RelatedInstance, error) {
	type relatedRow struct {
		ID         uint
		ModelId    uint
		ModelName  string
		ModelAlias string
		Data       datatypes.JSON
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	var rows []relatedRow
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
		Select("i.id, i.model_id, m.name as model_name, m.alias as model_alias, i.data, i.created_at, i.updated_at").
		Joins("INNER JOIN instance_relation ir ON i.id = "+joinCol).
		Joins("INNER JOIN model m ON m.id = i.model_id").
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

	if err := query.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("查询出错-%s", err.Error())
	}
	result := make([]resp.RelatedInstance, 0, len(rows))
	for _, r := range rows {
		result = append(result, resp.RelatedInstance{
			ID:         r.ID,
			ModelId:    r.ModelId,
			ModelName:  r.ModelName,
			ModelAlias: r.ModelAlias,
			Data:       r.Data,
			CreatedAt:  r.CreatedAt.Format(time.DateTime),
			UpdatedAt:  r.UpdatedAt.Format(time.DateTime),
		})
	}
	return result, nil
}

// DetailInstance
//
//	@Description: 查询单个实例详情
//	@receiver i
//	@param idStr 实例ID字符串
//	@return *resp.DetailInstance
//	@return error
func (i *instance) DetailInstance(idStr string) (*resp.DetailInstance, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的实例ID:%s", idStr)
	}
	var inst models.Instance
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": uint(id)}).
		First(&inst).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("实例不存在:%d", id)
		}
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	var m models.Model
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": inst.ModelId}).First(&m)
	return &resp.DetailInstance{
		ID:         inst.ID,
		ModelId:    inst.ModelId,
		ModelName:  m.Name,
		ModelAlias: m.Alias,
		Data:       inst.Data,
		CreatedAt:  time.Time(inst.CreatedAt).Format(time.DateTime),
		UpdatedAt:  time.Time(inst.UpdatedAt).Format(time.DateTime),
	}, nil
}

// TopologyInstance
//
//	@Description: 查询实例拓扑关系（上下游）
//	@receiver i
//	@param idStr 实例ID字符串
//	@return *resp.TopologyInfo
//	@return error
func (i *instance) TopologyInstance(idStr string, modelAlias string) (*resp.TopologyInfo, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的实例ID:%s", idStr)
	}
	// 查询实例
	var inst models.Instance
	if err := database.DB.Model(&models.Instance{}).
		Where(map[string]any{"id": uint(id)}).
		First(&inst).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("实例不存在:%d", id)
		}
		return nil, fmt.Errorf("查询失败-%s", err.Error())
	}
	var m models.Model
	database.DB.Model(&models.Model{}).Where(map[string]any{"id": inst.ModelId}).First(&m)

	detail := resp.DetailInstance{
		ID:         inst.ID,
		ModelId:    inst.ModelId,
		ModelName:  m.Name,
		ModelAlias: m.Alias,
		Data:       inst.Data,
		CreatedAt:  time.Time(inst.CreatedAt).Format(time.DateTime),
		UpdatedAt:  time.Time(inst.UpdatedAt).Format(time.DateTime),
	}

	// 查询下游（当前实例作为 source）
	downstream, err := i.queryTopologyRelations(uint(id), "target", modelAlias)
	if err != nil {
		return nil, err
	}
	// 查询上游（当前实例作为 target）
	upstream, err := i.queryTopologyRelations(uint(id), "source", modelAlias)
	if err != nil {
		return nil, err
	}

	return &resp.TopologyInfo{
		Instance:   detail,
		Upstream:   upstream,
		Downstream: downstream,
	}, nil
}

// queryTopologyRelations 查询拓扑关系, direction="target"查下游, "source"查上游
func (i *instance) queryTopologyRelations(instanceId uint, direction string, modelAlias string) ([]resp.TopologyRelation, error) {
	type relRow struct {
		InstanceId uint
		ModelId    uint
		ModelName  string
		ModelAlias string
		TypeId     uint
		S2T        string
		T2S        string
		Data       datatypes.JSON
	}

	// 可选模型过滤
	var filterModelId uint
	if modelAlias != "" {
		if err := database.DB.Model(&models.Model{}).Select("id").Where("alias = ?", modelAlias).First(&filterModelId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("未找到模型:%s", modelAlias)
			}
			return nil, fmt.Errorf("查询模型失败-%s", err.Error())
		}
	}

	var rows []relRow
	// direction="target": 当前实例是source，查target（下游）
	// direction="source": 当前实例是target，查source（上游）
	if direction == "target" {
		q := database.DB.Table("instance_relation ir").
			Select("ir.target_id as instance_id, m.id as model_id, m.name as model_name, m.alias as model_alias, mrt.id as type_id, mrt.s2t, mrt.t2s, i.data").
			Joins("INNER JOIN instance i ON i.id = ir.target_id").
			Joins("INNER JOIN model m ON m.id = ir.target_model_id").
			Joins("LEFT JOIN model_relation_type mrt ON mrt.id = (SELECT type_id FROM model_relation WHERE (source_id = ir.source_model_id AND target_id = ir.target_model_id) OR (source_id = ir.target_model_id AND target_id = ir.source_model_id) LIMIT 1)").
			Where("ir.source_id = ?", instanceId)
		if filterModelId > 0 {
			q = q.Where("ir.target_model_id = ?", filterModelId)
		}
		if err := q.Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("查询下游关系失败-%s", err.Error())
		}
	} else {
		q := database.DB.Table("instance_relation ir").
			Select("ir.source_id as instance_id, m.id as model_id, m.name as model_name, m.alias as model_alias, mrt.id as type_id, mrt.s2t, mrt.t2s, i.data").
			Joins("INNER JOIN instance i ON i.id = ir.source_id").
			Joins("INNER JOIN model m ON m.id = ir.source_model_id").
			Joins("LEFT JOIN model_relation_type mrt ON mrt.id = (SELECT type_id FROM model_relation WHERE (source_id = ir.source_model_id AND target_id = ir.target_model_id) OR (source_id = ir.target_model_id AND target_id = ir.source_model_id) LIMIT 1)").
			Where("ir.target_id = ?", instanceId)
		if filterModelId > 0 {
			q = q.Where("ir.source_model_id = ?", filterModelId)
		}
		if err := q.Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("查询上游关系失败-%s", err.Error())
		}
	}

	result := make([]resp.TopologyRelation, 0, len(rows))
	for _, r := range rows {
		relType := r.S2T
		if direction == "source" {
			relType = r.T2S
		}
		result = append(result, resp.TopologyRelation{
			InstanceId:   r.InstanceId,
			ModelId:      r.ModelId,
			ModelName:    r.ModelName,
			ModelAlias:   r.ModelAlias,
			RelationType: relType,
			Data:         r.Data,
		})
	}
	return result, nil
}
