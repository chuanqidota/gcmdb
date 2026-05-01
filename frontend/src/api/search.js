import openapiRequest from './openapiRequest'

// 获取所有模型
export const getAllModels = () => openapiRequest.get('/model/all')

// 获取单个模型详情（含字段、分组、关联等）
export const getModelDetail = (id) => openapiRequest.get('/model/single', { params: { id } })

// 全文检索
export const fulltextSearch = (data) => openapiRequest.post('/instance/fulltext', data)

// 结构化实例搜索
export const searchInstance = (data) => openapiRequest.post('/instance/search', data)
