import { ref } from 'vue'

export function usePagination(loadFn, defaultLimit = 10) {
  const list = ref([])
  const total = ref(0)
  const loading = ref(false)
  const page = ref(1)
  const pageSize = ref(defaultLimit)

  async function loadData() {
    loading.value = true
    try {
      const offset = (page.value - 1) * pageSize.value
      const res = await loadFn({ limit: pageSize.value, offset })
      const data = res.data || {}
      list.value = data.results || []
      total.value = data.count || 0
    } finally {
      loading.value = false
    }
  }

  function handlePageChange(p) {
    page.value = p
    loadData()
  }

  function handleSizeChange(s) {
    pageSize.value = s
    page.value = 1
    loadData()
  }

  return { list, total, loading, page, pageSize, loadData, handlePageChange, handleSizeChange }
}
