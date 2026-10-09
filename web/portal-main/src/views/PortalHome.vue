<template>
  <div class="portal-home">
    <div class="home-grid">
      <!-- ============ 左侧主区 ============ -->
      <div class="main-col">
        <!-- 1. 问候卡片 -->
        <section class="card greeting-card">
          <div class="greet-left">
            <div class="sun-icon">
              <el-icon :size="30"><Sunny /></el-icon>
            </div>
            <div>
              <div class="greet-title">
                Hi~ {{ me?.username || '陈元康' }} <span class="greet-suffix">{{ me?.title || '老师' }}</span>
              </div>
              <div class="greet-sub">
                今天是您在新东方相伴的第 <b class="days">{{ me?.join_days ?? 806 }}</b> 天，愿您工作顺利，万事顺心！
              </div>
            </div>
          </div>
          <div class="greet-date">
            <div class="gd-main">{{ todayStr }}</div>
            <div class="gd-week">{{ weekStr }}</div>
          </div>
        </section>

        <!-- 2. 运营 Banner -->
        <section v-if="!bannerCollapsed" class="banner-wrap">
          <el-carousel height="170px" arrow="hover" indicator-position="outside">
            <el-carousel-item v-for="(b, i) in banners" :key="i">
              <div class="banner-slide" :style="{ backgroundImage: `url(${b.img})` }">
                <div class="banner-mask"></div>
                <div class="banner-text">
                  <div class="bt-title">{{ b.title }}</div>
                  <div class="bt-sub">{{ b.sub }}</div>
                </div>
              </div>
            </el-carousel-item>
          </el-carousel>
          <div class="banner-collapse" @click="bannerCollapsed = true">
            收起 <el-icon><ArrowUpBold /></el-icon>
          </div>
        </section>
        <div v-else class="banner-collapsed-bar" @click="bannerCollapsed = false">
          <span class="bc-dot"></span> 运营活动
          <span class="bc-expand">展开 <el-icon><ArrowDownBold /></el-icon></span>
        </div>

        <!-- 3. 办事大厅 -->
        <section class="card hall-card">
          <div class="hall-header">
            <el-radio-group v-model="hallTab" size="default">
              <el-radio-button value="hall">办事大厅</el-radio-button>
              <el-radio-button value="guide">系统指南</el-radio-button>
            </el-radio-group>
            <el-input
              v-model="appKeyword"
              placeholder="搜索应用名称"
              clearable
              class="hall-search"
              :prefix-icon="Search"
            />
          </div>

          <!-- 办事大厅 -->
          <template v-if="hallTab === 'hall'">
            <div class="cat-tabs">
              <span
                v-for="c in categories"
                :key="c"
                class="cat-pill"
                :class="{ active: activeCategory === c }"
                @click="activeCategory = c"
              >{{ c }}</span>
            </div>
            <div class="app-grid">
              <div
                v-for="app in filteredApps"
                :key="app.name"
                class="app-item"
                @click="openApp(app)"
              >
                <div class="app-icon" :style="{ background: app.color + '18', color: app.color }">
                  <el-icon :size="24"><component :is="app.icon" /></el-icon>
                </div>
                <div class="app-name">{{ app.name }}</div>
              </div>
              <el-empty v-if="filteredApps.length === 0" description="未找到相关应用" :image-size="70" />
            </div>
          </template>

          <!-- 系统指南 -->
          <template v-else>
            <div class="guide-list">
              <div v-for="g in guides" :key="g" class="guide-item" @click="ElMessage.info(`《${g}》即将开放查看`)">
                <el-icon class="guide-icon"><Document /></el-icon>
                <span class="guide-name">{{ g }}</span>
                <el-button link type="primary">查看</el-button>
              </div>
            </div>
          </template>
        </section>

        <!-- 4. 审批流程 -->
        <section class="card approval-card">
          <div class="approval-controls">
            <div class="ap-tabs">
              <span
                class="ap-pill"
                :class="{ active: approvalType === '' }"
                @click="approvalType = ''"
              >全部({{ counts.all || 0 }})</span>
              <span
                v-for="t in ['EHR', '财务', 'OA']"
                :key="t"
                class="ap-pill"
                :class="{ active: approvalType === t }"
                @click="approvalType = t"
              >{{ t }}</span>
            </div>
            <el-input
              v-model="approvalKeyword"
              placeholder="流水号 / 关键字"
              clearable
              class="ap-search"
              @keyup.enter="loadApprovals"
              :prefix-icon="Search"
            />
            <el-radio-group v-model="approvalStatus" size="small">
              <el-radio-button value="todo">待办</el-radio-button>
              <el-radio-button value="done">已办</el-radio-button>
              <el-radio-button value="mine">我申请的</el-radio-button>
            </el-radio-group>
          </div>

          <div class="approval-list">
            <div v-for="a in approvals" :key="a.id" class="approval-row">
              <el-tag :type="typeTag(a.type)" size="small" effect="light">{{ a.type }}</el-tag>
              <div class="ap-main">
                <div class="ap-title">{{ a.title }}</div>
                <div class="ap-meta">流水号：{{ a.serial_no }} · 申请人：{{ a.applicant }} · 当前节点：{{ a.current_node }}</div>
              </div>
              <div class="ap-time">{{ formatTime(a.created_at) }}</div>
              <el-button type="primary" size="small" plain @click="ElMessage.info('审批详情功能开发中')">处理</el-button>
            </div>
            <el-empty v-if="approvals.length === 0" description="暂无相关审批" />
          </div>
        </section>
      </div>

      <!-- ============ 右侧边栏 ============ -->
      <aside class="side-col">
        <!-- 我的常用 -->
        <section class="card fav-card">
          <div class="side-header">
            <span class="side-title">我的常用</span>
            <el-icon class="side-edit" :class="{ editing: editFav }" @click="editFav = !editFav">
              <Delete />
            </el-icon>
          </div>
          <div class="fav-grid">
            <div v-for="f in favorites" :key="f.id" class="fav-item" @click="clickFavorite(f)">
              <span v-if="editFav" class="fav-remove" @click.stop="removeFavorite(f.id)">
                <el-icon><CloseBold /></el-icon>
              </span>
              <div class="fav-icon">
                <el-icon :size="22"><component :is="f.icon || 'Star'" /></el-icon>
              </div>
              <div class="fav-name">{{ f.name }}</div>
            </div>
            <div class="fav-item fav-add" @click="addFavorite">
              <div class="fav-icon add-icon">
                <el-icon :size="22"><Plus /></el-icon>
              </div>
              <div class="fav-name">添加</div>
            </div>
          </div>
        </section>

        <!-- 推广图 -->
        <section class="card promo-card">
          <el-carousel :height="'170px'" arrow="never" indicator-position="outside" :autoplay="true" :interval="5000">
            <el-carousel-item v-for="(p, i) in promos" :key="i">
              <div class="promo-slide">
                <img :src="p.img" :alt="p.title" />
                <div class="promo-caption">{{ p.title }}</div>
              </div>
            </el-carousel-item>
          </el-carousel>
        </section>

        <!-- 我的日程 -->
        <section class="card schedule-card">
          <div class="side-header">
            <span class="side-title">我的日程</span>
          </div>
          <el-calendar v-model="calDate">
            <template #header>
              <div class="cal-header">
                <el-button circle size="small" @click="changeMonth(-1)">
                  <el-icon><ArrowLeftBold /></el-icon>
                </el-button>
                <span class="cal-title">{{ calTitle }}</span>
                <el-button circle size="small" @click="changeMonth(1)">
                  <el-icon><ArrowRightBold /></el-icon>
                </el-button>
                <el-button link type="primary" size="small" @click="goToday">今天</el-button>
              </div>
            </template>
            <template #date-cell="{ data }">
              <div
                class="cal-cell"
                :class="{
                  'is-other': data.type !== 'current-month',
                  'is-today': data.day === todayDay
                }"
              >
                <span class="cal-num">{{ Number(data.day.split('-')[2]) }}</span>
                <div v-for="ev in eventMap[data.day] || []" :key="ev.id" class="cal-event">{{ ev.title }}</div>
              </div>
            </template>
          </el-calendar>
        </section>
      </aside>
    </div>

    <!-- 悬浮 400 小助手 -->
    <el-popover placement="top-end" :width="280" trigger="click">
      <template #reference>
        <div class="float-assistant">
          <el-icon :size="28" color="#fff"><Service /></el-icon>
          <span class="fa-badge">400</span>
        </div>
      </template>
      <div class="assistant-panel">
        <div class="as-title">
          <el-icon color="#00b578"><Service /></el-icon> 400 小助手
        </div>
        <div class="as-msg">您好，我是 400 小助手，请问有什么可以帮您？</div>
        <div class="as-quicks">
          <el-tag size="small" class="as-q" effect="plain" @click="ElMessage.info('正在为您转接人工客服…')">转人工客服</el-tag>
          <el-tag size="small" class="as-q" effect="plain" @click="ElMessage.info('常见问题即将上线')">常见问题</el-tag>
        </div>
      </div>
    </el-popover>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Document, Delete } from '@element-plus/icons-vue'
