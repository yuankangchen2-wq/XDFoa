<template>
  <div class="page">
    <el-card shadow="never">
      <el-tabs v-model="tab" @tab-change="onTab">
        <!-- ============ Tab1：导入/激活日志 ============ -->
        <el-tab-pane label="导入/激活日志" name="logs">
          <el-form :model="lf" inline class="filter-form">
            <el-form-item label="活动名称"><el-input v-model="lf.activity_name" clearable style="width:140px" @keyup.enter="searchLogs" /></el-form-item>
            <el-form-item label="任务名称"><el-input v-model="lf.task_name" clearable style="width:140px" @keyup.enter="searchLogs" /></el-form-item>
            <el-form-item label="操作人"><el-input v-model="lf.operator" clearable style="width:110px" @keyup.enter="searchLogs" /></el-form-item>
            <el-form-item label="流水号"><el-input v-model="lf.serial_no" clearable style="width:130px" @keyup.enter="searchLogs" /></el-form-item>
            <el-form-item label="批次号"><el-input v-model="lf.batch_no" clearable style="width:130px" @keyup.enter="searchLogs" /></el-form-item>
            <el-form-item label="导入时间">
              <el-date-picker v-model="lf.range" type="daterange" range-separator="至"
                start-placeholder="开始" end-placeholder="结束" value-format="YYYY-MM-DD" style="width:230px" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="searchLogs">搜索</el-button>
              <el-button @click="resetLogs">重置</el-button>
            </el-form-item>
          </el-form>

          <div class="op-bar">
            <div>
              <el-button @click="exportLogs"><el-icon><Download /></el-icon>导出</el-button>
            </div>
            <el-button text @click="loadLogs"><el-icon><Refresh /></el-icon>刷新</el-button>
          </div>

          <el-table :data="logs" border size="small" height="52vh">
            <el-table-column type="index" label="序号" width="55" />
            <el-table-column prop="activity_name" label="活动名称" min-width="170" show-overflow-tooltip />
            <el-table-column prop="task_name" label="任务名称" min-width="150" show-overflow-tooltip />
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag size="small" :type="statusType(row.status)" effect="plain">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="total_count" label="导入条数" width="80" align="center" />
            <el-table-column prop="success_count" label="成功" width="70" align="center" />
            <el-table-column prop="new_count" label="纯新" width="70" align="center" />
            <el-table-column prop="active_count" label="激活" width="70" align="center" />
            <el-table-column label="失败条数" width="85" align="center">
              <template #default="{ row }">
                <el-link v-if="row.fail_count > 0" type="success" :underline="false" @click="downloadFail(row.id)">
                  {{ row.fail_count }} 下载
                </el-link>
                <span v-else>0</span>
              </template>
            </el-table-column>
            <el-table-column prop="duplicate_count" label="重复" width="70" align="center" />
            <el-table-column prop="invalid_count" label="无效" width="70" align="center" />
            <el-table-column prop="uniqueness_fail_count" label="不符合唯一性" width="100" align="center" />
            <el-table-column label="去重规则去重" width="110" align="center">
              <template #default="{ row }">
                <el-link v-if="row.dedup_hit_count > 0" type="primary" :underline="false" @click="showDedup(row)">
                  {{ row.dedup_hit_count }} · 详情
                </el-link>
                <span v-else>{{ row.dedup_hit_count }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="operator" label="操作人" width="90" />
            <el-table-column prop="serial_no" label="流水号" width="120" />
            <el-table-column prop="batch_no" label="批次号" width="120" />
            <el-table-column label="导入时间" width="160">
              <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>

          <el-pagination
            class="pager"
            layout="total, prev, pager, next, sizes"
            :total="logsTotal" :page-size="logsSize" :current-page="logsPage" :page-sizes="[20, 50, 100]"
            @current-change="(p: number) => { logsPage = p; loadLogs() }"
            @size-change="(s: number) => { logsSize = s; logsPage = 1; loadLogs() }"
          />
        </el-tab-pane>

        <!-- ============ Tab2：数据进入记录 ============ -->
        <el-tab-pane label="数据进入记录" name="records">
          <el-form :model="rf" inline class="filter-form">
            <el-form-item label="姓名"><el-input v-model="rf.name" clearable style="width:120px" @keyup.enter="searchRecords" /></el-form-item>
            <el-form-item label="手机号"><el-input v-model="rf.phone" clearable style="width:130px" @keyup.enter="searchRecords" /></el-form-item>
            <el-form-item label="年级">
              <el-select v-model="rf.grade" clearable placeholder="全部" style="width:110px">
                <el-option v-for="o in gradeOpts" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="学校/活动"><el-input v-model="rf.school" clearable style="width:140px" @keyup.enter="searchRecords" /></el-form-item>
            <el-form-item label="学号"><el-input v-model="rf.student_no" clearable style="width:130px" @keyup.enter="searchRecords" /></el-form-item>
            <el-form-item label="进入结果">
              <el-select v-model="rf.success" clearable placeholder="全部" style="width:120px">
                <el-option label="成功进入" value="success" />
                <el-option label="进入失败" value="fail" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="searchRecords">搜索</el-button>
              <el-button @click="resetRecords">重置</el-button>
            </el-form-item>
          </el-form>

          <div class="op-bar">
            <div>
              <el-button @click="exportRecords"><el-icon><Download /></el-icon>导出</el-button>
              <el-button @click="adjustVisible = true" :disabled="!recIds.length">批量调整信息</el-button>
              <el-button type="success" plain @click="receiverVisible = true" :disabled="!recIds.length">
                交接给校责任人
              </el-button>
            </div>
            <el-button text @click="loadRecords"><el-icon><Refresh /></el-icon>刷新</el-button>
          </div>

          <el-table :data="records" border size="small" height="52vh"
            @selection-change="(s: any[]) => recIds = s.map(x => x.id)">
            <el-table-column type="selection" width="44" fixed />
            <el-table-column prop="activity_name" label="活动名称" min-width="160" show-overflow-tooltip fixed />
            <el-table-column prop="related_name" label="关联学员名称" width="110" />
            <el-table-column label="关联手机号" width="130">
              <template #default="{ row }">{{ maskPhone(row.related_phone) }}</template>
            </el-table-column>
            <el-table-column prop="related_student_no" label="学员号" width="130" />
            <el-table-column label="年级" width="120">
              <template #default="{ row }">
                <div>活动：{{ row.activity_grade || '—' }}</div>
                <div>关联：{{ row.related_grade || '—' }}</div>
              </template>
            </el-table-column>
            <el-table-column prop="import_type" label="数据类型" width="90" />
            <el-table-column prop="enter_way" label="进入方式" width="110" />
            <el-table-column label="线索进入类型" width="130">
              <template #default="{ row }">
                <el-tag size="small" :type="row.lead_enter_type?.startsWith('成功') ? 'success' : 'danger'" effect="plain">
                  {{ row.lead_enter_type }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="渠道" min-width="150" show-overflow-tooltip>
              <template #default="{ row }">{{ [row.channel_l1, row.channel_l2, row.channel_l3].filter(Boolean).join(' / ') }}</template>
            </el-table-column>
            <el-table-column prop="intention_dept" label="意向部门" width="100" />
            <el-table-column prop="intention_project" label="意向项目" width="100" />
            <el-table-column label="归属/处理" width="120">
              <template #default="{ row }">
                <el-tag size="small" :type="row.handled ? 'success' : 'warning'" effect="plain">
                  {{ row.handled ? (row.owner_name || '已处理') : '待处理' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="fail_reason" label="失败原因/备注" min-width="180" show-overflow-tooltip />
            <el-table-column label="进入时间" width="160">
              <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>

          <el-pagination
            class="pager"
            layout="total, prev, pager, next, sizes"
            :total="recordsTotal" :page-size="recordsSize" :current-page="recordsPage" :page-sizes="[20, 50, 100]"
            @current-change="(p: number) => { recordsPage = p; loadRecords() }"
            @size-change="(s: number) => { recordsSize = s; recordsPage = 1; loadRecords() }"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 去重明细 -->
    <el-dialog v-model="dedupVisible" title="去重规则命中明细" width="460px">
      <el-table :data="dedupRows" border size="small">
        <el-table-column type="index" label="#" width="55" />
        <el-table-column prop="rule" label="去重条件" min-width="220" />
        <el-table-column prop="count" label="命中条数" width="100" align="center" />
      </el-table>
    </el-dialog>

    <!-- 批量调整信息 -->
    <el-dialog v-model="adjustVisible" title="批量调整数据信息" width="500px">
      <el-form label-width="110px">
        <el-form-item label="关联学员名"><el-input v-model="adjustForm.related_name" placeholder="不修改留空" /></el-form-item>
        <el-form-item label="关联手机号"><el-input v-model="adjustForm.related_phone" placeholder="不修改留空" /></el-form-item>
        <el-form-item label="关联年级">
          <el-select v-model="adjustForm.related_grade" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in gradeOpts" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="意向部门">
          <el-select v-model="adjustForm.intention_dept" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in deptOpts" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="意向项目">
          <el-select v-model="adjustForm.intention_project" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in projectOpts" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdjust">确认调整（{{ recIds.length }} 条）</el-button>
      </template>
    </el-dialog>

    <ReceiverDialog v-model="receiverVisible" title="交接给校责任人" @confirm="submitHandover" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Download, Refresh } from '@element-plus/icons-vue'
import { scrmApi, maskPhone, fmtTime, type ImportTask, type DataRecord } from '../api/scrm'
import { downloadCsv } from '../utils/export'
import ReceiverDialog from '../components/ReceiverDialog.vue'

const tab = ref('logs')
const gradeOpts = ref<any[]>([])
const deptOpts = ref<any[]>([])
const projectOpts = ref<any[]>([])

// ---- Tab1 ----
const logs = ref<ImportTask[]>([])
const logsTotal = ref(0)
const logsPage = ref(1)
const logsSize = ref(20)
const lf = reactive<any>({ activity_name: '', task_name: '', operator: '', serial_no: '', batch_no: '', range: null })

function logParams() {
  return {
    page: logsPage.value, page_size: logsSize.value,
    activity_name: lf.activity_name, task_name: lf.task_name, operator: lf.operator,
    serial_no: lf.serial_no, batch_no: lf.batch_no,
    start: lf.range?.[0] || '', end: lf.range?.[1] || ''
  }
}
async function loadLogs() {
  const res = await scrmApi.importLogs(logParams())
  logs.value = res.items
  logsTotal.value = res.total
}
function searchLogs() { logsPage.value = 1; loadLogs() }
function resetLogs() {
  Object.keys(lf).forEach(k => (lf as any)[k] = k === 'range' ? null : '')
  searchLogs()
}
function statusType(s: string) {
  if (!s) return 'info'
  if (s.includes('失败')) return 'danger'
  if (s.includes('完成')) return 'success'
  return 'warning'
}
function downloadFail(id: number) {
  window.open(scrmApi.failFileUrl(id), '_blank')
}

const dedupVisible = ref(false)
const dedupRows = ref<{ rule: string; count: number }[]>([])
async function showDedup(row: ImportTask) {
  const res = await scrmApi.dedupDetails(row.id)
  try {
    dedupRows.value = typeof res.details === 'string' ? JSON.parse(res.details || '[]') : (res.details || [])
  } catch {
    dedupRows.value = []
  }
  dedupVisible.value = true
}

function exportLogs() {
  downloadCsv('导入激活日志.csv',
    ['活动名称', '任务名称', '状态', '导入条数', '成功', '纯新', '激活', '失败', '重复', '无效', '不符合唯一性', '去重去重', '操作人', '流水号', '批次号'],
    logs.value.map(r => [r.activity_name, r.task_name, r.status, r.total_count, r.success_count, r.new_count,
      r.active_count, r.fail_count, r.duplicate_count, r.invalid_count, r.uniqueness_fail_count,
      r.dedup_hit_count, r.operator, r.serial_no, r.batch_no]))
  ElMessage.success(`已导出当前页 ${logs.value.length} 条`)
}

// ---- Tab2 ----
const records = ref<DataRecord[]>([])
const recordsTotal = ref(0)
const recordsPage = ref(1)
const recordsSize = ref(20)
const recIds = ref<number[]>([])
const rf = reactive({ name: '', phone: '', grade: '', school: '', student_no: '', success: '' })

async function loadRecords() {
  const res = await scrmApi.dataRecords({ page: recordsPage.value, page_size: recordsSize.value, ...rf })
  records.value = res.items
  recordsTotal.value = res.total
}
function searchRecords() { recordsPage.value = 1; loadRecords() }
function resetRecords() {
  Object.keys(rf).forEach(k => (rf as any)[k] = '')
  searchRecords()
}

function exportRecords() {
  downloadCsv('数据进入记录.csv',
    ['活动名称', '关联学员', '手机号', '学员号', '活动年级', '关联年级', '数据类型', '进入方式', '进入类型', '渠道', '意向部门', '意向项目', '归属人', '失败原因'],
    records.value.map(r => [r.activity_name, r.related_name, r.related_phone, r.related_student_no,
      r.activity_grade, r.related_grade, r.import_type, r.enter_way, r.lead_enter_type,
      [r.channel_l1, r.channel_l2, r.channel_l3].filter(Boolean).join('/'),
      r.intention_dept, r.intention_project, r.owner_name, r.fail_reason]))
  ElMessage.success(`已导出当前页 ${records.value.length} 条`)
}

const adjustVisible = ref(false)
const adjustForm = reactive<any>({ related_name: '', related_phone: '', related_grade: '', intention_dept: '', intention_project: '' })
async function submitAdjust() {
  const fields: Record<string, string> = {}
  Object.entries(adjustForm).forEach(([k, v]) => { if (v) fields[k] = String(v) })
  if (Object.keys(fields).length === 0) return ElMessage.warning('请至少填写一个要修正的字段')
  await scrmApi.adjustDataRecords(recIds.value, fields)
  ElMessage.success('调整完成')
  adjustVisible.value = false
  Object.keys(adjustForm).forEach(k => (adjustForm as any)[k] = '')
  loadRecords()
}

const receiverVisible = ref(false)
async function submitHandover(emp: any) {
  await scrmApi.handoverDataRecords(recIds.value, emp.name)
  receiverVisible.value = false
  ElMessage.success(`已交接给 ${emp.name}`)
  loadRecords()
}

function onTab(name: string) {
  if (name === 'records' && !records.value.length) loadRecords()
}

onMounted(async () => {
  gradeOpts.value = await scrmApi.dict('grade')
  deptOpts.value = await scrmApi.dict('intention_dept')
  projectOpts.value = await scrmApi.dict('intention_project')
  loadLogs()
})
</script>

<style scoped>
.filter-form :deep(.el-form-item) { margin-bottom: 8px; }
.op-bar { display: flex; justify-content: space-between; align-items: center; margin: 4px 0 10px; }
.pager { margin-top: 10px; justify-content: flex-end; }
</style>
