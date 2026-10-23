<template>
  <div class="page">
    <div class="chip-bar">
      <el-tag closable type="info" effect="plain">创建时间：2026-08-14 00:00:00 ~ 2026-10-14 23:59:59</el-tag>
    </div>

    <el-card shadow="never">
      <!-- 筛选区 -->
      <el-form :model="f" inline class="filter-form">
        <el-form-item label="关键词">
          <el-input v-model="f.keyword" placeholder="姓名/手机/学员号" clearable style="width:180px" @keyup.enter="search" />
        </el-form-item>
        <el-form-item v-if="mode === 'all'" label="线索状态">
          <el-select v-model="f.status" clearable placeholder="全部状态" style="width:120px">
            <el-option v-for="o in opts.lead_status || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="活动名称">
          <el-input v-model="f.activity_name" clearable style="width:170px" @keyup.enter="search" />
        </el-form-item>
        <el-form-item label="学员标签">
          <el-select v-model="f.student_tag" clearable style="width:110px">
            <el-option v-for="o in opts.student_tag || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker v-model="createdRange" type="daterange" value-format="YYYY-MM-DD"
            start-placeholder="开始" end-placeholder="结束" style="width:240px" />
        </el-form-item>

        <template v-if="expanded">
          <el-form-item label="姓名"><el-input v-model="f.name" style="width:120px" /></el-form-item>
          <el-form-item label="手机号"><el-input v-model="f.phone" style="width:130px" /></el-form-item>
          <el-form-item label="学员号"><el-input v-model="f.student_no" style="width:130px" /></el-form-item>
          <el-form-item label="新老生">
            <el-select v-model="f.new_old" clearable style="width:90px">
              <el-option v-for="o in opts.new_old || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="在读年级">
            <el-select v-model="f.grade" clearable style="width:90px">
              <el-option v-for="o in opts.grade || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="意向年级">
            <el-select v-model="f.intention_grade" clearable style="width:90px">
              <el-option v-for="o in opts.grade || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="区县">
            <el-select v-model="f.district" clearable style="width:100px">
              <el-option v-for="o in opts.district || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="所在学校"><el-input v-model="f.school" style="width:140px" /></el-form-item>
          <el-form-item label="线索类型">
            <el-select v-model="f.lead_type" clearable style="width:110px">
              <el-option v-for="o in opts.lead_type || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="意向等级">
            <el-select v-model="f.intention_level" clearable style="width:90px">
              <el-option v-for="o in opts.intention_level || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="一级渠道">
            <el-select v-model="f.channel_l1" clearable style="width:120px">
              <el-option v-for="o in channelL1" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="二级渠道">
            <el-select v-model="f.channel_l2" clearable style="width:120px">
              <el-option v-for="o in channelL2" :key="o" :label="o" :value="o" />
            </el-select>
          </el-form-item>
          <el-form-item label="三级渠道">
            <el-select v-model="f.channel_l3" clearable style="width:120px">
              <el-option v-for="o in channelL3" :key="o" :label="o" :value="o" />
            </el-select>
          </el-form-item>
          <el-form-item label="意向部门">
            <el-select v-model="f.intention_dept" clearable style="width:110px">
              <el-option v-for="o in opts.intention_dept || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="意向项目">
            <el-select v-model="f.intention_project" clearable style="width:110px">
              <el-option v-for="o in opts.intention_project || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="意向校区">
            <el-select v-model="f.intention_campus" clearable style="width:110px">
              <el-option v-for="o in opts.intention_campus || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="归属人"><el-input v-model="f.owner_name" style="width:120px" /></el-form-item>
          <el-form-item label="创建人"><el-input v-model="f.creator" style="width:110px" /></el-form-item>
          <el-form-item label="跟进次数">
            <el-input-number v-model="f.follow_min" :min="0" controls-position="right" size="small" style="width:90px" /> ~
            <el-input-number v-model="f.follow_max" :min="0" controls-position="right" size="small" style="width:90px" />
          </el-form-item>
          <el-form-item v-if="mode === 'all'" label="交接时间">
            <el-date-picker v-model="handoverRange" type="daterange" value-format="YYYY-MM-DD"
              start-placeholder="开始" end-placeholder="结束" style="width:230px" />
          </el-form-item>
          <el-form-item v-if="mode === 'all'" label="归档方式">
            <el-select v-model="f.archive_way" clearable style="width:110px">
              <el-option v-for="o in opts.archive_way || []" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="mode === 'all'" label="备注"><el-input v-model="f.remark" style="width:140px" /></el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="reset">重置</el-button>
          <el-button link type="primary" @click="expanded = !expanded">
            {{ expanded ? '收起' : '展开' }}<el-icon><ArrowUp v-if="expanded" /><ArrowDown v-else /></el-icon>
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 操作区 -->
      <div class="op-bar">
        <div class="op-left">
          <el-button @click="exportXlsx"><el-icon><Download /></el-icon>导出</el-button>
          <el-button @click="adjustVisible = true" :disabled="!ids.length">批量调整信息</el-button>
          <el-button type="primary" plain @click="handoverVisible = true" :disabled="!ids.length">线索交接</el-button>
          <template v-if="mode === 'all'">
            <el-button @click="archiveVisible = true" :disabled="!ids.length">批量归档</el-button>
            <el-button type="success" plain @click="activate" :disabled="!ids.length">激活</el-button>
          </template>
          <el-button @click="showStat">统计</el-button>
        </div>
        <div class="op-right">
          <span class="sel-tip" v-if="ids.length">已选 {{ ids.length }} 条</span>
          <el-button text @click="load"><el-icon><Refresh /></el-icon>刷新</el-button>
        </div>
      </div>

      <!-- 待交接线索表格 -->
      <el-table v-if="mode === 'handover'" :data="rows" border size="small" height="54vh"
        @selection-change="(s: any[]) => ids = s.map(x => x.id)">
        <el-table-column type="selection" width="44" />
        <el-table-column prop="name" label="姓名" width="90">
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="openDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column label="手机" width="130">
          <template #default="{ row }">{{ maskPhone(row.phone) }}</template>
        </el-table-column>
        <el-table-column prop="student_tag" label="学员标签" width="90" />
        <el-table-column prop="student_no" label="学员号" width="130" />
        <el-table-column prop="school" label="所在学校" min-width="130" show-overflow-tooltip />
        <el-table-column prop="grade" label="在读年级" width="80" />
        <el-table-column prop="district" label="区县" width="80" />
        <el-table-column prop="code" label="线索编码" min-width="160" show-overflow-tooltip />
        <el-table-column prop="lead_type" label="线索类型" width="100" />
        <el-table-column prop="activity_name" label="活动名称" min-width="170" show-overflow-tooltip />
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="intention_grade" label="意向年级" width="80" />
        <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
      </el-table>

      <!-- 全部线索表格 -->
      <el-table v-else :data="rows" border size="small" height="54vh"
        @selection-change="(s: any[]) => ids = s.map(x => x.id)">
        <el-table-column type="selection" width="44" />
        <el-table-column prop="name" label="姓名" width="90" fixed>
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="openDetail(row)">{{ row.name }}</el-link>
            <el-tag v-if="row.is_duplicate" size="small" type="warning" effect="plain" style="margin-left:4px">重复</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="手机号" width="130">
          <template #default="{ row }">{{ maskPhone(row.phone) }}</template>
        </el-table-column>
        <el-table-column prop="student_tag" label="学员标签" width="90" />
        <el-table-column prop="student_no" label="学员号" width="130" />
        <el-table-column prop="code" label="线索编码" min-width="160" show-overflow-tooltip />
        <el-table-column label="线索状态" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="statusType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进入类型" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="row.enter_type === '激活线索' ? 'warning' : 'success'" effect="plain">{{ row.enter_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="follow_count" label="跟进次数" width="90" sortable :sort-by="'follow_count'" align="center" />
        <el-table-column prop="intention_level" label="意向等级" width="90" align="center" />
        <el-table-column label="最近跟进时间" width="160">
          <template #default="{ row }">{{ fmtTime(row.last_follow_at) }}</template>
        </el-table-column>
        <el-table-column prop="last_follow_content" label="最近跟进内容" min-width="180" show-overflow-tooltip />
        <el-table-column label="交接时间" width="160">
          <template #default="{ row }">{{ fmtTime(row.handover_time) }}</template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pager"
        layout="total, prev, pager, next, sizes"
        :total="total" :page-size="pageSize" :current-page="page" :page-sizes="[20, 50, 100]"
        @current-change="(p: number) => { page = p; load() }"
        @size-change="(s: number) => { pageSize = s; page = 1; load() }"
      />
    </el-card>

    <!-- 详情抽屉 -->
    <el-drawer v-model="detailVisible" title="线索详情" size="460px">
      <template v-if="detail">
        <div class="detail-ops">
          <el-button type="success" size="small" @click="openCall(detail)">
            <el-icon><Phone /></el-icon>&nbsp;云座席双呼
          </el-button>
          <el-tag size="small" :type="detail.call_status === '已接通' ? 'success' : detail.call_status === '未接通' ? 'info' : 'warning'" effect="plain">
            接线状态：{{ detail.call_status || '未拨打' }}
          </el-tag>
        </div>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="姓名">{{ detail.name }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ detail.phone }}（列表页默认脱敏）</el-descriptions-item>
          <el-descriptions-item label="学员号">{{ detail.student_no }}</el-descriptions-item>
          <el-descriptions-item label="线索编码">{{ detail.code }}</el-descriptions-item>
          <el-descriptions-item label="线索状态">{{ detail.status }}</el-descriptions-item>
          <el-descriptions-item label="进入类型">{{ detail.enter_type }}</el-descriptions-item>
          <el-descriptions-item label="学员标签">{{ detail.student_tag }}</el-descriptions-item>
          <el-descriptions-item label="学校 / 年级">{{ detail.school }} / {{ detail.grade }}</el-descriptions-item>
          <el-descriptions-item label="意向">{{ detail.intention_grade }} · {{ detail.intention_level }} · {{ detail.intention_campus }} · {{ detail.intention_project }}</el-descriptions-item>
          <el-descriptions-item label="渠道">{{ [detail.channel_l1, detail.channel_l2, detail.channel_l3].filter(Boolean).join(' / ') }}</el-descriptions-item>
          <el-descriptions-item label="所属活动">{{ detail.activity_name }}</el-descriptions-item>
          <el-descriptions-item label="归属人">{{ detail.owner_name || '未分配' }}</el-descriptions-item>
          <el-descriptions-item label="跟进次数">{{ detail.follow_count }}</el-descriptions-item>
          <el-descriptions-item label="外呼记录">自外呼 {{ detail.self_call_count || 0 }} 次 / 总计 {{ detail.total_call_count || 0 }} 次；最近：{{ fmtTime(detail.last_call_at || '') || '—' }}</el-descriptions-item>
          <el-descriptions-item label="最近跟进">{{ fmtTime(detail.last_follow_at) }} {{ detail.last_follow_content }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ fmtTime(detail.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="备注">{{ detail.remark || '—' }}</el-descriptions-item>
        </el-descriptions>

        <div class="logs-title">跟进日志（含外呼/通话记录）</div>
        <div v-if="!followLogs.length" class="logs-empty">暂无跟进记录</div>
        <el-timeline v-else class="logs-tl">
          <el-timeline-item
            v-for="lg in followLogs"
            :key="lg.id"
            :type="lg.result === '已接通' ? 'success' : lg.result === '未接通' ? 'info' : 'primary'"
            :timestamp="fmtTime(lg.created_at)"
          >
            <div class="lg-line">
              <el-tag size="small" effect="plain">{{ lg.type }}</el-tag>
              <el-tag size="small" :type="lg.result === '已接通' ? 'success' : 'info'" effect="plain">{{ lg.result || '事件' }}</el-tag>
              <span v-if="lg.duration_sec" class="lg-dur">{{ mmss(lg.duration_sec) }}</span>
            </div>
            <div class="lg-content">{{ lg.content }}</div>
            <div class="lg-operator">{{ lg.operator }}</div>
          </el-timeline-item>
        </el-timeline>
      </template>
    </el-drawer>

    <!-- 批量调整 -->
    <el-dialog v-model="adjustVisible" title="批量调整信息" width="480px">
      <el-form label-width="100px">
        <el-form-item label="学员标签">
          <el-select v-model="adjust.student_tag" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in opts.student_tag || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="意向年级">
          <el-select v-model="adjust.intention_grade" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in opts.grade || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="意向部门">
          <el-select v-model="adjust.intention_dept" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in opts.intention_dept || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="意向项目">
          <el-select v-model="adjust.intention_project" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in opts.intention_project || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="一级渠道">
          <el-select v-model="adjust.channel_l1" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in channelL1" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="adjust.remark" placeholder="不修改则留空" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdjust">确认调整（{{ ids.length }} 条）</el-button>
      </template>
    </el-dialog>

    <!-- 批量归档 -->
    <el-dialog v-model="archiveVisible" title="批量归档" width="480px">
      <el-form label-width="100px">
        <el-form-item label="归档方式" required>
          <el-radio-group v-model="archiveForm.way">
            <el-radio value="无效归档">无效归档</el-radio>
            <el-radio value="入库归档">入库归档</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="归档规则" required>
          <el-select v-model="archiveForm.rule" placeholder="请选择归档规则" style="width:100%">
            <el-option v-for="o in opts.archive_rule || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="原因说明">
          <el-input v-model="archiveForm.reason" type="textarea" :rows="3" placeholder="请填写归档原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="archiveVisible = false">取消</el-button>
        <el-button type="danger" @click="submitArchive">确认归档（{{ ids.length }} 条）</el-button>
      </template>
    </el-dialog>

    <!-- 统计 -->
    <el-dialog v-model="statVisible" title="当前筛选结果统计" width="460px">
      <el-statistic title="符合筛选条件的线索总数" :value="stat.total" />
      <el-divider />
      <div v-for="(c, k) in stat.status" :key="k" class="stat-row">
        <el-tag size="small" :type="statusType(String(k))">{{ k }}</el-tag>
        <span class="stat-num">{{ c }}</span>
      </div>
    </el-dialog>

    <ReceiverDialog v-model="handoverVisible" title="线索交接" :count="ids.length" @confirm="submitHandover" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Refresh, ArrowUp, ArrowDown, Phone } from '@element-plus/icons-vue'
import { scrmApi, maskPhone, fmtTime, type Lead, type FollowLog } from '../api/scrm'
import { openCallBar, callState } from '../store/call'
import ReceiverDialog from '../components/ReceiverDialog.vue'

const route = useRoute()
const mode = computed(() => (route.path === '/handover-leads' ? 'handover' : 'all'))

const opts = ref<Record<string, any[]>>({})
const channelL1 = ref<any[]>([])
const channelL2 = ref<string[]>([])
const channelL3 = ref<string[]>([])
const rows = ref<Lead[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const expanded = ref(false)
const ids = ref<number[]>([])
const createdRange = ref<string[]>([])
const handoverRange = ref<string[]>([])

const f = reactive<any>({
  keyword: '', status: '', activity_name: '', student_tag: '', name: '', phone: '',
  student_no: '', new_old: '', grade: '', intention_grade: '', district: '', school: '',
  lead_type: '', intention_level: '', channel_l1: '', channel_l2: '', channel_l3: '',
  intention_dept: '', intention_project: '', intention_campus: '', owner_name: '',
  creator: '', follow_min: 0, follow_max: 0, archive_way: '', remark: '',
  activity_id: 0, revisit: '', dup: ''
})

function buildParams() {
  const p: any = {
    page: page.value, page_size: pageSize.value,
    view: mode.value === 'handover' ? 'handover' : undefined,
    ...f
  }
  if (createdRange.value?.length === 2) {
    p.created_start = createdRange.value[0]
    p.created_end = createdRange.value[1]
  }
  if (handoverRange.value?.length === 2) {
    p.handover_start = handoverRange.value[0]
    p.handover_end = handoverRange.value[1]
  }
  // 去掉空值，避免 0 值干扰
  Object.keys(p).forEach(k => {
    if (p[k] === '' || p[k] === 0 || p[k] === undefined) delete p[k]
  })
  return p
}

async function load() {
  const res = await scrmApi.leads(buildParams())
  rows.value = res.items
  total.value = res.total
}
function search() { page.value = 1; load() }
function reset() {
  Object.keys(f).forEach(k => (f as any)[k] = '')
  f.follow_min = 0; f.follow_max = 0
  createdRange.value = []
  handoverRange.value = []
  search()
}

function statusType(s: string) {
  const m: Record<string, string> = {
    '待分配': 'info', '待跟进': 'warning', '跟进中': '', '待交接': 'danger',
    '已交接': 'success', '已转商': 'success', '无效归档': 'info', '入库归档': 'info', '已回收': 'info'
  }
  return m[s] || ''
}

// 导出：后端按当前筛选导出 xlsx
function exportXlsx() {
  const params: Record<string, string> = {}
  const p = buildParams()
  delete p.page; delete p.page_size
  Object.entries(p).forEach(([k, v]) => (params[k] = String(v)))
  window.open(scrmApi.exportLeadsUrl(params), '_blank')
}

// 详情
const detailVisible = ref(false)
const detail = ref<Lead | null>(null)
const followLogs = ref<FollowLog[]>([])

async function openDetail(row: Lead) {
  detail.value = row
  detailVisible.value = true
  followLogs.value = []
  try {
    // 拉取最新详情（含通话维度）与跟进日志
    const [fresh, logs] = await Promise.all([
      scrmApi.leadDetail(row.id),
      scrmApi.followLogs(row.id)
    ])
    detail.value = fresh
    followLogs.value = logs
  } catch { /* 用列表行数据兜底 */ }
}

function openCall(row: Lead) {
  openCallBar(detail.value && detail.value.id === row.id ? detail.value : row)
}

function mmss(sec: number) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

// 工作条关闭后刷新当前抽屉的通话记录
watch(() => callState.visible, v => {
  if (!v && detailVisible.value && detail.value) {
    openDetail(detail.value)
    load()
  }
})

// 批量调整
const adjustVisible = ref(false)
const adjust = reactive<any>({ student_tag: '', intention_grade: '', intention_dept: '', intention_project: '', channel_l1: '', remark: '' })
async function submitAdjust() {
  const fields: Record<string, string> = {}
  Object.entries(adjust).forEach(([k, v]) => { if (v) fields[k] = String(v) })
  if (Object.keys(fields).length === 0) return ElMessage.warning('请至少选择一个要调整的字段')
  await scrmApi.batchAdjustLeads(ids.value, fields)
  ElMessage.success(`已调整 ${ids.value.length} 条线索`)
  adjustVisible.value = false
  load()
}

// 交接
const handoverVisible = ref(false)
async function submitHandover(emp: any) {
  await scrmApi.handoverLeads(ids.value, emp.id, emp.name)
  ElMessage.success(`已交接给 ${emp.name}`)
  handoverVisible.value = false
  ids.value = []
  load()
}

// 归档 / 激活
const archiveVisible = ref(false)
const archiveForm = reactive({ way: '无效归档', rule: '', reason: '' })
async function submitArchive() {
  if (!archiveForm.rule) return ElMessage.warning('请选择归档规则')
  await scrmApi.batchArchiveLeads(ids.value, archiveForm.rule, archiveForm.way, archiveForm.reason)
  ElMessage.success(`已归档 ${ids.value.length} 条线索`)
  archiveVisible.value = false
  archiveForm.rule = ''; archiveForm.reason = ''
  load()
}
async function activate() {
  await ElMessageBox.confirm(`确认激活选中的 ${ids.value.length} 条线索？激活后状态回到待分配，进入类型标记为「激活线索」。`, '激活确认', { type: 'warning' })
  await scrmApi.activateLeads(ids.value)
  ElMessage.success('激活成功')
  load()
}

// 统计
const statVisible = ref(false)
const stat = ref<{ total: number; status: Record<string, number> }>({ total: 0, status: {} })
async function showStat() {
  const p = buildParams()
  delete p.page; delete p.page_size
  stat.value = await scrmApi.leadStat(p)
  statVisible.value = true
}

// 渠道联动
async function syncChannels() {
  const tree: any[] = await scrmApi.channels()
  channelL1.value = tree
}
function onL1(v: string) {
  channelL2.value = (channelL1.value.find((x: any) => x.value === v)?.children || []).map((x: any) => x.value)
  f.channel_l2 = ''; f.channel_l3 = ''
  channelL3.value = []
}
function onL2(v: string) {
  const l1 = channelL1.value.find((x: any) => x.value === f.channel_l1)
  channelL3.value = (l1?.children?.find((x: any) => x.value === v)?.children || []).map((x: any) => x.value)
  f.channel_l3 = ''
}

// 全局搜索深链：?lead_id=123 直接打开线索详情
async function maybeOpenDeepLink() {
  const id = Number(route.query.lead_id)
  if (!id) return
  try {
    const lead = await scrmApi.leadDetail(id)
    openDetail(lead)
  } catch { /* 线索不存在时忽略 */ }
}

onMounted(async () => {
  // 路由参数预填（待办卡片下钻 / 活动下钻）
  const q = route.query
  f.keyword = q.keyword || ''
  f.status = q.status || ''
  f.revisit = q.revisit || ''
  f.dup = q.dup || ''
  f.activity_id = Number(q.activity_id) || 0
  f.activity_name = q.activity_name || ''
  if (f.activity_id) {
    // activity_id 由 buildParams 直接携带
  }

  const kinds = ['lead_status', 'student_tag', 'new_old', 'grade', 'district', 'lead_type',
    'intention_level', 'intention_dept', 'intention_project', 'intention_campus', 'archive_rule', 'archive_way']
  await Promise.all(kinds.map(async k => { opts.value[k] = await scrmApi.dict(k) }))
  await syncChannels()
  load()
  maybeOpenDeepLink()
})

watch(() => f.channel_l1, onL1)
watch(() => f.channel_l2, onL2)
watch(() => route.fullPath, () => {
  // 从待交接 / 全部线索菜单切换或搜索深链跳转时重新按视图加载
  page.value = 1
  load()
  maybeOpenDeepLink()
})
</script>

<style scoped>
.chip-bar { margin-bottom: 10px; }
.filter-form :deep(.el-form-item) { margin-bottom: 8px; }
.op-bar { display: flex; justify-content: space-between; align-items: center; margin: 4px 0 10px; }
.op-right { display: flex; align-items: center; gap: 8px; }
.sel-tip { font-size: 12px; color: #00a86b; }
.pager { margin-top: 10px; justify-content: flex-end; }
.stat-row { display: flex; justify-content: space-between; align-items: center; padding: 6px 0; border-bottom: 1px dashed #f0f0f0; }
.stat-num { font-weight: 700; color: #00a86b; }
.detail-ops { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; }
.logs-title { font-weight: 700; font-size: 14px; margin: 18px 0 10px; color: #1f2d3d; }
.logs-empty { color: #909399; font-size: 13px; text-align: center; padding: 20px 0; }
.lg-line { display: flex; align-items: center; gap: 6px; }
.lg-dur { font-size: 12px; color: #00a86b; font-variant-numeric: tabular-nums; }
.lg-content { font-size: 13px; color: #475669; margin-top: 4px; }
.lg-operator { font-size: 11px; color: #c0c4cc; margin-top: 2px; }
</style>
