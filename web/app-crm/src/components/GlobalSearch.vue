<template>
  <div class="gs-wrap" ref="wrapRef">
    <el-input
      v-model="kw"
      class="gs-input"
      placeholder="搜索：姓名 / 手机（支持密文）/ 线索号 / 学员号"
      :prefix-icon="Search"
      clearable
      @focus="showPanel = true"
      @keyup.enter="enterJump"
    />
    <div v-if="showPanel" class="gs-panel">
      <div v-if="!kw.trim()" class="gs-tip">
        支持明文/密文电话（如 138****0137）、姓名、线索号、学员号，商机号将在第三期支持
      </div>
      <div v-else-if="loading" class="gs-tip">搜索中…</div>
      <template v-else-if="result">
        <template v-if="!result.leads.length && !result.students.length">
          <div class="gs-empty">未找到与「{{ kw }}」相关的线索或学员</div>
        </template>
        <template v-else>
          <div v-if="result.leads.length" class="gs-group">
            <div class="gs-group-title">线索 <span>{{ result.leads.length }}</span></div>
            <div
              v-for="l in result.leads"
              :key="'l' + l.id"
              class="gs-row"
              @click="goLead(l.id)"
            >
              <div class="gs-main">
                <span class="gs-name">{{ l.name }}</span>
                <el-tag size="small" type="success" effect="plain">{{ l.status }}</el-tag>
                <span v-if="l.is_duplicate" class="gs-dup">重复</span>
              </div>
              <div class="gs-sub">{{ maskPhone(l.phone) }} · {{ l.code }} · {{ l.owner_name || '未分配' }}</div>
            </div>
          </div>
          <div v-if="result.students.length" class="gs-group">
            <div class="gs-group-title">学员 <span>{{ result.students.length }}</span></div>
            <div
              v-for="s in result.students"
              :key="'s' + s.student_no"
              class="gs-row"
              @click="goLead(s.lead_id)"
            >
              <div class="gs-main">
                <span class="gs-name">{{ s.name }}</span>
                <el-tag size="small" type="info" effect="plain">学员</el-tag>
              </div>
              <div class="gs-sub">学员号 {{ s.student_no }} · {{ maskPhone(s.phone) }} · {{ s.school }} / {{ s.grade }}</div>
            </div>
          </div>
          <div class="gs-group">
            <div class="gs-group-title gs-future">商机（第三期上线）</div>
          </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { scrmApi, maskPhone, type SearchResult, type Lead } from '../api/scrm'

const router = useRouter()
const kw = ref('')
const showPanel = ref(false)
const loading = ref(false)
const result = ref<SearchResult | null>(null)
const wrapRef = ref<HTMLElement | null>(null)
let timer: number | undefined

async function doSearch() {
  const q = kw.value.trim()
  if (!q) {
    result.value = null
    return
  }
  loading.value = true
  try {
    result.value = await scrmApi.search(q)
  } catch {
    result.value = null
  } finally {
    loading.value = false
  }
}

watch(kw, () => {
  window.clearTimeout(timer)
  timer = window.setTimeout(doSearch, 250)
})

function goLead(id: number) {
  showPanel.value = false
  router.push({ path: '/all-leads', query: { lead_id: String(id) } })
}

function enterJump() {
  const first = result.value?.leads[0] || (result.value?.students[0] ? { id: result.value.students[0].lead_id } as Lead : null)
  if (first) goLead(first.id)
}

function onDocClick(e: MouseEvent) {
  if (wrapRef.value && !wrapRef.value.contains(e.target as Node)) {
    showPanel.value = false
  }
}

onMounted(() => document.addEventListener('click', onDocClick))
onBeforeUnmount(() => document.removeEventListener('click', onDocClick))
</script>

<style scoped>
.gs-wrap { position: relative; width: 320px; }
.gs-panel {
  position: absolute; top: 40px; right: 0; width: 420px;
  background: #fff; border: 1px solid #e4e7ed; border-radius: 8px;
  box-shadow: 0 6px 24px rgba(0,0,0,.12); z-index: 2200;
  max-height: 460px; overflow-y: auto; padding: 6px 0;
}
.gs-tip, .gs-empty { padding: 16px 18px; font-size: 12px; color: #909399; line-height: 1.7; }
.gs-group-title {
  padding: 6px 16px 4px; font-size: 12px; color: #909399; font-weight: 600;
  display: flex; align-items: center; gap: 6px;
}
.gs-group-title span { color: #c0c4cc; font-weight: 400; }
.gs-future { color: #c0c4cc; }
.gs-row { padding: 8px 16px; cursor: pointer; }
.gs-row:hover { background: #f0fbf7; }
.gs-main { display: flex; align-items: center; gap: 8px; }
.gs-name { font-size: 13px; color: #1f2d3d; font-weight: 600; }
.gs-sub { font-size: 12px; color: #909399; margin-top: 2px; }
.gs-dup {
  font-size: 11px; color: #e6a23c; border: 1px solid #f5dab1;
  border-radius: 3px; padding: 0 4px;
}
</style>
