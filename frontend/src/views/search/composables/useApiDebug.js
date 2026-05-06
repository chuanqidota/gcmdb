import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getModelDetail } from '../../../api/search'
import { listSearchDirectSql } from '../../../api/searchDirectSql'
import openapiRequest from '../../../api/openapiRequest'
import { formatJson } from '../../../utils/format'

export function useApiDebug(allModels) {
  const apiSavedSqls = ref([])
  const apiModelFields = ref([])
  const apiEndpointGroups = reactive([
    {
      name: '模型查询',
      collapsed: false,
      endpoints: [
        { key: 'model-all', name: '所有模型', method: 'GET', path: '/model/all', params: [] },
        { key: 'model-group', name: '模型分组', method: 'GET', path: '/model/group', params: [] },
        { key: 'model-single', name: '模型详情', method: 'GET', path: '/model/single', params: [
          { key: 'alias', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' }
        ]},
        { key: 'model-relation-type', name: '关系类型', method: 'GET', path: '/model/relation-type', params: [] },
        { key: 'model-relation', name: '模型关系', method: 'GET', path: '/model/relation', params: [] },
      ]
    },
    {
      name: '实例查询',
      collapsed: false,
      endpoints: [
        { key: 'instance-fulltext', name: '全文搜索', method: 'POST', path: '/instance/fulltext', params: [
          { key: 'search', label: '搜索词', type: 'input', required: true },
          { key: 'model_alias', label: '模型过滤', type: 'select', source: 'models', placeholder: '可选' },
          { key: 'limit', label: '数量', type: 'number' },
          { key: 'offset', label: '偏移', type: 'number' },
        ]},
        { key: 'instance-search', name: '结构化搜索', method: 'POST', path: '/instance/search', params: [
          { key: 'model', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' },
          { key: 'fields', label: '返回字段', type: 'multi-select', source: 'fields', placeholder: '不选则返回全部' },
          { key: 'where', label: '查询条件', type: 'condition-builder' },
          { key: 'order', label: '排序', type: 'select', source: 'orderFields', placeholder: '可选' },
          { key: 'limit', label: '数量', type: 'number' },
          { key: 'offset', label: '偏移', type: 'number' },
        ]},
        { key: 'instance-detail', name: '实例详情', method: 'GET', path: '/instance/detail', params: [
          { key: 'id', label: '实例ID', type: 'input', required: true, placeholder: '输入实例ID' },
        ]},
        { key: 'instance-topology', name: '拓扑关系', method: 'GET', path: '/instance/topology', params: [
          { key: 'id', label: '实例ID', type: 'input', required: true, placeholder: '输入实例ID' },
          { key: 'model', label: '模型过滤', type: 'select', source: 'models', placeholder: '可选' },
        ]},
        { key: 'instance-source', name: '查询上游实例', method: 'GET', path: '/instance/source', params: [
          { key: 'id', label: '目标实例ID', type: 'input', required: true },
          { key: 'model', label: '源模型', type: 'select', required: true, source: 'models', placeholder: '选择源模型' },
        ]},
        { key: 'instance-target', name: '查询下游实例', method: 'GET', path: '/instance/target', params: [
          { key: 'id', label: '源实例ID', type: 'input', required: true },
          { key: 'model', label: '目标模型', type: 'select', required: true, source: 'models', placeholder: '选择目标模型' },
        ]},
        { key: 'instance-direct', name: '执行SQL查询', method: 'GET', path: '/instance/direct', params: [
          { key: 'uuid', label: '查询', type: 'select', source: 'sqls', required: true, placeholder: '选择已保存的查询' },
        ]},
      ]
    },
    {
      name: '实例写入',
      collapsed: false,
      endpoints: [
        { key: 'instance-create', name: '创建实例', method: 'POST', path: '/instance/create', params: [
          { key: 'model_alias', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' },
          { key: 'data', label: '数据 (JSON)', type: 'textarea', required: true, placeholder: '{"field1": "value1", ...}' },
        ]},
        { key: 'instance-update', name: '更新实例', method: 'POST', path: '/instance/update', params: [
          { key: 'id', label: '实例ID', type: 'input', required: true },
          { key: 'data', label: '数据 (JSON)', type: 'textarea', required: true, placeholder: '{"field1": "new_value", ...}' },
        ]},
        { key: 'instance-delete', name: '删除实例', method: 'POST', path: '/instance/delete', params: [
          { key: 'model_alias', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' },
          { key: 'id', label: '实例ID', type: 'input', required: true },
        ]},
      ]
    },
    {
      name: '实例关系',
      collapsed: false,
      endpoints: [
        { key: 'relation-create', name: '创建关系', method: 'POST', path: '/instance-relation/create', params: [
          { key: 'source_model', label: '源模型', type: 'select', required: true, source: 'models', placeholder: '选择源模型' },
          { key: 'target_model', label: '目标模型', type: 'select', required: true, source: 'models', placeholder: '选择目标模型' },
          { key: 'source_id', label: '源实例ID', type: 'input', required: true },
          { key: 'target_ids', label: '目标实例ID', type: 'input', required: true, placeholder: '多个用逗号分隔' },
        ]},
        { key: 'relation-delete', name: '删除关系', method: 'POST', path: '/instance-relation/delete', params: [
          { key: 'source_model', label: '源模型', type: 'select', required: true, source: 'models', placeholder: '选择源模型' },
          { key: 'target_model', label: '目标模型', type: 'select', required: true, source: 'models', placeholder: '选择目标模型' },
          { key: 'source_id', label: '源实例ID', type: 'input', required: true },
          { key: 'target_id', label: '目标实例ID', type: 'input', required: true },
        ]},
      ]
    },
  ])

  const apiSelectedKey = ref('')
  const apiSelected = ref(null)
  const apiParamValues = reactive({})
  const apiConditions = ref([])
  const apiSending = ref(false)
  const apiRequestUrl = ref('')
  const apiRequestBody = ref('')
  const apiResponseData = ref(null)
  const apiResponseStatus = ref(null)
  const apiResponseTime = ref(0)

  const apiOperators = [
    { label: '等于 (=)', value: 'eq' },
    { label: '不等于 (!=)', value: 'ne' },
    { label: '大于 (>)', value: 'gt' },
    { label: '大于等于 (>=)', value: 'ge' },
    { label: '小于 (<)', value: 'lt' },
    { label: '小于等于 (<=)', value: 'le' },
    { label: '包含 (in)', value: 'in' },
    { label: '模糊匹配', value: 'contains' },
    { label: '前缀匹配', value: 'startswith' },
    { label: '后缀匹配', value: 'endswith' },
  ]

  const apiOrderFields = [
    { label: 'ID 升序', value: 'id' },
    { label: 'ID 降序', value: '-id' },
    { label: '创建时间 升序', value: 'created_at' },
    { label: '创建时间 降序', value: '-created_at' },
    { label: '更新时间 升序', value: 'updated_at' },
    { label: '更新时间 降序', value: '-updated_at' },
  ]

  const getApiOptions = (p) => {
    switch (p.source) {
      case 'models':
        return allModels.value.map(m => ({ label: `${m.name} (${m.alias})`, value: m.alias }))
      case 'fields':
        return apiModelFields.value.map(f => ({ label: `${f.name} (${f.alias})`, value: f.alias }))
      case 'sqls':
        return apiSavedSqls.value.map(s => ({ label: `${s.name} (${s.uuid?.slice(0, 8)}...)`, value: s.uuid }))
      case 'orderFields':
        return apiOrderFields
      default:
        return []
    }
  }

  const selectApiEndpoint = (ep) => {
    apiSelectedKey.value = ep.key
    apiSelected.value = ep
    for (const key of Object.keys(apiParamValues)) {
      apiParamValues[key] = undefined
    }
    ep.params.forEach(p => {
      if (p.type === 'number') apiParamValues[p.key] = undefined
      else if (p.type === 'multi-select') apiParamValues[p.key] = []
      else if (p.type === 'textarea') apiParamValues[p.key] = ''
      else apiParamValues[p.key] = ''
    })
    apiConditions.value = []
    apiModelFields.value = []
    updateApiPreview()
  }

  watch(() => apiParamValues.model, async (alias) => {
    if (!alias) { apiModelFields.value = []; return }
    try {
      const res = await getModelDetail(alias)
      apiModelFields.value = (res.data?.model_fields || [])
    } catch {
      apiModelFields.value = []
    }
  })

  const buildGetParams = (params, paramValues, conditions) => {
    const qs = new URLSearchParams()
    for (const p of params) {
      const val = paramValues[p.key]
      if (p.type === 'condition-builder') {
        const where = conditions.filter(c => c.field && c.op).map(c => {
          let v = c.value
          if (c.op === 'in') v = v.split(',').map(s => s.trim())
          return { [c.op]: { [c.field]: v } }
        })
        if (where.length) qs.set('where', JSON.stringify(where))
      } else if (p.type === 'multi-select') {
        if (val?.length) qs.set(p.key, val.join(','))
      } else if (p.type === 'number') {
        if (val !== undefined && val !== null && val !== '') qs.set(p.key, String(val))
      } else {
        if (val) qs.set(p.key, val)
      }
    }
    const str = qs.toString()
    return str ? `?${str}` : ''
  }

  const buildPostBody = (ep, paramValues, conditions) => {
    const body = {}
    const hasConditionBuilder = ep.params.some(p => p.type === 'condition-builder')
    for (const p of ep.params) {
      const val = paramValues[p.key]
      if (p.type === 'condition-builder') {
        const where = conditions.filter(c => c.field && c.op).map(c => {
          let v = c.value
          if (c.op === 'in') v = v.split(',').map(s => s.trim())
          return { [c.op]: { [c.field]: v } }
        })
        if (where.length) {
          if (!body.__condition) body.__condition = {}
          body.__condition.where = where
        }
      } else if (p.type === 'textarea') {
        if (val) {
          try { body[p.key] = JSON.parse(val) } catch { body[p.key] = val }
        }
      } else if (p.type === 'number') {
        if (val !== undefined && val !== null && val !== '') {
          if (hasConditionBuilder && (p.key === 'limit' || p.key === 'offset')) {
            if (!body.__condition) body.__condition = {}
            body.__condition[p.key] = val
          } else {
            body[p.key] = val
          }
        }
      } else if (p.type === 'multi-select') {
        if (val?.length) body[p.key] = val
      } else {
        if (val) {
          if (hasConditionBuilder && p.key === 'order') {
            if (!body.__condition) body.__condition = {}
            body.__condition.order = [val]
          } else {
            body[p.key] = val
          }
        }
      }
    }
    if (body.target_ids && typeof body.target_ids === 'string') {
      body.target_ids = body.target_ids.split(',').map(s => parseInt(s.trim(), 10)).filter(n => !isNaN(n))
    }
    if (body.source_id) body.source_id = parseInt(body.source_id, 10) || body.source_id
    if (body.target_id) body.target_id = parseInt(body.target_id, 10) || body.target_id
    if (body.id && ep.key !== 'instance-update') body.id = parseInt(body.id, 10) || body.id
    return body
  }

  const buildApiRequest = () => {
    if (!apiSelected.value) return { url: '', body: null, method: '' }
    const ep = apiSelected.value
    const method = ep.method
    let url = ep.path
    let body = null
    if (method === 'GET') {
      url = ep.path + buildGetParams(ep.params, apiParamValues, apiConditions.value)
    } else {
      body = buildPostBody(ep, apiParamValues, apiConditions.value)
    }
    return { url, body, method }
  }

  const updateApiPreview = () => {
    const { url, body, method } = buildApiRequest()
    apiRequestUrl.value = method ? `${method} ${url}` : ''
    apiRequestBody.value = body ? JSON.stringify(body, null, 2) : ''
  }

  watch(apiParamValues, updateApiPreview, { deep: true })
  watch(apiConditions, updateApiPreview, { deep: true })

  const isApiReadOnly = computed(() => {
    const key = apiSelectedKey.value
    return key === 'instance-create' || key === 'instance-update' || key === 'instance-delete' ||
           key === 'relation-create' || key === 'relation-delete'
  })

  const sendApiRequest = async () => {
    if (!apiSelected.value || apiSending.value || isApiReadOnly.value) return
    const req = buildApiRequest()
    if (!req.url) return
    const { url, body, method } = req

    apiSending.value = true
    apiResponseData.value = null
    apiResponseStatus.value = null
    apiResponseTime.value = 0

    const startTime = Date.now()
    try {
      let res
      if (method === 'GET') {
        res = await openapiRequest.get(url)
      } else {
        res = await openapiRequest.post(url, body)
      }
      apiResponseTime.value = Date.now() - startTime
      apiResponseStatus.value = 200
      apiResponseData.value = res
    } catch (err) {
      apiResponseTime.value = Date.now() - startTime
      apiResponseStatus.value = err.response?.status || err.status || 0
      apiResponseData.value = err.response?.data || { error: err.message || '请求失败' }
    } finally {
      apiSending.value = false
    }
  }

  const copyApiRequest = () => {
    const { url, body, method } = buildApiRequest()
    let text = `${method} ${url}`
    if (body) text += `\n\n${JSON.stringify(body, null, 2)}`
    navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  }

  const copyApiResponse = () => {
    if (apiResponseData.value === null) return
    navigator.clipboard.writeText(formatJson(apiResponseData.value))
    ElMessage.success('已复制到剪贴板')
  }

  async function loadApiSavedSqls() {
    try {
      const sqlRes = await listSearchDirectSql({ limit: 1000 })
      apiSavedSqls.value = sqlRes?.data?.results || []
    } catch {
      apiSavedSqls.value = []
    }
  }

  return {
    apiEndpointGroups, apiSelectedKey, apiSelected, apiParamValues, apiConditions,
    apiSending, apiRequestUrl, apiRequestBody, apiResponseData, apiResponseStatus, apiResponseTime,
    apiOperators, apiModelFields, apiSavedSqls,
    getApiOptions, selectApiEndpoint, sendApiRequest, copyApiRequest, copyApiResponse,
    isApiReadOnly, loadApiSavedSqls, formatJson,
  }
}
