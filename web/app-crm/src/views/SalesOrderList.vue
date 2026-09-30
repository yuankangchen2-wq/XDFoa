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
          <el-button type="primary" @click="openDialog">新建订单</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:140px">
            <el-option label="待审批" :value="1" />
            <el-option label="已审批" :value="2" />
            <el-option label="已发货" :value="3" />
            <el-option label="已完成" :value="4" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="order_no" label="订单号" width="200" />
        <el-table-column prop="customer_id" label="客户ID" width="80" />
        <el-table-column prop="total_amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.total_amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button v-if="row.status === 1" size="small" type="primary" @click="approve(row)">审批</el-button>
            <el-button v-if="row.status === 2" size="small" type="warning" @click="ship(row)">发货</el-button>
            <el-button v-if="row.status === 3" size="small" type="success" @click="complete(row)">完成</el-button>
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

    <el-dialog v-model="dialogVisible" title="新建销售订单" width="450px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="客户ID">
          <el-input-number v-model="form.customer_id" :min="1" />
        </el-form-item>
        <el-form-item label="订单金额">
          <el-input-number v-model="form.total_amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { salesApi } from '../api/sales'

const statusLabel: Record<number, string> = { 1: '待审批', 2: '已审批', 3: '已发货', 4: '已完成', 9: '已取消' }
const statusType: Record<number, string> = { 1: 'warning', 2: '', 3: 'info', 4: 'success', 9: 'danger' }

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: null as any })

const dialogVisible = ref(false)
const form = reactive({ customer_id: 1, total_amount: 0, remark: '' })

const loadList = async () => {
  const res: any = await salesApi.listOrders(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openDialog = () => {
  Object.assign(form, { customer_id: 1, total_amount: 0, remark: '' })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await salesApi.createOrder(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const approve = async (row: any) => {
  await ElMessageBox.confirm('确认审批该订单？', '提示', { type: 'warning' })
  try { await salesApi.approveOrder(row.id); ElMessage.success('审批成功'); loadList() }
  catch (e: any) { ElMessage.error(e.response?.data?.error || '审批失败') }
}

const ship = async (row: any) => {
  try { await salesApi.shipOrder(row.id); ElMessage.success('已发货'); loadList() }
  catch (e: any) { ElMessage.error(e.response?.data?.error || '发货失败') }
}

const complete = async (row: any) => {
  try { await salesApi.completeOrder(row.id); ElMessage.success('已完成'); loadList() }
  catch (e: any) { ElMessage.error(e.response?.data?.error || '操作失败') }
}

onMounted(loadList)
</script>
