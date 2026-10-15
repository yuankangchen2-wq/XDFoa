import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/home' },
  { path: '/home', component: () => import('../views/HomePage.vue'), meta: { title: '首页' } },
  { path: '/lead-import', component: () => import('../views/LeadImport.vue'), meta: { title: '线索导入' } },
  { path: '/activity', component: () => import('../views/ActivityList.vue'), meta: { title: '活动管理' } },
  { path: '/handover-leads', component: () => import('../views/LeadList.vue'), meta: { title: '待交接线索' } },
  { path: '/all-leads', component: () => import('../views/LeadList.vue'), meta: { title: '全部线索' } },
  { path: '/exception-leads', component: () => import('../views/ExceptionLeads.vue'), meta: { title: '异常线索' } },
  { path: '/data-record', component: () => import('../views/DataRecord.vue'), meta: { title: '数据记录' } },
  { path: '/coming-soon', component: () => import('../views/ComingSoon.vue'), meta: { title: '功能建设中' } },

  // 既有 CRM 页面（保留）
  { path: '/customers', component: () => import('../views/CustomerList.vue'), meta: { title: '客户列表' } },
  { path: '/customers/:id', component: () => import('../views/CustomerDetail.vue'), meta: { title: '客户详情' } },
  { path: '/opportunities', component: () => import('../views/OpportunityList.vue'), meta: { title: '商机管理' } },
  { path: '/sales/orders', component: () => import('../views/SalesOrderList.vue'), meta: { title: '销售订单' } },
  { path: '/sales/contracts', component: () => import('../views/ContractList.vue'), meta: { title: '合同管理' } },
  { path: '/tickets', component: () => import('../views/TicketList.vue'), meta: { title: '工单列表' } },
  { path: '/tickets/:id', component: () => import('../views/TicketDetail.vue'), meta: { title: '工单详情' } },
  { path: '/analytics/dashboard', component: () => import('../views/AnalyticsDashboard.vue'), meta: { title: '分析看板' } },
  { path: '/analytics/sales', component: () => import('../views/SalesReport.vue'), meta: { title: '销售报表' } },
  { path: '/analytics/funnel', component: () => import('../views/FunnelReport.vue'), meta: { title: '转化漏斗' } },
  { path: '/analytics/customers', component: () => import('../views/CustomerAnalytics.vue'), meta: { title: '客户分析' } },
  { path: '/:pathMatch(.*)*', redirect: '/home' }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
