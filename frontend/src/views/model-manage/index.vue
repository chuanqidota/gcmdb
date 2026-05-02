<template>
  <div class="model-manage-page">
    <el-row :gutter="16" class="main-row">
      <!-- 左侧：分组列表 -->
      <el-col :span="5" class="left-panel">
        <el-card class="group-card">
          <template #header>
            <div class="card-header">
              <span>模型分组</span>
              <el-button type="primary" size="small" @click="showGroupDialog()">
                <el-icon><Plus /></el-icon>
              </el-button>
            </div>
          </template>
          <div class="group-list">
            <div class="group-item" :class="{ active: selectedGroupId === null }" @click="selectedGroupId = null">
              <span class="group-item-name">全部模型</span>
            </div>
            <div
              v-for="g in groups" :key="g.id"
              class="group-item"
              :class="{ active: selectedGroupId === g.id }"
              @click="selectedGroupId = g.id"
            >
              <span class="group-item-name">{{ g.name }}</span>
              <div class="group-item-actions">
                <el-button link size="small" @click.stop="showGroupDialog(g)"><el-icon><Edit /></el-icon></el-button>
                <el-button link size="small" type="danger" @click.stop="handleDeleteGroup(g)"><el-icon><Delete /></el-icon></el-button>
              </div>
            </div>
            <div
              v-if="ungroupedCount"
              class="group-item"
              :class="{ active: selectedGroupId === 'ungrouped' }"
              @click="selectedGroupId = 'ungrouped'"
            >
              <span class="group-item-name">未分组</span>
              <span class="group-item-count">{{ ungroupedCount }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：模型列表 + 详情展开 -->
      <el-col :span="19" class="right-panel">
        <!-- 模型表格 -->
        <el-card>
          <template #header>
            <div class="card-header">
              <div class="header-left">
                <span class="header-title">{{ currentGroupName }}</span>
                <el-input v-model="search" placeholder="搜索模型..." style="width: 240px; margin-left: 16px" clearable @clear="loadModels" @keyup.enter="loadModels">
                  <template #prefix><el-icon><Search /></el-icon></template>
                </el-input>
              </div>
              <el-button type="primary" @click="showModelDialog()">
                <el-icon><Plus /></el-icon>新增模型
              </el-button>
            </div>
          </template>
          <el-table :data="models" stripe v-loading="modelsLoading" @row-click="onRowClick" row-class-name="model-row">
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="alias" label="别名" width="130" />
            <el-table-column prop="name" label="名称" width="140" />
            <el-table-column prop="description" label="描述" show-overflow-tooltip />
            <el-table-column prop="is_usable" label="状态" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="row.is_usable ? 'success' : 'danger'" size="small" effect="light">{{ row.is_usable ? '启用' : '停用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click.stop="toggleDetail(row)">{{ expandedModelId === row.id ? '收起' : '详情' }}</el-button>
                <el-button link type="primary" @click.stop="showModelDialog(row)">编辑</el-button>
                <el-button link type="primary" @click.stop="showMoveGroupDialog(row)">移组</el-button>
                <el-button link type="danger" @click.stop="handleDeleteModel(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="pagination-wrap">
            <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next" @change="loadModels" />
          </div>
        </el-card>

        <!-- 详情展开面板 -->
        <el-card v-if="expandedModelId" class="detail-card">
          <template #header>
            <div class="card-header">
              <div class="detail-header">
                <span class="detail-title">{{ expandedModel.name }}</span>
                <el-tag size="small" type="info">{{ expandedModel.alias }}</el-tag>
              </div>
              <el-button text @click="expandedModelId = null"><el-icon><Close /></el-icon></el-button>
            </div>
          </template>
          <el-tabs v-model="detailTab">
            <!-- 字段管理 Tab -->
            <el-tab-pane label="字段" name="fields">
              <el-row :gutter="16">
                <el-col :span="6">
                  <div class="field-group-list">
                    <div class="field-group-header">
                      <span>字段分组</span>
                      <el-button size="small" type="primary" text @click="showFieldGroupDialog()"><el-icon><Plus /></el-icon></el-button>
                    </div>
                    <div class="group-item" :class="{ active: !selectedFieldGroup }" @click="selectedFieldGroup = null">
                      <span class="group-item-name">全部字段</span>
                    </div>
                    <div
                      v-for="g in detailData.fieldGroups" :key="g.id"
                      class="group-item"
                      :class="{ active: selectedFieldGroup === g.id }"
                      @click="selectedFieldGroup = g.id"
                    >
                      <span class="group-item-name">{{ g.name }}</span>
                      <div class="group-item-actions">
                        <el-button link size="small" @click.stop="showFieldGroupDialog(g)"><el-icon><Edit /></el-icon></el-button>
                        <el-button link size="small" type="danger" @click.stop="handleDeleteFieldGroup(g)"><el-icon><Delete /></el-icon></el-button>
                      </div>
                    </div>
                  </div>
                </el-col>
                <el-col :span="18">
                  <div class="card-header section-header">
                    <span>字段列表</span>
                    <el-button size="small" type="primary" @click="showFieldDialog()"><el-icon><Plus /></el-icon>新增字段</el-button>
                  </div>
                  <el-table :data="currentFields" stripe size="small">
                    <el-table-column prop="id" label="ID" width="50" />
                    <el-table-column prop="alias" label="别名" width="120" />
                    <el-table-column prop="name" label="名称" width="120" />
                    <el-table-column prop="type" label="类型" width="80">
                      <template #default="{ row }"><el-tag size="small" effect="plain">{{ row.type }}</el-tag></template>
                    </el-table-column>
                    <el-table-column prop="is_required" label="必填" width="60" align="center">
                      <template #default="{ row }">
                        <el-icon v-if="row.is_required" color="var(--color-destructive)"><CircleCheckFilled /></el-icon>
                        <span v-else>-</span>
                      </template>
                    </el-table-column>
                    <el-table-column prop="order" label="排序" width="60" />
                    <el-table-column label="操作" width="120" fixed="right">
                      <template #default="{ row }">
                        <el-button link type="primary" size="small" @click="showFieldDialog(row)">编辑</el-button>
                        <el-button link type="danger" size="small" @click="handleDeleteField(row)">删除</el-button>
                      </template>
                    </el-table-column>
                  </el-table>
                </el-col>
              </el-row>
            </el-tab-pane>

            <!-- 唯一约束 Tab -->
            <el-tab-pane label="唯一约束" name="uniques">
              <div class="card-header section-header">
                <span>唯一约束</span>
                <el-button size="small" type="primary" @click="showUniqueDialog()"><el-icon><Plus /></el-icon>新增</el-button>
              </div>
              <el-table :data="detailData.uniques" stripe size="small" v-if="detailData.uniques.length">
                <el-table-column prop="id" label="ID" width="60" />
                <el-table-column prop="fields" label="字段组合">
                  <template #default="{ row }">
                    <el-tag v-for="f in (row.fields || '').split(',')" :key="f" size="small" style="margin-right: 4px">{{ f }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="description" label="描述" />
                <el-table-column label="操作" width="80">
                  <template #default="{ row }">
                    <el-button link type="danger" size="small" @click="handleDeleteUnique(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-empty v-else description="暂无唯一约束" :image-size="48" />
            </el-tab-pane>

            <!-- 模型关系 Tab -->
            <el-tab-pane label="模型关系" name="relations">
              <div class="card-header section-header">
                <span>模型关系</span>
                <el-button size="small" type="primary" @click="showRelationDialog()"><el-icon><Plus /></el-icon>新增关系</el-button>
              </div>
              <el-table :data="detailData.relations" stripe size="small" v-loading="detailData.relationsLoading">
                <el-table-column prop="id" label="ID" width="60" />
                <el-table-column prop="source_display" label="源模型" />
                <el-table-column prop="target_display" label="目标模型" />
                <el-table-column prop="type_display" label="关系类型" width="120" />
                <el-table-column prop="description" label="描述" show-overflow-tooltip />
                <el-table-column label="操作" width="80">
                  <template #default="{ row }">
                    <el-button link type="danger" size="small" @click="handleDeleteRelation(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>

            <!-- 关系类型 Tab -->
            <el-tab-pane label="关系类型" name="relationTypes">
              <div class="card-header section-header">
                <span>关系类型</span>
                <el-button size="small" type="primary" @click="showRtDialog()"><el-icon><Plus /></el-icon>新增</el-button>
              </div>
              <el-table :data="relationTypes" stripe size="small">
                <el-table-column prop="id" label="ID" width="60" />
                <el-table-column prop="name" label="名称" width="140" />
                <el-table-column prop="s2t" label="源 → 目标" />
                <el-table-column prop="t2s" label="目标 → 源" />
                <el-table-column prop="created_at" label="创建时间" width="170" />
                <el-table-column label="操作" width="120">
                  <template #default="{ row }">
                    <el-button link type="primary" size="small" @click="showRtDialog(row)">编辑</el-button>
                    <el-button link type="danger" size="small" @click="handleRtDelete(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>

    <!-- 分组对话框 -->
    <el-dialog v-model="groupDialogVisible" :title="editingGroup ? '编辑分组' : '新增分组'" width="420px">
      <el-form :model="groupForm" label-width="70px">
        <el-form-item label="别名" required><el-input v-model="groupForm.alias" :disabled="!!editingGroup" placeholder="英文标识，如 datacenter" /></el-form-item>
        <el-form-item label="名称" required><el-input v-model="groupForm.name" placeholder="中文名称" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="groupForm.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="groupForm.order" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="groupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveGroup">确定</el-button>
      </template>
    </el-dialog>

    <!-- 模型对话框 -->
    <el-dialog v-model="modelDialogVisible" :title="editingModel ? '编辑模型' : '新增模型'" width="520px">
      <el-form :model="modelForm" label-width="70px">
        <el-form-item label="分组" required>
          <el-select v-model="modelForm.group_id" placeholder="选择所属分组" style="width: 100%">
            <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="别名" required><el-input v-model="modelForm.alias" :disabled="!!editingModel" placeholder="英文标识，如 host" /></el-form-item>
        <el-form-item label="名称" required><el-input v-model="modelForm.name" placeholder="中文名称" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="modelForm.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="图标"><el-input v-model="modelForm.icon" placeholder="图标名称" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="modelForm.order" :min="0" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="modelForm.is_usable" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modelDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveModel">确定</el-button>
      </template>
    </el-dialog>

    <!-- 移动分组对话框 -->
    <el-dialog v-model="moveGroupDialogVisible" title="移动分组" width="420px">
      <el-form label-width="70px">
        <el-form-item label="目标分组">
          <el-select v-model="moveToGroupId" placeholder="选择目标分组" style="width: 100%">
            <el-option v-for="g in groups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="moveGroupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleMoveGroup">确定</el-button>
      </template>
    </el-dialog>

    <!-- 字段分组对话框 -->
    <el-dialog v-model="fieldGroupDialogVisible" :title="editingFieldGroup ? '编辑字段分组' : '新增字段分组'" width="420px">
      <el-form :model="fieldGroupForm" label-width="60px">
        <el-form-item label="名称" required><el-input v-model="fieldGroupForm.name" placeholder="分组名称" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="fieldGroupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveFieldGroup">确定</el-button>
      </template>
    </el-dialog>

    <!-- 字段对话框 -->
    <el-dialog v-model="fieldDialogVisible" :title="editingField ? '编辑字段' : '新增字段'" width="520px">
      <el-form :model="fieldForm" label-width="70px">
        <el-form-item label="别名" required><el-input v-model="fieldForm.alias" :disabled="!!editingField" placeholder="英文标识，如 hostname" /></el-form-item>
        <el-form-item label="名称" required><el-input v-model="fieldForm.name" placeholder="中文名称" /></el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="fieldForm.type" :disabled="!!editingField" style="width: 100%">
            <el-option v-for="t in fieldTypes" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组" required>
          <el-select v-model="fieldForm.field_group_id" style="width: 100%">
            <el-option v-for="g in detailData.fieldGroups" :key="g.id" :label="g.name" :value="g.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="必填"><el-switch v-model="fieldForm.is_required" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="fieldForm.order" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="fieldDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveField">确定</el-button>
      </template>
    </el-dialog>

    <!-- 唯一约束对话框 -->
    <el-dialog v-model="uniqueDialogVisible" title="新增唯一约束" width="420px">
      <el-form :model="uniqueForm" label-width="60px">
        <el-form-item label="字段" required>
          <el-select v-model="uniqueForm.fields" multiple style="width: 100%" placeholder="选择构成唯一约束的字段">
            <el-option v-for="f in detailData.allFields" :key="f.alias" :label="`${f.name} (${f.alias})`" :value="f.alias" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="uniqueForm.description" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uniqueDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveUnique">确定</el-button>
      </template>
    </el-dialog>

    <!-- 关系对话框 -->
    <el-dialog v-model="relationDialogVisible" title="新增模型关系" width="420px">
      <el-form :model="relationForm" label-width="70px">
        <el-form-item label="源模型" required>
          <el-select v-model="relationForm.source_id" placeholder="选择源模型" style="width: 100%" filterable>
            <el-option v-for="m in allModelsForSelect" :key="m.id" :label="`${m.name} (${m.alias})`" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标模型" required>
          <el-select v-model="relationForm.target_id" placeholder="选择目标模型" style="width: 100%" filterable>
            <el-option v-for="m in targetModelOptions" :key="m.id" :label="`${m.name} (${m.alias})`" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="关系类型" required>
          <el-select v-model="relationForm.type_id" placeholder="选择关系类型" style="width: 100%">
            <el-option v-for="t in relationTypes" :key="t.id" :label="`${t.name} (${t.s2t})`" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="relationForm.description" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="relationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveRelation">确定</el-button>
      </template>
    </el-dialog>

    <!-- 关系类型对话框 -->
    <el-dialog v-model="rtDialogVisible" :title="rtEditing ? '编辑关系类型' : '新增关系类型'" width="420px">
      <el-form :model="rtForm" label-width="70px">
        <el-form-item label="名称" required>
          <el-input v-model="rtForm.name" placeholder="如：运行、部署、依赖" />
        </el-form-item>
        <el-form-item label="源→目标" required>
          <el-input v-model="rtForm.s2t" placeholder="如：运行在" />
        </el-form-item>
        <el-form-item label="目标→源" required>
          <el-input v-model="rtForm.t2s" placeholder="如：被运行于" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rtDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRtSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Search, Close, CircleCheckFilled } from '@element-plus/icons-vue'
import { listModelGroup, createModelGroup, patchModelGroup, deleteModelGroup } from '../../api/modelGroup'
import { listModel, createModel, updateModel, deleteModel, retrieveModel, patchModelGroupId } from '../../api/model'
import { listModelRelation, createModelRelation, deleteModelRelation } from '../../api/modelRelation'
import { listModelRelationType, createModelRelationType, updateModelRelationType, deleteModelRelationType } from '../../api/modelRelationType'
import { createModelField, updateModelField, deleteModelField } from '../../api/modelField'
import { createModelFieldGroup, updateModelFieldGroup, deleteModelFieldGroup } from '../../api/modelFieldGroup'
import { listModelFieldUnique, createModelFieldUnique, deleteModelFieldUnique } from '../../api/modelFieldUnique'

const fieldTypes = ['string', 'number', 'bool', 'date', 'datetime', 'json']

// ===== 分组 =====
const groups = ref([])
const selectedGroupId = ref(null)
const ungroupedCount = ref(0)

const loadGroups = async () => {
  const res = await listModelGroup({ limit: 1000 })
  groups.value = res.data.results || []
}

const currentGroupName = computed(() => {
  if (selectedGroupId.value === null) return '全部模型'
  if (selectedGroupId.value === 'ungrouped') return '未分组'
  return groups.value.find(g => g.id === selectedGroupId.value)?.name || '模型列表'
})

// 分组对话框
const groupDialogVisible = ref(false)
const editingGroup = ref(null)
const groupForm = ref({ alias: '', name: '', description: '', order: 0 })

const showGroupDialog = (group) => {
  editingGroup.value = group || null
  groupForm.value = group
    ? { alias: group.alias, name: group.name, description: group.description, order: group.order }
    : { alias: '', name: '', description: '', order: 0 }
  groupDialogVisible.value = true
}

const handleSaveGroup = async () => {
  if (editingGroup.value) {
    await patchModelGroup(editingGroup.value.id, { name: groupForm.value.name, description: groupForm.value.description })
  } else {
    await createModelGroup(groupForm.value)
  }
  groupDialogVisible.value = false
  ElMessage.success('操作成功')
  loadGroups()
}

const handleDeleteGroup = async (group) => {
  await ElMessageBox.confirm(`确定删除分组「${group.name}」？`, '确认删除', { type: 'warning' })
  await deleteModelGroup(group.id)
  ElMessage.success('删除成功')
  if (selectedGroupId.value === group.id) selectedGroupId.value = null
  loadGroups()
  loadModels()
}

// ===== 模型列表 =====
const models = ref([])
const modelsLoading = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const loadModels = async () => {
  modelsLoading.value = true
  try {
    const params = { search: search.value, limit: pageSize.value, offset: (page.value - 1) * pageSize.value }
    if (selectedGroupId.value === 'ungrouped') {
      params.ungrouped = true
    } else if (selectedGroupId.value !== null) {
      params.group_id = selectedGroupId.value
    }
    const res = await listModel(params)
    models.value = res.data.results || []
    total.value = res.data.count || 0
  } finally {
    modelsLoading.value = false
  }
}

const loadUngroupedCount = async () => {
  const res = await listModel({ ungrouped: true, limit: 1 })
  ungroupedCount.value = res.data.count || 0
}

watch(selectedGroupId, () => {
  page.value = 1
  loadModels()
})

// 模型对话框
const modelDialogVisible = ref(false)
const editingModel = ref(null)
const modelForm = ref({ group_id: null, alias: '', name: '', description: '', icon: '', order: 0, is_usable: true })
const allModelsForSelect = ref([])

const showModelDialog = (row) => {
  editingModel.value = row || null
  modelForm.value = row
    ? { group_id: row.group_id, alias: row.alias, name: row.name, description: row.description, icon: row.icon, order: row.order, is_usable: row.is_usable }
    : { group_id: selectedGroupId.value && selectedGroupId.value !== 'ungrouped' ? selectedGroupId.value : null, alias: '', name: '', description: '', icon: '', order: 0, is_usable: true }
  modelDialogVisible.value = true
}

const handleSaveModel = async () => {
  if (editingModel.value) {
    await updateModel(editingModel.value.id, { name: modelForm.value.name, description: modelForm.value.description, is_usable: modelForm.value.is_usable, icon: modelForm.value.icon, order: modelForm.value.order })
  } else {
    await createModel(modelForm.value)
  }
  modelDialogVisible.value = false
  ElMessage.success('操作成功')
  loadModels()
  loadUngroupedCount()
}

const handleDeleteModel = async (row) => {
  await ElMessageBox.confirm(`确定删除模型「${row.name}」？删除后其字段和字段分组也将被删除。`, '确认删除', { type: 'warning' })
  await deleteModel(row.id)
  ElMessage.success('删除成功')
  if (expandedModelId.value === row.id) expandedModelId.value = null
  loadModels()
  loadUngroupedCount()
}

// 移动分组
const moveGroupDialogVisible = ref(false)
const moveTargetModel = ref(null)
const moveToGroupId = ref(null)

const showMoveGroupDialog = (row) => {
  moveTargetModel.value = row
  moveToGroupId.value = row.group_id
  moveGroupDialogVisible.value = true
}

const handleMoveGroup = async () => {
  if (!moveTargetModel.value || !moveToGroupId.value) return
  await patchModelGroupId(moveTargetModel.value.id, { group_id: moveToGroupId.value })
  moveGroupDialogVisible.value = false
  ElMessage.success('移动成功')
  loadModels()
  loadUngroupedCount()
}

// ===== 模型详情展开 =====
const expandedModelId = ref(null)
const expandedModel = ref({})
const detailTab = ref('fields')
const selectedFieldGroup = ref(null)

const detailData = ref({
  fieldGroups: [],
  allFields: [],
  uniques: [],
  relations: [],
  relationsLoading: false,
})

const currentFields = computed(() => {
  if (!selectedFieldGroup.value) return detailData.value.allFields
  return detailData.value.allFields.filter(f => f.field_group_id === selectedFieldGroup.value)
})

const toggleDetail = async (row) => {
  if (expandedModelId.value === row.id) {
    expandedModelId.value = null
    return
  }
  expandedModelId.value = row.id
  expandedModel.value = row
  detailTab.value = 'fields'
  selectedFieldGroup.value = null
  await loadDetail(row.id)
}

const onRowClick = (row) => {
  toggleDetail(row)
}

const loadDetail = async (modelId) => {
  // 加载字段
  try {
    const res = await retrieveModel(modelId)
    const d = res.data
    detailData.value.fieldGroups = (d.groups || []).map(g => ({ id: g.group.id, name: g.group.name, order: g.group.order }))
    detailData.value.allFields = (d.groups || []).flatMap(g => g.fields || [])
  } catch {
    detailData.value.fieldGroups = []
    detailData.value.allFields = []
  }
  // 加载唯一约束
  try {
    const uRes = await listModelFieldUnique(modelId)
    detailData.value.uniques = uRes.data || []
  } catch {
    detailData.value.uniques = []
  }
  // 加载关系
  loadRelations(modelId)
}

const loadRelations = async (modelId) => {
  detailData.value.relationsLoading = true
  try {
    const res = await listModelRelation({ source_id: modelId, limit: 1000 })
    detailData.value.relations = res.data.results || []
  } catch {
    detailData.value.relations = []
  } finally {
    detailData.value.relationsLoading = false
  }
}

// ===== 字段分组 CRUD =====
const fieldGroupDialogVisible = ref(false)
const editingFieldGroup = ref(null)
const fieldGroupForm = ref({ name: '' })

const showFieldGroupDialog = (group) => {
  editingFieldGroup.value = group || null
  fieldGroupForm.value = group ? { name: group.name } : { name: '' }
  fieldGroupDialogVisible.value = true
}

const handleSaveFieldGroup = async () => {
  if (editingFieldGroup.value) {
    await updateModelFieldGroup(editingFieldGroup.value.id, { name: fieldGroupForm.value.name })
  } else {
    await createModelFieldGroup({ model_id: expandedModelId.value, name: fieldGroupForm.value.name })
  }
  fieldGroupDialogVisible.value = false
  ElMessage.success('操作成功')
  loadDetail(expandedModelId.value)
}

const handleDeleteFieldGroup = async (group) => {
  await ElMessageBox.confirm(`确定删除分组「${group.name}」？`, '确认删除', { type: 'warning' })
  await deleteModelFieldGroup(group.id)
  ElMessage.success('删除成功')
  if (selectedFieldGroup.value === group.id) selectedFieldGroup.value = null
  loadDetail(expandedModelId.value)
}

// ===== 字段 CRUD =====
const fieldDialogVisible = ref(false)
const editingField = ref(null)
const fieldForm = ref({ alias: '', name: '', type: 'string', is_required: false, order: 0, field_group_id: null })

const showFieldDialog = (field) => {
  editingField.value = field || null
  fieldForm.value = field
    ? { alias: field.alias, name: field.name, type: field.type, is_required: field.is_required, order: field.order, field_group_id: field.field_group_id }
    : { alias: '', name: '', type: 'string', is_required: false, order: 0, field_group_id: selectedFieldGroup.value || detailData.value.fieldGroups[0]?.id }
  fieldDialogVisible.value = true
}

const handleSaveField = async () => {
  if (editingField.value) {
    await updateModelField(editingField.value.id, { field_group_id: fieldForm.value.field_group_id, name: fieldForm.value.name, is_required: fieldForm.value.is_required, order: fieldForm.value.order })
  } else {
    await createModelField({ ...fieldForm.value, model_id: expandedModelId.value })
  }
  fieldDialogVisible.value = false
  ElMessage.success('操作成功')
  loadDetail(expandedModelId.value)
}

const handleDeleteField = async (field) => {
  await ElMessageBox.confirm(`确定删除字段「${field.name}」？`, '确认删除', { type: 'warning' })
  await deleteModelField(field.id)
  ElMessage.success('删除成功')
  loadDetail(expandedModelId.value)
}

// ===== 唯一约束 CRUD =====
const uniqueDialogVisible = ref(false)
const uniqueForm = ref({ fields: [], description: '' })

const showUniqueDialog = () => {
  uniqueForm.value = { fields: [], description: '' }
  uniqueDialogVisible.value = true
}

const handleSaveUnique = async () => {
  await createModelFieldUnique({ model_id: expandedModelId.value, fields: uniqueForm.value.fields.join(','), description: uniqueForm.value.description })
  uniqueDialogVisible.value = false
  ElMessage.success('创建成功')
  loadDetail(expandedModelId.value)
}

const handleDeleteUnique = async (row) => {
  await ElMessageBox.confirm('确定删除该唯一约束？', '确认删除', { type: 'warning' })
  await deleteModelFieldUnique(row.id)
  ElMessage.success('删除成功')
  loadDetail(expandedModelId.value)
}

// ===== 模型关系 CRUD =====
const relationDialogVisible = ref(false)
const relationForm = ref({ source_id: null, target_id: null, type_id: null, description: '' })
const relationTypes = ref([])

const targetModelOptions = computed(() => {
  const sourceId = expandedModelId.value
  const existingTargetIds = new Set((detailData.value.relations || []).map(r => r.target_id))
  return allModelsForSelect.value.filter(m => m.id !== sourceId && !existingTargetIds.has(m.id))
})

const loadRelationOptions = async () => {
  const [m, t] = await Promise.all([
    listModel({ limit: 1000 }),
    listModelRelationType({ limit: 1000 }),
  ])
  allModelsForSelect.value = m.data.results || []
  relationTypes.value = t.data.results || []
}

const showRelationDialog = () => {
  relationForm.value = { source_id: expandedModelId.value, target_id: null, type_id: null, description: '' }
  relationDialogVisible.value = true
}

const handleSaveRelation = async () => {
  await createModelRelation(relationForm.value)
  relationDialogVisible.value = false
  ElMessage.success('创建成功')
  loadRelations(expandedModelId.value)
}

const handleDeleteRelation = async (row) => {
  await ElMessageBox.confirm('确定删除该关系？', '确认删除', { type: 'warning' })
  await deleteModelRelation(row.id)
  ElMessage.success('删除成功')
  loadRelations(expandedModelId.value)
}

// ===== 关系类型 CRUD =====
const rtDialogVisible = ref(false)
const rtEditing = ref(null)
const rtForm = ref({ name: '', s2t: '', t2s: '' })

const showRtDialog = (row) => {
  rtEditing.value = row || null
  rtForm.value = row ? { name: row.name, s2t: row.s2t, t2s: row.t2s } : { name: '', s2t: '', t2s: '' }
  rtDialogVisible.value = true
}

const handleRtSave = async () => {
  if (rtEditing.value) {
    await updateModelRelationType(rtEditing.value.id, rtForm.value)
  } else {
    await createModelRelationType(rtForm.value)
  }
  rtDialogVisible.value = false
  ElMessage.success('操作成功')
  loadRelationOptions()
}

const handleRtDelete = async (row) => {
  await ElMessageBox.confirm(`确定删除「${row.name}」？`, '确认删除', { type: 'warning' })
  await deleteModelRelationType(row.id)
  ElMessage.success('删除成功')
  loadRelationOptions()
}

// ===== 初始化 =====
onMounted(async () => {
  await Promise.all([loadGroups(), loadModels(), loadUngroupedCount(), loadRelationOptions()])
})
</script>

<style scoped>
.model-manage-page {
  height: var(--page-height);
}

.main-row {
  height: 100%;
}

.left-panel {
  height: 100%;
}

.group-card {
  height: 100%;
}

.group-card :deep(.el-card__body) {
  padding: 8px 12px;
  overflow-y: auto;
  max-height: calc(var(--page-height) - 120px);
}

.right-panel {
  height: 100%;
  overflow-y: auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-title {
  font-weight: 600;
  font-size: 15px;
  color: var(--color-text-primary);
  white-space: nowrap;
}

/* 分组列表 */
.group-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.group-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 9px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 13px;
}

.group-item:hover {
  background: var(--color-muted);
}

.group-item.active {
  background: var(--color-primary-lighter);
}

.group-item.active .group-item-name {
  color: var(--color-primary);
  font-weight: 600;
}

.group-item-name {
  font-weight: 500;
  color: var(--color-text-primary);
}

.group-item-count {
  font-size: 11px;
  color: var(--color-text-muted);
  background: var(--color-muted);
  padding: 1px 6px;
  border-radius: 10px;
}

.group-item-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.15s ease;
}

.group-item:hover .group-item-actions {
  opacity: 1;
}

/* 分页 */
.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

/* 详情面板 */
.detail-card {
  margin-top: 12px;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-primary);
}

/* 字段分组列表 */
.field-group-list {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 8px;
}

.field-group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 8px 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  border-bottom: 1px solid var(--color-border);
  margin-bottom: 4px;
}

.model-row {
  cursor: pointer;
}
</style>
