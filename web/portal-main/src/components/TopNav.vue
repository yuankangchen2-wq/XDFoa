<template>
  <header class="top-nav">
    <div class="nav-inner">
      <!-- 左侧 Logo -->
      <div class="brand" @click="$router.push('/')">
        <div class="brand-logo">
          <el-icon :size="22"><Promotion /></el-icon>
        </div>
        <div class="brand-text">
          <span class="brand-name">新东方</span>
          <span class="brand-divider">｜</span>
          <span class="brand-portal">集团统一门户</span>
        </div>
      </div>

      <!-- 右侧操作区 -->
      <div class="nav-actions">
        <el-autocomplete
          v-model="empKeyword"
          :fetch-suggestions="querySearch"
          placeholder="员工 / 邮箱"
          clearable
          class="emp-search"
          @select="handleSelect"
        >
          <template #default="{ item }">
            <div class="emp-item">
              <el-avatar :size="28" class="emp-avatar">{{ (item.real_name || item.username || '?').charAt(0) }}</el-avatar>
              <div class="emp-meta">
                <div class="emp-name">{{ item.real_name || item.username }}</div>
                <div class="emp-email">{{ item.email || item.phone || '—' }}</div>
              </div>
            </div>
          </template>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-autocomplete>

        <el-button type="primary" class="download-btn" @click="goDownload">
          <el-icon style="margin-right:4px"><Download /></el-icon>
          软件下载中心
        </el-button>

        <el-tooltip content="账号与安全" placement="bottom">
          <div class="icon-btn" @click="goIam">
            <el-icon :size="19"><Lock /></el-icon>
          </div>
        </el-tooltip>

        <el-dropdown @command="handleCommand">
          <div class="user-box">
            <el-avatar :size="32" class="user-avatar">{{ avatarText }}</el-avatar>
            <span class="user-name">{{ me?.username || '加载中...' }}</span>
            <el-icon><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="iam">
                <el-icon><Lock /></el-icon>权限中心
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { bffApi, type MeInfo, type Employee } from '../api/bff'

const router = useRouter()
const me = ref<MeInfo | null>(null)
const empKeyword = ref('')

const avatarText = computed(() => (me.value?.username || '陈').charAt(0))

const loadMe = async () => {
  try { me.value = await bffApi.me() } catch { /* 忽略 */ }
}

const querySearch = async (queryString: string, cb: (items: Employee[]) => void) => {
  try {
    const res = await bffApi.searchEmployees(queryString || '')
    cb(res.items || [])
  } catch {
    cb([])
  }
}

const handleSelect = (item: Employee) => {
  ElMessage.info(`员工：${item.real_name || item.username}（${item.email || '未绑定邮箱'}）`)
}

const goDownload = () => ElMessage.info('软件下载中心即将上线，敬请期待')
const goIam = () => router.push('/iam')

const handleCommand = (cmd: string) => {
  if (cmd === 'iam') router.push('/iam')
  if (cmd === 'logout') ElMessage.info('演示环境暂不支持退出')
}

onMounted(loadMe)
</script>

<style scoped>
.top-nav {
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 100;
}
.nav-inner {
  max-width: 1440px;
  margin: 0 auto;
  height: 100%;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.brand { display: flex; align-items: center; cursor: pointer; user-select: none; }
.brand-logo {
  width: 36px; height: 36px; border-radius: 9px;
  background: linear-gradient(135deg, #00c987, #00a06b);
  color: #fff; display: flex; align-items: center; justify-content: center;
  margin-right: 10px;
}
.brand-text { font-size: 18px; font-weight: 700; letter-spacing: .5px; }
.brand-name { color: #00b578; }
.brand-divider { color: #cfd6de; margin: 0 4px; font-weight: 400; }
.brand-portal { color: #1f2937; }

.nav-actions { display: flex; align-items: center; gap: 14px; }
.emp-search { width: 210px; }
.emp-item { display: flex; align-items: center; gap: 10px; padding: 2px 0; }
.emp-avatar { background: #00b578; font-size: 13px; }
.emp-meta { line-height: 1.3; }
.emp-name { font-size: 13px; color: #1f2937; font-weight: 600; }
.emp-email { font-size: 12px; color: #909399; }

.download-btn { font-weight: 600; }
.icon-btn {
  width: 34px; height: 34px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  color: #5f6b7a; cursor: pointer; transition: background .2s;
}
.icon-btn:hover { background: #ecfaf6; color: #00b578; }
.user-box { display: flex; align-items: center; gap: 8px; cursor: pointer; }
.user-avatar { background: #00b578; font-weight: 600; }
.user-name { font-size: 14px; color: #1f2937; }
</style>
