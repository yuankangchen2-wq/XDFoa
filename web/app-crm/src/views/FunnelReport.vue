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
          <el-button type="primary" @click="openDialog">录入漏斗阶段</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="期间">
          <el-input v-model="period" placeholder="2026-09" style="width:140px" />
        </el-form-item>
        <el-button type="primary" @click="loadFunnel">查询</el-button>
      </el-form>

      <!-- 漏斗可视化 -->
      <div v-if="funnel.length" style="margin:20px 0">
        <div v-for="(stage, i) in funnel" :key="stage.id" style="display:flex;align-items:center;margin-bottom:10px">
          <div style="width:120px;text-align:right;padding-right:12px">{{ stage.stage }}</div>
          <div style="flex:1;max-width:600px">
            <div :style="{
              width: widthPct(i) + '%',
              background: colors[i % colors.length],
              color: '#fff',
              padding: '12px 16px',
              borderRadius: '4px',
              display: 'flex',
              justifyContent: 'space-between'
            }">
              <span>商机数：{{ stage.opportunity_count }}</span>
              <span>金额：¥{{ stage.amount.toFixed(2) }}</span>
            </div>
          </div>
          <div style="width:100px;text-align:center;color:#909399">
            {{ stage.conversion_rate > 0 ? stage.conversion_rate.toFixed(1) + '%' : '-' }}
          </div>
        </div>
      </div>
      <el-empty v-else description="暂无漏斗数据" />

      <el-table :data="funnel" border stripe style="margin-top:16px">
        <el-table-column prop="period" label="期间" width="120" />
        <el-table-column prop="stage" label="阶段" />
        <el-table-column prop="opportunity_count" label="商机数" width="120" />
        <el-table-column prop="amount" label="金额" width="160" />
        <el-table-column prop="conversion_rate" label="转化率" width="120" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="录入漏斗阶段" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="期间">
          <el-input v-model="form.period" placeholder="2026-09" />
        </el-form-item>
        <el-form-item label="阶段">
          <el-input v-model="form.stage" placeholder="线索/商机/方案/成交" />
        </el-form-item>
        <el-form-item label="商机数">
          <el-input-number v-model="form.opportunity_count" :min="0" />
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="form.amount" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item label="转化率">
          <el-input-number v-model="form.conversion_rate" :min="0" :max="100" :precision="2" style="width:100%" />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { analyticsApi } from '../api/analytics'

const funnel = ref<any[]>([])
const period = ref('')

const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399']

const maxCount = computed(() => {
  return funnel.value.length ? Math.max(...funnel.value.map(f => f.opportunity_count)) : 0
})

const widthPct = (i: number) => {
  if (!maxCount.value) return 100
  return Math.max((funnel.value[i].opportunity_count / maxCount.value) * 100, 15)
}

const dialogVisible = ref(false)
const form = reactive({ period: '', stage: '', opportunity_count: 0, amount: 0, conversion_rate: 0 })

const loadFunnel = async () => {
  try {
    const res: any = await analyticsApi.listFunnel({ period: period.value })
    funnel.value = res.items || []
  } catch { funnel.value = [] }
}

const openDialog = () => {
  Object.assign(form, { period: '', stage: '', opportunity_count: 0, amount: 0, conversion_rate: 0 })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await analyticsApi.recordFunnel(form)
    ElMessage.success('录入成功')
    dialogVisible.value = false
    loadFunnel()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '录入失败')
  }
}

onMounted(loadFunnel)
</script>
