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
          <el-button type="primary" @click="openDialog">新建到货</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="receipt_no" label="到货单号" width="200" />
        <el-table-column prop="order_id" label="采购单ID" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 2 ? 'success' : 'warning'">
              {{ row.status === 1 ? '待入库' : '已入库' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="到货时间" width="180" />
        <el-table-column label="明细" width="300">
          <template #default="{ row }">
            <span v-for="item in row.items" :key="item.id" style="margin-right:12px">
              {{ item.sku_code }} × {{ item.quantity }}
            </span>
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

    <el-dialog v-model="dialogVisible" title="新建到货单" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="采购单ID">
          <el-input-number v-model="form.order_id" :min="1" />
        </el-form-item>
        <el-form-item label="到货明细">
          <el-table :data="form.items" border size="small">
            <el-table-column label="SKU编码" width="160">
              <template #default="{ row }">
                <el-input v-model="row.sku_code" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="数量" width="120">
              <template #default="{ row }">
                <el-input-number v-model="row.quantity" size="small" :min="1" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="70">
              <template #default="{ $index }">
                <el-button size="small" type="danger" link @click="form.items.splice($index, 1)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button size="small" style="margin-top:8px" @click="form.items.push({ sku_code: '', quantity: 1 })">添加明细</el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">确认到货入库</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { purchaseApi } from '../api/purchase'

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20 })

const dialogVisible = ref(false)
const form = reactive({ order_id: 1, items: [{ sku_code: '', quantity: 1 }] })

const loadList = async () => {
  const res: any = await purchaseApi.listReceipts(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openDialog = () => {
  form.order_id = 1
  form.items = [{ sku_code: '', quantity: 1 }]
  dialogVisible.value = true
}

const save = async () => {
  try {
    await purchaseApi.receiveGoods(form)
    ElMessage.success('到货入库成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '到货失败')
  }
}

onMounted(loadList)
</script>
