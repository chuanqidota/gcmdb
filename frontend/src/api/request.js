import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/v1/cmdb',
  timeout: 10000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('api_key')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
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

export default request
