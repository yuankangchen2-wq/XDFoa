<template>
  <div>
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <span>商品管理</span>
          <el-button type="primary" @click="openDialog()">新增商品</el-button>
        </div>
      </template>

      <el-form :inline="true" @submit.prevent>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="商品名称/SKU" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="query.category" placeholder="分类" clearable />
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="sku_code" label="SKU编码" width="140" />
        <el-table-column prop="product_name" label="商品名称" />
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="spec" label="规格" width="120" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '上架' : '下架' }}
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑商品' : '新增商品'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="SKU编码">
          <el-input v-model="form.sku_code" />
        </el-form-item>
        <el-form-item label="商品名称">
          <el-input v-model="form.product_name" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="form.category" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="form.unit" />
        </el-form-item>
        <el-form-item label="规格">
          <el-input v-model="form.spec" />
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
import { productApi } from '../api/mdm'

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, keyword: '', category: '' })
const dialogVisible = ref(false)
const form = reactive<any>({})

const loadList = async () => {
  const res: any = await productApi.list(query)
  list.value = res.data?.items || []
  total.value = res.data?.total || 0
}

const openDialog = (row?: any) => {
  Object.assign(form, row || {})
  dialogVisible.value = true
}

const save = async () => {
  if (form.id) {
    await productApi.update(form.id, form)
    ElMessage.success('更新成功')
  } else {
    await productApi.create(form)
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
  loadList()
}

const remove = async (row: any) => {
  await ElMessageBox.confirm('确认删除该商品？', '提示', { type: 'warning' })
  await productApi.remove(row.id)
  ElMessage.success('删除成功')
  loadList()
}

onMounted(loadList)
</script>
