<template>
  <div class="lead-import">
    <!-- 顶部时间 chip -->
    <div class="chip-bar">
      <el-tag closable type="info" effect="plain">创建时间：2026-08-14 00:00:00 ~ 2026-10-14 23:59:59</el-tag>
    </div>

    <!-- ① 基础信息 -->
    <el-card shadow="never" class="sec">
      <template #header><span class="sec-title">① 基础信息</span></template>
      <el-form :model="form" label-width="110px" label-position="right">
        <el-form-item label="导入方式" required>
          <el-radio-group v-model="form.import_method">
            <el-radio value="phone">手机号导入</el-radio>
            <el-radio value="student_no">学员号导入</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="数据来源" required>
          <el-select v-model="form.data_source" placeholder="请选择数据来源" style="width:260px">
            <el-option v-for="o in dict.data_source || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="市场渠道" required>
          <el-cascader
            v-model="channelValue"
            :options="channelOptions"
            :props="{ expandTrigger: 'hover' }"
            placeholder="请选择一级 / 二级 / 三级渠道"
            style="width:420px"
            @change="onChannelChange"
          />
        </el-form-item>
        <el-form-item>
          <el-alert type="info" :closable="false" show-icon style="max-width:680px">
            导入后系统按规则优先确定线索渠道（与本任务渠道可能不一样），渠道不同时以本次渠道为准。
            <el-link type="primary" :underline="false" @click="ruleDlg = true">点击了解填写规则</el-link>
          </el-alert>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- ② 活动关联 -->
    <el-card shadow="never" class="sec">
      <template #header><span class="sec-title">② 活动关联</span></template>
      <el-form :model="form" label-width="110px">
        <el-row :gutter="12">
          <el-col :span="10">
            <el-form-item label="活动名称">
              <el-select
                v-model="form.activity_id"
                filterable
                remote
                :remote-method="searchActivities"
                placeholder="输入关键词联想活动"
                style="width:100%"
                @change="onActivityChange"
              >
                <el-option
                  v-for="a in activities"
                  :key="a.id"
                  :label="a.name"
                  :value="a.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="活动类别" label-width="90px">
              <el-select v-model="form.activity_category" placeholder="活动类别" style="width:100%">
                <el-option v-for="o in dict.activity_category || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="活动类型" label-width="90px">
              <el-select v-model="form.activity_type" placeholder="活动类型" style="width:100%">
                <el-option v-for="o in dict.activity_type || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="10">
            <el-form-item label="活动意向部门">
              <el-select v-model="form.intention_dept" placeholder="组织架构 / 意向部门" style="width:100%">
                <el-option v-for="o in dict.intention_dept || []" :key="o.value" :label="o.label" :value="o.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-button type="primary" plain @click="actDlg = true">
              <el-icon><Connection /></el-icon>&nbsp;关联已有活动
            </el-button>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- ③ 去重规则 -->
    <el-card shadow="never" class="sec">
      <template #header>
        <div class="sec-head-row">
          <span class="sec-title">③ 去重规则</span>
          <el-switch v-model="form.dedup.enable" />
        </div>
      </template>
      <div v-show="form.dedup.enable" class="dedup-panel">
        <el-checkbox v-model="form.dedup.followed_checked">
          指定时间段内跟进过的线索去重
        </el-checkbox>
        <div v-show="form.dedup.followed_checked" class="dedup-sub">
          最近
          <el-input-number v-model="form.dedup.followed_dur_days" :min="1" :max="365" size="small" />
          天内跟进过的线索不导入
        </div>

        <el-checkbox v-model="form.dedup.old_student_checked">老学生学去重</el-checkbox>
        <div v-show="form.dedup.old_student_checked" class="dedup-sub">
          <el-select v-model="form.dedup.old_student.status" placeholder="在读状态" size="small" style="width:120px">
            <el-option v-for="o in dict.in_school_status || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-select v-model="form.dedup.old_student.fresh_type" placeholder="新生类型" size="small" style="width:120px">
            <el-option v-for="o in dict.fresh_type || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-select v-model="form.dedup.old_student.project" placeholder="在读项目" size="small" style="width:120px">
            <el-option v-for="o in dict.intention_project || []" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-select v-model="form.dedup.old_student.in_year" placeholder="当年在校" size="small" style="width:120px">
            <el-option label="是" value="是" /><el-option label="否" value="否" />
          </el-select>
          <el-checkbox v-model="form.dedup.old_student.loss_enabled">流失天数</el-checkbox>
          <el-input-number v-model="form.dedup.old_student.loss_days" :min="1" :max="3650" size="small" :disabled="!form.dedup.old_student.loss_enabled" />
          天
        </div>

        <el-checkbox disabled>手机号保护</el-checkbox>
        <div class="dedup-tip">新数据与指定场次内做过的手机号数据不会导入（系统强制开启）</div>

        <el-checkbox v-model="form.dedup.opportunity_dedup">活动商机去重</el-checkbox>
        <div class="dedup-tip">默认只去重意向的商机</div>

        <el-checkbox v-model="form.dedup.task_dedup">指定任务去重</el-checkbox>
        <div class="dedup-tip">与已发起的重任任务下的线索去重</div>
      </div>
    </el-card>

    <!-- ④ 导入为商机 -->
    <el-card shadow="never" class="sec">
      <template #header>
        <div class="sec-head-row">
          <span class="sec-title">④ 导入为商机</span>
          <el-switch v-model="form.as_opportunity" />
        </div>
      </template>
      <div class="dedup-tip">开启后，本次导入数据在生成线索的同时生成商机记录。</div>
    </el-card>

    <!-- ⑤ 数据提交 -->
    <el-card shadow="never" class="sec">
      <template #header><span class="sec-title">⑤ 数据提交</span></template>
      <el-radio-group v-model="form.submit_type" class="submit-radio">
        <el-radio value="supervisor">交接给主管</el-radio>
        <el-radio value="employee">交接给员工</el-radio>
        <el-radio value="none">暂不分配</el-radio>
      </el-radio-group>
      <el-alert type="error" :closable="false" show-icon class="warn-box">
        如果上传的 Excel 表格中「归属人」有归属员工，则不论选择交接给谁，都会按照上传的归属人自动分配，直接分配到所属员工。如果选择暂不分配，则数据不会按照 Excel 表格中的归属人进行分配。
      </el-alert>
    </el-card>

    <!-- ⑥ 任务与文件 -->
    <el-card shadow="never" class="sec">
      <template #header><span class="sec-title">⑥ 任务与文件</span></template>
      <el-form label-width="110px">
        <el-form-item label="任务名称" required>
          <el-select
            v-model="form.task_name"
            filterable
            allow-create
            default-first-option
            placeholder="请输入关键词，并选择创建好的可用任务"
            style="width:420px"
          >
            <el-option v-for="t in templates" :key="t.id" :label="t.task_name" :value="t.task_name" />
          </el-select>
          <el-button type="primary" plain style="margin-left:10px" @click="saveTemplate">创建任务</el-button>
        </el-form-item>
        <el-form-item label="上传模板">
          <el-link type="primary" :underline="false" :href="scrmApi.templateUrl(form.import_method)" target="_blank">
            <el-icon><Download /></el-icon>&nbsp;下载模板（{{ form.import_method === 'phone' ? '手机号' : '学员号' }}版）
          </el-link>
        </el-form-item>
        <el-form-item label="选择文件" required>
          <el-upload
            :auto-upload="false"
            :limit="1"
            accept=".xlsx,.xls"
            :on-change="onFileChange"
            :on-remove="() => file = null"
          >
            <el-button type="primary"><el-icon><UploadFilled /></el-icon>&nbsp;选择待导入 Excel</el-button>
            <template #tip>
              <div class="upload-tip">仅支持 .xlsx 文件，请严格使用模板填写，单次最多 5000 行。</div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" :loading="uploading" @click="doUpload">
            <el-icon><Promotion /></el-icon>&nbsp;开始导入
          </el-button>
          <el-button size="large" @click="goRecord">查看数据记录</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 导入结果 -->
    <el-dialog v-model="resultDlg" title="导入结果" width="520px">
      <div v-if="result" class="result-box">
        <el-result icon="success" title="导入处理完成" :sub-title="'流水号 ' + result.serial_no" />
        <el-descriptions :column="2" border>
          <el-descriptions-item label="总条数">{{ result.total_count }}</el-descriptions-item>
          <el-descriptions-item label="成功">{{ result.success_count }}</el-descriptions-item>
          <el-descriptions-item label="纯新线索">{{ result.new_count }}</el-descriptions-item>
          <el-descriptions-item label="激活线索">{{ result.active_count }}</el-descriptions-item>
          <el-descriptions-item label="失败条数">
            <el-link v-if="result.fail_count > 0" type="danger" :href="scrmApi.failFileUrl(result.id)" target="_blank">
              {{ result.fail_count }}（下载失败明细）
            </el-link>
            <span v-else>0</span>
          </el-descriptions-item>
          <el-descriptions-item label="重复线索">{{ result.duplicate_count }}</el-descriptions-item>
          <el-descriptions-item label="无效数据">{{ result.invalid_count }}</el-descriptions-item>
          <el-descriptions-item label="去重规则去重">{{ result.dedup_hit_count }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="resultDlg = false">关闭</el-button>
        <el-button type="primary" @click="goRecord">前往数据记录</el-button>
      </template>
    </el-dialog>

    <!-- 渠道规则弹窗 -->
    <el-dialog v-model="ruleDlg" title="渠道填写规则" width="560px">
      <ol class="rule-list">
        <li>渠道分为一级 / 二级 / 三级，请逐级选择，未选到三级时按最末级落库。</li>
        <li>导入后系统会按手机号 / 学员号回溯历史线索渠道，若与本次渠道不一致，以本次填写渠道为准。</li>
        <li>「市场地推-学校地推」适用于线下驻点、班级宣讲、家长会等场景。</li>
        <li>「大数据图文」适用于信息流投放、短信触达等付费渠道。</li>
        <li>渠道数据用于业绩归因与渠道 ROI 分析，请谨慎选择。</li>
      </ol>
    </el-dialog>

    <!-- 关联活动弹窗 -->
    <el-dialog v-model="actDlg" title="选择已有活动" width="720px">
      <div style="margin-bottom:10px">
        <el-input v-model="actKw" placeholder="搜索活动名称" clearable style="width:260px" @input="searchActivities(actKw)">
          <template #append><el-button @click="searchActivities(actKw)">搜索</el-button></template>
        </el-input>
      </div>
      <el-table :data="activities" height="360" highlight-current-row @row-dblclick="pickActivity">
        <el-table-column prop="name" label="活动名称" min-width="200" />
        <el-table-column prop="category" label="类别" width="90" />
        <el-table-column prop="type" label="类型" width="80" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="['info','success','warning'][row.status]">
              {{ ['未开始', '进行中', '已结束'][row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button link type="primary" @click="pickActivity(row)">选择</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { UploadFile } from 'element-plus'
import { Connection, Download, UploadFilled, Promotion } from '@element-plus/icons-vue'
import { scrmApi, type Activity } from '../api/scrm'

const route = useRoute()
const router = useRouter()

const dict = ref<Record<string, any[]>>({})
const channelOptions = ref<any[]>([])
const activities = ref<Activity[]>([])
const templates = ref<any[]>([])

const channelValue = ref<string[]>([])
const file = ref<File | null>(null)
const uploading = ref(false)
const resultDlg = ref(false)
const result = ref<any>(null)
const ruleDlg = ref(false)
const actDlg = ref(false)
const actKw = ref('')

const form = reactive({
  import_method: 'phone',
  data_source: '市场地推',
  activity_id: 0 as number,
  activity_name: '',
  activity_category: '',
  activity_type: '',
  intention_dept: '',
  dedup: {
    enable: true,
    followed_checked: false,
    followed_dur_days: 30,
    old_student_checked: true,
    old_student: {
      status: '在读', fresh_type: '新生', project: '', in_year: '是',
      loss_enabled: false, loss_days: 180
    },
    opportunity_dedup: false,
    task_dedup: false
  },
  as_opportunity: false,
  submit_type: 'supervisor',
  task_name: ''
})

async function loadDicts() {
  const kinds = ['data_source', 'activity_category', 'activity_type', 'intention_dept',
    'intention_project', 'in_school_status', 'fresh_type']
  await Promise.all(kinds.map(async k => {
    dict.value[k] = await scrmApi.dict(k)
  }))
  channelOptions.value = await scrmApi.channels()
}

async function searchActivities(kw: string) {
  activities.value = await scrmApi.allActivities(kw || '')
}
async function loadTemplates() {
  templates.value = await scrmApi.taskTemplates('')
}

function onChannelChange(v: string[]) {
  // 三级渠道值已在级联选择器中，提交时使用
}

function onActivityChange(id: number) {
  const a = activities.value.find(x => x.id === id)
  if (a) {
    form.activity_name = a.name
    form.activity_category = a.category
    form.activity_type = a.type
    form.intention_dept = a.intention_dept
  }
}
function pickActivity(a: Activity) {
  form.activity_id = a.id
  onActivityChange(a.id)
  actDlg.value = false
}

async function saveTemplate() {
  if (!form.task_name.trim()) {
    ElMessage.warning('请先填写任务名称')
    return
  }
  await scrmApi.createTaskTemplate(buildPayload())
  ElMessage.success('任务已保存，可在后续导入中复用')
  await loadTemplates()
}

function buildPayload() {
  const [l1 = '', l2 = '', l3 = ''] = channelValue.value
  return {
    task_name: form.task_name,
    activity_id: form.activity_id,
    activity_name: form.activity_name,
    import_method: form.import_method,
    data_source: form.data_source,
    channel_l1: l1, channel_l2: l2, channel_l3: l3,
    submit_type: form.submit_type,
    as_opportunity: form.as_opportunity,
    dedup: form.dedup
  }
}

function onFileChange(f: UploadFile) {
  file.value = (f.raw as File) || null
}

async function doUpload() {
  if (!form.data_source) return ElMessage.warning('请选择数据来源')
  if (channelValue.value.length === 0) return ElMessage.warning('请选择市场渠道')
  if (!form.activity_id) return ElMessage.warning('请关联活动')
  if (!form.task_name.trim()) return ElMessage.warning('请选择或输入任务名称')
  if (!file.value) return ElMessage.warning('请选择要导入的 Excel 文件')

  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file.value)
    const payload = buildPayload()
    Object.entries(payload).forEach(([k, v]) => {
      fd.append(k, typeof v === 'object' ? JSON.stringify(v) : String(v))
    })
    result.value = await scrmApi.uploadImport(fd)
    resultDlg.value = true
  } catch (e: any) {
    ElMessage.error(e.message || '导入失败')
  } finally {
    uploading.value = false
  }
}

function goRecord() {
  router.push('/data-record')
}

onMounted(async () => {
  await Promise.all([loadDicts(), searchActivities(''), loadTemplates()])
  // 从活动管理页「导入数据」跳转过来时预填活动
  const aid = Number(route.query.activity_id)
  if (aid) {
    form.activity_id = aid
    form.activity_name = (route.query.name as string) || ''
  }
})
</script>

<style scoped>
.chip-bar { margin-bottom: 10px; }
.sec { margin-bottom: 12px; }
.sec-title { font-weight: 700; font-size: 14px; }
.sec-head-row { display: flex; align-items: center; justify-content: space-between; padding-right: 12px; }
.dedup-panel { display: flex; flex-direction: column; gap: 10px; padding-left: 4px; }
.dedup-sub {
  display: flex; align-items: center; gap: 10px; flex-wrap: wrap;
  background: #f7f9fb; border-radius: 6px; padding: 10px 12px; margin: 4px 0 4 24px;
}
.dedup-tip { color: #909399; font-size: 12px; margin: -4px 0 0 24px; }
.submit-radio { margin-bottom: 14px; }
.warn-box { max-width: 860px; line-height: 1.7; }
.upload-tip { color: #909399; font-size: 12px; margin-top: 4px; }
.result-box :deep(.el-result__icon) { margin-bottom: 0; }
.rule-list { line-height: 2; padding-left: 20px; color: #475669; }
</style>
