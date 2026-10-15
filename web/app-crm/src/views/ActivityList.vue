<template>
  <div class="page">
    <!-- 时间 chip -->
    <div class="chip-bar">
      <el-tag closable type="info" effect="plain" @close="f.created_start = ''; f.created_end = ''; load()">
        创建时间：{{ f.created_start || '2026-08-14' }} 00:00:00 ~ {{ f.created_end || '2026-10-14' }} 23:59:59
      </el-tag>
    </div>

    <el-card shadow="never">
      <!-- 筛选区 -->
      <el-form :model="f" inline class="filter-form">
        <el-form-item label="活动名称"><el-input v-model="f.name" placeholder="活动名称" clearable style="width:170px" @keyup.enter="search" /></el-form-item>
        <el-form-item label="活动类型">
          <el-select v-model="f.type" placeholder="全部" clearable style="width:120px">
            <el-option v-for="o in opts.activity_type || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源系统">
          <el-select v-model="f.source_system" clearable style="width:110px">
            <el-option v-for="o in opts.source_system || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="线索状态">
          <el-select v-model="f.lead_status" clearable style="width:110px">
            <el-option v-for="o in opts.lead_status || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <template v-if="expanded">
          <el-form-item label="活动编码"><el-input v-model="f.code" style="width:160px" /></el-form-item>
          <el-form-item label="活动ID"><el-input v-model="f.external_id" style="width:120px" /></el-form-item>
          <el-form-item label="创建人"><el-input v-model="f.creator" style="width:120px" /></el-form-item>
          <el-form-item label="创建部门"><el-input v-model="f.creator_dept" style="width:120px" /></el-form-item>
          <el-form-item label="创建时间">
            <el-date-picker v-model="createdRange" type="daterange" value-format="YYYY-MM-DD"
              start-placeholder="开始" end-placeholder="结束" style="width:240px" />
          </el-form-item>
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
          <el-button @click="exportCsv"><el-icon><Download /></el-icon>导出</el-button>
          <el-button @click="editSelected" :disabled="selectedRow == null">编辑活动</el-button>
          <el-button type="primary" plain @click="importData" :disabled="selectedRow == null">导入数据</el-button>
          <el-button type="primary" @click="openCreate"><el-icon><Plus /></el-icon>新建活动</el-button>
          <el-button @click="handoverVisible = true" :disabled="selectedRow == null">线索交接</el-button>
        </div>
        <el-button text @click="load"><el-icon><Refresh /></el-icon>刷新</el-button>
      </div>

      <!-- 表格 -->
      <el-table
        :data="rows"
        border
        size="small"
        height="56vh"
        highlight-current-row
        @current-change="(r: any) => selectedRow = r"
        @row-dblclick="drillLeads"
      >
        <el-table-column width="48" align="center">
          <template #default="{ row }">
            <el-radio :model-value="selectedRow?.id" :value="row.id" @change="selectedRow = row"><span /></el-radio>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="活动名称" min-width="200" fixed>
          <template #default="{ row }">
            <el-link type="primary" :underline="false" @click="drillLeads(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="source_system" label="来源系统" width="90" />
        <el-table-column prop="category" label="活动类别" width="100" />
        <el-table-column prop="type" label="活动类型" width="80" />
        <el-table-column label="明细总数" width="80" align="right">
          <template #default="{ row }">{{ row.stats?.total || 0 }}</template>
        </el-table-column>
        <el-table-column label="线索总数" width="80" align="right">
          <template #default="{ row }">{{ row.stats?.total || 0 }}</template>
        </el-table-column>
        <el-table-column label="纯新线索" width="80" align="right">
          <template #default="{ row }">{{ row.stats?.['enter:纯新线索'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="激活线索" width="80" align="right">
          <template #default="{ row }">{{ row.stats?.['enter:激活线索'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="待交接" width="70" align="right">
          <template #default="{ row }">
            <span class="num-link" @click="drillLeads(row, '待交接')">{{ row.stats?.['status:待交接'] || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="待分配" width="70" align="right">
          <template #default="{ row }">{{ row.stats?.['status:待分配'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="待跟进" width="70" align="right">
          <template #default="{ row }">{{ row.stats?.['status:待跟进'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="跟进中" width="70" align="right">
          <template #default="{ row }">{{ row.stats?.['status:跟进中'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="已回收" width="70" align="right">
          <template #default="{ row }">{{ row.stats?.['status:已回收'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="无效归档" width="80" align="right">
          <template #default="{ row }">{{ row.stats?.['status:无效归档'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="入库归档" width="80" align="right">
          <template #default="{ row }">{{ row.stats?.['status:入库归档'] || 0 }}</template>
        </el-table-column>
        <el-table-column label="已转商" width="70" align="right">
          <template #default="{ row }">
            <span style="color:#f56c6c;font-weight:600">{{ row.stats?.['status:已转商'] || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="活动状态" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="['info', 'success', 'warning'][row.status]">
              {{ ['未开始', '进行中', '已结束'][row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target_object" label="活动对象" width="90" />
        <el-table-column prop="intention_dept" label="意向部门" width="90" />
        <el-table-column prop="intention_project" label="意向项目" width="90" />
        <el-table-column label="开始时间" width="110">
          <template #default="{ row }">{{ fmt(row.start_time) }}</template>
        </el-table-column>
        <el-table-column label="结束时间" width="110">
          <template #default="{ row }">{{ fmt(row.end_time) }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">{{ fmt(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="creator" label="创建人" width="90" />
        <el-table-column prop="creator_dept" label="创建部门" width="90" />
        <el-table-column prop="remark" label="活动备注" min-width="140" show-overflow-tooltip />
        <el-table-column prop="external_id" label="活动Id" width="90" />
        <el-table-column prop="code" label="活动编码" min-width="150" show-overflow-tooltip />
      </el-table>

      <el-pagination
        class="pager"
        layout="total, prev, pager, next, sizes"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        :page-sizes="[20, 50, 100]"
        @current-change="(p: number) => { page = p; load() }"
        @size-change="(s: number) => { pageSize = s; page = 1; load() }"
      />
    </el-card>

    <!-- 新建/编辑活动 -->
    <el-dialog v-model="dlg" :title="editing ? '编辑活动' : '新建活动'" width="640px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="活动名称" required>
          <el-input v-model="form.name" placeholder="如 20261014-高中-迂回转学" />
        </el-form-item>
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="活动类别">
              <el-select v-model="form.category" style="width:100%">
                <el-option v-for="o in opts.activity_category || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="活动类型">
              <el-select v-model="form.type" style="width:100%">
                <el-option v-for="o in opts.activity_type || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="活动对象">
              <el-select v-model="form.target_object" style="width:100%">
                <el-option v-for="o in opts.target_object || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="活动状态">
              <el-select v-model="form.status" style="width:100%">
                <el-option :value="0" label="未开始" />
                <el-option :value="1" label="进行中" />
                <el-option :value="2" label="已结束" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="意向部门">
              <el-select v-model="form.intention_dept" style="width:100%">
                <el-option v-for="o in opts.intention_dept || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="意向项目">
              <el-select v-model="form.intention_project" style="width:100%">
                <el-option v-for="o in opts.intention_project || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="起止时间">
          <el-date-picker v-model="dateRange" type="daterange" value-format="YYYY-MM-DD"
            start-placeholder="开始日期" end-placeholder="结束日期" style="width:100%" />
        </el-form-item>
        <el-form-item label="活动备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>

    <ReceiverDialog v-model="handoverVisible" title="活动线索交接" :count="activityLeadCount" @confirm="doHandover" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Plus, Refresh, ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import { scrmApi, fmtTime, type Activity } from '../api/scrm'
import { downloadCsv } from '../utils/export'
import ReceiverDialog from '../components/ReceiverDialog.vue'

const router = useRouter()
const opts = ref<Record<string, any[]>>({})
const rows = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const expanded = ref(false)
const selectedRow = ref<Activity | null>(null)
const createdRange = ref<string[]>([])

const f = reactive({
  name: '', type: '', source_system: '', lead_status: '', code: '', external_id: '',
  creator: '', creator_dept: '', created_start: '', created_end: ''
})

async function load() {
  f.created_start = createdRange.value?.[0] || ''
  f.created_end = createdRange.value?.[1] || ''
  const res = await scrmApi.activities({ page: page.value, page_size: pageSize.value, ...f })
  rows.value = res.items
  total.value = res.total
}
function search() { page.value = 1; load() }
function reset() {
  Object.keys(f).forEach(k => (f as any)[k] = '')
  createdRange.value = []
  search()
}

function fmt(s: string) { return fmtTime(s) }

function drillLeads(row: Activity, status?: string) {
  router.push({ path: '/all-leads', query: { activity_id: String(row.id), activity_name: row.name, ...(status ? { status } : {}) } })
}

function exportCsv() {
  downloadCsv('活动列表.csv',
    ['活动名称', '来源系统', '类别', '类型', '线索总数', '纯新', '激活', '待交接', '待跟进', '跟进中', '已转商', '状态', '创建人', '活动Id', '活动编码'],
    rows.value.map((r: any) => [r.name, r.source_system, r.category, r.type, r.stats?.total || 0,
      r.stats?.['enter:纯新线索'] || 0, r.stats?.['enter:激活线索'] || 0,
      r.stats?.['status:待交接'] || 0, r.stats?.['status:待跟进'] || 0, r.stats?.['status:跟进中'] || 0,
      r.stats?.['status:已转商'] || 0, ['未开始', '进行中', '已结束'][r.status], r.creator, r.external_id, r.code]))
  ElMessage.success(`已导出当前页 ${rows.value.length} 条活动`)
}

function importData() {
  if (!selectedRow.value) return
  router.push({ path: '/lead-import', query: { activity_id: String(selectedRow.value.id), name: selectedRow.value.name } })
}

// 新建/编辑
const dlg = ref(false)
const editing = ref(false)
const dateRange = ref<string[]>([])
const form = reactive<any>({ name: '', category: '', type: '', target_object: '', status: 0, intention_dept: '', intention_project: '', remark: '' })

function openCreate() {
  editing.value = false
  Object.assign(form, { id: 0, name: '', category: '', type: '', target_object: '', status: 0, intention_dept: '', intention_project: '', remark: '' })
  dateRange.value = []
  dlg.value = true
}
function editSelected() {
  if (!selectedRow.value) return
  const r = selectedRow.value
  editing.value = true
  Object.assign(form, {
    id: r.id, name: r.name, category: r.category, type: r.type, target_object: r.target_object,
    status: r.status, intention_dept: r.intention_dept, intention_project: r.intention_project, remark: r.remark
  })
  dateRange.value = r.start_time ? [fmtTime(r.start_time).slice(0, 10), fmtTime(r.end_time).slice(0, 10)] : []
  dlg.value = true
}
async function submit() {
  if (!form.name.trim()) return ElMessage.warning('请填写活动名称')
  const payload: any = { ...form }
  if (dateRange.value?.length === 2) {
    payload.start_time = dateRange.value[0] + ' 00:00:00'
    payload.end_time = dateRange.value[1] + ' 23:59:59'
  }
  if (editing.value) {
    await scrmApi.updateActivity(form.id, payload)
    ElMessage.success('活动已更新')
  } else {
    await scrmApi.createActivity(payload)
    ElMessage.success('活动已创建')
  }
  dlg.value = false
  load()
}

// 活动线索交接
const handoverVisible = ref(false)
const activityLeadCount = ref(0)
async function doHandover(emp: any) {
  const res = await scrmApi.leads({ view: 'handover', activity_id: selectedRow.value?.id, page: 1, page_size: 500 })
  activityLeadCount.value = res.total
  const ids = res.items.map((x: any) => x.id)
  if (ids.length === 0) {
    ElMessage.warning('该活动下暂无待交接线索')
    handoverVisible.value = false
    return
  }
  await ElMessageBox.confirm(`该活动下有 ${ids.length} 条待交接线索，确认交接给 ${emp.name}？`, '确认', { type: 'warning' })
  await scrmApi.handoverLeads(ids, emp.id, emp.name)
  ElMessage.success('交接完成')
  handoverVisible.value = false
  load()
}

onMounted(async () => {
  const kinds = ['activity_type', 'source_system', 'lead_status', 'activity_category', 'target_object', 'intention_dept', 'intention_project']
  await Promise.all(kinds.map(async k => { opts.value[k] = await scrmApi.dict(k) }))
  load()
})
</script>

<style scoped>
.chip-bar { margin-bottom: 10px; }
.filter-form :deep(.el-form-item) { margin-bottom: 8px; }
.op-bar { display: flex; justify-content: space-between; align-items: center; margin: 4px 0 10px; }
.pager { margin-top: 10px; justify-content: flex-end; }
.num-link { color: #00a86b; cursor: pointer; }
.num-link:hover { text-decoration: underline; }
</style>
