<template>
  <div class="portal-root">
    <TopNav />

    <!-- 门户首页 -->
    <main v-show="isHome">
      <router-view />
    </main>

    <!-- qiankun 微前端子应用容器（始终保留节点） -->
    <div v-show="!isHome" class="subapp-wrap">
      <div class="subapp-bar">
        <el-button text @click="$router.push('/')">
          <el-icon style="margin-right:4px"><ArrowLeftBold /></el-icon>
          返回门户首页
        </el-button>
      </div>
      <div id="subapp-container"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import TopNav from './components/TopNav.vue'
import { startQiankun } from './qiankun'

const route = useRoute()
const isHome = computed(() => route.path === '/')

onMounted(() => {
  startQiankun()
})
</script>

<style>
/* ===== 全局绿色主题（仅主应用，子应用样式沙箱隔离） ===== */
:root {
  --el-color-primary: #00b578;
  --el-color-primary-light-3: #4dc99c;
  --el-color-primary-light-5: #80d7b8;
  --el-color-primary-light-7: #b3e6d4;
  --el-color-primary-light-8: #cceee3;
  --el-color-primary-light-9: #e6f7f1;
  --el-color-primary-dark-2: #009e68;
}

html, body, #app {
  margin: 0;
  padding: 0;
  min-height: 100%;
  background: #f2f4f7;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC',
    'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
}

* { box-sizing: border-box; }

.portal-root { min-height: 100vh; }

.subapp-wrap { min-height: calc(100vh - 60px); }
.subapp-bar {
  max-width: 100%;
  padding: 6px 24px;
  background: #fff;
  border-bottom: 1px solid #eef0f3;
}
#subapp-container { min-height: calc(100vh - 100px); }
</style>
