import axios from 'axios'

const http = axios.create({
  baseURL: '/api/mdm',
  headers: { 'X-Tenant-Id': 'default' }
})

// 客户
export const customerApi = {
  list: (params: any) => http.get('/customers', { params }).then(r => r.data),
  get: (id: number) => http.get(`/customers/${id}`).then(r => r.data),
  create: (data: any) => http.post('/customers', data).then(r => r.data),
  update: (id: number, data: any) => http.put(`/customers/${id}`, data).then(r => r.data),
  remove: (id: number) => http.delete(`/customers/${id}`).then(r => r.data)
}

// 商品
export const productApi = {
  list: (params: any) => http.get('/products', { params }).then(r => r.data),
  get: (id: number) => http.get(`/products/${id}`).then(r => r.data),
  create: (data: any) => http.post('/products', data).then(r => r.data),
  update: (id: number, data: any) => http.put(`/products/${id}`, data).then(r => r.data),
  remove: (id: number) => http.delete(`/products/${id}`).then(r => r.data)
}

// 供应商
export const supplierApi = {
  list: (params: any) => http.get('/suppliers', { params }).then(r => r.data),
  get: (id: number) => http.get(`/suppliers/${id}`).then(r => r.data),
  create: (data: any) => http.post('/suppliers', data).then(r => r.data),
  update: (id: number, data: any) => http.put(`/suppliers/${id}`, data).then(r => r.data),
  remove: (id: number) => http.delete(`/suppliers/${id}`).then(r => r.data)
}

// 组织
export const orgApi = {
  list: (parentId = 0) => http.get('/organizations', { params: { parent_id: parentId } }).then(r => r.data),
  get: (id: number) => http.get(`/organizations/${id}`).then(r => r.data),
  create: (data: any) => http.post('/organizations', data).then(r => r.data),
  update: (id: number, data: any) => http.put(`/organizations/${id}`, data).then(r => r.data),
  remove: (id: number) => http.delete(`/organizations/${id}`).then(r => r.data)
}
