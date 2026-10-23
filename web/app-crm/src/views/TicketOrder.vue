<template>
  <div class="page">
    <div class="title">退调转 / 订单</div>

    <el-tabs v-model="mainTab">
      <!-- 班课业务办理工单 -->
      <el-tab-pane label="班课业务办理工单" name="ticket">
        <el-tabs v-model="subTab" @tab-change="reloadTickets">
          <el-tab-pane :label="`我发起的 (${cnt.mine})`" name="mine" />
          <el-tab-pane :label="`待我处理 (${cnt.todo})`" name="todo" />
          <el-tab-pane :label="`全部 (${cnt.all})`" name="all" />
        </el-tabs>

        <div class="toolbar">
          <el-radio-group v-model="overtime" size="small" @change="reloadTickets">
            <el-radio-button value="">全部工单</el-radio-button>
            <el-radio-button value="cs24">超24H客服未处理</el-radio-button>
            <el-radio-button value="instant">提单即时审核</el-radio-button>
            <el-radio-button value="day7">超7天未完成</el-radio-button>
          </el-radio-group>
          <el-select v-model="bizType" placeholder="业务类型" clearable size="small" style="width:120px" @change="reloadTickets">
            <el-option label="退班" value="退班" /><el-option label="调班" value="调班" /><el-option label="转班" value="转班" />
          </el-select>
          <el-input v-model="kw" placeholder="工单编号/学员" size="small" style="width:200px" clearable
            @keyup.enter="reloadTickets" @clear="reloadTickets" />
          <el-button type="primary" size="small" @click="reloadTickets">查询</el-button>
        </div>

        <el-table :data="tickets" border size="small" v-loading="loading">
          <el-table-column prop="code" label="工单编号" width="160" />
          <el-table-column prop="submit_name" label="提交人" width="100" />
          <el-table-column prop="biz_type" label="业务类型" width="90" />
          <el-table-column prop="student_name" label="学员姓名" width="100" />
          <el-table-column prop="student_no" label="学员号" width="120" />
          <el-table-column prop="status" label="工单状态" width="110">
            <template #default="{ row }">
              <el-tag :type="ticketTag(row.status)" size="small">{{ row.status }}</el-tag>
              <el-tag v-if="isOvertime(row)" type="danger" size="small" effect="plain" style="margin-left:4px">超时</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="source" label="提单来源" width="100" />
          <el-table-column prop="cs_account" label="分配客服" width="100">
            <template #default="{ row }">{{ row.cs_account || '待分配' }}</template>
          </el-table-column>
          <el-table-column prop="submit_at" label="提交时间" width="160">
            <template #default="{ row }">{{ fmtTime(row.submit_at) }}</template>
          </el-table-column>
        </el-table>
        <el-pagination class="pager" background layout="total,prev,pager,next" :total="ticketTotal"
          v-model:current-page="page" @current-change="reloadTickets" />
      </el-tab-pane>

      <!-- 我的订单（积分商城） -->
      <el-tab-pane label="我的订单" name="order">
        <div class="toolbar">
          <el-form inline size="small">
            <el-form-item label="下单账号">
              <el-input model-value="chenyuankang" disabled style="width:150px" />
            </el-form-item>
          </el-form>
          <el-radio-group v-model="orderChip" size="small" @change="reloadOrders">
            <el-radio-button value="1m">近一个月</el-radio-button>
            <el-radio-button value="3m">近三个月</el-radio-button>
            <el-radio-button value="all">全部</el-radio-button>
          </el-radio-group>
          <el-input v-model="orderKw" placeholder="学员/产品/订单号" size="small" style="width:210px" clearable
            @keyup.enter="reloadOrders" @clear="reloadOrders" />
          <el-button type="primary" size="small" @click="reloadOrders">查询</el-button>
        </div>

        <el-table :data="orders" border size="small" v-loading="orderLoading">
          <el-table-column prop="order_no" label="订单号" width="160" />
          <el-table-column prop="batch_no" label="批次号" width="120" />
          <el-table-column prop="student_code" label="学员编号" width="120" />
          <el-table-column prop="student_name" label="学员姓名" width="100" />
          <el-table-column prop="product_type" label="商品类型" width="100" />
          <el-table-column prop="product_code" label="商品编码" width="100" />
          <el-table-column prop="product_name" label="商品名称" min-width="180" show-overflow-tooltip />
          <el-table-column prop="sale_type" label="销售类型" width="90" />
          <el-table-column prop="biz_type" label="业务类型" width="90" />
          <el-table-column prop="status" label="订单状态" width="100">
            <template #default="{ row }"><el-tag :type="orderTag(row.status)" size="small">{{ row.status }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="created_at" label="下单时间" width="160">
            <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="80" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="detail = row; dlg = true">详情</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination class="pager" background layout="total,prev,pager,next,jumper" :total="orderTotal"
          :page-size="10" v-model:current-page="orderPage" @current-change="reloadOrders" />
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="dlg" title="订单详情" width="460px">
      <el-descriptions v-if="detail" :column="1" border size="small">
        <el-descriptions-item label="订单号">{{ detail.order_no }}</el-descriptions-item>
        <el-descriptions-item label="商品">{{ detail.product_name }}</el-descriptions-item>
        <el-descriptions-item label="学员">{{ detail.student_name }}（{{ detail.student_code }}）</el-descriptions-item>
        <el-descriptions-item label="下单账号">{{ detail.order_account }}</el-descriptions-item>
        <el-descriptions-item label="归属部门">{{ detail.flag_dept }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ detail.status }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ fmtTime(detail.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { scrmApi, fmtTime } from '../api/scrm'

const mainTab = ref('ticket')

// ---- 工单 ----
const subTab = ref('mine')
const overtime = ref('')
const bizType = ref('')
const kw = ref('')
const tickets = ref<any[]>([])
const ticketTotal = ref(0)
const page = ref(1)
const loading = ref(false)
const cnt = ref({ mine: 0, todo: 0, all: 0 })

function ticketTag(s: string) {
  return { 已完成: 'success', 已驳回: 'danger', 待审核: 'warning', 待客服处理: 'primary' }[s] || 'info'
}
function isOvertime(row: any) {
  const hours = (Date.now() - new Date(row.submit_at).getTime()) / 3600000
  return (row.status === '待客服处理' && hours > 24) || (row.status !== '已完成' && hours > 24 * 7)
}

async function loadCounts() {
  const [mine, todo, all] = await Promise.all([
    scrmApi.tickets({ scope: 'mine', page_size: 1 }),
    scrmApi.tickets({ scope: 'all', status: '待客服处理', page_size: 1 }),
    scrmApi.tickets({ scope: 'all', page_size: 1 })
  ])
  cnt.value = { mine: mine.total, todo: todo.total, all: all.total }
}

async function reloadTickets() {
  loading.value = true
  try {
    // 超时过滤在前端按提交时间计算（后端已带超时种子）
    const scope = subTab.value === 'mine' ? 'mine' : 'all'
    const status = subTab.value === 'todo' ? '待客服处理' : ''
    const res = await scrmApi.tickets({
      scope, status, biz_type: bizType.value, keyword: kw.value,
      page: page.value, page_size: 50
    })
    let items = res.items
    if (overtime.value === 'cs24') items = items.filter((r: any) => r.status === '待客服处理' && (Date.now() - new Date(r.submit_at).getTime()) / 3600000 > 24)
    if (overtime.value === 'day7') items = items.filter((r: any) => r.status !== '已完成' && (Date.now() - new Date(r.submit_at).getTime()) / 3600000 > 24 * 7)
    if (overtime.value === 'instant') items = items.filter((r: any) => r.auto_handle === 1 || r.source === '线上自助')
    tickets.value = items
    ticketTotal.value = overtime.value ? items.length : res.total
  } finally { loading.value = false }
}

// ---- 订单 ----
const orderChip = ref('1m')
const orderKw = ref('')
const orders = ref<any[]>([])
const orderTotal = ref(0)
const orderPage = ref(1)
const orderLoading = ref(false)
const dlg = ref(false)
const detail = ref<any>(null)

function orderTag(s: string) {
  return ({ 已完成: 'success', 已发货: 'primary', 待发货: 'warning', 已取消: 'info' } as Record<string, string>)[s] || 'info'
}

async function reloadOrders() {
  orderLoading.value = true
  try {
    const params: Record<string, any> = {
      keyword: orderKw.value, page: orderPage.value, page_size: 10
    }
    if (orderChip.value !== 'all') {
      const months = orderChip.value === '1m' ? 1 : 3
      const d = new Date(); d.setMonth(d.getMonth() - months)
      params.start = d.toISOString().slice(0, 10)
    }
    const res = await scrmApi.mallOrders(params)
    orders.value = res.items; orderTotal.value = res.total
  } finally { orderLoading.value = false }
}

onMounted(() => { loadCounts(); reloadTickets() })
</script>

<style scoped>
.page { padding: 12px 16px; }
.title { font-size: 16px; font-weight: 600; margin-bottom: 8px; }
.toolbar { display: flex; gap: 8px; margin: 8px 0; flex-wrap: wrap; align-items: center; }
.pager { margin-top: 10px; justify-content: flex-end; }
</style>
