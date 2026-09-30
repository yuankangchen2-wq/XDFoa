import axios from 'axios'

const http = axios.create({
  baseURL: '/api/finance',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const financeApi = {
  // 应收
  listReceivables: (params: any) => http.get('/receivables', { params }).then(r => r.data),
  createReceivable: (data: any) => http.post('/receivables', data).then(r => r.data),
  // 应付
  listPayables: (params: any) => http.get('/payables', { params }).then(r => r.data),
  createPayable: (data: any) => http.post('/payables', data).then(r => r.data),
  // 收付款
  receive: (data: any) => http.post('/payments/receive', data).then(r => r.data),
  makePayment: (data: any) => http.post('/payments/make', data).then(r => r.data),
  listPayments: (params: any) => http.get('/payments', { params }).then(r => r.data)
}
