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
          <el-button type="primary" @click="openCreateDialog">新建客户</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="分级">
          <el-select v-model="query.level" placeholder="全部" clearable style="width:120px">
            <el-option label="A 级" value="A" />
            <el-option label="B 级" value="B" />
            <el-option label="C 级" value="C" />
            <el-option label="D 级" value="D" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:120px">
            <el-option label="潜在" :value="1" />
            <el-option label="意向" :value="2" />
            <el-option label="成交" :value="3" />
            <el-option label="流失" :value="4" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="name" label="客户名称" width="160" />
        <el-table-column prop="contact_person" label="联系人" width="100" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column prop="level" label="分级" width="80">
          <template #default="{ row }">
            <el-tag :type="levelType[row.level]">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="industry" label="行业" width="100" />
        <el-table-column prop="owner_id" label="归属销售" width="100" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="$router.push(`/customers/${row.id}`)">详情</el-button>
            <el-button size="small" type="danger" @click="del(row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" title="新建客户" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="客户名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="form.contact_person" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" />
        </el-form-item>
        <el-form-item label="分级">
          <el-select v-model="form.level" style="width:100%">
            <el-option label="A" value="A" />
            <el-option label="B" value="B" />
            <el-option label="C" value="C" />
            <el-option label="D" value="D" />
          </el-select>
        </el-form-item>
        <el-form-item label="行业">
          <el-input v-model="form.industry" />
        </el-form-item>
        <el-form-item label="归属销售">
          <el-input v-model="form.owner_id" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { crmApi } from '../api/crm'

const levelType: Record<string, string> = { A: 'danger', B: 'warning', C: '', D: 'info' }
const statusLabel: Record<number, string> = { 1: '潜在', 2: '意向', 3: '成交', 4: '流失' }
const statusType: Record<number, string> = { 1: 'info', 2: 'warning', 3: 'success', 4: 'danger' }

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, level: '', status: null as any })

const dialogVisible = ref(false)
const form = reactive({
  name: '', contact_person: '', phone: '', email: '', address: '',
  level: 'B', industry: '', owner_id: '', remark: ''
})

const loadList = async () => {
  const res: any = await crmApi.listCustomers(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openCreateDialog = () => {
  Object.assign(form, { name: '', contact_person: '', phone: '', email: '', address: '', level: 'B', industry: '', owner_id: '', remark: '' })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await crmApi.createCustomer(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const del = async (row: any) => {
  await ElMessageBox.confirm(`确认删除客户"${row.name}"？`, '提示', { type: 'warning' })
  try {
    await crmApi.deleteCustomer(row.id)
    ElMessage.success('已删除')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '删除失败')
  }
}

onMounted(loadList)
</script>
