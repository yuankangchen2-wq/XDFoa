<template>
  <div class="page">
    <div class="head">
      <div class="title">商机备单</div>
      <el-button size="small" @click="resetFilter">重置｜筛选</el-button>
    </div>

    <div class="tags">
      <el-check-tag v-for="t in tagDefs" :key="t.v" :checked="activeTags.includes(t.v)"
        @change="toggleTag(t.v)">{{ t.t }}</el-check-tag>
    </div>

    <el-table :data="pagedRows" border size="small" v-loading="loading">
      <el-table-column prop="code" label="商机编码" width="150" />
      <el-table-column prop="name" label="姓名" width="90" />
      <el-table-column prop="phone" label="手机" width="120">
        <template #default="{ row }">{{ maskPhone(row.phone) }}</template>
      </el-table-column>
      <el-table-column prop="grade" label="年级" width="80" />
      <el-table-column prop="level" label="意向等级" width="90">
        <template #default="{ row }">
          <el-tag v-if="activeTags.includes('hightoday')" type="info" size="small" disabled>高（锁定）</el-tag>
          <el-tag v-else :type="row.level === '高' ? 'danger' : row.level === '中' ? 'warning' : 'info'" size="small">{{ row.level }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="grade2" label="商机等级" width="90">
        <template #default="{ row }">
          <el-tag :type="gradeTag(row.grade2)" size="small" effect="dark">{{ row.grade2 }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="owner_name" label="所属人" width="100" />
      <el-table-column prop="promise_visit_at" label="诺访时间" width="160">
        <template #default="{ row }">{{ fmtTime(row.promise_visit_at) }}</template>
      </el-table-column>
      <el-table-column prop="visit_converted" label="到访转化" width="90">
        <template #default="{ row }">
          <el-tag :type="row.visit_converted ? 'success' : 'info'" size="small">{{ row.visit_converted ? '已转化' : '未转化' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_follow_content" label="最近跟进" min-width="200" show-overflow-tooltip />
      <el-table-column prop="next_follow_at" label="下次回访" width="160">
        <template #default="{ row }">{{ fmtTime(row.next_follow_at) }}</template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" background layout="total,prev,pager,next,jumper" :total="rows.length"
      :page-size="10" v-model:current-page="page" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { scrmApi, maskPhone, fmtTime, type Opportunity } from '../api/scrm'

const tagDefs = [
  { v: 'promiseweek', t: '本周诺访' },
  { v: 'hightoday', t: '今日高意向' },
  { v: 'visitnotconv', t: '本周到访未转化' },
  { v: 'all', t: '全部备单' }
]
const activeTags = ref<string[]>(['promiseweek'])
const allRows = ref<Opportunity[]>([])
const rows = ref<Opportunity[]>([])
const page = ref(1)
const loading = ref(false)

const pagedRows = computed(() => rows.value.slice((page.value - 1) * 10, page.value * 10))

function gradeTag(g: string) { return { S: 'danger', A: 'warning', B: 'success', C: 'info' }[g] || 'info' }

function toggleTag(v: string) {
  const idx = activeTags.value.indexOf(v)
  if (v === 'all') {
    activeTags.value = activeTags.value.includes('all') ? [] : ['all']
  } else {
    const next = activeTags.value.filter(x => x !== 'all')
    if (idx >= 0) next.splice(next.indexOf(v), 1)
    else next.push(v)
    activeTags.value = next.length ? next : ['all']
  }
  load()
}

function resetFilter() {
  activeTags.value = ['promiseweek']
  page.value = 1
  load()
}

async function load() {
  page.value = 1
  loading.value = true
  try {
    const tags = activeTags.value.includes('all') ? ['pending'] : activeTags.value
    const results = await Promise.all(tags.map(q =>
      scrmApi.opportunities({ scope: 'dept', quick: q === 'pending' ? 'pending' : q, page: 1, page_size: 200 })
    ))
    const map = new Map<number, Opportunity>()
    results.forEach(r => r.items.forEach(o => map.set(o.id, o)))
    rows.value = [...map.values()]
    allRows.value = rows.value
  } finally { loading.value = false }
}

onMounted(load)
</script>

<style scoped>
.page { padding: 12px 16px; }
.head { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 16px; font-weight: 600; }
.tags { display: flex; gap: 10px; margin: 12px 0; }
.pager { margin-top: 10px; justify-content: flex-end; }
</style>
