<template>
  <div class="search-page">
    <!-- Tab 栏 -->
    <div class="tab-bar">
      <span
        v-for="tab in tabs"
        :key="tab.name"
        class="tab-item"
        :class="{ active: activeTab === tab.name }"
        @click="switchTab(tab.name)"
      >{{ tab.label }}</span>
    </div>

    <!-- 搜索区 -->
    <div class="search-hero">
      <!-- 全文检索 -->
      <div v-if="activeTab === 'fulltext'" class="search-box">
        <el-input
          v-model="ft.keyword"
          placeholder="输入关键词搜索全部模型数据..."
          clearable
          @keyup.enter="doFulltextSearch"
          class="search-input"
        />
        <el-select v-model="ft.modelAlias" placeholder="全部模型" clearable filterable class="search-model-select">
          <el-option label="全部模型" value="" />
          <el-option v-for="m in allModels" :key="m.id" :label="m.name" :value="m.alias" />
        </el-select>
        <el-button type="primary" class="search-btn" @click="doFulltextSearch" :loading="ft.loading">搜索</el-button>
      </div>

      <!-- 实例搜索：模型选择 + 操作按钮（与全文检索同款居中布局） -->
      <div v-if="activeTab === 'instance'" class="instance-search-box">
        <el-select
          v-model="inst.modelId"
          placeholder="选择模型开始搜索..."
          filterable
          class="instance-model-select"
          @change="onModelChange"
        >
          <el-option
            v-for="m in allModels"
            :key="m.id"
            :label="`${m.name} (${m.alias})`"
            :value="m.id"
          />
        </el-select>
        <el-button type="primary" class="search-btn" @click="doInstanceSearch" :loading="inst.loading" :disabled="!inst.modelId">搜索</el-button>
        <el-button class="search-btn" @click="resetInstance">重置</el-button>
      </div>

      <!-- 模型浏览 -->
      <div v-if="activeTab === 'model'" class="search-box">
        <el-input
          v-model="modelSearch"
          placeholder="搜索模型名称、别名..."
          clearable
          class="search-input"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
      </div>
    </div>

    <!-- 结果区 -->
    <div class="results-area">

      <!-- 全文检索结果 -->
      <template v-if="activeTab === 'fulltext' && (ft.results.length > 0 || ft.loading)">
        <div class="result-count" v-if="ft.count > 0">找到 <strong>{{ ft.count }}</strong> 条结果</div>
        <el-table :data="ft.results" stripe v-loading="ft.loading" size="small">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="model_name" label="模型" width="120">
            <template #default="{ row }"><el-tag size="small">{{ row.model_name }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="model_alias" label="别名" width="100" />
          <el-table-column label="数据" min-width="300">
            <template #default="{ row }">
              <span v-for="(v, k) in (row.data || {})" :key="k" class="data-tag">
                <strong>{{ k }}:</strong> {{ v }}
              </span>
            </template>
          </el-table-column>
        </el-table>
        <div v-if="ft.count > 0" class="pagination-wrap">
          <el-pagination
            v-model:current-page="ft.page"
            v-model:page-size="ft.pageSize"
            :total="ft.count"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @current-change="doFulltextSearch"
            @size-change="doFulltextSearch"
          />
        </div>
      </template>

      <!-- 实例搜索：条件构建器 + 结果 -->
      <template v-if="activeTab === 'instance'">
        <!-- 条件构建器 -->
        <div v-if="inst.modelId" class="condition-builder">
          <div v-for="(c, i) in inst.conditions" :key="i" class="condition-row">
            <el-select v-model="c.field" placeholder="字段" class="cond-field">
              <el-option v-for="f in inst.fields" :key="f.alias" :label="f.name" :value="f.alias" />
            </el-select>
            <el-select v-model="c.op" class="cond-op">
              <el-option label="=" value="eq" />
              <el-option label="!=" value="ne" />
              <el-option label=">" value="gt" />
              <el-option label=">=" value="ge" />
              <el-option label="<" value="lt" />
              <el-option label="<=" value="le" />
              <el-option label="包含" value="contains" />
              <el-option label="开头是" value="startswith" />
              <el-option label="结尾是" value="endswith" />
              <el-option label="在列表" value="in" />
            </el-select>
            <el-input v-model="c.val" placeholder="值" class="cond-val" />
            <el-button :icon="Delete" circle @click="inst.conditions.splice(i, 1)" />
          </div>
          <div class="condition-actions">
            <el-button @click="inst.conditions.push({ field: '', op: 'eq', val: '' })">+ 添加条件</el-button>
          </div>
        </div>
        <div class="result-count" v-if="inst.count > 0">找到 <strong>{{ inst.count }}</strong> 条实例</div>
        <el-table :data="inst.results" stripe v-loading="inst.loading" size="small">
          <el-table-column
            v-for="col in inst.columns"
            :key="col"
            :prop="col"
            :label="inst.columnLabelMap[col] || col"
            :min-width="calcColWidth(inst.columnLabelMap[col] || col)"
            show-overflow-tooltip
          />
        </el-table>
        <div v-if="inst.count > 0" class="pagination-wrap">
          <el-pagination
            v-model:current-page="inst.page"
            v-model:page-size="inst.pageSize"
            :total="inst.count"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @current-change="doInstanceSearch"
            @size-change="doInstanceSearch"
          />
        </div>
      </template>

      <!-- 模型浏览结果 -->
      <template v-if="activeTab === 'model'">
        <div v-for="group in groupedModels" :key="group.id" class="model-group-section">
          <div class="model-group-header">{{ group.name }}</div>
          <div class="model-list">
            <template v-for="m in group.models" :key="m.id">
              <div
                class="model-card"
                :class="{ expanded: expandedModelId === m.id }"
                @click="toggleFields(m)"
              >
                <div class="model-card-header">
                  <div class="model-card-info">
                    <span class="model-card-name">{{ m.name }}</span>
                    <el-tag size="small" type="info">{{ m.alias }}</el-tag>
                  </div>
                  <el-icon class="expand-icon" :class="{ rotated: expandedModelId === m.id }"><ArrowDown /></el-icon>
                </div>
                <div v-if="m.description" class="model-card-desc">{{ m.description }}</div>
              </div>
              <!-- 字段表内联展开（按分组） -->
              <div v-if="expandedModelId === m.id" class="fields-inline">
                <template v-if="modelFieldsCache[m.id]?.loading">
                  <div class="fields-loading">加载中...</div>
                </template>
                <template v-else-if="modelFieldsCache[m.id]?.fieldGroups?.length">
                  <div v-for="group in modelFieldsCache[m.id].fieldGroups" :key="group.name" class="field-group">
                    <div class="field-group-title">{{ group.name || '未分组' }}</div>
                    <el-table :data="group.fields" stripe size="small" :show-header="group._showHeader">
                      <el-table-column prop="alias" label="别名" width="150">
                        <template #default="{ row }">
                          <span>{{ row.alias }}</span>
                          <el-tag v-if="modelFieldsCache[m.id]?.uniqueAliases?.includes(row.alias)" size="small" type="warning" class="unique-tag">唯一</el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="name" label="名称" width="150" />
                      <el-table-column prop="type" label="类型" width="100">
                        <template #default="{ row }">
                          <el-tag size="small" :type="typeTag[row.type] || 'info'">{{ row.type }}</el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="is_required" label="必填" width="80" align="center">
                        <template #default="{ row }">
                          <el-icon v-if="row.is_required" color="var(--color-destructive)"><CircleCheckFilled /></el-icon>
                          <span v-else>-</span>
                        </template>
                      </el-table-column>
                      <el-table-column prop="order" label="排序" width="80" align="center" />
                    </el-table>
                  </div>
                </template>
                <template v-else>
                  <div class="fields-loading">暂无字段</div>
                </template>
                <!-- 模型关系 -->
                <div v-if="modelFieldsCache[m.id]?.relations?.length" class="model-relations">
                  <div class="relations-title">模型关系</div>
                  <div v-for="rel in modelFieldsCache[m.id].relations" :key="rel.id" class="relation-item">
                    <span class="relation-type">{{ resolveRelation(rel).typeLabel }}</span>
                    <span class="relation-target">{{ resolveRelation(rel).targetName }}</span>
                    <span v-if="resolveRelation(rel).description" class="relation-desc">{{ resolveRelation(rel).description }}</span>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Search, Delete, CircleCheckFilled, ArrowDown } from '@element-plus/icons-vue'
