<template>
  <div class="api-test-page">
    <!-- 左侧：接口选择 + 参数面板 -->
    <div class="left-panel">
      <div class="panel-title">OpenAPI 接口</div>
      <div class="endpoint-groups">
        <div v-for="group in endpointGroups" :key="group.name" class="endpoint-group">
          <div class="group-title" @click="group.collapsed = !group.collapsed">
            <el-icon><ArrowRight v-if="group.collapsed" /><ArrowDown v-else /></el-icon>
            {{ group.name }}
          </div>
          <div v-show="!group.collapsed" class="group-items">
            <div
              v-for="ep in group.endpoints"
              :key="ep.key"
              :class="['endpoint-item', { active: selectedKey === ep.key }]"
              @click="selectEndpoint(ep)"
            >
              <span class="method-tag" :class="ep.method.toLowerCase()">{{ ep.method }}</span>
              <span class="ep-name">{{ ep.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 参数表单 -->
      <div v-if="selected" class="param-section">
        <div class="panel-title">参数</div>
        <div class="param-form">
          <div v-for="p in selected.params" :key="p.key" class="param-row">
            <label class="param-label">
              {{ p.label }}
              <span v-if="p.required" class="required">*</span>
            </label>

            <!-- 普通输入 -->
            <el-input v-if="p.type === 'input'" v-model="paramValues[p.key]" :placeholder="p.placeholder || ''" size="small" />

            <!-- 数字输入 -->
            <el-input-number v-else-if="p.type === 'number'" v-model="paramValues[p.key]" :min="0" :placeholder="p.label" size="small" style="width: 100%" />

            <!-- 单选下拉 -->
            <el-select v-else-if="p.type === 'select'" v-model="paramValues[p.key]" :placeholder="p.placeholder || '请选择'" size="small" style="width: 100%" filterable>
              <el-option v-for="opt in getOptions(p)" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>

            <!-- 多选下拉 -->
            <el-select v-else-if="p.type === 'multi-select'" v-model="paramValues[p.key]" :placeholder="p.placeholder || '可多选'" size="small" style="width: 100%" multiple filterable collapse-tags collapse-tags-tooltip>
              <el-option v-for="opt in getOptions(p)" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>

            <!-- 条件构建器 -->
            <div v-else-if="p.type === 'condition-builder'" class="condition-builder">
              <div v-for="(cond, idx) in conditions" :key="idx" class="cond-row">
                <el-select v-model="cond.field" placeholder="字段" size="small" style="width: 100px" filterable>
                  <el-option v-for="f in modelFields" :key="f.alias" :label="f.name" :value="f.alias" />
                </el-select>
                <el-select v-model="cond.op" placeholder="操作符" size="small" style="width: 110px">
                  <el-option v-for="op in operators" :key="op.value" :label="op.label" :value="op.value" />
                </el-select>
                <el-input v-model="cond.value" placeholder="值" size="small" style="flex: 1" />
                <el-button link type="danger" size="small" @click="conditions.splice(idx, 1)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
              <el-button size="small" @click="conditions.push({ field: '', op: 'eq', value: '' })">
                <el-icon><Plus /></el-icon> 添加条件
              </el-button>
            </div>
          </div>

          <el-button type="primary" style="width: 100%; margin-top: 12px" @click="sendRequest" :loading="sending">
            发送请求
          </el-button>
        </div>
      </div>
    </div>

    <!-- 右侧：请求预览 + 响应 -->
    <div class="right-panel">
      <!-- 请求预览 -->
      <div class="request-preview">
        <div class="preview-header">
          <span class="panel-title">请求</span>
          <el-button size="small" link @click="copyRequest">
            <el-icon><CopyDocument /></el-icon> 复制
          </el-button>
        </div>
        <div class="preview-url">{{ requestUrl }}</div>
        <pre v-if="requestBody" class="preview-body">{{ requestBody }}</pre>
      </div>

      <!-- 响应结果 -->
      <div class="response-section">
        <div class="response-header">
          <span class="panel-title">响应</span>
          <span v-if="responseStatus !== null" class="response-meta">
            <el-tag :type="responseStatus === 200 ? 'success' : 'danger'" size="small">{{ responseStatus }}</el-tag>
            <span class="response-time">{{ responseTime }}ms</span>
          </span>
        </div>
        <div v-if="responseData === null" class="response-empty">
          点击「发送请求」查看响应结果
        </div>
        <pre v-else class="response-body"><code>{{ formatJson(responseData) }}</code></pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowRight, ArrowDown, Delete, Plus, CopyDocument } from '@element-plus/icons-vue'
import { getAllModels, getModelGroups, getModelRelationTypes, getAllModelRelations, getModelDetail } from '../../api/search'
import { listSearchDirectSql } from '../../api/searchDirectSql'
import openapiRequest from '../../api/openapiRequest'

// ===== 数据源 =====
const allModels = ref([])
const allGroups = ref([])
const allRelationTypes = ref([])
const allRelations = ref([])
const savedSqls = ref([])
const modelFields = ref([])

// ===== 接口定义 =====
const endpointGroups = reactive([
  {
    name: '模型查询',
    collapsed: false,
    endpoints: [
      { key: 'model-all', name: '获取所有模型', method: 'GET', path: '/openapi/model/all', params: [] },
      { key: 'model-group', name: '获取所有分组', method: 'GET', path: '/openapi/model/group', params: [] },
      { key: 'model-relation-type', name: '获取关系类型', method: 'GET', path: '/openapi/model/relation-type', params: [] },
      { key: 'model-relation', name: '获取所有关系', method: 'GET', path: '/openapi/model/relation', params: [] },
      {
        key: 'model-single', name: '获取模型详情', method: 'GET', path: '/openapi/model/single',
        params: [{ key: 'id', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' }]
      },
    ]
  },
  {
    name: '实例查询',
    collapsed: false,
    endpoints: [
      {
        key: 'instance-fulltext', name: '全文搜索', method: 'POST', path: '/openapi/instance/fulltext',
        params: [
          { key: 'search', label: '搜索词', type: 'input', required: true },
          { key: 'model_alias', label: '模型过滤', type: 'multi-select', source: 'models', placeholder: '可选，不选则搜索全部' },
          { key: 'limit', label: '数量', type: 'number' },
          { key: 'offset', label: '偏移', type: 'number' },
        ]
      },
      {
        key: 'instance-search', name: '结构化搜索', method: 'POST', path: '/openapi/instance/search',
        params: [
          { key: 'model', label: '模型', type: 'select', required: true, source: 'models', placeholder: '选择模型' },
          { key: 'fields', label: '返回字段', type: 'multi-select', source: 'fields', placeholder: '不选则返回全部' },
          { key: 'where', label: '查询条件', type: 'condition-builder' },
          { key: 'order', label: '排序字段', type: 'select', source: 'orderFields', placeholder: '可选' },
          { key: 'limit', label: '数量', type: 'number' },
          { key: 'offset', label: '偏移', type: 'number' },
        ]
      },
      {
        key: 'instance-source', name: '查询源实例', method: 'POST', path: '/openapi/instance/source',
        params: [
          { key: 'id', label: '目标实例ID', type: 'input', required: true },
          { key: 'model', label: '源模型', type: 'select', source: 'models', placeholder: '选择源模型' },
        ]
      },
      {
        key: 'instance-target', name: '查询目标实例', method: 'POST', path: '/openapi/instance/target',
        params: [
          { key: 'id', label: '源实例ID', type: 'input', required: true },
          { key: 'model', label: '目标模型', type: 'select', source: 'models', placeholder: '选择目标模型' },
        ]
      },
      {
        key: 'instance-direct', name: '执行保存的SQL', method: 'POST', path: '/openapi/instance/direct',
        params: [
          { key: 'uuid', label: '查询', type: 'select', source: 'sqls', required: true, placeholder: '选择已保存的查询' },
        ]
      },
    ]
  }
])

// ===== 状态 =====
const selectedKey = ref('')
const selected = ref(null)
const paramValues = reactive({})
const conditions = ref([])
const sending = ref(false)
const requestUrl = ref('')
const requestBody = ref('')
const responseData = ref(null)
const responseStatus = ref(null)
const responseTime = ref(0)

const operators = [
  { label: '等于 (=)', value: 'eq' },
  { label: '不等于 (!=)', value: 'ne' },
  { label: '大于 (>)', value: 'gt' },
  { label: '大于等于 (>=)', value: 'ge' },
  { label: '小于 (<)', value: 'lt' },
  { label: '小于等于 (<=)', value: 'le' },
  { label: '包含 (in)', value: 'in' },
  { label: '模糊匹配 (contains)', value: 'contains' },
  { label: '前缀匹配', value: 'startswith' },
  { label: '后缀匹配', value: 'endswith' },
]

const orderFields = [
  { label: 'ID 升序', value: 'id' },
  { label: 'ID 降序', value: '-id' },
  { label: '创建时间 升序', value: 'created_at' },
  { label: '创建时间 降序', value: '-created_at' },
  { label: '更新时间 升序', value: 'updated_at' },
  { label: '更新时间 降序', value: '-updated_at' },
]

// ===== 下拉选项 =====
const getOptions = (p) => {
  switch (p.source) {
    case 'models':
      return allModels.value.map(m => ({ label: `${m.name} (${m.alias})`, value: m.alias }))
    case 'fields':
      return modelFields.value.map(f => ({ label: `${f.name} (${f.alias})`, value: f.alias }))
    case 'sqls':
      return savedSqls.value.map(s => ({ label: `${s.name} (${s.uuid?.slice(0, 8)}...)`, value: s.uuid }))
    case 'orderFields':
      return orderFields
    default:
      return []
  }
}

// ===== 接口选择 =====
const selectEndpoint = (ep) => {
  selectedKey.value = ep.key
  selected.value = ep
  // 重置参数
  Object.keys(paramValues).forEach(k => delete paramValues[k])
  conditions.value = []
  modelFields.value = []
  // 设置默认值
  ep.params.forEach(p => {
    if (p.type === 'number') paramValues[p.key] = undefined
    else if (p.type === 'multi-select') paramValues[p.key] = []
    else paramValues[p.key] = ''
  })
  updatePreview()
}

// ===== 监听模型选择变化（加载字段列表） =====
watch(() => paramValues.model, async (alias) => {
  if (!alias) { modelFields.value = []; return }
  const model = allModels.value.find(m => m.alias === alias)
  if (model) {
    try {
      const res = await getModelDetail(model.id)
      modelFields.value = (res.data?.model_fields || [])
    } catch {
      modelFields.value = []
    }
  }
})

// ===== 构建请求 =====
const buildRequest = () => {
  if (!selected.value) return { url: '', body: null, method: '' }

  const ep = selected.value
  const method = ep.method
  let url = ep.path
  let body = null

  if (ep.key === 'model-single') {
    url = `${ep.path}?id=${paramValues.id || ''}`
  } else if (ep.key === 'instance-source' || ep.key === 'instance-target') {
    // source/target 用 query params
    url = `${ep.path}?id=${paramValues.id || ''}&model=${paramValues.model || ''}`
  } else if (ep.key === 'instance-fulltext') {
    body = { search: paramValues.search || '' }
    if (paramValues.model_alias?.length) body.model_alias = paramValues.model_alias.join(',')
    if (paramValues.limit) body.limit = paramValues.limit
    if (paramValues.offset) body.offset = paramValues.offset
  } else if (ep.key === 'instance-search') {
    body = { model: paramValues.model || '' }
    if (paramValues.fields?.length) body.fields = paramValues.fields
    const condition = {}
    if (conditions.value.length) {
      condition.where = conditions.value.filter(c => c.field && c.op).map(c => {
        let val = c.value
        if (c.op === 'in') {
          val = val.split(',').map(v => v.trim())
        }
        return { [c.op]: { [c.field]: val } }
      })
    }
    if (paramValues.order) condition.order = [paramValues.order]
    if (paramValues.limit) condition.limit = paramValues.limit
    if (paramValues.offset) condition.offset = paramValues.offset
    if (Object.keys(condition).length) body.__condition = condition
  } else if (ep.key === 'instance-direct') {
    body = { uuid: paramValues.uuid || '' }
  }

  return { url, body, method }
}

const updatePreview = () => {
  const { url, body, method } = buildRequest()
  requestUrl.value = method ? `${method} ${url}` : ''
  requestBody.value = body ? JSON.stringify(body, null, 2) : ''
}

// 监听参数变化更新预览
watch(paramValues, updatePreview, { deep: true })
watch(conditions, updatePreview, { deep: true })

// ===== 发送请求 =====
const sendRequest = async () => {
  if (!selected.value) return

  const { url, body, method } = buildRequest()
  if (!url) return

  sending.value = true
  responseData.value = null
  responseStatus.value = null
  responseTime.value = 0

  const startTime = Date.now()
  try {
    let res
    if (method === 'GET') {
      res = await openapiRequest.get(url)
    } else {
      res = await openapiRequest.post(url, body)
    }
    responseTime.value = Date.now() - startTime
    responseStatus.value = 200
    responseData.value = res
  } catch (err) {
    responseTime.value = Date.now() - startTime
    responseStatus.value = err.response?.status || 0
    responseData.value = err.response?.data || { error: err.message }
  } finally {
    sending.value = false
  }
}

// ===== 复制请求 =====
const copyRequest = () => {
  const { url, body, method } = buildRequest()
  let text = `${method} ${url}`
  if (body) text += `\n\n${JSON.stringify(body, null, 2)}`
  navigator.clipboard.writeText(text)
  ElMessage.success('已复制到剪贴板')
}

// ===== 格式化 JSON =====
const formatJson = (data) => {
  try {
    return JSON.stringify(data, null, 2)
  } catch {
    return String(data)
  }
}

// ===== 初始化 =====
onMounted(async () => {
  const [modelsRes, groupsRes, typesRes, relationsRes, sqlRes] = await Promise.allSettled([
    getAllModels(),
    getModelGroups(),
    getModelRelationTypes(),
    getAllModelRelations(),
    listSearchDirectSql({ limit: 1000 }),
  ])
  allModels.value = modelsRes.status === 'fulfilled' ? (modelsRes.value.data || []) : []
  allGroups.value = groupsRes.status === 'fulfilled' ? (groupsRes.value.data || []) : []
  allRelationTypes.value = typesRes.status === 'fulfilled' ? (typesRes.value.data || []) : []
  allRelations.value = relationsRes.status === 'fulfilled' ? (relationsRes.value.data || []) : []
  savedSqls.value = sqlRes.status === 'fulfilled' ? (sqlRes.value.data?.results || []) : []
})
</script>

<style scoped>
.api-test-page {
  display: flex;
  gap: 16px;
  height: var(--page-height, calc(100vh - 116px));
}

.left-panel {
  width: 320px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
  overflow-y: auto;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
}

.panel-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
  padding: 12px 16px 8px;
}

.endpoint-groups {
  padding: 0 8px;
}

.group-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  padding: 8px;
  cursor: pointer;
  border-radius: 6px;
}

.group-title:hover {
  background: var(--color-muted);
}

.group-items {
  padding: 0 0 4px 0;
}

.endpoint-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px 6px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: var(--color-text-secondary);
  transition: all 0.15s;
}

