<template>
  <div class="wb-page">
    <div class="wb-head">
      <div class="wb-title">{{ isSuper ? '主管工作台 · 线索' : '员工工作台 · 线索' }}</div>
      <el-radio-group v-if="isLeadTab" v-model="dateChip" size="small" @change="loadLeads">
        <el-radio-button v-for="c in dateChips" :key="c.v" :value="c.v">{{ c.t }}</el-radio-button>
      </el-radio-group>
    </div>

    <el-tabs v-model="tab" @tab-change="onTabChange">
      <el-tab-pane v-for="t in tabs" :key="t.v" :label="t.t" :name="t.v" />
    </el-tabs>

    <!-- 线索列表类 Tab -->
    <template v-if="isLeadTab">
      <div class="toolbar">
        <el-input v-model="keyword" placeholder="姓名/手机/线索编码" clearable size="small" style="width:220px"
          @keyup.enter="loadLeads" @clear="loadLeads" />
        <el-select v-model="fGrade" placeholder="在读年级" clearable size="small" style="width:120px" @change="loadLeads">
          <el-option v-for="g in grades" :key="g" :label="g" :value="g" />
        </el-select>
        <el-button type="primary" size="small" @click="loadLeads">查询</el-button>
        <div class="tb-spacer" />
        <el-popover trigger="click" width="320" placement="bottom-end">
          <template #reference>
            <el-button size="small">列配置</el-button>
          </template>
          <el-checkbox-group v-model="visibleCols">
            <el-checkbox v-for="c in allColumns" :key="c.p" :value="c.p" style="display:block">{{ c.label }}</el-checkbox>
          </el-checkbox-group>
        </el-popover>
        <template v-if="isSuper">
          <el-button size="small" @click="openAssign">分配</el-button>
          <el-button size="small" type="warning" @click="recycleVisible = true">回收</el-button>
          <el-button size="small" @click="archiveVisible = true">批量归档</el-button>
          <el-button size="small" @click="onActivate">激活</el-button>
        </template>
      </div>

      <el-table :data="rows" v-loading="loading" border size="small" height="58vh"
        @selection-change="(v: Lead[]) => selection = v" @header-dragend="onDragEnd">
        <el-table-column type="selection" width="42" fixed />
        <el-table-column v-for="c in shownColumns" :key="c.p" :prop="c.p" :label="c.label"
          :width="widths[c.p] || c.w" :min-width="c.mw || 100" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag v-if="c.p === 'status'" :type="statusTag(row.status)" size="small">{{ row.status }}</el-tag>
            <span v-else-if="c.p === 'phone'">{{ maskPhone(row.phone) }}</span>
            <span v-else-if="c.p === 'deduct_wechat'">{{ row.deduct_wechat ? '是' : '否' }}</span>
            <span v-else-if="c.p === 'created_at' || c.p === 'last_follow_at' || c.p === 'last_call_at'
                || c.p === 'next_follow_at' || c.p === 'assign_time' || c.p === 'archive_time'
                || c.p === 'ai_call_at' || c.p === 'recycle_at'">{{ fmtTime(row[c.p] || '') }}</span>
            <span v-else>{{ row[c.p] ?? '' }}</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination class="pager" background layout="total,sizes,prev,pager,next,jumper"
        :total="total" v-model:current-page="page" v-model:page-size="pageSize"
        :page-sizes="[20, 50, 100]" @change="loadLeads" />
    </template>

    <!-- 任务列表（按任务聚合 9 状态计数） -->
    <template v-else-if="tab === 'tasks'">
      <el-table :data="taskAgg" border size="small" v-loading="loading">
        <el-table-column prop="task_name" label="任务名称" min-width="260" show-overflow-tooltip />
        <el-table-column prop="task_type" label="任务类型" width="110" />
        <el-table-column prop="total" label="线索总数" width="90" />
        <el-table-column v-for="s in leadStatuses" :key="s" :label="s" width="92">
          <template #default="{ row }">{{ row.counts[s] || 0 }}</template>
        </el-table-column>
      </el-table>
    </template>

    <!-- 员工任务（任务 × 员工） -->
    <template v-else>
      <el-table :data="empTasks" border size="small" v-loading="loading">
        <el-table-column prop="task_name" label="任务名称" min-width="240" show-overflow-tooltip />
        <el-table-column prop="owner_name" label="所属员工" width="110" />
        <el-table-column prop="owner_dept" label="所属部门" width="150" />
        <el-table-column prop="task_type" label="任务类型" width="100" />
        <el-table-column prop="total_count" label="线索总数" width="90" />
        <el-table-column v-for="s in leadStatuses" :key="s" :label="s" width="92">
          <template #default="{ row }">{{ parseCounts(row.status_counts)[s] || 0 }}</template>
        </el-table-column>
      </el-table>
    </template>

    <!-- 回收弹窗 -->
    <el-dialog v-model="recycleVisible" title="线索回收" width="420px">
      <el-form label-width="90px">
        <el-form-item label="回收原因">
          <el-select v-model="recycleReason" placeholder="请选择回收原因" style="width:100%">
            <el-option v-for="r in recycleReasons" :key="r" :label="r" :value="r" />
          </el-select>
        </el-form-item>
        <el-form-item label="已选线索">{{ selection.length }} 条</el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="recycleVisible = false">取消</el-button>
        <el-button type="primary" @click="onRecycle">确认回收</el-button>
      </template>
    </el-dialog>

    <!-- 分配弹窗 -->
    <el-dialog v-model="assignVisible" title="分配线索" width="420px">
      <el-select v-model="receiver" placeholder="选择接收员工" filterable style="width:100%">
        <el-option v-for="e in employees" :key="e.id" :label="e.name" :value="e.name" />
      </el-select>
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" @click="onAssign">确认分配</el-button>
      </template>
    </el-dialog>

    <!-- 归档弹窗 -->
    <el-dialog v-model="archiveVisible" title="批量归档" width="460px">
      <el-form label-width="90px">
        <el-form-item label="归档方式">
          <el-radio-group v-model="archiveForm.way">
            <el-radio value="无效归档">无效归档</el-radio>
            <el-radio value="入库归档">入库归档</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="归档原因">
          <el-select v-model="archiveForm.reason" placeholder="请选择" style="width:100%">
            <el-option v-for="r in archiveReasons" :key="r" :label="r" :value="r" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="archiveVisible = false">取消</el-button>
        <el-button type="primary" @click="onArchive">确认归档</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { scrmApi, maskPhone, fmtTime, type Lead, type Employee } from '../api/scrm'

