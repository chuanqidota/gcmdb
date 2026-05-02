<template>
  <div class="instance-page">
    <el-row :gutter="20">
      <el-col :span="treeCollapsed ? 1 : 5">
        <el-card class="model-tree-card" :class="{ collapsed: treeCollapsed }">
          <template #header>
            <div class="tree-header">
              <span v-show="!treeCollapsed">选择模型</span>
              <el-icon class="tree-collapse-btn" @click="treeCollapsed = !treeCollapsed">
                <Fold v-if="!treeCollapsed" /><Expand v-else />
              </el-icon>
            </div>
          </template>
          <template v-if="!treeCollapsed">
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
          </template>
        </el-card>
      </el-col>
      <el-col :span="treeCollapsed ? 23 : 19">
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

          <el-table :data="instances" stripe v-loading="loading" @selection-change="onSelectionChange" @expand-change="onExpandChange" style="margin-top: 12px" row-key="id">
            <el-table-column type="expand" width="45">
              <template #default="{ row }">
                <div class="relation-panel" v-loading="relationCache[row.id]?.loading">

                  <template v-if="relationCache[row.id]?.groups?.length">
                    <div class="relation-panel-header">
                      <el-tabs v-model="relationCache[row.id].activeTab" @tab-change="(id) => onTabChange(row, Number(id))" class="relation-tabs-inline">
                        <el-tab-pane v-for="g in relationCache[row.id].groups" :key="g.model_id" :name="String(g.model_id)">
                          <template #label>
                            <span class="tab-label">
                              <span :class="['tab-direction', g.direction]">{{ g.direction === 'source' ? '→' : '←' }}</span>
                              {{ g.model_display }}
                              <el-badge :value="g.count" class="tab-badge" />
                            </span>
                          </template>
                        </el-tab-pane>
                      </el-tabs>
                      <div class="relation-panel-actions">
                        <el-button size="small" type="primary" @click="showRelationDialog(row, Number(relationCache[row.id].activeTab))"><el-icon><Plus /></el-icon>添加关系</el-button>
                        <el-button size="small" type="danger" :disabled="!relationCache[row.id]?.tabData?.[Number(relationCache[row.id].activeTab)]?.selectedIds?.length" @click="handleBatchDeleteRelations(row, relationCache[row.id].groups.find(g => String(g.model_id) === relationCache[row.id].activeTab))">批量删除 ({{ relationCache[row.id]?.tabData?.[Number(relationCache[row.id].activeTab)]?.selectedIds?.length || 0 }})</el-button>
                      </div>
                    </div>

                    <!-- 当前 Tab 的表格内容 -->
                    <template v-for="g in relationCache[row.id].groups" :key="'content-' + g.model_id">
                      <div v-if="String(g.model_id) === relationCache[row.id].activeTab">
                        <el-table
                          :data="getPaginatedInstances(row.id, g.model_id)"
                          size="small"
                          stripe
                          row-key="id"
                          v-loading="relationCache[row.id]?.tabData?.[g.model_id]?.loading"
                        >
                          <el-table-column width="40">
                            <template #header>
                              <el-checkbox
                                :model-value="isAllSelected(row.id, g.model_id)"
                                :indeterminate="isIndeterminate(row.id, g.model_id)"
                                @change="(val) => toggleSelectAll(row.id, g.model_id, val)"
                              />
                            </template>
                            <template #default="{ row: inst }">
                              <el-checkbox
                                :model-value="isInstanceSelected(row.id, g.model_id, inst.id)"
                                @change="(val) => toggleInstanceSelection(row.id, g.model_id, inst.id, val)"
                              />
                            </template>
                          </el-table-column>
                          <el-table-column prop="id" label="ID" width="60" />
                          <el-table-column
                            v-for="f in (relationCache[row.id]?.tabData?.[g.model_id]?.fields || [])"
                            :key="f.alias"
                            :label="f.name"
                            show-overflow-tooltip
                            min-width="100"
                          >
                            <template #default="{ row: inst }">{{ inst.data?.[f.alias] ?? '-' }}</template>
                          </el-table-column>
                          <el-table-column prop="created_at" label="创建时间" width="170" />
                          <el-table-column label="操作" width="80" fixed="right">
                            <template #default="{ row: inst }">
                              <el-button link type="danger" size="small" @click="handleDeleteSingleRelation(row, g, inst.id)">删除</el-button>
                            </template>
                          </el-table-column>
                        </el-table>

                        <div class="relation-pagination">
                          <el-pagination
                            v-model:current-page="relationCache[row.id].tabData[g.model_id].page"
                            v-model:page-size="relationCache[row.id].tabData[g.model_id].pageSize"
                            :total="relationCache[row.id].tabData[g.model_id].instances.length"
                            :page-sizes="[10, 20, 50]"
                            layout="total, sizes, prev, pager, next"
                            small
                            background
                          />
                        </div>
                      </div>
                    </template>
                  </template>

                  <template v-else-if="!relationCache[row.id]?.loading">
                    <div class="relation-empty">未定义模型关系，请先在模型管理中配置</div>
                  </template>
                </div>
              </template>
            </el-table-column>
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

    <!-- 新建实例关系对话框 -->
    <el-dialog v-model="relDialogVisible" title="新建实例关系" width="420px">
      <el-form :model="relForm" label-width="80px">
        <el-form-item label="当前实例">
          <el-input :model-value="relForm.source_display" disabled />
        </el-form-item>
        <el-form-item label="关联模型">
          <el-input :model-value="relForm.related_model_display" disabled />
        </el-form-item>
        <el-form-item label="选择实例" required>
          <el-select v-model="relForm.target_instance_ids" placeholder="选择要关联的实例" style="width: 100%" multiple filterable :loading="targetInstancesLoading">
            <el-option v-for="inst in targetInstances" :key="inst.id" :label="`#${inst.id} ${getInstanceLabel(inst)}`" :value="inst.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="relDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateRelation" :loading="relSaving">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogVisible" :title="editingInstance ? '编辑实例' : '新增实例'" width="520px">
      <el-form :model="instanceForm" label-width="100px" class="instance-form">
        <el-form-item v-for="f in modelFields" :key="f.alias" :label="f.name" :required="f.is_required">
          <el-input v-if="f.type === 'string'" v-model="instanceForm[f.alias]" :placeholder="f.alias" />
          <el-input-number v-else-if="f.type === 'number'" v-model="instanceForm[f.alias]" style="width: 100%" />
          <el-switch v-else-if="f.type === 'bool'" v-model="instanceForm[f.alias]" />
          <el-date-picker v-else-if="f.type === 'date'" v-model="instanceForm[f.alias]" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
          <el-date-picker v-else-if="f.type === 'datetime'" v-model="instanceForm[f.alias]" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
          <el-select v-else-if="f.type === 'enum'" v-model="instanceForm[f.alias]" style="width: 100%">
            <el-option v-for="opt in parseOptions(f.options)" :key="opt" :label="opt" :value="opt" />
          </el-select>
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
import { Plus, Delete, Search, Fold, Expand } from '@element-plus/icons-vue'
import { listModel, retrieveModel } from '../../api/model'
import { listModelGroup } from '../../api/modelGroup'
import { listModelRelation } from '../../api/modelRelation'
import { listInstance, retrieveInstance, createInstance, updateInstance, deleteInstance } from '../../api/instance'
import { sourceTargetRelation, targetSourceRelation, createInstanceRelation, deleteInstanceRelationByKeys } from '../../api/instanceRelation'

