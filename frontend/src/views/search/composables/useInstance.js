import { ref, reactive, computed } from 'vue'
import { getModelDetail, searchInstance, getInstanceTopology } from '../../../api/search'

export function useInstance(allModels) {
  const inst = ref({
    modelId: null,
    fields: [],
    conditions: [],
    globalSearch: '',
    results: [],
    columns: [],
    columnLabelMap: {},
    count: 0,
    page: 1,
    pageSize: 10,
    loading: false,
    topoLoading: false,
  })

  const instRelationCache = reactive({})
  const instRelations = ref([])

  // ===== 关联抽屉 =====
  const relDrawerVisible = ref(false)
  const relDrawerStack = ref([])
  const relDrawerCurrent = computed(() => relDrawerStack.value[relDrawerStack.value.length - 1] || null)

  async function openRelDrawer(instanceId, modelName, modelAlias) {
    relDrawerStack.value = []
    await pushRelDrawer(instanceId, modelName, modelAlias)
    relDrawerVisible.value = true
  }

  async function pushRelDrawer(instanceId, modelName, modelAlias) {
    relDrawerStack.value.push({ instanceId, modelName, modelAlias, fields: [], instanceData: {}, loading: true })
    const idx = relDrawerStack.value.length - 1
    try {
      const [detailRes, topoRes] = await Promise.all([
        getModelDetail(modelAlias),
        getInstanceTopology(instanceId),
      ])
      relDrawerStack.value[idx].fields = (detailRes.data?.model_fields || []).slice(0, 10)
      relDrawerStack.value[idx].instanceData = topoRes.data?.instance?.data || {}
    } catch {
      relDrawerStack.value[idx].fields = []
    } finally {
      relDrawerStack.value[idx].loading = false
    }
    // 加载该实例的关联
    loadRelDrawerRelations(instanceId)
  }

  function popRelDrawer() {
    relDrawerStack.value.pop()
  }

  function popToRelDrawer(idx) {
    relDrawerStack.value = relDrawerStack.value.slice(0, idx + 1)
  }

  function closeRelDrawer() {
    relDrawerStack.value = []
    relDrawerVisible.value = false
  }

  // 抽屉内关联数据
  const relDrawerRelationCache = reactive({})

  async function loadRelDrawerRelations(instanceId) {
    relDrawerRelationCache[instanceId] = { groups: [], loading: true, activeTab: '', tabData: {} }
    try {
      const res = await getInstanceTopology(instanceId)
      const topoData = res.data || {}
      const upstream = (topoData.upstream || []).map(r => ({ ...r, direction: 'upstream' }))
      const downstream = (topoData.downstream || []).map(r => ({ ...r, direction: 'downstream' }))
      const allRels = [...upstream, ...downstream]

      const groupMap = new Map()
      for (const rel of allRels) {
        const key = rel.model_id
        if (!groupMap.has(key)) {
          groupMap.set(key, { model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, direction: rel.direction, instances: [] })
        }
        groupMap.get(key).instances.push(rel)
      }

      const groups = [...groupMap.values()]
      relDrawerRelationCache[instanceId].groups = groups
      relDrawerRelationCache[instanceId].loading = false

      if (groups.length) {
        const tabToLoad = String(groups[0].model_id)
        relDrawerRelationCache[instanceId].activeTab = tabToLoad
        onRelDrawerTabChange(instanceId, groups[0].model_id)
      }
    } catch {
      relDrawerRelationCache[instanceId] = { groups: [], loading: false, activeTab: '', tabData: {} }
    }
  }

  async function onRelDrawerTabChange(instanceId, modelId) {
    const cache = relDrawerRelationCache[instanceId]
    if (!cache || cache.tabData[modelId]) return

    const group = cache.groups.find(g => g.model_id === modelId)
    if (!group) return

    cache.tabData[modelId] = { instances: group.instances, fields: [], loading: true, page: 1, pageSize: 10 }

    try {
      const detailRes = await getModelDetail(group.model_alias)
      const fields = (detailRes.data?.model_fields || []).slice(0, 8)
      cache.tabData[modelId].fields = fields
    } catch {
      cache.tabData[modelId].fields = []
    } finally {
      cache.tabData[modelId].loading = false
    }
  }

  function handleRelDrawerTabChange(instanceId, modelId) {
    onRelDrawerTabChange(instanceId, modelId)
  }

  // ===== 搜索相关 =====

  async function onModelChange(id) {
    inst.value.fields = []
    inst.value.conditions = []
    inst.value.results = []
    inst.value.columns = []
    inst.value.count = 0
    inst.value.page = 1
    Object.keys(instRelationCache).forEach(k => delete instRelationCache[k])
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
      const group1 = []
      const orGroups = {}
      if (inst.value.globalSearch) {
        group1.push({ search: inst.value.globalSearch })
      }
      for (const c of inst.value.conditions) {
        if (!c.op || c.val === '' || c.val === undefined) continue
        let cond = null
        if (c.op === 'search') {
          cond = { search: c.val }
        } else if (!c.field) {
          continue
        } else if (c.op === 'in') {
          const vals = c.val.split(',').map(v => v.trim()).filter(Boolean)
          if (vals.length) cond = { in: { [c.field]: vals } }
        } else {
          cond = { [c.op]: { [c.field]: c.val } }
        }
        if (!cond) continue
        const g = c.group || 1
        if (g === 1) {
          group1.push(cond)
        } else {
          if (!orGroups[g]) orGroups[g] = []
          orGroups[g].push(cond)
        }
      }
      where.push(...group1)
      for (const g of Object.keys(orGroups)) {
        if (orGroups[g].length) where.push({ or: orGroups[g] })
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
    loadInstanceRelations()
  }

  function calcColWidth(label) {
    let w = 0
    for (const ch of label) {
      w += /[一-鿿]/.test(ch) ? 14 : 8
    }
    return Math.max(w + 24, 80)
  }

  function resetInstance() {
    inst.value.conditions = []
    inst.value.globalSearch = ''
    inst.value.results = []
    inst.value.columns = []
    inst.value.columnLabelMap = {}
    inst.value.count = 0
    inst.value.page = 1
    Object.keys(instRelationCache).forEach(k => delete instRelationCache[k])
  }

  function addOrGroup() {
    const maxGroup = inst.value.conditions.reduce((max, c) => Math.max(max, c.group || 1), 0)
    inst.value.conditions.push({ group: maxGroup + 1, field: '', op: 'eq', val: '' })
  }

  // Expand row: load related instances
  const onInstExpandChange = (row, expanded) => {
    if (expanded.length && !instRelationCache[row.id]) {
      loadInstRelations(row)
    }
  }

  const loadInstRelations = async (row) => {
    instRelationCache[row.id] = { groups: [], loading: true, activeTab: '', tabData: {} }
    try {
      const res = await getInstanceTopology(row.id)
      const topoData = res.data || {}
      const upstream = (topoData.upstream || []).map(r => ({ ...r, direction: 'upstream' }))
      const downstream = (topoData.downstream || []).map(r => ({ ...r, direction: 'downstream' }))
      const allRels = [...upstream, ...downstream]

      const groupMap = new Map()
      for (const rel of allRels) {
        const key = rel.model_id
        if (!groupMap.has(key)) {
          groupMap.set(key, { model_id: rel.model_id, model_name: rel.model_name, model_alias: rel.model_alias, direction: rel.direction, instances: [] })
        }
        groupMap.get(key).instances.push(rel)
      }

      const groups = [...groupMap.values()]
      instRelationCache[row.id].groups = groups
      instRelationCache[row.id].loading = false

      if (groups.length) {
        const tabToLoad = String(groups[0].model_id)
        instRelationCache[row.id].activeTab = tabToLoad
        onInstTabChange(row, Number(tabToLoad))
      }
    } catch {
      instRelationCache[row.id] = { groups: [], loading: false, activeTab: '', tabData: {} }
    }
  }

  const onInstTabChange = async (row, modelId) => {
    const cache = instRelationCache[row.id]
    if (!cache || cache.tabData[modelId]) return

    const group = cache.groups.find(g => g.model_id === modelId)
    if (!group) return

    cache.tabData[modelId] = { instances: group.instances, fields: [], loading: true, page: 1, pageSize: 10 }

    try {
      const detailRes = await getModelDetail(group.model_alias)
      const fields = (detailRes.data?.model_fields || []).slice(0, 8)
      cache.tabData[modelId].fields = fields
    } catch {
      cache.tabData[modelId].fields = []
    } finally {
      cache.tabData[modelId].loading = false
    }
  }

  function handleTabChange(instanceId, modelId) {
    onInstTabChange({ id: instanceId }, modelId)
  }

  // Load relations between current page instances
  async function loadInstanceRelations() {
    if (!inst.value.results.length) { instRelations.value = []; return }
    const model = allModels.value.find(m => m.id === inst.value.modelId)
    if (!model) return
    const idSet = new Set(inst.value.results.map(r => r.id))
    const topoPromises = inst.value.results.map(i => getInstanceTopology(i.id).catch(() => null))
    const topoResults = await Promise.all(topoPromises)
    const relations = []
    for (const res of topoResults) {
      if (!res?.data) continue
      const { instance: srcInst, upstream, downstream } = res.data
      if (!srcInst) continue
      for (const rel of (downstream || [])) {
        if (idSet.has(rel.instance_id)) {
          relations.push({ sourceId: srcInst.id, sourceLabel: `${model.name} #${srcInst.id}`, relationType: rel.relation_type || '', targetId: rel.instance_id, targetLabel: `${rel.model_name} #${rel.instance_id}` })
        }
      }
    }
    instRelations.value = relations
  }

  return {
    inst, instRelationCache, instRelations,
    onModelChange, doInstanceSearch, calcColWidth, resetInstance, addOrGroup,
    onInstExpandChange, onInstTabChange, handleTabChange,
    // 关联抽屉
    relDrawerVisible, relDrawerStack, relDrawerCurrent, relDrawerRelationCache,
    openRelDrawer, pushRelDrawer, popRelDrawer, popToRelDrawer, closeRelDrawer,
    handleRelDrawerTabChange,
  }
}
