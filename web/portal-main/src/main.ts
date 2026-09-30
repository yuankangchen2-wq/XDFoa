import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/portal' },
    { path: '/portal/:pathMatch(.*)*', component: { template: '<div></div>' } },
    { path: '/iam/:pathMatch(.*)*', component: { template: '<div></div>' } },
    { path: '/mdm/:pathMatch(.*)*', component: { template: '<div></div>' } },
    { path: '/erp/:pathMatch(.*)*', component: { template: '<div></div>' } },
    { path: '/crm/:pathMatch(.*)*', component: { template: '<div></div>' } }
  ]
})

const app = createApp(App)
for (const [key, comp] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, comp as any)
}
app.use(router)
app.use(ElementPlus)
app.mount('#app')
