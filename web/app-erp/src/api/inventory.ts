import axios from 'axios'

const http = axios.create({
  baseURL: '/api/inventory',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const warehouseApi = {
  list: (params: any) => http.get('/warehouses', { params }).then(r => r.data),
  create: (data: any) => http.post('/warehouses', data).then(r => r.data)
}

export const stockApi = {
  list: (params: any) => http.get('/stocks', { params }).then(r => r.data),
  allocate: (data: any) => http.post('/stocks/allocate', data).then(r => r.data),
  release: (data: any) => http.post('/stocks/release', data).then(r => r.data),
  deduct: (data: any) => http.post('/stocks/deduct', data).then(r => r.data),
  stockIn: (data: any) => http.post('/stocks/in', data).then(r => r.data)
}

export const journalApi = {
  list: (params: any) => http.get('/journals', { params }).then(r => r.data)
}

export const transferApi = {
  list: (params: any) => http.get('/transfers', { params }).then(r => r.data),
  create: (data: any) => http.post('/transfers', data).then(r => r.data),
  confirmOut: (id: number) => http.post(`/transfers/${id}/out`).then(r => r.data),
  confirmIn: (id: number) => http.post(`/transfers/${id}/in`).then(r => r.data)
}

export const stocktakeApi = {
  list: (params: any) => http.get('/stocktakes', { params }).then(r => r.data),
  create: (data: any) => http.post('/stocktakes', data).then(r => r.data)
}
