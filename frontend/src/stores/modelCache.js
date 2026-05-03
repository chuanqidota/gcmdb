import { defineStore } from 'pinia'
import { ref } from 'vue'
import { listModelGroup } from '../api/modelGroup'
import { listModel } from '../api/model'

const TTL = 5 * 60 * 1000 // 5 minutes

export const useModelCacheStore = defineStore('modelCache', () => {
  const groups = ref([])
  const models = ref([])
  const lastFetch = ref(0)

  async function fetchAll(force = false) {
    if (!force && Date.now() - lastFetch.value < TTL && groups.value.length > 0) {
      return { groups: groups.value, models: models.value }
    }
    const [g, m] = await Promise.all([
      listModelGroup({ limit: 1000 }),
      listModel({ limit: 1000 }),
    ])
    groups.value = g.data.results || []
    models.value = m.data.results || []
    lastFetch.value = Date.now()
    return { groups: groups.value, models: models.value }
  }

  function invalidate() {
    lastFetch.value = 0
  }

  return { groups, models, fetchAll, invalidate }
})
