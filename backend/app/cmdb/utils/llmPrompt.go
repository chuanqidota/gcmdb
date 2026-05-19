package utils

import (
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"strings"
)

type llmPrompt struct{}

var LLMPrompt = new(llmPrompt)

// BuildSystemPrompt 构建 system prompt，注入模型元数据
func (lp *llmPrompt) BuildSystemPrompt(modelFilter []string) (string, error) {
	// 查询所有可用模型
	var modelsList []models.Model
	db := database.DB.Model(&models.Model{}).Where("is_usable = ?", true)
	if len(modelFilter) > 0 {
		db = db.Where("alias IN ?", modelFilter)
	}
	if err := db.Order("`order` asc").Scan(&modelsList).Error; err != nil {
		return "", fmt.Errorf("查询模型失败: %s", err.Error())
	}

	if len(modelsList) == 0 {
		return "", fmt.Errorf("未找到可用模型")
	}

	var sb strings.Builder
	sb.WriteString("你是 CMDB 智能检索助手。用户的 CMDB 系统包含以下模型和数据结构：\n\n")

	// 构建模型列表
	sb.WriteString("【模型列表】\n")
	modelIdMap := make(map[uint]string) // model_id -> alias
	for _, m := range modelsList {
		modelIdMap[m.ID] = m.Alias
		sb.WriteString(fmt.Sprintf("- %s (%s): ", m.Name, m.Alias))

		// 查询模型字段
		var fields []models.ModelField
		database.DB.Model(&models.ModelField{}).
			Where("model_id = ?", m.ID).
			Order("`order` asc").
			Scan(&fields)

		fieldDescs := make([]string, 0, len(fields))
		for _, f := range fields {
			desc := fmt.Sprintf("%s(%s):%s", f.Name, f.Alias, f.Type)
			if f.IsRequired {
				desc += ",必填"
			}
			if f.IsIndexed {
				desc += ",已索引"
			}
			fieldDescs = append(fieldDescs, desc)
		}
		sb.WriteString(strings.Join(fieldDescs, "; "))
		sb.WriteString("\n")
	}

	// 构建模型关系
	sb.WriteString("\n【模型关系】\n")
	var relations []models.ModelRelation
	database.DB.Model(&models.ModelRelation{}).Scan(&relations)

	var relationTypes []models.ModelRelationType
	database.DB.Model(&models.ModelRelationType{}).Scan(&relationTypes)
	typeMap := make(map[uint]models.ModelRelationType)
	for _, rt := range relationTypes {
		typeMap[rt.ID] = rt
	}

	if len(relations) == 0 {
		sb.WriteString("暂无模型关系\n")
	} else {
		for _, rel := range relations {
			sourceAlias := modelIdMap[rel.SourceId]
			targetAlias := modelIdMap[rel.TargetId]
			if sourceAlias == "" || targetAlias == "" {
				continue
			}
			rt := typeMap[rel.TypeId]
			sb.WriteString(fmt.Sprintf("- %s → %s (%s / %s)\n", sourceAlias, targetAlias, rt.S2T, rt.T2S))
		}
	}

	// 任务说明
	sb.WriteString(`
【你的任务】
1. 理解用户自然语言查询意图
2. 返回 JSON 格式的查询计划，格式如下：
{
  "thinking": "你的推理过程",
  "queries": [
    {
      "model": "模型别名",
      "conditions": [
        {"field": "字段别名", "op": "操作符", "value": "值"}
      ],
      "fields": ["字段别名1", "字段别名2"],
      "limit": 50
    }
  ],
  "analysis": "分析请求（如果需要对结果做统计分析则填写，否则留空）"
}

支持的操作符（op）：
- eq: 等于 (=)
- ne: 不等于 (!=)
- gt: 大于 (>)
- ge: 大于等于 (>=)
- lt: 小于 (<)
- le: 小于等于 (<=)
- contains: 包含
- startswith: 开头是
- endswith: 结尾是
- search: 全文模糊搜索

注意事项：
- fields 为空时返回该模型全部字段
- conditions 为空时返回该模型全部实例
- limit 默认 50，最大 200
- 支持同时查询多个模型（queries 数组可包含多个元素）
- analysis 字段用于请求你对查询结果进行统计分析，留空则不分析
- 只返回 JSON，不要添加其他文字
`)

	return sb.String(), nil
}

// BuildUserMessage 构建用户消息（带历史对话）
func (lp *llmPrompt) BuildUserMessage(message string, history []ChatMessage) []ChatMessage {
	messages := make([]ChatMessage, 0, len(history)+1)
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: message,
	})
	return messages
}
