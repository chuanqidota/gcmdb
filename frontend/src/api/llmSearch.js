import request from './request'

// 大模型搜索
export const llmSearch = (data) => request.post('/llm-search', data, { timeout: 60000 })

// 清空会话
export const llmClearSession = (data) => request.post('/llm-search/clear', data)