import { getAllModels, getModelGroups, getModelRelationTypes, getModelDetail, fulltextSearch, searchInstance } from '../../api/search'

const typeTag = { string: '', number: 'success', bool: 'warning', date: 'info', datetime: 'info', json: 'danger' }

const tabs = [
  { name: 'fulltext', label: '全文检索' },
  { name: 'instance', label: '实例搜索' },
  { name: 'model', label: '模型浏览' },
]

const activeTab = ref('fulltext')
const allModels = ref([])

function switchTab(name) {
  activeTab.value = name
}

// ===== 全文检索 =====
const ft = ref({
  keyword: '',
  modelAlias: '',
  results: [],
  count: 0,
  page: 1,
  pageSize: 10,
  loading: false,
})

async function doFulltextSearch() {
  const hasKeyword = ft.value.keyword.trim().length > 0
  const hasModel = !!ft.value.modelAlias
  // 无关键词且无模型选择时不搜索
  if (!hasKeyword && !hasModel) return

  ft.value.loading = true
  try {
    // 有关键词 → 全文检索
    if (hasKeyword) {
      const params = {
        search: ft.value.keyword.trim(),
        limit: ft.value.pageSize,
        offset: (ft.value.page - 1) * ft.value.pageSize,
      }
      if (hasModel) params.model_alias = ft.value.modelAlias
      const res = await fulltextSearch(params)
      ft.value.results = res.data?.result || []
      ft.value.count = res.data?.count || 0
    } else {
      // 无关键词但选了模型 → 用实例搜索列出全部
      const body = {
        model: ft.value.modelAlias,
        __condition: {
          limit: ft.value.pageSize,
          offset: (ft.value.page - 1) * ft.value.pageSize,
          order: ['-id'],
        },
      }
      const res = await searchInstance(body)
      const results = res.data?.results || []
      ft.value.results = results.map(r => ({
        id: r.id,
        model_name: ft.value.modelAlias,
        model_alias: ft.value.modelAlias,
        data: Object.fromEntries(Object.entries(r).filter(([k]) => !['id', 'model_id', 'model_alias'].includes(k))),
      }))
      ft.value.count = res.data?.count || 0
    }
  } catch {
    ft.value.results = []
    ft.value.count = 0
  } finally {
    ft.value.loading = false
  }
}

