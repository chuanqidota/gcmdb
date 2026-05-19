import { ref, computed } from 'vue'
import { llmSearch, llmClearSession } from '../../../api/llmSearch'

export function useLlmSearch(allModels) {
  const llm = ref({
    message: '',
    sessionId: '',
    modelFilter: [],
    loading: false,
    aiAnalysis: '',
    results: [],
    thinking: '',
    tokenUsed: 0,
    history: [], // 对话历史 [{ role, message, results, aiAnalysis }]
    showHistory: false,
  })

  const hasResults = computed(() => llm.value.results.length > 0 || llm.value.aiAnalysis)

  async function doLlmSearch() {
    const msg = llm.value.message.trim()
    if (!msg || llm.value.loading) return

    llm.value.loading = true
    llm.value.thinking = ''
    llm.value.aiAnalysis = ''
    llm.value.results = []

    try {
      const res = await llmSearch({
        message: msg,
        session_id: llm.value.sessionId,
        model_filter: llm.value.modelFilter.length > 0 ? llm.value.modelFilter : undefined,
      })

      const data = res.data || {}
      llm.value.sessionId = data.session_id || ''
      llm.value.aiAnalysis = data.ai_analysis || ''
      llm.value.results = data.results || []
      llm.value.thinking = data.thinking || ''
      llm.value.tokenUsed = data.token_used || 0

      // 保存到历史
      llm.value.history.unshift({
        role: 'user',
        message: msg,
        results: data.results || [],
        aiAnalysis: data.ai_analysis || '',
        timestamp: new Date().toLocaleTimeString(),
      })
    } catch (err) {
      llm.value.aiAnalysis = ''
      llm.value.results = []
    } finally {
      llm.value.loading = false
      llm.value.message = ''
    }
  }

  async function doClearSession() {
    if (llm.value.sessionId) {
      try {
        await llmClearSession({ session_id: llm.value.sessionId })
      } catch {}
    }
    llm.value.sessionId = ''
    llm.value.aiAnalysis = ''
    llm.value.results = []
    llm.value.thinking = ''
    llm.value.tokenUsed = 0
    llm.value.history = []
  }

  return {
    llm,
    hasResults,
    doLlmSearch,
    doClearSession,
  }
}
