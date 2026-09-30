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
          <el-button type="primary" @click="openDialog">新建BOM</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="bom_name" label="BOM名称" width="160" />
        <el-table-column prop="product_sku" label="成品SKU" width="140" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 2 ? 'success' : 'info'">
              {{ row.status === 1 ? '草稿' : '已启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="子件明细" min-width="300">
          <template #default="{ row }">
            <el-tag v-for="item in row.items" :key="item.id" style="margin:2px" type="warning">
              {{ item.component_sku }} × {{ item.quantity }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button v-if="row.status === 1" size="small" type="primary" @click="enable(row)">启用</el-button>
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

    <el-dialog v-model="dialogVisible" title="新建BOM" width="550px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="BOM名称">
          <el-input v-model="form.bom_name" />
        </el-form-item>
        <el-form-item label="成品SKU">
          <el-input v-model="form.product_sku" />
        </el-form-item>
        <el-form-item label="子件明细">
          <el-table :data="form.items" border size="small">
            <el-table-column label="子件SKU" width="160">
              <template #default="{ row }">
                <el-input v-model="row.component_sku" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="用量" width="120">
              <template #default="{ row }">
                <el-input-number v-model="row.quantity" size="small" :min="1" controls-position="right" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="70">
              <template #default="{ $index }">
                <el-button size="small" type="danger" link @click="form.items.splice($index, 1)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button size="small" style="margin-top:8px" @click="form.items.push({ component_sku: '', quantity: 1 })">添加子件</el-button>
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
import { ElMessage } from 'element-plus'
import { productionApi } from '../api/production'

const list = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20 })

const dialogVisible = ref(false)
const form = reactive({ bom_name: '', product_sku: '', items: [{ component_sku: '', quantity: 1 }] })

const loadList = async () => {
  const res: any = await productionApi.listBoms(query)
  list.value = res.items || []
  total.value = res.total || 0
}

const openDialog = () => {
  form.bom_name = ''
  form.product_sku = ''
  form.items = [{ component_sku: '', quantity: 1 }]
  dialogVisible.value = true
}

const save = async () => {
  try {
    await productionApi.createBom(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '创建失败')
  }
}

const enable = async (row: any) => {
  try {
    await productionApi.enableBom(row.id)
    ElMessage.success('已启用')
    loadList()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '启用失败')
  }
}

onMounted(loadList)
</script>
