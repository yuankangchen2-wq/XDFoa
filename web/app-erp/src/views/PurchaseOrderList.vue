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
          <el-button type="primary" @click="openDialog">新建采购单</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:140px">
            <el-option label="草稿" :value="1" />
            <el-option label="已审批" :value="2" />
            <el-option label="部分到货" :value="3" />
            <el-option label="全部到货" :value="4" />
            <el-option label="已结算" :value="5" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="order_no" label="采购单号" width="180" />
        <el-table-column prop="supplier_code" label="供应商" width="120" />
        <el-table-column prop="warehouse_id" label="入库仓" width="80" />
        <el-table-column prop="total_amount" label="总金额" width="120">
          <template #default="{ row }">¥{{ row.total_amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">明细</el-button>
            <el-button v-if="row.status === 1" size="small" type="primary" @click="approve(row)">审批</el-button>
            <el-button v-if="row.status === 4" size="small" type="success" @click="settle(row)">结算</el-button>
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

    <!-- 新建采购单弹窗 -->
    <el-dialog v-model="dialogVisible" title="新建采购订单" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="供应商">
          <el-input v-model="form.supplier_code" placeholder="供应商编码" />
        </el-form-item>
        <el-form-item label="入库仓库">
          <el-input-number v-model="form.warehouse_id" :min="1" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
        <el-form-item label="明细">
          <el-table :data="form.items" border size="small">
            <el-table-column label="SKU编码" width="140">
              <template #default="{ row }">
                <el-input v-model="row.sku_code" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="数量" width="100">
              <template #default="{ row }">
                <el-input-number v-model="row.quantity" size="small" :min="1" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="单价" width="120">
              <template #default="{ row }">
                <el-input-number v-model="row.unit_price" size="small" :min="0" :precision="2" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="70">
              <template #default="{ $index }">
                <el-button size="small" type="danger" link @click="form.items.splice($index, 1)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button size="small" style="margin-top:8px" @click="form.items.push({ sku_code: '', quantity: 1, unit_price: 0 })">添加明细</el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">提交</el-button>
      </template>
    </el-dialog>

    <!-- 明细弹窗 -->
    <el-dialog v-model="detailVisible" title="订单明细" width="600px">
      <el-table :data="detailItems" border size="small">
        <el-table-column prop="sku_code" label="SKU" width="140" />
        <el-table-column prop="quantity" label="采购数量" width="100" />
        <el-table-column prop="received_qty" label="已到货" width="100" />
        <el-table-column prop="unit_price" label="单价" width="100" />
        <el-table-column label="小计" width="120">
          <template #default="{ row }">¥{{ (row.quantity * row.unit_price).toFixed(2) }}</template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { purchaseApi } from '../api/purchase'

const statusLabel: Record<number, string> = { 1: '草稿', 2: '已审批', 3: '部分到货', 4: '全部到货', 5: '已结算', 9: '已取消' }
const statusType: Record<number, string> = { 1: 'info', 2: 'warning', 3: 'warning', 4: 'success', 5: '', 9: 'danger' }

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: null as any })

const dialogVisible = ref(false)
const form = reactive({ supplier_code: '', warehouse_id: 1, remark: '', items: [{ sku_code: '', quantity: 1, unit_price: 0 }] })

const detailVisible = ref(false)
const detailItems = ref<any[]>([])

const loadList = async () => {
  const res: any = await purchaseApi.listOrders(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openDialog = () => {
  form.supplier_code = ''
  form.warehouse_id = 1
  form.remark = ''
  form.items = [{ sku_code: '', quantity: 1, unit_price: 0 }]
  dialogVisible.value = true
}

const save = async () => {
  try {
    await purchaseApi.createOrder(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const approve = async (row: any) => {
  await ElMessageBox.confirm('确认审批该采购单？', '提示', { type: 'warning' })
  try {
    await purchaseApi.approveOrder(row.id)
    ElMessage.success('审批成功')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '审批失败')
  }
}

const settle = async (row: any) => {
  await ElMessageBox.confirm('确认结算该采购单？', '提示', { type: 'warning' })
  try {
    await purchaseApi.settleOrder(row.id)
    ElMessage.success('结算成功')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '结算失败')
  }
}

const viewDetail = async (row: any) => {
  const res: any = await purchaseApi.getOrder(row.id)
  detailItems.value = res.items || []
  detailVisible.value = true
}

onMounted(loadList)
</script>
