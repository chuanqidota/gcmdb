<template>
  <div class="sql-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>SQL 查询</span>
          <el-button type="primary" @click="showDialog()">
            <el-icon><Plus /></el-icon>新增查询
          </el-button>
        </div>
      </template>
      <el-skeleton :loading="loading" animated :count="5">
        <template #template>
          <div style="padding: 12px 0"><el-skeleton-item v-for="i in 5" :key="i" variant="text" style="height: 36px; margin-bottom: 6px" /></div>
        </template>
        <template #default>
          <el-table :data="list" stripe highlight-current-row>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="name" label="名称" width="160" />
            <el-table-column prop="sql" label="SQL 语句" show-overflow-tooltip />
            <el-table-column prop="created_at" label="创建时间" width="170" />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button link type="success" @click="handleExecute(row)">执行</el-button>
                <el-button link type="primary" @click="showDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty description="暂无 SQL 查询">
                <el-button type="primary" @click="showDialog()">创建第一个查询</el-button>
              </el-empty>
            </template>
          </el-table>
        </template>
      </el-skeleton>
    </el-card>

    <el-card v-if="queryResult !== null" class="result-card">
      <template #header>
        <div class="card-header">
          <span>查询结果</span>
          <el-tag type="info" size="small">{{ queryResult.length }} 条记录</el-tag>
        </div>
      </template>
      <div class="result-table-wrap">
        <el-table :data="queryResult" stripe highlight-current-row max-height="400" v-if="queryResult.length">
          <el-table-column v-for="col in resultColumns" :key="col" :label="col" show-overflow-tooltip min-width="120">
            <template #default="{ row }">
              <template v-if="row[col] !== null && typeof row[col] === 'object'">
                <span class="json-cell">{{ JSON.stringify(row[col]) }}</span>
              </template>
              <template v-else>{{ row[col] }}</template>
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="查询无结果" :image-size="60" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑查询' : '新增查询'" width="520px">
      <el-form :model="form" label-width="60px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="查询名称" />
        </el-form-item>
        <el-form-item label="SQL" required>
          <el-input v-model="form.sql" type="textarea" :rows="8" placeholder="SELECT ..." class="sql-textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saveLoading" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { listSearchDirectSql, createSearchDirectSql, executeSearchDirectSql, updateSearchDirectSql, deleteSearchDirectSql } from '../../api/searchDirectSql'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editing = ref(null)
const form = ref({ name: '', sql: '' })
const queryResult = ref(null)
const saveLoading = ref(false)

const resultColumns = computed(() => {
  if (!queryResult.value?.length) return []
  return Object.keys(queryResult.value[0])
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await listSearchDirectSql({ limit: 1000 })
    list.value = res.data.results
  } finally {
    loading.value = false
  }
}

const showDialog = (row) => {
  editing.value = row || null
  form.value = row ? { name: row.name, sql: row.sql } : { name: '', sql: '' }
  dialogVisible.value = true
}

const handleSave = async () => {
  saveLoading.value = true
  try {
    if (editing.value) {
      await updateSearchDirectSql(editing.value.id, form.value)
    } else {
      await createSearchDirectSql(form.value)
    }
    dialogVisible.value = false
    ElMessage.success('操作成功')
    loadData()
  } catch {
    // error shown by interceptor
  } finally {
    saveLoading.value = false
  }
}

const handleExecute = async (row) => {
  try {
    const res = await executeSearchDirectSql(row.id)
    queryResult.value = Array.isArray(res.data) ? res.data : []
    ElMessage.success(`查询完成，返回 ${queryResult.value.length} 条记录`)
  } catch {
    queryResult.value = null
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm(`确定删除「${row.name}」？`, '确认删除', { type: 'warning' })
  await deleteSearchDirectSql(row.id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.sql-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-card {
  animation: slideUp 0.2s ease;
}

.result-table-wrap {
  overflow-x: auto;
}

.json-cell {
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 12px;
  color: var(--color-text-secondary);
  word-break: break-all;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.sql-textarea :deep(textarea) {
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
}
</style>
