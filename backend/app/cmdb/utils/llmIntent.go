package utils

import (
	"encoding/json"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/resp"
	"gcmdb/pkg/database"
	"strings"

	"gorm.io/gorm"
)

type llmIntent struct{}

var LLMIntent = new(llmIntent)

// LLMQuery 大模型返回的单个查询
type LLMQuery struct {
	Model      string              `json:"model"`
	Conditions []LLMQueryCondition `json:"conditions"`
	Fields     []string            `json:"fields"`
	Limit      int                 `json:"limit"`
}

// LLMQueryCondition 查询条件
type LLMQueryCondition struct {
	Field string      `json:"field"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

// LLMIntentResponse 大模型返回的意图
type LLMIntentResponse struct {
	Thinking string     `json:"thinking"`
	Queries  []LLMQuery `json:"queries"`
	Analysis string     `json:"analysis"`
}

// ParseIntent 解析大模型返回的 JSON 意图
func (li *llmIntent) ParseIntent(content string) (*LLMIntentResponse, error) {
	// 尝试提取 JSON（大模型可能在 JSON 前后添加文字）
	content = strings.TrimSpace(content)

	// 尝试直接解析
	var intent LLMIntentResponse
	if err := json.Unmarshal([]byte(content), &intent); err != nil {
		// 尝试提取 JSON 块
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")
		if start == -1 || end == -1 || end <= start {
			return nil, fmt.Errorf("无法解析大模型返回的格式: %s", content[:minInt(200, len(content))])
		}
		jsonStr := content[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &intent); err != nil {
			return nil, fmt.Errorf("JSON 解析失败: %s", err.Error())
		}
	}

	if len(intent.Queries) == 0 {
		return nil, fmt.Errorf("大模型未生成有效的查询计划")
	}

	return &intent, nil
}

// ExecuteQueries 执行查询计划
func (li *llmIntent) ExecuteQueries(queries []LLMQuery) ([]resp.LLMQueryResult, error) {
	results := make([]resp.LLMQueryResult, 0, len(queries))

	for _, q := range queries {
		if q.Model == "" {
			continue
		}

		// 查询模型
		var model models.Model
		if err := database.DB.Model(&models.Model{}).
			Where("alias = ? AND is_usable = ?", q.Model, true).
			First(&model).Error; err != nil {
			continue // 跳过不存在的模型
		}

		// 构建查询
		db := database.DB.Model(&models.Instance{}).Where("model_id = ?", model.ID)

		// 应用条件
		for _, cond := range q.Conditions {
			if cond.Field == "" || cond.Op == "" {
				continue
			}
			db = li.applyCondition(db, cond)
		}

		// 计数
		var count int64
		db.Count(&count)

		// 限制返回数量
		limit := q.Limit
		if limit <= 0 || limit > 200 {
			limit = 50
		}

		// 查询数据
		var instances []models.Instance
		if err := db.Limit(limit).Scan(&instances).Error; err != nil {
			continue
		}

		// 构建结果
		data := make([]map[string]interface{}, 0, len(instances))
		for _, inst := range instances {
			row := make(map[string]interface{})
			row["id"] = inst.ID

			// 解析 JSON 数据
			var jsonData map[string]interface{}
			if err := json.Unmarshal(inst.Data, &jsonData); err == nil {
				// 如果指定了字段，只返回指定字段
				if len(q.Fields) > 0 {
					filtered := make(map[string]interface{})
					filtered["id"] = inst.ID
					for _, f := range q.Fields {
						if v, ok := jsonData[f]; ok {
							filtered[f] = v
						}
					}
					row = filtered
				} else {
					for k, v := range jsonData {
						row[k] = v
					}
				}
			}
			data = append(data, row)
		}

		results = append(results, resp.LLMQueryResult{
			Model:      model.Alias,
			ModelName:  model.Name,
			ModelAlias: model.Alias,
			Count:      count,
			Data:       data,
		})
	}

	return results, nil
}

// applyCondition 应用单个查询条件到 GORM 查询
func (li *llmIntent) applyCondition(db *gorm.DB, cond LLMQueryCondition) *gorm.DB {
	value := fmt.Sprintf("%v", cond.Value)

	switch cond.Op {
	case "eq":
		return db.Where(fmt.Sprintf("data->>'$.%s' = ?", cond.Field), value)
	case "ne":
		return db.Where(fmt.Sprintf("data->>'$.%s' != ?", cond.Field), value)
	case "gt":
		return db.Where(fmt.Sprintf("data->>'$.%s' > ?", cond.Field), value)
	case "ge":
		return db.Where(fmt.Sprintf("data->>'$.%s' >= ?", cond.Field), value)
	case "lt":
		return db.Where(fmt.Sprintf("data->>'$.%s' < ?", cond.Field), value)
	case "le":
		return db.Where(fmt.Sprintf("data->>'$.%s' <= ?", cond.Field), value)
	case "contains":
		return db.Where(fmt.Sprintf("data->>'$.%s' LIKE ?", cond.Field), "%"+value+"%")
	case "startswith":
		return db.Where(fmt.Sprintf("data->>'$.%s' LIKE ?", cond.Field), value+"%")
	case "endswith":
		return db.Where(fmt.Sprintf("data->>'$.%s' LIKE ?", cond.Field), "%"+value)
	case "search":
		return db.Where("JSON_SEARCH(data, 'one', ?) IS NOT NULL", "%"+value+"%")
	default:
		return db
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FormatResultsForAnalysis 将查询结果格式化为文本供大模型分析
func (li *llmIntent) FormatResultsForAnalysis(results []resp.LLMQueryResult) string {
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("模型: %s (%s), 共 %d 条记录\n", r.ModelName, r.ModelAlias, r.Count))
		sb.WriteString("数据样例:\n")
		for i, d := range r.Data {
			if i >= 10 {
				sb.WriteString(fmt.Sprintf("... 共 %d 条\n", len(r.Data)))
				break
			}
			jsonBytes, _ := json.Marshal(d)
			sb.WriteString(string(jsonBytes) + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
