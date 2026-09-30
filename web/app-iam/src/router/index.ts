import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/users' },
  { path: '/login', component: () => import('../views/Login.vue') },
  { path: '/users', component: () => import('../views/User.vue') },
  { path: '/roles', component: () => import('../views/Role.vue') },
  { path: '/permissions', component: () => import('../views/Permission.vue') },
  { path: '/audit', component: () => import('../views/AuditLog.vue') }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
