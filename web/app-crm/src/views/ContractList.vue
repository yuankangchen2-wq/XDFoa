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
          <el-button type="primary" @click="openContractDialog">新建合同</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:120px">
            <el-option label="草稿" :value="1" />
            <el-option label="生效" :value="3" />
            <el-option label="归档" :value="4" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="contract_no" label="合同编号" width="200" />
        <el-table-column prop="name" label="合同名称" width="160" />
        <el-table-column prop="customer_id" label="客户ID" width="80" />
        <el-table-column prop="amount" label="合同金额" width="120">
          <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="已回款" width="120">
          <template #default="{ row }">¥{{ row.paid_amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="回款进度" width="160">
          <template #default="{ row }">
            <el-progress :percentage="row.amount > 0 ? Math.round(row.paid_amount / row.amount * 100) : 0" :stroke-width="12" />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sign_date" label="签订日期" width="120" />
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button v-if="row.status === 1" size="small" type="primary" @click="approve(row)">审批</el-button>
            <el-button v-if="row.status === 3" size="small" type="success" @click="openPaymentDialog(row)">回款</el-button>
            <el-button v-if="row.status === 3" size="small" type="info" @click="archive(row)">归档</el-button>
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

    <!-- 新建合同弹窗 -->
    <el-dialog v-model="contractDialogVisible" title="新建合同" width="500px">
      <el-form :model="contractForm" label-width="100px">
        <el-form-item label="合同名称">
          <el-input v-model="contractForm.name" />
        </el-form-item>
        <el-form-item label="客户ID">
          <el-input-number v-model="contractForm.customer_id" :min="1" />
        </el-form-item>
        <el-form-item label="合同金额">
          <el-input-number v-model="contractForm.amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="签订日期">
          <el-input v-model="contractForm.sign_date" placeholder="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="contractForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="contractDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveContract">保存</el-button>
      </template>
    </el-dialog>

    <!-- 回款弹窗 -->
    <el-dialog v-model="paymentDialogVisible" title="添加回款" width="450px">
      <el-form :model="paymentForm" label-width="100px">
        <el-form-item label="合同名称">{{ currentContract?.name }}</el-form-item>
        <el-form-item label="合同金额">¥{{ currentContract?.amount?.toFixed(2) }}</el-form-item>
        <el-form-item label="已回款">¥{{ currentContract?.paid_amount?.toFixed(2) }}</el-form-item>
        <el-form-item label="回款金额">
          <el-input-number v-model="paymentForm.amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="回款日期">
          <el-input v-model="paymentForm.pay_date" placeholder="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="付款方式">
          <el-select v-model="paymentForm.pay_method" style="width:100%">
            <el-option label="银行转账" value="银行转账" />
            <el-option label="支付宝" value="支付宝" />
            <el-option label="微信" value="微信" />
            <el-option label="现金" value="现金" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="paymentForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="paymentDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="savePayment">确认回款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { salesApi } from '../api/sales'

const statusLabel: Record<number, string> = { 1: '草稿', 2: '审批中', 3: '生效', 4: '归档' }
const statusType: Record<number, string> = { 1: 'info', 2: 'warning', 3: 'success', 4: '' }

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, status: null as any })

const contractDialogVisible = ref(false)
const contractForm = reactive({ name: '', customer_id: 1, amount: 0, sign_date: '', remark: '' })

const paymentDialogVisible = ref(false)
const currentContract = ref<any>(null)
const paymentForm = reactive({ amount: 0, pay_date: '', pay_method: '银行转账', remark: '' })

const loadList = async () => {
  const res: any = await salesApi.listContracts(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openContractDialog = () => {
  Object.assign(contractForm, { name: '', customer_id: 1, amount: 0, sign_date: '', remark: '' })
  contractDialogVisible.value = true
}

const saveContract = async () => {
  try {
    await salesApi.createContract(contractForm)
    ElMessage.success('创建成功')
    contractDialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const approve = async (row: any) => {
  await ElMessageBox.confirm('确认审批该合同？', '提示', { type: 'warning' })
  try { await salesApi.approveContract(row.id); ElMessage.success('审批通过，合同已生效'); loadList() }
  catch (e: any) { ElMessage.error(e.response?.data?.error || '审批失败') }
}

const archive = async (row: any) => {
  try { await salesApi.archiveContract(row.id); ElMessage.success('已归档'); loadList() }
  catch (e: any) { ElMessage.error(e.response?.data?.error || '归档失败') }
}

const openPaymentDialog = (row: any) => {
  currentContract.value = row
  Object.assign(paymentForm, { amount: 0, pay_date: new Date().toISOString().slice(0, 10), pay_method: '银行转账', remark: '' })
  paymentDialogVisible.value = true
}

const savePayment = async () => {
  try {
    await salesApi.addPayment(currentContract.value.id, paymentForm)
    ElMessage.success('回款成功')
    paymentDialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '回款失败')
  }
}

onMounted(loadList)
</script>
