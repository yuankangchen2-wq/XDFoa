<template>
  <el-dialog
    :model-value="modelValue"
    :title="title"
    width="440px"
    @update:model-value="$emit('update:modelValue', $event)"
    @open="loadEmployees"
  >
    <el-form label-width="80px">
      <el-form-item label="接收人" required>
        <el-select v-model="receiver" filterable placeholder="选择员工 / 主管" style="width:100%">
          <el-option
            v-for="e in employees"
            :key="e.id"
            :label="e.name + '（' + (e.role === 'supervisor' ? '主管' : '员工') + ' · ' + e.campus + '）'"
            :value="e.name"
          />
        </el-select>
      </el-form-item>
      <el-form-item v-if="showCount">
        <span class="tip">本次将交接 <b>{{ count }}</b> 条数据给所选接收人，并写入交接记录。</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="$emit('update:modelValue', false)">取消</el-button>
      <el-button type="primary" @click="confirm">确认交接</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { scrmApi, type Employee } from '../api/scrm'

const props = withDefaults(defineProps<{
  modelValue: boolean
  title?: string
  count?: number
  showCount?: boolean
}>(), {
  title: '线索交接',
  count: 0,
  showCount: true
})
const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'confirm', employee: Employee): void
}>()

const employees = ref<Employee[]>([])
const receiver = ref('')

async function loadEmployees() {
  if (employees.value.length === 0) {
    employees.value = await scrmApi.employees()
  }
}

watch(() => props.modelValue, v => {
  if (v) receiver.value = ''
})

function confirm() {
  if (!receiver.value) {
    ElMessage.warning('请选择接收人')
    return
  }
  const emp = employees.value.find(e => e.name === receiver.value)
  if (emp) emit('confirm', emp)
}
</script>

<style scoped>
.tip { color: #606f82; font-size: 13px; }
.tip b { color: #00a86b; }
</style>
