<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>关系类型</span>
        <el-button type="primary" @click="showDialog()">
          <el-icon><Plus /></el-icon>新增
        </el-button>
      </div>
    </template>
    <el-table :data="list" stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="名称" width="140" />
      <el-table-column prop="s2t" label="源 → 目标" />
      <el-table-column prop="t2s" label="目标 → 源" />
      <el-table-column prop="created_at" label="创建时间" width="170" />
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="showDialog(row)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="dialogVisible" :title="editing ? '编辑关系类型' : '新增关系类型'" width="420px">
    <el-form :model="form" label-width="70px">
      <el-form-item label="名称" required>
        <el-input v-model="form.name" placeholder="如：运行、部署、依赖" />
      </el-form-item>
      <el-form-item label="源→目标" required>
        <el-input v-model="form.s2t" placeholder="如：运行在" />
      </el-form-item>
      <el-form-item label="目标→源" required>
        <el-input v-model="form.t2s" placeholder="如：被运行于" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { listModelRelationType, createModelRelationType, updateModelRelationType, deleteModelRelationType } from '../../api/modelRelationType'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editing = ref(null)
const form = ref({ name: '', s2t: '', t2s: '' })

const loadData = async () => {
  loading.value = true
  try {
    const res = await listModelRelationType({ limit: 1000 })
    list.value = res.data.results
  } finally {
    loading.value = false
  }
}

const showDialog = (row) => {
  editing.value = row || null
  form.value = row ? { name: row.name, s2t: row.s2t, t2s: row.t2s } : { name: '', s2t: '', t2s: '' }
  dialogVisible.value = true
}

const handleSave = async () => {
  if (editing.value) {
    await updateModelRelationType(editing.value.id, form.value)
  } else {
    await createModelRelationType(form.value)
  }
  dialogVisible.value = false
  ElMessage.success('操作成功')
  loadData()
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm(`确定删除「${row.name}」？`, '确认删除', { type: 'warning' })
  await deleteModelRelationType(row.id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
