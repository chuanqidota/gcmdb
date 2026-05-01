<template>
  <div class="instance-page">
    <el-row :gutter="20">
      <el-col :span="5">
        <el-card class="model-tree-card">
          <template #header>选择模型</template>
          <el-input v-model="modelSearch" placeholder="搜索模型..." clearable size="small" style="margin-bottom: 12px">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-tree
            :data="treeData"
            :props="{ label: 'name', children: 'children' }"
            :filter-node-method="filterNode"
            node-key="id"
            ref="treeRef"
            highlight-current
            @node-click="onNodeClick"
            class="model-tree"
          />
        </el-card>
      </el-col>
      <el-col :span="19">
        <el-card v-if="!currentModel">
          <el-empty description="请在左侧选择一个模型查看实例" />
        </el-card>
        <el-card v-else>
          <template #header>
            <div class="card-header">
              <div class="card-header-left">
                <span class="card-title">{{ currentModel.name }}</span>
                <el-tag size="small" type="info">{{ currentModel.alias }}</el-tag>
              </div>
              <div class="card-header-right">
                <el-button type="primary" size="small" @click="showCreateDialog">
                  <el-icon><Plus /></el-icon>新增
                </el-button>
                <el-button type="danger" size="small" :disabled="!selectedIds.length" @click="handleMulDelete">
                  <el-icon><Delete /></el-icon>批量删除 ({{ selectedIds.length }})
                </el-button>
              </div>
            </div>
          </template>

          <div class="search-bar">
            <el-input v-model="searchText" placeholder="模糊搜索..." style="width: 200px" clearable @clear="loadInstances" @keyup.enter="loadInstances" size="small">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-select v-model="filterField" placeholder="筛选字段" clearable style="width: 140px" size="small">
              <el-option v-for="f in modelFields" :key="f.alias" :label="f.name" :value="f.alias" />
            </el-select>
            <el-select v-model="filterCompare" placeholder="条件" clearable style="width: 90px" size="small">
              <el-option label="等于" value="=" /><el-option label="包含" value="like" /><el-option label="不等于" value="!=" />
            </el-select>
            <el-input v-model="filterValue" placeholder="值" style="width: 140px" clearable size="small" @keyup.enter="loadInstances" />
            <el-button type="primary" size="small" @click="loadInstances">搜索</el-button>
          </div>

          <el-table :data="instances" stripe v-loading="loading" @selection-change="onSelectionChange" style="margin-top: 12px">
            <el-table-column type="selection" width="45" />
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column v-for="col in tableColumns" :key="col.alias" :label="col.name" show-overflow-tooltip min-width="120">
              <template #default="{ row }">{{ row.data?.[col.alias] ?? '-' }}</template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="170" />
            <el-table-column label="操作" width="130" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="showEditDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="pagination-wrap">
            <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next" @change="loadInstances" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="dialogVisible" :title="editingInstance ? '编辑实例' : '新增实例'" width="600px">
      <el-form :model="instanceForm" label-width="100px" class="instance-form">
        <el-form-item v-for="f in modelFields" :key="f.alias" :label="f.name" :required="f.is_required">
          <el-input v-if="f.type === 'string'" v-model="instanceForm[f.alias]" :placeholder="f.alias" />
          <el-input-number v-else-if="f.type === 'number'" v-model="instanceForm[f.alias]" style="width: 100%" />
          <el-switch v-else-if="f.type === 'bool'" v-model="instanceForm[f.alias]" />
          <el-date-picker v-else-if="f.type === 'date'" v-model="instanceForm[f.alias]" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
          <el-date-picker v-else-if="f.type === 'datetime'" v-model="instanceForm[f.alias]" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
          <el-input v-else v-model="instanceForm[f.alias]" type="textarea" :rows="3" placeholder="JSON 格式" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, Search } from '@element-plus/icons-vue'
