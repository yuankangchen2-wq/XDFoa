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
          <el-button type="primary" @click="openCreateDialog">新建商机</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="阶段">
          <el-select v-model="query.stage" placeholder="全部" clearable style="width:120px">
            <el-option v-for="s in stages" :key="s.value" :label="s.label" :value="s.value" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="name" label="商机名称" width="160" />
        <el-table-column prop="customer_id" label="客户ID" width="80" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="阶段" width="100">
          <template #default="{ row }">
            <el-tag :type="stageType[row.stage]">{{ stageLabel[row.stage] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="赢率" width="120">
          <template #default="{ row }">
            <el-progress :percentage="row.win_rate" :stroke-width="12" />
          </template>
        </el-table-column>
        <el-table-column prop="expected_close_date" label="预计成交" width="120" />
        <el-table-column prop="owner_id" label="归属销售" width="100" />
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-dropdown @command="(cmd) => transition(row, cmd)" trigger="click">
              <el-button size="small" type="primary">阶段流转<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="s in stages" :key="s.value" :command="s.value" :disabled="row.stage === s.value">
                    转为{{ s.label }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button v-if="row.stage === 5" size="small" type="success" @click="convertOrder(row)">转订单</el-button>
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

    <el-dialog v-model="dialogVisible" title="新建商机" width="550px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="商机名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="客户ID">
          <el-input-number v-model="form.customer_id" :min="1" />
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="form.amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="预计成交日期">
          <el-input v-model="form.expected_close_date" placeholder="YYYY-MM-DD" />
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
import { ArrowDown } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { oppApi } from '../api/opportunity'

const stages = [
  { value: 1, label: '线索' },
  { value: 2, label: '意向' },
  { value: 3, label: '报价' },
  { value: 4, label: '谈判' },
  { value: 5, label: '成交' },
  { value: 6, label: '输单' }
]
const stageLabel: Record<number, string> = { 1: '线索', 2: '意向', 3: '报价', 4: '谈判', 5: '成交', 6: '输单' }
const stageType: Record<number, string> = { 1: 'info', 2: '', 3: 'warning', 4: 'warning', 5: 'success', 6: 'danger' }

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, stage: null as any })

const dialogVisible = ref(false)
const form = reactive({ name: '', customer_id: 1, amount: 0, expected_close_date: '', owner_id: '', remark: '' })

const loadList = async () => {
  const res: any = await oppApi.list(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openCreateDialog = () => {
  Object.assign(form, { name: '', customer_id: 1, amount: 0, expected_close_date: '', owner_id: '', remark: '' })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await oppApi.create(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const transition = async (row: any, stage: number) => {
  await ElMessageBox.confirm(`确认将商机转为"${stageLabel[stage]}"？`, '提示', { type: 'warning' })
  try {
    await oppApi.transition(row.id, { to_stage: stage })
    ElMessage.success('阶段已更新')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '流转失败')
  }
}

const convertOrder = async (row: any) => {
  try {
    await oppApi.convertOrder(row.id, { user_id: 1001, warehouse_id: 1 })
    ElMessage.success('已转为销售订单')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '转订单失败')
  }
}

const del = async (row: any) => {
  await ElMessageBox.confirm(`确认删除商机"${row.name}"？`, '提示', { type: 'warning' })
  try {
    await oppApi.remove(row.id)
    ElMessage.success('已删除')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '删除失败')
  }
}

onMounted(loadList)
</script>
