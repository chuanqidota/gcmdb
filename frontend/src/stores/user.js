import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMe } from '../api/auth'

export const useUserStore = defineStore('user', () => {
  const username = ref('')
  const isAdmin = ref(false)
  const token = ref('')
  const userId = ref(null)

  async function fetchUser() {
    try {
      const res = await getMe()
      const data = res.data || {}
      username.value = data.username || ''
      isAdmin.value = data.is_admin || false
      token.value = data.token || ''
      userId.value = data.id || null
    } catch {
      // will be redirected by 401 interceptor
    }
  }

  function clearUser() {
    username.value = ''
    isAdmin.value = false
    token.value = ''
    userId.value = null
  }

  return { username, isAdmin, token, userId, fetchUser, clearUser }
})
