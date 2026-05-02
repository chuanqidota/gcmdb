<template>
  <div class="home">
    <!-- 统计卡片 -->
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

    <!-- 快捷入口 -->
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

    <!-- 模型分组概览 -->
    <div class="section-title">模型分组</div>
    <div class="group-grid">
      <div v-for="g in groupCards" :key="g.id" class="group-card">
        <div class="group-card-header">
          <div class="group-icon" :style="{ background: g.bg }">
            <el-icon :size="20" :color="g.color"><Grid /></el-icon>
          </div>
          <div class="group-name">{{ g.name }}</div>
        </div>
        <div class="group-desc">{{ g.description || '暂无描述' }}</div>
        <div class="group-meta">
          <span class="group-count">{{ g.model_count }}</span>
          <span class="group-count-label">个模型</span>
        </div>
      </div>
      <el-empty v-if="!groupCards.length && !loading" description="暂无模型分组" :image-size="48" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Grid, List, Search, Share } from '@element-plus/icons-vue'
import { listModelGroup } from '../../api/modelGroup'
import { listModel } from '../../api/model'
import { listInstance } from '../../api/instance'

const loading = ref(true)
const statCards = ref([])

const shortcuts = [
  { label: '综合检索', desc: '搜索 CMDB 中的资源', path: '/search', icon: Search, color: '#2563EB', bg: '#DBEAFE' },
  { label: '实例管理', desc: '管理各模型的实例数据', path: '/instance', icon: List, color: '#059669', bg: '#D1FAE5' },
  { label: '模型管理', desc: '配置模型、字段和关系', path: '/model-manage', icon: Grid, color: '#D97706', bg: '#FEF3C7' },
  { label: '模型拓扑', desc: '可视化模型关系图', path: '/model-topology', icon: Share, color: '#7C3AED', bg: '#EDE9FE' },
]

const groupCards = ref([])
const groupColors = [
  { color: '#2563EB', bg: '#DBEAFE' },
  { color: '#059669', bg: '#D1FAE5' },
  { color: '#D97706', bg: '#FEF3C7' },
  { color: '#7C3AED', bg: '#EDE9FE' },
  { color: '#DC2626', bg: '#FEE2E2' },
  { color: '#0891B2', bg: '#CFFAFE' },
]

onMounted(async () => {
  const [g, m] = await Promise.all([
    listModelGroup({ limit: 1000 }),
    listModel({ limit: 1000 }),
  ])
  const groups = g.data.results || []
  const models = m.data.results || []

  // 统计每个分组的模型数量
  const countMap = {}
  models.forEach(model => {
    countMap[model.group_id] = (countMap[model.group_id] || 0) + 1
  })

  groupCards.value = groups.map((group, i) => ({
    ...group,
    model_count: countMap[group.id] || 0,
    ...groupColors[i % groupColors.length],
  }))

  // 统计卡片
  let totalInstances = 0
  if (models.length > 0) {
    const countResults = await Promise.allSettled(
      models.map(model => listInstance(model.id, { limit: 1 }))
    )
    countResults.forEach(r => {
      if (r.status === 'fulfilled') totalInstances += r.value.data.count || 0
    })
  }

  statCards.value = [
    { label: '模型分组', value: g.data.count || 0, icon: Grid, color: '#2563EB', bg: '#DBEAFE' },
    { label: '模型数量', value: m.data.count || 0, icon: Grid, color: '#059669', bg: '#D1FAE5' },
    { label: '实例总数', value: totalInstances, icon: List, color: '#D97706', bg: '#FEF3C7' },
  ]

  loading.value = false
})
</script>

<style scoped>
.home {
  max-width: 1200px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 28px;
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
  margin-bottom: 28px;
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

.group-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}

.group-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 20px;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.group-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.group-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.group-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.group-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.group-desc {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-bottom: 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.group-meta {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.group-count {
  font-size: 22px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.group-count-label {
  font-size: 12px;
  color: var(--color-text-muted);
}
</style>
