<template>
  <div style="padding:20px;max-width:900px;margin:0 auto">
    <el-page-header @back="$router.back()" title="返回列表" />
    <el-card style="margin-top:16px" v-if="ticket">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <div>
            <span style="font-size:18px;font-weight:bold">{{ ticket.title }}</span>
            <el-tag style="margin-left:12px" :type="statusType[ticket.status]">{{ statusLabel[ticket.status] }}</el-tag>
            <el-tag style="margin-left:8px" :type="priorityType[ticket.priority]">{{ priorityLabel[ticket.priority] }}</el-tag>
          </div>
          <span style="color:#999;font-size:13px">{{ ticket.ticket_no }}</span>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="客户ID">{{ ticket.customer_id }}</el-descriptions-item>
        <el-descriptions-item label="处理人">{{ ticket.assignee_id || '未派单' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ ticket.created_at }}</el-descriptions-item>
        <el-descriptions-item label="满意度">
          <el-rate v-if="ticket.satisfaction > 0" v-model="ticket.satisfaction" disabled show-score text-color="#ff9900" score-template="{value}" />
          <span v-else>未评价</span>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ ticket.description }}</el-descriptions-item>
      </el-descriptions>

      <!-- 操作栏 -->
      <div style="margin:16px 0;display:flex;gap:8px;flex-wrap:wrap">
        <el-input v-model="assigneeId" placeholder="输入处理人ID" style="width:200px" />
        <el-button type="primary" @click="assign">派单</el-button>
        <el-button v-if="ticket.status === 1 || ticket.status === 2" type="success" @click="updateStatus(3)">标记解决</el-button>
        <el-button v-if="ticket.status === 3" type="info" @click="updateStatus(4)">关闭工单</el-button>
        <el-select v-model="rateScore" placeholder="满意度评分" style="width:140px">
          <el-option v-for="n in 5" :key="n" :label="`${n} 星`" :value="n" />
        </el-select>
        <el-button v-if="ticket.status === 3" type="warning" @click="rate">提交评价</el-button>
      </div>

      <!-- 回复区 -->
      <h4 style="margin:20px 0 12px">处理记录</h4>
      <el-timeline>
        <el-timeline-item
          v-for="r in ticket.replies"
          :key="r.id"
          :timestamp="r.created_at"
          placement="top"
        >
          <el-card shadow="never">
            <div style="color:#409eff;font-weight:bold">{{ r.user_id }}</div>
            <p style="margin:8px 0 0">{{ r.content }}</p>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <!-- 添加回复 -->
      <div style="margin-top:16px">
        <el-input v-model="replyContent" type="textarea" :rows="3" placeholder="输入回复内容..." />
        <div style="margin-top:8px;text-align:right">
          <el-button type="primary" @click="addReply">发送回复</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ticketApi } from '../api/ticket'

const route = useRoute()
const ticket = ref<any>(null)
const assigneeId = ref('')
const replyContent = ref('')
const rateScore = ref(5)

const statusLabel: Record<number, string> = { 1: '待处理', 2: '处理中', 3: '已解决', 4: '已关闭' }
const statusType: Record<number, string> = { 1: 'warning', 2: '', 3: 'success', 4: 'info' }
const priorityLabel: Record<number, string> = { 1: '低', 2: '中', 3: '高', 4: '紧急' }
const priorityType: Record<number, string> = { 1: 'info', 2: '', 3: 'warning', 4: 'danger' }

const loadTicket = async () => {
  const id = Number(route.params.id)
  ticket.value = await ticketApi.get(id)
}

const assign = async () => {
  if (!assigneeId.value) {
    ElMessage.warning('请输入处理人ID')
    return
  }
  try {
    await ticketApi.assign(ticket.value.id, assigneeId.value)
    ElMessage.success('派单成功')
    loadTicket()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '派单失败')
  }
}

const updateStatus = async (status: number) => {
  try {
    await ticketApi.updateStatus(ticket.value.id, status)
    ElMessage.success('状态已更新')
    loadTicket()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  }
}

const addReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  try {
    await ticketApi.addReply(ticket.value.id, { user_id: 'current_user', content: replyContent.value })
    ElMessage.success('回复成功')
    replyContent.value = ''
    loadTicket()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '回复失败')
  }
}

const rate = async () => {
  try {
    await ticketApi.rate(ticket.value.id, rateScore.value)
    ElMessage.success('评价成功')
    loadTicket()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '评价失败')
  }
}

onMounted(loadTicket)
</script>
