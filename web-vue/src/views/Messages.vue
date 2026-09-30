<template>
  <div class="messages-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>消息中心</span>
          <el-button type="primary" text @click="markAllRead">全部已读</el-button>
        </div>
      </template>

      <div v-for="msg in messages" :key="msg.id" class="message-item" :class="{ unread: !msg.read }" @click="msg.read = true">
        <el-avatar :size="40" :style="{ background: msg.color }">{{ msg.from?.charAt(0) || '系' }}</el-avatar>
        <div class="message-content">
          <div class="message-header">
            <span class="message-from">{{ msg.from }}</span>
            <span class="message-time">{{ msg.time }}</span>
          </div>
          <div class="message-title">{{ msg.title }}</div>
          <div class="message-body">{{ msg.content }}</div>
        </div>
        <el-badge v-if="!msg.read" :value="1" type="danger" />
      </div>

      <el-empty v-if="messages.length === 0" description="暂无消息" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const messages = ref([
  { id: 1, from: '系统通知', title: '欢迎使用统一门户', content: '您已成功登录统一门户系统', time: '10分钟前', color: '#409EFF', read: false },
  { id: 2, from: 'OA系统', title: '您有一条新的待办', content: '张三提交了请假申请，请审批', time: '1小时前', color: '#67C23A', read: false },
  { id: 3, from: 'ERP系统', title: '采购订单提醒', content: '采购订单 #20240926 待您审批', time: '2小时前', color: '#E6A23C', read: false },
  { id: 4, from: 'HR系统', title: '考勤异常提醒', content: '您昨日考勤异常，请及时处理', time: '昨天', color: '#F56C6C', read: true }
])

function markAllRead() {
  messages.value.forEach(m => m.read = true)
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.message-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background 0.2s;
}
.message-item:hover {
  background: #f5f7fa;
}
.message-item.unread {
  background: #ecf5ff;
}
.message-content {
  flex: 1;
  overflow: hidden;
}
.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}
.message-from {
  font-weight: 600;
  color: #303133;
}
.message-time {
  font-size: 12px;
  color: #909399;
}
.message-title {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
}
.message-body {
  font-size: 13px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
