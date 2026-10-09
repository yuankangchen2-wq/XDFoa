import axios from 'axios'

const http = axios.create({ baseURL: '/api/bff', timeout: 10000 })

http.interceptors.response.use(
  (res) => res.data,
  (err) => {
    console.error('[bff api]', err?.message)
    return Promise.reject(err)
  }
)

export interface MeInfo {
  user_id: string
  username: string
  title: string
  join_days: number
  join_date: string
  roles: string[]
}

export interface Employee {
  id: number
  username: string
  real_name: string
  email: string
  phone: string
}

export interface Approval {
  id: number
  title: string
  type: string
  serial_no: string
  applicant: string
  applicant_id: string
  status: number
  current_node: string
  created_at: string
}

export interface ScheduleEvent {
  id: number
  title: string
  event_date: string
}

export interface Favorite {
  id: number
  name: string
  icon: string
  path: string
  sort: number
}

export const bffApi = {
  me: () => http.get<any, MeInfo>('/me'),
  health: () => http.get<any, { items: any[] }>('/health'),
  searchEmployees: (keyword: string) =>
    http.get<any, { items: Employee[] }>('/employees', { params: { keyword } }),

  listApprovals: (params: { type?: string; status?: string; keyword?: string }) =>
    http.get<any, { items: Approval[] }>('/approvals', { params }),
  approvalCounts: (type = '') =>
    http.get<any, Record<string, number>>('/approvals/counts', { params: { type } }),

  listSchedules: (month: string) =>
    http.get<any, { items: ScheduleEvent[] }>('/schedules', { params: { month } }),

  listFavorites: () => http.get<any, { items: Favorite[] }>('/favorites'),
  addFavorite: (data: { name: string; icon?: string; path?: string }) =>
    http.post<any, Favorite>('/favorites', data),
  deleteFavorite: (id: number) => http.delete(`/favorites/${id}`)
}