const treeRef = ref(null)
const treeCollapsed = ref(false)
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

// ===== 实例关系 =====
// relationCache[instanceId] = { groups, loading, activeTab, tabData: { [modelId]: { fields, instances, loading, page, selectedIds } } }
const relationCache = ref({})
const relDialogVisible = ref(false)
const relSaving = ref(false)
const relForm = ref({ source_instance_id: null, source_display: '', related_model_id: null, related_model_display: '', target_instance_ids: [] })
const targetInstances = ref([])
const targetInstancesLoading = ref(false)
const allModels = ref([])
const modelRelations = ref([]) // 当前模型的模型关系列表

const tableColumns = computed(() => modelFields.value.slice(0, 6))

const parseOptions = (options) => {
  if (!options) return []
  try { return JSON.parse(options) } catch { return [] }
}


const filterNode = (value, data) => {
  if (!value) return true
  return data.name.includes(value)
}

watch(modelSearch, (val) => { treeRef.value?.filter(val) })

const onNodeClick = async (data) => {
  if (!data.isModel) return
  currentModel.value = data
  const [modelRes, relRes] = await Promise.all([
    retrieveModel(data.id),
    listModelRelation({ source_id: data.id }),
  ])
  modelFields.value = (modelRes.data.groups || []).flatMap((g) => g.fields || [])
  modelRelations.value = relRes.data.results || []
  page.value = 1
  relationCache.value = {}
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
  try {
    if (editingInstance.value) {
      await updateInstance(editingInstance.value.id, { data: instanceForm.value })
    } else {
      await createInstance({ model_id: currentModel.value.id, data: instanceForm.value })
    }
    dialogVisible.value = false
    ElMessage.success('操作成功')
    loadInstances()
  } catch {
    // error shown by interceptor
  }
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

// ===== 实例关系函数 =====
const onExpandChange = (row, expanded) => {
  if (expanded.length && !relationCache.value[row.id]) {
    loadRelations(row)
  }
}

const loadRelations = async (row) => {
  const existing = relationCache.value[row.id]
  const activeTab = existing?.activeTab || ''
  const tabData = { ...(existing?.tabData || {}) }
  if (activeTab && tabData[activeTab]) {
    delete tabData[activeTab]
  }
  relationCache.value[row.id] = { groups: [], loading: true, activeTab, tabData }
  try {
    const [srcRes, tgtRes] = await Promise.all([
      sourceTargetRelation(row.id),
      targetSourceRelation(row.id),
    ])

    // 收集已有实例关系的模型分组
    const instanceGroupMap = new Map()
    for (const g of (srcRes.data || [])) {
      instanceGroupMap.set(g.model_id, { direction: 'source', ...g, instanceIds: g.instances })
    }
    for (const g of (tgtRes.data || [])) {
      instanceGroupMap.set(g.model_id, { direction: 'target', ...g, instanceIds: g.instances })
    }

    // 合并模型关系定义中的所有关联模型
    const groups = []
    for (const mr of modelRelations.value) {
      const isSource = mr.source_id === currentModel.value.id
      const relatedModelId = isSource ? mr.target_id : mr.source_id
      const direction = isSource ? 'source' : 'target'
      const existingGroup = instanceGroupMap.get(relatedModelId)
      if (existingGroup) {
        groups.push(existingGroup)
      } else {
        const relatedModel = allModels.value.find(m => m.id === relatedModelId)
        groups.push({
          direction,
          model_id: relatedModelId,
          model_display: relatedModel ? relatedModel.name : `模型#${relatedModelId}`,
          count: 0,
          instanceIds: [],
        })
      }
    }

    relationCache.value[row.id].groups = groups
    relationCache.value[row.id].loading = false

    if (groups.length) {
      const tabToLoad = groups.find(g => String(g.model_id) === activeTab) ? activeTab : String(groups[0].model_id)
      relationCache.value[row.id].activeTab = tabToLoad
      onTabChange(row, Number(tabToLoad))
    }
  } catch {
    relationCache.value[row.id] = { groups: [], loading: false, activeTab: '', tabData: {} }
  }
}

const onTabChange = async (row, modelId) => {
  const cache = relationCache.value[row.id]
  if (!cache || cache.tabData[modelId]) return  // 已缓存

  const group = cache.groups.find(g => g.model_id === modelId)
  if (!group) return

  cache.tabData[modelId] = { fields: [], instances: [], loading: true, page: 1, pageSize: 10, selectedIds: [] }
  relationCache.value = { ...relationCache.value }  // 触发响应式更新

  try {
    // 获取模型字段定义
    const modelRes = await retrieveModel(modelId)
    const fields = (modelRes.data.groups || []).flatMap(g => g.fields || [])
    cache.tabData[modelId].fields = fields.slice(0, 8)

    // 并行获取所有实例详情
    const instResults = await Promise.allSettled(
      group.instanceIds.map(id => retrieveInstance(id))
    )
    cache.tabData[modelId].instances = instResults
      .filter(r => r.status === 'fulfilled')
      .map(r => r.value.data)
  } catch {
    // ignore
  } finally {
    cache.tabData[modelId].loading = false
    relationCache.value = { ...relationCache.value }
  }
}

const getPaginatedInstances = (instanceId, modelId) => {
  const tabData = relationCache.value[instanceId]?.tabData?.[modelId]
  if (!tabData) return []
  const size = tabData.pageSize || 10
  const start = (tabData.page - 1) * size
  return tabData.instances.slice(start, start + size)
}

const getInstanceLabel = (inst) => {
  if (!inst.data) return ''
  const keys = Object.keys(inst.data).slice(0, 2)
  return keys.map(k => inst.data[k]).filter(Boolean).join(', ')
}

// 查找当前模型与目标模型的关系定义，确定源/目标方向
const findModelRelation = (relatedModelId) => {
  return modelRelations.value.find(r =>
    (r.source_id === currentModel.value.id && r.target_id === relatedModelId) ||
    (r.target_id === currentModel.value.id && r.source_id === relatedModelId)
  )
}

const getRelatedInstanceIds = (instanceId, modelId) => {
  const cache = relationCache.value[instanceId]
  if (!cache) return new Set()
  const ids = new Set()
  for (const g of cache.groups) {
    if (g.model_id === modelId) {
      g.instanceIds.forEach(id => ids.add(id))
    }
  }
  return ids
}

const showRelationDialog = async (row, relatedModelId) => {
  const mr = findModelRelation(relatedModelId)
  const relatedModel = allModels.value.find(m => m.id === relatedModelId)
  const relatedDisplay = relatedModel ? `${relatedModel.name} (${relatedModel.alias})` : `模型#${relatedModelId}`
  relForm.value = {
    source_instance_id: row.id,
    source_display: `#${row.id} (模型: ${currentModel.value.name})`,
    related_model_id: relatedModelId,
    related_model_display: relatedDisplay,
    model_relation: mr,
    target_instance_ids: [],
  }
  targetInstances.value = []
  relDialogVisible.value = true
  // 加载目标模型的实例，过滤已关联的
  targetInstancesLoading.value = true
  try {
    const res = await listInstance(relatedModelId, { limit: 1000 })
    const all = res.data.results || []
    const relatedIds = getRelatedInstanceIds(row.id, relatedModelId)
    targetInstances.value = all.filter(inst => !relatedIds.has(inst.id))
  } catch {
    targetInstances.value = []
  } finally {
    targetInstancesLoading.value = false
  }
}

const handleCreateRelation = async () => {
  if (!relForm.value.related_model_id || !relForm.value.target_instance_ids.length) return
  const mr = relForm.value.model_relation
  if (!mr) { ElMessage.error('未找到模型关系定义'); return }
  relSaving.value = true
  try {
    const currentId = currentModel.value.id
    const relatedId = relForm.value.related_model_id
    const isCurrentSource = mr.source_id === currentId
    const sourceModelId = isCurrentSource ? currentId : relatedId
    const targetModelId = isCurrentSource ? relatedId : currentId
    const sourceInstanceId = relForm.value.source_instance_id
    const results = await Promise.allSettled(
      relForm.value.target_instance_ids.map(targetId =>
        createInstanceRelation({
          source_model_id: sourceModelId,
          target_model_id: targetModelId,
          source_instance_id: isCurrentSource ? sourceInstanceId : targetId,
          target_instance_id: isCurrentSource ? targetId : sourceInstanceId,
        })
      )
    )
    const failed = results.filter(r => r.status === 'rejected')
    if (failed.length) {
      ElMessage.warning(`${failed.length} 条创建失败`)
    } else {
      ElMessage.success('创建成功')
    }
    relDialogVisible.value = false
    const row = instances.value.find(i => i.id === sourceInstanceId)
    if (row) loadRelations(row)
  } catch {
    ElMessage.error('创建失败')
  } finally {
    relSaving.value = false
  }
}

// ===== 手动复选框管理 =====
const isInstanceSelected = (instanceId, modelId, instId) => {
  return relationCache.value[instanceId]?.tabData?.[modelId]?.selectedIds?.includes(instId) || false
}

const isAllSelected = (instanceId, modelId) => {
  const tabData = relationCache.value[instanceId]?.tabData?.[modelId]
  if (!tabData) return false
  const pageInstances = getPaginatedInstances(instanceId, modelId)
  return pageInstances.length > 0 && pageInstances.every(inst => tabData.selectedIds.includes(inst.id))
}

const isIndeterminate = (instanceId, modelId) => {
  const tabData = relationCache.value[instanceId]?.tabData?.[modelId]
  if (!tabData) return false
  const pageInstances = getPaginatedInstances(instanceId, modelId)
  const selected = pageInstances.filter(inst => tabData.selectedIds.includes(inst.id))
  return selected.length > 0 && selected.length < pageInstances.length
}

const toggleInstanceSelection = (instanceId, modelId, instId, checked) => {
  const tabData = relationCache.value[instanceId]?.tabData?.[modelId]
  if (!tabData) return
  if (checked) {
    if (!tabData.selectedIds.includes(instId)) {
      tabData.selectedIds = [...tabData.selectedIds, instId]
    }
  } else {
    tabData.selectedIds = tabData.selectedIds.filter(id => id !== instId)
  }
}

const toggleSelectAll = (instanceId, modelId, checked) => {
  const tabData = relationCache.value[instanceId]?.tabData?.[modelId]
  if (!tabData) return
  const pageInstances = getPaginatedInstances(instanceId, modelId)
  if (checked) {
    const newIds = pageInstances.map(inst => inst.id).filter(id => !tabData.selectedIds.includes(id))
    tabData.selectedIds = [...tabData.selectedIds, ...newIds]
  } else {
    const pageIds = new Set(pageInstances.map(inst => inst.id))
    tabData.selectedIds = tabData.selectedIds.filter(id => !pageIds.has(id))
  }
}

const removeRelationFromCache = (row, group, ids) => {
  const cache = relationCache.value[row.id]
  const tabData = cache.tabData[group.model_id]
  if (tabData) {
    tabData.instances = tabData.instances.filter(inst => !ids.includes(inst.id))
    tabData.selectedIds = (tabData.selectedIds || []).filter(id => !ids.includes(id))
  }
  group.instanceIds = group.instanceIds.filter(id => !ids.includes(id))
  group.count = group.instanceIds.length
  if (group.count === 0) {
    cache.groups = cache.groups.filter(g => g.model_id !== group.model_id)
    delete cache.tabData[group.model_id]
    if (cache.groups.length) {
      cache.activeTab = String(cache.groups[0].model_id)
    }
  }
  relationCache.value = { ...relationCache.value }
}

const handleDeleteSingleRelation = async (row, group, targetInstId) => {
  await ElMessageBox.confirm(`确定删除实例 #${row.id} 与 #${targetInstId} 的关系？`, '确认删除', { type: 'warning' })
  const sourceModelId = group.direction === 'source' ? currentModel.value.id : group.model_id
  const targetModelId = group.direction === 'source' ? group.model_id : currentModel.value.id
  const sourceId = group.direction === 'source' ? row.id : targetInstId
  const targetId = group.direction === 'source' ? targetInstId : row.id
  await deleteInstanceRelationByKeys({ source_model_id: sourceModelId, target_model_id: targetModelId, source_id: sourceId, target_id: targetId })
  ElMessage.success('删除成功')
  removeRelationFromCache(row, group, [targetInstId])
}

const handleBatchDeleteRelations = async (row, group) => {
  const tabData = relationCache.value[row.id]?.tabData?.[group.model_id]
  const ids = tabData?.selectedIds || []
  if (!ids.length) return
  await ElMessageBox.confirm(`确定删除选中的 ${ids.length} 条关系？`, '确认删除', { type: 'warning' })
  for (const targetInstId of ids) {
    const sourceModelId = group.direction === 'source' ? currentModel.value.id : group.model_id
    const targetModelId = group.direction === 'source' ? group.model_id : currentModel.value.id
    const sourceId = group.direction === 'source' ? row.id : targetInstId
    const targetId = group.direction === 'source' ? targetInstId : row.id
    try {
      await deleteInstanceRelationByKeys({ source_model_id: sourceModelId, target_model_id: targetModelId, source_id: sourceId, target_id: targetId })
    } catch { /* skip */ }
  }
  ElMessage.success('批量删除成功')
  removeRelationFromCache(row, group, ids)
}

// ===== 初始化 =====
const buildTree = async () => {
  const [gRes, mRes] = await Promise.all([
    listModelGroup({ limit: 1000 }),
    listModel({ limit: 1000 }),
  ])
  allModels.value = mRes.data.results || []
  const groups = gRes.data.results
  treeData.value = groups.map((g) => ({
    id: `g-${g.id}`,
    name: g.name,
    children: allModels.value.filter((m) => m.group_id === g.id).map((m) => ({ id: m.id, name: m.name, alias: m.alias, isModel: true })),
  }))
}

onMounted(buildTree)
</script>

<style scoped>
.instance-page {
  height: var(--page-height);
}

.model-tree-card {
  height: 100%;
  transition: all 0.2s ease;
}

.model-tree-card.collapsed :deep(.el-card__body) {
  padding: 0;
  display: none;
}

.model-tree-card :deep(.el-card__body) {
  padding: 12px;
}

.tree-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tree-collapse-btn {
  cursor: pointer;
  font-size: 16px;
  color: var(--color-text-secondary);
  padding: 2px;
  border-radius: 4px;
  transition: all 0.15s ease;
}

.tree-collapse-btn:hover {
  background: var(--color-muted);
  color: var(--color-primary);
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

/* 实例关系面板 */
.relation-panel {
  padding: 12px 20px;
  background: var(--color-muted);
  border-radius: var(--radius-md);
  margin: 4px 0;
}

.relation-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.relation-panel-header :deep(.el-tabs__header) {
  margin: 0;
}

.relation-tabs-inline {
  flex: 1;
  min-width: 0;
}

.relation-panel-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  margin-left: 16px;
}

.tab-label {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.tab-direction {
  font-weight: 700;
  font-size: 13px;
}

.tab-direction.source {
  color: var(--color-primary);
}

.tab-direction.target {
  color: var(--color-warning);
}

.tab-badge {
  margin-left: 4px;
}

.tab-badge :deep(.el-badge__content) {
  font-size: 10px;
}

.relation-pagination {
  display: flex;
  justify-content: center;
  margin-top: 8px;
}

.relation-empty {
  text-align: center;
  padding: 24px;
  font-size: 13px;
  color: var(--color-text-muted);
}
</style>
