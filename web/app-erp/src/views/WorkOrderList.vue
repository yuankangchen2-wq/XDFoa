<template>
  <div style="padding:20px">
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <el-menu mode="horizontal" :default-active="$route.path" router style="border:none">
            <el-menu-item index="/stocks">库存管理</el-menu-item>
            <el-menu-item index="/purchase/orders">采购订单</el-menu-item>
            <el-menu-item index="/production/work-orders">生产工单</el-menu-item>
            <el-menu-item index="/finance/receivables">应收账款</el-menu-item>
            <el-menu-item index="/finance/payables">应付账款</el-menu-item>
            <el-menu-item index="/finance/payments">收付款记录</el-menu-item>
          </el-menu>
          <el-button type="primary" @click="openDialog">新建工单</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:140px">
            <el-option label="待领料" :value="1" />
            <el-option label="生产中" :value="2" />
            <el-option label="已完工" :value="3" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="order_no" label="工单号" width="200" />
        <el-table-column prop="product_sku" label="产品SKU" width="140" />
        <el-table-column prop="quantity" label="计划数量" width="100" />
        <el-table-column prop="produced_qty" label="已生产" width="100" />
        <el-table-column prop="warehouse_id" label="入库仓" width="80" />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button v-if="row.status === 1" size="small" type="warning" @click="issue(row)">领料</el-button>
            <el-button v-if="row.status === 2" size="small" type="success" @click="openReportDialog(row)">报工</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.page_size"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadList"
        style="margin-top:16px;justify-content:flex-end"
      />
    </el-card>

    <!-- 新建工单 -->
    <el-dialog v-model="dialogVisible" title="新建生产工单" width="450px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="产品SKU">
          <el-input v-model="form.product_sku" />
        </el-form-item>
        <el-form-item label="计划数量">
          <el-input-number v-model="form.quantity" :min="1" />
        </el-form-item>
        <el-form-item label="入库仓库">
          <el-input-number v-model="form.warehouse_id" :min="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">创建</el-button>
      </template>
    </el-dialog>

    <!-- 报工 -->
    <el-dialog v-model="reportDialogVisible" title="生产报工" width="400px">
      <el-form label-width="100px">
        <el-form-item label="工单号">{{ currentOrder?.order_no }}</el-form-item>
        <el-form-item label="产品">{{ currentOrder?.product_sku }}</el-form-item>
        <el-form-item label="计划数量">{{ currentOrder?.quantity }}</el-form-item>
        <el-form-item label="已生产">{{ currentOrder?.produced_qty }}</el-form-item>
        <el-form-item label="本次报工">
          <el-input-number v-model="reportQty" :min="1" :max="remainingQty" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reportDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmReport">确认报工</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { productionApi } from '../api/production'

const statusLabel: Record<number, string> = { 1: '待领料', 2: '生产中', 3: '已完工', 9: '已取消' }
const statusType: Record<number, string> = { 1: 'info', 2: 'warning', 3: 'success', 9: 'danger' }

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: null as any })

const dialogVisible = ref(false)
const form = reactive({ product_sku: '', quantity: 1, warehouse_id: 1 })

const reportDialogVisible = ref(false)
const currentOrder = ref<any>(null)
const reportQty = ref(1)

const remainingQty = computed(() => {
  if (!currentOrder.value) return 0
  return currentOrder.value.quantity - currentOrder.value.produced_qty
})

const loadList = async () => {
  const res: any = await productionApi.listWorkOrders(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openDialog = () => {
  form.product_sku = ''
  form.quantity = 1
  form.warehouse_id = 1
  dialogVisible.value = true
}

const save = async () => {
  try {
    await productionApi.createWorkOrder(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const issue = async (row: any) => {
  await ElMessageBox.confirm('确认领料？将展开BOM并扣减物料库存。', '提示', { type: 'warning' })
  try {
    await productionApi.issueMaterial(row.id)
    ElMessage.success('领料成功')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '领料失败')
  }
}

const openReportDialog = (row: any) => {
  currentOrder.value = row
  reportQty.value = Math.min(row.quantity - row.produced_qty, row.quantity)
  reportDialogVisible.value = true
}

const confirmReport = async () => {
  try {
    await productionApi.reportProduction(currentOrder.value.id, { produced_qty: reportQty.value })
    ElMessage.success('报工成功')
    reportDialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '报工失败')
  }
}

onMounted(loadList)
</script>
