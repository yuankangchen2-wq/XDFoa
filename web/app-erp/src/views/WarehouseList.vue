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
          <el-button type="primary" @click="openDialog">新增仓库</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="warehouse_code" label="仓库编码" width="160" />
        <el-table-column prop="warehouse_name" label="仓库名称" width="160" />
        <el-table-column prop="address" label="地址" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增仓库" width="450px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="仓库编码">
          <el-input v-model="form.warehouse_code" />
        </el-form-item>
        <el-form-item label="仓库名称">
          <el-input v-model="form.warehouse_name" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { warehouseApi } from '../api/inventory'

const list = ref<any[]>([])
const dialogVisible = ref(false)
const form = reactive({ warehouse_code: '', warehouse_name: '', address: '' })

const loadList = async () => {
  const res: any = await warehouseApi.list({ page: 1, page_size: 100 })
  list.value = res.items || []
}

const openDialog = () => { dialogVisible.value = true }

const save = async () => {
  try {
    await warehouseApi.create(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

onMounted(loadList)
</script>
