import { ref } from 'vue'
import { fulltextSearch, searchInstance } from '../../../api/search'

export function useFulltext(allModels) {
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
    if (!hasKeyword && !hasModel) return

    ft.value.loading = true
    try {
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
        const selectedModel = allModels?.value?.find(m => m.alias === ft.value.modelAlias)
        const modelName = selectedModel?.name || ft.value.modelAlias
        ft.value.results = results.map(r => ({
          id: r.id,
          model_name: modelName,
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

  return { ft, doFulltextSearch }
}
