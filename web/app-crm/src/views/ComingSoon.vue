<template>
  <div class="coming">
    <el-result icon="info" :title="'「' + name + '」模块将在后续版本开放'" sub-title="第一期范围：首页、线索创建（线索导入 / 活动管理 / 待交接线索 / 全部线索 / 异常线索 / 数据记录）">
      <template #extra>
        <el-button type="primary" @click="$router.push('/home')">返回首页</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const name = computed(() => (route.query.name as string) || '该功能')
</script>

<style scoped>
.coming { background: #fff; border-radius: 8px; min-height: 400px; display: flex; align-items: center; justify-content: center; }
</style>
