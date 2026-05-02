import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { listModelGroup } from '../api/modelGroup'
import { listModel } from '../api/model'

export const useModelStore = defineStore('model', () => {
  const groups = ref([])
  const models = ref([])
  const loaded = ref(false)

  async function fetchAll() {
    if (loaded.value) return
    try {
      const [groupRes, modelRes] = await Promise.all([
        listModelGroup({ limit: 1000 }),
        listModel({ limit: 1000 }),
      ])
      groups.value = groupRes.data?.results || []
      models.value = modelRes.data?.results || []
      loaded.value = true
    } catch {
      // ignore
    }
  }

  function getModelById(id) {
    return models.value.find(m => m.id === id)
  }

  function getModelByAlias(alias) {
    return models.value.find(m => m.alias === alias)
  }

  const modelMap = computed(() => {
    const map = {}
    for (const m of models.value) {
      map[m.id] = m
    }
    return map
  })

  function $reset() {
    groups.value = []
    models.value = []
    loaded.value = false
  }

  return { groups, models, loaded, fetchAll, getModelById, getModelByAlias, modelMap, $reset }
})
