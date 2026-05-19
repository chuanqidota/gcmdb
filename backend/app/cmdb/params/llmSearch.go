package params

// LLMSearchRequest LLM 搜索请求
type LLMSearchRequest struct {
	Message     string   `json:"message" binding:"required" label:"用户消息"`
	SessionId   string   `json:"session_id" label:"会话ID"`
	ModelFilter []string `json:"model_filter" label:"限定模型别名范围"`
}

// LLMClearSessionRequest 清空会话请求
type LLMClearSessionRequest struct {
	SessionId string `json:"session_id" binding:"required" label:"会话ID"`
}
