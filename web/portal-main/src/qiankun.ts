import { registerMicroApps, start } from 'qiankun'

export const microApps = [
  {
    name: 'app-portal',
    entry: '//localhost:5177',
    container: '#subapp-container',
    activeRule: '/portal'
  },
  {
    name: 'app-iam',
    entry: '//localhost:5175',
    container: '#subapp-container',
    activeRule: '/iam'
  },
  {
    name: 'app-mdm',
    entry: '//localhost:5174',
    container: '#subapp-container',
    activeRule: '/mdm'
  },
  {
    name: 'app-erp',
    entry: '//localhost:5176',
    container: '#subapp-container',
    activeRule: '/erp'
  },
  {
    name: 'app-crm',
    entry: '//localhost:5178',
    container: '#subapp-container',
    activeRule: '/crm'
  }
]

export function startQiankun() {
  registerMicroApps(microApps, {
    beforeLoad: [app => {
      console.log('[qiankun] before load', app.name)
      return Promise.resolve()
    }],
    beforeMount: [app => {
      console.log('[qiankun] before mount', app.name)
      return Promise.resolve()
    }],
    afterUnmount: [app => {
      console.log('[qiankun] after unmount', app.name)
      return Promise.resolve()
    }]
  })

  start({
    sandbox: { experimentalStyleIsolation: true },
    prefetch: false
  })
}
