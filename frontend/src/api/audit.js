import request from './request'

export const listAuditLog = (params) => request.get('/audit-log', { params })
export const retrieveAuditLog = (resourceType, resourceId) => request.get(`/audit-log/${resourceType}/${resourceId}`)
