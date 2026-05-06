<template>
  <div class="topology-page">
    <!-- 图区域 -->
    <div ref="graphContainer" class="graph-container" />

    <!-- 右侧详情面板 -->
    <transition name="slide">
      <div v-if="selectedModel" class="detail-panel">
        <div class="panel-header">
          <div class="panel-title-row">
            <span class="panel-model-name">{{ selectedModel.name }}</span>
            <el-tag size="small" type="info">{{ selectedModel.alias }}</el-tag>
          </div>
          <el-icon class="panel-close" @click="selectedModel = null"><Close /></el-icon>
        </div>
        <div v-if="selectedModel.description" class="panel-desc">{{ selectedModel.description }}</div>

        <!-- 字段 -->
        <div class="panel-section">
          <div class="panel-section-title">字段</div>
          <template v-if="detailLoading">
            <div class="panel-loading">加载中...</div>
          </template>
          <template v-else-if="detailFields.length">
            <el-table :data="detailFields" size="small" stripe highlight-current-row :show-header="true" max-height="320">
              <el-table-column prop="alias" label="别名" width="120" />
              <el-table-column prop="name" label="名称" width="120" />
              <el-table-column prop="type" label="类型" width="80">
                <template #default="{ row }">
                  <el-tag size="small" :type="typeTag[row.type] || 'info'">{{ row.type }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="is_required" label="必填" width="60" align="center">
                <template #default="{ row }">
                  <el-icon v-if="row.is_required" color="var(--color-destructive)"><CircleCheckFilled /></el-icon>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>
          </template>
          <template v-else>
            <div class="panel-loading">暂无字段</div>
          </template>
        </div>

        <!-- 关系 -->
        <div class="panel-section">
          <div class="panel-section-title">模型关系</div>
          <template v-if="detailRelations.length">
            <div v-for="rel in detailRelations" :key="rel.id" class="relation-row">
              <span class="relation-label">{{ rel.typeLabel }}</span>
              <span class="relation-target">{{ rel.targetName }}</span>
            </div>
          </template>
          <template v-else>
            <div class="panel-loading">暂无关系</div>
          </template>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { Graph } from '@antv/g6'
import { Close, CircleCheckFilled } from '@element-plus/icons-vue'
import { getAllModels, getModelGroups, getModelRelationTypes, getModelDetail, getAllModelRelations } from '../../api/search'
import { GROUP_COLORS } from '../../utils/colors'

const typeTag = { string: '', number: 'success', bool: 'warning', date: 'info', datetime: 'info', json: 'danger' }

const graphContainer = ref(null)
const selectedModel = ref(null)
const detailFields = ref([])
const detailRelations = ref([])
const detailLoading = ref(false)

let graph = null
let allModels = []
let allRelationTypes = []
let allGroups = []

onMounted(async () => {
  const [modelsRes, groupsRes, typesRes, relationsRes] = await Promise.all([
    getAllModels(),
    getModelGroups(),
    getModelRelationTypes(),
    getAllModelRelations(),
  ])

  allModels = (modelsRes.data || []).filter(m => m.is_usable !== false)
  allGroups = groupsRes.data || []
  allRelationTypes = typesRes.data || []
  const relations = relationsRes.data || []

  // 分组颜色映射
  const groupColorMap = {}
  allGroups.forEach((g, i) => {
    groupColorMap[g.id] = GROUP_COLORS[i % GROUP_COLORS.length]
  })

  // 关系类型映射
  const relTypeMap = {}
  allRelationTypes.forEach(t => { relTypeMap[t.id] = t })

  // 构建节点
  const modelIdSet = new Set(allModels.map(m => m.id))
  const nodes = allModels.map(m => ({
    id: String(m.id),
    data: {
      label: m.name,
      alias: m.alias,
      model: m,
    },
    style: {
      fill: groupColorMap[m.group_id] || '#94A3B8',
      stroke: '#fff',
      lineWidth: 2,
      radius: 24,
    },
  }))

  // 构建边 (只包含两端都在模型集合中的关系)
  const edges = relations
    .filter(r => modelIdSet.has(r.source_id) && modelIdSet.has(r.target_id))
    .map(r => ({
      id: `edge-${r.id}`,
      source: String(r.source_id),
      target: String(r.target_id),
      data: {
        relType: relTypeMap[r.type_id],
      },
      style: {
        endArrow: true,
        stroke: '#C8D0DE',
        lineWidth: 1.5,
      },
    }))

  await nextTick()
  if (!graphContainer.value) return

  const rect = graphContainer.value.getBoundingClientRect()

  graph = new Graph({
    container: graphContainer.value,
    width: rect.width,
    height: rect.height,
    data: { nodes, edges },
    node: {
      style: {
        labelText: (d) => d.data?.label || '',
        labelFontSize: 12,
        labelFontWeight: 600,
        labelFill: '#1E293B',
        labelPlacement: 'bottom',
        labelOffsetY: 8,
        cursor: 'pointer',
      },
    },
    edge: {
      style: {
        labelText: (d) => d.data?.relType?.s2t || '',
        labelFontSize: 10,
        labelFill: '#64748B',
        endArrow: true,
        endArrowSize: 6,
        stroke: '#C8D0DE',
        lineWidth: 1.5,
      },
    },
    layout: {
      type: 'force',
      preventOverlap: true,
      nodeSize: 48,
      linkDistance: 180,
      strength: {
        charge: -400,
      },
    },
    behaviors: ['drag-canvas', 'zoom-canvas', 'drag-element'],
    autoFit: 'view',
  })

  await graph.render()

  // 节点点击
  let prevSelected = null
  graph.on('node:click', (e) => {
    const nodeId = e.target.id
    const model = allModels.find(m => String(m.id) === nodeId)
    if (!model) return

    // 高亮选中节点
    if (prevSelected && prevSelected !== nodeId) {
      const prevModel = allModels.find(m => String(m.id) === prevSelected)
      graph.updateNodeData([{ id: prevSelected, style: { stroke: '#fff', lineWidth: 2 } }])
    }
    graph.updateNodeData([{ id: nodeId, style: { stroke: '#3B82F6', lineWidth: 3 } }])
    prevSelected = nodeId

    selectedModel.value = model
    loadModelDetail(model.alias)
  })

  // 画布点击关闭面板
  graph.on('canvas:click', () => {
    if (prevSelected) {
      graph.updateNodeData([{ id: prevSelected, style: { stroke: '#fff', lineWidth: 2 } }])
      prevSelected = null
    }
    selectedModel.value = null
  })
})

onBeforeUnmount(() => {
  if (graph) {
    graph.destroy()
    graph = null
  }
})

async function loadModelDetail(alias) {
  detailLoading.value = true
  detailFields.value = []
  detailRelations.value = []
  try {
    const res = await getModelDetail(alias)
    const data = res.data || {}
    const modelId = data.model?.id
    detailFields.value = data.model_fields || []

    // 解析关系
    const relations = data.model_relations || []
    detailRelations.value = relations.map(rel => {
      const isSource = rel.source_id === modelId
      const targetId = isSource ? rel.target_id : rel.source_id
      const targetModel = allModels.find(m => m.id === targetId)
      const relType = allRelationTypes.find(t => t.id === rel.type_id)
      const direction = isSource ? 's2t' : 't2s'
      return {
        id: rel.id,
        targetName: targetModel ? `${targetModel.name}(${targetModel.alias})` : `#${targetId}`,
        typeLabel: relType ? relType[direction] : '',
      }
    })
  } catch {
    detailFields.value = []
    detailRelations.value = []
  } finally {
    detailLoading.value = false
  }
}
</script>

<style scoped>
.topology-page {
  display: flex;
  height: var(--page-height);
  position: relative;
  overflow: hidden;
}

.graph-container {
  flex: 1;
  background: var(--color-background);
  border-radius: var(--radius-md);
}

/* 详情面板 */
.detail-panel {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 380px;
  z-index: 10;
  background: var(--color-surface);
  border-left: 1px solid var(--color-border);
  padding: 20px;
  overflow-y: auto;
  box-shadow: -4px 0 16px rgba(0, 0, 0, 0.06);
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.25s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.panel-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 12px;
}

.panel-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.panel-model-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.panel-close {
  cursor: pointer;
  font-size: 18px;
  color: var(--color-text-muted);
  padding: 4px;
  border-radius: 4px;
  flex-shrink: 0;
}

.panel-close:hover {
  background: var(--color-muted);
  color: var(--color-text-primary);
}

.panel-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: 16px;
  line-height: 1.5;
}

.panel-section {
  margin-bottom: 20px;
}

.panel-section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 10px;
  padding-bottom: 6px;
  border-bottom: 1px solid var(--color-border);
}

.panel-loading {
  text-align: center;
  padding: 16px;
  font-size: 13px;
  color: var(--color-text-muted);
}

/* 关系行 */
.relation-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  font-size: 13px;
}

.relation-label {
  color: var(--color-primary);
  font-weight: 500;
}

.relation-target {
  font-weight: 600;
  color: var(--color-text-primary);
}
</style>
