import axios from 'axios'

const http = axios.create({
  baseURL: '/api/audit',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const auditApi = {
  list: (params: any) => http.get('/logs', { params }).then(r => r.data),
  get: (id: number) => http.get(`/logs/${id}`).then(r => r.data),
  record: (data: any) => http.post('/logs', data).then(r => r.data),
  moduleStats: () => http.get('/stats/modules').then(r => r.data)
}
