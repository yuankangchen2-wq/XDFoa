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
          <el-button type="primary" @click="openDialog">创建调拨</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="transfer_no" label="调拨单号" width="180" />
        <el-table-column prop="from_warehouse_id" label="源仓库" width="100" />
        <el-table-column prop="to_warehouse_id" label="目标仓库" width="100" />
        <el-table-column prop="sku_code" label="SKU" width="140" />
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button v-if="row.status === 1" size="small" type="primary" @click="confirmOut(row)">确认出库</el-button>
            <el-button v-if="row.status === 2" size="small" type="success" @click="confirmIn(row)">确认入库</el-button>
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

    <el-dialog v-model="dialogVisible" title="创建调拨单" width="450px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="源仓库">
          <el-select v-model="form.from_warehouse_id" style="width:100%">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.warehouse_name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标仓库">
          <el-select v-model="form.to_warehouse_id" style="width:100%">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.warehouse_name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="SKU编码">
          <el-input v-model="form.sku_code" />
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="form.quantity" :min="1" />
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
import { transferApi, warehouseApi } from '../api/inventory'

const statusLabel: Record<number, string> = { 1: '待出库', 2: '已出库', 3: '已入库', 4: '已取消' }
const statusType: Record<number, string> = { 1: 'warning', 2: '', 3: 'success', 4: 'info' }

const list = ref<any[]>([])
const total = ref(0)
const warehouses = ref<any[]>([])
const query = reactive({ page: 1, page_size: 20 })
const dialogVisible = ref(false)
const form = reactive({ from_warehouse_id: null as any, to_warehouse_id: null as any, sku_code: '', quantity: 1 })

const loadList = async () => {
  const res: any = await transferApi.list(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const loadWarehouses = async () => {
  const res: any = await warehouseApi.list({ page: 1, page_size: 100 })
  warehouses.value = res.items || []
}

const openDialog = () => { dialogVisible.value = true }

const save = async () => {
  try {
    await transferApi.create(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const confirmOut = async (row: any) => {
  await ElMessageBox.confirm('确认从源仓库出库？', '提示', { type: 'warning' })
  try {
    await transferApi.confirmOut(row.id)
    ElMessage.success('出库成功')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '出库失败')
  }
}

const confirmIn = async (row: any) => {
  await ElMessageBox.confirm('确认入目标仓库？', '提示', { type: 'warning' })
  try {
    await transferApi.confirmIn(row.id)
    ElMessage.success('入库成功')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '入库失败')
  }
}

onMounted(() => { loadList(); loadWarehouses() })
</script>
