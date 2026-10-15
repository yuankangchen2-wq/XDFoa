<template>
  <div class="home-page">
    <el-tabs v-model="role" class="role-tabs" @tab-change="loadTodos">
      <el-tab-pane label="员工首页" name="employee" />
      <el-tab-pane label="主管首页" name="supervisor" />
    </el-tabs>

    <el-row :gutter="12">
      <!-- 左：待办事项 -->
      <el-col :span="16">
        <el-card shadow="never" class="todo-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">待办事项</span>
              <el-button text :icon="Refresh" @click="loadAll">刷新</el-button>
            </div>
          </template>
          <div v-for="g in todoGroups" :key="g.name" class="todo-group">
            <div class="group-name">{{ g.name }}</div>
            <div class="todo-grid">
              <div
                v-for="t in g.items"
                :key="t.key"
                class="todo-item"
                @click="gotoTodo(t)"
              >
                <div class="todo-icon" :style="{ background: groupColor(g.name) }">
                  <el-icon><component :is="groupIcon(g.name)" /></el-icon>
                </div>
                <span class="todo-label">{{ t.label }}</span>
                <span v-if="t.count > 0" class="todo-badge">{{ t.count > 99 ? '99+' : t.count }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右：日程表 -->
      <el-col :span="8">
        <el-card shadow="never" class="sched-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">日程表</span>
              <div>
                <el-button text type="primary" @click="openCreate()">新建日程</el-button>
                <el-button text type="primary" @click="selected = todayStr; loadDayEvents()">管理日程</el-button>
              </div>
            </div>
          </template>
          <div class="cal-head">
            <el-button text :icon="ArrowLeft" @click="shiftMonth(-1)" />
            <span class="cal-ym">{{ calYear }} 年 {{ calMonth + 1 }} 月</span>
            <el-button text :icon="ArrowRight" @click="shiftMonth(1)" />
            <el-button text size="small" @click="goThisMonth">今天</el-button>
          </div>
          <div class="cal-grid cal-row">
            <div v-for="w in weekDays" :key="w" class="cal-week">{{ w }}</div>
          </div>
          <div class="cal-grid">
            <div
              v-for="(d, i) in calDays"
              :key="i"
              class="cal-cell"
              :class="{
                muted: !d.inMonth,
                today: d.iso === todayStr,
                selected: d.iso === selected
              }"
              @click="pickDay(d)"
            >
              <span class="cal-num">{{ d.day }}</span>
              <span class="cal-dots">
                <i
                  v-for="(e, ei) in (eventsMap[d.iso] || []).slice(0, 3)"
                  :key="ei"
                  class="dot"
                  :class="dotClass(e.status)"
                />
              </span>
            </div>
          </div>
          <div class="day-events">
            <div class="de-title">{{ selected }} 的日程</div>
            <div v-if="dayEvents.length === 0" class="de-empty">暂无日程</div>
            <div v-for="e in dayEvents" :key="e.id" class="de-row">
              <span class="de-time">{{ e.event_time || '全天' }}</span>
              <span class="de-name">{{ e.title }}</span>
              <el-tag size="small" :type="statusTag(e.status)" effect="plain">{{ e.status }}</el-tag>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="12" style="margin-top:12px">
      <!-- 顾问排名 -->
      <el-col :span="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">顾问排名</span>
              <div class="filters">
                <el-select v-model="campus" size="small" style="width:120px" @change="loadRanking">
                  <el-option label="全部校区" value="全部校区" />
                  <el-option v-for="c in campusList" :key="c" :label="c" :value="c" />
                </el-select>
                <el-select v-model="month" size="small" style="width:130px" @change="loadRanking">
                  <el-option label="按月" value="" disabled />
                  <el-option v-for="m in monthList" :key="m" :label="m" :value="m" />
                </el-select>
                <el-button text :icon="Refresh" @click="loadRanking" />
              </div>
            </div>
          </template>
          <el-table :data="ranking" size="small" stripe>
            <el-table-column label="排名" width="80" align="center">
              <template #default="{ $index }">
                <span class="rank-badge" :class="'rank-' + ($index + 1)">{{ $index + 1 }}</span>
              </template>
            </el-table-column>
            <el-table-column label="姓名">
              <template #default="{ row }">
                <el-avatar :size="24" class="rank-avatar">{{ row.name.slice(0, 1) }}</el-avatar>
                <span style="margin-left:8px">{{ row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="campus" label="校区" width="110" />
            <el-table-column label="成单金额" sortable :sort-by="'deal_amount'" align="right">
              <template #default="{ row }">
                <span class="money">¥{{ formatAmount(row.deal_amount) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="deal_persons" label="成单人数" sortable align="center" />
            <el-table-column prop="deal_visits" label="成单人次" sortable align="center" />
          </el-table>
        </el-card>
      </el-col>

      <!-- 成单年级分布 -->
      <el-col :span="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">成单年级分布</span>
              <div class="filters">
                <el-select v-model="distCampus" size="small" style="width:110px" @change="loadGradeDist">
                  <el-option label="全部" value="" />
                  <el-option v-for="c in campusList" :key="c" :label="c" :value="c" />
                </el-select>
                <el-date-picker
                  v-model="distMonthDate"
                  type="month"
                  size="small"
                  style="width:130px"
                  value-format="YYYY-MM"
                  @change="loadGradeDist"
                />
              </div>
            </div>
          </template>
          <div v-if="gradeDist.length === 0" class="chart-empty">
            <el-empty description="暂无数据" :image-size="70" />
          </div>
          <div v-else class="bar-chart">
            <div v-for="g in gradeDist" :key="g.grade" class="bar-row">
              <span class="bar-label">{{ g.grade }}</span>
              <div class="bar-track">
                <div class="bar-fill" :style="{ width: barWidth(g.count) + '%' }" />
              </div>
              <span class="bar-value">{{ g.count }} 单</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 新建日程弹窗 -->
    <el-dialog v-model="schedDlg" title="新建日程" width="420px">
      <el-form label-width="80px">
        <el-form-item label="事项" required>
          <el-input v-model="schedForm.title" placeholder="请输入日程事项" />
        </el-form-item>
        <el-form-item label="日期" required>
          <el-date-picker v-model="schedForm.date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="时间">
          <el-time-picker v-model="schedForm.time" value-format="HH:mm" format="HH:mm" placeholder="如 14:00" style="width:100%" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="schedForm.type" style="width:100%">
            <el-option label="日程" value="日程" />
            <el-option label="体验课" value="体验课" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="schedDlg = false">取消</el-button>
        <el-button type="primary" @click="submitSchedule">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Refresh, ArrowLeft, ArrowRight,
  Money, Timer, Calendar, Warning, Tickets, User
} from '@element-plus/icons-vue'
import { scrmApi, type TodoItem, type ScheduleEvent } from '../api/scrm'

const router = useRouter()
const role = ref('employee')
const todos = ref<TodoItem[]>([])

const groupOrder = ['商机类', '体验课类', '日程类', '线索类', '任务/其他']
const todoGroups = computed(() =>
  groupOrder
    .map(name => ({ name, items: todos.value.filter(t => t.group === name) }))
    .filter(g => g.items.length > 0)
)

function groupColor(g: string) {
  const map: Record<string, string> = {
    '商机类': 'linear-gradient(135deg,#36d1dc,#5b86e5)',
    '体验课类': 'linear-gradient(135deg,#ff9966,#ff5e62)',
    '日程类': 'linear-gradient(135deg,#43cea2,#185a9d)',
    '线索类': 'linear-gradient(135deg,#11998e,#38ef7d)',
    '任务/其他': 'linear-gradient(135deg,#c471f5,#fa71cd)'
  }
  return map[g] || '#00a86b'
}
function groupIcon(g: string) {
  return { '商机类': Money, '体验课类': User, '日程类': Calendar, '线索类': Warning, '任务/其他': Tickets }[g] || Timer
}

async function loadTodos() {
  todos.value = await scrmApi.todos(role.value)
}

function gotoTodo(t: TodoItem) {
  const q: Record<string, string> = {}
  if (t.query) {
    t.query.split('&').forEach(kv => {
      const [k, v] = kv.split('=')
      q[k] = decodeURIComponent(v || '')
    })
  }
  router.push({ path: t.route, query: q })
}

// ---------------- 日程 ----------------
const weekDays = ['一', '二', '三', '四', '五', '六', '日']
const now0 = new Date()
const calYear = ref(now0.getFullYear())
const calMonth = ref(now0.getMonth())
const todayStr = isoDate(now0)
const selected = ref(todayStr)
const schedules = ref<ScheduleEvent[]>([])
const dayEvents = ref<ScheduleEvent[]>([])

function isoDate(d: Date) {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const eventsMap = computed(() => {
  const m: Record<string, ScheduleEvent[]> = {}
  for (const e of schedules.value) {
    const k = (e.event_date || '').slice(0, 10)
    if (!m[k]) m[k] = []
    m[k].push(e)
  }
  return m
})

const calDays = computed(() => {
  const first = new Date(calYear.value, calMonth.value, 1)
  const startWeekday = (first.getDay() + 6) % 7 // 周一为 0
  const cells: { day: number; iso: string; inMonth: boolean }[] = []
  for (let i = 0; i < 42; i++) {
    const d = new Date(calYear.value, calMonth.value, i - startWeekday + 1)
    cells.push({
      day: d.getDate(),
      iso: isoDate(d),
      inMonth: d.getMonth() === calMonth.value
    })
  }
  return cells
})

function shiftMonth(n: number) {
  let m = calMonth.value + n
  let y = calYear.value
  if (m < 0) { m = 11; y-- }
  if (m > 11) { m = 0; y++ }
  calMonth.value = m
  calYear.value = y
  loadSchedules()
}
function goThisMonth() {
  calYear.value = now0.getFullYear()
  calMonth.value = now0.getMonth()
  selected.value = todayStr
  loadSchedules()
}
function pickDay(d: { iso: string; inMonth: boolean }) {
  if (!d.inMonth) return
  selected.value = d.iso
  loadDayEvents()
}
function dotClass(status: string) {
  if (status === '待确访') return 'dot-orange'
  if (status === '到访未转化') return 'dot-red'
  return 'dot-green'
}
function statusTag(status: string) {
  if (status === '待确访') return 'warning'
  if (status === '到访未转化') return 'danger'
  return 'success'
}

async function loadSchedules() {
  const mm = `${calYear.value}-${String(calMonth.value + 1).padStart(2, '0')}`
  schedules.value = await scrmApi.schedules(mm)
  loadDayEvents()
}
function loadDayEvents() {
  dayEvents.value = eventsMap.value[selected.value] || []
}

const schedDlg = ref(false)
const schedForm = ref({ title: '', date: todayStr, time: '', type: '日程' })
function openCreate() {
  schedForm.value = { title: '', date: selected.value, time: '', type: '日程' }
  schedDlg.value = true
}
async function submitSchedule() {
  if (!schedForm.value.title.trim()) {
    ElMessage.warning('请填写事项')
    return
  }
  await scrmApi.createSchedule(schedForm.value)
  ElMessage.success('日程已创建')
  schedDlg.value = false
  await loadSchedules()
}

// ---------------- 排名 / 年级分布 ----------------
const campusList = ['武汉校区', '武昌校区', '汉口校区', '江汉校区']
const campus = ref('全部校区')
const month = ref(now0.toISOString().slice(0, 7))
const monthList = [month.value, '2026-09']
const ranking = ref<any[]>([])
async function loadRanking() {
  ranking.value = await scrmApi.ranking(month.value || todayStr.slice(0, 7), campus.value)
}
function formatAmount(v: string | number) {
  return Number(v || 0).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const distCampus = ref('')
const distMonthDate = ref(todayStr.slice(0, 7))
const gradeDist = ref<{ grade: string; count: number }[]>([])
const maxCount = computed(() => Math.max(1, ...gradeDist.value.map(g => g.count)))
function barWidth(c: number) { return Math.round((c / maxCount.value) * 100) }
async function loadGradeDist() {
  gradeDist.value = await scrmApi.gradeDist(distMonthDate.value, distCampus.value)
}

async function loadAll() {
  await Promise.all([loadTodos(), loadSchedules(), loadRanking(), loadGradeDist()])
}

onMounted(loadAll)
</script>

<style scoped>
.role-tabs { margin-bottom: 8px; }
.role-tabs :deep(.el-tabs__header) { margin: 0; }

.card-header { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-weight: 700; font-size: 15px; }
.filters { display: flex; gap: 8px; align-items: center; }

.todo-group { margin-bottom: 14px; }
.group-name { font-size: 13px; color: #8a94a6; margin-bottom: 8px; }
.todo-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 10px;
}
.todo-item {
  position: relative;
  border: 1px solid #eef0f3;
  border-radius: 8px;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all .15s;
  background: #fff;
}
.todo-item:hover { border-color: #00a86b; box-shadow: 0 2px 10px rgba(0,168,107,.12); transform: translateY(-2px); }
.todo-icon {
  width: 38px; height: 38px; border-radius: 10px;
  display: flex; align-items: center; justify-content: center;
  color: #fff; font-size: 19px;
}
.todo-label { font-size: 12px; color: #475669; text-align: center; line-height: 1.3; }
.todo-badge {
  position: absolute; top: 2px; right: 6px;
  min-width: 18px; height: 18px; padding: 0 5px;
  background: #f56c6c; color: #fff; border-radius: 9px;
  font-size: 11px; line-height: 18px; text-align: center;
}

.cal-head { display: flex; align-items: center; gap: 6px; margin-bottom: 8px; }
.cal-ym { font-weight: 600; flex: 1; text-align: center; }
.cal-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 2px; }
.cal-row { margin-bottom: 4px; }
.cal-week { text-align: center; font-size: 12px; color: #909399; padding: 4px 0; }
.cal-cell {
  height: 42px; border-radius: 6px;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  cursor: pointer; font-size: 13px; position: relative;
}
.cal-cell:hover { background: #f0fbf7; }
.cal-cell.muted { color: #c0c4cc; }
.cal-cell.today .cal-num {
  background: #00a86b; color: #fff; width: 22px; height: 22px;
  border-radius: 50%; line-height: 22px; text-align: center;
}
.cal-cell.selected { background: #e2f7ef; }
.cal-dots { height: 6px; display: flex; gap: 2px; }
.dot { width: 5px; height: 5px; border-radius: 50%; display: inline-block; }
.dot-green { background: #00a86b; }
.dot-orange { background: #e6a23c; }
.dot-red { background: #f56c6c; }

.day-events { margin-top: 10px; border-top: 1px dashed #ebeef5; padding-top: 8px; }
.de-title { font-size: 12px; color: #8a94a6; margin-bottom: 6px; }
.de-empty { font-size: 12px; color: #c0c4cc; text-align: center; padding: 8px 0; }
.de-row { display: flex; align-items: center; gap: 8px; padding: 4px 0; font-size: 13px; }
.de-time { color: #00a86b; width: 44px; flex-shrink: 0; }
.de-name { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.rank-badge {
  display: inline-block; width: 22px; height: 22px; line-height: 22px;
  border-radius: 50%; font-size: 12px; font-weight: 700;
  background: #e4e7ed; color: #909399;
}
.rank-1 { background: #f7c948; color: #fff; }
.rank-2 { background: #b0c4d9; color: #fff; }
.rank-3 { background: #d99a6c; color: #fff; }
.rank-avatar { background: #00a86b; font-size: 12px; }
.money { color: #f56c6c; font-weight: 600; }

.chart-empty { padding: 30px 0; }
.bar-chart { display: flex; flex-direction: column; gap: 14px; padding: 10px 4px; }
.bar-row { display: flex; align-items: center; gap: 10px; }
.bar-label { width: 42px; font-size: 13px; color: #475669; text-align: right; flex-shrink: 0; }
.bar-track { flex: 1; background: #f4f6f8; border-radius: 10px; height: 18px; overflow: hidden; }
.bar-fill {
  height: 100%; border-radius: 10px;
  background: linear-gradient(90deg,#00c07a,#00a86b);
  transition: width .5s;
}
.bar-value { width: 48px; font-size: 12px; color: #606f82; flex-shrink: 0; }
</style>