import { listModel, retrieveModel } from '../../api/model'
import { listModelGroup } from '../../api/modelGroup'
import { listInstance, createInstance, updateInstance, deleteInstance } from '../../api/instance'

const treeRef = ref(null)
const modelSearch = ref('')
const treeData = ref([])
const currentModel = ref(null)
const modelFields = ref([])
const instances = ref([])
const loading = ref(false)
const searchText = ref('')
const filterField = ref('')
const filterCompare = ref('')
const filterValue = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedIds = ref([])
const dialogVisible = ref(false)
const editingInstance = ref(null)
const instanceForm = ref({})

const tableColumns = computed(() => modelFields.value.slice(0, 6))

const buildTree = async () => {
  const [gRes, mRes] = await Promise.all([
    listModelGroup({ limit: 1000 }),
    listModel({ limit: 1000 }),
  ])
  const groups = gRes.data.results
  const models = mRes.data.results
  treeData.value = groups.map((g) => ({
    id: `g-${g.id}`,
    name: g.name,
    children: models.filter((m) => m.group_id === g.id).map((m) => ({ id: m.id, name: m.name, alias: m.alias, isModel: true })),
  }))
}

const filterNode = (value, data) => {
  if (!value) return true
  return data.name.includes(value)
}

watch(modelSearch, (val) => { treeRef.value?.filter(val) })

const onNodeClick = async (data) => {
  if (!data.isModel) return
  currentModel.value = data
  const res = await retrieveModel(data.id)
  modelFields.value = (res.data.groups || []).flatMap((g) => g.fields || [])
  page.value = 1
  loadInstances()
}

const loadInstances = async () => {
  if (!currentModel.value) return
  loading.value = true
  try {
    const params = { limit: pageSize.value, offset: (page.value - 1) * pageSize.value }
    if (searchText.value) params.search = searchText.value
    if (filterField.value && filterCompare.value && filterValue.value) {
      params.field = filterField.value
      params.compare = filterCompare.value
      params.value = filterValue.value
    }
    const res = await listInstance(currentModel.value.id, params)
    instances.value = res.data.results
    total.value = res.data.count
  } finally {
    loading.value = false
  }
}

const onSelectionChange = (rows) => { selectedIds.value = rows.map((r) => r.id) }

const showCreateDialog = () => {
  editingInstance.value = null
  const form = {}
  modelFields.value.forEach((f) => { form[f.alias] = f.type === 'bool' ? false : f.type === 'number' ? 0 : '' })
  instanceForm.value = form
  dialogVisible.value = true
}

const showEditDialog = (row) => {
  editingInstance.value = row
  instanceForm.value = { ...row.data }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (editingInstance.value) {
    await updateInstance(editingInstance.value.id, { data: instanceForm.value })
  } else {
    await createInstance({ model_id: currentModel.value.id, model_alias: currentModel.value.alias, data: instanceForm.value })
  }
  dialogVisible.value = false
  ElMessage.success('操作成功')
  loadInstances()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该实例？', '确认删除', { type: 'warning' })
  await deleteInstance(row.id)
  ElMessage.success('删除成功')
  loadInstances()
}

const handleMulDelete = async () => {
  await ElMessageBox.confirm(`确定删除选中的 ${selectedIds.value.length} 条实例？`, '确认批量删除', { type: 'warning' })
  const results = await Promise.allSettled(selectedIds.value.map((id) => deleteInstance(id)))
  const failed = results.filter(r => r.status === 'rejected')
  if (failed.length) {
    ElMessage.warning(`${failed.length} 条删除失败`)
  } else {
    ElMessage.success('批量删除成功')
  }
  loadInstances()
}

onMounted(buildTree)
</script>

<style scoped>
.instance-page {
  height: calc(100vh - 140px);
}

.model-tree-card {
  height: 100%;
}

.model-tree-card :deep(.el-card__body) {
  padding: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-header-right {
  display: flex;
  gap: 8px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
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

.instance-form :deep(.el-form-item) {
  margin-bottom: 18px;
}
</style>
