<template>
  <div class="home">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <template v-if="loading">
        <div class="stat-card" v-for="i in 4" :key="i">
          <el-skeleton :loading="true" animated>
            <template #template>
              <div style="display: flex; align-items: center; gap: 16px; width: 100%">
                <el-skeleton-item variant="circle" style="width: 48px; height: 48px; flex-shrink: 0" />
                <div style="flex: 1">
                  <el-skeleton-item variant="h1" style="width: 60%; height: 28px" />
                  <el-skeleton-item variant="text" style="width: 40%; height: 14px; margin-top: 6px" />
                </div>
              </div>
            </template>
          </el-skeleton>
        </div>
      </template>
      <template v-else>
        <div class="stat-card" v-for="s in statCards" :key="s.label">
          <div class="stat-icon" :style="{ background: s.bg }">
            <el-icon :size="22" :color="s.color"><component :is="s.icon" /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ s.value }}</div>
            <div class="stat-label">{{ s.label }}</div>
          </div>
        </div>
      </template>
    </div>

    <!-- 实例数量分布图 -->
    <div class="section-title">实例数量分布</div>
    <div class="chart-card">
      <div v-if="loading" style="padding: 20px">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else-if="chartData.length" class="bar-chart">
        <div v-for="item in chartData" :key="item.label" class="bar-row">
          <div class="bar-label" :title="item.label">{{ item.label }}</div>
          <div class="bar-track">
            <div
              class="bar-fill"
              :style="{ width: item.pct + '%', background: item.color }"
            />
          </div>
          <div class="bar-value">{{ item.value }}</div>
        </div>
      </div>
      <el-empty v-else description="暂无实例数据" :image-size="48" />
    </div>

    <!-- 模型分组概览 -->
    <div class="section-title">模型分组</div>
    <div class="group-grid">
      <div v-for="g in groupCards" :key="g.name" class="group-card">
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
          <span class="group-sep">/</span>
          <span class="group-count">{{ g.instance_count }}</span>
          <span class="group-count-label">个实例</span>
        </div>
      </div>
      <el-empty v-if="!groupCards.length && !loading" description="暂无模型分组" :image-size="48" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Grid, List, Share } from '@element-plus/icons-vue'
import { getStats } from '../../api/stats'
import { listModelGroup } from '../../api/modelGroup'

const loading = ref(true)
const statCards = ref([])
const modelInstanceCounts = ref([])
const groupCards = ref([])

const BAR_COLORS = ['#3B82F6', '#10B981', '#F59E0B', '#8B5CF6', '#EF4444', '#06B6D4', '#EC4899', '#84CC16', '#F97316', '#6366F1']

const chartData = computed(() => {
  const all = modelInstanceCounts.value
  if (!all.length) return []
  const TOP_N = 10
  const top = all.slice(0, TOP_N)
  const rest = all.slice(TOP_N)
  const maxVal = top[0]?.instance_count || 1

  const items = top.map((m, i) => ({
    label: `${m.model_name} (${m.model_alias})`,
    value: m.instance_count,
    pct: Math.max(2, (m.instance_count / maxVal) * 100),
    color: BAR_COLORS[i % BAR_COLORS.length],
  }))

  if (rest.length) {
    const restTotal = rest.reduce((s, m) => s + m.instance_count, 0)
    items.push({
      label: `其他 (${rest.length} 个模型)`,
      value: restTotal,
      pct: Math.max(2, (restTotal / maxVal) * 100),
      color: '#94A3B8',
    })
  }

  return items
})

onMounted(async () => {
  try {
    const [statsRes, groupsRes] = await Promise.all([
      getStats(),
      listModelGroup({ limit: 1000 }),
    ])
    const stats = statsRes.data || {}

    statCards.value = [
      { label: '模型分组', value: stats.model_group_count || 0, icon: Grid, color: '#2563EB', bg: '#DBEAFE' },
      { label: '模型数量', value: stats.model_count || 0, icon: Grid, color: '#059669', bg: '#D1FAE5' },
      { label: '实例总数', value: stats.instance_count || 0, icon: List, color: '#D97706', bg: '#FEF3C7' },
      { label: '关系总数', value: stats.relation_count || 0, icon: Share, color: '#7C3AED', bg: '#EDE9FE' },
    ]

    modelInstanceCounts.value = stats.model_instance_counts || []

    // 模型分组概览
    const groups = groupsRes.data?.results || []
    const instanceCountByGroup = {}
    const modelCountByGroup = {}
    for (const m of (stats.model_instance_counts || [])) {
      const g = m.group_name || '未分组'
      instanceCountByGroup[g] = (instanceCountByGroup[g] || 0) + m.instance_count
      modelCountByGroup[g] = (modelCountByGroup[g] || 0) + 1
    }

    const groupColors = [
      { color: '#2563EB', bg: '#DBEAFE' },
      { color: '#059669', bg: '#D1FAE5' },
      { color: '#D97706', bg: '#FEF3C7' },
      { color: '#7C3AED', bg: '#EDE9FE' },
      { color: '#DC2626', bg: '#FEE2E2' },
      { color: '#0891B2', bg: '#CFFAFE' },
    ]

    groupCards.value = groups.map((group, i) => ({
      ...group,
      model_count: modelCountByGroup[group.name] || 0,
      instance_count: instanceCountByGroup[group.name] || 0,
      ...groupColors[i % groupColors.length],
    }))
  } catch {
    // error shown by interceptor
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.home {
  max-width: 1200px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
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

/* 柱状图 */
.chart-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 20px;
  margin-bottom: 28px;
}

.bar-chart {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.bar-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bar-label {
  width: 180px;
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bar-track {
  flex: 1;
  height: 24px;
  background: var(--color-muted);
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.6s ease;
  min-width: 4px;
}

.bar-value {
  width: 60px;
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-primary);
  text-align: right;
}

/* 分组卡片 */
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

.group-sep {
  color: var(--color-border);
  margin: 0 6px;
  font-size: 14px;
}
</style>
