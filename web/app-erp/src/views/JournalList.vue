<template>
  <div style="padding:20px">
    <el-card>
      <template #header>
        <el-menu mode="horizontal" :default-active="$route.path" router style="border:none">
          <el-menu-item index="/stocks">库存管理</el-menu-item>
          <el-menu-item index="/purchase/orders">采购订单</el-menu-item>
          <el-menu-item index="/production/work-orders">生产工单</el-menu-item>
          <el-menu-item index="/finance/receivables">应收账款</el-menu-item>
          <el-menu-item index="/finance/payables">应付账款</el-menu-item>
          <el-menu-item index="/finance/payments">收付款记录</el-menu-item>
        </el-menu>
      </template>

      <el-form :inline="true">
        <el-form-item label="SKU编码">
          <el-input v-model="query.sku_code" placeholder="精确搜索" clearable />
        </el-form-item>
        <el-form-item label="仓库">
          <el-select v-model="query.warehouse_id" placeholder="全部" clearable style="width:160px">
            <el-option v-for="w in warehouses" :key="w.id" :label="w.warehouse_name" :value="w.id" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="流水号" width="80" />
        <el-table-column prop="sku_code" label="SKU编码" width="140" />
        <el-table-column prop="warehouse_id" label="仓库" width="80" />
        <el-table-column prop="change_type" label="变动类型" width="130">
          <template #default="{ row }">
            <el-tag :type="typeMap[row.change_type]">{{ typeLabel[row.change_type] || row.change_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="变动数量" width="120">
          <template #default="{ row }">
            <span :style="{ color: row.quantity > 0 ? '#67c23a' : '#f56c6c' }">
              {{ row.quantity > 0 ? '+' : '' }}{{ row.quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="before_qty" label="变动前" width="100" />
        <el-table-column prop="after_qty" label="变动后" width="100" />
        <el-table-column prop="reference_id" label="关联单号" width="140" />
        <el-table-column prop="created_at" label="时间" width="180" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { journalApi, warehouseApi } from '../api/inventory'

const typeLabel: Record<string, string> = {
  in: '入库', out: '出库', allocate: '预扣', release: '释放',
  transfer_in: '调拨入库', transfer_out: '调拨出库', adjust: '盘点调整'
}
const typeMap: Record<string, string> = {
  in: 'success', out: 'danger', allocate: 'warning', release: 'info',
  transfer_in: 'success', transfer_out: 'danger', adjust: ''
}

const list = ref<any[]>([])
const total = ref(0)
const warehouses = ref<any[]>([])
const query = reactive({ page: 1, page_size: 20, sku_code: '', warehouse_id: null as any })

const loadList = async () => {
  const res: any = await journalApi.list(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const loadWarehouses = async () => {
  const res: any = await warehouseApi.list({ page: 1, page_size: 100 })
  warehouses.value = res.items || []
}

onMounted(() => { loadList(); loadWarehouses() })
</script>
