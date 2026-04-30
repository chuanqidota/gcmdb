import request from './request'

export const listModelFieldUnique = (modelId) => request.get(`/models-field-unique/${modelId}`)
export const createModelFieldUnique = (data) => request.post('/models-field-unique', data)
export const deleteModelFieldUnique = (id) => request.delete(`/models-field-unique/${id}`)