import { bffApi, type MeInfo, type Approval, type ScheduleEvent, type Favorite } from '../api/bff'

const router = useRouter()

// ---------- 图片资源（统一走文本生成图服务） ----------
const img = (prompt: string, size: string) =>
  `https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=${encodeURIComponent(prompt)}&image_size=${size}`

const banners = [
  {
    img: img('企业大学线上培训活动横幅 深绿色科技感背景 抽象光线 现代商务风格', 'landscape_16_9'),
    title: '新员工训练营 · 2026 秋季班',
    sub: '赋能成长 · 链接未来 · 即刻报名'
  },
  {
    img: img('集团年度盛典颁奖典礼 深色背景金色粒子 商务大气横幅海报', 'landscape_16_9'),
    title: 'FY2026 年度评优正式启动',
    sub: '致敬每一位努力奋斗的新东方人'
  }
]

const promos = [
  { img: img('钉钉视频会议使用规范宣传图 蓝绿色扁平化商务插画 在线会议场景', 'landscape_4_3'), title: '钉钉视频会议使用规范' },
  { img: img('企业直播使用规范指南 绿色简约商务风 培训直播场景插画', 'landscape_4_3'), title: '直播使用规范指南' }
]

// ---------- 问候 ----------
const me = ref<MeInfo | null>(null)
const now = new Date()
const weekMap = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
const todayStr = `${now.getFullYear()} 年 ${now.getMonth() + 1} 月 ${now.getDate()} 日`
const weekStr = weekMap[now.getDay()]
const todayDay =
  `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`