// ===== 实例搜索 =====
const inst = ref({
  modelId: null,
  fields: [],
  conditions: [],
  results: [],
  columns: [],
  columnLabelMap: {},
  count: 0,
  page: 1,
  pageSize: 10,
  loading: false,
})

async function onModelChange(id) {
  inst.value.fields = []
  inst.value.conditions = []
  inst.value.results = []
  inst.value.columns = []
  inst.value.count = 0
  inst.value.page = 1
  if (!id) return
  try {
    const m = allModels.value.find(x => x.id === id)
    const res = await getModelDetail(m?.alias || id)
    inst.value.fields = res.data?.model_fields || []
  } catch {
    inst.value.fields = []
  }
}

async function doInstanceSearch() {
  const m = allModels.value.find(x => x.id === inst.value.modelId)
  if (!m) return
  inst.value.loading = true
  try {
    const where = []
    for (const c of inst.value.conditions) {
      if (!c.field || !c.op) continue
      if (c.op === 'in') {
        const vals = c.val.split(',').map(v => v.trim()).filter(Boolean)
        if (vals.length) where.push({ in: { [c.field]: vals } })
      } else if (c.val !== '') {
        where.push({ [c.op]: { [c.field]: c.val } })
      }
    }
    const body = {
      model: m.alias,
      fields: inst.value.fields.map(f => f.alias),
      __condition: {
        limit: inst.value.pageSize,
        offset: (inst.value.page - 1) * inst.value.pageSize,
        order: ['-id'],
      },
    }
    if (where.length) body.__condition.where = where
    const res = await searchInstance(body)
    inst.value.results = res.data?.results || []
    inst.value.count = res.data?.count || 0
    if (inst.value.results.length) {
      inst.value.columns = Object.keys(inst.value.results[0]).filter(k => !['model_id', 'model_alias'].includes(k))
      // 构建列标签映射: alias → "name(alias)"
      const labelMap = { id: 'ID' }
      for (const f of inst.value.fields) {
        labelMap[f.alias] = f.name ? `${f.name}(${f.alias})` : f.alias
      }
      inst.value.columnLabelMap = labelMap
    }
  } catch {
    inst.value.results = []
    inst.value.count = 0
  } finally {
    inst.value.loading = false
  }
}

function calcColWidth(label) {
  // 中文字符按 14px，英文/数字按 8px，加上 padding
  let w = 0
  for (const ch of label) {
    w += /[一-鿿]/.test(ch) ? 14 : 8
  }
  return Math.max(w + 24, 80)
}

