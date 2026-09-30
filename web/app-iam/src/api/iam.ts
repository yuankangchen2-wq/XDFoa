import axios from 'axios'

const http = axios.create({
  baseURL: '/api/iam',
  headers: { 'X-Tenant-Id': 'default' }
})

// 请求拦截器：注入 token
http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 认证 API
export const authApi = {
  login: (data: any) => http.post('/auth/login', data).then(r => r.data),
  refresh: (refreshToken: string) => http.post('/auth/refresh', { refresh_token: refreshToken }).then(r => r.data),
  verify: (token: string) => http.post('/auth/verify', { token }).then(r => r.data)
}

// 用户 API
export const userApi = {
  list: (params: any) => http.get('/users', { params }).then(r => r.data),
  get: (id: number) => http.get(`/users/${id}`).then(r => r.data),
  create: (data: any) => http.post('/users', data).then(r => r.data),
  update: (id: number, data: any) => http.put(`/users/${id}`, data).then(r => r.data),
  remove: (id: number) => http.delete(`/users/${id}`).then(r => r.data)
}

// 角色 API
export const roleApi = {
  list: (params: any) => http.get('/roles', { params }).then(r => r.data),
  create: (data: any) => http.post('/roles', data).then(r => r.data),
  remove: (id: number) => http.delete(`/roles/${id}`).then(r => r.data)
}

// 权限 API
export const permApi = {
  list: () => http.get('/permissions').then(r => r.data),
  listPolicies: (roleCode: string) => http.get('/policies', { params: { role_code: roleCode } }).then(r => r.data),
  addPolicy: (data: any) => http.post('/policies', data).then(r => r.data),
  removePolicy: (data: any) => http.delete('/policies', { data }).then(r => r.data)
}
