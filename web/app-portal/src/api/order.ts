import axios from 'axios'

const http = axios.create({
  baseURL: '/api/order',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const orderApi = {
  list: (params: any) => http.get('/orders', { params }).then(r => r.data),
  get: (id: number) => http.get(`/orders/${id}`).then(r => r.data),
  create: (data: any) => http.post('/orders', data).then(r => r.data),
  pay: (id: number) => http.post(`/orders/${id}/pay`).then(r => r.data),
  cancel: (id: number) => http.post(`/orders/${id}/cancel`).then(r => r.data),
  ship: (id: number) => http.post(`/orders/${id}/ship`).then(r => r.data),
  complete: (id: number) => http.post(`/orders/${id}/complete`).then(r => r.data)
}
