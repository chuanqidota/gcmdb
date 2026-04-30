import request from './request'

export const listInstance = (modelId, params) => request.get(`/instances/${modelId}`, { params })
export const createInstance = (data) => request.post('/instance', data)
export const retrieveInstance = (id) => request.get(`/instance/${id}`)
export const updateInstance = (id, data) => request.put(`/instance/${id}`, data)
export const deleteInstance = (id) => request.delete(`/instance/${id}`)
