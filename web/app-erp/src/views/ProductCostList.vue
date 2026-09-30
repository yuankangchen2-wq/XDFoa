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
            <el-menu-item index="/cost/centers">成本中心</el-menu-item>
            <el-menu-item index="/cost/products">产品成本</el-menu-item>
            <el-menu-item index="/cost/records">成本记录</el-menu-item>
          </el-menu>
          <div>
            <el-button type="success" @click="openCollect">成本归集</el-button>
            <el-button type="primary" @click="openCalc">手动录入</el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="产品ID">
          <el-input-number v-model="query.product_id" :min="0" style="width:120px" />
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="product_id" label="产品ID" width="100" />
        <el-table-column prop="period" label="期间" width="110" />
        <el-table-column prop="material_cost" label="物料成本" width="130" />
        <el-table-column prop="labor_cost" label="人工成本" width="130" />
        <el-table-column prop="overhead_cost" label="制造费用" width="130" />
        <el-table-column label="总成本" width="140">
          <template #default="{ row }">
            <span style="font-weight:bold;color:#f56c6c">{{ row.total_cost.toFixed(4) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="110" />
        <el-table-column label="单位成本" width="130">
          <template #default="{ row }">
            <span style="color:#409eff;font-weight:bold">{{ row.unit_cost.toFixed(4) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="200" />
      </el-table>
    </el-card>

    <!-- 手动录入成本 -->
    <el-dialog v-model="calcVisible" title="手动录入产品成本" width="560px">
      <el-form :model="calcForm" label-width="100px">
        <el-form-item label="产品ID">
          <el-input-number v-model="calcForm.product_id" :min="1" />
        </el-form-item>
        <el-form-item label="物料成本">
          <el-input-number v-model="calcForm.material_cost" :min="0" :precision="4" style="width:100%" />
        </el-form-item>
        <el-form-item label="人工成本">
          <el-input-number v-model="calcForm.labor_cost" :min="0" :precision="4" style="width:100%" />
        </el-form-item>
        <el-form-item label="制造费用">
          <el-input-number v-model="calcForm.overhead_cost" :min="0" :precision="4" style="width:100%" />
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="calcForm.quantity" :min="0.01" :precision="4" style="width:100%" />
        </el-form-item>
        <el-form-item label="期间">
          <el-input v-model="calcForm.period" placeholder="2026-09" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="calcVisible = false">取消</el-button>
        <el-button type="primary" @click="doCalc">计算并保存</el-button>
      </template>
    </el-dialog>

    <!-- 成本归集 -->
    <el-dialog v-model="collectVisible" title="成本归集（按产品+期间汇总待归集记录）" width="450px">
      <el-form label-width="100px">
        <el-form-item label="产品ID">
          <el-input-number v-model="collectForm.product_id" :min="1" />
        </el-form-item>
        <el-form-item label="期间">
          <el-input v-model="collectForm.period" placeholder="2026-09" />
        </el-form-item>
        <el-form-item label="数量">
          <el-input-number v-model="collectForm.quantity" :min="0.01" :precision="4" style="width:100%" />
        </el-form-item>
      </el-form>
      <el-alert type="info" :closable="false" style="margin-bottom:12px">
        将汇总该产品在指定期间内所有「待归集」成本记录，生成产品成本，并把记录标记为「已归集」。
      </el-alert>
      <template #footer>
        <el-button @click="collectVisible = false">取消</el-button>
        <el-button type="success" @click="doCollect">执行归集</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { costApi } from '../api/cost'

const list = ref<any[]>([])
const query = reactive({ product_id: undefined as any })

const calcVisible = ref(false)
const calcForm = reactive({ product_id: 1, material_cost: 0, labor_cost: 0, overhead_cost: 0, quantity: 1, period: '' })

const collectVisible = ref(false)
const collectForm = reactive({ product_id: 1, period: '', quantity: 1 })

const loadList = async () => {
  const res: any = await costApi.listProductCosts(query)
  list.value = res.items || []
}

const openCalc = () => {
  Object.assign(calcForm, { product_id: 1, material_cost: 0, labor_cost: 0, overhead_cost: 0, quantity: 1, period: '' })
  calcVisible.value = true
}

const doCalc = async () => {
  try {
    await costApi.calculateProductCost(calcForm)
    ElMessage.success('成本计算完成')
    calcVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '计算失败')
  }
}

const openCollect = () => {
  Object.assign(collectForm, { product_id: 1, period: '', quantity: 1 })
  collectVisible.value = true
}

const doCollect = async () => {
  try {
    const res: any = await costApi.collectCost(collectForm)
    ElMessage.success(`归集成功，总成本 ${res.total_cost.toFixed(4)}，单位成本 ${res.unit_cost.toFixed(4)}`)
    collectVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '归集失败')
  }
}

onMounted(loadList)
</script>