const route = useRoute()
const isSuper = computed(() => route.query.mode !== 'employee')

const tabs = computed(() => isSuper.value
  ? [
      { v: 'active', t: '活跃线索' }, { v: 'all', t: '全部线索' },
      { v: 'recycled', t: '已回收' }, { v: 'archived', t: '已归档' },
      { v: 'tasks', t: '任务列表' }, { v: 'empTasks', t: '员工任务' }
    ]
  : [
      { v: 'active', t: '活跃线索' }, { v: 'all', t: '全部线索' }, { v: 'tasks', t: '任务列表' }
    ])
const tab = ref('active')
const leadTabs = ['active', 'all', 'recycled', 'archived']
const isLeadTab = computed(() => leadTabs.includes(tab.value))

const grades = ['高三', '高二', '高一', '初三', '初二']
const leadStatuses = ['待分配', '待跟进', '跟进中', '已回收', '无效归档', '入库归档', '已转商', '已成单']
const dateChips = [
  { v: 'today', t: '今天' }, { v: '7d', t: '近7天' }, { v: '30d', t: '近30天' },
  { v: '3m', t: '近三个月' }, { v: '6m', t: '近半年' }, { v: 'all', t: '全部' }
]
const dateChip = ref('6m')
const recycleReasons = ['30天未跟进自动回收', '主管手动回收-线索质量差', '撞单回收', '恶意信息回收', '渠道重复数据']
const archiveReasons = ['空号停机', '明确拒绝', '重复数据入库', '竞品已报名', '30天未跟进归档']

