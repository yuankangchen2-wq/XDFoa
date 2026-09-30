import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/customers' },
  { path: '/customers', component: () => import('../views/Customer.vue') },
  { path: '/products', component: () => import('../views/Product.vue') },
  { path: '/suppliers', component: () => import('../views/Supplier.vue') },
  { path: '/organizations', component: () => import('../views/Organization.vue') }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
