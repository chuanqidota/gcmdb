import request from './request'

export const listModelGroup = (params) => request.get('/models-group', { params })
export const createModelGroup = (data) => request.post('/models-group', data)
export const patchModelGroup = (id, data) => request.patch(`/models-group/${id}`, data)
export const deleteModelGroup = (id) => request.delete(`/models-group/${id}`)
