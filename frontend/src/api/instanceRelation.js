import request from './request'

export const createInstanceRelation = (data) => request.post('/instance-relation', data)
export const sourceTargetRelation = (sourceId) => request.get(`/source-target-relation/${sourceId}`)
export const targetSourceRelation = (targetId) => request.get(`/target-source-relation/${targetId}`)
export const deleteInstanceRelationByKeys = (params) => request.delete('/instance-relation-by-keys', { params })
