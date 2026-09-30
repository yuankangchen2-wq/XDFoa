import axios from 'axios'

const http = axios.create({
  baseURL: '/api/opp',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const oppApi = {
  list: (params: any) => http.get('/opportunities', { params }).then(r => r.data),
  get: (id: number) => http.get(`/opportunities/${id}`).then(r => r.data),
  create: (data: any) => http.post('/opportunities', data).then(r => r.data),
  update: (id: number, data: any) => http.put(`/opportunities/${id}`, data).then(r => r.data),
  remove: (id: number) => http.delete(`/opportunities/${id}`).then(r => r.data),
  transition: (id: number, data: any) => http.post(`/opportunities/${id}/transition`, data).then(r => r.data),
  convertOrder: (id: number, data: any) => http.post(`/opportunities/${id}/convert-order`, data).then(r => r.data)
}
