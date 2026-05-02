import axios from 'axios'
import { ElMessage } from 'element-plus'

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
    const msg = err.response?.data?.msg || err.message || '网络错误'
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)

export const login = (data) => authRequest.post('/login', data)
export const logout = () => authRequest.post('/logout')
export const getMe = () => authRequest.get('/me')
