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
          <el-button type="primary" @click="openDialog">新建成本中心</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="类型">
          <el-select v-model="query.type" placeholder="全部" clearable style="width:140px">
            <el-option label="部门" :value="1" />
            <el-option label="车间" :value="2" />
            <el-option label="其他" :value="3" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="编码" width="140" />
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ typeLabel[row.type] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" />
        <el-table-column prop="created_at" label="创建时间" width="200" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新建成本中心" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="编码">
          <el-input v-model="form.code" />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type" style="width:100%">
            <el-option label="部门" :value="1" />
            <el-option label="车间" :value="2" />
            <el-option label="其他" :value="3" />
          </el-select>
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

const typeLabel: Record<number, string> = { 1: '部门', 2: '车间', 3: '其他' }

const list = ref<any[]>([])
const query = reactive({ type: null as any })

const dialogVisible = ref(false)
const form = reactive({ code: '', name: '', type: 1, remark: '' })

const loadList = async () => {
  const res: any = await costApi.listCenters(query)
  list.value = res.items || []
}

const openDialog = () => {
  Object.assign(form, { code: '', name: '', type: 1, remark: '' })
  dialogVisible.value = true
}

const save = async () => {
  try {
    await costApi.createCenter(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

onMounted(loadList)
</script>
