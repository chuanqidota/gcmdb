import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const authRequest = axios.create({
  baseURL: '/v1/auth',
  timeout: 10000,
  withCredentials: true,
})

authRequest.interceptors.response.use(
  (res) => {
    if (res.data.code === 0) {
      ElMessage.error(res.data.msg || '请求失败')
      return Promise.reject(new Error(res.data.msg))
    }
    return res.data
  },
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('gcmdb_logged_in')
      router.push('/login')
      ElMessage.error('会话已过期，请重新登录')
      return Promise.reject(err)
    }
    const msg = err.response?.data?.msg || err.message || '网络错误'
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)

export const login = (data) => authRequest.post('/login', data)
export const logout = () => authRequest.post('/logout')
export const getMe = () => authRequest.get('/me')
export const changePassword = (data) => authRequest.post('/change-password', data)
export const resetPassword = (data) => authRequest.post('/reset-password', data)

// 用户管理（管理员）
export const listUsers = () => authRequest.get('/users')
export const createUser = (data) => authRequest.post('/users', data)
export const patchUser = (id, data) => authRequest.patch(`/users/${id}`, data)
