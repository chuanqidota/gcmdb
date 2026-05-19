package resp

// LLMQueryResult 单个模型的查询结果
type LLMQueryResult struct {
	Model      string                   `json:"model"`
	ModelName  string                   `json:"model_name"`
	ModelAlias string                   `json:"model_alias"`
	Count      int64                    `json:"count"`
	Data       []map[string]interface{} `json:"data"`
}

// LLMSearchResponse LLM 搜索响应
type LLMSearchResponse struct {
	AIAnalysis string           `json:"ai_analysis"`
	Results    []LLMQueryResult `json:"results"`
	SessionId  string           `json:"session_id"`
	TokenUsed  int              `json:"token_used"`
	Thinking   string           `json:"thinking"`
}
