import axios from 'axios'

const http = axios.create({
  baseURL: '/api/analytics',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const analyticsApi = {
  dashboard: () => http.get('/dashboard').then(r => r.data),
  // 销售
  listSales: (params: any) => http.get('/sales', { params }).then(r => r.data),
  recordSales: (data: any) => http.post('/sales', data).then(r => r.data),
  salesSummary: () => http.get('/sales/summary').then(r => r.data),
  // 漏斗
  listFunnel: (params: any) => http.get('/funnel', { params }).then(r => r.data),
  recordFunnel: (data: any) => http.post('/funnel', data).then(r => r.data),
  // 客户统计
  listCustomerStats: (params: any) => http.get('/customers', { params }).then(r => r.data),
  recordCustomerStat: (data: any) => http.post('/customers', data).then(r => r.data)
}
