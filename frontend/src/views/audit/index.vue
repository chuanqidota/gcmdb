<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <span>审计日志</span>
      </div>
    </template>

    <div class="search-bar">
      <el-select v-model="filter.resource_type" placeholder="资源类型" clearable style="width: 120px" size="small" @change="loadData">
        <el-option label="instance" value="instance" />
        <el-option label="model" value="model" />
        <el-option label="model_group" value="model_group" />
        <el-option label="model_relation" value="model_relation" />
        <el-option label="model_field" value="model_field" />
      </el-select>
      <el-select v-model="filter.action" placeholder="操作类型" clearable style="width: 100px" size="small" @change="loadData">
        <el-option label="create" value="create" />
        <el-option label="update" value="update" />
        <el-option label="delete" value="delete" />
      </el-select>
      <el-button size="small" @click="resetFilter">重置</el-button>
    </div>

    <el-table :data="list" stripe v-loading="loading" style="margin-top: 12px">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="resource_type" label="资源类型" width="120">
        <template #default="{ row }">
          <el-tag size="small" effect="plain">{{ row.resource_type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="resource_id" label="资源ID" width="80" />
      <el-table-column prop="action" label="操作" width="80">
        <template #default="{ row }">
          <el-tag :type="actionTagType(row.action)" size="small">{{ row.action }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="method" label="方法" width="70" />
      <el-table-column prop="path" label="路径" show-overflow-tooltip />
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
import { ref, onMounted } from 'vue'
import { listAuditLog } from '../../api/audit'

const list = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filter = ref({ resource_type: '', action: '' })

const actionTagType = (action) => {
  if (action === 'create') return 'success'
  if (action === 'delete') return 'danger'
  if (action === 'update') return 'warning'
  return 'info'
}

const loadData = async () => {
  loading.value = true
  try {
    const params = { limit: pageSize.value, offset: (page.value - 1) * pageSize.value }
    if (filter.value.resource_type) params.resource_type = filter.value.resource_type
    if (filter.value.action) params.action = filter.value.action
    const res = await listAuditLog(params)
    list.value = res.data.results || []
    total.value = res.data.count || 0
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filter.value = { resource_type: '', action: '' }
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
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
