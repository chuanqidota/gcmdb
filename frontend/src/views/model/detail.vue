<template>
  <div class="model-detail" v-loading="loading">
    <el-page-header @back="router.push('/model')" style="margin-bottom: 24px">
      <template #content>
        <div class="page-header-content">
          <span class="model-name">{{ modelInfo.name }}</span>
          <el-tag size="small" type="info">{{ modelInfo.alias }}</el-tag>
        </div>
      </template>
    </el-page-header>

    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="field-group-card">
          <template #header>
            <div class="card-header">
              <span>字段分组</span>
              <el-button size="small" type="primary" @click="showGroupDialog()">
                <el-icon><Plus /></el-icon>
              </el-button>
            </div>
          </template>
          <div class="group-list">
            <div
              v-for="g in fieldGroups" :key="g.id"
              class="group-item"
              :class="{ active: selectedGroup === g.id }"
              @click="onSelectGroup(g.id)"
            >
              <span class="group-item-name">{{ g.name }}</span>
              <div class="group-item-actions">
                <el-button link size="small" @click.stop="showGroupDialog(g)"><el-icon><Edit /></el-icon></el-button>
                <el-button link size="small" type="danger" @click.stop="handleDeleteGroup(g)"><el-icon><Delete /></el-icon></el-button>
              </div>
            </div>
            <div class="group-item" :class="{ active: !selectedGroup }" @click="selectedGroup = null">
              <span class="group-item-name">全部字段</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="18">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>字段列表</span>
              <el-button size="small" type="primary" @click="showFieldDialog()">
                <el-icon><Plus /></el-icon>新增字段
              </el-button>
            </div>
          </template>
          <el-table :data="currentFields" stripe>
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="alias" label="别名" width="140" />
            <el-table-column prop="name" label="名称" width="140" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }"><el-tag size="small" effect="plain">{{ row.type }}</el-tag></template>
            </el-table-column>
            <el-table-column prop="is_required" label="必填" width="70">
              <template #default="{ row }">
                <el-tag :type="row.is_required ? 'danger' : 'info'" size="small" effect="light">{{ row.is_required ? '是' : '否' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="order" label="排序" width="70" />
            <el-table-column label="操作" width="130" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="showFieldDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="handleDeleteField(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-card style="margin-top: 16px">
          <template #header>
            <div class="card-header">
              <span>唯一约束</span>
              <el-button size="small" type="primary" @click="showUniqueDialog">
                <el-icon><Plus /></el-icon>新增
              </el-button>
            </div>
          </template>
          <el-table :data="uniqueConstraints" stripe v-if="uniqueConstraints.length">
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="fields" label="字段组合">
              <template #default="{ row }">
                <el-tag v-for="f in row.fields.split(',')" :key="f" size="small" style="margin-right: 4px">{{ f }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" />
            <el-table-column label="操作" width="80">
              <template #default="{ row }">
                <el-button link type="danger" @click="handleDeleteUnique(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-else description="暂无唯一约束" :image-size="48" />
        </el-card>
      </el-col>
    </el-row>
  </div>

  <el-dialog v-model="groupDialogVisible" :title="editingGroup ? '编辑分组' : '新增分组'" width="400px">
    <el-form :model="groupForm" label-width="60px">
      <el-form-item label="名称" required><el-input v-model="groupForm.name" placeholder="分组名称" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="groupDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSaveGroup">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="fieldDialogVisible" :title="editingField ? '编辑字段' : '新增字段'" width="500px">
    <el-form :model="fieldForm" label-width="70px">
      <el-form-item label="别名" required>
        <el-input v-model="fieldForm.alias" :disabled="!!editingField" placeholder="英文标识，如 hostname" />
      </el-form-item>
      <el-form-item label="名称" required>
        <el-input v-model="fieldForm.name" placeholder="中文名称" />
      </el-form-item>
      <el-form-item label="类型" required>
        <el-select v-model="fieldForm.type" :disabled="!!editingField" style="width: 100%">
          <el-option v-for="t in fieldTypes" :key="t" :label="t" :value="t" />
        </el-select>
      </el-form-item>
      <el-form-item label="分组" required>
        <el-select v-model="fieldForm.field_group_id" style="width: 100%">
          <el-option v-for="g in fieldGroups" :key="g.id" :label="g.name" :value="g.id" />
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

  <el-dialog v-model="uniqueDialogVisible" title="新增唯一约束" width="440px">
    <el-form :model="uniqueForm" label-width="60px">
      <el-form-item label="字段" required>
        <el-select v-model="uniqueForm.fields" multiple style="width: 100%" placeholder="选择构成唯一约束的字段">
          <el-option v-for="f in allFields" :key="f.alias" :label="`${f.name} (${f.alias})`" :value="f.alias" />
        </el-select>
      </el-form-item>
      <el-form-item label="描述"><el-input v-model="uniqueForm.description" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="uniqueDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSaveUnique">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { retrieveModel } from '../../api/model'
import { createModelField, updateModelField, deleteModelField } from '../../api/modelField'
import { createModelFieldGroup, updateModelFieldGroup, deleteModelFieldGroup } from '../../api/modelFieldGroup'
import { listModelFieldUnique, createModelFieldUnique, deleteModelFieldUnique } from '../../api/modelFieldUnique'

const route = useRoute()
const router = useRouter()
const modelId = Number(route.params.id)
const loading = ref(false)
const modelInfo = ref({})
const fieldGroups = ref([])
const allFields = ref([])
const selectedGroup = ref(null)
const uniqueConstraints = ref([])
const fieldTypes = ['string', 'number', 'bool', 'date', 'datetime', 'json']

const currentFields = computed(() => {
  if (!selectedGroup.value) return allFields.value
  return allFields.value.filter((f) => f.field_group_id === selectedGroup.value)
})

const groupDialogVisible = ref(false)
const editingGroup = ref(null)
const groupForm = ref({ name: '' })
const fieldDialogVisible = ref(false)
const editingField = ref(null)
const fieldForm = ref({ alias: '', name: '', type: 'string', is_required: false, order: 0, field_group_id: null })
const uniqueDialogVisible = ref(false)
const uniqueForm = ref({ fields: [], description: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await retrieveModel(modelId)
    const d = res.data
    modelInfo.value = { id: d.id, alias: d.alias, name: d.name, description: d.description }
    fieldGroups.value = (d.groups || []).map((g) => ({ id: g.group.id, name: g.group.name, order: g.group.order }))
    allFields.value = (d.groups || []).flatMap((g) => g.fields || [])
    const uRes = await listModelFieldUnique(modelId)
    uniqueConstraints.value = uRes.data || []
  } finally {
    loading.value = false
  }
}

const onSelectGroup = (id) => { selectedGroup.value = id }

const showGroupDialog = (group) => {
  editingGroup.value = group || null
  groupForm.value = group ? { name: group.name } : { name: '' }
  groupDialogVisible.value = true
}

const handleSaveGroup = async () => {
  if (editingGroup.value) {
    await updateModelFieldGroup(editingGroup.value.id, { name: groupForm.value.name })
  } else {
    await createModelFieldGroup({ model_id: modelId, name: groupForm.value.name })
  }
  groupDialogVisible.value = false
  ElMessage.success('操作成功')
  loadData()
}

const handleDeleteGroup = async (group) => {
  await ElMessageBox.confirm(`确定删除分组「${group.name}」？`, '确认删除', { type: 'warning' })
  await deleteModelFieldGroup(group.id)
  ElMessage.success('删除成功')
  if (selectedGroup.value === group.id) selectedGroup.value = null
  loadData()
}

const showFieldDialog = (field) => {
  editingField.value = field || null
  fieldForm.value = field
    ? { alias: field.alias, name: field.name, type: field.type, is_required: field.is_required, order: field.order, field_group_id: field.field_group_id }
    : { alias: '', name: '', type: 'string', is_required: false, order: 0, field_group_id: selectedGroup.value || fieldGroups.value[0]?.id }
  fieldDialogVisible.value = true
}

const handleSaveField = async () => {
  if (editingField.value) {
    await updateModelField(editingField.value.id, { field_group_id: fieldForm.value.field_group_id, name: fieldForm.value.name, is_required: fieldForm.value.is_required, order: fieldForm.value.order })
  } else {
    await createModelField({ ...fieldForm.value, model_id: modelId })
  }
  fieldDialogVisible.value = false
  ElMessage.success('操作成功')
  loadData()
}

const handleDeleteField = async (field) => {
  await ElMessageBox.confirm(`确定删除字段「${field.name}」？`, '确认删除', { type: 'warning' })
  await deleteModelField(field.id)
  ElMessage.success('删除成功')
  loadData()
}

const showUniqueDialog = () => {
  uniqueForm.value = { fields: [], description: '' }
  uniqueDialogVisible.value = true
}

const handleSaveUnique = async () => {
  await createModelFieldUnique({ model_id: modelId, fields: uniqueForm.value.fields.join(','), description: uniqueForm.value.description })
  uniqueDialogVisible.value = false
  ElMessage.success('创建成功')
  loadData()
}

const handleDeleteUnique = async (row) => {
  await ElMessageBox.confirm('确定删除该唯一约束？', '确认删除', { type: 'warning' })
  await deleteModelFieldUnique(row.id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.model-detail {
  max-width: 1200px;
}

.page-header-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.model-name {
  font-size: 16px;
  font-weight: 600;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

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

.group-item-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.15s ease;
}

.group-item:hover .group-item-actions {
  opacity: 1;
}
</style>
