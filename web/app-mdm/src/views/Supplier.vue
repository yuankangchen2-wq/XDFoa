<template>
  <div>
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <span>供应商管理</span>
          <el-button type="primary" @click="openDialog()">新增供应商</el-button>
        </div>
      </template>

      <el-form :inline="true" @submit.prevent>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="供应商名称/编码" clearable />
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="supplier_code" label="供应商编码" width="140" />
        <el-table-column prop="supplier_name" label="供应商名称" />
        <el-table-column prop="contact_name" label="联系人" width="100" />
        <el-table-column prop="contact_phone" label="电话" width="130" />
        <el-table-column prop="contact_email" label="邮箱" width="160" />
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑供应商' : '新增供应商'" width="550px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="供应商编码">
          <el-input v-model="form.supplier_code" />
        </el-form-item>
        <el-form-item label="供应商名称">
          <el-input v-model="form.supplier_name" />
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
import { supplierApi } from '../api/mdm'

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, keyword: '' })
const dialogVisible = ref(false)
const form = reactive<any>({})

const loadList = async () => {
  const res: any = await supplierApi.list(query)
  list.value = res.data?.items || []
  total.value = res.data?.total || 0
}

const openDialog = (row?: any) => {
  Object.assign(form, row || {})
  dialogVisible.value = true
}

const save = async () => {
  if (form.id) {
    await supplierApi.update(form.id, form)
    ElMessage.success('更新成功')
  } else {
    await supplierApi.create(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  loadList()
}

const remove = async (row: any) => {
  await ElMessageBox.confirm('确认删除该供应商？', '提示', { type: 'warning' })
  await supplierApi.remove(row.id)
  ElMessage.success('删除成功')
  loadList()
}

onMounted(loadList)
</script>
