<template>
  <div class="todos-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>待办中心</span>
          <el-radio-group v-model="filter" size="small">
            <el-radio-button label="all">全部</el-radio-button>
            <el-radio-button label="pending">待办</el-radio-button>
            <el-radio-button label="done">已办</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <el-table :data="filteredTodos" style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="300">
          <template #default="{ row }">
            <span class="todo-title" @click="handleTodo(row)">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源系统" width="120" />
        <el-table-column label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="row.priority === 'high' ? 'danger' : row.priority === 'medium' ? 'warning' : 'info'" size="small">
              {{ row.priority === 'high' ? '高' : row.priority === 'medium' ? '中' : '低' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'pending' ? 'warning' : 'success'" size="small">
              {{ row.status === 'pending' ? '待办' : '已办' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="time" label="时间" width="160" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

const filter = ref('all')

const todos = ref([
  { id: 1, title: '审批：张三的请假申请', source: 'OA', priority: 'high', status: 'pending', time: '2024-09-26 09:00' },
  { id: 2, title: '审批：采购订单 #20240926', source: 'ERP', priority: 'medium', status: 'pending', time: '2024-09-25 14:30' },
  { id: 3, title: '确认：合同签署', source: 'CRM', priority: 'high', status: 'pending', time: '2024-09-25 10:00' },
  { id: 4, title: '审批：李四的报销单', source: '财务', priority: 'low', status: 'done', time: '2024-09-24 16:00' },
  { id: 5, title: '会议：项目周例会', source: 'OA', priority: 'medium', status: 'done', time: '2024-09-24 10:00' }
])

const filteredTodos = computed(() => {
  if (filter.value === 'all') return todos.value
  return todos.value.filter(t => t.status === filter.value)
})

function handleTodo(row: any) {
  ElMessage.info(`打开待办：${row.title}`)
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.todo-title {
  cursor: pointer;
  color: #409EFF;
}
.todo-title:hover {
  text-decoration: underline;
}
</style>