function resetInstance() {
  inst.value.conditions = []
  inst.value.results = []
  inst.value.columns = []
  inst.value.columnLabelMap = {}
  inst.value.count = 0
  inst.value.page = 1
}

// ===== 模型浏览 =====
const modelSearch = ref('')
const expandedModelId = ref(null)
const allModelGroups = ref([])
const allRelationTypes = ref([])
// 缓存: { [modelId]: { fieldGroups, uniqueAliases, relations, loading } }
const modelFieldsCache = ref({})

const groupedModels = computed(() => {
  const kw = modelSearch.value.toLowerCase()
  const filtered = kw
    ? allModels.value.filter(m =>
        m.name.toLowerCase().includes(kw) ||
        m.alias.toLowerCase().includes(kw) ||
        (m.description || '').toLowerCase().includes(kw))
    : allModels.value

  const groupMap = new Map()
  for (const g of allModelGroups.value) {
    groupMap.set(g.id, { ...g, models: [] })
  }
  const ungrouped = { id: 0, name: '未分组', alias: '', models: [] }

  for (const m of filtered) {
    if (m.group_id && groupMap.has(m.group_id)) {
      groupMap.get(m.group_id).models.push(m)
    } else {
      ungrouped.models.push(m)
    }
  }

  const result = []
  for (const g of groupMap.values()) {
    if (g.models.length) result.push(g)
  }
  if (ungrouped.models.length) result.push(ungrouped)
  return result
})

async function toggleFields(m) {
  if (expandedModelId.value === m.id) {
    expandedModelId.value = null
    return
  }
  expandedModelId.value = m.id
  if (modelFieldsCache.value[m.id]) return
  modelFieldsCache.value[m.id] = { fieldGroups: [], uniqueAliases: [], loading: true }
  try {
    const res = await getModelDetail(m.alias)
    const data = res.data || {}
    const fields = data.model_fields || []
    const groups = data.model_field_groups || []
    const uniques = data.model_field_uniques || []

    // 唯一字段别名
    const uniqueSet = new Set()
    for (const u of uniques) {
      if (u.fields) {
        u.fields.split(',').map(f => f.trim()).filter(Boolean).forEach(f => uniqueSet.add(f))
      }
    }

    // 按分组归类字段
    const groupMap = new Map()
    for (const g of groups) {
      groupMap.set(g.id, { name: g.name, fields: [], _showHeader: true })
    }
    const ungrouped = { name: '', fields: [], _showHeader: true }

    for (const f of fields) {
      if (f.field_group_id && groupMap.has(f.field_group_id)) {
        groupMap.get(f.field_group_id).fields.push(f)
      } else {
        ungrouped.fields.push(f)
      }
    }

    const fieldGroups = []
    for (const g of groupMap.values()) {
      if (g.fields.length) fieldGroups.push(g)
    }
    if (ungrouped.fields.length) fieldGroups.push(ungrouped)

    // 只有第一组显示表头，其余组隐藏避免重复
    fieldGroups.forEach((g, i) => { g._showHeader = i === 0 })

    modelFieldsCache.value[m.id] = {
      fieldGroups,
      uniqueAliases: [...uniqueSet],
      relations: data.model_relations || [],
      loading: false,
    }
  } catch {
    modelFieldsCache.value[m.id] = { fieldGroups: [], uniqueAliases: [], relations: [], loading: false }
  }
}

function resolveRelation(rel) {
  const isSource = rel.source_id === expandedModelId.value
  const targetId = isSource ? rel.target_id : rel.source_id
  const targetModel = allModels.value.find(m => m.id === targetId)
  const relType = allRelationTypes.value.find(t => t.id === rel.type_id)
  const direction = isSource ? 's2t' : 't2s'
  return {
    targetName: targetModel ? `${targetModel.name}(${targetModel.alias})` : `#${targetId}`,
    typeLabel: relType ? relType[direction] : '',
    description: rel.description || '',
  }
}

// ===== 初始化 =====
onMounted(async () => {
  try {
    const [modelsRes, groupsRes, typesRes] = await Promise.all([
      getAllModels(),
      getModelGroups(),
      getModelRelationTypes(),
    ])
    allModels.value = (modelsRes.data || []).filter(m => m.is_usable !== false)
    allModelGroups.value = groupsRes.data || []
    allRelationTypes.value = typesRes.data || []
  } catch {}
})
</script>

<style scoped>
.search-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: var(--page-height);
}

/* ===== Tab 栏 ===== */
.tab-bar {
  display: flex;
  gap: 32px;
  margin-top: 40px;
  margin-bottom: 32px;
}

