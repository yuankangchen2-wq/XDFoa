<template>
  <el-popover
    :visible="visible"
    placement="bottom-end"
    :width="400"
    trigger="manual"
    popper-class="notify-pop"
  >
    <template #reference>
      <el-badge :value="unread" :max="99" :hidden="unread === 0" class="bell-badge">
        <el-icon class="top-link bell" @click="toggle"><Bell /></el-icon>
      </el-badge>
    </template>

    <div class="notify-head">
      <span class="notify-title">通知中心</span>
      <div class="notify-head-ops">
        <el-radio-group v-model="filter" size="small" @change="reload">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="1">未读</el-radio-button>
        </el-radio-group>
        <el-button link type="primary" size="small" @click="readAll">全部已读</el-button>
      </div>
    </div>

    <div class="notify-list">
      <div v-if="loading" class="notify-state">加载中…</div>
      <div v-else-if="!items.length" class="notify-state">暂无通知</div>
      <div
        v-for="n in items"
        :key="n.id"
        class="notify-item"
        :class="{ unread: n.is_read === 0 }"
        @click="openItem(n)"
      >
        <span class="dot" />
        <div class="notify-body">
          <div class="notify-line">
            <span class="n-title">{{ n.title }}</span>
            <span class="n-time">{{ relative(n.created_at) }}</span>
          </div>
          <div class="n-content">{{ n.content }}</div>
        </div>
      </div>
    </div>

    <div v-if="items.length < total" class="notify-more">
      <el-button link type="primary" size="small" :loading="loading" @click="loadMore">加载更多</el-button>
    </div>
  </el-popover>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { scrmApi, fmtTime, type NotificationItem } from '../api/scrm'

const router = useRouter()
const visible = ref(false)
const unread = ref(0)
const items = ref<NotificationItem[]>([])
const total = ref(0)
const page = ref(1)
const filter = ref('')
const loading = ref(false)
let poller: number | undefined

async function refreshUnread() {
  try { unread.value = await scrmApi.unreadCount() } catch { /* ignore */ }
}

async function reload() {
  page.value = 1
  await load()
}

async function load() {
  loading.value = true
  try {
    const res = await scrmApi.notifications({ page: page.value, page_size: 20, unread: filter.value })
    if (page.value === 1) items.value = res.items
    else items.value = items.value.concat(res.items)
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function loadMore() {
  page.value += 1
  load()
}

async function toggle() {
  visible.value = !visible.value
  if (visible.value) await reload()
}

async function openItem(n: NotificationItem) {
  if (n.is_read === 0) {
    try {
      await scrmApi.markNotificationsRead([n.id])
      n.is_read = 1
      refreshUnread()
    } catch { /* ignore */ }
  }
  visible.value = false
  if (n.route) {
    const [path, qs] = n.route.split('?')
    const query: Record<string, string> = {}
    if (qs) {
      new URLSearchParams(qs).forEach((v, k) => { query[k] = v })
    }
    router.push({ path, query })
  }
}

async function readAll() {
  try {
    await scrmApi.markNotificationsRead(null, true)
    ElMessage.success('已全部标记为已读')
    await Promise.all([refreshUnread(), reload()])
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

function relative(s: string) {
  const t = new Date(s).getTime()
  const diff = Date.now() - t
  if (diff < 3600_000) return Math.max(1, Math.floor(diff / 60000)) + ' 分钟前'
  if (diff < 86400_000) return Math.floor(diff / 3600000) + ' 小时前'
  return fmtTime(s).slice(5, 16)
}

onMounted(() => {
  refreshUnread()
  poller = window.setInterval(refreshUnread, 60000)
})
onBeforeUnmount(() => window.clearInterval(poller))
</script>

<style scoped>
.bell-badge { display: inline-flex; align-items: center; }
.bell { font-size: 18px; }
.notify-head {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 8px;
}
.notify-title { font-weight: 700; font-size: 14px; }
.notify-head-ops { display: flex; align-items: center; gap: 8px; }
.notify-list { max-height: 380px; overflow-y: auto; margin: 0 -8px; }
.notify-item {
  display: flex; gap: 10px; padding: 10px 12px; cursor: pointer;
  border-radius: 6px;
}
.notify-item:hover { background: #f5f7fa; }
.notify-item .dot {
  width: 8px; height: 8px; border-radius: 50%; background: transparent;
  margin-top: 6px; flex-shrink: 0;
}
.notify-item.unread .dot { background: #f56c6c; }
.notify-body { flex: 1; min-width: 0; }
.notify-line { display: flex; justify-content: space-between; align-items: center; gap: 8px; }
.n-title {
  font-size: 13px; color: #1f2d3d;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.notify-item.unread .n-title { font-weight: 700; }
.n-time { font-size: 11px; color: #c0c4cc; flex-shrink: 0; }
.n-content {
  font-size: 12px; color: #606f82; line-height: 1.6; margin-top: 3px;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}
.notify-state { text-align: center; color: #909399; font-size: 13px; padding: 30px 0; }
.notify-more { text-align: center; margin-top: 6px; border-top: 1px solid #f0f0f0; padding-top: 6px; }
</style>
