<template>
  <el-dialog
    v-model="callState.visible"
    title="云座席 · 双呼工作条"
    width="380px"
    :close-on-click-modal="false"
    @closed="reset"
  >
    <el-tabs v-model="tab">
      <el-tab-pane label="拨号" name="dial">
        <!-- 顾问侧 -->
        <div class="party advisor">
          <el-icon><User /></el-icon>
          <div class="party-info">
            <div class="party-name">{{ me.name }}（顾问）</div>
            <div class="party-phone">{{ me.phone ? maskPhone(me.phone) : '—' }}</div>
          </div>
          <el-tag size="small" type="success" effect="plain">本机</el-tag>
        </div>

        <div class="call-arrow">
          <el-icon><Top /></el-icon>
          <span class="arrow-text">{{ arrowText }}</span>
          <el-icon><Bottom /></el-icon>
        </div>

        <!-- 客户侧 -->
        <div class="party customer" :class="{ ringing: phase === 'customer' }">
          <el-icon><Iphone /></el-icon>
          <div v-if="callState.lead" class="party-info">
            <div class="party-name">{{ callState.lead.name }}</div>
            <div class="party-phone">{{ maskPhone(callState.lead.phone) }}</div>
          </div>
          <div v-else class="party-info">
            <el-input
              v-model="pickKw"
              placeholder="输入姓名/手机/线索号选择客户"
              size="small"
              class="pick-input"
            >
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <div v-if="picks.length" class="pick-list">
              <div
                v-for="l in picks"
                :key="l.id"
                class="pick-row"
                @click="chooseLead(l)"
              >
                {{ l.name }} · {{ maskPhone(l.phone) }} · {{ l.code }}
              </div>
            </div>
          </div>
          <el-tag v-if="callState.lead" size="small" :type="resultTag" effect="dark">{{ phaseLabel }}</el-tag>
        </div>

        <div class="call-actions">
          <template v-if="phase === 'idle'">
            <el-button
              type="success"
              class="call-btn"
              :disabled="!callState.lead"
              @click="startDoubleCall"
            >
              <el-icon><Phone /></el-icon>&nbsp;双呼上号
            </el-button>
            <el-button v-if="callState.lead" link @click="callState.lead = null">更换客户</el-button>
          </template>
          <template v-else-if="phase === 'advisor'">
            <el-button type="success" loading plain>正在呼起顾问手机…</el-button>
          </template>
          <template v-else-if="phase === 'customer'">
            <el-button type="success" @click="customerAnswer"><el-icon><Check /></el-icon>&nbsp;客户已接听</el-button>
            <el-button type="danger" plain @click="customerMiss"><el-icon><Close /></el-icon>&nbsp;客户未接听</el-button>
          </template>
          <template v-else-if="phase === 'connected'">
            <div class="timer">{{ mmss(elapsed) }}</div>
            <el-button type="danger" @click="hangup"><el-icon><PhoneFilled /></el-icon>&nbsp;挂断</el-button>
          </template>
          <template v-else>
            <el-result
              :icon="connected ? 'success' : 'warning'"
              :title="connected ? '通话结束' : '客户未接通'"
              :sub-title="connected ? `本次通话 ${mmss(elapsed)}，记录已写入线索跟进日志` : '拨号事件已记录到线索跟进日志'"
              style="padding: 8px 0"
            />
            <el-button type="primary" plain @click="reset">继续外呼</el-button>
          </template>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="`通话记录${calls.length ? '（' + calls.length + '）' : ''}`" name="logs">
        <div v-if="!calls.length" class="logs-empty">暂无通话记录</div>
        <el-timeline v-else class="logs-timeline">
          <el-timeline-item
            v-for="c in calls"
            :key="c.id"
            :type="c.result === '已接通' ? 'success' : 'info'"
            :timestamp="fmtTime(c.created_at)"
          >
            <div class="log-line">
              <span class="log-name">{{ c.lead_name }}</span>
              <el-tag size="small" :type="c.result === '已接通' ? 'success' : 'info'" effect="plain">
                {{ c.result || '拨号中' }}
              </el-tag>
              <span v-if="c.duration_sec" class="log-dur">{{ mmss(c.duration_sec) }}</span>
            </div>
            <div class="log-content">{{ c.content }}</div>
          </el-timeline-item>
        </el-timeline>
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import {
  User, Iphone, Phone, PhoneFilled, Search, Check, Close, Top, Bottom
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { scrmApi, maskPhone, fmtTime, type MeInfo, type Lead, type FollowLog } from '../api/scrm'
import { callState } from '../store/call'

const me = ref<MeInfo>({ name: '陈元素', title: '', city: '', store: '', unread: 0, roles: [] })
const tab = ref('dial')
const phase = ref<'idle' | 'advisor' | 'customer' | 'connected' | 'done'>('idle')
const connected = ref(false)
const elapsed = ref(0)
const logId = ref(0)
let timerId: number | undefined
let pickTimer: number | undefined

const pickKw = ref('')
const picks = ref<Lead[]>([])
const calls = ref<FollowLog[]>([])

const arrowText = computed(() => ({
  idle: '准备双呼', advisor: '呼起顾问中', customer: '呼起客户中', connected: '通话中', done: '通话结束'
} as Record<string, string>)[phase.value])

const phaseLabel = computed(() => ({
  idle: '待呼叫', advisor: '呼叫中', customer: '振铃中', connected: '通话中', done: connected.value ? '已接通' : '未接通'
} as Record<string, string>)[phase.value])

const resultTag = computed(() => {
  if (phase.value === 'connected') return 'danger'
  if (phase.value === 'done') return connected.value ? 'success' : 'info'
  return 'warning'
})

watch(pickKw, () => {
  window.clearTimeout(pickTimer)
  pickTimer = window.setTimeout(async () => {
    const q = pickKw.value.trim()
    picks.value = q ? (await scrmApi.search(q)).leads.slice(0, 6) : []
  }, 250)
})

watch(() => callState.visible, async v => {
  if (v) {
    if (!me.value.phone) {
      try { me.value = await scrmApi.me() } catch { /* 兜底 */ }
    }
    loadCalls()
  }
})

function chooseLead(l: Lead) {
  callState.lead = l
  picks.value = []
  pickKw.value = ''
}

async function startDoubleCall() {
  if (!callState.lead) return
  try {
    // 点击「双呼上号」即把拨号事件落入线索跟进日志
    const log = await scrmApi.callDial(callState.lead.id, callState.lead.phone)
    logId.value = log.id
  } catch (e: any) {
    ElMessage.error(e.message || '外呼失败')
    return
  }
  phase.value = 'advisor'
  window.setTimeout(() => {
    if (phase.value === 'advisor') phase.value = 'customer'
  }, 2000)
}

function customerAnswer() {
  phase.value = 'connected'
  connected.value = true
  elapsed.value = 0
  timerId = window.setInterval(() => { elapsed.value += 1 }, 1000)
}

function customerMiss() {
  connected.value = false
  finish()
}

function hangup() {
  if (timerId) window.clearInterval(timerId)
  finish()
}

async function finish() {
  if (logId.value) {
    try {
      await scrmApi.callComplete(logId.value, connected.value, elapsed.value)
    } catch (e: any) {
      ElMessage.error(e.message || '通话记录保存失败')
    }
  }
  phase.value = 'done'
  loadCalls()
}

async function loadCalls() {
  try { calls.value = await scrmApi.recentCalls() } catch { /* ignore */ }
}

function mmss(sec: number) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

function reset() {
  phase.value = 'idle'
  connected.value = false
  elapsed.value = 0
  logId.value = 0
  if (timerId) window.clearInterval(timerId)
  tab.value = 'dial'
  callState.lead = null
}

onBeforeUnmount(() => { if (timerId) window.clearInterval(timerId) })
</script>

<style scoped>
.party {
  display: flex; align-items: center; gap: 10px;
  border: 1px solid #e4e7ed; border-radius: 8px; padding: 10px 12px;
}
.party.advisor { background: #f0fbf7; border-color: #b7e8d2; }
.party.customer { background: #fafafa; }
.party.customer.ringing { border-color: #f56c6c; animation: blink 1s infinite; }
@keyframes blink { 50% { box-shadow: 0 0 0 3px rgba(245,108,108,.15); } }
.party-info { flex: 1; min-width: 0; }
.party-name { font-size: 14px; font-weight: 600; color: #1f2d3d; }
.party-phone { font-size: 13px; color: #606f82; margin-top: 2px; }
.call-arrow {
  display: flex; align-items: center; justify-content: center; gap: 8px;
  padding: 10px 0; color: #00a86b; font-size: 12px;
}
.arrow-text { color: #909399; }
.call-actions { text-align: center; padding: 12px 0 4px; min-height: 64px; }
.call-btn { width: 100%; }
.timer {
  font-size: 30px; font-weight: 700; color: #f56c6c;
  font-variant-numeric: tabular-nums; margin-bottom: 8px;
}
.pick-input { width: 100%; }
.pick-list {
  margin-top: 6px; max-height: 150px; overflow-y: auto;
  border: 1px solid #ebeef5; border-radius: 6px;
}
.pick-row { padding: 6px 10px; font-size: 12px; cursor: pointer; }
.pick-row:hover { background: #f0fbf7; }
.logs-empty { text-align: center; color: #909399; font-size: 13px; padding: 30px 0; }
.log-line { display: flex; align-items: center; gap: 8px; }
.log-name { font-weight: 600; font-size: 13px; }
.log-dur { font-size: 12px; color: #00a86b; font-variant-numeric: tabular-nums; }
.log-content { font-size: 12px; color: #909399; margin-top: 2px; }
:deep(.el-result__icon) { display: none; }
</style>
