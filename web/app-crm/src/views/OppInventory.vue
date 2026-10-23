<template>
  <div class="inv-page">
    <div class="inv-head">
      <el-radio-group v-model="dim" size="small" @change="load">
        <el-radio-button value="school">学校维度</el-radio-button>
        <el-radio-button value="dept">部门维度</el-radio-button>
        <el-radio-button value="consultant">顾问维度</el-radio-button>
      </el-radio-group>
      <el-radio-group v-model="viewMode" size="small">
        <el-radio-button value="both">表格和图表</el-radio-button>
        <el-radio-button value="table">仅表格</el-radio-button>
        <el-radio-button value="chart">仅图表</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 昨日统计 6 卡 -->
    <div v-if="viewMode !== 'table'" class="card-row">
      <div v-for="c in data.yesterday" :key="c.key" class="stat-card">
        <div class="sv">{{ c.value }}<span class="su">{{ c.unit }}</span></div>
        <div class="st">{{ c.title }}</div>
      </div>
    </div>

    <!-- 商机圈选 4 卡 -->
    <div v-if="viewMode !== 'table'" class="card-row">
      <div v-for="c in data.circle" :key="c.key" class="stat-card circle">
        <div class="sv">{{ c.key }}<span class="su"> {{ c.value }}{{ c.unit }}</span></div>
        <div class="st">{{ c.title }}</div>
      </div>
    </div>

    <!-- 全局指标库 -->
    <div v-if="viewMode !== 'table'" class="metric-global">
      <span class="mg-title">指标库命中率：</span>
      <div v-for="m in data.global_metric" :key="m.key" class="mg-item">
        <span>{{ m.title }} 平均分 {{ m.avg_score }} · 命中率 {{ m.hit_rate }}</span>
        <div class="bar"><div class="bar-in" :style="{ width: m.hit_rate }" /></div>
      </div>
    </div>

    <!-- 表格 -->
    <el-table v-if="viewMode !== 'chart'" :data="data.rows" border size="small" class="inv-table">
      <el-table-column prop="name" :label="dimLabel" fixed min-width="150" />
      <el-table-column prop="total" label="商机量" width="80" />
      <el-table-column label="S" width="60"><template #default="{ row }">{{ row.s }}</template></el-table-column>
      <el-table-column label="A" width="60"><template #default="{ row }">{{ row.a }}</template></el-table-column>
      <el-table-column label="B" width="60"><template #default="{ row }">{{ row.b }}</template></el-table-column>
      <el-table-column label="C" width="60"><template #default="{ row }">{{ row.c }}</template></el-table-column>
      <el-table-column prop="deals" label="成单量" width="80" />
      <el-table-column prop="conv" label="转化率" width="90" />
      <el-table-column v-for="m in metricMeta" :key="m.key" :label="m.title + '命中率'" width="120">
        <template #default="{ row }">{{ (row.metrics.find((x: any) => x.key === m.key) || {}).hit_rate || '0%' }}</template>
      </el-table-column>
    </el-table>

    <!-- 图表分组卡片 -->
    <template v-if="viewMode !== 'table'">
      <div v-for="row in data.rows" :key="row.name" class="group-card">
        <div class="gc-title">{{ row.name }}
          <span class="gc-sub">商机 {{ row.total }} · 成单 {{ row.deals }} · 转化率 {{ row.conv }}</span>
        </div>
        <div class="gc-grid">
          <!-- 转化率趋势（纯 CSS 折线） -->
          <div class="gc-block">
            <div class="blk-title">S/A/B/C 商机量·转化率（近6月）</div>
            <div class="months">
              <div v-for="(mo, i) in row.months" :key="i" class="month-col">
                <div class="stack">
                  <div class="seg seg-s" :style="{ height: mo.s * 14 + 'px' }" :title="`S ${mo.s}`" />
                  <div class="seg seg-a" :style="{ height: mo.a * 14 + 'px' }" :title="`A ${mo.a}`" />
                  <div class="seg seg-b" :style="{ height: mo.b * 14 + 'px' }" :title="`B ${mo.b}`" />
                  <div class="seg seg-c" :style="{ height: mo.c * 14 + 'px' }" :title="`C ${mo.c}`" />
                </div>
                <div class="mo-conv">{{ mo.conv }}</div>
                <div class="mo-lbl">{{ mo.month.slice(5) }}</div>
              </div>
            </div>
          </div>
          <!-- 等级占比 100% 堆叠条 -->
          <div class="gc-block">
            <div class="blk-title">等级占比</div>
            <div class="stacked-bar">
              <div class="seg-s" :style="{ width: pct(row.s, row.total) }">S {{ row.s }}</div>
              <div class="seg-a" :style="{ width: pct(row.a, row.total) }">A {{ row.a }}</div>
              <div class="seg-b" :style="{ width: pct(row.b, row.total) }">B {{ row.b }}</div>
              <div class="seg-c" :style="{ width: pct(row.c, row.total) }">C {{ row.c }}</div>
            </div>
            <!-- 指标分类命中率横向条形 -->
            <div class="blk-title" style="margin-top:14px">指标分类命中率</div>
            <div v-for="m in row.metrics" :key="m.key" class="hmetric">
              <span class="hm-lbl">{{ m.title }}</span>
              <div class="bar"><div class="bar-in" :style="{ width: m.hit_rate }" /></div>
              <span class="hm-val">{{ m.hit_rate }} / {{ m.avg_score }}分</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { scrmApi } from '../api/scrm'

