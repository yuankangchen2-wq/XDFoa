<template>
  <div class="page">
    <div class="head">
      <div class="title">审核中心</div>
      <div>
        <el-button type="success" size="small" @click="handle(true)">批量通过</el-button>
        <el-button type="danger" size="small" @click="handle(false)">批量驳回</el-button>
      </div>
    </div>

    <el-tabs v-model="tab" @tab-change="onTab">
      <el-tab-pane label="申请成单待审" name="deal" />
      <el-tab-pane label="业绩申诉待审" name="appeal" />
    </el-tabs>

    <div class="toolbar">
      <el-radio-group v-model="status" size="small" @change="reload">
        <el-radio-button :value="pendingStatus">{{ pendingLabel }}</el-radio-button>
        <el-radio-button value="已通过">已通过</el-radio-button>
        <el-radio-button value="已驳回">已驳回</el-radio-button>
      </el-radio-group>
      <el-input v-model="kw" placeholder="商机编码/编号" size="small" style="width:220px" clearable
        @keyup.enter="reload" @clear="reload" />
      <el-button type="primary" size="small" @click="reload">查询</el-button>
    </div>

    <!-- 成单审批 -->
    <el-table v-if="tab === 'deal'" :data="rows" border size="small" v-loading="loading"
      @selection-change="(v: any[]) => selection = v">
      <el-table-column type="selection" width="42" />
      <el-table-column prop="opp_code" label="商机编码" width="160" />
      <el-table-column prop="biz_type" label="业务类型" width="90" />
      <el-table-column prop="cert_no" label="订单/听课证号" width="160" />
      <el-table-column prop="amount" label="订单金额（元）" width="130">
        <template #default="{ row }">
          <span class="money">{{ row.amount }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="applicant" label="申请人" width="100" />
      <el-table-column prop="applicant_dept" label="申请部门" min-width="150" />
      <el-table-column prop="reason" label="申请说明" min-width="180" show-overflow-tooltip />
      <el-table-column prop="status" label="审批状态" width="100">
        <template #default="{ row }"><el-tag :type="tagType(row.status)" size="small">{{ row.status }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="created_at" label="申请时间" width="160">
        <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>

    <!-- 申诉审批 -->
    <el-table v-else :data="rows" border size="small" v-loading="loading"
      @selection-change="(v: any[]) => selection = v">
      <el-table-column type="selection" width="42" />
      <el-table-column prop="biz_line" label="业务线" width="90" />
      <el-table-column prop="opp_code" label="商机编码" width="160" />
      <el-table-column prop="cert_no" label="业务编号" width="150" />
      <el-table-column prop="owner_name" label="当前归属人" width="110" />
      <el-table-column prop="applicant" label="申诉人" width="100" />
      <el-table-column prop="adjust_mode" label="调整模式" width="100" />
      <el-table-column prop="reason" label="申诉原因" min-width="220" show-overflow-tooltip />
      <el-table-column prop="status" label="审批状态" width="100">
        <template #default="{ row }"><el-tag :type="tagType(row.status)" size="small">{{ row.status }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="created_at" label="申诉时间" width="160">
        <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>

    <el-pagination class="pager" background layout="total,prev,pager,next" :total="total"
      v-model:current-page="page" @current-change="reload" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { scrmApi, fmtTime } from '../api/scrm'

const tab = ref('deal')
const status = ref('待审核')
const kw = ref('')
const rows = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const selection = ref<any[]>([])

const pendingStatus = computed(() => tab.value === 'deal' ? '待审核' : '待审批')
const pendingLabel = computed(() => tab.value === 'deal' ? '待审核' : '待审批')

function tagType(s: string) {
  return s === '已通过' ? 'success' : (s === '已驳回' || s === '驳回') ? 'danger' : 'warning'
}

async function reload() {
  loading.value = true
  try {
    const res = tab.value === 'deal'
      ? await scrmApi.approvals({ keyword: kw.value, status: status.value, page: page.value, page_size: 20 })
      : await scrmApi.appeals({ keyword: kw.value, status: status.value, page: page.value, page_size: 20 })
    rows.value = res.items; total.value = res.total
  } finally { loading.value = false }
}

function onTab() {
  status.value = pendingStatus.value
  page.value = 1; kw.value = ''
  reload()
}

async function handle(pass: boolean) {
  if (!selection.value.length) return ElMessage.warning('请先勾选记录')
  await ElMessageBox.confirm(`确认${pass ? '通过' : '驳回'}选中的 ${selection.value.length} 条记录？`, '审批确认', {
    type: pass ? 'success' : 'warning'
  })
  const ids = selection.value.map(x => x.id)
  if (tab.value === 'deal') await scrmApi.handleApprovals(ids, pass)
  else await scrmApi.handleAppeals(ids, pass)
  ElMessage.success(pass ? '已通过，通过的成单申请已写入成单日志并更新商机状态' : '已驳回')
  reload()
}

onMounted(reload)
</script>

<style scoped>
.page { padding: 12px 16px; }
.head { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 16px; font-weight: 600; }
.toolbar { display: flex; gap: 8px; margin: 8px 0; }
.pager { margin-top: 10px; justify-content: flex-end; }
.money { color: #f56c6c; font-weight: 600; }
</style>
