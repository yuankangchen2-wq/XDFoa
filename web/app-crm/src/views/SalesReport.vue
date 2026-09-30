<template>
  <div style="padding:20px">
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <el-menu mode="horizontal" :default-active="$route.path" router style="border:none">
            <el-menu-item index="/customers">客户管理</el-menu-item>
            <el-menu-item index="/opportunities">商机管理</el-menu-item>
            <el-menu-item index="/sales/orders">销售订单</el-menu-item>
            <el-menu-item index="/sales/contracts">合同管理</el-menu-item>
            <el-menu-item index="/tickets">工单管理</el-menu-item>
            <el-menu-item index="/analytics/dashboard">数据分析</el-menu-item>
          </el-menu>
          <el-button type="primary" @click="openDialog">录入销售数据</el-button>
        </div>
      </template>

      <h4 style="margin:0 0 12px">按期间销售趋势</h4>
      <el-table :data="summary" border stripe>
        <el-table-column prop="period" label="期间" width="140" />
        <el-table-column prop="total_amount" label="销售额" width="160" />
        <el-table-column prop="order_count" label="订单数" width="120" />
        <el-table-column prop="avg_amount" label="客单价" width="160" />
      </el-table>

      <h4 style="margin:20px 0 12px">销售快照明细</h4>
      <el-form :inline="true">
        <el-form-item label="期间">
          <el-input v-model="query.period" placeholder="2026-09" style="width:140px" />
        </el-form-item>
        <el-button type="primary" @click="loadSales">查询</el-button>
      </el-form>
      <el-table :data="list" border stripe>
        <el-table-column prop="period" label="期间" width="120" />
        <el-table-column prop="product_id" label="产品ID" width="100" />
        <el-table-column prop="customer_id" label="客户ID" width="100" />
        <el-table-column prop="total_amount" label="销售额" width="140" />
        <el-table-column prop="order_count" label="订单数" width="100" />
        <el-table-column prop="avg_amount" label="客单价" width="140" />
        <el-table-column prop="created_at" label="创建时间" width="200" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="录入销售快照" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="期间">
          <el-input v-model="form.period" placeholder="2026-09" />
        </el-form-item>
        <el-form-item label="产品ID">
          <el-input-number v-model="form.product_id" :min="0" />
        </el-form-item>
        <el-form-item label="客户ID">
          <el-input-number v-model="form.customer_id" :min="0" />
        </el-form-item>
        <el-form-item label="销售额">
          <el-input-number v-model="form.total_amount" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item label="订单数">
          <el-input-number v-model="form.order_count" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { analyticsApi } from '../api/analytics'

const summary = ref<any[]>([])
const list = ref<any[]>([])
const query = reactive({ period: '' })

const dialogVisible = ref(false)
const form = reactive({ period: '', product_id: 0, customer_id: 0, total_amount: 0, order_count: 0 })

const loadSummary = async () => {
  try {
    const res: any = await analyticsApi.salesSummary()
    summary.value = res.items || []
  } catch { summary.value = [] }
}

const loadSales = async () => {
  try {
    const res: any = await analyticsApi.listSales(query)
    list.value = res.items || []
  } catch { list.value = [] }
}

const openDialog = () => {
  Object.assign(form, { period: '', product_id: 0, customer_id: 0, total_amount: 0, order_count: 0 })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await analyticsApi.recordSales(form)
    ElMessage.success('录入成功')
    dialogVisible.value = false
    loadSummary(); loadSales()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '录入失败')
  }
}

onMounted(() => { loadSummary(); loadSales() })
</script>
