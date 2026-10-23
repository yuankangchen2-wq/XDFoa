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
  phone?: string
  city: string
  store: string
  unread: number
  roles: string[]
}

export interface FollowLog {
  id: number
  lead_id: number
  lead_name: string
  type: string
  content: string
  result: string
  duration_sec: number
  customer_phone: string
  operator: string
  operator_id: string
  created_at: string
}

export interface NotificationItem {
  id: number
  type: string
  title: string
  content: string
  biz_type: string
  biz_id: string
  route: string
  is_read: number
  created_at: string
}

export interface StudentHit {
  lead_id: number
  name: string
  phone: string
  student_no: string
  school: string
  grade: string
}

export interface SearchResult {
  keyword: string
  leads: Lead[]
  students: StudentHit[]
  opportunities: { id: number; code: string; title: string }[]
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
  call_status?: string
  last_call_at?: string
  self_call_count?: number
  total_call_count?: number
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
  leadDetail: (id: number) => http.get<Lead>(`/leads/${id}`).then(r => r.data),
  followLogs: (leadId: number) => http.get<FollowLog[]>(`/leads/${leadId}/follow-logs`).then(r => r.data),
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
    http.post('/data-records/adjust', { ids, fields }).then(r => r.data),

  // 顶部栏：全局搜索
  search: (q: string) => http.get<SearchResult>('/search', { params: { q } }).then(r => r.data),

  // 云座席双呼
  callDial: (leadId: number, phone?: string) =>
    http.post<FollowLog>('/calls/dial', { lead_id: leadId, phone: phone || '' }).then(r => r.data),
  callComplete: (logId: number, connected: boolean, durationSec: number) =>
    http.post('/calls/complete', { log_id: logId, connected, duration_sec: durationSec }).then(r => r.data),
  recentCalls: () => http.get<FollowLog[]>('/calls').then(r => r.data),

  // 通知中心
  notifications: (params: any) =>
    http.get<PageResult<NotificationItem>>('/notifications', { params }).then(r => r.data),
  unreadCount: () => http.get<{ unread: number }>('/notifications/unread-count').then(r => r.data.unread),
  markNotificationsRead: (ids: number[] | null, all = false) =>
    http.post('/notifications/read', { ids: ids || [], all }).then(r => r.data),

  // ===== 二期：线索管理工作台 =====
  recycleLeads: (ids: number[], reason: string) =>
    http.post('/leads/recycle', { ids, reason }).then(r => r.data),
  leadTaskSnapshot: (scope: string) =>
    http.get<{ items: any[] }>('/lead-task-snapshot', { params: { scope } }).then(r => r.data.items),
  employeeTasks: (params: any) =>
    http.get<PageResult<any>>('/employee-tasks', { params }).then(r => r.data),

  // ===== 三期：商机管理 =====
  opportunities: (params: any) =>
    http.get<PageResult<Opportunity>>('/opportunities', { params }).then(r => r.data),
  createOpportunity: (data: any) => http.post('/opportunities', data).then(r => r.data),
  updateOpportunity: (id: number, fields: Record<string, any>) =>
    http.put(`/opportunities/${id}`, fields).then(r => r.data),
  scoreRules: () => http.get<Record<string, string>>('/opportunities/score-rules').then(r => r.data),
  assignOpportunities: (data: any) => http.post('/opportunities/assign', data).then(r => r.data),
  recycleOpportunities: (ids: number[], reason: string) =>
    http.post('/opportunities/recycle', { ids, reason }).then(r => r.data),
  archiveOpportunities: (ids: number[], rule: string) =>
    http.post('/opportunities/batch-archive', { ids, rule }).then(r => r.data),
  activateOpportunities: (ids: number[]) =>
    http.post('/opportunities/activate', { ids }).then(r => r.data),
  push91: (ids: number[]) => http.post('/opportunities/push-91', { ids }).then(r => r.data),
  pushHedgehog: (ids: number[]) => http.post('/opportunities/push-hedgehog', { ids }).then(r => r.data),
  applyDeal: (data: any) => http.post('/opportunities/apply-deal', data).then(r => r.data),
  assignRecords: (params: any) =>
    http.get<PageResult<any>>('/assign-records', { params }).then(r => r.data),
  inventory: (dim: string) => http.get<any>('/inventory', { params: { dim } }).then(r => r.data),

  // 业绩申诉 / 审核中心 / 成单日志
  appeals: (params: any) => http.get<PageResult<any>>('/appeals', { params }).then(r => r.data),
  createAppeal: (data: any) => http.post('/appeals', data).then(r => r.data),
  handleAppeals: (ids: number[], pass: boolean) =>
    http.post('/appeals/handle', { ids, pass }).then(r => r.data),
  approvals: (params: any) => http.get<PageResult<any>>('/approvals', { params }).then(r => r.data),
  handleApprovals: (ids: number[], pass: boolean) =>
    http.post('/approvals/handle', { ids, pass }).then(r => r.data),
  dealLogs: (params: any) => http.get<PageResult<any>>('/deal-logs', { params }).then(r => r.data),

  // 退调转工单 / 积分商城订单
  tickets: (params: any) => http.get<PageResult<any>>('/tickets', { params }).then(r => r.data),
  mallOrders: (params: any) => http.get<PageResult<any>>('/mall-orders', { params }).then(r => r.data),

  // 用户级列/筛选配置
  getGridConfig: (pageKey: string, cfgType = 'column') =>
    http.get<{ config: string }>('/grid-config', { params: { page_key: pageKey, cfg_type: cfgType } }).then(r => r.data.config),
  saveGridConfig: (pageKey: string, cfgType: string, config: string) =>
    http.post('/grid-config', { page_key: pageKey, cfg_type: cfgType, config }).then(r => r.data)
}

export interface Opportunity {
  id: number
  code: string
  name: string
  phone: string
  student_no: string
  new_old_type: string
  grade: string
  school: string
  district: string
  add_wechat: number
  deduct_wechat: number
  status: string
  enter_type: string
  opp_type: string
  category: string
  is_intl_travel: number
  level: string
  first_level: string
  grade2: string
  score: number
  init_score: number
  max_score: number
  ai_diagnosis: string
  intention_dept: string
  intention_project: string
  intention_campus: string
  intention_grade: string
  intention_subject: string
  contact_point: string
  activity_name: string
  channel_l1: string
  channel_l2: string
  channel_l3: string
  business_channel: string
  system_source: string
  create_method: string
  lead_code: string
  owner_id: string
  owner_name: string
  owner_dept: string
  assign_at: string
  follow_count: number
  last_follow_at: string
  last_follow_content: string
  last_follow_by: string
  next_follow_at: string
  promise_visit_at: string
  visit_converted: number
  archive_at: string
  archive_rule: string
  planned_recycle_at: string
  recycle_reason: string
  recycle_at: string
  deal_at: string
  remark: string
  creator: string
  created_at: string
  pushed_91: number
  push_91_result: string
  resource_no_91: string
  push_91_at: string
  follow_91: string
  advancer_91: string
  pushed_hedgehog: number
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
