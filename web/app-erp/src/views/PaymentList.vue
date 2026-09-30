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
          <span style="color:#666">收付款流水</span>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="类型">
          <el-select v-model="query.type" placeholder="全部" clearable style="width:140px">
            <el-option label="收款" :value="1" />
            <el-option label="付款" :value="2" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 1 ? 'success' : 'warning'">{{ row.type === 1 ? '收款' : '付款' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ref_id" label="关联ID" width="100" />
        <el-table-column label="金额" width="140">
          <template #default="{ row }">
            <span :style="{ color: row.type === 1 ? '#67c23a' : '#f56c6c', fontWeight: 'bold' }">
              {{ row.type === 1 ? '+' : '-' }}{{ row.amount.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="pay_method" label="支付方式" width="120" />
        <el-table-column prop="pay_date" label="日期" width="130" />
        <el-table-column prop="remark" label="备注" />
        <el-table-column prop="created_at" label="创建时间" width="200" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { financeApi } from '../api/finance'

const list = ref<any[]>([])
const query = reactive({ type: null as any })

const loadList = async () => {
  const res: any = await financeApi.listPayments(query)
  list.value = res.items || []
}

onMounted(loadList)
</script>
