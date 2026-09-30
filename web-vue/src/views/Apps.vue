<template>
  <div class="apps-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>应用中心</span>
          <el-input v-model="keyword" placeholder="搜索应用" style="width: 240px" prefix-icon="Search" clearable />
        </div>
      </template>

      <el-row :gutter="16">
        <el-col :span="6" v-for="app in filteredApps" :key="app.name">
          <div class="app-card" @click="openApp(app)">
            <div class="app-icon" :style="{ background: app.color }">
              {{ app.displayName?.charAt(0) || app.name?.charAt(0) || 'A' }}
            </div>
            <div class="app-info">
              <div class="app-name">{{ app.displayName || app.name }}</div>
              <div class="app-desc">{{ app.description || '暂无描述' }}</div>
            </div>
          </div>
        </el-col>
      </el-row>

      <el-empty v-if="filteredApps.length === 0" description="未找到应用" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { applicationApi } from '@/api'

const keyword = ref('')
const apps = ref<any[]>([])

const colors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#9C27B0', '#00BCD4', '#FF5722']

const filteredApps = computed(() => {
  if (!keyword.value) return apps.value
  return apps.value.filter(a =>
    (a.displayName || a.name || '').toLowerCase().includes(keyword.value.toLowerCase())
  )
})

onMounted(async () => {
  try {
    const res: any = await applicationApi.getOrganizationApplications('built-in')
    apps.value = (res.data || []).map((app: any, i: number) => ({
      ...app,
      color: colors[i % colors.length]
    }))
  } catch (e) {
    apps.value = [
      { name: 'CRM', displayName: '客户管理系统', description: '客户关系管理', color: '#409EFF' },
      { name: 'ERP', displayName: '企业资源计划', description: '企业资源规划系统', color: '#67C23A' },
      { name: 'OA', displayName: '办公自动化', description: '日常办公协作', color: '#E6A23C' },
      { name: 'HR', displayName: '人力资源', description: '人事管理系统', color: '#F56C6C' },
      { name: 'BI', displayName: '商业智能', description: '数据分析报表', color: '#9C27B0' },
      { name: '财务', displayName: '财务系统', description: '财务管理平台', color: '#00BCD4' }
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
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.app-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 16px;
}
.app-card:hover {
  border-color: #409EFF;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.app-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.app-info {
  flex: 1;
  overflow: hidden;
}
.app-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}
.app-desc {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