const dim = ref('school')
const viewMode = ref('both')
const data = ref<any>({ rows: [], yesterday: [], circle: [], global_metric: [] })
const metricMeta = [
  { key: 'client', title: '客户信息' }, { key: 'activity', title: '活动行为' },
  { key: 'communication', title: '沟通质量' }, { key: 'payment', title: '付费行为' }
]
const dimLabel = computed(() => ({ school: '学校', dept: '部门', consultant: '顾问' })[dim.value])

function pct(n: number, d: number) { return d ? (n / d * 100) + '%' : '0%' }

async function load() {
  data.value = await scrmApi.inventory(dim.value)
}
onMounted(load)
</script>

<style scoped>
.inv-page { padding: 12px 16px; }
.inv-head { display: flex; justify-content: space-between; margin-bottom: 10px; }
.card-row { display: flex; gap: 10px; margin-bottom: 10px; flex-wrap: wrap; }
.stat-card { flex: 1; min-width: 130px; background: #fff; border: 1px solid #ebeef5; border-radius: 8px; padding: 12px; text-align: center; }
.stat-card.circle { background: linear-gradient(135deg, #f0f7ff, #fff); }
.sv { font-size: 22px; font-weight: 700; color: #303133; }
.su { font-size: 12px; font-weight: 400; color: #909399; margin-left: 2px; }
.st { color: #909399; font-size: 12px; margin-top: 4px; }
.metric-global { display: flex; gap: 18px; align-items: center; background: #fff; border: 1px solid #ebeef5; border-radius: 8px; padding: 10px 14px; margin-bottom: 12px; flex-wrap: wrap; }
.mg-title { font-weight: 600; }
.mg-item { display: flex; align-items: center; gap: 8px; font-size: 12px; color: #606266; }
.bar { width: 90px; height: 8px; background: #f0f2f5; border-radius: 4px; overflow: hidden; }
.bar-in { height: 100%; background: linear-gradient(90deg, #409eff, #67c23a); }
.inv-table { margin-bottom: 14px; }
.group-card { background: #fff; border: 1px solid #ebeef5; border-radius: 8px; padding: 14px; margin-bottom: 12px; }
.gc-title { font-weight: 600; margin-bottom: 10px; }
.gc-sub { font-weight: 400; color: #909399; font-size: 12px; margin-left: 10px; }
.gc-grid { display: grid; grid-template-columns: 1.3fr 1fr; gap: 18px; }
.blk-title { font-size: 13px; color: #606266; margin-bottom: 8px; }
.months { display: flex; gap: 10px; align-items: flex-end; height: 170px; padding-top: 8px; }
.month-col { display: flex; flex-direction: column; align-items: center; flex: 1; }
.stack { display: flex; flex-direction: column-reverse; justify-content: flex-start; height: 120px; width: 100%; max-width: 38px; }
.seg { width: 100%; }
.seg-s { background: #f56c6c; }.seg-a { background: #e6a23c; }.seg-b { background: #67c23a; }.seg-c { background: #a8abb2; }
.mo-conv { font-size: 10px; color: #409eff; margin-top: 4px; }
.mo-lbl { font-size: 11px; color: #909399; }
.stacked-bar { display: flex; width: 100%; height: 26px; border-radius: 4px; overflow: hidden; color: #fff; font-size: 11px; line-height: 26px; text-align: center; }
.hmetric { display: flex; align-items: center; gap: 8px; margin-bottom: 6px; font-size: 12px; }
.hm-lbl { width: 64px; color: #606266; }.hm-val { width: 104px; color: #909399; }
</style>
