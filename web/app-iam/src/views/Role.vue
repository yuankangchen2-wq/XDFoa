<template>
  <div style="padding:20px">
    <el-card>
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <el-menu mode="horizontal" :default-active="$route.path" router style="border:none">
            <el-menu-item index="/users">用户管理</el-menu-item>
            <el-menu-item index="/roles">角色管理</el-menu-item>
            <el-menu-item index="/permissions">权限策略</el-menu-item>
            <el-menu-item index="/audit">审计日志</el-menu-item>
          </el-menu>
          <el-button type="primary" @click="openDialog()">新增角色</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="role_code" label="角色编码" width="160" />
        <el-table-column prop="role_name" label="角色名称" width="160" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="data_scope" label="数据范围" width="140">
          <template #default="{ row }">
            <el-tag>{{ scopeMap[row.data_scope] || row.data_scope }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑角色' : '新增角色'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="角色编码">
          <el-input v-model="form.role_code" />
        </el-form-item>
        <el-form-item label="角色名称">
          <el-input v-model="form.role_name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="数据范围">
          <el-select v-model="form.data_scope" style="width:100%">
            <el-option label="仅本人" value="self" />
            <el-option label="本部门" value="dept" />
            <el-option label="本部门及下级" value="dept_and_sub" />
            <el-option label="全部" value="all" />
          </el-select>
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
import { roleApi } from '../api/iam'

const scopeMap: Record<string, string> = {
  self: '仅本人',
  dept: '本部门',
  dept_and_sub: '本部门及下级',
  all: '全部'
}

const list = ref<any[]>([])
const dialogVisible = ref(false)
const form = reactive<any>({})

const loadList = async () => {
  const res: any = await roleApi.list({ page: 1, page_size: 100 })
  list.value = res.data?.items || []
}

const openDialog = (row?: any) => {
  Object.assign(form, row || { role_code: '', role_name: '', description: '', data_scope: 'self' })
  dialogVisible.value = true
}

const save = async () => {
  await roleApi.create(form)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadList()
}

const remove = async (row: any) => {
  await ElMessageBox.confirm('确认删除该角色？', '提示', { type: 'warning' })
  await roleApi.remove(row.id)
  ElMessage.success('删除成功')
  loadList()
}

onMounted(loadList)
</script>
