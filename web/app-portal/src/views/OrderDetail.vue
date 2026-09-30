<template>
  <div style="padding:20px;max-width:900px;margin:0 auto">
    <el-page-header @back="$router.back()" title="返回列表" />
    <el-card style="margin-top:16px" v-if="order">
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <span>订单号：{{ order.order_no }}</span>
          <el-tag :type="statusType[order.status]">{{ statusLabel[order.status] }}</el-tag>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ order.user_id }}</el-descriptions-item>
        <el-descriptions-item label="仓库ID">{{ order.warehouse_id }}</el-descriptions-item>
        <el-descriptions-item label="总金额">¥{{ order.total_amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ order.created_at }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ order.remark }}</el-descriptions-item>
      </el-descriptions>

      <h4 style="margin:20px 0 12px">商品明细</h4>
      <el-table :data="order.items" border>
        <el-table-column prop="sku_code" label="SKU编码" width="160" />
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column prop="unit_price" label="单价" width="120" />
        <el-table-column label="小计" width="140">
          <template #default="{ row }">¥{{ (row.quantity * row.unit_price).toFixed(2) }}</template>
        </el-table-column>
      </el-table>

      <div style="margin-top:20px;text-align:right">
        <el-button v-if="order.status === 1" type="success" @click="pay">立即支付</el-button>
        <el-button v-if="order.status === 1" type="danger" @click="cancel">取消订单</el-button>
        <el-button v-if="order.status === 2" type="primary" @click="ship">发货</el-button>
        <el-button v-if="order.status === 3" type="success" @click="complete">确认收货</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orderApi } from '../api/order'

const route = useRoute()
const order = ref<any>(null)

const statusLabel: Record<number, string> = { 1: '待支付', 2: '已支付', 3: '已发货', 4: '已完成', 9: '已取消' }
const statusType: Record<number, string> = { 1: 'warning', 2: '', 3: 'info', 4: 'success', 9: 'danger' }

const loadOrder = async () => {
  const id = Number(route.params.id)
  order.value = await orderApi.get(id)
}

const pay = async () => {
  await ElMessageBox.confirm('确认支付？', '提示', { type: 'warning' })
  try {
    await orderApi.pay(order.value.id)
    ElMessage.success('支付成功')
    loadOrder()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '支付失败')
  }
}

const cancel = async () => {
  await ElMessageBox.confirm('确认取消订单？', '提示', { type: 'warning' })
  try {
    await orderApi.cancel(order.value.id)
    ElMessage.success('已取消')
    loadOrder()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '取消失败')
  }
}

const ship = async () => {
  try {
    await orderApi.ship(order.value.id)
    ElMessage.success('已发货')
    loadOrder()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '发货失败')
  }
}

const complete = async () => {
  try {
    await orderApi.complete(order.value.id)
    ElMessage.success('已确认收货')
    loadOrder()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  }
}

onMounted(loadOrder)
</script>