// ---------- Banner 折叠 ----------
const bannerCollapsed = ref(false)

// ---------- 办事大厅 ----------
interface AppItem { name: string; icon: string; color: string; category: string; path?: string }
const hallTab = ref('hall')
const appKeyword = ref('')
const activeCategory = ref('推荐')
const categories = ['推荐', '人事', '行政', '财务法务', '市场营销', '员工福利', '智能办公', '业务工具']

const apps: AppItem[] = [
  { name: 'ERP', icon: 'Box', color: '#3b82f6', category: '推荐', path: '/erp' },
  { name: '新东方之家', icon: 'House', color: '#f59e0b', category: '推荐' },
  { name: 'SCRM', icon: 'ChatDotRound', color: '#00b578', category: '推荐', path: '/crm' },
  { name: '极客通', icon: 'Cellphone', color: '#8b5cf6', category: '推荐' },
  { name: '企业微信', icon: 'ChatLineRound', color: '#22c55e', category: '推荐' },
  { name: 'OA 系统', icon: 'Monitor', color: '#0ea5e9', category: '推荐', path: '/iam' },
  { name: '主数据', icon: 'Files', color: '#ef4444', category: '推荐', path: '/mdm' },
  { name: '工单系统', icon: 'Tickets', color: '#ec4899', category: '推荐', path: '/crm' },
  { name: 'EHR', icon: 'UserFilled', color: '#0ea5e9', category: '人事' },
  { name: '我的考勤', icon: 'Clock', color: '#f59e0b', category: '人事' },
  { name: '请假申请', icon: 'Document', color: '#10b981', category: '人事' },
  { name: '薪资查询', icon: 'Money', color: '#ef4444', category: '人事' },
  { name: '会议室预定', icon: 'Calendar', color: '#00b578', category: '行政' },
  { name: '用车申请', icon: 'Van', color: '#6366f1', category: '行政' },
  { name: '名片印刷', icon: 'Printer', color: '#f97316', category: '行政' },
  { name: '办公用品', icon: 'ShoppingBag', color: '#14b8a6', category: '行政' },
  { name: '我的报销', icon: 'Wallet', color: '#22c55e', category: '财务法务' },
  { name: '付款申请', icon: 'CreditCard', color: '#3b82f6', category: '财务法务' },
  { name: '合同审批', icon: 'Stamp', color: '#eab308', category: '财务法务' },
  { name: '法务咨询', icon: 'Reading', color: '#8b5cf6', category: '财务法务' },
  { name: '营销素材', icon: 'Picture', color: '#ec4899', category: '市场营销' },
  { name: '客户管理', icon: 'User', color: '#00b578', category: '市场营销', path: '/crm' },
  { name: '销售订单', icon: 'List', color: '#f59e0b', category: '市场营销', path: '/crm' },
  { name: '合同管理', icon: 'Notebook', color: '#06b6d4', category: '市场营销', path: '/crm' },
  { name: '福利商城', icon: 'Present', color: '#ef4444', category: '员工福利' },
  { name: '年度体检', icon: 'FirstAidKit', color: '#10b981', category: '员工福利' },
  { name: '健身房', icon: 'Basketball', color: '#f97316', category: '员工福利' },
  { name: '食堂菜单', icon: 'Food', color: '#84cc16', category: '员工福利' },
  { name: '智能填表', icon: 'EditPen', color: '#6366f1', category: '智能办公' },
  { name: '流程中心', icon: 'Share', color: '#0ea5e9', category: '智能办公' },
  { name: '知识库', icon: 'Collection', color: '#eab308', category: '智能办公' },
  { name: '审批中心', icon: 'Finished', color: '#00b578', category: '智能办公' },
  { name: '数据看板', icon: 'DataAnalysis', color: '#3b82f6', category: '业务工具' },
  { name: '库存管理', icon: 'Box', color: '#14b8a6', category: '业务工具', path: '/erp' },
  { name: '采购订单', icon: 'ShoppingCart', color: '#f59e0b', category: '业务工具', path: '/erp' },
  { name: '生产工单', icon: 'Setting', color: '#8b5cf6', category: '业务工具', path: '/erp' }
]

