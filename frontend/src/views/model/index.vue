<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <el-input v-model="search" placeholder="搜索模型名称或别名..." style="width: 320px" clearable @clear="loadData" @keyup.enter="loadData">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-button type="primary" @click="showDialog()">
          <el-icon><Plus /></el-icon>新增模型
        </el-button>
      </div>
    </template>
    <el-table :data="list" stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="alias" label="别名" width="140" />
      <el-table-column prop="name" label="名称" width="140" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column prop="is_usable" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_usable ? 'success' : 'danger'" size="small" effect="light">
            {{ row.is_usable ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="router.push(`/model/${row.id}`)">字段管理</el-button>
          <el-button link type="primary" @click="showDialog(row)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination-wrap">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next" @change="loadData" />
    </div>
  </el-card>

  <el-dialog v-model="dialogVisible" :title="editing ? '编辑模型' : '新增模型'" width="520px">
    <el-form :model="form" label-width="70px">
      <el-form-item label="分组" required>
        <el-select v-model="form.group_id" placeholder="选择所属分组" style="width: 100%">
          <el-option v-for="g in groupOptions" :key="g.id" :label="g.name" :value="g.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="别名" required>
        <el-input v-model="form.alias" :disabled="!!editing" placeholder="英文标识，如 host" />
      </el-form-item>
      <el-form-item label="名称" required>
        <el-input v-model="form.name" placeholder="中文名称" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="图标">
        <el-input v-model="form.icon" placeholder="图标名称" />
      </el-form-item>
      <el-form-item label="排序">
        <el-input-number v-model="form.order" :min="0" />
      </el-form-item>
      <el-form-item label="启用">
        <el-switch v-model="form.is_usable" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleSave">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { listModel, createModel, updateModel, deleteModel } from '../../api/model'
import { listModelGroup } from '../../api/modelGroup'

const router = useRouter()
const list = ref([])
const loading = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dialogVisible = ref(false)
const editing = ref(null)
const groupOptions = ref([])
const form = ref({ group_id: null, alias: '', name: '', description: '', icon: '', order: 0, is_usable: true })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listModel({ search: search.value, limit: pageSize.value, offset: (page.value - 1) * pageSize.value })
    list.value = res.data.results
    total.value = res.data.count
  } finally {
    loading.value = false
  }
}

const loadGroups = async () => {
  const res = await listModelGroup({ limit: 1000 })
  groupOptions.value = res.data.results
}

const showDialog = (row) => {
  editing.value = row || null
  form.value = row
    ? { group_id: row.group_id, alias: row.alias, name: row.name, description: row.description, icon: row.icon, order: row.order, is_usable: row.is_usable }
    : { group_id: null, alias: '', name: '', description: '', icon: '', order: 0, is_usable: true }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (editing.value) {
    await updateModel(editing.value.id, { name: form.value.name, description: form.value.description, is_usable: form.value.is_usable, icon: form.value.icon, order: form.value.order })
  } else {
    await createModel(form.value)
  }
  dialogVisible.value = false
  ElMessage.success('操作成功')
  loadData()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm(`确定删除模型「${row.name}」？`, '确认删除', { type: 'warning' })
  await deleteModel(row.id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(() => { loadData(); loadGroups() })
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
