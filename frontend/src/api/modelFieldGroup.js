import request from './request'

export const createModelFieldGroup = (data) => request.post('/models-field-group', data)
export const updateModelFieldGroup = (id, data) => request.put(`/models-field-group/${id}`, data)
export const deleteModelFieldGroup = (id) => request.delete(`/models-field-group/${id}`)
