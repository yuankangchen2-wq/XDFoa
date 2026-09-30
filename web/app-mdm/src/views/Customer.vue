<template>
  <div>
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <span>客户管理</span>
          <el-button type="primary" @click="openDialog()">新增客户</el-button>
        </div>
      </template>

      <el-form :inline="true" @submit.prevent>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="客户名称/编码" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="query.customer_type" placeholder="全部" clearable style="width:120px">
            <el-option label="企业" value="enterprise" />
            <el-option label="个人" value="individual" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="customer_code" label="客户编码" width="120" />
        <el-table-column prop="customer_name" label="客户名称" />
        <el-table-column prop="customer_type" label="类型" width="100" />
        <el-table-column prop="industry" label="行业" width="120" />
        <el-table-column prop="contact_name" label="联系人" width="100" />
        <el-table-column prop="contact_phone" label="电话" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.page_size"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="loadList"
        @size-change="loadList"
        style="margin-top:16px;justify-content:flex-end"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑客户' : '新增客户'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="客户编码">
          <el-input v-model="form.customer_code" />
        </el-form-item>
        <el-form-item label="客户名称">
          <el-input v-model="form.customer_name" />
        </el-form-item>
        <el-form-item label="客户类型">
          <el-select v-model="form.customer_type" style="width:100%">
            <el-option label="企业" value="enterprise" />
            <el-option label="个人" value="individual" />
          </el-select>
        </el-form-item>
        <el-form-item label="行业">
          <el-input v-model="form.industry" />
        </el-form-item>
        <el-form-item label="区域">
          <el-input v-model="form.region" />
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="form.contact_name" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="form.contact_phone" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.contact_email" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { customerApi } from '../api/mdm'

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, keyword: '', customer_type: '' })
const dialogVisible = ref(false)
const form = reactive<any>({})

const loadList = async () => {
  const res: any = await customerApi.list(query)
  list.value = res.data?.items || []
  total.value = res.data?.total || 0
}

const openDialog = (row?: any) => {
  Object.assign(form, row || {})
  dialogVisible.value = true
}

const save = async () => {
  if (form.id) {
    await customerApi.update(form.id, form)
    ElMessage.success('更新成功')
  } else {
    await customerApi.create(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  loadList()
}

const remove = async (row: any) => {
  await ElMessageBox.confirm('确认删除该客户？', '提示', { type: 'warning' })
  await customerApi.remove(row.id)
  ElMessage.success('删除成功')
  loadList()
}

onMounted(loadList)
</script>
