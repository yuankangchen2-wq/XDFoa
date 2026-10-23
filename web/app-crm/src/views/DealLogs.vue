<template>
  <div class="page">
    <div class="head">
      <div class="title">成单日志</div>
      <div>
        <el-button size="small" @click="exportRows">导出</el-button>
        <el-button type="primary" size="small" @click="toFinance">到库入财</el-button>
      </div>
    </div>

    <el-tabs v-model="bizLine" @tab-change="reload">
      <el-tab-pane v-for="b in bizLines" :key="b" :label="b" :name="b" />
    </el-tabs>

    <div class="toolbar">
      <el-radio-group v-model="dateChip" size="small" @change="reload">
        <el-radio-button value="3m">近三个月</el-radio-button>
        <el-radio-button value="6m">近半年</el-radio-button>
        <el-radio-button value="1y">近一年</el-radio-button>
        <el-radio-button value="all">全部</el-radio-button>
      </el-radio-group>
      <el-input v-model="kw" placeholder="商机编码/学员/编号" size="small" style="width:230px" clearable
        @keyup.enter="reload" @clear="reload" />
      <el-button type="primary" size="small" @click="reload">查询</el-button>
    </div>

    <el-table :data="rows" border size="small" v-loading="loading" @selection-change="(v: any[]) => selection = v">
      <el-table-column type="selection" width="42" />
      <el-table-column prop="opp_code" label="商机编码" width="160" />
      <el-table-column prop="owner_name" label="当前所属人" width="110" />
      <el-table-column prop="order_person" label="下单人" width="100" />
      <el-table-column prop="student_no" label="学员号" width="120" />
      <el-table-column prop="student_name" label="姓名" width="90" />
      <el-table-column prop="product_name" :label="productCol" min-width="180" show-overflow-tooltip />
      <el-table-column prop="amount" label="报班金额（元）" width="140">
        <template #default="{ row }"><span class="money">{{ row.amount }}</span></template>
      </el-table-column>
      <el-table-column prop="cert_no" :label="certCol" width="160" />
      <el-table-column prop="order_no" label="订单编号" width="150" />
      <el-table-column prop="deal_type" label="成单类型" width="90" />
      <el-table-column prop="pay_at" label="支付时间" width="160">
        <template #default="{ row }">{{ fmtTime(row.pay_at) }}</template>
      </el-table-column>
      <el-table-column prop="match_at" label="订单匹配商机时间" width="170">
        <template #default="{ row }">{{ fmtTime(row.match_at) }}</template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" background layout="total,prev,pager,next" :total="total"
      v-model:current-page="page" @current-change="reload" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { scrmApi, fmtTime } from '../api/scrm'

const bizLines = ['班课', '高端', '会员卡', '课时包', '国游']
const bizLine = ref('班课')
const dateChip = ref('6m')
const kw = ref('')
const rows = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const selection = ref<any[]>([])

const productCol = computed(() => ({ 班课: '班级名称', 国游: '游学产品', 高端: '高端产品' })[bizLine.value] || '产品名称')
const certCol = computed(() => ({ 班课: '听课证号', 高端: '合同编号', 会员卡: '交易编码', 课时包: '交易编号', 国游: '资源编号' })[bizLine.value])

function range() {
  if (dateChip.value === 'all') return null
  const months = { '3m': 3, '6m': 6, '1y': 12 }[dateChip.value] || 6
  const d = new Date(); d.setMonth(d.getMonth() - months)
  return [d.toISOString().slice(0, 10), '']
}

async function reload() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      biz_line: bizLine.value, keyword: kw.value, page: page.value, page_size: 20
    }
    const rg = range()
    if (rg) { params.start = rg[0]; params.end = rg[1] }
    const res = await scrmApi.dealLogs(params)
    rows.value = res.items; total.value = res.total
  } finally { loading.value = false }
}

function exportRows() { ElMessage.info(`已按当前条件导出「${bizLine.value}」成单日志（演示）`) }
function toFinance() {
  if (!selection.value.length) return ElMessage.warning('请先勾选记录')
  ElMessage.success(`已将 ${selection.value.length} 条成单记录推送至库入财（演示）`)
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
