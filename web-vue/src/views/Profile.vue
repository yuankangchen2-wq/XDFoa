<template>
  <div class="profile-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="never" class="profile-card">
          <div class="avatar-section">
            <el-avatar :size="100" :src="user?.avatar">{{ user?.name?.charAt(0) || 'U' }}</el-avatar>
            <h2>{{ user?.displayName || user?.name || '用户' }}</h2>
            <p>{{ user?.email || '未设置邮箱' }}</p>
          </div>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">{{ user?.name }}</el-descriptions-item>
            <el-descriptions-item label="姓名">{{ user?.displayName }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ user?.email || '-' }}</el-descriptions-item>
            <el-descriptions-item label="手机">{{ user?.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="组织">{{ user?.owner || 'built-in' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card shadow="never">
          <template #header><span>基本信息</span></template>
          <el-form :model="form" label-width="100px" style="max-width: 500px">
            <el-form-item label="姓名">
              <el-input v-model="form.displayName" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="form.email" />
            </el-form-item>
            <el-form-item label="手机">
              <el-input v-model="form.phone" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="save">保存修改</el-button>
              <el-button @click="changePassword">修改密码</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const user = computed(() => userStore.userInfo)

const form = reactive({
  displayName: '',
  email: '',
  phone: ''
})

onMounted(() => {
  if (user.value) {
    form.displayName = user.value.displayName || ''
    form.email = user.value.email || ''
    form.phone = user.value.phone || ''
  }
})

function save() {
  ElMessage.success('保存成功（演示）')
}

function changePassword() {
  ElMessage.info('修改密码功能开发中')
}
</script>

<style scoped>
.profile-card {
  text-align: center;
}
.avatar-section {
  margin-bottom: 20px;
}
.avatar-section h2 {
  margin: 12px 0 4px;
}
.avatar-section p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}
</style>
