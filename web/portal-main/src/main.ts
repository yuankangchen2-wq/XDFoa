import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import PortalHome from './views/PortalHome.vue'

const Empty = { template: '<div></div>' }

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: PortalHome },
    { path: '/portal/:pathMatch(.*)*', component: Empty },
    { path: '/iam/:pathMatch(.*)*', component: Empty },
    { path: '/mdm/:pathMatch(.*)*', component: Empty },
    { path: '/erp/:pathMatch(.*)*', component: Empty },
    { path: '/crm/:pathMatch(.*)*', component: Empty },
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ]
})

const app = createApp(App)
for (const [key, comp] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, comp as any)
}
app.use(router)
app.use(ElementPlus)
app.mount('#app')
