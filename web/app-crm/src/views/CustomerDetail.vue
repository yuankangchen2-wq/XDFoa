<template>
  <div style="padding:20px;max-width:1100px;margin:0 auto">
    <el-page-header @back="$router.back()" title="返回列表" />
    <el-card style="margin-top:16px" v-if="customer">
      <template #header>
        <div style="display:flex;justify-content:space-between">
          <h3 style="margin:0">{{ customer.name }}</h3>
          <el-tag :type="statusType[customer.status]">{{ statusLabel[customer.status] }}</el-tag>
        </div>
      </template>

      <el-descriptions :column="3" border>
        <el-descriptions-item label="联系人">{{ customer.contact_person }}</el-descriptions-item>
        <el-descriptions-item label="电话">{{ customer.phone }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ customer.email }}</el-descriptions-item>
        <el-descriptions-item label="分级">
          <el-tag :type="levelType[customer.level]">{{ customer.level }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="行业">{{ customer.industry }}</el-descriptions-item>
        <el-descriptions-item label="归属销售">{{ customer.owner_id }}</el-descriptions-item>
        <el-descriptions-item label="地址" :span="3">{{ customer.address }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="3">{{ customer.remark }}</el-descriptions-item>
      </el-descriptions>

      <!-- 联系人 -->
      <div style="margin-top:24px">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
          <h4 style="margin:0">联系人</h4>
          <el-button size="small" type="primary" @click="openContactDialog">添加联系人</el-button>
        </div>
        <el-table :data="customer.contacts" border size="small">
          <el-table-column prop="name" label="姓名" width="120" />
          <el-table-column prop="position" label="职位" width="120" />
          <el-table-column prop="phone" label="电话" width="140" />
          <el-table-column prop="email" label="邮箱" width="180" />
          <el-table-column label="主联系人" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.is_primary" type="success">是</el-tag>
              <span v-else>否</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80">
            <template #default="{ row }">
              <el-button size="small" type="danger" link @click="deleteContact(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 跟进记录 -->
      <div style="margin-top:24px">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
          <h4 style="margin:0">跟进记录</h4>
          <el-button size="small" type="primary" @click="openFollowUpDialog">添加跟进</el-button>
        </div>
        <el-timeline>
          <el-timeline-item
            v-for="f in customer.follow_ups"
            :key="f.id"
            :timestamp="f.follow_up_at"
            placement="top"
          >
            <el-card shadow="never">
              <div style="display:flex;justify-content:space-between">
                <el-tag>{{ f.type }}</el-tag>
                <span style="color:#999">{{ f.user_id }}</span>
              </div>
              <p style="margin:8px 0 0">{{ f.content }}</p>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-card>

    <!-- 添加联系人弹窗 -->
    <el-dialog v-model="contactDialogVisible" title="添加联系人" width="450px">
      <el-form :model="contactForm" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="contactForm.name" />
        </el-form-item>
        <el-form-item label="职位">
          <el-input v-model="contactForm.position" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="contactForm.phone" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="contactForm.email" />
        </el-form-item>
        <el-form-item label="主联系人">
          <el-switch v-model="contactForm.is_primary" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="contactDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveContact">保存</el-button>
      </template>
    </el-dialog>

    <!-- 添加跟进弹窗 -->
    <el-dialog v-model="followUpDialogVisible" title="添加跟进记录" width="500px">
      <el-form :model="followUpForm" label-width="80px">
        <el-form-item label="类型">
          <el-select v-model="followUpForm.type" style="width:100%">
            <el-option label="电话" value="电话" />
            <el-option label="拜访" value="拜访" />
            <el-option label="邮件" value="邮件" />
            <el-option label="微信" value="微信" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="followUpForm.content" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="跟进人">
          <el-input v-model="followUpForm.user_id" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="followUpDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveFollowUp">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { crmApi } from '../api/crm'

const route = useRoute()
const customer = ref<any>(null)

const levelType: Record<string, string> = { A: 'danger', B: 'warning', C: '', D: 'info' }
const statusLabel: Record<number, string> = { 1: '潜在', 2: '意向', 3: '成交', 4: '流失' }
const statusType: Record<number, string> = { 1: 'info', 2: 'warning', 3: 'success', 4: 'danger' }

const contactDialogVisible = ref(false)
const contactForm = reactive({ name: '', position: '', phone: '', email: '', is_primary: 0 })

const followUpDialogVisible = ref(false)
const followUpForm = reactive({ type: '电话', content: '', user_id: '' })

const loadCustomer = async () => {
  const id = Number(route.params.id)
  customer.value = await crmApi.getCustomer(id)
}

const openContactDialog = () => {
  Object.assign(contactForm, { name: '', position: '', phone: '', email: '', is_primary: 0 })
  contactDialogVisible.value = true
}

const saveContact = async () => {
  try {
    await crmApi.addContact(customer.value.id, contactForm)
    ElMessage.success('添加成功')
    contactDialogVisible.value = false
    loadCustomer()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '添加失败')
  }
}

const deleteContact = async (row: any) => {
  await ElMessageBox.confirm('确认删除该联系人？', '提示', { type: 'warning' })
  await crmApi.deleteContact(row.id)
  ElMessage.success('已删除')
  loadCustomer()
}

const openFollowUpDialog = () => {
  Object.assign(followUpForm, { type: '电话', content: '', user_id: '' })
  followUpDialogVisible.value = true
}

const saveFollowUp = async () => {
  try {
    await crmApi.addFollowUp(customer.value.id, { ...followUpForm, follow_up_at: new Date().toISOString() })
    ElMessage.success('添加成功')
    followUpDialogVisible.value = false
    loadCustomer()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '添加失败')
  }
}

onMounted(loadCustomer)
</script>
