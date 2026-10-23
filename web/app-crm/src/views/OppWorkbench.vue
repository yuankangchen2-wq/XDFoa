<template>
  <div class="wb-page">
    <div class="wb-head">
      <div class="wb-title">{{ isSuper ? '主管工作台 · 商机' : '员工工作台 · 商机' }}</div>
      <el-radio-group v-if="isOppTab" v-model="dateChip" size="small" @change="loadOpps">
        <el-radio-button v-for="c in dateChips" :key="c.v" :value="c.v">{{ c.t }}</el-radio-button>
      </el-radio-group>
    </div>

    <el-tabs v-model="tab" @tab-change="onTabChange">
      <el-tab-pane v-for="t in tabs" :key="t.v" :label="t.t" :name="t.v" />
    </el-tabs>

    <!-- ===== 商机列表类 Tab ===== -->
    <template v-if="isOppTab">
      <div class="toolbar">
        <el-button-group>
          <el-button :type="quick === 'pending' ? 'primary' : ''" size="small" @click="setQuick('pending')">待跟进</el-button>
          <el-button :type="quick === 'today' ? 'primary' : ''" size="small" @click="setQuick('today')">今日回访</el-button>
          <el-button :type="quick === 'recycle3' ? 'primary' : ''" size="small" @click="setQuick('recycle3')">3天即回收</el-button>
        </el-button-group>
        <el-input v-model="keyword" placeholder="商机编码/姓名/手机/学员号" clearable size="small" style="width:230px"
          @keyup.enter="loadOpps" @clear="loadOpps" />
        <el-button size="small" @click="filterVisible = true">筛选配置</el-button>
        <el-button type="primary" size="small" @click="loadOpps">查询</el-button>
        <el-button size="small" @click="resetQuery">重置</el-button>
        <div class="tb-spacer" />
        <el-button size="small" @click="colVisible = true">列配置</el-button>
        <el-button type="primary" size="small" @click="openCreate">+ 新建</el-button>
        <el-button size="small" @click="onAssign">分配</el-button>
        <el-button size="small" @click="openFollowPlan">跟进计划</el-button>
        <el-button size="small" @click="onExport">导出</el-button>
        <el-button size="small" type="success" @click="openApplyDeal">申请成单</el-button>
        <el-button size="small" @click="onPushHedgehog">推送至刺猬</el-button>
        <el-button size="small" type="warning" @click="onPush91">推送到91</el-button>
        <template v-if="isSuper">
          <el-button size="small" @click="recycleVisible = true">回收</el-button>
          <el-button size="small" @click="onArchive">批量归档</el-button>
          <el-button size="small" @click="onActivate">激活</el-button>
        </template>
      </div>

      <el-table :data="rows" v-loading="loading" border size="small" height="56vh"
        @selection-change="(v: Opportunity[]) => selection = v">
        <el-table-column type="selection" width="42" fixed />
        <el-table-column v-for="c in shownColumns" :key="c.p" :prop="c.p" :label="c.label"
          :width="c.w" :min-width="c.mw || 100" :fixed="c.pin ? 'left' : undefined" show-overflow-tooltip>
          <template #header>
            <span>
              {{ c.label }}
              <el-tooltip v-if="c.score && scoreTip[c.score]" :content="scoreTip[c.score]" placement="top">
                <span class="ii">②</span>
              </el-tooltip>
            </span>
          </template>
          <template #default="{ row }">
            <el-tag v-if="c.p === 'status'" :type="oppStatusTag(row.status)" size="small">{{ row.status }}</el-tag>
            <el-tag v-else-if="c.p === 'grade2'" :type="gradeTag(row.grade2)" size="small" effect="dark">{{ row.grade2 }}</el-tag>
            <el-tag v-else-if="c.p === 'level'" :type="row.level === '高' ? 'danger' : row.level === '中' ? 'warning' : 'info'" size="small">{{ row.level }}</el-tag>
            <el-tooltip v-else-if="c.p === 'ai_diagnosis' && row.ai_diagnosis" :content="row.ai_diagnosis">
              <span class="diag">{{ row.ai_diagnosis }}</span>
            </el-tooltip>
            <span v-else-if="c.p === 'phone'">{{ maskPhone(row.phone) }}</span>
            <span v-else-if="c.p === 'add_wechat'">{{ row.add_wechat ? '已加' : '未加' }}</span>
            <span v-else-if="c.p === 'pushed_91'">{{ row.pushed_91 ? '已推送' : '未推送' }}</span>
            <span v-else-if="c.p === 'pushed_hedgehog'">{{ row.pushed_hedgehog ? '已推送' : '—' }}</span>
            <span v-else-if="timeCols.includes(c.p)">{{ fmtTime(row[c.p] || '') }}</span>
            <span v-else>{{ row[c.p] ?? '' }}</span>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination class="pager" background layout="total,prev,pager,next,jumper"
        :total="total" v-model:current-page="page" @current-change="loadOpps" />
    </template>

    <!-- ===== 已分配商机（分配流水） ===== -->
    <template v-else-if="tab === 'assigned'">
      <div class="toolbar">
        <el-input v-model="assignKw" placeholder="商机编码/姓名/学员号" size="small" style="width:230px" @keyup.enter="loadAssigns" />
        <el-select v-model="assignAuto" placeholder="分配方式" clearable size="small" style="width:130px" @change="loadAssigns">
          <el-option label="自动分配" value="1" /><el-option label="手动分配" value="0" />
        </el-select>
        <el-button type="primary" size="small" @click="loadAssigns">查询</el-button>
      </div>
      <el-table :data="assignRows" border size="small" v-loading="loading">
        <el-table-column prop="code" label="商机编码" width="150" />
        <el-table-column prop="name" label="姓名" width="90" />
        <el-table-column prop="student_no" label="学员号" width="120" />
        <el-table-column prop="grade" label="年级" width="80" />
        <el-table-column prop="school" label="学校" min-width="130" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="auto_assign" label="自动分配" width="90">
          <template #default="{ row }">
            <el-tag :type="row.auto_assign ? 'success' : 'info'" size="small">{{ row.auto_assign ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="receiver" label="接收人" width="100" />
        <el-table-column prop="receiver_dept" label="接收部门" min-width="140" />
        <el-table-column prop="status" label="商机状态" width="100" />
        <el-table-column prop="op_at" label="分配时间" width="160">
          <template #default="{ row }">{{ fmtTime(row.op_at) }}</template>
        </el-table-column>
      </el-table>
      <el-pagination class="pager" background layout="total,prev,pager,next" :total="assignTotal"
        v-model:current-page="assignPage" @current-change="loadAssigns" />
    </template>

    <!-- ===== 任务列表 / 员工任务（主管） ===== -->
    <template v-else>
      <el-table :data="tab === 'tasks' ? taskAgg : empTasks" border size="small" v-loading="loading">
        <el-table-column prop="task_name" label="任务名称" min-width="240" show-overflow-tooltip />
        <el-table-column v-if="tab === 'empTasks'" prop="owner_name" label="所属员工" width="110" />
        <el-table-column v-if="tab === 'empTasks'" prop="owner_dept" label="所属部门" width="150" />
        <el-table-column prop="task_type" label="任务类型" width="110" />
        <el-table-column prop="total_count" label="商机总数" width="90" />
        <el-table-column v-for="s in oppStatuses" :key="s" :label="s" width="92">
          <template #default="{ row }">{{ parseCounts(row.status_counts)[s] || 0 }}</template>
        </el-table-column>
      </el-table>
    </template>

    <!-- ===== 列配置弹窗（≡排序 + 📌置顶 + 👁显隐 + 重置） ===== -->
    <el-dialog v-model="colVisible" title="列配置" width="520px">
      <div class="colcfg-tip">拖拽 ≡ 排序，📌 置顶固定，👁 控制显隐；配置自动保存</div>
      <draggable-list v-model="colCfg" />
      <template #footer>
        <el-button @click="resetCols">重置</el-button>
        <el-button type="primary" @click="saveCols">保存</el-button>
      </template>
    </el-dialog>

    <!-- ===== 筛选配置弹窗 ===== -->
    <el-dialog v-model="filterVisible" title="筛选配置" width="560px">
      <el-form :model="filters" label-width="110px" size="small">
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="意向等级">
            <el-select v-model="filters.level" clearable style="width:100%">
              <el-option label="高" value="高" /><el-option label="中" value="中" /><el-option label="低" value="低" />
            </el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="商机等级">
            <el-select v-model="filters.grade2" clearable style="width:100%">
              <el-option v-for="g in ['S','A','B','C']" :key="g" :label="g" :value="g" />
            </el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="在读年级">
            <el-select v-model="filters.grade" clearable style="width:100%">
              <el-option v-for="g in grades" :key="g" :label="g" :value="g" />
            </el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="新老生">
            <el-select v-model="filters.new_old" clearable style="width:100%">
              <el-option label="新生" value="新生" /><el-option label="老学员" value="老学员" />
            </el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="接触点">
            <el-select v-model="filters.contact_point" clearable style="width:100%">
              <el-option label="呼入" value="呼入" /><el-option label="呼出" value="呼出" />
            </el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="进入类型">
            <el-select v-model="filters.enter_type" clearable style="width:100%">
              <el-option v-for="e in ['直接创建商机','激活商机','优秀跟进创建']" :key="e" :label="e" :value="e" />
            </el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="意向部门">
            <el-input v-model="filters.intention_dept" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="意向项目">
            <el-input v-model="filters.intention_project" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="所属人">
            <el-input v-model="filters.owner_name" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="学校">
            <el-input v-model="filters.school" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="分数区间">
            <div style="display:flex;gap:6px"><el-input-number v-model="filters.score_min" :min="0" :max="100" controls-position="right" style="flex:1" />
            <el-input-number v-model="filters.score_max" :min="0" :max="100" controls-position="right" style="flex:1" /></div>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="91推送">
            <el-select v-model="filters.pushed_91" clearable style="width:100%">
              <el-option label="已推送" value="1" /><el-option label="未推送" value="0" />
            </el-select></el-form-item></el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="resetFilters">清空</el-button>
        <el-button type="primary" @click="applyFilters">应用筛选</el-button>
      </template>
    </el-dialog>

    <!-- ===== 新建商机 ===== -->
    <el-dialog v-model="createVisible" title="新建商机" width="600px">
      <el-form :model="form" label-width="100px" size="small">
        <el-row :gutter="12">
          <el-col :span="12"><el-form-item label="姓名"><el-input v-model="form.name" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="手机"><el-input v-model="form.phone" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="在读年级">
            <el-select v-model="form.grade" style="width:100%"><el-option v-for="g in grades" :key="g" :label="g" :value="g" /></el-select>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="学校"><el-input v-model="form.school" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="意向等级">
            <el-radio-group v-model="form.level"><el-radio value="高">高</el-radio><el-radio value="中">中</el-radio><el-radio value="低">低</el-radio></el-radio-group>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="接触点">
            <el-radio-group v-model="form.contact_point"><el-radio value="呼入">呼入</el-radio><el-radio value="呼出">呼出</el-radio></el-radio-group>
          </el-form-item></el-col>
          <el-col :span="12"><el-form-item label="意向项目">
            <el-select v-model="form.intention_project" style="width:100%">
              <el-option v-for="p in ['一对一','小班课','大班课','集训营']" :key="p" :label="p" :value="p" />
            </el-select></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="所属人">
            <el-select v-model="form.owner_name" filterable style="width:100%">
              <el-option v-for="e in employees" :key="e.id" :label="e.name" :value="e.name" />
            </el-select></el-form-item></el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>

    <!-- ===== 跟进计划 ===== -->
    <el-dialog v-model="planVisible" title="跟进计划" width="420px">
      <el-form label-width="100px" size="small">
        <el-form-item label="已选商机">{{ selection.length }} 条</el-form-item>
        <el-form-item label="下次回访"><el-date-picker v-model="planAt" type="datetime" style="width:100%" /></el-form-item>
        <el-form-item label="跟进备注"><el-input v-model="planContent" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="planVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPlan">保存</el-button>
      </template>
    </el-dialog>

    <!-- ===== 申请成单 ===== -->
    <el-dialog v-model="dealVisible" title="申请成单" width="480px">
      <el-form :model="dealForm" label-width="100px" size="small">
        <el-form-item label="业务类型">
          <el-select v-model="dealForm.biz_type" style="width:100%">
            <el-option v-for="b in ['班课','高端','会员卡','课时包']" :key="b" :label="b" :value="b" />
          </el-select>
        </el-form-item>
        <el-form-item label="订单金额"><el-input v-model="dealForm.amount" placeholder="单位：元" /></el-form-item>
        <el-form-item label="订单编号"><el-input v-model="dealForm.cert_no" /></el-form-item>
        <el-form-item label="申请说明"><el-input v-model="dealForm.reason" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dealVisible = false">取消</el-button>
        <el-button type="primary" @click="submitDeal">提交审批</el-button>
      </template>
    </el-dialog>

    <!-- ===== 分配 / 回收 / 归档 ===== -->
    <el-dialog v-model="assignVisible" title="分配商机" width="420px">
      <el-select v-model="receiver" placeholder="选择接收员工" filterable style="width:100%">
        <el-option v-for="e in employees" :key="e.id" :label="e.name" :value="e.name" />
      </el-select>
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAssign">确认</el-button>
      </template>
    </el-dialog>
    <el-dialog v-model="recycleVisible" title="商机回收" width="420px">
      <el-select v-model="recycleReason" placeholder="回收原因" style="width:100%">
        <el-option v-for="r in recycleReasons" :key="r" :label="r" :value="r" />
      </el-select>
      <template #footer>
        <el-button @click="recycleVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRecycle">确认回收</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, reactive, onMounted, h } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { scrmApi, maskPhone, fmtTime, type Opportunity, type Employee } from '../api/scrm'

const route = useRoute()
const isSuper = computed(() => route.query.mode !== 'employee')

const tabs = computed(() => isSuper.value
  ? [
      { v: 'active', t: '活跃商机' }, { v: 'all', t: '全部商机' },
      { v: 'recycled', t: '已回收' }, { v: 'archived', t: '已归档' },
      { v: 'tasks', t: '任务列表' }, { v: 'empTasks', t: '员工任务' }
    ]
  : [
      { v: 'active', t: '活跃商机' }, { v: 'all', t: '全部商机' },
      { v: 'assigned', t: '已分配商机' }, { v: 'travel', t: '国内国外游商机' }
    ])
const tab = ref('active')
const oppTabs = ['active', 'all', 'recycled', 'archived', 'travel']
const isOppTab = computed(() => oppTabs.includes(tab.value))
const oppStatuses = ['待分配', '待跟进', '跟进中', '已回收', '死单归档', '已成单']
const grades = ['高三', '高二', '高一', '初三', '初二']
const dateChips = [
  { v: 'today', t: '今天' }, { v: '7d', t: '近7天' }, { v: '30d', t: '近30天' },
  { v: '3m', t: '近三个月' }, { v: '6m', t: '近半年' }, { v: 'all', t: '全部' }
]
const recycleReasons = ['30天未跟进自动回收', '主管手动回收-质量差', '撞单回收', '长期无意向回收']
const timeCols = ['created_at', 'last_follow_at', 'next_follow_at', 'assign_at', 'promise_visit_at',
  'planned_recycle_at', 'recycle_at', 'archive_at', 'deal_at', 'push_91_at']

const keyword = ref('')
const dateChip = ref('6m')
const quick = ref('pending')
const rows = ref<Opportunity[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const selection = ref<Opportunity[]>([])
const employees = ref<Employee[]>([])
const scoreTip = ref<Record<string, string>>({})

const filters = reactive<any>({
  level: '', grade2: '', grade: '', new_old: '', contact_point: '', enter_type: '',
  intention_dept: '', intention_project: '', owner_name: '', school: '',
  score_min: undefined, score_max: undefined, pushed_91: ''
})

function dateRange(chip: string): [string, string] | [] {
  if (chip === 'all') return []
  const days: Record<string, number> = { today: 0, '7d': 7, '30d': 30, '3m': 90, '6m': 180 }
  const d = new Date(); const end = d.toISOString().slice(0, 10)
  d.setDate(d.getDate() - (days[chip] ?? 180))
  return [d.toISOString().slice(0, 10), end]
}

async function loadOpps() {
  loading.value = true
  try {
    if (tab.value === 'travel') { await loadTravel(); return }
    const params: Record<string, any> = {
      page: page.value, page_size: 20, scope: isSuper.value ? 'dept' : 'mine',
      keyword: keyword.value, quick: quick.value || '',
      level: filters.level, grade2: filters.grade2, grade: filters.grade,
      new_old: filters.new_old, contact_point: filters.contact_point, enter_type: filters.enter_type,
      intention_dept: filters.intention_dept, intention_project: filters.intention_project,
      owner_name: filters.owner_name, school: filters.school,
      score_min: filters.score_min || 0, score_max: filters.score_max || 0, pushed_91: filters.pushed_91
    }
    const rg = dateRange(dateChip.value)
    if (rg.length) { params.created_start = rg[0]; params.created_end = rg[1] }
    if (tab.value === 'active' && !quick.value) params.statuses = '待跟进,跟进中'
    if (tab.value === 'recycled') { params.status = '已回收'; params.quick = '' }
    if (tab.value === 'archived') { params.status = '死单归档'; params.quick = '' }
    const res = await scrmApi.opportunities(params)
    rows.value = res.items; total.value = res.total
  } finally { loading.value = false }
}

async function loadTravel() {
  const [a, b] = await Promise.all([
    scrmApi.opportunities({ scope: 'dept', category: '国际游学', page: 1, page_size: 100 }),
    scrmApi.opportunities({ scope: 'dept', category: '国内游学', page: 1, page_size: 100 })
  ])
  rows.value = [...a.items, ...b.items]
  total.value = a.total + b.total
  loading.value = false
}

function setQuick(q: string) { quick.value = quick.value === q ? '' : q; page.value = 1; loadOpps() }
function resetQuery() {
  keyword.value = ''; quick.value = ''; dateChip.value = '6m'; page.value = 1
  Object.assign(filters, { level: '', grade2: '', grade: '', new_old: '', contact_point: '', enter_type: '', intention_dept: '', intention_project: '', owner_name: '', school: '', score_min: undefined, score_max: undefined, pushed_91: '' })
  loadOpps()
}
function applyFilters() { filterVisible.value = false; page.value = 1; loadOpps() }
function resetFilters() {
  Object.keys(filters).forEach(k => { (filters as any)[k] = k.startsWith('score') ? undefined : '' })
}

function oppStatusTag(s: string) {
  return ({ 跟进中: 'success', 待跟进: 'warning', 已回收: 'info', 已成单: 'danger' } as Record<string, string>)[s] || 'info'
}
function gradeTag(g: string) { return { S: 'danger', A: 'warning', B: 'success', C: 'info' }[g] || 'info' }

// ---- 列定义与列配置（后端持久化） ----
interface Col { p: string; label: string; w?: number; mw?: number; pin?: boolean; show?: boolean; score?: string }
const defaultCols: Col[] = [
  { p: 'code', label: '商机编码', w: 150 },
  { p: 'name', label: '姓名', w: 90 },
  { p: 'phone', label: '手机', w: 120 },
  { p: 'student_no', label: '学员号', w: 120 },
  { p: 'new_old_type', label: '新老生', w: 80 },
  { p: 'grade', label: '在读年级', w: 90 },
  { p: 'school', label: '学校', mw: 120 },
  { p: 'district', label: '区县', w: 90 },
  { p: 'status', label: '商机状态', w: 100 },
  { p: 'enter_type', label: '进入类型', w: 120 },
  { p: 'create_method', label: '创建方式', w: 110 },
  { p: 'contact_point', label: '接触点', w: 80 },
  { p: 'level', label: '意向等级', w: 90 },
  { p: 'grade2', label: '商机等级', w: 90, score: 'level' },
  { p: 'score', label: '商机分数', w: 90, score: 'score' },
  { p: 'init_score', label: '初始分数', w: 90, score: 'init' },
  { p: 'max_score', label: '最高分', w: 80, score: 'max' },
  { p: 'ai_diagnosis', label: 'AI诊断', mw: 160, score: 'ai' },
  { p: 'intention_dept', label: '意向部门', w: 100 },
  { p: 'intention_project', label: '意向项目', w: 100 },
  { p: 'intention_campus', label: '意向校区', w: 110 },
  { p: 'intention_grade', label: '意向年级', w: 90 },
  { p: 'intention_subject', label: '意向科目', w: 90 },
  { p: 'channel_l1', label: '数据渠道', w: 100 },
  { p: 'business_channel', label: '业绩渠道', w: 100 },
  { p: 'system_source', label: '系统来源', w: 100 },
  { p: 'activity_name', label: '关联活动', mw: 160 },
  { p: 'lead_code', label: '线索编码', w: 150 },
  { p: 'owner_name', label: '当前所属人', w: 100 },
  { p: 'owner_dept', label: '所属部门', w: 140 },
  { p: 'assign_at', label: '分配时间', w: 160 },
  { p: 'follow_count', label: '跟进次数', w: 90 },
  { p: 'last_follow_at', label: '最近跟进时间', w: 150 },
  { p: 'last_follow_content', label: '最近跟进内容', mw: 170 },
  { p: 'last_follow_by', label: '跟进人', w: 90 },
  { p: 'next_follow_at', label: '下次回访', w: 150 },
  { p: 'promise_visit_at', label: '诺访时间', w: 150 },
  { p: 'add_wechat', label: '是否加微', w: 90 },
  { p: 'planned_recycle_at', label: '计划回收时间', w: 150 },
  { p: 'recycle_reason', label: '回收原因', mw: 140 },
  { p: 'archive_at', label: '归档时间', w: 150 },
  { p: 'archive_rule', label: '归档原因', mw: 140 },
  { p: 'pushed_91', label: '91推送状态', w: 100 },
  { p: 'resource_no_91', label: '91资源编号', w: 150 },
  { p: 'push_91_at', label: '推送91时间', w: 150 },
  { p: 'follow_91', label: '91跟进人', w: 100 },
  { p: 'advancer_91', label: '91推进人', w: 100 },
  { p: 'pushed_hedgehog', label: '刺猬推送', w: 90 },
  { p: 'deal_at', label: '成单时间', w: 150 },
  { p: 'creator', label: '创建人', w: 90 },
  { p: 'created_at', label: '创建时间', w: 150 },
  { p: 'remark', label: '备注', mw: 140 }
]
const colCfg = ref<Col[]>(defaultCols.map(c => ({ ...c, show: true })))
const colVisible = ref(false)
const filterVisible = ref(false)
const PAGE_KEY = 'opp-workbench'

const shownColumns = computed(() => colCfg.value.filter(c => c.show !== false))

async function loadCols() {
  try {
    const raw = await scrmApi.getGridConfig(PAGE_KEY)
    if (raw) {
      const saved = JSON.parse(raw) as { p: string; show: boolean; pin: boolean }[]
      const map = new Map(saved.map(s => [s.p, s]))
      colCfg.value = saved.map(s => {
        const d = defaultCols.find(x => x.p === s.p) || { p: s.p, label: s.p }
        return { ...d, show: s.show, pin: s.pin }
      })
      defaultCols.forEach(d => { if (!map.has(d.p)) colCfg.value.push({ ...d, show: true }) })
    }
  } catch { /* 无配置时用默认 */ }
}
async function saveCols() {
  const payload = colCfg.value.map(c => ({ p: c.p, show: c.show !== false, pin: !!c.pin }))
  try { await scrmApi.saveGridConfig(PAGE_KEY, 'column', JSON.stringify(payload)) } catch { /* ignore */ }
  colVisible.value = false
  ElMessage.success('列配置已保存')
}
function resetCols() { colCfg.value = defaultCols.map(c => ({ ...c, show: true })) }

// 简易拖拽排序（不引第三方库）
const DraggableList = (props: any) => h('div', props.modelValue.map((c: Col, idx: number) =>
  h('div', {
    class: 'colcfg-row', draggable: 'true', key: c.p,
    onDragstart: (): any => dragIdx.value = idx,
    onDragover: (e: DragEvent) => e.preventDefault(),
    onDrop: () => {
      const arr = props.modelValue as Col[]
      const [item] = arr.splice(dragIdx.value, 1)
      arr.splice(idx, 0, item)
    }
  }, [
    h('span', { class: 'drag-handle' }, '≡'),
    h('el-checkbox', { modelValue: c.show !== false, 'onUpdate:modelValue': (v: boolean) => { c.show = v } }, () => ''),
    h('span', { class: 'colcfg-label' }, c.label + (c.score ? ' ②' : '')),
    h('el-button', {
      size: 'small', text: true, type: c.pin ? 'primary' : '',
      onClick: () => { c.pin = !c.pin }
    }, () => '📌')
  ])
))
const dragIdx = ref(-1)

// ---- 新建 ----
const createVisible = ref(false)
const form = reactive<any>({ name: '', phone: '', grade: '高三', school: '', level: '中', contact_point: '呼出', intention_project: '一对一', owner_name: '' })
function openCreate() { createVisible.value = true }
async function submitCreate() {
  if (!form.name) return ElMessage.warning('请填写姓名')
  await scrmApi.createOpportunity({ ...form, owner_id: '1' })
  ElMessage.success('商机已创建')
  createVisible.value = false
  loadOpps()
}

// ---- 跟进计划 ----
const planVisible = ref(false)
const planAt = ref<any>('')
const planContent = ref('')
function openFollowPlan() {
  if (!selection.value.length) return ElMessage.warning('请先选择商机')
  planVisible.value = true
}
async function submitPlan() {
  if (!planAt.value) return ElMessage.warning('请选择回访时间')
  for (const o of selection.value) {
    await scrmApi.updateOpportunity(o.id, {
      next_follow_at: new Date(planAt.value).toISOString(),
      last_follow_content: planContent.value
    })
  }
  ElMessage.success('跟进计划已保存')
  planVisible.value = false
  loadOpps()
}

// ---- 申请成单 ----
const dealVisible = ref(false)
const dealForm = reactive<any>({ biz_type: '班课', amount: '', cert_no: '', reason: '' })
function openApplyDeal() {
  if (!selection.value.length) return ElMessage.warning('请先选择商机')
  dealVisible.value = true
}
async function submitDeal() {
  if (!dealForm.amount) return ElMessage.warning('请填写订单金额')
  await Promise.all(selection.value.map(o => scrmApi.applyDeal({
    opp_code: o.code, biz_type: dealForm.biz_type, amount: dealForm.amount,
    cert_no: dealForm.cert_no, reason: dealForm.reason, order_creator: '陈元素'
  })))
  ElMessage.success('已提交审核中心审批')
  dealVisible.value = false
}

// ---- 分配 / 回收 / 归档 / 激活 / 推送 ----
const assignVisible = ref(false)
const receiver = ref('')
function onAssign() { if (!selection.value.length) return ElMessage.warning('请先选择商机'); assignVisible.value = true }
async function submitAssign() {
  if (!receiver.value) return ElMessage.warning('请选择接收员工')
  await scrmApi.assignOpportunities({ ids: selection.value.map(o => o.id), receiver_id: '1', receiver_name: receiver.value })
  ElMessage.success('分配成功')
  assignVisible.value = false
  loadOpps()
}
const recycleVisible = ref(false)
const recycleReason = ref('')
async function submitRecycle() {
  if (!recycleReason.value) return ElMessage.warning('请选择回收原因')
  await scrmApi.recycleOpportunities(selection.value.map(o => o.id), recycleReason.value)
  ElMessage.success('已回收'); recycleVisible.value = false; loadOpps()
}
async function onArchive() {
  if (!selection.value.length) return ElMessage.warning('请先选择商机')
  await scrmApi.archiveOpportunities(selection.value.map(o => o.id), '死单归档-主管批量归档')
  ElMessage.success('已归档'); loadOpps()
}
async function onActivate() {
  if (!selection.value.length) return ElMessage.warning('请先选择商机')
  await ElMessageBox.confirm(`确认激活 ${selection.value.length} 条商机？`, '提示', { type: 'warning' })
  await scrmApi.activateOpportunities(selection.value.map(o => o.id))
  ElMessage.success('已激活'); loadOpps()
}
async function onPush91() {
  if (!selection.value.length) return ElMessage.warning('请先选择可连线商机')
  const r = await scrmApi.push91(selection.value.map(o => o.id))
  ElMessage.success(`已推送 ${r.pushed} 条至91国游，资源编号已回捞`)
  loadOpps()
}
async function onPushHedgehog() {
  if (!selection.value.length) return ElMessage.warning('请先选择商机')
  const r = await scrmApi.pushHedgehog(selection.value.map(o => o.id))
  ElMessage.success(`已推送 ${r.pushed} 条至刺猬`)
  loadOpps()
}
function onExport() {
  ElMessage.info('已按当前筛选条件导出商机列表（演示）')
}

// ---- 已分配商机流水 ----
const assignRows = ref<any[]>([])
const assignTotal = ref(0)
const assignPage = ref(1)
const assignKw = ref('')
const assignAuto = ref('')
async function loadAssigns() {
  loading.value = true
  try {
    const res = await scrmApi.assignRecords({
      scope: 'all', page: assignPage.value, page_size: 20,
      keyword: assignKw.value, auto: assignAuto.value
    })
    assignRows.value = res.items; assignTotal.value = res.total
  } finally { loading.value = false }
}

// ---- 任务 ----
const empTasks = ref<any[]>([])
const taskAgg = ref<any[]>([])
function parseCounts(s: string): Record<string, number> { try { return JSON.parse(s || '{}') } catch { return {} } }
async function loadTasks() {
  loading.value = true
  try {
    const res = await scrmApi.employeeTasks({ biz_kind: 'opp', page: 1, page_size: 200 })
    empTasks.value = res.items
    const byTask = new Map<string, any>()
    for (const t of res.items) {
      const cur = byTask.get(t.task_name) || { task_name: t.task_name, task_type: t.task_type, total_count: 0, status_counts: '{}' }
      cur.total_count += t.total_count
      const mc = parseCounts(cur.status_counts); const c = parseCounts(t.status_counts)
      for (const k of Object.keys(c)) mc[k] = (mc[k] || 0) + c[k]
      cur.status_counts = JSON.stringify(mc)
      byTask.set(t.task_name, cur)
    }
    taskAgg.value = [...byTask.values()]
  } finally { loading.value = false }
}

function onTabChange() {
  page.value = 1; selection.value = []
  if (isOppTab.value) { quick.value = tab.value === 'active' ? 'pending' : ''; loadOpps() }
  else if (tab.value === 'assigned') loadAssigns()
  else loadTasks()
}

onMounted(async () => {
  try { employees.value = await scrmApi.employees() } catch { /* ignore */ }
  try { scoreTip.value = await scrmApi.scoreRules() } catch { /* ignore */ }
  await loadCols()
  loadOpps()
})
</script>

<style scoped>
.wb-page { padding: 12px 16px; }
.wb-head { display: flex; align-items: center; justify-content: space-between; }
.wb-title { font-size: 16px; font-weight: 600; }
.toolbar { display: flex; gap: 8px; align-items: center; margin: 8px 0; flex-wrap: wrap; }
.tb-spacer { flex: 1; }
.pager { margin-top: 10px; justify-content: flex-end; }
.ii { color: #409eff; cursor: help; font-size: 12px; }
.diag { color: #e6a23c; cursor: help; }
.colcfg-tip { color: #909399; font-size: 12px; margin-bottom: 8px; }
:deep(.colcfg-row) { display: flex; align-items: center; gap: 8px; padding: 6px 8px; border: 1px solid #ebeef5; border-radius: 4px; margin-bottom: 6px; background: #fafafa; }
:deep(.drag-handle) { cursor: move; color: #909399; width: 16px; }
:deep(.colcfg-label) { flex: 1; }
</style>
