import request from './request'

export const listModelRelation = (params) => request.get('/models-relation', { params })
export const createModelRelation = (data) => request.post('/models-relation', data)
export const deleteModelRelation = (id) => request.delete(`/models-relation/${id}`)
