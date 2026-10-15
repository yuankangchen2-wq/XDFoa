import axios from 'axios'

const http = axios.create({
  baseURL: '/api/scrm',
  headers: { 'X-Tenant-Id': 'default' }
})

http.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

http.interceptors.response.use(
  r => r,
  err => {
    const msg = err.response?.data?.error || err.message || '请求失败'
    return Promise.reject(new Error(msg))
  }
)

export interface MeInfo {
  name: string
  title: string
  city: string
  store: string
  unread: number
  roles: string[]
}

export interface TodoItem {
  key: string
  label: string
  group: string
  count: number
  route: string
  query?: string
}

export interface Activity {
  id: number
  name: string
  source_system: string
  category: string
  type: string
  target_object: string
  intention_dept: string
  intention_project: string
  start_time: string
  end_time: string
  remark: string
  external_id: string
  code: string
  status: number
  creator: string
  creator_dept: string
  created_at: string
  stats?: Record<string, string>
}

export interface Lead {
  id: number
  code: string
  name: string
  phone: string
  student_no: string
  student_tag: string
  new_old_type: string
  school: string
  grade: string
  district: string
  lead_type: string
  enter_type: string
  status: string
  intention_level: string
  intention_grade: string
  intention_dept: string
  intention_project: string
  intention_campus: string
  intention_subject: string
  channel_l1: string
  channel_l2: string
  channel_l3: string
  activity_name: string
  task_name: string
  owner_name: string
  assign_time: string
  handover_time: string
  last_follow_at: string
  last_follow_content: string
  follow_count: number
  next_follow_at: string
  archive_rule: string
  archive_way: string
  is_duplicate: number
  creator: string
  remark: string
  created_at: string
}

export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface ImportTask {
  id: number
  task_name: string
  activity_name: string
  import_method: string
  data_source: string
  operator: string
  serial_no: string
  batch_no: string
  status: string
  file_name: string
  total_count: number
  success_count: number
  new_count: number
  active_count: number
  fail_count: number
  duplicate_count: number
  invalid_count: number
  uniqueness_fail_count: number
  dedup_hit_count: number
  dedup_condition_count: number
  dedup_details: string
  created_at: string
}

export interface DataRecord {
  id: number
  related_name: string
  related_phone: string
  related_student_no: string
  activity_student_name: string
  activity_phone: string
  import_type: string
  activity_grade: string
  related_grade: string
  enter_way: string
  lead_enter_type: string
  fail_reason: string
  channel_l1: string
  channel_l2: string
  channel_l3: string
  activity_name: string
  intention_dept: string
  intention_project: string
  owner_name: string
  handled: number
}

export interface ExceptionLead {
  id: number
  related_name: string
  related_phone: string
  related_student_no: string
  activity_student_name: string
  activity_phone: string
  channel_l1: string
  channel_l2: string
  channel_l3: string
  system_source_l1: string
  system_source_l2: string
  activity_name: string
  data_remark: string
  intention_dept: string
  intention_project: string
  status: string
}

export interface Employee {
  id: string
  name: string
  role: string
  campus: string
}

export interface ScheduleEvent {
  id: number
  title: string
  type: string
  status: string
  event_date: string
  event_time: string
}

export const scrmApi = {
  me: () => http.get<MeInfo>('/me').then(r => r.data),
  dict: (kind: string) => http.get('/dict', { params: { kind } }).then(r => r.data),
  channels: () => http.get('/dict?kind=channels').then(r => r.data),
  employees: () => http.get<Employee[]>('/employees').then(r => r.data),

  todos: (role: string) => http.get<TodoItem[]>('/home/todos', { params: { role } }).then(r => r.data),
  ranking: (month: string, campus: string) =>
    http.get('/home/ranking', { params: { month, campus } }).then(r => r.data),
  gradeDist: (month: string, campus: string) =>
    http.get('/home/grade-dist', { params: { month, campus } }).then(r => r.data),

  schedules: (month: string) => http.get<ScheduleEvent[]>('/schedules', { params: { month } }).then(r => r.data),
  createSchedule: (data: any) => http.post('/schedules', data).then(r => r.data),

  activities: (params: any) => http.get<PageResult<Activity>>('/activities', { params }).then(r => r.data),
  allActivities: (keyword = '') => http.get<Activity[]>('/activities/all', { params: { keyword } }).then(r => r.data),
  createActivity: (data: any) => http.post('/activities', data).then(r => r.data),
  updateActivity: (id: number, data: any) => http.put(`/activities/${id}`, data).then(r => r.data),

  leads: (params: any) => http.get<PageResult<Lead>>('/leads', { params }).then(r => r.data),
  leadStat: (params: any) => http.get('/leads/stat', { params }).then(r => r.data),
  exportLeadsUrl: (params: any) => `/api/scrm/leads/export?${new URLSearchParams(params).toString()}`,
  handoverLeads: (ids: number[], receiver_id: string, receiver_name: string) =>
    http.post('/leads/handover', { ids, receiver_id, receiver_name }).then(r => r.data),
  batchAdjustLeads: (ids: number[], fields: Record<string, string>) =>
    http.post('/leads/batch-adjust', { ids, fields }).then(r => r.data),
  batchArchiveLeads: (ids: number[], rule: string, way: string, reason: string) =>
    http.post('/leads/batch-archive', { ids, rule, way, reason }).then(r => r.data),
  activateLeads: (ids: number[]) => http.post('/leads/activate', { ids }).then(r => r.data),

  exceptions: (params: any) => http.get<PageResult<ExceptionLead>>('/exceptions', { params }).then(r => r.data),
  ignoreExceptions: (ids: number[]) => http.post('/exceptions/ignore', { ids }).then(r => r.data),
  adjustExceptions: (ids: number[], fields: Record<string, string>) =>
    http.post('/exceptions/adjust', { ids, fields }).then(r => r.data),

  taskTemplates: (keyword = '') =>
    http.get<ImportTask[]>('/task-templates', { params: { keyword } }).then(r => r.data),
  createTaskTemplate: (data: any) => http.post('/task-templates', data).then(r => r.data),
  importLogs: (params: any) => http.get<PageResult<ImportTask>>('/import-logs', { params }).then(r => r.data),
  dedupDetails: (id: number) => http.get(`/import-logs/${id}/dedup`).then(r => r.data),
  failFileUrl: (id: number) => `/api/scrm/import-logs/${id}/fail-file`,
  templateUrl: (method: string) => `/api/scrm/imports/template?method=${method}`,
  uploadImport: (form: FormData) =>
    http.post('/imports', form, { headers: { 'Content-Type': 'multipart/form-data' } }).then(r => r.data),

  dataRecords: (params: any) => http.get<PageResult<DataRecord>>('/data-records', { params }).then(r => r.data),
  handoverDataRecords: (ids: number[], receiver: string) =>
    http.post('/data-records/handover', { ids, receiver }).then(r => r.data),
  adjustDataRecords: (ids: number[], fields: Record<string, string>) =>
    http.post('/data-records/adjust', { ids, fields }).then(r => r.data)
}

// 手机号脱敏：130****3868
export function maskPhone(phone: string): string {
  if (!phone || phone.length < 7) return phone || ''
  return phone.slice(0, 3) + '****' + phone.slice(-4)
}

export function fmtTime(s: string): string {
  if (!s) return ''
  return s.replace('T', ' ').slice(0, 19)
}
