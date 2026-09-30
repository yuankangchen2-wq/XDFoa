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
          <el-button type="primary" @click="openDialog">新增成本记录</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:140px">
            <el-option label="待归集" :value="1" />
            <el-option label="已归集" :value="2" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="product_id" label="产品ID" width="100" />
        <el-table-column prop="work_order_id" label="工单ID" width="100" />
        <el-table-column prop="cost_center_id" label="成本中心ID" width="120" />
        <el-table-column label="要素" width="100">
          <template #default="{ row }">
            <el-tag :type="elemType[row.element]">{{ elemLabel[row.element] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="130" />
        <el-table-column prop="period" label="期间" width="110" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'warning' : 'success'">
              {{ row.status === 1 ? '待归集' : '已归集' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column prop="created_at" label="创建时间" width="200" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增成本记录" width="500px">
      <el-form :model="form" label-width="110px">
        <el-form-item label="产品ID">
          <el-input-number v-model="form.product_id" :min="1" />
        </el-form-item>
        <el-form-item label="工单ID">
          <el-input-number v-model="form.work_order_id" :min="0" />
        </el-form-item>
        <el-form-item label="成本中心ID">
          <el-input-number v-model="form.cost_center_id" :min="0" />
        </el-form-item>
        <el-form-item label="成本要素">
          <el-select v-model="form.element" style="width:100%">
            <el-option label="物料" :value="1" />
            <el-option label="人工" :value="2" />
            <el-option label="制造费用" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="form.amount" :min="0.01" :precision="4" style="width:100%" />
        </el-form-item>
        <el-form-item label="期间">
          <el-input v-model="form.period" placeholder="2026-09" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
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
import { costApi } from '../api/cost'

const elemLabel: Record<number, string> = { 1: '物料', 2: '人工', 3: '制造费用' }
const elemType: Record<number, string> = { 1: '', 2: 'success', 3: 'warning' }

const list = ref<any[]>([])
const query = reactive({ status: null as any })

const dialogVisible = ref(false)
const form = reactive({ product_id: 1, work_order_id: 0, cost_center_id: 0, element: 1, amount: 0, period: '', remark: '' })

const loadList = async () => {
  const res: any = await costApi.listRecords(query)
  list.value = res.items || []
}

const openDialog = () => {
  Object.assign(form, { product_id: 1, work_order_id: 0, cost_center_id: 0, element: 1, amount: 0, period: '', remark: '' })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await costApi.addRecord(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

onMounted(loadList)
</script>
