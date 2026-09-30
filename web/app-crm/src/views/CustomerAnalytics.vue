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
          <el-button type="primary" @click="openDialog">录入客户统计</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="period" label="期间" width="120" />
        <el-table-column prop="new_customers" label="新增客户" width="120" />
        <el-table-column prop="active_customers" label="活跃客户" width="120" />
        <el-table-column label="活跃率" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.new_customers > 0" type="success">
              {{ ((row.active_customers / row.new_customers) * 100).toFixed(1) }}%
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_contribution" label="总贡献金额" width="160" />
        <el-table-column prop="top_customer_id" label="TOP客户ID" width="120" />
        <el-table-column prop="top_contribution" label="TOP贡献金额" width="160" />
        <el-table-column prop="created_at" label="创建时间" width="200" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="录入客户统计" width="500px">
      <el-form :model="form" label-width="110px">
        <el-form-item label="期间">
          <el-input v-model="form.period" placeholder="2026-09" />
        </el-form-item>
        <el-form-item label="新增客户">
          <el-input-number v-model="form.new_customers" :min="0" />
        </el-form-item>
        <el-form-item label="活跃客户">
          <el-input-number v-model="form.active_customers" :min="0" />
        </el-form-item>
        <el-form-item label="总贡献金额">
          <el-input-number v-model="form.total_contribution" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item label="TOP客户ID">
          <el-input-number v-model="form.top_customer_id" :min="0" />
        </el-form-item>
        <el-form-item label="TOP贡献金额">
          <el-input-number v-model="form.top_contribution" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { analyticsApi } from '../api/analytics'

const list = ref<any[]>([])

const dialogVisible = ref(false)
const form = reactive({
  period: '', new_customers: 0, active_customers: 0,
  total_contribution: 0, top_customer_id: 0, top_contribution: 0
})

const loadList = async () => {
  try {
    const res: any = await analyticsApi.listCustomerStats({})
    list.value = res.items || []
  } catch { list.value = [] }
}

const openDialog = () => {
  Object.assign(form, { period: '', new_customers: 0, active_customers: 0, total_contribution: 0, top_customer_id: 0, top_contribution: 0 })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await analyticsApi.recordCustomerStat(form)
    ElMessage.success('录入成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '录入失败')
  }
}

onMounted(loadList)
</script>
