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
          <el-button type="primary" @click="openDialog()">新增用户</el-button>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="用户名/姓名" clearable />
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="real_name" label="姓名" width="100" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column prop="org_path" label="组织路径" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '正常' : '禁用' }}</el-tag>
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
        layout="total, prev, pager, next"
        @current-change="loadList"
        style="margin-top:16px;justify-content:flex-end"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑用户' : '新增用户'" width="550px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="用户名" v-if="!form.id">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码" v-if="!form.id">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="form.real_name" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width:100%">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role_ids" multiple placeholder="选择角色" style="width:100%">
            <el-option v-for="r in roles" :key="r.id" :label="r.role_name" :value="r.id" />
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
import { userApi, roleApi } from '../api/iam'

const list = ref<any[]>([])
const total = ref(0)
const roles = ref<any[]>([])
const query = reactive({ page: 1, page_size: 20, keyword: '' })
const dialogVisible = ref(false)
const form = reactive<any>({ role_ids: [] })

const loadList = async () => {
  const res: any = await userApi.list(query)
  list.value = res.data?.items || []
  total.value = res.data?.total || 0
}

const loadRoles = async () => {
  const res: any = await roleApi.list({ page: 1, page_size: 100 })
  roles.value = res.data?.items || []
}

const openDialog = (row?: any) => {
  Object.assign(form, row || { username: '', password: '', real_name: '', email: '', phone: '', status: 1, role_ids: [] })
  dialogVisible.value = true
}

const save = async () => {
  if (form.id) {
    await userApi.update(form.id, form)
  } else {
    await userApi.create(form)
  }
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadList()
}

const remove = async (row: any) => {
  await ElMessageBox.confirm('确认删除该用户？', '提示', { type: 'warning' })
  await userApi.remove(row.id)
  ElMessage.success('删除成功')
  loadList()
}

onMounted(() => { loadList(); loadRoles() })
</script>
