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
          <el-button type="primary" @click="openDialog">新建应收</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:140px">
            <el-option label="未收" :value="1" />
            <el-option label="部分已收" :value="2" />
            <el-option label="已收清" :value="3" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="customer_id" label="客户ID" width="100" />
        <el-table-column prop="amount" label="应收金额" width="130" />
        <el-table-column prop="received_amount" label="已收金额" width="130" />
        <el-table-column label="未收" width="130">
          <template #default="{ row }">
            <span style="color:#f56c6c">{{ (row.amount - row.received_amount).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="due_date" label="到期日" width="120" />
        <el-table-column label="操作" width="130">
          <template #default="{ row }">
            <el-button size="small" v-if="row.status < 3" type="success" @click="openReceive(row)">收款</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新建应收 -->
    <el-dialog v-model="dialogVisible" title="新建应收账款" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="客户ID">
          <el-input-number v-model="form.customer_id" :min="1" />
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="form.amount" :min="0" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item label="到期日">
          <el-input v-model="form.due_date" placeholder="2026-12-31" />
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

    <!-- 收款 -->
    <el-dialog v-model="receiveVisible" title="收款" width="450px">
      <div v-if="receiveRow" style="margin-bottom:12px">
        <p>应收金额：<b>{{ receiveRow.amount }}</b></p>
        <p>已收金额：<b>{{ receiveRow.received_amount }}</b></p>
        <p style="color:#f56c6c">未收金额：<b>{{ (receiveRow.amount - receiveRow.received_amount).toFixed(2) }}</b></p>
      </div>
      <el-form label-width="100px">
        <el-form-item label="收款金额">
          <el-input-number v-model="receiveForm.amount" :min="0.01" :precision="2" style="width:100%" />
        </el-form-item>
        <el-form-item label="收款方式">
          <el-input v-model="receiveForm.pay_method" placeholder="转账/现金/..." />
        </el-form-item>
        <el-form-item label="收款日期">
          <el-input v-model="receiveForm.pay_date" placeholder="2026-09-29" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="receiveVisible = false">取消</el-button>
        <el-button type="primary" @click="doReceive">确认收款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { financeApi } from '../api/finance'

const statusLabel: Record<number, string> = { 1: '未收', 2: '部分已收', 3: '已收清' }
const statusType: Record<number, string> = { 1: 'warning', 2: '', 3: 'success' }

const list = ref<any[]>([])
const query = reactive({ status: null as any })

const dialogVisible = ref(false)
const form = reactive({ customer_id: 1, amount: 0, due_date: '', remark: '' })

const receiveVisible = ref(false)
const receiveRow = ref<any>(null)
const receiveForm = reactive({ amount: 0, pay_method: '转账', pay_date: '' })

const loadList = async () => {
  const res: any = await financeApi.listReceivables(query)
  list.value = res.items || []
}

const openDialog = () => {
  Object.assign(form, { customer_id: 1, amount: 0, due_date: '', remark: '' })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await financeApi.createReceivable(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const openReceive = (row: any) => {
  receiveRow.value = row
  Object.assign(receiveForm, { amount: 0, pay_method: '转账', pay_date: '' })
  receiveVisible.value = true
}

const doReceive = async () => {
  try {
    await financeApi.receive({
      ref_id: receiveRow.value.id,
      amount: receiveForm.amount,
      pay_method: receiveForm.pay_method,
      pay_date: receiveForm.pay_date
    })
    ElMessage.success('收款成功')
    receiveVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '收款失败')
  }
}

onMounted(loadList)
</script>
