import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/orders' },
  { path: '/orders', component: () => import('../views/OrderList.vue') },
  { path: '/orders/:id', component: () => import('../views/OrderDetail.vue') }
]

export default createRouter({
  history: createWebHashHistory(),
  routes
})
