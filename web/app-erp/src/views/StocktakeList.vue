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
          <el-button type="primary" @click="openDialog">创建盘点</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="stocktake_no" label="盘点单号" width="200" />
        <el-table-column prop="warehouse_id" label="仓库" width="100" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'warning' : 'success'">
              {{ row.status === 1 ? '进行中' : '已完成' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
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

    <el-dialog v-model="dialogVisible" title="创建盘点单" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="仓库">
          <el-select v-model="form.warehouse_id" style="width:100%">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.warehouse_name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="盘点明细">
          <el-table :data="form.items" border size="small">
            <el-table-column label="SKU编码" width="160">
              <template #default="{ row }">
                <el-input v-model="row.sku_code" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="系统数量" width="110">
              <template #default="{ row }">
                <el-input-number v-model="row.system_qty" size="small" :min="0" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="实盘数量" width="110">
              <template #default="{ row }">
                <el-input-number v-model="row.actual_qty" size="small" :min="0" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80">
              <template #default="{ $index }">
                <el-button size="small" type="danger" link @click="form.items.splice($index, 1)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button size="small" style="margin-top:8px" @click="form.items.push({ sku_code: '', system_qty: 0, actual_qty: 0 })">添加明细</el-button>
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
import { ElMessage } from 'element-plus'
import { stocktakeApi, warehouseApi } from '../api/inventory'

const list = ref<any[]>([])
const total = ref(0)
const warehouses = ref<any[]>([])
const query = reactive({ page: 1, page_size: 20 })
const dialogVisible = ref(false)
const form = reactive({ warehouse_id: null as any, items: [] as any[] })

const loadList = async () => {
  const res: any = await stocktakeApi.list(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const loadWarehouses = async () => {
  const res: any = await warehouseApi.list({ page: 1, page_size: 100 })
  warehouses.value = res.items || []
}

const openDialog = () => {
  form.warehouse_id = null
  form.items = [{ sku_code: '', system_qty: 0, actual_qty: 0 }]
  dialogVisible.value = true
}

const save = async () => {
  try {
    await stocktakeApi.create(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

onMounted(() => { loadList(); loadWarehouses() })
</script>
