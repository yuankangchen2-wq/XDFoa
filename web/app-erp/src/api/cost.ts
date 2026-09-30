import axios from 'axios'

const http = axios.create({
  baseURL: '/api/cost',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const costApi = {
  // 成本中心
  listCenters: (params: any) => http.get('/centers', { params }).then(r => r.data),
  createCenter: (data: any) => http.post('/centers', data).then(r => r.data),
  // 产品成本
  listProductCosts: (params: any) => http.get('/products', { params }).then(r => r.data),
  calculateProductCost: (data: any) => http.post('/products', data).then(r => r.data),
  collectCost: (data: any) => http.post('/collect', data).then(r => r.data),
  // 成本记录
  listRecords: (params: any) => http.get('/records', { params }).then(r => r.data),
  addRecord: (data: any) => http.post('/records', data).then(r => r.data)
}
