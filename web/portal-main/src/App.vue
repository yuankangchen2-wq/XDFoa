<template>
  <el-container style="height:100vh">
    <el-aside width="220px" style="background:#001529">
      <div style="color:#fff;font-size:18px;font-weight:bold;padding:20px;text-align:center;border-bottom:1px solid #1f2d3d">
        集团 OA 门户
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#001529"
        text-color="#fff"
        active-text-color="#409eff"
        @select="handleMenuSelect"
      >
        <el-menu-item v-for="m in menus" :key="m.path" :index="m.path">
          <el-icon><component :is="m.icon" /></el-icon>
          <span>{{ m.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header style="background:#fff;border-bottom:1px solid #eee;display:flex;justify-content:space-between;align-items:center">
        <div style="font-size:16px;font-weight:bold">{{ currentTitle }}</div>
        <div style="display:flex;align-items:center;gap:16px">
          <el-tag size="small" type="success" @click="showHealth = true" style="cursor:pointer">系统状态</el-tag>
          <el-dropdown>
            <span style="cursor:pointer">
              <el-avatar :size="30" style="background:#409eff">{{ user.username.charAt(0) }}</el-avatar>
              <span style="margin-left:8px">{{ user.username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人中心</el-dropdown-item>
                <el-dropdown-item divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main style="padding:0;background:#f5f5f5;position:relative">
        <div v-if="!activeApp" style="padding:24px">
          <el-card>
            <template #header><b>欢迎使用集团 OA 门户</b></template>
            <p>请从左侧菜单选择功能模块。</p>
            <el-row :gutter="16" style="margin-top:16px">
              <el-col :span="8" v-for="m in menus" :key="m.path">
                <el-card shadow="hover" style="text-align:center;cursor:pointer" @click="handleMenuSelect(m.path)">
                  <el-icon :size="32" color="#409eff"><component :is="m.icon" /></el-icon>
                  <div style="margin-top:8px;font-weight:bold">{{ m.title }}</div>
                </el-card>
              </el-col>
            </el-row>
          </el-card>
        </div>
        <div id="subapp-container"></div>
      </el-main>
    </el-container>

    <!-- 系统健康状态弹窗 -->
    <el-dialog v-model="showHealth" title="系统服务状态" width="600px">
      <el-table :data="healthList" border>
        <el-table-column prop="service" label="服务" />
        <el-table-column prop="url" label="地址" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'up' ? 'success' : 'danger'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="latency" label="耗时(ms)" width="100" />
      </el-table>
    </el-dialog>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { startQiankun, microApps } from './qiankun'
import axios from 'axios'
import {
  HomeFilled, Lock, Files, Box, UserFilled
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const menus = [
  { key: 'portal', title: '工作台', path: '/portal', icon: 'HomeFilled' },
  { key: 'iam', title: '权限中心', path: '/iam', icon: 'Lock' },
  { key: 'mdm', title: '主数据', path: '/mdm', icon: 'Files' },
  { key: 'erp', title: 'ERP', path: '/erp', icon: 'Box' },
  { key: 'crm', title: 'CRM', path: '/crm', icon: 'UserFilled' }
]

const user = ref({ username: '管理员' })
const showHealth = ref(false)
const healthList = ref<any[]>([])

const activeMenu = computed(() => {
  const p = route.path
  for (const m of menus) {
    if (p.startsWith(m.path)) return m.path
  }
  return ''
})

const currentTitle = computed(() => {
  const m = menus.find(m => activeMenu.value === m.path)
  return m ? m.title : '集团 OA 门户'
})

const activeApp = computed(() => {
  return microApps.some(a => route.path.startsWith(a.activeRule))
})

const handleMenuSelect = (path: string) => {
  router.push(path)
}

const loadHealth = async () => {
  try {
    const res = await axios.get('/api/bff/health')
    healthList.value = res.data.items || []
  } catch { healthList.value = [] }
}

onMounted(() => {
  startQiankun()
  loadHealth()
})
</script>

<style>
html, body, #app { margin: 0; padding: 0; height: 100%; }
</style>
