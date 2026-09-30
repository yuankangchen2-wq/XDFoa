import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000,
  withCredentials: true
})

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('casdoor_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

request.interceptors.response.use(
  (response) => {
    const res = response.data
    // Casdoor 返回格式: { status: 'ok'|'error', msg, data }
    if (res.status && res.status === 'error') {
      ElMessage.error(res.msg || '请求失败')
      return Promise.reject(new Error(res.msg))
    }
    return res
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('casdoor_token')
      localStorage.removeItem('casdoor_user')
      router.push('/login')
    } else {
      ElMessage.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

// Casdoor API 封装
export const authApi = {
  // 登录
  login: (data: { application: string; username: string; password: string; organization?: string; signinMethod?: string }) =>
    request.post('/login', data),
  // 登出
  logout: () => request.post('/logout'),
  // 获取当前用户
  getAccount: () => request.get('/get-account')
}

export const applicationApi = {
  // 获取应用列表
  getApplications: (params?: { owner?: string; page?: number; pageSize?: number }) =>
    request.get('/get-applications', { params }),
  // 获取组织的应用
  getOrganizationApplications: (organization: string) =>
    request.get('/get-organization-applications', { params: { organization } })
}

export const dashboardApi = {
  // 仪表盘数据
  getDashboard: () => request.get('/get-dashboard')
}

export default request
