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
        <el-menu-item index="/home">
          <el-icon><HomeFilled /></el-icon>
          <template #title>首页</template>
        </el-menu-item>
        <el-sub-menu index="model-mgmt">
          <template #title>
            <el-icon><Grid /></el-icon>
            <span>模型管理</span>
          </template>
          <el-menu-item index="/model-group">模型分组</el-menu-item>
          <el-menu-item index="/model">模型列表</el-menu-item>
          <el-menu-item index="/model-relation-type">关系类型</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/model-relation">
          <el-icon><Connection /></el-icon>
          <template #title>模型关系</template>
        </el-menu-item>
        <el-menu-item index="/instance">
          <el-icon><List /></el-icon>
          <template #title>实例管理</template>
        </el-menu-item>
        <el-menu-item index="/search-direct-sql">
          <el-icon><Search /></el-icon>
          <template #title>SQL 查询</template>
        </el-menu-item>
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
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { HomeFilled, Grid, Connection, List, Search, Fold, Expand } from '@element-plus/icons-vue'

const route = useRoute()
const isCollapsed = ref(false)
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
</style>
