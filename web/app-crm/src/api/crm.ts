import axios from 'axios'

const http = axios.create({
  baseURL: '/api/crm',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const crmApi = {
  // 客户
  listCustomers: (params: any) => http.get('/customers', { params }).then(r => r.data),
  getCustomer: (id: number) => http.get(`/customers/${id}`).then(r => r.data),
  createCustomer: (data: any) => http.post('/customers', data).then(r => r.data),
  updateCustomer: (id: number, data: any) => http.put(`/customers/${id}`, data).then(r => r.data),
  deleteCustomer: (id: number) => http.delete(`/customers/${id}`).then(r => r.data),
  // 联系人
  addContact: (customerId: number, data: any) => http.post(`/customers/${customerId}/contacts`, data).then(r => r.data),
  listContacts: (customerId: number) => http.get(`/customers/${customerId}/contacts`).then(r => r.data),
  deleteContact: (id: number) => http.delete(`/contacts/${id}`).then(r => r.data),
  // 跟进
  addFollowUp: (customerId: number, data: any) => http.post(`/customers/${customerId}/follow-ups`, data).then(r => r.data),
  listFollowUps: (customerId: number) => http.get(`/customers/${customerId}/follow-ups`).then(r => r.data)
}
