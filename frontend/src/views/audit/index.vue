<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>审计日志</span>
      </div>
    </template>

    <div class="search-bar">
      <el-select v-model="filter.category" placeholder="分类" clearable style="width: 120px" @change="onCategoryChange">
        <el-option v-for="(_, k) in categoryMap" :key="k" :label="k" :value="k" />
      </el-select>
      <el-select v-model="filter.resource_type" placeholder="资源类型" clearable style="width: 160px" @change="onPageChange">
        <el-option v-for="t in currentResourceTypes" :key="t" :label="resourceLabels[t] || t" :value="t" />
      </el-select>
      <el-select v-model="filter.action" placeholder="操作类型" clearable style="width: 100px" @change="onPageChange">
        <el-option label="创建" value="create" />
        <el-option label="更新" value="update" />
        <el-option label="删除" value="delete" />
        <el-option label="执行" value="execute" />
      </el-select>
      <el-input v-model="filter.username" placeholder="操作人" clearable style="width: 120px" @clear="onPageChange" @keyup.enter="onPageChange" />
      <el-input v-model="filter.resource_id" placeholder="资源ID" clearable style="width: 100px" @clear="onPageChange" @keyup.enter="onPageChange" />
      <el-input v-model="filter.keyword" placeholder="关键词搜索" clearable style="width: 150px" @clear="onPageChange" @keyup.enter="onPageChange" />
      <el-date-picker v-model="filter.timeRange" type="daterange" range-separator="-" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" style="width: 260px" @change="onPageChange" />
      <el-button @click="resetFilter">重置</el-button>
    </div>

    <el-table :data="list" stripe v-loading="loading" style="margin-top: 12px" row-key="id">
      <el-table-column type="expand">
        <template #default="{ row }">
          <div class="expand-content">
            <div v-if="!row.before_data && !row.after_data" class="expand-empty">无变更记录</div>
            <template v-else>
              <div class="diff-block" v-if="row.before_data">
                <div class="diff-label diff-label--before">变更前</div>
                <pre class="diff-json">{{ formatJson(row.before_data) }}</pre>
              </div>
              <div class="diff-block" v-if="row.after_data">
                <div class="diff-label diff-label--after">变更后</div>
                <pre class="diff-json">{{ formatJson(row.after_data) }}</pre>
              </div>
            </template>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="resource_type" label="资源类型" width="140">
        <template #default="{ row }">
          <el-tag size="small" effect="plain">{{ resourceLabels[row.resource_type] || row.resource_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="resource_id" label="资源ID" width="80" />
      <el-table-column prop="action" label="操作" width="80">
        <template #default="{ row }">
          <el-tag :type="actionTagType(row.action)" size="small">{{ actionLabels[row.action] || row.action }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="username" label="操作人" width="100" />
      <el-table-column prop="source_ip" label="来源 IP" width="130" />
      <el-table-column prop="status_code" label="状态码" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status_code < 400 ? 'success' : 'danger'" size="small">{{ row.status_code }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="170" />
    </el-table>

    <div class="pagination-wrap">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @change="loadData"
      />
    </div>
  </el-card>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { listAuditLog } from '../../api/audit'

const categoryMap = {
  '模型配置': ['model_group', 'model', 'model_field_group', 'model_field', 'model_field_unique', 'model_field_relation', 'model_relation_type', 'model_relation'],
  '实例数据': ['instance', 'instance_relation'],
  '其他': ['search_direct_sql', 'task'],
}

const resourceLabels = {
  model_group: '模型分组',
  model: '模型',
  model_field_group: '字段分组',
  model_field: '字段',
  model_field_unique: '唯一字段',
  model_field_relation: '字段关系',
  model_relation_type: '关系类型',
  model_relation: '模型关系',
  instance: '实例',
  instance_relation: '实例关系',
  search_direct_sql: 'SQL查询',
  task: '任务',
}

const actionLabels = { create: '创建', update: '更新', delete: '删除', execute: '执行' }

const list = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filter = ref({ category: '', resource_type: '', action: '', username: '', resource_id: '', keyword: '', timeRange: null })

const currentResourceTypes = computed(() => {
  if (filter.value.category) return categoryMap[filter.value.category] || []
  return Object.keys(resourceLabels)
})

const actionTagType = (action) => {
  if (action === 'create') return 'success'
  if (action === 'delete') return 'danger'
  if (action === 'update') return 'warning'
  return 'info'
}

const formatJson = (str) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

const onCategoryChange = () => {
  filter.value.resource_type = ''
  onPageChange()
}

const onPageChange = () => {
  page.value = 1
  loadData()
}

const loadData = async () => {
  loading.value = true
  try {
    const params = { limit: pageSize.value, offset: (page.value - 1) * pageSize.value }
    if (filter.value.resource_type) params.resource_type = filter.value.resource_type
    if (filter.value.action) params.action = filter.value.action
    if (filter.value.username) params.username = filter.value.username
    if (filter.value.resource_id) params.resource_id = filter.value.resource_id
    if (filter.value.keyword) params.keyword = filter.value.keyword
    if (filter.value.timeRange?.[0]) params.start_time = filter.value.timeRange[0]
    if (filter.value.timeRange?.[1]) params.end_time = filter.value.timeRange[1] + ' 23:59:59'
    const res = await listAuditLog(params)
    list.value = res.data.results || []
    total.value = res.data.count || 0
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filter.value = { category: '', resource_type: '', action: '', username: '', resource_id: '', keyword: '', timeRange: null }
  page.value = 1
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.search-bar {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.expand-content {
  padding: 12px 24px;
  display: flex;
  gap: 24px;
}

.expand-empty {
  color: var(--color-text-muted);
  font-size: 13px;
}

.diff-block {
  flex: 1;
  min-width: 0;
}

.diff-label {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 6px;
  padding: 2px 8px;
  border-radius: 4px;
  display: inline-block;
}

.diff-label--before {
  background: #FEE2E2;
  color: #DC2626;
}

.diff-label--after {
  background: #D1FAE5;
  color: #059669;
}

.diff-json {
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 12px;
  font-size: 12px;
  line-height: 1.6;
  overflow-x: auto;
  max-height: 400px;
  margin: 0;
  font-family: 'Menlo', 'Monaco', 'Consolas', monospace;
}
</style>
