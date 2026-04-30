import request from './request'

export const listSearchDirectSql = (params) => request.get('/search-direct-sql', { params })
export const createSearchDirectSql = (data) => request.post('/search-direct-sql', data)
export const executeSearchDirectSql = (id) => request.get(`/search-direct-sql/${id}`)
export const updateSearchDirectSql = (id, data) => request.put(`/search-direct-sql/${id}`, data)
export const deleteSearchDirectSql = (id) => request.delete(`/search-direct-sql/${id}`)