const keyword = ref('')
const fGrade = ref('')
const rows = ref<Lead[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const selection = ref<Lead[]>([])
const employees = ref<Employee[]>([])

// ---- 列配置：显隐 + 宽度本地记忆 ----
interface Col { p: string; label: string; w?: number; mw?: number; sup?: boolean }
const allColumns = ref<Col[]>([
  { p: 'task_name', label: '任务名称', mw: 180 },
  { p: 'name', label: '姓名', w: 90 },
  { p: 'phone', label: '手机', w: 120 },
  { p: 'grade', label: '在读年级', w: 90 },
  { p: 'student_tag', label: '学员标签', w: 100 },
  { p: 'call_status', label: '接线状态', w: 90 },
  { p: 'last_follow_at', label: '最近跟进时间', w: 150 },
  { p: 'last_follow_content', label: '最近跟进内容', mw: 180 },
  { p: 'school', label: '学校', mw: 120 },
  { p: 'district', label: '区县', w: 90 },
  { p: 'status', label: '线索状态', w: 100 },
  { p: 'follow_count', label: '自跟进次数', w: 90 },
  { p: 'self_call_count', label: '自外呼次数', w: 90 },
  { p: 'total_follow_count', label: '总跟进次数', w: 90 },
  { p: 'ai_call_count', label: 'AI外呼次数', w: 95, sup: true },
  { p: 'ai_call_at', label: 'AI外呼时间', w: 150, sup: true },
  { p: 'deduct_wechat', label: '是否扣微', w: 80 },
  { p: 'last_call_at', label: '最近外呼时间', w: 150 },
  { p: 'next_follow_at', label: '下次回访', w: 150 },
  { p: 'assign_time', label: '最近分配', w: 150 },
  { p: 'activity_name', label: '关联活动', mw: 180 },
  { p: 'recycle_reason', label: '回收原因', mw: 150 },
  { p: 'archive_way', label: '归档方式', w: 90 },
  { p: 'archive_reason', label: '归档原因', mw: 140 },
  { p: 'archive_operator', label: '归档操作人', w: 100 },
  { p: 'archive_time', label: '归档时间', w: 150 },
  { p: 'owner_name', label: '所属人', w: 100 },
  { p: 'remark', label: '备注', mw: 140 },
  { p: 'created_at', label: '创建时间', w: 150 }
])
const storeKey = computed(() => `lead-wb-cols-${isSuper.value ? 'sup' : 'emp'}`)
const widths = reactive<Record<string, number>>(JSON.parse(localStorage.getItem(storeKey.value + '-w') || '{}'))
const visibleCols = ref<string[]>([])

function initCols() {
  const saved = localStorage.getItem(storeKey.value)
  const base = allColumns.value.filter(c => !c.sup || isSuper.value).map(c => c.p)
  visibleCols.value = saved ? JSON.parse(saved) : base
}
const shownColumns = computed(() => allColumns.value.filter(c => visibleCols.value.includes(c.p) && (!c.sup || isSuper.value)))

function onDragEnd(newW: number, _old: number, col: any) {
  widths[col.property] = newW
  localStorage.setItem(storeKey.value + '-w', JSON.stringify(widths))
}
function persistVisible() {
  localStorage.setItem(storeKey.value, JSON.stringify(visibleCols.value))
}
import { watch } from 'vue'
watch(visibleCols, persistVisible, { deep: true })

function dateRange(chip: string): [string, string] | [] {
  if (chip === 'all') return []
  const days: Record<string, number> = { today: 0, '7d': 7, '30d': 30, '3m': 90, '6m': 180 }
  const d = new Date()
  const end = d.toISOString().slice(0, 10)
  d.setDate(d.getDate() - (days[chip] ?? 180))
  return [d.toISOString().slice(0, 10), end]
}

async function loadLeads() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: page.value, page_size: pageSize.value,
      keyword: keyword.value, grade: fGrade.value
    }
    const rg = dateRange(dateChip.value)
    if (rg.length) { params.created_start = rg[0]; params.created_end = rg[1] }
    if (tab.value === 'active') params.statuses = '待跟进,跟进中'
    if (tab.value === 'recycled') params.status = '已回收'
    if (tab.value === 'archived') params.statuses = '无效归档,入库归档'
    const res = await scrmApi.leads(params)
    rows.value = res.items
    total.value = res.total
  } finally { loading.value = false }
}