const guides = [
  '新员工入职操作指南',
  'EHR 人事系统使用手册',
  '报销填报流程说明',
  '会议室预定操作指引',
  '视频会议 / 直播使用规范'
]

const filteredApps = computed(() => {
  let list = apps.filter((a) => a.category === activeCategory.value)
  if (appKeyword.value.trim()) {
    const kw = appKeyword.value.trim()
    list = apps.filter((a) => a.name.includes(kw))
  }
  return list
})

const openApp = (app: AppItem) => {
  if (app.path) {
    router.push(app.path).catch(() => { /* qiankun 加载失败时忽略 */ })
  } else {
    ElMessage.info(`「${app.name}」应用即将上线，敬请期待`)
  }
}

// ---------- 审批 ----------
const approvalType = ref('')
const approvalStatus = ref('todo')
const approvalKeyword = ref('')
const approvals = ref<Approval[]>([])
const counts = ref<Record<string, number>>({ all: 0, todo: 0, done: 0, mine: 0 })

const loadApprovals = async () => {
  try {
    const [list, c] = await Promise.all([
      bffApi.listApprovals({
        type: approvalType.value,
        status: approvalStatus.value,
        keyword: approvalKeyword.value
      }),
      bffApi.approvalCounts(approvalType.value)
    ])
    approvals.value = list.items || []
    counts.value = c || {}
  } catch {
    approvals.value = []
  }
}

watch([approvalType, approvalStatus], loadApprovals)

