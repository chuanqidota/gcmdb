<template>
  <div class="user-manage">
    <div class="section-header">
      <span class="section-title">用户管理</span>
      <el-button type="primary" @click="openCreateDialog">创建用户</el-button>
    </div>

    <el-table :data="list" stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column label="Token" min-width="280">
        <template #default="{ row }">
          <div class="token-cell">
            <code class="token-text">{{ row.token }}</code>
            <el-button link type="primary" size="small" @click="copyToken(row.token)">复制</el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
            {{ row.is_active ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="170" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="openResetPwd(row)">重置密码</el-button>
          <el-button
            link
            :type="row.is_active ? 'danger' : 'success'"
            size="small"
            @click="toggleActive(row)"
          >
            {{ row.is_active ? '禁用' : '启用' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 创建用户对话框 -->
    <el-dialog v-model="createDialog.visible" title="创建用户" width="400px" :close-on-click-modal="false">
      <el-form ref="createFormRef" :model="createDialog.form" :rules="createRules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="createDialog.form.username" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="createDialog.form.password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="createDialog.loading" @click="submitCreate">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog v-model="resetDialog.visible" title="重置密码" width="400px" :close-on-click-modal="false">
      <p class="reset-hint">为用户 <strong>{{ resetDialog.username }}</strong> 设置新密码</p>
      <el-form ref="resetFormRef" :model="resetDialog.form" :rules="resetRules" label-width="80px">
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="resetDialog.form.new_password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="resetDialog.loading" @click="submitResetPwd">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listUsers, createUser, patchUser, resetPassword } from '../../api/auth'

const list = ref([])
const loading = ref(false)
const createFormRef = ref(null)
const resetFormRef = ref(null)

const createDialog = ref({
  visible: false,
  loading: false,
  form: { username: '', password: '' },
})

const resetDialog = ref({
  visible: false,
  loading: false,
  userId: null,
  username: '',
  form: { new_password: '' },
})

const createRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
}

const resetRules = {
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await listUsers()
    list.value = res.data || []
  } finally {
    loading.value = false
  }
}

const copyToken = (token) => {
  navigator.clipboard.writeText(token)
  ElMessage.success('已复制到剪贴板')
}

const openCreateDialog = () => {
  createDialog.value.form = { username: '', password: '' }
  createDialog.value.visible = true
}

const submitCreate = async () => {
  await createFormRef.value.validate()
  createDialog.value.loading = true
  try {
    const res = await createUser(createDialog.value.form)
    ElMessage.success(`用户创建成功，Token: ${res.data?.token || ''}`)
    createDialog.value.visible = false
    loadData()
  } catch {
    // error shown by interceptor
  } finally {
    createDialog.value.loading = false
  }
}

const openResetPwd = (row) => {
  resetDialog.value.userId = row.id
  resetDialog.value.username = row.username
  resetDialog.value.form = { new_password: '' }
  resetDialog.value.visible = true
}

const submitResetPwd = async () => {
  await resetFormRef.value.validate()
  resetDialog.value.loading = true
  try {
    await resetPassword({
      user_id: resetDialog.value.userId,
      new_password: resetDialog.value.form.new_password,
    })
    ElMessage.success('密码重置成功')
    resetDialog.value.visible = false
  } catch {
    // error shown by interceptor
  } finally {
    resetDialog.value.loading = false
  }
}

const toggleActive = async (row) => {
  const action = row.is_active ? '禁用' : '启用'
  await ElMessageBox.confirm(`确定${action}用户 ${row.username}？`, '确认', { type: 'warning' })
  try {
    await patchUser(row.id, { is_active: !row.is_active })
    ElMessage.success(`${action}成功`)
    loadData()
  } catch {
    // error shown by interceptor
  }
}

onMounted(loadData)
</script>

<style scoped>
.user-manage {
  max-width: 1200px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.token-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.token-text {
  font-size: 12px;
  color: var(--color-text-secondary);
  word-break: break-all;
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.reset-hint {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: var(--color-text-secondary);
}
</style>
