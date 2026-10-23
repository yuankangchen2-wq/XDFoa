<template>
  <div class="page">
    <div class="head">
      <div class="title">业绩申诉</div>
      <el-button type="primary" size="small" @click="openDialog">申诉</el-button>
    </div>

    <el-tabs v-model="bizLine" @tab-change="reload">
      <el-tab-pane label="班课（听课证号）" name="班课" />
      <el-tab-pane label="高端（合同编号）" name="高端" />
      <el-tab-pane label="会员卡（交易编码）" name="会员卡" />
      <el-tab-pane label="课时包（交易编号）" name="课时包" />
    </el-tabs>

    <div class="toolbar">
      <el-input v-model="kw" placeholder="商机编码/编号" size="small" style="width:220px" clearable
        @keyup.enter="reload" @clear="reload" />
      <el-select v-model="st" placeholder="审批状态" clearable size="small" style="width:130px" @change="reload">
        <el-option label="待审批" value="待审批" />
        <el-option label="已通过" value="已通过" />
        <el-option label="已驳回" value="已驳回" />
      </el-select>
      <el-button type="primary" size="small" @click="reload">查询</el-button>
    </div>

    <el-table :data="rows" border size="small" v-loading="loading">
      <el-table-column prop="opp_code" label="商机编码" width="160" />
      <el-table-column :label="certLabel" prop="cert_no" width="150" />
      <el-table-column prop="owner_name" label="当前归属人" width="110" />
      <el-table-column prop="owner_dept" label="归属部门" min-width="150" />
      <el-table-column prop="reason" label="申诉原因" min-width="220" show-overflow-tooltip />
      <el-table-column prop="adjust_mode" label="调整方式" width="100" />
      <el-table-column prop="applicant" label="申诉人" width="100" />
      <el-table-column prop="status" label="审批状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === '已通过' ? 'success' : row.status === '已驳回' ? 'danger' : 'warning'" size="small">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="申诉时间" width="160">
        <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>
    <el-pagination class="pager" background layout="total,prev,pager,next" :total="total"
      v-model:current-page="page" @current-change="reload" />

    <el-dialog v-model="dialog" title="发起业绩申诉" width="500px">
      <el-form :model="form" label-width="100px" size="small">
        <el-form-item label="业务类型">
          <el-input :model-value="bizLine" disabled />
        </el-form-item>
        <el-form-item label="选择商机" required>
          <el-select v-model="form.opp_code" filterable placeholder="选择需申诉的商机" style="width:100%">
            <el-option v-for="o in oppOptions" :key="o.id"
              :label="`${o.code}｜${o.name}｜${o.owner_name}`" :value="o.code" />
          </el-select>
        </el-form-item>
        <el-form-item :label="certLabel">
          <el-input v-model="form.cert_no" :placeholder="`请输入${certLabel}`" />
        </el-form-item>
        <el-form-item label="调整方式">
          <el-radio-group v-model="form.adjust_mode">
            <el-radio value="调整归属">调整归属</el-radio>
            <el-radio value="业绩拆分">业绩拆分</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="申诉原因" required>
          <el-input v-model="form.reason" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="凭证图片">
          <el-upload action="#" list-type="picture-card" :auto-upload="false" :limit="3"
            :on-change="onFile" accept="image/*">
            <span>+ 上传</span>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog = false">取消</el-button>
        <el-button type="primary" @click="submit">提交申诉</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { scrmApi, fmtTime, type Opportunity } from '../api/scrm'

const bizLine = ref('班课')
const certLabel = computed(() => ({ 班课: '听课证号', 高端: '合同编号', 会员卡: '交易编码', 课时包: '交易编号' })[bizLine.value] || '')
const rows = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const kw = ref('')
const st = ref('')

async function reload() {
  loading.value = true
  try {
    const res = await scrmApi.appeals({ biz_line: bizLine.value, keyword: kw.value, status: st.value, page: page.value, page_size: 20 })
    rows.value = res.items; total.value = res.total
  } finally { loading.value = false }
}

const dialog = ref(false)
const form = reactive<any>({ opp_code: '', cert_no: '', adjust_mode: '调整归属', reason: '' })
const oppOptions = ref<Opportunity[]>([])
async function openDialog() {
  Object.assign(form, { opp_code: '', cert_no: '', adjust_mode: '调整归属', reason: '' })
  const res = await scrmApi.opportunities({ scope: 'dept', page: 1, page_size: 100 })
  oppOptions.value = res.items
  dialog.value = true
}
function onFile(_f: any) { /* 演示环境仅记录文件名 */ }
async function submit() {
  if (!form.opp_code) return ElMessage.warning('请选择商机')
  if (!form.reason) return ElMessage.warning('请填写申诉原因')
  await scrmApi.createAppeal({ ...form, biz_line: bizLine.value })
  ElMessage.success('申诉已提交，等待审核中心审批')
  dialog.value = false
  reload()
}

onMounted(reload)
</script>

<style scoped>
.page { padding: 12px 16px; }
.head { display: flex; justify-content: space-between; align-items: center; }
.title { font-size: 16px; font-weight: 600; }
.toolbar { display: flex; gap: 8px; margin: 8px 0; }
.pager { margin-top: 10px; justify-content: flex-end; }
</style>
