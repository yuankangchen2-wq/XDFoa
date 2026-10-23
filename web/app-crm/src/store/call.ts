import { reactive } from 'vue'
import type { Lead } from '../api/scrm'

// 云呼叫工作条全局状态：任意页面都可带线索上下文唤起
export const callState = reactive({
  visible: false,
  lead: null as Lead | null
})

export function openCallBar(lead?: Lead | null) {
  callState.lead = lead || null
  callState.visible = true
}

export function closeCallBar() {
  callState.visible = false
}
