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
          <div>
            <el-button @click="openInDialog">入库</el-button>
            <el-button type="warning" @click="openOutDialog">出库</el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="SKU编码">
          <el-input v-model="query.sku_code" placeholder="模糊搜索" clearable />
        </el-form-item>
        <el-form-item label="仓库">
          <el-select v-model="query.warehouse_id" placeholder="全部" clearable style="width:160px">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.warehouse_name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="sku_code" label="SKU编码" width="140" />
        <el-table-column prop="warehouse_id" label="仓库ID" width="100" />
        <el-table-column prop="quantity" label="可用库存" width="120">
          <template #default="{ row }">
            <el-tag :type="row.quantity < row.safety_stock ? 'danger' : 'success'">{{ row.quantity }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="allocated" label="已预扣" width="100" />
        <el-table-column prop="safety_stock" label="安全库存" width="100" />
        <el-table-column prop="batch_no" label="批次" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.quantity < row.safety_stock" type="danger">预警</el-tag>
            <el-tag v-else type="success">正常</el-tag>
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

    <!-- 入库弹窗 -->
    <el-dialog v-model="inDialog" title="入库" width="450px">
      <el-form :model="inForm" label-width="100px">
        <el-form-item label="SKU编码">
          <el-input v-model="inForm.sku_code" />
        </el-form-item>
        <el-form-item label="仓库">
          <el-select v-model="inForm.warehouse_id" style="width:100%">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.warehouse_name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="inForm.quantity" :min="1" />
        </el-form-item>
        <el-form-item label="批次号">
          <el-input v-model="inForm.batch_no" />
        </el-form-item>
        <el-form-item label="关联单号">
          <el-input v-model="inForm.reference_id" placeholder="采购单号等" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="inDialog = false">取消</el-button>
        <el-button type="primary" @click="stockIn">确认入库</el-button>
      </template>
    </el-dialog>

    <!-- 出库弹窗 -->
    <el-dialog v-model="outDialog" title="出库" width="450px">
      <el-form :model="outForm" label-width="100px">
        <el-form-item label="SKU编码">
          <el-input v-model="outForm.sku_code" />
        </el-form-item>
        <el-form-item label="仓库">
          <el-select v-model="outForm.warehouse_id" style="width:100%">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.warehouse_name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="outForm.quantity" :min="1" />
        </el-form-item>
        <el-form-item label="关联单号">
          <el-input v-model="outForm.reference_id" placeholder="订单号等" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="outDialog = false">取消</el-button>
        <el-button type="primary" @click="stockOut">确认出库</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { stockApi, warehouseApi } from '../api/inventory'

const list = ref<any[]>([])
const total = ref(0)
const warehouses = ref<any[]>([])
const query = reactive({ page: 1, page_size: 20, sku_code: '', warehouse_id: null as any })

const inDialog = ref(false)
const inForm = reactive({ sku_code: '', warehouse_id: null as any, quantity: 1, batch_no: '', reference_id: '' })

const outDialog = ref(false)
const outForm = reactive({ sku_code: '', warehouse_id: null as any, quantity: 1, reference_id: '' })

const loadList = async () => {
  const res: any = await stockApi.list(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const loadWarehouses = async () => {
  const res: any = await warehouseApi.list({ page: 1, page_size: 100 })
  warehouses.value = res.items || []
}

const openInDialog = () => { inDialog.value = true }
const openOutDialog = () => { outDialog.value = true }

const stockIn = async () => {
  try {
    await stockApi.stockIn(inForm)
    ElMessage.success('入库成功')
    inDialog.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '入库失败')
  }
}

const stockOut = async () => {
  try {
    await stockApi.deduct(outForm)
    ElMessage.success('出库成功')
    outDialog.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '出库失败')
  }
}

onMounted(() => { loadList(); loadWarehouses() })
</script>
