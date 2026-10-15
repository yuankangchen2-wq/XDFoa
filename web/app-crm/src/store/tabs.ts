import { reactive } from 'vue'

export interface WorkTab {
  title: string
  path: string
  fullPath: string
  closable: boolean
}

export const tabStore = reactive({
  tabs: [
    { title: '首页', path: '/home', fullPath: '/home', closable: false }
  ] as WorkTab[],
  activeFullPath: '/home'
})

export function openTab(t: WorkTab) {
  const exists = tabStore.tabs.some(x => x.fullPath === t.fullPath)
  if (!exists) {
    tabStore.tabs.push(t)
  }
  tabStore.activeFullPath = t.fullPath
}

export function closeTab(fullPath: string): string | null {
  const idx = tabStore.tabs.findIndex(x => x.fullPath === fullPath)
  if (idx < 0) return null
  const removed = tabStore.tabs[idx]
  if (!removed.closable) return null
  tabStore.tabs.splice(idx, 1)
  if (tabStore.activeFullPath === fullPath) {
    const next = tabStore.tabs[Math.max(0, idx - 1)]
    tabStore.activeFullPath = next.fullPath
    return next.fullPath
  }
  return null
}

export function syncActive(fullPath: string, path: string, title: string, query?: Record<string, string>) {
  tabStore.activeFullPath = fullPath
  const exists = tabStore.tabs.some(x => x.fullPath === fullPath)
  if (!exists) {
    tabStore.tabs.push({ title, path, fullPath, closable: path !== '/home' })
  }
  if (query && Object.keys(query).length === 0) return
}
