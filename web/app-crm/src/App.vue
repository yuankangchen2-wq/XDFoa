<template>
  <div class="scrm-layout">
    <!-- 左侧导航 -->
    <aside class="sidebar" :class="{ collapsed }">
      <div class="logo">
        <span class="logo-mark">S</span>
        <span v-show="!collapsed" class="logo-text">新东方 SCRM</span>
      </div>
      <el-scrollbar class="menu-scroll">
        <div
          v-for="m in menus"
          :key="m.key"
          class="menu-item-wrap"
        >
          <div
            class="menu-item"
            :class="{ active: isMenuActive(m) }"
            @click="clickTopMenu(m)"
          >
            <el-icon class="menu-icon"><component :is="iconComp(m.icon)" /></el-icon>
            <span v-show="!collapsed" class="menu-title">{{ m.title }}</span>
            <el-icon
              v-show="!collapsed && m.children?.length"
              class="menu-arrow"
              :class="{ open: opened[m.key] || isMenuActive(m) }"
            >
              <ArrowDown />
            </el-icon>
          </div>
          <div
            v-show="!collapsed && (opened[m.key] || isMenuActive(m)) && m.children?.length"
            class="sub-list"
          >
            <div
              v-for="s in m.children"
              :key="s.path"
              class="sub-item"
              :class="{ active: route.path === s.path }"
              @click="clickSubMenu(m, s)"
            >
              <span class="sub-dot" />{{ s.title }}
            </div>
          </div>
        </div>
      </el-scrollbar>
    </aside>

    <!-- 右侧主区域 -->
    <div class="main-area">
      <!-- 顶部栏 -->
      <header class="topbar">
        <div class="topbar-left">
          <el-icon class="collapse-btn" @click="collapsed = !collapsed">
            <Fold v-if="!collapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title && route.path !== '/home'">
              {{ route.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="topbar-right">
          <span class="top-link"><el-icon><QuestionFilled /></el-icon>帮助文档</span>
          <span class="top-link" title="云座席双呼" @click="openCallBar()"><el-icon><Phone /></el-icon>云座席</span>
          <span class="top-link"><el-icon><ChatDotRound /></el-icon>企业微信</span>
          <GlobalSearch />
          <NotificationCenter />
          <span class="top-link city">
            <el-icon><LocationInformation /></el-icon>
            {{ me.city }}·{{ me.store }}
          </span>
          <el-dropdown>
            <span class="user-chip">
              <el-avatar :size="26" class="user-avatar">{{ me.name.slice(0, 1) }}</el-avatar>
              {{ me.name }}<el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人中心</el-dropdown-item>
                <el-dropdown-item>切换角色：主管视角</el-dropdown-item>
                <el-dropdown-item divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- 多标签页 -->
      <div class="tabbar">
        <div
          v-for="t in tabStore.tabs"
          :key="t.fullPath"
          class="work-tab"
          :class="{ active: t.fullPath === route.fullPath }"
          @click="goTab(t)"
        >
          {{ t.title }}
          <el-icon
            v-if="t.closable"
            class="tab-close"
            @click.stop="closeWorkTab(t)"
          ><Close /></el-icon>
        </div>
      </div>

      <!-- 内容区 -->
      <main class="content">
        <router-view />
      </main>
    </div>

    <!-- 云座席双呼工作条（全局唤起） -->
    <CallBar />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowDown, Close, Fold, Expand, QuestionFilled, Phone, ChatDotRound,
  LocationInformation,
  HomeFilled, Promotion, List, TrendCharts, Shop, User, Grid,
  DataAnalysis, Microphone, FolderOpened, Tickets, RefreshLeft, Setting
} from '@element-plus/icons-vue'
import { menus, resolveMenu, type TopMenu, type SubMenu } from './config/menus'
import { tabStore, openTab, closeTab } from './store/tabs'
import { scrmApi, type MeInfo } from './api/scrm'
import { openCallBar } from './store/call'
import GlobalSearch from './components/GlobalSearch.vue'
import CallBar from './components/CallBar.vue'
import NotificationCenter from './components/NotificationCenter.vue'

const route = useRoute()
const router = useRouter()
const collapsed = ref(false)
const opened = reactive<Record<string, boolean>>({})
const me = ref<MeInfo>({ name: '陈元素', title: '', city: '武汉', store: '江汉校区', unread: 0, roles: [] })

const iconMap: Record<string, any> = {
  HomeFilled, Promotion, List, TrendCharts, ChatDotRound, Shop, User, Grid,
  DataAnalysis, Microphone, FolderOpened, Tickets, RefreshLeft, Setting
}
const iconComp = (name: string) => iconMap[name] || HomeFilled

onMounted(async () => {
  try {
    me.value = await scrmApi.me()
  } catch { /* 离线兜底 */ }
  syncRouteTab()
})

function isMenuActive(m: TopMenu) {
  if (m.children) {
    return m.children.some(s => route.path.startsWith(s.path))
  }
  return route.path === m.path
}

function clickTopMenu(m: TopMenu) {
  if (m.children?.length) {
    if (collapsed.value) {
      collapsed.value = false
    } else {
      opened[m.key] = !(opened[m.key] || isMenuActive(m))
    }
    return
  }
  const r = resolveMenu(m)
  router.push({ path: r.path, query: r.query as any })
}

function clickSubMenu(_m: TopMenu, s: SubMenu) {
  router.push({ path: s.path, query: s.query as any })
}

function goTab(t: { fullPath: string }) {
  if (t.fullPath !== route.fullPath) router.push(t.fullPath)
}

function closeWorkTab(t: { fullPath: string }) {
  const redirect = closeTab(t.fullPath)
  if (redirect) router.push(redirect)
}

function syncRouteTab() {
  const title = (route.meta.title as string) || '未命名'
  openTab({
    title,
    path: route.path,
    fullPath: route.fullPath,
    closable: route.path !== '/home'
  })
}

// 路由变化时自动开/切标签
watch(() => route.fullPath, () => syncRouteTab())
</script>

<style scoped>
.scrm-layout {
  display: flex;
  height: 100vh;
  background: #f0f2f5;
  --scrm-green: #00a86b;
}

.sidebar {
  width: 220px;
  background: #fff;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
  transition: width .2s;
  flex-shrink: 0;
}
.sidebar.collapsed { width: 60px; }

.logo {
  height: 56px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 14px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}
.logo-mark {
  width: 30px; height: 30px; border-radius: 8px;
  background: var(--scrm-green);
  color: #fff; font-weight: 700; font-size: 18px;
  display: flex; align-items: center; justify-content: center;
}
.logo-text { font-weight: 700; font-size: 15px; color: #1f2d3d; white-space: nowrap; }

.menu-scroll { flex: 1; }
.menu-item {
  height: 44px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  cursor: pointer;
  color: #475669;
  font-size: 14px;
  white-space: nowrap;
  position: relative;
}
.menu-item:hover { background: #f3faf7; color: var(--scrm-green); }
.menu-item.active { color: var(--scrm-green); background: #eaf8f2; font-weight: 600; }
.menu-item.active::before {
  content: '';
  position: absolute; left: 0; top: 8px; bottom: 8px;
  width: 3px; border-radius: 0 3px 3px 0; background: var(--scrm-green);
}
.menu-icon { font-size: 17px; flex-shrink: 0; }
.menu-title { flex: 1; }
.menu-arrow { font-size: 12px; transition: transform .2s; }
.menu-arrow.open { transform: rotate(180deg); }

.sub-list { background: #fafbfc; }
.sub-item {
  height: 38px; line-height: 38px; padding-left: 46px;
  font-size: 13px; color: #606f82; cursor: pointer; white-space: nowrap;
}
.sub-item:hover { color: var(--scrm-green); }
.sub-item.active { color: var(--scrm-green); font-weight: 600; }
.sub-dot {
  display: inline-block; width: 5px; height: 5px; border-radius: 50%;
  background: #c0ccda; margin-right: 8px; vertical-align: middle;
}
.sub-item.active .sub-dot { background: var(--scrm-green); }

.main-area { flex: 1; display: flex; flex-direction: column; min-width: 0; }

.topbar {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  flex-shrink: 0;
}
.topbar-left { display: flex; align-items: center; gap: 14px; }
.collapse-btn { font-size: 19px; cursor: pointer; color: #5e6d82; }
.topbar-right { display: flex; align-items: center; gap: 16px; }
.top-link {
  display: inline-flex; align-items: center; gap: 4px;
  font-size: 13px; color: #5e6d82; cursor: pointer; white-space: nowrap;
}
.top-link:hover { color: var(--scrm-green); }
.global-search { width: 230px; }
.msg-badge { display: inline-flex; align-items: center; }
.city { max-width: 200px; overflow: hidden; text-overflow: ellipsis; }
.user-chip {
  display: inline-flex; align-items: center; gap: 6px;
  cursor: pointer; font-size: 13px; color: #1f2d3d;
}
.user-avatar { background: var(--scrm-green); }

.tabbar {
  height: 40px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  padding: 0 10px;
  gap: 6px;
  overflow-x: auto;
  flex-shrink: 0;
}
.work-tab {
  height: 28px;
  padding: 0 10px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 13px;
  color: #475669;
  cursor: pointer;
  white-space: nowrap;
  background: #fff;
}
.work-tab.active { color: var(--scrm-green); border-color: var(--scrm-green); background: #f0fbf7; }
.tab-close { font-size: 12px; border-radius: 50%; }
.tab-close:hover { background: #e4e7ed; }

.content { flex: 1; overflow: auto; padding: 12px; }
</style>
