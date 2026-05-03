import openapiRequest from './openapiRequest'

// 获取所有模型
export const getAllModels = () => openapiRequest.get('/model/all')

// 获取所有模型分组
export const getModelGroups = () => openapiRequest.get('/model/group')

// 获取单个模型详情（含字段、分组、关联等）
export const getModelDetail = (alias) => openapiRequest.get('/model/single', { params: { alias } })

// 全文检索
export const fulltextSearch = (data) => openapiRequest.post('/instance/fulltext', data)

// 结构化实例搜索
export const searchInstance = (data) => openapiRequest.post('/instance/search', data)

// 获取所有模型关系类型
export const getModelRelationTypes = () => openapiRequest.get('/model/relation-type')

// 获取所有模型关系
export const getAllModelRelations = () => openapiRequest.get('/model/relation')

// 获取实例拓扑（上下游关系）
export const getInstanceTopology = (id, model) => openapiRequest.get('/instance/topology', { params: { id, model } })
