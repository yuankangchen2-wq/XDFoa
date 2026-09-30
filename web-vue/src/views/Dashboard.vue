<template>
  <div class="dashboard">
    <!-- 欢迎横幅 -->
    <el-card class="welcome-card" shadow="never">
      <div class="welcome-content">
        <div>
          <h2>你好，{{ displayName }} 👋</h2>
          <p>今天是 {{ today }}，祝你工作顺利！</p>
        </div>
        <div class="welcome-stats">
          <div class="stat-item">
            <span class="stat-num">{{ todoCount }}</span>
            <span class="stat-label">待办事项</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">{{ appCount }}</span>
            <span class="stat-label">可用应用</span>
          </div>
          <div class="stat-item">
            <span class="stat-num">{{ msgCount }}</span>
            <span class="stat-label">未读消息</span>
          </div>
        </div>
      </div>
    </el-card>

    <el-row :gutter="20">
      <!-- 常用应用 -->
      <el-col :span="16">
        <el-card shadow="never" class="section-card">
          <template #header>
            <div class="card-header">
              <span>常用应用</span>
              <el-button text type="primary" @click="$router.push('/apps')">全部应用 →</el-button>
            </div>
          </template>
          <el-row :gutter="16">
            <el-col :span="6" v-for="app in apps" :key="app.name">
              <div class="app-card" @click="openApp(app)">
                <div class="app-icon" :style="{ background: app.color }">
                  {{ app.name?.charAt(0) || 'A' }}
                </div>
                <div class="app-name">{{ app.displayName || app.name }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <!-- 待办事项 -->
      <el-col :span="8">
        <el-card shadow="never" class="section-card">
          <template #header>
            <div class="card-header">
              <span>待办事项</span>
              <el-badge :value="todoCount" class="badge">
                <el-button text type="primary" @click="$router.push('/todos')">处理 →</el-button>
              </el-badge>
            </div>
          </template>
          <el-empty v-if="todos.length === 0" description="暂无待办" :image-size="60" />
          <div v-else>
            <div v-for="todo in todos" :key="todo.id" class="todo-item">
              <el-tag :type="todo.priority === 'high' ? 'danger' : todo.priority === 'medium' ? 'warning' : 'info'" size="small">
                {{ todo.priority === 'high' ? '高' : todo.priority === 'medium' ? '中' : '低' }}
              </el-tag>
              <span class="todo-title">{{ todo.title }}</span>
              <span class="todo-time">{{ todo.time }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 公告 -->
      <el-col :span="16">
        <el-card shadow="never" class="section-card">
          <template #header>
            <div class="card-header">
              <span>公告通知</span>
              <el-button text type="primary">更多 →</el-button>
            </div>
          </template>
          <div v-for="item in announcements" :key="item.id" class="announcement-item">
            <el-tag type="danger" size="small" v-if="item.top">置顶</el-tag>
            <span class="announcement-title">{{ item.title }}</span>
            <span class="announcement-date">{{ item.date }}</span>
          </div>
        </el-card>
      </el-col>

      <!-- 快捷入口 -->
      <el-col :span="8">
        <el-card shadow="never" class="section-card">
          <template #header>
            <div class="card-header">
              <span>快捷入口</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button class="quick-btn" @click="$router.push('/apps')">
              <el-icon><Grid /></el-icon>
              <span>应用中心</span>
            </el-button>
            <el-button class="quick-btn" @click="$router.push('/todos')">
              <el-icon><List /></el-icon>
              <span>我的待办</span>
            </el-button>
            <el-button class="quick-btn" @click="$router.push('/messages')">
              <el-icon><Bell /></el-icon>
              <span>消息中心</span>
            </el-button>
            <el-button class="quick-btn" @click="$router.push('/profile')">
              <el-icon><User /></el-icon>
              <span>个人中心</span>
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { applicationApi } from '@/api'

const userStore = useUserStore()

const displayName = computed(() => userStore.userInfo?.displayName || userStore.userInfo?.name || '用户')
const today = new Date().toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' })

const apps = ref<any[]>([])
const todos = ref<any[]>([
  { id: 1, title: '审批：XXX的请假申请', priority: 'high', time: '2小时前' },
  { id: 2, title: '审批：采购订单 #20240926', priority: 'medium', time: '昨天' },
  { id: 3, title: '会议：周例会', priority: 'low', time: '明天 10:00' }
])
const announcements = ref<any[]>([
  { id: 1, title: '关于国庆假期安排的通知', date: '2024-09-25', top: true },
  { id: 2, title: '系统升级维护公告', date: '2024-09-24', top: false },
  { id: 3, title: '新员工入职培训通知', date: '2024-09-23', top: false }
])

const todoCount = computed(() => todos.value.length)
const appCount = computed(() => apps.value.length)
const msgCount = ref(3)

const colors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#9C27B0', '#00BCD4', '#FF5722']

onMounted(async () => {
  try {
    const res: any = await applicationApi.getOrganizationApplications('built-in')
    apps.value = (res.data || []).slice(0, 8).map((app: any, i: number) => ({
      ...app,
      color: colors[i % colors.length]
    }))
  } catch (e) {
    // 后端未启动时使用示例数据
    apps.value = [
      { name: 'CRM', displayName: '客户管理系统', color: '#409EFF' },
      { name: 'ERP', displayName: '企业资源计划', color: '#67C23A' },
      { name: 'OA', displayName: '办公自动化', color: '#E6A23C' },
      { name: 'HR', displayName: '人力资源', color: '#F56C6C' },
      { name: 'BI', displayName: '商业智能', color: '#9C27B0' },
      { name: '财务', displayName: '财务系统', color: '#00BCD4' },
      { name: '项目', displayName: '项目管理', color: '#FF5722' },
      { name: '邮箱', displayName: '企业邮箱', color: '#909399' }
    ]
  }
})

function openApp(app: any) {
  if (app.homeUrl) {
    window.open(app.homeUrl, '_blank')
  }
}
</script>

<style scoped>
.dashboard {
  padding: 0;
}
.welcome-card {
  margin-bottom: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}
.welcome-card :deep(.el-card__body) {
  padding: 24px;
}
.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #fff;
}
.welcome-content h2 {
  margin: 0 0 8px;
  font-size: 22px;
}
.welcome-content p {
  margin: 0;
  opacity: 0.9;
}
.welcome-stats {
  display: flex;
  gap: 40px;
}
.stat-item {
  text-align: center;
}
.stat-num {
  display: block;
  font-size: 28px;
  font-weight: bold;
}
.stat-label {
  font-size: 13px;
  opacity: 0.9;
}
.section-card {
  height: 100%;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}
.app-card {
  text-align: center;
  padding: 16px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 12px;
}
.app-card:hover {
  background: #f5f7fa;
  transform: translateY(-2px);
}
.app-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 8px;
  border-radius: 12px;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
}
.app-name {
  font-size: 13px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.todo-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}
.todo-title {
  flex: 1;
  font-size: 14px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.todo-time {
  font-size: 12px;
  color: #909399;
}
.announcement-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}
.announcement-title {
  flex: 1;
  font-size: 14px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.announcement-date {
  font-size: 12px;
  color: #909399;
}
.quick-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.quick-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px;
  height: auto;
}
.quick-btn .el-icon {
  font-size: 22px;
}
</style>
