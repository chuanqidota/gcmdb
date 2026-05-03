<template>
  <el-container style="height: 100vh">
    <el-aside :width="isCollapsed ? '64px' : '220px'" class="sidebar" :class="{ collapsed: isCollapsed }">
      <div class="logo">
        <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
          <rect width="28" height="28" rx="8" fill="var(--color-primary)"/>
          <path d="M8 10h12M8 14h8M8 18h10" stroke="#fff" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <span v-show="!isCollapsed" class="logo-text">GCMDB</span>
      </div>
      <el-menu
        :default-active="route.path"
        :collapse="isCollapsed"
        background-color="transparent"
        text-color="#94A3B8"
        active-text-color="#FFFFFF"
        router
        class="nav-menu"
      >
        <el-tooltip :disabled="!isCollapsed" content="首页" placement="right" :show-after="300">
          <el-menu-item index="/home">
            <el-icon><HomeFilled /></el-icon>
            <template #title>首页</template>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip :disabled="!isCollapsed" content="综合检索" placement="right" :show-after="300">
          <el-menu-item index="/search">
            <el-icon><Search /></el-icon>
            <template #title>综合检索</template>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip :disabled="!isCollapsed" content="实例管理" placement="right" :show-after="300">
          <el-menu-item index="/instance">
            <el-icon><List /></el-icon>
            <template #title>实例管理</template>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip :disabled="!isCollapsed" content="模型管理" placement="right" :show-after="300">
          <el-menu-item index="/model-manage">
            <el-icon><Grid /></el-icon>
            <template #title>模型管理</template>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip :disabled="!isCollapsed" content="模型拓扑" placement="right" :show-after="300">
          <el-menu-item index="/model-topology">
            <el-icon><Share /></el-icon>
            <template #title>模型拓扑</template>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip :disabled="!isCollapsed" content="SQL 查询" placement="right" :show-after="300">
          <el-menu-item index="/search-direct-sql">
            <el-icon><Document /></el-icon>
            <template #title>SQL 查询</template>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip :disabled="!isCollapsed" content="审计日志" placement="right" :show-after="300">
          <el-menu-item index="/audit">
            <el-icon><Notebook /></el-icon>
            <template #title>审计日志</template>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip v-if="isAdmin" :disabled="!isCollapsed" content="用户管理" placement="right" :show-after="300">
          <el-menu-item index="/user-manage">
            <el-icon><User /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
        </el-tooltip>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="app-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="isCollapsed = !isCollapsed">
            <Fold v-if="!isCollapsed" /><Expand v-else />
          </el-icon>
          <span class="page-title">{{ route.meta.title }}</span>
        </div>
        <el-dropdown @command="handleCommand">
          <span class="user-dropdown">
            <el-icon><User /></el-icon>
            <span class="username">{{ username }}</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="token">查看 Token</el-dropdown-item>
              <el-dropdown-item command="password">修改密码</el-dropdown-item>
              <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>

  <!-- 修改密码对话框 -->
  <el-dialog v-model="pwdDialog.visible" title="修改密码" width="400px" :close-on-click-modal="false">
    <el-form ref="pwdFormRef" :model="pwdDialog.form" :rules="pwdRules" label-width="80px">
      <el-form-item label="当前密码" prop="old_password">
        <el-input v-model="pwdDialog.form.old_password" type="password" show-password />
      </el-form-item>
      <el-form-item label="新密码" prop="new_password">
        <el-input v-model="pwdDialog.form.new_password" type="password" show-password />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirm_password">
        <el-input v-model="pwdDialog.form.confirm_password" type="password" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="pwdDialog.visible = false">取消</el-button>
      <el-button type="primary" :loading="pwdDialog.loading" @click="submitChangePassword">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { HomeFilled, Grid, List, Search, Fold, Expand, Share, Document, Notebook, User } from '@element-plus/icons-vue'
import { getMe, logout, changePassword } from '../api/auth'
import { useUserStore } from '../stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isCollapsed = ref(false)
const username = ref('')
const isAdmin = ref(false)
const pwdFormRef = ref(null)

const pwdDialog = ref({
  visible: false,
  loading: false,
  form: { old_password: '', new_password: '', confirm_password: '' },
})

const pwdRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdDialog.value.form.new_password) {
          callback(new Error('两次密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

const handleResize = () => {
  isCollapsed.value = window.innerWidth < 768
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

onMounted(async () => {
  try {
    const res = await getMe()
    username.value = res.data?.username || ''
    isAdmin.value = res.data?.is_admin || false
    userStore.username = username.value
    userStore.isAdmin = isAdmin.value
    userStore.token = res.data?.token || ''
    userStore.userId = res.data?.id || null
  } catch {
    // will be redirected by 401 interceptor
  }
})

const handleCommand = async (cmd) => {
  if (cmd === 'logout') {
    await ElMessageBox.confirm('确定退出登录？', '确认', { type: 'warning' })
    try { await logout() } catch {}
    localStorage.removeItem('gcmdb_logged_in')
    router.push('/login')
    ElMessage.success('已退出登录')
  } else if (cmd === 'token') {
    try {
      const res = await getMe()
      const token = res.data?.token || ''
      await ElMessageBox.alert(`Token: <code style="word-break:break-all">${token}</code>`, 'API Token', { dangerouslyUseHTMLString: true, confirmButtonText: '复制并关闭' })
      navigator.clipboard.writeText(token)
      ElMessage.success('已复制到剪贴板')
    } catch {}
  } else if (cmd === 'password') {
    pwdDialog.value.form = { old_password: '', new_password: '', confirm_password: '' }
    pwdDialog.value.visible = true
  }
}

const submitChangePassword = async () => {
  await pwdFormRef.value.validate()
  pwdDialog.value.loading = true
  try {
    await changePassword({
      old_password: pwdDialog.value.form.old_password,
      new_password: pwdDialog.value.form.new_password,
    })
    ElMessage.success('密码修改成功')
    pwdDialog.value.visible = false
  } catch {
    // error shown by interceptor
  } finally {
    pwdDialog.value.loading = false
  }
}
</script>

<style scoped>
.sidebar {
  background: #0F172A;
  transition: width 0.2s ease;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.logo-text {
  color: #F8FAFC;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.02em;
  white-space: nowrap;
}

.nav-menu {
  flex: 1;
  padding: 8px 0;
  border: none !important;
}

.nav-menu:not(.el-menu--collapse) {
  width: 220px;
}

.nav-menu .el-menu-item,
.nav-menu :deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  margin: 2px 8px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.15s ease;
}

.nav-menu .el-menu-item:hover,
.nav-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.06) !important;
  color: #E2E8F0 !important;
}

.nav-menu .el-menu-item.is-active {
  background: var(--color-primary) !important;
  color: #FFFFFF !important;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}

.nav-menu :deep(.el-sub-menu .el-menu-item) {
  padding-left: 52px !important;
  font-size: 13px;
}

.app-header {
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 56px;
  box-shadow: var(--shadow-sm);
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  font-size: 18px;
  cursor: pointer;
  color: var(--color-text-secondary);
  padding: 6px;
  border-radius: 6px;
  transition: all 0.15s ease;
}

.collapse-btn:hover {
  background: var(--color-muted);
  color: var(--color-primary);
}

.page-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.app-main {
  background: var(--color-background);
  padding: 20px;
  overflow-y: auto;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 13px;
  color: var(--color-text-secondary);
  padding: 6px 12px;
  border-radius: var(--radius-md);
  transition: all 0.15s ease;
}

.user-dropdown:hover {
  background: var(--color-muted);
  color: var(--color-primary);
}

.username {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