const typeTag = (t: string) => (t === '财务' ? 'warning' : t === 'OA' ? 'success' : 'primary')

const formatTime = (t: string) => {
  if (!t) return ''
  const d = new Date(t)
  return `${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

// ---------- 我的常用 ----------
const favorites = ref<Favorite[]>([])
const editFav = ref(false)

const loadFavorites = async () => {
  try {
    const res = await bffApi.listFavorites()
    favorites.value = res.items || []
  } catch { /* 忽略 */ }
}

const clickFavorite = (f: Favorite) => {
  if (editFav.value) return
  if (f.path) router.push(f.path)
  else ElMessage.info(`「${f.name}」即将上线，敬请期待`)
}

const removeFavorite = async (id: number) => {
  try {
    await bffApi.deleteFavorite(id)
    favorites.value = favorites.value.filter((f) => f.id !== id)
    ElMessage.success('已移除')
  } catch { ElMessage.error('移除失败') }
}

const addFavorite = async () => {
  try {
    const { value } = await ElMessageBox.prompt('请输入常用应用名称', '添加常用', {
      confirmButtonText: '添加',
      cancelButtonText: '取消'
    })
    if (!value) return
    const f = await bffApi.addFavorite({ name: value })
    favorites.value.push(f)
  } catch { /* 取消 */ }
}

// ---------- 日程 ----------
const calDate = ref(new Date())
const schedules = ref<ScheduleEvent[]>([])
const eventMap = computed(() => {
  const map: Record<string, ScheduleEvent[]> = {}
  for (const ev of schedules.value) {
    const day = ev.event_date.slice(0, 10)
    ;(map[day] ||= []).push(ev)
  }
  return map
})

const calTitle = computed(() => `${calDate.value.getFullYear()} 年 ${calDate.value.getMonth() + 1} 月`)

const changeMonth = (delta: number) => {
  const d = new Date(calDate.value)
  d.setDate(1)
  d.setMonth(d.getMonth() + delta)
  calDate.value = d
}

const goToday = () => { calDate.value = new Date() }

const loadSchedules = async () => {
  const month = `${calDate.value.getFullYear()}-${String(calDate.value.getMonth() + 1).padStart(2, '0')}`
  try {
    const res = await bffApi.listSchedules(month)
    schedules.value = res.items || []
  } catch { schedules.value = [] }
}

watch(calDate, loadSchedules)

onMounted(async () => {
  try { me.value = await bffApi.me() } catch { /* 忽略 */ }
  loadApprovals()
  loadFavorites()
  loadSchedules()
})
</script>

<style scoped>
.portal-home {
  max-width: 1440px;
  margin: 0 auto;
  padding: 18px 24px 40px;
  position: relative;
}
.home-grid {
  display: grid;
  grid-template-columns: minmax(0, 3fr) minmax(300px, 1fr);
  gap: 16px;
}
.main-col { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.side-col { display: flex; flex-direction: column; gap: 16px; min-width: 0; }

.card {
  background: #fff;
  border-radius: 12px;
  padding: 18px 20px;
  box-shadow: 0 1px 4px rgba(17, 34, 51, 0.05);
}

/* ===== 问候 ===== */
.greeting-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(120deg, #ffffff 0%, #f3fdf9 100%);
}
.greet-left { display: flex; align-items: center; gap: 14px; }
.sun-icon {
  width: 54px; height: 54px; border-radius: 50%;
  background: linear-gradient(135deg, #ffd66b, #ff9d2e);
  display: flex; align-items: center; justify-content: center;
  color: #fff; box-shadow: 0 6px 16px rgba(255, 157, 46, .35);
}
.greet-title { font-size: 22px; font-weight: 700; color: #162432; }
.greet-suffix { font-size: 16px; font-weight: 500; color: #00b578; margin-left: 2px; }
.greet-sub { font-size: 13px; color: #8a94a0; margin-top: 6px; }
.greet-sub .days { color: #00b578; font-size: 15px; margin: 0 2px; }
.greet-date { text-align: right; }
.gd-main { font-size: 15px; font-weight: 600; color: #162432; }
.gd-week { font-size: 12px; color: #97a0ab; margin-top: 4px; }

/* ===== Banner ===== */
.banner-wrap { position: relative; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 4px rgba(17,34,51,.05); }
.banner-slide {
  position: relative;
  height: 170px;
  background-size: cover;
  background-position: center;
  border-radius: 12px;
  overflow: hidden;
}
.banner-mask {
  position: absolute; inset: 0;
  background: linear-gradient(90deg, rgba(0, 20, 14, .72) 0%, rgba(0, 20, 14, .18) 55%, rgba(0,20,14,0) 100%);
}
.banner-text { position: absolute; left: 26px; top: 50%; transform: translateY(-50%); color: #fff; z-index: 2; }
.bt-title { font-size: 24px; font-weight: 700; letter-spacing: 1px; text-shadow: 0 2px 8px rgba(0,0,0,.3); }
.bt-sub { font-size: 14px; opacity: .9; margin-top: 10px; }
.banner-collapse {
  position: absolute; right: 14px; top: 10px; z-index: 3;
  color: #fff; font-size: 12px; cursor: pointer;
  background: rgba(0,0,0,.28); padding: 3px 10px; border-radius: 20px;
  display: flex; align-items: center; gap: 2px;
}
.banner-collapse:hover { background: rgba(0,0,0,.5); }
.banner-collapsed-bar {
  background: #fff; border-radius: 12px; padding: 10px 18px; font-size: 13px; color: #5f6b7a;
  display: flex; align-items: center; cursor: pointer;
  box-shadow: 0 1px 4px rgba(17,34,51,.05);
}
.bc-dot { width: 8px; height: 8px; border-radius: 50%; background: #00b578; margin-right: 8px; }
.bc-expand { margin-left: auto; color: #00b578; display: flex; align-items: center; gap: 2px; }

/* ===== 办事大厅 ===== */
.hall-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
.hall-search { width: 220px; }
.cat-tabs { display: flex; flex-wrap: wrap; gap: 10px; margin-bottom: 16px; }
.cat-pill {
  padding: 5px 16px; border-radius: 20px; font-size: 13px; cursor: pointer;
  color: #5f6b7a; background: #f3f5f8; transition: all .2s;
  user-select: none;
}
.cat-pill:hover { color: #00b578; }
.cat-pill.active { background: #00b578; color: #fff; font-weight: 600; box-shadow: 0 3px 8px rgba(0,181,120,.28); }
.app-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 14px 8px;
}
.app-item { display: flex; flex-direction: column; align-items: center; gap: 8px; cursor: pointer; }
.app-item:hover .app-icon { transform: translateY(-3px); box-shadow: 0 6px 14px rgba(0,181,120,.18); }
.app-icon {
  width: 52px; height: 52px; border-radius: 14px;
  display: flex; align-items: center; justify-content: center;
  transition: all .25s;
}
.app-name { font-size: 12px; color: #41505f; text-align: center; }
.guide-list { display: flex; flex-direction: column; }
.guide-item {
  display: flex; align-items: center; gap: 10px;
  padding: 12px 4px; border-bottom: 1px solid #f1f3f6; cursor: pointer;
}
.guide-item:last-child { border-bottom: none; }
.guide-icon { color: #00b578; font-size: 17px; }
.guide-name { flex: 1; font-size: 14px; color: #33404f; }

/* ===== 审批 ===== */
.approval-controls { display: flex; align-items: center; gap: 14px; margin-bottom: 14px; flex-wrap: wrap; }
.ap-tabs { display: flex; gap: 8px; }
.ap-pill {
  padding: 5px 14px; border-radius: 20px; font-size: 13px; cursor: pointer;
  color: #5f6b7a; background: #f3f5f8; user-select: none;
}
.ap-pill.active { background: #00b578; color: #fff; font-weight: 600; }
.ap-search { width: 210px; }
.approval-controls .el-radio-group { margin-left: auto; }
.approval-list { min-height: 120px; }
.approval-row {
  display: flex; align-items: center; gap: 14px;
  padding: 12px 6px; border-bottom: 1px solid #f1f3f6;
}
.approval-row:last-child { border-bottom: none; }
.ap-main { flex: 1; min-width: 0; }
.ap-title { font-size: 14px; color: #1f2d3d; font-weight: 600; }
.ap-meta { font-size: 12px; color: #97a0ab; margin-top: 4px; }
.ap-time { font-size: 12px; color: #97a0ab; white-space: nowrap; }

/* ===== 右侧通用 ===== */
.side-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
.side-title { font-size: 15px; font-weight: 700; color: #162432; }
.side-edit { color: #b0b8c2; cursor: pointer; font-size: 16px; }
.side-edit:hover { color: #ef4444; }
.side-edit.editing { color: #ef4444; }

/* 我的常用 */
.fav-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.fav-item {
  position: relative;
  display: flex; flex-direction: column; align-items: center; gap: 8px;
  cursor: pointer; padding: 6px 0;
}
.fav-icon {
  width: 46px; height: 46px; border-radius: 13px;
  background: #ecfaf6; color: #00b578;
  display: flex; align-items: center; justify-content: center;
}
.fav-item:hover .fav-icon { background: #d9f4ec; }
.fav-name { font-size: 12px; color: #41505f; text-align: center; }
.add-icon { background: #f3f5f8 !important; color: #97a0ab !important; border: 1px dashed #cdd4dc; }
.fav-remove {
  position: absolute; top: -2px; right: 8px;
  width: 18px; height: 18px; border-radius: 50%;
  background: #ef4444; color: #fff; font-size: 11px;
  display: flex; align-items: center; justify-content: center; z-index: 2;
}

/* 推广 */
.promo-slide { height: 170px; position: relative; border-radius: 8px; overflow: hidden; }
.promo-slide img { width: 100%; height: 100%; object-fit: cover; border-radius: 8px; }
.promo-caption {
  position: absolute; left: 0; right: 0; bottom: 0;
  background: linear-gradient(transparent, rgba(0,0,0,.7));
  color: #fff; font-size: 12px; padding: 18px 10px 8px;
}

/* 日历 */
.cal-header { display: flex; align-items: center; gap: 10px; }
.cal-title { font-size: 14px; font-weight: 700; color: #162432; min-width: 108px; text-align: center; }
:deep(.el-calendar__body) { padding: 8px 4px 4px; }
:deep(.el-calendar-table .el-calendar-day) { height: 44px; padding: 2px; }
.cal-cell { height: 100%; box-sizing: border-box; }
.cal-num { font-size: 12px; color: #33404f; }
.cal-cell.is-other .cal-num { color: #c6cdd5; }
.cal-cell.is-today .cal-num {
  display: inline-flex; align-items: center; justify-content: center;
  width: 22px; height: 22px; border-radius: 50%;
  background: #00b578; color: #fff; font-weight: 700;
}
.cal-event {
  font-size: 10px; color: #00a06b; background: #ecfaf6;
  border-radius: 3px; padding: 0 3px; margin-top: 2px;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}

/* 悬浮助手 */
.float-assistant {
  position: fixed; right: 28px; bottom: 42px;
  width: 58px; height: 58px; border-radius: 50%;
  background: linear-gradient(135deg, #00c987, #00a06b);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; z-index: 200;
  box-shadow: 0 8px 22px rgba(0, 181, 120, .42);
  transition: transform .2s;
}
.float-assistant:hover { transform: translateY(-3px) scale(1.04); }
.fa-badge {
  position: absolute; right: -4px; top: -4px;
  background: #ef4444; color: #fff; font-size: 10px; font-weight: 700;
  border-radius: 10px; padding: 1px 6px; border: 2px solid #fff;
}
.assistant-panel .as-title { font-weight: 700; display: flex; align-items: center; gap: 6px; margin-bottom: 10px; }
.as-msg { font-size: 13px; color: #41505f; background: #ecfaf6; border-radius: 8px; padding: 10px; }
.as-quicks { margin-top: 10px; display: flex; gap: 8px; }
.as-q { cursor: pointer; }

@media (max-width: 1100px) {
  .home-grid { grid-template-columns: 1fr; }
  .app-grid { grid-template-columns: repeat(6, 1fr); }
}
</style>
