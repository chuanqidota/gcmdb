import { ref } from 'vue'

export const OPERATORS = [
  { label: '=', value: 'eq' },
  { label: '!=', value: 'ne' },
  { label: '>', value: 'gt' },
  { label: '>=', value: 'ge' },
  { label: '<', value: 'lt' },
  { label: '<=', value: 'le' },
  { label: '包含', value: 'contains' },
  { label: '开头是', value: 'startswith' },
  { label: '结尾是', value: 'endswith' },
  { label: '在列表', value: 'in' },
  { label: '模糊搜索', value: 'search' },
]

export function useConditionBuilder() {
  const conditions = ref([])

  function addCondition(group = 1) {
    conditions.value.push({ group, field: '', op: 'eq', val: '' })
  }

  function addOrGroup() {
    const maxGroup = conditions.value.reduce((max, c) => Math.max(max, c.group || 1), 0)
    conditions.value.push({ group: maxGroup + 1, field: '', op: 'eq', val: '' })
  }

  function removeCondition(index) {
    conditions.value.splice(index, 1)
  }

  function resetConditions() {
    conditions.value = []
  }

  function sortConditions() {
    conditions.value.sort((a, b) => (a.group || 1) - (b.group || 1))
  }

  function buildWhere(overrideConditions) {
    const conds = overrideConditions || conditions.value
    conds.sort((a, b) => (a.group || 1) - (b.group || 1))
    const where = []
    const group1 = []
    const orGroups = {}
    for (const c of conds) {
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
    return where
  }

  return {
    conditions,
    addCondition,
    addOrGroup,
    removeCondition,
    resetConditions,
    sortConditions,
    buildWhere,
    OPERATORS,
  }
}
