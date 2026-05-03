import request from './request'

export const listAuditLog = (params) => request.get('/audit-log', { params })
