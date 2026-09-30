import axios from 'axios'

const http = axios.create({
  baseURL: '/api/ticket',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

export const ticketApi = {
  list: (params: any) => http.get('/tickets', { params }).then(r => r.data),
  get: (id: number) => http.get(`/tickets/${id}`).then(r => r.data),
  create: (data: any) => http.post('/tickets', data).then(r => r.data),
  updateStatus: (id: number, status: number) => http.put(`/tickets/${id}/status`, { status }).then(r => r.data),
  assign: (id: number, assigneeId: string) => http.put(`/tickets/${id}/assign`, { assignee_id: assigneeId }).then(r => r.data),
  addReply: (id: number, data: any) => http.post(`/tickets/${id}/replies`, data).then(r => r.data),
  listReplies: (id: number) => http.get(`/tickets/${id}/replies`).then(r => r.data),
  rate: (id: number, score: number) => http.post(`/tickets/${id}/rate`, { score }).then(r => r.data)
}
