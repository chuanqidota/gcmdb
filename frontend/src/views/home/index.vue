<template>
  <div class="home">
    <div class="stats-grid">
      <div class="stat-card" v-for="s in statCards" :key="s.label">
        <div class="stat-icon" :style="{ background: s.bg }">
          <el-icon :size="22" :color="s.color"><component :is="s.icon" /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ s.value }}</div>
          <div class="stat-label">{{ s.label }}</div>
        </div>
      </div>
    </div>

    <div class="section-title">快捷入口</div>
    <div class="shortcuts-grid">
      <router-link v-for="item in shortcuts" :key="item.path" :to="item.path" class="shortcut-card">
        <div class="shortcut-icon" :style="{ background: item.bg }">
          <el-icon :size="24" :color="item.color"><component :is="item.icon" /></el-icon>
        </div>
        <div class="shortcut-label">{{ item.label }}</div>
        <div class="shortcut-desc">{{ item.desc }}</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Grid, List, Connection, Search } from '@element-plus/icons-vue'
import { listModelGroup } from '../../api/modelGroup'
import { listModel } from '../../api/model'
import { listModelRelationType } from '../../api/modelRelationType'

const stats = ref({ groups: 0, models: 0, relationTypes: 0 })

const statCards = computed(() => [
  { label: '模型分组', value: stats.value.groups, icon: Grid, color: '#2563EB', bg: '#DBEAFE' },
  { label: '模型数量', value: stats.value.models, icon: List, color: '#059669', bg: '#D1FAE5' },
  { label: '关系类型', value: stats.value.relationTypes, icon: Connection, color: '#D97706', bg: '#FEF3C7' },
])

const shortcuts = [
  { label: '模型分组', desc: '管理模型分组', path: '/model-group', icon: Grid, color: '#2563EB', bg: '#DBEAFE' },
  { label: '模型列表', desc: '查看所有模型', path: '/model', icon: List, color: '#059669', bg: '#D1FAE5' },
  { label: '模型关系', desc: '管理模型间关系', path: '/model-relation', icon: Connection, color: '#D97706', bg: '#FEF3C7' },
  { label: 'SQL 查询', desc: '执行自定义查询', path: '/search-direct-sql', icon: Search, color: '#DC2626', bg: '#FEE2E2' },
]

onMounted(async () => {
  const [g, m, r] = await Promise.all([
    listModelGroup({ limit: 1 }),
    listModel({ limit: 1 }),
    listModelRelationType({ limit: 1 }),
  ])
  stats.value = { groups: g.data.count, models: m.data.count, relationTypes: r.data.count }
})
</script>

<style scoped>
.home {
  max-width: 1100px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-top: 2px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 16px;
}

.shortcuts-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.shortcut-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px 20px;
  text-align: center;
  text-decoration: none;
  transition: all 0.2s ease;
  cursor: pointer;
}

.shortcut-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
  border-color: var(--color-primary-lighter);
}

.shortcut-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
}

.shortcut-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 4px;
}

.shortcut-desc {
  font-size: 12px;
  color: var(--color-text-muted);
}
</style>
