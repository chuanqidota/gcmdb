<template>
  <div class="model-group-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="group-card">
          <template #header>
            <div class="card-header">
              <span>模型分组</span>
              <el-button type="primary" size="small" @click="showGroupDialog()">
                <el-icon><Plus /></el-icon>新增
              </el-button>
            </div>
          </template>
          <div class="group-list">
            <div
              v-for="g in groups" :key="g.id"
              class="group-item"
              :class="{ active: selectedGroup === g.id }"
              @click="onSelectGroup(g.id)"
            >
              <div class="group-item-info">
                <div class="group-item-name">{{ g.name }}</div>
                <div class="group-item-alias">{{ g.alias }}</div>
              </div>
              <div class="group-item-actions">
                <el-button link size="small" @click.stop="showGroupDialog(g)">
                  <el-icon><Edit /></el-icon>
                </el-button>
                <el-button link size="small" type="danger" @click.stop="handleDeleteGroup(g)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
            <el-empty v-if="!groups.length" description="暂无分组" :image-size="60" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>{{ selectedGroupName || '请选择分组查看模型' }}</span>
          </template>
          <el-table :data="groupModels" stripe v-if="selectedGroup">
            <el-table-column prop="id" label="ID" width="60" />
            <el-table-column prop="alias" label="别名" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="description" label="描述" show-overflow-tooltip />
            <el-table-column prop="is_usable" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.is_usable ? 'success' : 'danger'" size="small" effect="light">
                  {{ row.is_usable ? '启用' : '停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button link type="primary" @click="router.push(`/model/${row.id}`)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-else description="请在左侧选择一个分组" />
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="dialogVisible" :title="editingGroup ? '编辑分组' : '新增分组'" width="460px">
      <el-form :model="form" label-width="70px">
        <el-form-item label="别名" required>
          <el-input v-model="form.alias" :disabled="!!editingGroup" placeholder="英文标识，如 datacenter" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="中文名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="可选描述" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveGroup">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { listModelGroup, createModelGroup, patchModelGroup, deleteModelGroup, retrieveModelGroup } from '../../api/modelGroup'

const router = useRouter()
const groups = ref([])
const selectedGroup = ref(null)
const selectedGroupName = ref('')
const groupModels = ref([])
const dialogVisible = ref(false)
const editingGroup = ref(null)
const form = ref({ alias: '', name: '', description: '', order: 0 })

const loadGroups = async () => {
  const res = await listModelGroup({ limit: 1000 })
  groups.value = res.data.results
}

const onSelectGroup = async (id) => {
  selectedGroup.value = id
  const g = groups.value.find((g) => g.id === id)
  selectedGroupName.value = g ? g.name : ''
  const res = await retrieveModelGroup(id)
  groupModels.value = res.data.models || []
}

const showGroupDialog = (group) => {
  editingGroup.value = group || null
  form.value = group
    ? { alias: group.alias, name: group.name, description: group.description, order: group.order }
    : { alias: '', name: '', description: '', order: 0 }
  dialogVisible.value = true
}

const handleSaveGroup = async () => {
  if (editingGroup.value) {
    await patchModelGroup(editingGroup.value.id, { name: form.value.name, description: form.value.description })
  } else {
    await createModelGroup(form.value)
  }
  dialogVisible.value = false
  ElMessage.success('操作成功')
  loadGroups()
}

const handleDeleteGroup = async (group) => {
  await ElMessageBox.confirm(`确定删除分组「${group.name}」？`, '确认删除', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' })
  await deleteModelGroup(group.id)
  ElMessage.success('删除成功')
  if (selectedGroup.value === group.id) {
    selectedGroup.value = null
    groupModels.value = []
  }
  loadGroups()
}

onMounted(loadGroups)
</script>

<style scoped>
.model-group-page {
  max-width: 1200px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.group-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.group-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
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
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.group-item-alias {
  font-size: 11px;
  color: var(--color-text-muted);
  margin-top: 1px;
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
