<template>
  <div class="page">
    <el-card shadow="never">
      <!-- 筛选区：关联学员 vs 活动数据 双视角 -->
      <el-form :model="f" inline class="filter-form">
        <el-form-item label="关联学员名"><el-input v-model="f.related_name" clearable style="width:130px" @keyup.enter="search" /></el-form-item>
        <el-form-item label="关联手机号"><el-input v-model="f.related_phone" clearable style="width:140px" @keyup.enter="search" /></el-form-item>
        <el-form-item label="活动学员名"><el-input v-model="f.activity_name" clearable style="width:130px" @keyup.enter="search" /></el-form-item>
        <el-form-item label="活动手机号"><el-input v-model="f.activity_phone" clearable style="width:140px" @keyup.enter="search" /></el-form-item>
        <el-form-item label="关联学员号"><el-input v-model="f.student_no" clearable style="width:140px" @keyup.enter="search" /></el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 操作区 -->
      <div class="op-bar">
        <div class="op-left">
          <el-button @click="adjustVisible = true" :disabled="!ids.length">批量调整</el-button>
          <el-button @click="exportCsv"><el-icon><Download /></el-icon>导出</el-button>
          <el-button type="warning" plain @click="ignoreAll" :disabled="!ids.length">忽略异常数据</el-button>
        </div>
        <el-button text @click="load"><el-icon><Refresh /></el-icon>刷新</el-button>
      </div>

      <el-alert type="warning" :closable="false" show-icon style="margin-bottom:10px">
        异常线索为导入/生成过程中发生数据冲突或信息不完整的线索，表格并列展示「库内关联学员」与「活动数据」两侧信息，请人工核对后调整或忽略。
      </el-alert>

      <el-table :data="rows" border size="small" height="56vh"
        @selection-change="(s: any[]) => ids = s.map(x => x.id)">
        <el-table-column type="selection" width="44" fixed />
        <el-table-column label="关联学员名称" width="120" fixed>
          <template #default="{ row }">
            <el-link type="primary" :underline="false">{{ row.related_name || '—' }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="related_phone" label="关联手机号" width="130" />
        <el-table-column prop="activity_student_name" label="活动学员名称" width="120">
          <template #default="{ row }">
            <span :class="{ conflict: row.activity_student_name && row.related_name && row.activity_student_name !== row.related_name }">
              {{ row.activity_student_name || '—' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="活动手机号" width="130">
          <template #default="{ row }">
            <span :class="{ conflict: row.activity_phone && row.related_phone && row.activity_phone !== row.related_phone }">
              {{ row.activity_phone || '—' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="related_student_no" label="关联学员号" width="140" />
        <el-table-column prop="channel_l1" label="一级渠道" width="100" />
        <el-table-column prop="channel_l2" label="二级渠道" width="100" />
        <el-table-column prop="channel_l3" label="三级渠道" width="100" />
        <el-table-column prop="system_source_l1" label="一级系统来源" width="110" />
        <el-table-column prop="system_source_l2" label="二级系统来源" width="110" />
        <el-table-column prop="activity_name" label="关联活动" min-width="180" show-overflow-tooltip />
        <el-table-column prop="data_remark" label="数据备注" min-width="200" show-overflow-tooltip />
        <el-table-column prop="intention_dept" label="意向部门" width="100" />
        <el-table-column prop="intention_project" label="意向项目" width="100" />
        <el-table-column label="状态" width="90" fixed="right">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status === '已忽略' ? 'info' : 'danger'" effect="plain">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pager"
        layout="total, prev, pager, next, sizes"
        :total="total" :page-size="pageSize" :current-page="page" :page-sizes="[20, 50, 100]"
        @current-change="(p: number) => { page = p; load() }"
        @size-change="(s: number) => { pageSize = s; page = 1; load() }"
      />
    </el-card>

    <!-- 批量调整 -->
    <el-dialog v-model="adjustVisible" title="批量调整异常数据" width="500px">
      <el-alert type="info" :closable="false" show-icon style="margin-bottom:14px">
        修正字段后数据将标记为已核对，可重新进入线索分配流程。
      </el-alert>
      <el-form label-width="110px">
        <el-form-item label="关联学员名"><el-input v-model="form.related_name" placeholder="不修改留空" /></el-form-item>
        <el-form-item label="关联手机号"><el-input v-model="form.related_phone" placeholder="不修改留空" /></el-form-item>
        <el-form-item label="关联学员号"><el-input v-model="form.related_student_no" placeholder="不修改留空" /></el-form-item>
        <el-form-item label="意向部门">
          <el-select v-model="form.intention_dept" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in opts.intention_dept || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="意向项目">
          <el-select v-model="form.intention_project" clearable placeholder="不修改" style="width:100%">
            <el-option v-for="o in opts.intention_project || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="数据备注"><el-input v-model="form.data_remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdjust">确认调整（{{ ids.length }} 条）</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Refresh } from '@element-plus/icons-vue'
import { scrmApi, type ExceptionLead } from '../api/scrm'
import { downloadCsv } from '../utils/export'

const opts = ref<Record<string, any[]>>({})
const rows = ref<ExceptionLead[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const ids = ref<number[]>([])

const f = reactive({ related_name: '', related_phone: '', activity_name: '', activity_phone: '', student_no: '' })
const form = reactive<any>({ related_name: '', related_phone: '', related_student_no: '', intention_dept: '', intention_project: '', data_remark: '' })
const adjustVisible = ref(false)

async function load() {
  const res = await scrmApi.exceptions({ page: page.value, page_size: pageSize.value, ...f })
  rows.value = res.items
  total.value = res.total
}
function search() { page.value = 1; load() }
function reset() {
  Object.keys(f).forEach(k => (f as any)[k] = '')
  search()
}

function exportCsv() {
  downloadCsv('异常线索.csv',
    ['关联学员名称', '关联手机号', '活动学员名称', '活动手机号', '关联学员号', '渠道', '关联活动', '数据备注', '状态'],
    rows.value.map(r => [r.related_name, r.related_phone, r.activity_student_name, r.activity_phone,
      r.related_student_no, [r.channel_l1, r.channel_l2, r.channel_l3].filter(Boolean).join('/'),
      r.activity_name, r.data_remark, r.status]))
  ElMessage.success(`已导出当前页 ${rows.value.length} 条`)
}

async function ignoreAll() {
  await ElMessageBox.confirm(`确认忽略选中的 ${ids.value.length} 条异常数据？忽略后可通过筛选找回。`, '忽略确认', { type: 'warning' })
  await scrmApi.ignoreExceptions(ids.value)
  ElMessage.success('已标记为已忽略')
  load()
}

async function submitAdjust() {
  const fields: Record<string, string> = {}
  Object.entries(form).forEach(([k, v]) => { if (v) fields[k] = String(v) })
  if (Object.keys(fields).length === 0) return ElMessage.warning('请至少填写一个要修正的字段')
  await scrmApi.adjustExceptions(ids.value, fields)
  ElMessage.success('调整完成')
  adjustVisible.value = false
  Object.keys(form).forEach(k => (form as any)[k] = '')
  load()
}

onMounted(async () => {
  opts.value.intention_dept = await scrmApi.dict('intention_dept')
  opts.value.intention_project = await scrmApi.dict('intention_project')
  load()
})
</script>

<style scoped>
.filter-form :deep(.el-form-item) { margin-bottom: 8px; }
.op-bar { display: flex; justify-content: space-between; align-items: center; margin: 4px 0 10px; }
.pager { margin-top: 10px; justify-content: flex-end; }
.conflict { color: #f56c6c; font-weight: 600; text-decoration: underline wavy #f56c6c; }
</style>
