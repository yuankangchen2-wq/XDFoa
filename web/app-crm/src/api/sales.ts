import axios from 'axios'

const http = axios.create({
  baseURL: '/api/sales',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const salesApi = {
  // 销售订单
  listOrders: (params: any) => http.get('/orders', { params }).then(r => r.data),
  getOrder: (id: number) => http.get(`/orders/${id}`).then(r => r.data),
  createOrder: (data: any) => http.post('/orders', data).then(r => r.data),
  approveOrder: (id: number) => http.post(`/orders/${id}/approve`).then(r => r.data),
  shipOrder: (id: number) => http.post(`/orders/${id}/ship`).then(r => r.data),
  completeOrder: (id: number) => http.post(`/orders/${id}/complete`).then(r => r.data),
  // 合同
  listContracts: (params: any) => http.get('/contracts', { params }).then(r => r.data),
  getContract: (id: number) => http.get(`/contracts/${id}`).then(r => r.data),
  createContract: (data: any) => http.post('/contracts', data).then(r => r.data),
  approveContract: (id: number) => http.post(`/contracts/${id}/approve`).then(r => r.data),
  archiveContract: (id: number) => http.post(`/contracts/${id}/archive`).then(r => r.data),
  // 回款
  addPayment: (contractId: number, data: any) => http.post(`/contracts/${contractId}/payments`, data).then(r => r.data),
  listPayments: (contractId: number) => http.get(`/contracts/${contractId}/payments`).then(r => r.data)
}
