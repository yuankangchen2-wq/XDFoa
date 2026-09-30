import axios from 'axios'

const http = axios.create({
  baseURL: '/api/production',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const productionApi = {
  // BOM
  listBoms: (params: any) => http.get('/boms', { params }).then(r => r.data),
  createBom: (data: any) => http.post('/boms', data).then(r => r.data),
  enableBom: (id: number) => http.post(`/boms/${id}/enable`).then(r => r.data),
  // 工单
  listWorkOrders: (params: any) => http.get('/work-orders', { params }).then(r => r.data),
  createWorkOrder: (data: any) => http.post('/work-orders', data).then(r => r.data),
  getWorkOrder: (id: number) => http.get(`/work-orders/${id}`).then(r => r.data),
  // 领料
  issueMaterial: (id: number) => http.post(`/work-orders/${id}/issue`).then(r => r.data),
  listMaterialIssues: (params: any) => http.get('/material-issues', { params }).then(r => r.data),
  // 报工
  reportProduction: (id: number, data: any) => http.post(`/work-orders/${id}/report`, data).then(r => r.data)
}