// ---- 任务列表 / 员工任务 ----
const empTasks = ref<any[]>([])
const taskAgg = ref<any[]>([])
function parseCounts(s: string): Record<string, number> {
  try { return JSON.parse(s || '{}') } catch { return {} }
}
async function loadTasks() {
  loading.value = true
  try {
    // 任务列表/员工任务默认近半年（任务创建时间）
    const res = await scrmApi.employeeTasks({ biz_kind: 'lead', page: 1, page_size: 200 })
    empTasks.value = res.items
    const byTask = new Map<string, any>()
    for (const t of res.items) {
      const cur = byTask.get(t.task_name) || { task_name: t.task_name, task_type: t.task_type, total: 0, counts: {} }
      cur.total += t.total_count
      const c = parseCounts(t.status_counts)
      for (const k of Object.keys(c)) cur.counts[k] = (cur.counts[k] || 0) + c[k]
      byTask.set(t.task_name, cur)
    }
    taskAgg.value = [...byTask.values()]
  } finally { loading.value = false }
}

function onTabChange() {
  page.value = 1
  selection.value = []
  if (isLeadTab.value) loadLeads()
  else loadTasks()
}

function statusTag(s: string) {
  return ({ 跟进中: 'success', 待跟进: 'warning', 已回收: 'info', 已转商: 'success', 已成单: 'danger' } as Record<string, string>)[s] || 'info'
}

// ---- 主管动作 ----
const recycleVisible = ref(false)
const recycleReason = ref('')
async function onRecycle() {
  if (!selection.value.length) return ElMessage.warning('请先选择线索')
  if (!recycleReason.value) return ElMessage.warning('请选择回收原因')
  await scrmApi.recycleLeads(selection.value.map(x => x.id), recycleReason.value)
  ElMessage.success('已回收到回收池')
  recycleVisible.value = false
  loadLeads()
}

const assignVisible = ref(false)
const receiver = ref('')
async function onAssign() {
  if (!selection.value.length) return ElMessage.warning('请先选择线索')
  if (!receiver.value) return ElMessage.warning('请选择接收员工')
  await scrmApi.handoverLeads(selection.value.map(x => x.id), '', receiver.value)
  ElMessage.success('分配成功')
  assignVisible.value = false
  loadLeads()
}
function openAssign() {
  if (!selection.value.length) return ElMessage.warning('请先选择线索')
  assignVisible.value = true
}

const archiveVisible = ref(false)
const archiveForm = reactive({ way: '无效归档', reason: '' })
async function onArchive() {
  if (!selection.value.length) return ElMessage.warning('请先选择线索')
  if (!archiveForm.reason) return ElMessage.warning('请选择归档原因')
  await scrmApi.batchArchiveLeads(selection.value.map(x => x.id), archiveForm.reason, archiveForm.way, archiveForm.reason)
  ElMessage.success('归档成功')
  archiveVisible.value = false
  loadLeads()
}

async function onActivate() {
  if (!selection.value.length) return ElMessage.warning('请先选择线索')
  await ElMessageBox.confirm(`确认激活选中的 ${selection.value.length} 条线索？`, '提示', { type: 'warning' })
  await scrmApi.activateLeads(selection.value.map(x => x.id))
  ElMessage.success('已激活')
  loadLeads()
}

onMounted(async () => {
  initCols()
  employees.value = await scrmApi.employees()
  loadLeads()
})
</script>

<style scoped>
.wb-page { padding: 12px 16px; }
.wb-head { display: flex; align-items: center; justify-content: space-between; }
.wb-title { font-size: 16px; font-weight: 600; }
.toolbar { display: flex; gap: 8px; align-items: center; margin: 8px 0; flex-wrap: wrap; }
.tb-spacer { flex: 1; }
.pager { margin-top: 10px; justify-content: flex-end; }
</style>
