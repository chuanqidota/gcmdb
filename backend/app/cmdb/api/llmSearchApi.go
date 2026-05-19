package api

import (
	"fmt"
	"gcmdb/app/cmdb/params"
	"gcmdb/app/cmdb/resp"
	"gcmdb/app/cmdb/utils"
	"gcmdb/pkg/logger"
	"gcmdb/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type llmSearch struct{}

var LLMSearch = new(llmSearch)

// Search 大模型搜索
func (ls *llmSearch) Search(c *gin.Context) {
	var body params.LLMSearchRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}

	startTime := time.Now()

	// 生成或使用已有的 session ID
	sessionId := body.SessionId
	if sessionId == "" {
		sessionId = uuid.New().String()
	}

	// 构建 system prompt
	systemPrompt, err := utils.LLMPrompt.BuildSystemPrompt(body.ModelFilter)
	if err != nil {
		response.Fail(c, fmt.Sprintf("构建提示词失败-%s", err.Error()))
		return
	}

	// 获取或创建会话
	session := utils.LLMSession.GetOrCreate(sessionId, systemPrompt)

	// 添加用户消息
	utils.LLMSession.AddMessage(sessionId, "user", body.Message)

	// 调用大模型
	content, tokenUsed, err := utils.LLMClient.Chat(session.Messages)
	if err != nil {
		logger.Error(fmt.Sprintf("LLM 调用失败: %s", err.Error()))
		response.Fail(c, err.Error())
		return
	}

	// 解析意图
	intent, err := utils.LLMIntent.ParseIntent(content)
	if err != nil {
		logger.Error(fmt.Sprintf("意图解析失败: %s, 原始内容: %s", err.Error(), content))
		// 返回原始回复供用户查看
		response.Success(c, "执行成功", resp.LLMSearchResponse{
			AIAnalysis: content,
			Results:    []resp.LLMQueryResult{},
			SessionId:  sessionId,
			TokenUsed:  tokenUsed,
		})
		return
	}

	// 执行查询
	results, err := utils.LLMIntent.ExecuteQueries(intent.Queries)
	if err != nil {
		logger.Error(fmt.Sprintf("查询执行失败: %s", err.Error()))
		response.Fail(c, fmt.Sprintf("查询执行失败-%s", err.Error()))
		return
	}

	// 添加助手回复到会话历史
	utils.LLMSession.AddMessage(sessionId, "assistant", content)

	// 如果有分析请求，再次调用大模型
	aiAnalysis := ""
	if intent.Analysis != "" && len(results) > 0 {
		analysisPrompt := fmt.Sprintf("用户请求: %s\n\n查询结果:\n%s\n\n请根据以上数据进行分析：%s",
			body.Message,
			utils.LLMIntent.FormatResultsForAnalysis(results),
			intent.Analysis)

		analysisMessages := []utils.ChatMessage{
			{Role: "system", Content: "你是数据分析助手，请根据提供的数据进行简洁的分析总结。"},
			{Role: "user", Content: analysisPrompt},
		}

		analysisResult, analysisTokens, err := utils.LLMClient.Chat(analysisMessages)
		if err == nil {
			aiAnalysis = analysisResult
			tokenUsed += analysisTokens
		}
	}

	elapsed := time.Since(startTime)
	logger.Info(fmt.Sprintf("LLM 搜索完成: session=%s, 耗时=%v, tokens=%d", sessionId, elapsed, tokenUsed))

	response.Success(c, "执行成功", resp.LLMSearchResponse{
		AIAnalysis: aiAnalysis,
		Results:    results,
		SessionId:  sessionId,
		TokenUsed:  tokenUsed,
		Thinking:   intent.Thinking,
	})
}

// ClearSession 清空会话
func (ls *llmSearch) ClearSession(c *gin.Context) {
	var body params.LLMClearSessionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}

	utils.LLMSession.ClearSession(body.SessionId)
	response.Success(c, "会话已清空", nil)
}
