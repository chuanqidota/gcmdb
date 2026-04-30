import request from './request'

export const listModelFieldRelation = (sourceModelId) => request.get(`/models-field-relation/${sourceModelId}`)
export const createModelFieldRelation = (data) => request.post('/models-field-relation', data)
export const deleteModelFieldRelation = (id) => request.delete(`/models-field-relation/${id}`)
