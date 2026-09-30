import axios from 'axios'

const http = axios.create({
  baseURL: '/api/purchase',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const purchaseApi = {
  listOrders: (params: any) => http.get('/orders', { params }).then(r => r.data),
  getOrder: (id: number) => http.get(`/orders/${id}`).then(r => r.data),
  createOrder: (data: any) => http.post('/orders', data).then(r => r.data),
  approveOrder: (id: number) => http.post(`/orders/${id}/approve`).then(r => r.data),
  settleOrder: (id: number) => http.post(`/orders/${id}/settle`).then(r => r.data),
  listReceipts: (params: any) => http.get('/receipts', { params }).then(r => r.data),
  receiveGoods: (data: any) => http.post('/receipts', data).then(r => r.data)
}
