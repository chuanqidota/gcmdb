import request from './request'

export const listModel = (params) => request.get('/models', { params })
export const createModel = (data) => request.post('/models', data)
export const retrieveModel = (id) => request.get(`/models/${id}`)
export const updateModel = (id, data) => request.put(`/models/${id}`, data)
export const deleteModel = (id) => request.delete(`/models/${id}`)
