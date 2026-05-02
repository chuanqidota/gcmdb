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

    <!-- 模型实例分布 -->
    <div class="dist-section">
      <div class="section-title">模型实例分布</div>
      <div class="model-dist">
        <div v-for="m in modelStats" :key="m.id" class="model-dist-item">
          <div class="model-dist-label">{{ m.name }}</div>
          <div class="model-dist-bar-wrap">
            <div class="model-dist-bar" :style="{ width: m.percent + '%' }"></div>
          </div>
          <div class="model-dist-count">{{ m.count }}</div>
        </div>
        <el-empty v-if="!modelStats.length && !loading" description="暂无模型数据" :image-size="48" />
      </div>
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

const modelStats = ref([])

onMounted(async () => {
  const [g, m] = await Promise.all([
    listModelGroup({ limit: 1 }),
    listModel({ limit: 1000 }),
  ])
  const models = m.data.results || []

  // 并行获取每个模型的实例数
  const countResults = await Promise.allSettled(
    models.map(model => listInstance(model.id, { limit: 1 }))
  )

  let totalInstances = 0
  const modelData = models.map((model, i) => {
    const count = countResults[i].status === 'fulfilled' ? (countResults[i].value.data.count || 0) : 0
    totalInstances += count
    return { id: model.id, name: model.name, count }
  })

  const max = Math.max(...modelData.map(m => m.count), 1)
  modelStats.value = modelData
    .sort((a, b) => b.count - a.count)
    .map(m => ({ ...m, percent: Math.round((m.count / max) * 100) }))

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

.dist-section {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 20px;
}

.model-dist {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px 24px;
  max-height: 400px;
  overflow-y: auto;
}

.model-dist-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.model-dist-label {
  width: 80px;
  font-size: 13px;
  color: var(--color-text-primary);
  text-align: right;
  flex-shrink: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.model-dist-bar-wrap {
  flex: 1;
  height: 20px;
  background: var(--color-muted);
  border-radius: 4px;
  overflow: hidden;
}

.model-dist-bar {
  height: 100%;
  background: var(--color-primary);
  border-radius: 4px;
  transition: width 0.4s ease;
  min-width: 2px;
}

.model-dist-count {
  width: 40px;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
  text-align: right;
  flex-shrink: 0;
}
</style>
