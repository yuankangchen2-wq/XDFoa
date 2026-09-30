<template>
  <div style="padding:20px;max-width:1200px;margin:0 auto">
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <h3 style="margin:0">我的订单</h3>
          <el-button type="primary" @click="openCreateDialog">下单</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:140px">
            <el-option label="待支付" :value="1" />
            <el-option label="已支付" :value="2" />
            <el-option label="已发货" :value="3" />
            <el-option label="已完成" :value="4" />
            <el-option label="已取消" :value="9" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="order_no" label="订单号" width="200" />
        <el-table-column label="商品" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="item in row.items" :key="item.id" style="margin:2px">
              {{ item.sku_code }} × {{ item.quantity }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.total_amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="180" />
        <el-table-column label="操作" width="260">
          <template #default="{ row }">
            <el-button size="small" @click="$router.push(`/orders/${row.id}`)">详情</el-button>
            <el-button v-if="row.status === 1" size="small" type="success" @click="pay(row)">支付</el-button>
            <el-button v-if="row.status === 1" size="small" type="danger" @click="cancel(row)">取消</el-button>
            <el-button v-if="row.status === 2" size="small" type="primary" @click="ship(row)">发货</el-button>
            <el-button v-if="row.status === 3" size="small" type="success" @click="complete(row)">确认收货</el-button>
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

    <!-- 下单弹窗 -->
    <el-dialog v-model="dialogVisible" title="创建订单" width="550px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="仓库ID">
          <el-input-number v-model="form.warehouse_id" :min="1" />
        </el-form-item>
        <el-form-item label="商品明细">
          <el-table :data="form.items" border size="small">
            <el-table-column label="SKU编码" width="160">
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
          <el-button size="small" style="margin-top:8px" @click="form.items.push({ sku_code: '', quantity: 1, unit_price: 0 })">添加商品</el-button>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createOrder">提交订单</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { orderApi } from '../api/order'

const statusLabel: Record<number, string> = { 1: '待支付', 2: '已支付', 3: '已发货', 4: '已完成', 9: '已取消' }
const statusType: Record<number, string> = { 1: 'warning', 2: '', 3: 'info', 4: 'success', 9: 'danger' }

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: null as any, user_id: 1001 })

const dialogVisible = ref(false)
const form = reactive({ warehouse_id: 1, remark: '', items: [{ sku_code: '', quantity: 1, unit_price: 0 }] })

const loadList = async () => {
  const res: any = await orderApi.list(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openCreateDialog = () => {
  form.warehouse_id = 1
  form.remark = ''
  form.items = [{ sku_code: '', quantity: 1, unit_price: 0 }]
  dialogVisible.value = true
}

const createOrder = async () => {
  try {
    await orderApi.create(form)
    ElMessage.success('下单成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '下单失败')
  }
}

const pay = async (row: any) => {
  await ElMessageBox.confirm('确认支付？', '提示', { type: 'warning' })
  try {
    await orderApi.pay(row.id)
    ElMessage.success('支付成功')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '支付失败')
  }
}

const cancel = async (row: any) => {
  await ElMessageBox.confirm('确认取消订单？', '提示', { type: 'warning' })
  try {
    await orderApi.cancel(row.id)
    ElMessage.success('已取消')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '取消失败')
  }
}

const ship = async (row: any) => {
  try {
    await orderApi.ship(row.id)
    ElMessage.success('已发货')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '发货失败')
  }
}

const complete = async (row: any) => {
  try {
    await orderApi.complete(row.id)
    ElMessage.success('已确认收货')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  }
}

onMounted(loadList)
</script>
