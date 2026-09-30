<template>
  <div style="display:flex;justify-content:center;align-items:center;height:100vh;background:#f5f7fa">
    <el-card style="width:400px">
      <h2 style="text-align:center;margin-bottom:24px">统一登录</h2>
      <el-form :model="form" label-width="80px" @submit.prevent>
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="admin" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" placeholder="admin123" show-password />
        </el-form-item>
        <el-form-item label="租户">
          <el-input v-model="form.tenant_id" placeholder="default" />
        </el-form-item>
        <el-button type="primary" style="width:100%" @click="login">登 录</el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { authApi } from '../api/iam'

const router = useRouter()
const form = reactive({ username: 'admin', password: 'admin123', tenant_id: 'default' })

const login = async () => {
  try {
    const res: any = await authApi.login(form)
    if (res.code === 0) {
      localStorage.setItem('access_token', res.data.access_token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      ElMessage.success('登录成功')
      router.push('/users')
    } else {
      ElMessage.error(res.message)
    }
  } catch (e: any) {
    ElMessage.error(e.message || '登录失败')
  }
}
</script>