.endpoint-item:hover {
  background: var(--color-muted);
  color: var(--color-text-primary);
}

.endpoint-item.active {
  background: var(--color-primary);
  color: #fff;
}

.method-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 3px;
  flex-shrink: 0;
}

.method-tag.get {
  background: #DBEAFE;
  color: #1D4ED8;
}

.method-tag.post {
  background: #FEF3C7;
  color: #92400E;
}

.endpoint-item.active .method-tag {
  background: rgba(255,255,255,0.25);
  color: #fff;
}

.ep-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.param-section {
  border-top: 1px solid var(--color-border);
  margin-top: 4px;
}

.param-form {
  padding: 0 16px 16px;
}

.param-row {
  margin-bottom: 12px;
}

.param-label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}

.required {
  color: #EF4444;
}

.condition-builder {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.cond-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 右侧 */
.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}

.request-preview {
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
}

.preview-header .panel-title {
  padding: 12px 0 8px;
}

.preview-url {
  padding: 0 16px 8px;
  font-family: 'SF Mono', 'Menlo', 'Consolas', monospace;
  font-size: 13px;
  color: var(--color-primary);
  word-break: break-all;
}

.preview-body {
  margin: 0;
  padding: 12px 16px;
  background: #1E293B;
  color: #E2E8F0;
  font-family: 'SF Mono', 'Menlo', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.6;
  overflow-x: auto;
  max-height: 200px;
}

.response-section {
  flex: 1;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.response-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
}

.response-header .panel-title {
  padding: 12px 0 8px;
}

.response-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.response-time {
  font-size: 12px;
  color: var(--color-text-muted);
}

.response-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  font-size: 13px;
}

.response-body {
  flex: 1;
  margin: 0;
  padding: 12px 16px;
  background: #1E293B;
  color: #E2E8F0;
  font-family: 'SF Mono', 'Menlo', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.6;
  overflow: auto;
}
</style>
