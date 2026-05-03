import { ref, computed } from 'vue'
import { getModelDetail } from '../../../api/search'

export function useModelBrowse(allModels, allModelGroups, allRelationTypes) {
  const modelSearch = ref('')
  const expandedModelId = ref(null)
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

      const uniqueSet = new Set()
      for (const u of uniques) {
        if (u.fields) {
          u.fields.split(',').map(f => f.trim()).filter(Boolean).forEach(f => uniqueSet.add(f))
        }
      }

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

  return { modelSearch, expandedModelId, modelFieldsCache, groupedModels, toggleFields, resolveRelation }
}
