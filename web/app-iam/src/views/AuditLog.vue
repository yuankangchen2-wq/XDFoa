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
          <el-button @click="loadStats">刷新统计</el-button>
        </div>
      </template>

      <!-- 模块统计 -->
      <el-row :gutter="12" style="margin-bottom:16px">
        <el-col :span="8" v-for="s in stats" :key="s.module">
          <el-card shadow="hover" style="text-align:center">
            <div style="color:#909399;font-size:13px">{{ s.module }}</div>
            <div style="font-size:24px;font-weight:bold;color:#409eff;margin-top:4px">{{ s.count }}</div>
          </el-card>
        </el-col>
        <el-col :span="8" v-if="!stats.length">
          <el-empty description="暂无统计数据" :image-size="60" />
        </el-col>
      </el-row>

      <!-- 筛选 -->
      <el-form :inline="true">
        <el-form-item label="模块">
          <el-input v-model="query.module" placeholder="如 user/order" style="width:140px" />
        </el-form-item>
        <el-form-item label="操作">
          <el-input v-model="query.action" placeholder="如 CREATE/LOGIN" style="width:150px" />
        </el-form-item>
        <el-form-item label="操作人">
          <el-input v-model="query.user_id" style="width:140px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:120px">
            <el-option label="成功" :value="1" />
            <el-option label="失败" :value="2" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="module" label="模块" width="110" />
        <el-table-column prop="action" label="操作" width="100" />
        <el-table-column prop="username" label="操作人" width="120" />
        <el-table-column prop="target_name" label="操作对象" min-width="140" />
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration_ms" label="耗时(ms)" width="100" />
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button size="small" @click="showDetail(row)">详情</el-button>
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

    <!-- 详情抽屉 -->
    <el-drawer v-model="detailVisible" title="审计详情" size="500px">
      <el-descriptions v-if="current" :column="1" border>
        <el-descriptions-item label="ID">{{ current.id }}</el-descriptions-item>
        <el-descriptions-item label="模块">{{ current.module }}</el-descriptions-item>
        <el-descriptions-item label="操作">{{ current.action }}</el-descriptions-item>
        <el-descriptions-item label="操作人ID">{{ current.user_id }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ current.username }}</el-descriptions-item>
        <el-descriptions-item label="对象ID">{{ current.target_id }}</el-descriptions-item>
        <el-descriptions-item label="对象名称">{{ current.target_name }}</el-descriptions-item>
        <el-descriptions-item label="IP">{{ current.ip }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ current.status === 1 ? '成功' : '失败' }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ current.duration_ms }} ms</el-descriptions-item>
        <el-descriptions-item label="时间">{{ current.created_at }}</el-descriptions-item>
        <el-descriptions-item label="请求参数">
          <pre style="white-space:pre-wrap;max-height:200px;overflow:auto">{{ current.request_params }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应结果">
          <pre style="white-space:pre-wrap;max-height:200px;overflow:auto">{{ current.result }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { auditApi } from '../api/audit'

const list = ref<any[]>([])
const stats = ref<any[]>([])
const total = ref(0)
const query = reactive({ page: 1, page_size: 20, module: '', action: '', user_id: '', status: null as any })

const detailVisible = ref(false)
const current = ref<any>(null)

const loadList = async () => {
  try {
    const res: any = await auditApi.list(query)
    list.value = res.items || []
    total.value = res.total || 0
  } catch {
    list.value = []; total.value = 0
  }
}

const loadStats = async () => {
  try {
    const res: any = await auditApi.moduleStats()
    stats.value = res.items || []
  } catch { stats.value = [] }
}

const showDetail = (row: any) => {
  current.value = row
  detailVisible.value = true
}

onMounted(() => { loadList(); loadStats() })
</script>