.tab-item {
  font-size: 16px;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding-bottom: 8px;
  border-bottom: 3px solid transparent;
  transition: all 0.2s ease;
  user-select: none;
}

.tab-item:hover {
  color: var(--color-text-primary);
}

.tab-item.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  font-weight: 600;
}

/* ===== 搜索区 ===== */
.search-hero {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-bottom: 40px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 640px;
}

.search-input {
  flex: 1;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 4px 20px;
  height: 40px;
  transition: box-shadow 0.2s ease;
}

.search-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 2px 16px rgba(37, 99, 235, 0.15);
}

.search-model-select {
  width: 160px;
}

.search-model-select :deep(.el-input__wrapper) {
  border-radius: 24px;
  height: 40px;
}

.search-btn {
  border-radius: 24px;
  padding: 0 28px;
  height: 40px;
  font-weight: 600;
  flex-shrink: 0;
}

/* 实例搜索：复用全文检索的 hero 居中布局 */
.instance-search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 640px;
}

.instance-search-box .instance-model-select {
  flex: 1;
}

.instance-search-box .instance-model-select :deep(.el-input__wrapper) {
  border-radius: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 4px 20px;
  height: 40px;
  transition: box-shadow 0.2s ease;
}

.instance-search-box .instance-model-select :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 2px 16px rgba(37, 99, 235, 0.15);
}

.instance-search-box .search-btn {
  border-radius: 24px;
  padding: 0 28px;
  height: 40px;
  font-weight: 600;
  flex-shrink: 0;
}

/* ===== 结果区 ===== */
.results-area {
  width: 100%;
  max-width: 960px;
}

.result-count {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-bottom: 12px;
}

.result-count strong {
  color: var(--color-primary);
}

/* 数据标签 */
.data-tag {
  display: inline-block;
  font-size: 12px;
  padding: 1px 6px;
  margin: 2px 4px 2px 0;
  background: var(--color-muted);
  border-radius: 4px;
  color: var(--color-text-secondary);
}

.data-tag strong {
  color: var(--color-text-primary);
}

/* 分页 */
.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

/* 条件构建器外层间距 */

/* 条件构建器 */
.condition-builder {
  margin-bottom: 16px;
  padding: 16px;
  background: var(--color-muted);
  border-radius: var(--radius-md);
}

.condition-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.cond-field {
  width: 180px;
}

.cond-op {
  width: 120px;
}

.cond-val {
  flex: 1;
}

.condition-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

/* 模型分组 */
.model-group-section {
  margin-bottom: 24px;
}

.model-group-header {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 2px solid var(--color-primary);
}

/* 模型列表 */
.model-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.model-card {
  padding: 12px 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.15s ease;
}

.model-card:hover {
  border-color: var(--color-primary);
  background: var(--color-muted);
}

.model-card.expanded {
  border-color: var(--color-primary);
  background: rgba(37, 99, 235, 0.03);
}

.model-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.model-card-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.model-card-name {
  font-weight: 600;
  color: var(--color-text-primary);
}

.model-card-desc {
  margin-top: 4px;
  font-size: 12px;
  color: var(--color-text-muted);
}

.expand-icon {
  transition: transform 0.2s ease;
  color: var(--color-text-muted);
}

.expand-icon.rotated {
  transform: rotate(180deg);
}

/* 字段内联展开 */
.fields-inline {
  margin: -2px 0 6px 0;
  padding: 12px 16px;
  border: 1px solid var(--color-primary);
  border-top: none;
  border-radius: 0 0 var(--radius-md) var(--radius-md);
  background: rgba(37, 99, 235, 0.02);
}

.fields-loading {
  text-align: center;
  padding: 16px;
  font-size: 13px;
  color: var(--color-text-muted);
}

/* 字段分组 */
.field-group {
  margin-bottom: 12px;
}

.field-group:last-child {
  margin-bottom: 0;
}

.field-group-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
  padding-left: 2px;
}

.unique-tag {
  margin-left: 6px;
  vertical-align: middle;
}

/* 模型关系 */
.model-relations {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed var(--color-border);
}

.relations-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}

.relation-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  font-size: 13px;
}

.relation-type {
  color: var(--color-primary);
  font-weight: 500;
}

.relation-target {
  font-weight: 600;
  color: var(--color-text-primary);
}

.relation-desc {
  color: var(--color-text-muted);
  font-size: 12px;
}
</style>
