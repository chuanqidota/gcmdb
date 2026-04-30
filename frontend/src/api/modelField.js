import request from './request'

export const listModelField = (modelId) => request.get(`/models-field/${modelId}`)
export const createModelField = (data) => request.post('/models-field', data)
export const updateModelField = (id, data) => request.put(`/models-field/${id}`, data)
export const deleteModelField = (id) => request.delete(`/models-field/${id}`)
