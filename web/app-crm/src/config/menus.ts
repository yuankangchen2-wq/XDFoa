export interface SubMenu {
  title: string
  path: string
  query?: Record<string, string>
}

export interface TopMenu {
  key: string
  title: string
  icon: string
  children?: SubMenu[]
  path?: string
  implemented?: boolean
}

// 全局 14 个一级菜单（截图逆向）
export const menus: TopMenu[] = [
  {
    key: 'home', title: '首页', icon: 'HomeFilled', path: '/home'
  },
  {
    key: 'lead-create', title: '线索创建', icon: 'Promotion',
    children: [
      { title: '线索导入', path: '/lead-import' },
      { title: '活动管理', path: '/activity' },
      { title: '待交接线索', path: '/handover-leads' },
      { title: '全部线索', path: '/all-leads' },
      { title: '异常线索', path: '/exception-leads' },
      { title: '数据记录', path: '/data-record' }
    ]
  },
  {
    key: 'lead-mgr', title: '线索管理', icon: 'List',
    children: [
      { title: '员工工作台(线索)', path: '/lead-workbench', query: { mode: 'employee' } },
      { title: '主管工作台(线索)', path: '/lead-workbench', query: { mode: 'supervisor' } }
    ]
  },
  {
    key: 'opp-mgr', title: '商机管理', icon: 'TrendCharts',
    children: [
      { title: '员工工作台(商机)', path: '/opp-workbench', query: { mode: 'employee' } },
      { title: '主管工作台(商机)', path: '/opp-workbench', query: { mode: 'supervisor' } },
      { title: '商机盘点', path: '/opp-inventory' },
      { title: '业绩申诉', path: '/appeals' },
      { title: '审核中心', path: '/approvals' },
      { title: '商机备单', path: '/opp-standby' },
      { title: '成单日志', path: '/deal-logs' },
      { title: '退调转/订单', path: '/ticket-orders' }
    ]
  },
  { key: 'wechat', title: '微信管理', icon: 'ChatDotRound' },
  { key: 'visit', title: '到店管理', icon: 'Shop' },
  { key: 'student', title: '学员服务', icon: 'User' },
  { key: 'daily', title: '日常应用', icon: 'Grid' },
  { key: 'report', title: '报表管理', icon: 'DataAnalysis', path: '/analytics/funnel' },
  { key: 'qc', title: '智能质检', icon: 'Microphone' },
  { key: 'nonfile', title: '非档案任务', icon: 'FolderOpened' },
  { key: 'affair', title: '事务管理', icon: 'Tickets' },
  { key: 'refund', title: '退换管理', icon: 'RefreshLeft' },
  { key: 'admin', title: '管理工具', icon: 'Setting' }
]

// 已落地页面的路由白名单（其余菜单进入占位页）
export function resolveMenu(menu: TopMenu, sub?: SubMenu) {
  if (sub) {
    return { path: sub.path, query: sub.query, title: sub.title }
  }
  if (menu.path) {
    return { path: menu.path, title: menu.title }
  }
  return { path: '/coming-soon', query: { name: menu.title }, title: menu.title }
}
