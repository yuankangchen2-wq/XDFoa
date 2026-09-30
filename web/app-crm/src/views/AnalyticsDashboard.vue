<template>
  <div style="padding:20px">
    <el-card>
      <template #header>
        <el-menu mode="horizontal" :default-active="$route.path" router style="border:none">
          <el-menu-item index="/customers">客户管理</el-menu-item>
          <el-menu-item index="/opportunities">商机管理</el-menu-item>
          <el-menu-item index="/sales/orders">销售订单</el-menu-item>
          <el-menu-item index="/sales/contracts">合同管理</el-menu-item>
          <el-menu-item index="/tickets">工单管理</el-menu-item>
          <el-menu-item index="/analytics/dashboard">数据分析</el-menu-item>
        </el-menu>
      </template>

      <div v-if="dashboard">
        <el-row :gutter="16">
          <el-col :span="6">
            <el-card shadow="hover" style="text-align:center;background:linear-gradient(135deg,#409eff,#66b1ff);color:#fff">
              <div style="font-size:14px;opacity:0.9">销售总额</div>
              <div style="font-size:28px;font-weight:bold;margin-top:8px">¥{{ dashboard.total_sales.toFixed(2) }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" style="text-align:center;background:linear-gradient(135deg,#67c23a,#85ce61);color:#fff">
              <div style="font-size:14px;opacity:0.9">订单总数</div>
              <div style="font-size:28px;font-weight:bold;margin-top:8px">{{ dashboard.total_orders }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" style="text-align:center;background:linear-gradient(135deg,#e6a23c,#ebb563);color:#fff">
              <div style="font-size:14px;opacity:0.9">客单价</div>
              <div style="font-size:28px;font-weight:bold;margin-top:8px">¥{{ dashboard.avg_ticket.toFixed(2) }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" style="text-align:center;background:linear-gradient(135deg,#f56c6c,#f78989);color:#fff">
              <div style="font-size:14px;opacity:0.9">转化率</div>
              <div style="font-size:28px;font-weight:bold;margin-top:8px">{{ dashboard.conversion_rate.toFixed(2) }}%</div>
            </el-card>
          </el-col>
        </el-row>

        <el-row :gutter="16" style="margin-top:16px">
          <el-col :span="12">
            <el-card>
              <template #header><b>客户指标</b></template>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="新增客户">{{ dashboard.new_customers }}</el-descriptions-item>
                <el-descriptions-item label="待处理工单">{{ dashboard.open_tickets }}</el-descriptions-item>
              </el-descriptions>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card>
              <template #header><b>快速入口</b></template>
              <el-space wrap>
                <el-button type="primary" @click="$router.push('/analytics/sales')">销售报表</el-button>
                <el-button type="success" @click="$router.push('/analytics/funnel')">销售漏斗</el-button>
                <el-button type="warning" @click="$router.push('/analytics/customers')">客户分析</el-button>
              </el-space>
            </el-card>
          </el-col>
        </el-row>
      </div>
      <el-skeleton v-else :rows="6" animated />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { analyticsApi } from '../api/analytics'

const dashboard = ref<any>(null)

onMounted(async () => {
  try {
    dashboard.value = await analyticsApi.dashboard()
  } catch {
    // 后端未启动时展示默认空数据
    dashboard.value = {
      total_sales: 0, total_orders: 0, new_customers: 0,
      open_tickets: 0, conversion_rate: 0, avg_ticket: 0
    }
  }
})
</script>
