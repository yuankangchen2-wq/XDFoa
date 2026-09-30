import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/customers' },
  { path: '/customers', component: () => import('../views/CustomerList.vue') },
  { path: '/customers/:id', component: () => import('../views/CustomerDetail.vue') },
  { path: '/opportunities', component: () => import('../views/OpportunityList.vue') },
  { path: '/sales/orders', component: () => import('../views/SalesOrderList.vue') },
  { path: '/sales/contracts', component: () => import('../views/ContractList.vue') },
  { path: '/tickets', component: () => import('../views/TicketList.vue') },
  { path: '/tickets/:id', component: () => import('../views/TicketDetail.vue') },
  { path: '/analytics/dashboard', component: () => import('../views/AnalyticsDashboard.vue') },
  { path: '/analytics/sales', component: () => import('../views/SalesReport.vue') },
  { path: '/analytics/funnel', component: () => import('../views/FunnelReport.vue') },
  { path: '/analytics/customers', component: () => import('../views/CustomerAnalytics.vue') }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
