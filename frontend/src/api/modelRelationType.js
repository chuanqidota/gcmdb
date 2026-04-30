import request from './request'

export const listModelRelationType = (params) => request.get('/models-relation-type', { params })
export const createModelRelationType = (data) => request.post('/models-relation-type', data)
export const updateModelRelationType = (id, data) => request.put(`/models-relation-type/${id}`, data)
export const deleteModelRelationType = (id) => request.delete(`/models-relation-type/${id}`)
