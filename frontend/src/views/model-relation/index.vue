<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>模型关系</span>
        <el-button type="primary" @click="showDialog">
          <el-icon><Plus /></el-icon>新增关系
        </el-button>
      </div>
    </template>
    <el-table :data="list" stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="source_display" label="源模型" />
      <el-table-column prop="target_display" label="目标模型" />
      <el-table-column prop="type_display" label="关系类型" width="120" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column prop="created_at" label="创建时间" width="170" />
      <el-table-column label="操作" width="80" fixed="right">
        <template #default="{ row }">
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="dialogVisible" title="新增模型关系" width="480px">
    <el-form :model="form" label-width="70px">
      <el-form-item label="源模型" required>
        <el-select v-model="form.source_id" placeholder="选择源模型" style="width: 100%" filterable>
          <el-option v-for="m in sourceModelOptions" :key="m.id" :label="`${m.name} (${m.alias})`" :value="m.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="目标模型" required>
        <el-select v-model="form.target_id" placeholder="选择目标模型" style="width: 100%" filterable>
          <el-option v-for="m in targetModelOptions" :key="m.id" :label="`${m.name} (${m.alias})`" :value="m.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="关系类型" required>
        <el-select v-model="form.type_id" placeholder="选择关系类型" style="width: 100%">
          <el-option v-for="t in typeOptions" :key="t.id" :label="`${t.name} (${t.s2t})`" :value="t.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="3" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSave">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { listModelRelation, createModelRelation, deleteModelRelation } from '../../api/modelRelation'
import { listModel } from '../../api/model'
import { listModelRelationType } from '../../api/modelRelationType'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const modelOptions = ref([])
const typeOptions = ref([])
const form = ref({ source_id: null, target_id: null, type_id: null, description: '' })

const sourceModelOptions = computed(() => {
  const targetId = form.value.target_id
  if (!targetId) return modelOptions.value
  const existingSourceIds = new Set(list.value.filter(r => r.target_id === targetId).map(r => r.source_id))
  return modelOptions.value.filter(m => m.id !== targetId && !existingSourceIds.has(m.id))
})

const targetModelOptions = computed(() => {
  const sourceId = form.value.source_id
  if (!sourceId) return modelOptions.value
  const existingTargetIds = new Set(list.value.filter(r => r.source_id === sourceId).map(r => r.target_id))
  return modelOptions.value.filter(m => m.id !== sourceId && !existingTargetIds.has(m.id))
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await listModelRelation({ limit: 1000 })
    list.value = res.data.results
  } finally {
    loading.value = false
  }
}

const loadOptions = async () => {
  const [m, t] = await Promise.all([
    listModel({ limit: 1000 }),
    listModelRelationType({ limit: 1000 }),
  ])
  modelOptions.value = m.data.results
  typeOptions.value = t.data.results
}

const showDialog = () => {
  form.value = { source_id: null, target_id: null, type_id: null, description: '' }
  dialogVisible.value = true
}

const handleSave = async () => {
  await createModelRelation(form.value)
  dialogVisible.value = false
  ElMessage.success('创建成功')
  loadData()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该关系？', '确认删除', { type: 'warning' })
  await deleteModelRelation(row.id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(() => { loadData(); loadOptions() })
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
