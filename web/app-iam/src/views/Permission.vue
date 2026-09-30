<template>
  <div style="padding:20px">
    <el-card>
      <template #header>
        <el-menu mode="horizontal" :default-active="$route.path" router style="border:none">
          <el-menu-item index="/users">用户管理</el-menu-item>
          <el-menu-item index="/roles">角色管理</el-menu-item>
          <el-menu-item index="/permissions">权限策略</el-menu-item>
          <el-menu-item index="/audit">审计日志</el-menu-item>
        </el-menu>
      </template>

      <el-row :gutter="20">
        <el-col :span="8">
          <h4>权限列表</h4>
          <el-table :data="permissions" border stripe size="small" max-height="500">
            <el-table-column prop="perm_code" label="权限编码" show-overflow-tooltip />
            <el-table-column prop="perm_name" label="名称" width="120" />
          </el-table>
        </el-col>
        <el-col :span="16">
          <h4>角色权限策略</h4>
          <el-form :inline="true">
            <el-form-item label="选择角色">
              <el-select v-model="selectedRole" placeholder="请选择角色" @change="loadPolicies" style="width:200px">
                <el-option v-for="r in roles" :key="r.id" :label="r.role_name" :value="r.role_code" />
              </el-select>
            </el-form-item>
            <el-button type="primary" @click="openAddDialog">添加策略</el-button>
          </el-form>

          <el-table :data="policies" border stripe size="small" style="margin-top:12px">
            <el-table-column type="index" label="#" width="50" />
            <el-table-column prop="0" label="角色" width="120" />
            <el-table-column prop="1" label="权限/资源" show-overflow-tooltip />
            <el-table-column prop="2" label="操作" width="100" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="removePolicy(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-col>
      </el-row>
    </el-card>

    <el-dialog v-model="addDialogVisible" title="添加策略" width="450px">
      <el-form :model="addForm" label-width="100px">
        <el-form-item label="角色">
          <el-input :model-value="selectedRole" disabled />
        </el-form-item>
        <el-form-item label="权限">
          <el-select v-model="addForm.perm_code" placeholder="选择权限" style="width:100%">
            <el-option v-for="p in permissions" :key="p.id" :label="`${p.perm_name} (${p.perm_code})`" :value="p.perm_code" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作">
          <el-select v-model="addForm.action" style="width:100%">
            <el-option label="读取 (read)" value="read" />
            <el-option label="写入 (write)" value="write" />
            <el-option label="全部 (*)" value="*" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="addPolicy">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { permApi, roleApi } from '../api/iam'

const permissions = ref<any[]>([])
const roles = ref<any[]>([])
const selectedRole = ref('')
const policies = ref<string[][]>([])
const addDialogVisible = ref(false)
const addForm = reactive({ perm_code: '', action: 'read' })

const loadPermissions = async () => {
  const res: any = await permApi.list()
  permissions.value = res.data || []
}

const loadRoles = async () => {
  const res: any = await roleApi.list({ page: 1, page_size: 100 })
  roles.value = res.data?.items || []
}

const loadPolicies = async () => {
  if (!selectedRole.value) return
  const res: any = await permApi.listPolicies(selectedRole.value)
  policies.value = res.data || []
}

const openAddDialog = () => {
  if (!selectedRole.value) {
    ElMessage.warning('请先选择角色')
    return
  }
  addForm.perm_code = ''
  addForm.action = 'read'
  addDialogVisible.value = true
}

const addPolicy = async () => {
  await permApi.addPolicy({ role_code: selectedRole.value, perm_code: addForm.perm_code, action: addForm.action })
  ElMessage.success('添加成功')
  addDialogVisible.value = false
  loadPolicies()
}

const removePolicy = async (row: string[]) => {
  await ElMessageBox.confirm('确认删除该策略？', '提示', { type: 'warning' })
  await permApi.removePolicy({ role_code: row[0], perm_code: row[1], action: row[2] })
  ElMessage.success('删除成功')
  loadPolicies()
}

onMounted(() => { loadPermissions(); loadRoles() })
</script>
